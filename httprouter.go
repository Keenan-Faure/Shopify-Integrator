package main

import (
	"api"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"iocsv"
	"log"
	"net/http"
	"objects"
	"strconv"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

// GET /api/inventory
func (dbconfig *DbConfig) GetWarehouseLocations(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	shopify_locations, err := dbconfig.DB.GetShopifyLocations(r.Context(), database.GetShopifyLocationsParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	locations := []database.ShopifyLocation{}
	if len(shopify_locations) > 0 {
		locations = append(locations, shopify_locations...)
	}
	RespondWithJSON(w, http.StatusOK, locations)
}

// DELETE /api/inventory
func (dbconfig *DbConfig) RemoveWarehouseLocation(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	id := chi.URLParam(r, "id")
	err := IDValidation(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	delete_id, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+id)
		return
	}
	err = dbconfig.DB.RemoveShopifyLocationMap(r.Context(), delete_id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "Deleted",
	})
}

// POST /api/inventory
func (dbconfig *DbConfig) AddWarehouseLocationMap(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	location_map, err := DecodeInventoryMap(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	if InventoryMapValidation(location_map) != nil {
		RespondWithError(w, http.StatusBadRequest, "data validation error")
		return
	}
	result, err := dbconfig.DB.CreateShopifyLocation(r.Context(), database.CreateShopifyLocationParams{
		ID:                   uuid.New(),
		ShopifyWarehouseName: location_map.ShopifyWarehouseName,
		ShopifyLocationID:    location_map.LocationID,
		WarehouseName:        location_map.WarehouseName,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusCreated, result)
}

// GET /api/products/export
func (dbconfig *DbConfig) ExportProductsHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	product_ids, err := dbconfig.DB.GetProductIDs(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	products := []objects.Product{}
	for _, product_id := range product_ids {
		product, err := CompileProductData(dbconfig, product_id, r.Context(), false)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		products = append(products, product)
	}
	csv_data := [][]string{}
	headers := []string{}
	if len(products) > 0 {
		headers = iocsv.CSVProductHeaders(products[0])
	}
	// Adds (distinct) pricing and warehouses headers
	price_tiers, err := AddPricingHeaders(dbconfig, r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	headers = append(headers, price_tiers...)
	qty_warehouses, err := AddQtyHeaders(dbconfig, r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	headers = append(headers, qty_warehouses...)

	images_max, err := IOGetMax(dbconfig, r.Context(), "image")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pricing_max, err := IOGetMax(dbconfig, r.Context(), "price")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	qty_max, err := IOGetMax(dbconfig, r.Context(), "qty")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	image_headers := iocsv.GetProductImagesCSV(products[0].ProductImages, int(images_max), true)
	headers = append(headers, image_headers...)
	csv_data = append(csv_data, headers)
	for _, product := range products {
		for _, variant := range product.Variants {
			row := iocsv.CSVProductValuesByVariant(product, variant, pricing_max, qty_max)
			csv_data = append(csv_data, row)
		}
	}
	file_name, err := iocsv.WriteFile(csv_data, "")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: file_name,
	})
	// use javascript to return that file to be sent on the browser
}

// POST /api/products/import?file_name={{file}}
func (dbconfig *DbConfig) ProductImportHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	file_name := r.URL.Query().Get("file_name")
	csv_products, err := iocsv.ReadFile(file_name)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	processed_counter := 0
	failure_counter := 0
	products_added := 0
	products_updated := 0
	variants_updated := 0
	variants_added := 0
	for _, csv_product := range csv_products {
		product_exists := false
		product, err := dbconfig.DB.CreateProduct(r.Context(), database.CreateProductParams{
			ID:          uuid.New(),
			ProductCode: csv_product.ProductCode,
			Active:      "1",
			Title:       utils.ConvertStringToSQL(csv_product.Title),
			BodyHtml:    utils.ConvertStringToSQL(csv_product.BodyHTML),
			Category:    utils.ConvertStringToSQL(csv_product.Category),
			Vendor:      utils.ConvertStringToSQL(csv_product.Vendor),
			ProductType: utils.ConvertStringToSQL(csv_product.ProductType),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})
		if err != nil {
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				product_exists = true
				err := dbconfig.DB.UpdateProduct(r.Context(), database.UpdateProductParams{
					Active:      "1",
					ProductCode: csv_product.ProductCode,
					Title:       utils.ConvertStringToSQL(csv_product.Title),
					BodyHtml:    utils.ConvertStringToSQL(csv_product.BodyHTML),
					Category:    utils.ConvertStringToSQL(csv_product.Category),
					Vendor:      utils.ConvertStringToSQL(csv_product.Vendor),
					ProductType: utils.ConvertStringToSQL(csv_product.ProductType),
					UpdatedAt:   time.Now().UTC(),
				})
				if err != nil {
					log.Println(err)
					failure_counter++
					continue
				}
				products_updated++
			} else {
				log.Println(err)
				failure_counter++
				continue
			}
		}
		if !product_exists {
			products_added++
		}
		if !product_exists {
			option_names := CreateOptionNamesMap(csv_product)
			for key, option_name := range option_names {
				if option_name != "" {
					_, err = dbconfig.DB.CreateProductOption(r.Context(), database.CreateProductOptionParams{
						ID:        uuid.New(),
						ProductID: product.ID,
						Name:      option_name,
						Position:  int32(key + 1),
					})
					if err != nil {
						log.Println(err)
						failure_counter++
						continue
					}
				}
			}
		}
		if product.ID == uuid.Nil {
			product.ID, err = dbconfig.DB.GetProductIDByCode(r.Context(), csv_product.ProductCode)
			if err != nil {
				log.Println(err)
				failure_counter++
				continue
			}
		}
		// Update product options
		option_names := CreateOptionNamesMap(csv_product)
		for key, option_name := range option_names {
			if option_name != "" {
				_, err = dbconfig.DB.UpdateProductOption(r.Context(), database.UpdateProductOptionParams{
					Name:       option_name,
					Position:   int32(key + 1),
					ProductID:  product.ID,
					Position_2: int32(key + 1),
				})
				if err != nil {
					log.Println(err)
					failure_counter++
					continue
				}
			}
		}
		// add images to product
		// overwrite ones with the same position
		images := CreateImageMap(csv_product)
		for key := range images {
			if images[key] != "" {
				err = AddImagery(dbconfig, product.ID, images[key], key+1)
				if err != nil {
					log.Println(err)
					failure_counter++
					continue
				}
			}
		}
		// create variant
		variant, err := dbconfig.DB.CreateVariant(r.Context(), database.CreateVariantParams{
			ID:        uuid.New(),
			ProductID: product.ID,
			Sku:       csv_product.SKU,
			Option1:   utils.ConvertStringToSQL(csv_product.Option1Value),
			Option2:   utils.ConvertStringToSQL(csv_product.Option2Value),
			Option3:   utils.ConvertStringToSQL(csv_product.Option3Value),
			Barcode:   utils.ConvertStringToSQL(csv_product.Barcode),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			// if it already exists, then we update it
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				err := dbconfig.DB.UpdateVariant(r.Context(), database.UpdateVariantParams{
					Option1:   utils.ConvertStringToSQL(csv_product.Option1Value),
					Option2:   utils.ConvertStringToSQL(csv_product.Option2Value),
					Option3:   utils.ConvertStringToSQL(csv_product.Option3Value),
					Barcode:   utils.ConvertStringToSQL(csv_product.Barcode),
					UpdatedAt: time.Now().UTC(),
					Sku:       csv_product.SKU,
				})
				if err != nil {
					log.Println(err)
					failure_counter++
					continue
				}
				variants_updated++
				continue
			}
			log.Println(err)
			failure_counter++
			continue
		}
		variants_added++
		if variant.ID == uuid.Nil {
			variant.ID, err = dbconfig.DB.GetVariantIDByCode(r.Context(), csv_product.SKU)
			if err != nil {
				log.Println(err)
				failure_counter++
				continue
			}
		}
		for _, pricing_value := range csv_product.Pricing {
			err = AddPricing(dbconfig, csv_product.SKU, variant.ID, pricing_value.Name, pricing_value.Value)
			if err != nil {
				log.Println(err)
				failure_counter++
				continue
			}
		}
		for _, qty_value := range csv_product.Warehouses {
			err = AddWarehouse(dbconfig, variant.Sku, variant.ID, qty_value.Name, qty_value.Value)
			if err != nil {
				log.Println(err)
				failure_counter++
				continue
			}
		}
		if err != nil {
			log.Println(err)
			failure_counter++
			continue
		}
		processed_counter++
	}
	err = iocsv.RemoveFile(file_name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ImportResponse{
		ProcessedCounter: processed_counter,
		FailCounter:      failure_counter,
		ProductsAdded:    products_added,
		ProductsUpdated:  products_updated,
		VariantsAdded:    variants_added,
		VariantsUpdated:  variants_updated,
	})
}

// POST /api/customers/
func (dbconfig *DbConfig) PostCustomerHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	customer_body, err := DecodeCustomerRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	if CustomerValidation(customer_body) != nil {
		RespondWithError(w, http.StatusBadRequest, "data validation error")
		return
	}
	customer, err := dbconfig.DB.CreateCustomer(r.Context(), database.CreateCustomerParams{
		ID:        uuid.New(),
		FirstName: customer_body.FirstName,
		LastName:  customer_body.LastName,
		Email:     utils.ConvertStringToSQL(customer_body.Email),
		Phone:     utils.ConvertStringToSQL(customer_body.Phone),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	for key := range customer_body.Address {
		_, err := dbconfig.DB.CreateAddress(r.Context(), database.CreateAddressParams{
			ID:         uuid.New(),
			CustomerID: customer.ID,
			Name:       utils.ConvertStringToSQL("default"),
			FirstName:  customer_body.Address[key].FirstName,
			LastName:   customer_body.Address[key].LastName,
			Address1:   utils.ConvertStringToSQL(customer_body.Address[key].Address1),
			Address2:   utils.ConvertStringToSQL(customer_body.Address[key].Address2),
			Suburb:     utils.ConvertStringToSQL(""),
			City:       utils.ConvertStringToSQL(customer_body.Address[key].City),
			Province:   utils.ConvertStringToSQL(customer_body.Address[key].Province),
			PostalCode: utils.ConvertStringToSQL(customer_body.Address[key].PostalCode),
			Company:    utils.ConvertStringToSQL(customer_body.Address[key].Company),
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	customer_data, err := CompileCustomerData(dbconfig, customer.ID, r.Context(), false)
	if err != nil {
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	RespondWithJSON(w, http.StatusCreated, customer_data)
}

// POST /api/orders?token={{token}}&api_key={{key}}
// ngrok exposed url
func (dbconfig *DbConfig) PostOrderHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	web_token := r.URL.Query().Get("token")
	if TokenValidation(web_token) != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}
	_, err := dbconfig.DB.ValidateWebhookByUser(r.Context(), database.ValidateWebhookByUserParams{
		WebhookToken: web_token,
		ApiKey:       dbUser.ApiKey,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	order_body, err := DecodeOrderRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(order_body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	db_order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), utils.ConvertStringToSQL(fmt.Sprint(order_body.OrderNumber)))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return
		}
	}
	if db_order.WebCode.String == fmt.Sprint(order_body.OrderNumber) {
		response_payload, err := dbconfig.QueueHelper(objects.RequestQueueHelper{
			Type:        "order",
			Status:      "in-queue",
			Instruction: "update_order",
			Endpoint:    "queue",
			ApiKey:      dbUser.ApiKey,
			Method:      http.MethodPost,
			Object:      order_body,
		}, nil)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
			return
		}
		RespondWithJSON(w, http.StatusOK, response_payload)
	} else {
		response_payload, err := dbconfig.QueueHelper(objects.RequestQueueHelper{
			Type:        "order",
			Status:      "in-queue",
			Instruction: "add_order",
			Endpoint:    "queue",
			ApiKey:      dbUser.ApiKey,
			Method:      http.MethodPost,
			Object:      order_body,
		}, &buffer)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
			return
		}
		RespondWithJSON(w, http.StatusCreated, response_payload)
	}
}

// POST /api/products
func (dbconfig *DbConfig) PostProductHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	params, err := DecodeProductRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	validation := ProductValidation(params)
	if validation != nil {
		RespondWithError(w, http.StatusBadRequest, validation.Error())
		return
	}
	err = ValidateDuplicateOption(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	err = ValidateDuplicateSKU(params, dbconfig, r)
	if err != nil {
		RespondWithError(w, http.StatusConflict, utils.ConfirmError(err))
		return
	}
	err = DuplicateOptionValues(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	csv_products := ConvertProductToCSV(params)
	for _, csv_product := range csv_products {
		err = ProductValidationDatabase(csv_product, dbconfig, r)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	// add product to database
	product, err := dbconfig.DB.CreateProduct(r.Context(), database.CreateProductParams{
		ID:          uuid.New(),
		Active:      "1",
		ProductCode: params.ProductCode,
		Title:       utils.ConvertStringToSQL(params.Title),
		BodyHtml:    utils.ConvertStringToSQL(params.BodyHTML),
		Category:    utils.ConvertStringToSQL(params.Category),
		Vendor:      utils.ConvertStringToSQL(params.Vendor),
		ProductType: utils.ConvertStringToSQL(params.ProductType),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		log.Println("1: " + err.Error())
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	for key := range params.ProductOptions {
		_, err := dbconfig.DB.CreateProductOption(r.Context(), database.CreateProductOptionParams{
			ID:        uuid.New(),
			ProductID: product.ID,
			Name:      params.ProductOptions[key].Value,
			Position:  int32(key + 1),
		})
		if err != nil {
			log.Println("2: " + err.Error())
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	// add variants
	for key := range params.Variants {
		variant, err := dbconfig.DB.CreateVariant(r.Context(), database.CreateVariantParams{
			ID:        uuid.New(),
			ProductID: product.ID,
			Sku:       params.Variants[key].Sku,
			Option1:   utils.ConvertStringToSQL(params.Variants[key].Option1),
			Option2:   utils.ConvertStringToSQL(params.Variants[key].Option2),
			Option3:   utils.ConvertStringToSQL(params.Variants[key].Option3),
			Barcode:   utils.ConvertStringToSQL(params.Variants[key].Barcode),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Println("3: " + err.Error())
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		// variant pricing & variant qty
		for key_price := range params.Variants[key].VariantPricing {
			_, err := dbconfig.DB.CreateVariantPricing(r.Context(), database.CreateVariantPricingParams{
				ID:        uuid.New(),
				VariantID: variant.ID,
				Name:      params.Variants[key].VariantPricing[key_price].Name,
				Value:     utils.ConvertStringToSQL(params.Variants[key].VariantPricing[key_price].Value),
				Isdefault: params.Variants[key].VariantPricing[key_price].IsDefault,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println("4: " + err.Error())
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
			_, err = dbconfig.DB.CreateVariantQty(r.Context(), database.CreateVariantQtyParams{
				ID:        uuid.New(),
				VariantID: variant.ID,
				Name:      params.Variants[key].VariantQuantity[key_price].Name,
				Isdefault: params.Variants[key].VariantQuantity[key_price].IsDefault,
				Value:     utils.ConvertIntToSQL(params.Variants[key].VariantQuantity[key_price].Value),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				log.Println("5: " + err.Error())
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
		if err != nil {
			log.Println("6: " + err.Error())
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	// TODO is it necessary to respond with the created product data
	product_added, err := CompileProductData(dbconfig, product.ID, r.Context(), false)
	if err != nil {
		log.Println("7: " + err.Error())
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusCreated, product_added)
}

// GET /api/customers/search?q=value
func (dbconfig *DbConfig) CustomerSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	customers_by_name, err := dbconfig.DB.GetCustomersByName(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
	}
	RespondWithJSON(w, http.StatusOK, customers_by_name)
}

// GET /api/customers/{id}
func (dbconfig *DbConfig) CustomerHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	customer_id := chi.URLParam(r, "id")
	err := IDValidation(customer_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	customer_uuid, err := uuid.Parse(customer_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+customer_id)
		return
	}
	customer, err := CompileCustomerData(dbconfig, customer_uuid, r.Context(), false)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, customer)
}

// GET /api/customers?page=1
func (dbconfig *DbConfig) CustomersHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	dbCustomers, err := dbconfig.DB.GetCustomers(r.Context(), database.GetCustomersParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	customers := []objects.Customer{}
	for _, value := range dbCustomers {
		cust, err := CompileCustomerData(dbconfig, value.ID, r.Context(), true)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		customers = append(customers, cust)
	}
	RespondWithJSON(w, http.StatusOK, customers)
}

// GET /api/orders/search?q=value
func (dbconfig *DbConfig) OrderSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	customer_orders, err := dbconfig.DB.GetOrdersSearchByCustomer(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
	}
	webcode_orders, err := dbconfig.DB.GetOrdersSearchWebCode(r.Context(), search_query)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
	}
	RespondWithJSON(w, http.StatusOK, CompileOrderSearchResult(customer_orders, webcode_orders))
}

// GET /api/orders/{id}
func (dbconfig *DbConfig) OrderHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	order_id := chi.URLParam(r, "id")
	err := IDValidation(order_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	order_uuid, err := uuid.Parse(order_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+order_id)
		return
	}
	order_data, err := CompileOrderData(dbconfig, order_uuid, r.Context(), false)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, order_data)
}

// GET /api/orders?page=1
func (dbconfig *DbConfig) OrdersHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	dbOrders, err := dbconfig.DB.GetOrders(r.Context(), database.GetOrdersParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	orders := []objects.Order{}
	for _, value := range dbOrders {
		ord, err := CompileOrderData(dbconfig, value.ID, r.Context(), true)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		orders = append(orders, ord)
	}
	RespondWithJSON(w, http.StatusOK, orders)
}

// GET /api/products/filter?key=value&page=1
func (dbconfig *DbConfig) ProductFilterHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	query_param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
	query_param_category := utils.ConfirmFilters(r.URL.Query().Get("category"))
	query_param_vendor := utils.ConfirmFilters(r.URL.Query().Get("vendor"))
	response, err := CompileFilterSearch(dbconfig, r.Context(), page, query_param_type, query_param_category, query_param_vendor)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// GET /api/products/search?q=value
func (dbconfig *DbConfig) ProductSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query == "" || len(search_query) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid search param")
		return
	}
	sku_search, err := dbconfig.DB.GetProductsSearchSKU(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	title_search, err := dbconfig.DB.GetProductsSearchTitle(r.Context(), search_query)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, CompileSearchResult(sku_search, title_search))
}

// GET /api/products/{id}
func (dbconfig *DbConfig) ProductHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	product_id := chi.URLParam(r, "id")
	err := IDValidation(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	product_uuid, err := uuid.Parse(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+product_id)
		return
	}
	product_data, err := CompileProductData(dbconfig, product_uuid, r.Context(), false)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, product_data)
}

// GET /api/products?page=1
func (dbconfig *DbConfig) ProductsHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	dbProducts, err := dbconfig.DB.GetProducts(r.Context(), database.GetProductsParams{
		Limit:  10,
		Offset: int32((page - 1) * 10),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	products := []objects.Product{}
	for _, value := range dbProducts {
		prod, err := CompileProductData(dbconfig, value.ID, r.Context(), false)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		products = append(products, prod)
	}
	RespondWithJSON(w, http.StatusOK, products)
}

// POST /api/login
func (dbconfig *DbConfig) LoginHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/preregister
func (dbconfig *DbConfig) PreRegisterHandle(w http.ResponseWriter, r *http.Request) {
	email := utils.LoadEnv("email")
	email_psw := utils.LoadEnv("email_psw")
	if email == "" || email_psw == "" {
		RespondWithError(w, http.StatusInternalServerError, "invalid email or email password")
		return
	}
	request_body, err := DecodePreRegisterRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	err = PreRegisterValidation(request_body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	token_value := uuid.UUID{}
	token_value, exists, err := dbconfig.CheckTokenExists(request_body, r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	if !exists {
		token, err := dbconfig.DB.CreateToken(r.Context(), database.CreateTokenParams{
			Token:     uuid.New(),
			ID:        uuid.New(),
			Name:      request_body.Name,
			Email:     request_body.Email,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		token_value = token.Token
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	err = SendEmail(token_value, request_body.Email, request_body.Name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
		Message: "email sent",
	})
}

// POST /api/register
func (dbconfig *DbConfig) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeUserRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	err = ValidateTokenValidation(body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	token, err := dbconfig.DB.GetTokenValidation(r.Context(), body.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	request_token, err := uuid.Parse(body.Token)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode feed_id: "+body.Token)
		return
	}
	if token.Token != request_token {
		RespondWithError(w, http.StatusNotFound, "invalid token for user")
		return
	}
	if UserValidation(body) != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	exists, err := dbconfig.CheckUserExist(body.Name, r)
	if exists {
		RespondWithError(w, http.StatusConflict, utils.ConfirmError(err))
		return
	}
	user, err := dbconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      body.Name,
		Email:     body.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusCreated, user)
}

// GET /api/endpoints
func (dbconfig *DbConfig) EndpointsHandle(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, api.Endpoints())
}

// puppy function for a handler
//func Handle{{Name}} {
// 1.) decode params
// 2.) validate requestBody
// 3.) Check for duplicates in db
// 4.) (fi) create/insert new record
// 4.) (if) return record
// 5.) (else) return error
//}

// GET /api/ready
func (dbconfig *DbConfig) ReadyHandle(w http.ResponseWriter, r *http.Request) {
	if dbconfig.Valid {
		RespondWithJSON(w, 200, objects.ResponseString{
			Message: "OK",
		})
	} else {
		RespondWithJSON(w, 503, objects.ResponseString{
			Message: "Error",
		})
	}
}

// JSON helper functions
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	if payload == nil {
		response, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		w.Write(response)
		return nil
	} else {
		response, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		w.Write(response)
		return nil
	}
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}

// Middleware that determines which headers, http methods and orgins are allowed
func MiddleWare() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
}
