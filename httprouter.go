package main

import (
	"api"
	"database/sql"
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

// POST /api/products/import?file_name={{file}}
func (dbconfig *DbConfig) ProductImport(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	file_name := r.URL.Query().Get("file_name")
	csv_products, err := iocsv.ReadFile(file_name)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	processed_counter := 0
	failure_counter := 0
	skip_counter := 0
	products_added := 0
	products_updated := 0
	variants_updated := 0
	variants_added := 0
	for _, csv_product := range csv_products {
		product_exists := false
		// err := ProductValidationDatabase(csv_product, dbconfig, r)
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
			fmt.Println("1: " + err.Error())
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				product_exists = true
				// update product
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
					fmt.Println("2: " + err.Error())
					failure_counter++
					continue
				}
				products_updated++
			} else {
				// TODO log messages to console?
				fmt.Println("3: " + err.Error())
				failure_counter++
				continue
			}
		}
		if !product_exists {
			products_added++
		}
		if !product_exists {
			option_names := CreateOptionNamesMap(csv_product)
			for _, option_name := range option_names {
				if option_name != "" {
					_, err = dbconfig.DB.CreateProductOption(r.Context(), database.CreateProductOptionParams{
						ID:        uuid.New(),
						ProductID: product.ID,
						Name:      option_name,
					})
					if err != nil {
						fmt.Println("4: " + err.Error())
						failure_counter++
						continue
					}
				}
			}
		}
		if product.ID == uuid.Nil {
			product.ID, err = dbconfig.DB.GetProductIDByCode(r.Context(), csv_product.ProductCode)
			if err != nil {
				fmt.Println("4.5: " + err.Error())
				failure_counter++
				continue
			}
		}
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
			fmt.Println("5: " + err.Error())
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
					fmt.Println("6: " + err.Error())
					failure_counter++
					continue
				}
				err = dbconfig.DB.UpdateVariantPricing(r.Context(), database.UpdateVariantPricingParams{
					Name:  csv_product.PriceName,
					Value: utils.ConvertStringToSQL(csv_product.PriceValue),
					Sku:   csv_product.SKU,
				})
				if err != nil {
					fmt.Println("7: " + err.Error())
					failure_counter++
					continue
				}
				err = dbconfig.DB.UpdateVariantQty(r.Context(), database.UpdateVariantQtyParams{
					Name:  csv_product.QtyName,
					Value: utils.ConvertIntToSQL(csv_product.QtyValue),
					Sku:   csv_product.SKU,
				})
				if err != nil {
					fmt.Println("8: " + err.Error())
					failure_counter++
					continue
				}
				variants_updated++
				continue
			}
			fmt.Println("9: " + err.Error())
			failure_counter++
			continue
		}
		variants_added++
		if variant.ID == uuid.Nil {
			variant.ID, err = dbconfig.DB.GetVariantIDByCode(r.Context(), csv_product.SKU)
			if err != nil {
				fmt.Println("9.5: " + err.Error())
				failure_counter++
				continue
			}
		}
		_, err = dbconfig.DB.CreateVariantPricing(r.Context(), database.CreateVariantPricingParams{
			ID:        uuid.New(),
			VariantID: variant.ID,
			Name:      csv_product.PriceName,
			Value:     utils.ConvertStringToSQL(csv_product.PriceValue),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			fmt.Println("10: " + err.Error())
			failure_counter++
			continue
		}
		_, err = dbconfig.DB.CreateVariantQty(r.Context(), database.CreateVariantQtyParams{
			ID:        uuid.New(),
			VariantID: variant.ID,
			Name:      csv_product.QtyName,
			Value:     utils.ConvertIntToSQL(csv_product.QtyValue),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			fmt.Println("11: " + err.Error())
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
		SkipCounter:      skip_counter,
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
		_, err = dbconfig.DB.CreateAddress(r.Context(), database.CreateAddressParams{
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
	RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
		Message: customer.ID.String(),
	})
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
	if OrderValidation(order_body) != nil {
		RespondWithError(w, http.StatusBadRequest, "data validation error")
		return
	}
	customer, err := dbconfig.DB.CreateCustomer(r.Context(), database.CreateCustomerParams{
		ID:        uuid.New(),
		FirstName: order_body.Customer.FirstName,
		LastName:  order_body.Customer.FirstName,
		Email:     utils.ConvertStringToSQL(order_body.Customer.Email),
		Phone:     utils.ConvertStringToSQL(order_body.Customer.Phone),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	_, err = dbconfig.DB.CreateAddress(r.Context(), CreateDefaultAddress(order_body, customer.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	_, err = dbconfig.DB.CreateAddress(r.Context(), CreateShippingAddress(order_body, customer.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	_, err = dbconfig.DB.CreateAddress(r.Context(), CreateBillingAddress(order_body, customer.ID))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	order, err := dbconfig.DB.CreateOrder(r.Context(), database.CreateOrderParams{
		ID:            uuid.New(),
		Notes:         utils.ConvertStringToSQL(""),
		WebCode:       utils.ConvertStringToSQL(order_body.Name),
		TaxTotal:      utils.ConvertStringToSQL(order_body.TotalTax),
		OrderTotal:    utils.ConvertStringToSQL(order_body.TotalPrice),
		ShippingTotal: utils.ConvertStringToSQL(order_body.TotalShippingPriceSet.ShopMoney.Amount),
		DiscountTotal: utils.ConvertStringToSQL(order_body.TotalDiscounts),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	for _, value := range order_body.LineItems {
		if len(value.TaxLines) > 0 {
			_, err := dbconfig.DB.CreateOrderLine(r.Context(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("product"),
				Sku:       value.Sku,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(value.Quantity),
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		} else {
			_, err := dbconfig.DB.CreateOrderLine(r.Context(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("product"),
				Sku:       value.Sku,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(value.Quantity),
				TaxRate:   sql.NullString{},
				TaxTotal:  sql.NullString{},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
	}
	for _, value := range order_body.ShippingLines {
		if len(value.TaxLines) > 0 {
			_, err := dbconfig.DB.CreateOrderLine(r.Context(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("shipping"),
				Sku:       value.Code,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(1),
				TaxRate:   utils.ConvertStringToSQL(fmt.Sprintf("%v", value.TaxLines[0].Rate)),
				TaxTotal:  utils.ConvertStringToSQL(value.TaxLines[0].Price),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		} else {
			_, err := dbconfig.DB.CreateOrderLine(r.Context(), database.CreateOrderLineParams{
				ID:        uuid.New(),
				OrderID:   order.ID,
				LineType:  utils.ConvertStringToSQL("shipping"),
				Sku:       value.Code,
				Price:     utils.ConvertStringToSQL(value.Price),
				Barcode:   utils.ConvertIntToSQL(0),
				Qty:       utils.ConvertIntToSQL(1),
				TaxRate:   sql.NullString{},
				TaxTotal:  sql.NullString{},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
	}
	err = dbconfig.DB.CreateCustomerOrder(r.Context(), database.CreateCustomerOrderParams{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		OrderID:    order.ID,
		UpdatedAt:  time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
		Message: order.ID.String(),
	})
}

// POST /api/products/
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
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
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
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	for key := range params.ProductOptions {
		_, err := dbconfig.DB.CreateProductOption(r.Context(), database.CreateProductOptionParams{
			ID:        uuid.New(),
			ProductID: product.ID,
			Name:      params.ProductOptions[key].Value,
		})
		if err != nil {
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
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
			_, err = dbconfig.DB.CreateVariantQty(r.Context(), database.CreateVariantQtyParams{
				ID:        uuid.New(),
				VariantID: variant.ID,
				Name:      params.Variants[key].VariantQuantity[key_price].Name,
				Value:     utils.ConvertIntToSQL(params.Variants[key].VariantQuantity[key_price].Value),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	// TODO is it necessary to respond with the created product data
	product_added, err := CompileProductData(dbconfig, product.ID, r, false)
	if err != nil {
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
	customer, err := CompileCustomerData(dbconfig, customer_uuid, r, false)
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
		log.Println("Error decoding page param:", err)
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
		cust, err := CompileCustomerData(dbconfig, value.ID, r, true)
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
	order_data, err := CompileOrderData(dbconfig, order_uuid, r, false)
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
		log.Println("Error decoding page param:", err)
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
		ord, err := CompileOrderData(dbconfig, value.ID, r, true)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		orders = append(orders, ord)
	}
	RespondWithJSON(w, http.StatusOK, orders)
}

// GET /api/products/filter?data=value&page=1
func (dbconfig *DbConfig) ProductFilterHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
		log.Println("Error decoding page param:", err)
	}
	query_param_type := utils.ConfirmFilters(r.URL.Query().Get("type"))
	query_param_category := utils.ConfirmFilters(r.URL.Query().Get("category"))
	query_param_vendor := utils.ConfirmFilters(r.URL.Query().Get("vendor"))
	response, err := CompileFilterSearch(dbconfig, r, page, query_param_type, query_param_category, query_param_vendor)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// GET /api/products/search?q=value
func (dbconfig *DbConfig) ProductSearchHandle(w http.ResponseWriter, r *http.Request, dbuser database.User) {
	search_query := r.URL.Query().Get("q")
	if search_query != "" || len(search_query) == 0 {
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
	product_data, err := CompileProductData(dbconfig, product_uuid, r, false)
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
		log.Println("Error decoding page param:", err)
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
		prod, err := CompileProductData(dbconfig, value.ID, r, true)
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
	request_body, err := DecodePreRegisterRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	if PreRegisterValidation(request_body) != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	exists, err := dbconfig.CheckTokenExists(request_body, r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	if exists {
		RespondWithError(w, http.StatusConflict, utils.ConfirmError(err))
		return
	}
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
	err = SendEmail(token.Token, request_body.Email, request_body.Name)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
	}
	RespondWithJSON(w, http.StatusCreated, []string{"email sent"})
}

// POST /api/register
func (dbconfig *DbConfig) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeUserRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	if ValidateTokenValidation(body) != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	token, err := dbconfig.DB.GetTokenValidation(r.Context(), database.GetTokenValidationParams{
		Name:  body.Name,
		Email: body.Email,
	})
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
