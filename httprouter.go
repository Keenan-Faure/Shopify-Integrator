package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"iocsv"
	"log"
	"net/http"
	"objects"
	"os"
	"shopify"
	"strconv"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

// PUT /api/push/restriction
func (dbconfig *DbConfig) PushRestrictionHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	restrictions, err := DecodeRestriction(dbconfig, r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = RestrictionValidation(restrictions)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for _, value := range restrictions {
		err = dbconfig.DB.UpdatePushRestriction(r.Context(), database.UpdatePushRestrictionParams{
			Flag:      value.Flag,
			UpdatedAt: time.Now().UTC(),
			Field:     value.Field,
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// GET /api/push/restriction
func (dbconfig *DbConfig) GetPushRestrictionHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	restrictions, err := dbconfig.DB.GetPushRestriction(context.Background())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, "no push restrictions found found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, restrictions)
}

// GET /api/fetch/restriction
func (dbconfig *DbConfig) GetFetchRestrictionHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	restrictions, err := dbconfig.DB.GetFetchRestriction(context.Background())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, "no fetch restrictions found found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, restrictions)
}

// PUT /api/fetch/restriction
func (dbconfig *DbConfig) FetchRestrictionHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	restrictions, err := DecodeRestriction(dbconfig, r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = RestrictionValidation(restrictions)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for _, value := range restrictions {
		err = dbconfig.DB.UpdateFetchRestriction(r.Context(), database.UpdateFetchRestrictionParams{
			Flag:      value.Flag,
			UpdatedAt: time.Now().UTC(),
			Field:     value.Field,
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// POST /api/worker/fetch
func (dbconfig *DbConfig) WorkerFetchProductsHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	// create database table containing the status of this
	shopifyConfig := shopify.InitConfigShopify()
	err := FetchShopifyProducts(dbconfig, shopifyConfig)
	if err != nil {
		if err.Error() == "worker is currently running" {
			RespondWithError(w, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// should never be used in production

// PUT /api/shopify/reset_fetch
func (dbconfig *DbConfig) ResetShopifyFetchHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	err := dbconfig.DB.ResetFetchWorker(context.Background(), "0")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// POST /api/inventory/warehouse?reindex=false
func (dbconfig *DbConfig) AddInventoryWarehouse(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	reindex := r.URL.Query().Get("reindex")
	warehouse, err := DecodeGlobalWarehouse(dbconfig, r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = GlobalWarehouseValidation(warehouse)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if reindex == "true" {
		// check if the warehouse exists internally
		warehouse_db, err := dbconfig.DB.GetWarehouseByName(r.Context(), warehouse.Name)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			} else {
				RespondWithError(w, http.StatusInternalServerError, "cannot reindex an invalid warehouse")
				return
			}
		}
		// reindex warehouse that was found
		err = InsertGlobalWarehouse(dbconfig, r.Context(), warehouse_db.Name, true)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
			Message: "success",
		})
		return
	}
	// check if a warehouse already exists
	warehouses_db, err := dbconfig.DB.GetWarehouses(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, warehouse_db := range warehouses_db {
		if warehouse_db.Name == warehouse.Name {
			RespondWithError(w, http.StatusBadRequest, "warehouse already exists")
			return
		}
	}
	err = dbconfig.DB.CreateWarehouse(r.Context(), database.CreateWarehouseParams{
		ID:        uuid.New(),
		Name:      warehouse.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = InsertGlobalWarehouse(dbconfig, r.Context(), warehouse.Name, false)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
		Message: "success",
	})
}

// GET /api/inventory/warehouse
func (dbconfig *DbConfig) GetInventoryWarehouses(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	warehouses, err := dbconfig.DB.GetWarehouses(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(warehouses) == 0 {
		warehouses = []database.GetWarehousesRow{}
	}
	RespondWithJSON(w, http.StatusOK, warehouses)
}

// GET /api/inventory/warehouse/{id}
func (dbconfig *DbConfig) GetInventoryWarehouse(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	warehouse_id := chi.URLParam(r, "id")
	err := IDValidation(warehouse_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	warehouse_uuid, err := uuid.Parse(warehouse_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode order id: "+warehouse_id)
		return
	}
	warehouse, err := dbconfig.DB.GetWarehouseByID(r.Context(), warehouse_uuid)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, warehouse)
}

// DELETE /api/inventory/warehouse/{id}
func (dbconfig *DbConfig) DeleteInventoryWarehouse(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	warehouse_id := chi.URLParam(r, "id")
	err := IDValidation(warehouse_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	warehouse_uuid, err := uuid.Parse(warehouse_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode order id: "+warehouse_id)
		return
	}
	err = dbconfig.DB.RemoveWarehouse(r.Context(), warehouse_uuid)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	// remove all variant warehouses

	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// PUT /api/products/{id}
func (dbconfig *DbConfig) UpdateProductHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	product_id := chi.URLParam(r, "id")
	err := IDValidation(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	product_uuid, err := uuid.Parse(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode product id: "+product_id)
		return
	}
	found := false
	_, err = dbconfig.DB.GetProductByID(r.Context(), product_uuid)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			found = false
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		found = true
	}

	if !found {
		RespondWithError(w, http.StatusNotFound, "could not find product id: "+product_id)
		return
	}

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
	err = DuplicateOptionValues(params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}

	// update product
	err = dbconfig.DB.UpdateProductByID(r.Context(), database.UpdateProductByIDParams{
		Active:      params.Active,
		Title:       utils.ConvertStringToSQL(params.Title),
		BodyHtml:    utils.ConvertStringToSQL(params.BodyHTML),
		Category:    utils.ConvertStringToSQL(params.Category),
		Vendor:      utils.ConvertStringToSQL(params.Vendor),
		ProductType: utils.ConvertStringToSQL(params.ProductType),
		UpdatedAt:   time.Now().UTC(),
		ID:          product_uuid,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}

	for key := range params.ProductOptions {
		// TODO Should we use the position in the POST Body or the key that is it inside the array?
		_, err = dbconfig.DB.UpdateProductOption(r.Context(), database.UpdateProductOptionParams{
			Name:       params.ProductOptions[key].Value,
			Position:   int32(key + 1),
			ProductID:  product_uuid,
			Position_2: int32(key + 1),
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	for _, variant := range params.Variants {
		err = dbconfig.DB.UpdateVariant(r.Context(), database.UpdateVariantParams{
			Option1:   utils.ConvertStringToSQL(variant.Option1),
			Option2:   utils.ConvertStringToSQL(variant.Option2),
			Option3:   utils.ConvertStringToSQL(variant.Option3),
			Barcode:   utils.ConvertStringToSQL(variant.Barcode),
			UpdatedAt: time.Now().UTC(),
			Sku:       variant.Sku,
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
		// update variant pricing and qty here
		for _, price_lists := range variant.VariantPricing {
			err = dbconfig.DB.UpdateVariantPricing(r.Context(), database.UpdateVariantPricingParams{
				Name:      price_lists.Name,
				Value:     utils.ConvertStringToSQL(price_lists.Value),
				Isdefault: price_lists.IsDefault,
				UpdatedAt: time.Now().UTC(),
				Sku:       variant.Sku,
				Name_2:    price_lists.Name,
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
		for _, warehouses := range variant.VariantQuantity {
			err = dbconfig.DB.UpdateVariantQty(r.Context(), database.UpdateVariantQtyParams{
				Name:      warehouses.Name,
				Value:     utils.ConvertIntToSQL(warehouses.Value),
				Isdefault: warehouses.IsDefault,
				Sku:       variant.Sku,
				Name_2:    warehouses.Name,
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
	}
	updated_data, err := CompileProductData(dbconfig, product_uuid, r.Context(), false)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	// only update if the active = 1
	if updated_data.Active == "1" {
		err = CompileInstructionProduct(dbconfig, updated_data, dbUser)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		for _, variant := range updated_data.Variants {
			err = CompileInstructionVariant(dbconfig, variant, updated_data, dbUser)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	RespondWithJSON(w, http.StatusOK, updated_data)
}

// GET /api/stats/fetch
func (dbconfig *DbConfig) GetFetchStats(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	data, err := dbconfig.DB.GetFetchStats(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	// convert data to include missing dates, and convert dates to appropriate values
	RespondWithJSON(w, http.StatusOK, ParseFetchStats(data))
}

// GET /api/stats/orders?status=paid
func (dbconfig *DbConfig) GetOrderStats(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	status := r.URL.Query().Get("status")
	if status == "paid" {
		data, err := dbconfig.DB.FetchOrderStatsPaid(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
			return
		}
		// convert data to include missing dates, and convert dates to appropriate values
		RespondWithJSON(w, http.StatusOK, ParseOrderStatsPaid(data))
	} else if status == "not_paid" {
		data, err := dbconfig.DB.FetchOrderStatsNotPaid(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
			return
		}
		// convert data to include missing dates, and convert dates to appropriate values
		RespondWithJSON(w, http.StatusOK, ParseOrderStatsNotPaid(data))
	} else {
		RespondWithError(w, http.StatusBadRequest, "invalid status type")
		return
	}
}

// POST /api/settings/webhook
func (dbconfig *DbConfig) GetWebhookURL(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	body, err := DecodeWebhookURL(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	// create webhook url
	webhook_url := body.Domain + "/api/orders?token=" + dbUser.WebhookToken + "&api_key=" + dbUser.ApiKey
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: webhook_url,
	})
}

// GET /api/inventory/config
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

// DELETE /api/inventory/config
func (dbconfig *DbConfig) RemoveWarehouseLocation(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	id := chi.URLParam(r, "id")
	err := IDValidation(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	delete_id, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode id: "+id)
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

// Returns all locations on Shopify and all warehouses in the app

// GET /api/inventory/map
func (dbconfig *DbConfig) ConfigLocationWarehouse(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	shopifyConfig := shopify.InitConfigShopify()
	if !shopifyConfig.Valid {
		RespondWithError(w, http.StatusInternalServerError, "invalid shopify config")
		return
	}
	locations, err := shopifyConfig.GetShopifyLocations()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	warehouses, err := dbconfig.DB.GetWarehouses(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseWarehouseLocation{
		Warehouses:       ConvertDatabaseToWarehouse(warehouses),
		ShopifyLocations: locations,
	})
}

// POST /api/inventory/config
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

// GET /api/products/export?test=true
func (dbconfig *DbConfig) ExportProductsHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	test := r.URL.Query().Get("test")
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
	if len(products) > 0 {
		image_headers := iocsv.GetProductImagesCSV(products[0].ProductImages, int(images_max), true)
		headers = append(headers, image_headers...)
	}
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
	// removes the product from the server if there are tests
	if test == "true" {
		defer os.Remove(file_name)
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: file_name,
	})
}

// POST /api/products/import?test=true
func (dbconfig *DbConfig) ProductImportHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	test := r.URL.Query().Get("test")
	file_name_global := "test_import.csv"
	if test == "true" {
		// generate the file for the test and ignore upload form
		data := [][]string{
			{"type", "active", "product_code", "title", "body_html", "category", "vendor", "product_type", "sku", "option1_name", "option1_value", "option2_name", "option2_value", "option3_name", "option3_value", "barcode", "price_default"},
			{"product", "1", "grouper", "test_title", "<p>I am a paragraph</p>", "test_category", "test_vendor", "test_product_type", "skubca", "size", "medium", "color", "blue", "", "", "", "1500.00"},
		}
		_, err := iocsv.WriteFile(data, "test_import")
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// if in production then expect form data &&
		// file to exist in import
		file_name, err := iocsv.UploadFile(r, "")
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		file_name_global = file_name
	}
	wd, err := os.Getwd()
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	csv_products, err := iocsv.ReadFile(wd + "/" + file_name_global)
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
			Active:      csv_product.Active,
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
					Active:      csv_product.Active,
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
	err = iocsv.RemoveFile(file_name_global)
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
			Type:       utils.ConvertStringToSQL(customer_body.Address[key].Type),
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
func (dbconfig *DbConfig) PostOrderHandle(w http.ResponseWriter, r *http.Request) {
	web_token := r.URL.Query().Get("token")
	if TokenValidation(web_token) != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid token")
		return
	}
	api_key := r.URL.Query().Get("api_key")
	if TokenValidation(api_key) != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid api_key")
		return
	}
	_, err := dbconfig.DB.ValidateWebhookByUser(r.Context(), database.ValidateWebhookByUserParams{
		WebhookToken: web_token,
		ApiKey:       api_key,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, "invalid token for user")
			return
		} else {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
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
	db_order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), utils.ConvertStringToSQL(fmt.Sprint(order_body.Name)))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
			return
		}
	}
	if db_order.WebCode.String == fmt.Sprint(order_body.Name) {
		response_payload, err := dbconfig.QueueHelper(objects.RequestQueueHelper{
			Type:        "order",
			Status:      "in-queue",
			Instruction: "update_order",
			Endpoint:    "queue",
			ApiKey:      api_key,
			Method:      http.MethodPost,
			Object:      order_body,
		})
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
			ApiKey:      api_key,
			Method:      http.MethodPost,
			Object:      order_body,
		})
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
			return
		}
		RespondWithJSON(w, http.StatusCreated, response_payload)
	}
}

// DELETE /api/products?id={{product_id}}
func (dbconfig *DbConfig) RemoveProductHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	product_id := chi.URLParam(r, "id")
	err := IDValidation(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	product_uuid, err := uuid.Parse(product_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode product id: "+product_id)
		return
	}
	err = dbconfig.DB.RemoveProduct(r.Context(), product_uuid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
}

// DELETE /api/products?variant_id={{variant_id}}
func (dbconfig *DbConfig) RemoveProductVariantHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	variant_id := chi.URLParam(r, "variant_id")
	err := IDValidation(variant_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, utils.ConfirmError(err))
		return
	}
	variant_uuid, err := uuid.Parse(variant_id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "could not decode variant id: "+variant_id)
		return
	}
	err = dbconfig.DB.RemoveVariant(r.Context(), variant_uuid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
		Message: "success",
	})
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
		Active:      params.Active,
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
			Position:  int32(key + 1),
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
		for key_pricing := range params.Variants[key].VariantPricing {
			_, err := dbconfig.DB.CreateVariantPricing(r.Context(), database.CreateVariantPricingParams{
				ID:        uuid.New(),
				VariantID: variant.ID,
				Name:      params.Variants[key].VariantPricing[key_pricing].Name,
				Value:     utils.ConvertStringToSQL(params.Variants[key].VariantPricing[key_pricing].Value),
				Isdefault: params.Variants[key].VariantPricing[key_pricing].IsDefault,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
		}
		for key_qty := range params.Variants[key].VariantQuantity {
			// check if the warehouse exists, then we update the quantity
			warehouse_name := params.Variants[key].VariantQuantity[key_qty].Name
			warehouse_qty := params.Variants[key].VariantQuantity[key_qty].Value
			_, err = dbconfig.DB.GetWarehouseByName(context.Background(), warehouse_name)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(w, http.StatusInternalServerError, "warehouse "+warehouse_name+" not found")
					return
				}
				RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
				return
			}
			// if warehouse is found, we update the qty, we cannot create a new one
			err = dbconfig.DB.UpdateVariantQty(context.Background(), database.UpdateVariantQtyParams{
				Name:      warehouse_name,
				Value:     utils.ConvertIntToSQL(warehouse_qty),
				Isdefault: false,
				UpdatedAt: time.Now().UTC(),
				Sku:       variant.Sku,
				Name_2:    warehouse_name,
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
	product_added, err := CompileProductData(dbconfig, product.ID, r.Context(), false)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// queue new products to be added to shopify
	err = CompileInstructionProduct(dbconfig, product_added, dbUser)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, variant := range product_added.Variants {
		err = CompileInstructionVariant(dbconfig, variant, product_added, dbUser)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	RespondWithJSON(w, http.StatusCreated, objects.ResponseString{
		Message: "success",
	})
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
		RespondWithError(w, http.StatusBadRequest, "could not decode customer id: "+customer_id)
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
	webcode_orders, err := dbconfig.DB.GetOrdersSearchWebCode(r.Context(), utils.ConvertStringToLike(search_query))
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
		RespondWithError(w, http.StatusBadRequest, "could not decode order id: "+order_id)
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
	response, err := CompileFilterSearch(
		dbconfig,
		r.Context(),
		page,
		utils.ConvertStringToLike(query_param_type),
		utils.ConvertStringToLike(query_param_category),
		utils.ConvertStringToLike(query_param_vendor),
	)
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
	search, err := dbconfig.DB.GetProductsSearch(r.Context(), utils.ConvertStringToLike(search_query))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	compiled, err := CompileSearchResult(dbconfig, r.Context(), search)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, utils.ConfirmError(err))
		return
	}
	RespondWithJSON(w, http.StatusOK, compiled)
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
		RespondWithError(w, http.StatusBadRequest, "could not decode product id: "+product_id)
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
	RespondWithJSON(w, http.StatusOK, objects.RequestString{
		Message: "success",
	})
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
		RespondWithError(w, http.StatusBadRequest, "could not decode token: "+body.Token)
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
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload == nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		_, err = w.Write(response)
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	} else {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		_, err = w.Write(response)
		if err != nil {
			log.Printf("Error writing JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	}
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

// Middleware that determines which headers, http methods and orgins are allowed
func MiddleWare() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}
