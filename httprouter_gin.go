package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Gets the paginated list of warehouse from the application

Route: /api/inventory/warehouse?page=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetInventoryWarehouses() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		warehouses, err := dbconfig.DB.GetWarehouses(c.Request.Context(), database.GetWarehousesParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if len(warehouses) == 0 {
			warehouses = []database.GetWarehousesRow{}
		}
		RespondWithJSON(c, http.StatusOK, warehouses)
	}
}

/*
Gets the specific warehouse from the application

Route: /api/inventory/warehouse/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetInventoryWarehouse() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouse_id := c.Param("id")
		err := IDValidation(warehouse_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		warehouse_uuid, err := uuid.Parse(warehouse_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode order id: "+warehouse_id)
			return
		}
		warehouse, err := dbconfig.DB.GetWarehouseByID(c.Request.Context(), warehouse_uuid)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, warehouse)
	}
}

/*
Removes the specific warehouse from the application

Route: /api/inventory/warehouse/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) DeleteInventoryWarehouse() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouse_id := c.Param("id")
		err := IDValidation(warehouse_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		warehouse_uuid, err := uuid.Parse(warehouse_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode order id: "+warehouse_id)
			return
		}
		err = dbconfig.DB.RemoveWarehouse(c.Request.Context(), warehouse_uuid)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// remove all variant warehouses

		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Updates a product by its ID

Route: /api/products/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) UpdateProductHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id: "+product_id)
			return
		}
		found := false
		_, err = dbconfig.DB.GetProductByID(c.Request.Context(), product_uuid)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				found = false
			} else {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			found = true
		}

		if !found {
			RespondWithError(c, http.StatusNotFound, "could not find product id: "+product_id)
			return
		}

		params, err := DecodeProductRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		validation := ProductValidation(dbconfig, params)
		if validation != nil {
			RespondWithError(c, http.StatusBadRequest, validation.Error())
			return
		}
		err = ValidateDuplicateOption(params)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = DuplicateOptionValues(params)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}

		// update product
		err = dbconfig.DB.UpdateProductByID(c.Request.Context(), database.UpdateProductByIDParams{
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
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		for key := range params.ProductOptions {
			// TODO Should we use the position in the POST Body or the key that is it inside the array?
			_, err = dbconfig.DB.UpdateProductOption(c.Request.Context(), database.UpdateProductOptionParams{
				Name:       params.ProductOptions[key].Value,
				Position:   int32(key + 1),
				ProductID:  product_uuid,
				Position_2: int32(key + 1),
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		for _, variant := range params.Variants {
			err = dbconfig.DB.UpdateVariant(c.Request.Context(), database.UpdateVariantParams{
				Option1:   utils.ConvertStringToSQL(variant.Option1),
				Option2:   utils.ConvertStringToSQL(variant.Option2),
				Option3:   utils.ConvertStringToSQL(variant.Option3),
				Barcode:   utils.ConvertStringToSQL(variant.Barcode),
				UpdatedAt: time.Now().UTC(),
				Sku:       variant.Sku,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			// update variant pricing and qty here
			for _, price_lists := range variant.VariantPricing {
				// check if the pricing is acceptable
				if price_lists.Name == "Selling Price" || price_lists.Name == "Compare At Price" {
					err = dbconfig.DB.UpdateVariantPricing(c.Request.Context(), database.UpdateVariantPricingParams{
						Name:      price_lists.Name,
						Value:     utils.ConvertStringToSQL(price_lists.Value),
						Isdefault: price_lists.IsDefault,
						UpdatedAt: time.Now().UTC(),
						Sku:       variant.Sku,
						Name_2:    price_lists.Name,
					})
					if err != nil {
						RespondWithError(c, http.StatusInternalServerError, err.Error())
						return
					}
				} else {
					RespondWithError(c, http.StatusInternalServerError, "invalid price tier "+price_lists.Name)
					return
				}
			}
			for _, warehouses := range variant.VariantQuantity {
				err = dbconfig.DB.UpdateVariantQty(c.Request.Context(), database.UpdateVariantQtyParams{
					Name:      warehouses.Name,
					Value:     utils.ConvertIntToSQL(warehouses.Value),
					Isdefault: warehouses.IsDefault,
					Sku:       variant.Sku,
					Name_2:    warehouses.Name,
				})
				if err != nil {
					RespondWithError(c, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		updated_data, err := CompileProductData(dbconfig, product_uuid, c.Request.Context(), false)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		// only update if the active = 1
		if updated_data.Active == "1" {
			api_key := c.GetString("api_key")
			err = CompileInstructionProduct(dbconfig, updated_data, api_key)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			for _, variant := range updated_data.Variants {
				err = CompileInstructionVariant(dbconfig, variant, updated_data, api_key)
				if err != nil {
					RespondWithError(c, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		RespondWithJSON(c, http.StatusOK, updated_data)
	}
}

/*
Returns the shopify fetch stats that are recorded internally.

Route: /api/stats/fetch

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetFetchStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := dbconfig.DB.GetFetchStats(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, ParseFetchStats(data))
	}
}

/*
Returns the internal stats of orders; either "paid" or "not_paid"

Route: /api/stats/orders?status=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetOrderStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		if status == "paid" {
			data, err := dbconfig.DB.FetchOrderStatsPaid(c.Request.Context())
			if err != nil {
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			// convert data to include missing dates, and convert dates to appropriate values
			RespondWithJSON(c, http.StatusOK, ParseOrderStatsPaid(data))
		} else if status == "not_paid" {
			data, err := dbconfig.DB.FetchOrderStatsNotPaid(c.Request.Context())
			if err != nil {
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			// convert data to include missing dates, and convert dates to appropriate values
			RespondWithJSON(c, http.StatusOK, ParseOrderStatsNotPaid(data))
		} else {
			RespondWithError(c, http.StatusBadRequest, "invalid status type")
			return
		}
	}
}

/*
Returns the location-warehouse map currently set internally

Route: /api/inventory/map

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) LocationWarehouseHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		shopify_locations, err := dbconfig.DB.GetShopifyLocations(c.Request.Context(), database.GetShopifyLocationsParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		locations := []database.ShopifyLocation{}
		if len(shopify_locations) > 0 {
			locations = append(locations, shopify_locations...)
		}
		RespondWithJSON(c, http.StatusOK, locations)
	}
}

/*
Adds a specific warehouse-location map internally

Route: /api/inventory/map

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddWarehouseLocationMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		location_map, err := DecodeInventoryMap(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		if InventoryMapValidation(location_map) != nil {
			RespondWithError(c, http.StatusBadRequest, "data validation error")
			return
		}
		result, err := dbconfig.DB.CreateShopifyLocation(c.Request.Context(), database.CreateShopifyLocationParams{
			ID:                   uuid.New(),
			ShopifyWarehouseName: location_map.ShopifyWarehouseName,
			ShopifyLocationID:    location_map.LocationID,
			WarehouseName:        location_map.WarehouseName,
			CreatedAt:            time.Now().UTC(),
			UpdatedAt:            time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, result)
	}
}

/*
Removes the specific warehouse-location map internally

Route: /api/inventory/map/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) RemoveWarehouseLocation() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := IDValidation(id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		delete_id, err := uuid.Parse(id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode id: "+id)
			return
		}
		err = dbconfig.DB.RemoveShopifyLocationMap(c.Request.Context(), delete_id)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Returns an object specifying the shopify locations and internal warehouses

Route: /api/inventory/config

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ConfigLocationWarehouseHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		shopifyConfig := shopify.InitConfigShopify()
		if !shopifyConfig.Valid {
			RespondWithError(c, http.StatusInternalServerError, "invalid shopify config")
			return
		}
		locations, err := shopifyConfig.GetShopifyLocations()
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		warehouses, err := dbconfig.DB.GetWarehouses(c.Request.Context(), database.GetWarehousesParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseWarehouseLocation{
			Warehouses:       ConvertDatabaseToWarehouse(warehouses),
			ShopifyLocations: locations,
		})
	}
}

/*
Creates and adds a new customer to the application

Route: /api/customers

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostCustomerHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		customer_body, err := DecodeCustomerRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		if CustomerValidation(customer_body) != nil {
			RespondWithError(c, http.StatusBadRequest, "data validation error")
			return
		}
		customer, err := dbconfig.DB.CreateCustomer(c.Request.Context(), database.CreateCustomerParams{
			ID:        uuid.New(),
			FirstName: customer_body.FirstName,
			LastName:  customer_body.LastName,
			Email:     utils.ConvertStringToSQL(customer_body.Email),
			Phone:     utils.ConvertStringToSQL(customer_body.Phone),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		for key := range customer_body.Address {
			_, err := dbconfig.DB.CreateAddress(c.Request.Context(), database.CreateAddressParams{
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
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		customer_data, err := CompileCustomerData(dbconfig, customer.ID, c.Request.Context(), false)
		if err != nil {
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		RespondWithJSON(c, http.StatusCreated, customer_data)
	}
}

/*
Returns the results of a search query by the customer name and web code of the order

Route: /api/customers/search?q=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomerSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query != "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		customers_by_name, err := dbconfig.DB.GetCustomersByName(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		RespondWithJSON(c, http.StatusOK, customers_by_name)
	}
}

/*
Returns the customer data having the specific id

Route: /api/customers/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomerIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		customer_id := c.Param("id")
		err := IDValidation(customer_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		customer_uuid, err := uuid.Parse(customer_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode customer id: "+customer_id)
			return
		}
		customer, err := CompileCustomerData(dbconfig, customer_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, customer)
	}
}

/*
Returns the respective page of customer data from the database

Route: /api/customers?page=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) CustomersHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbCustomers, err := dbconfig.DB.GetCustomers(c.Request.Context(), database.GetCustomersParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		customers := []objects.Customer{}
		for _, value := range dbCustomers {
			cust, err := CompileCustomerData(dbconfig, value.ID, c.Request.Context(), true)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			customers = append(customers, cust)
		}
		RespondWithJSON(c, http.StatusOK, customers)
	}
}

/*
Queues the respective order to be added to the application from Shopify

Route: /api/orders?token={token}&api_key={api_key}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostOrderHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		web_token := c.Query("token")
		if TokenValidation(web_token) != nil {
			RespondWithError(c, http.StatusBadRequest, "invalid token")
			return
		}
		api_key := c.Query("api_key")
		if TokenValidation(api_key) != nil {
			RespondWithError(c, http.StatusBadRequest, "invalid api_key")
			return
		}
		_, err := dbconfig.DB.ValidateWebhookByUser(c.Request.Context(), database.ValidateWebhookByUserParams{
			WebhookToken: web_token,
			ApiKey:       api_key,
		})
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusInternalServerError, "invalid token for user")
				return
			} else {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		order_body, err := DecodeOrderRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		var buffer bytes.Buffer
		err = json.NewEncoder(&buffer).Encode(order_body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		db_order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), utils.ConvertStringToSQL(fmt.Sprint(order_body.Name)))
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
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
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, response_payload)
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
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusCreated, response_payload)
		}
	}
}

/*
Returns the results of a search query by the customer name and web code of the order

Route: /api/orders/search?q=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrderSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query != "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		customer_orders, err := dbconfig.DB.GetOrdersSearchByCustomer(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		webcode_orders, err := dbconfig.DB.GetOrdersSearchWebCode(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
		}
		RespondWithJSON(c, http.StatusOK, CompileOrderSearchResult(customer_orders, webcode_orders))
	}
}

/*
Returns the order data having the specific id

Route: /api/orders/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrderIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		order_id := c.Param("id")
		err := IDValidation(order_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		order_uuid, err := uuid.Parse(order_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode order id: "+order_id)
			return
		}
		order_data, err := CompileOrderData(dbconfig, order_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, order_data)
	}
}

/*
Returns the respective page of order data from the database

Route: /api/orders?page=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) OrdersHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbOrders, err := dbconfig.DB.GetOrders(c.Request.Context(), database.GetOrdersParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		orders := []objects.Order{}
		for _, value := range dbOrders {
			ord, err := CompileOrderData(dbconfig, value.ID, c.Request.Context(), true)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			orders = append(orders, ord)
		}
		RespondWithJSON(c, http.StatusOK, orders)
	}
}

/*
Removes the specific product from the application

Route: /api/products/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) RemoveProductHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id: "+product_id)
			return
		}
		err = dbconfig.DB.RemoveProduct(c.Request.Context(), product_uuid)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Removes the specific variant from a product

Route: /api/products/{variant_id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) RemoveProductVariantHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		variant_id := c.Param("variant_id")
		err := IDValidation(variant_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		variant_uuid, err := uuid.Parse(variant_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode variant id: "+variant_id)
			return
		}
		err = dbconfig.DB.RemoveVariant(c.Request.Context(), variant_uuid)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Exports product data to a .CSV file.

Route: /api/products/export

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ExportProductsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		test := c.Query("test")
		product_ids, err := dbconfig.DB.GetProductIDs(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		products := []objects.Product{}
		for _, product_id := range product_ids {
			product, err := CompileProductData(dbconfig, product_id, c.Request.Context(), false)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
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
		price_tiers, err := AddPricingHeaders(dbconfig, c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		headers = append(headers, price_tiers...)
		qty_warehouses, err := AddQtyHeaders(dbconfig, c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		headers = append(headers, qty_warehouses...)

		images_max, err := IOGetMax(dbconfig, c.Request.Context(), "image")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		pricing_max, err := IOGetMax(dbconfig, c.Request.Context(), "price")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		qty_max, err := IOGetMax(dbconfig, c.Request.Context(), "qty")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if len(products) > 0 {
			image_headers := iocsv.GetProductImagesCSV(products[0].ProductImages, int(images_max), true)
			headers = append(headers, image_headers...)
		}
		csv_data = append(csv_data, headers)
		for _, product := range products {
			for _, variant := range product.Variants {
				row := iocsv.CSVProductValuesByVariant(product, variant, int(images_max), pricing_max, qty_max)
				csv_data = append(csv_data, row)
			}
		}
		file_name, err := iocsv.WriteFile(csv_data, "")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// removes the product from the server if there are tests
		if test == "true" {
			defer os.Remove(file_name)
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: file_name,
		})
	}
}

/*
Uses a .CSV file to bulk import products into the application. The file must be sent with the request.

Route: /api/products/import

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ProductImportHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		test := c.Query("test")
		file_name_global := "test_import.csv"
		if test == "true" {
			// generate the file for the test and ignore upload form
			data := [][]string{
				{"type", "active", "product_code", "title", "body_html", "category", "vendor", "product_type", "sku", "option1_name", "option1_value", "option2_name", "option2_value", "option3_name", "option3_value", "barcode", "price_Selling Price"},
				{"product", "1", "grouper", "test_title", "<p>I am a paragraph</p>", "test_category", "test_vendor", "test_product_type", "skubca", "size", "medium", "color", "blue", "", "", "", "1500.00"},
			}
			_, err := iocsv.WriteFile(data, "test_import")
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			// if in production then expect form data &&
			// file to exist in import
			file_name, err := iocsv.UploadFile(c.Request)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			file_name_global = file_name
		}
		wd, err := os.Getwd()
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		csv_products, err := iocsv.ReadFile(wd + "/" + file_name_global)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		processed_counter := 0
		failure_counter := 0
		products_added := 0
		products_updated := 0
		variants_updated := 0
		variants_added := 0
		for _, csv_product := range csv_products {
			product, err := dbconfig.DB.UpsertProduct(c.Request.Context(), database.UpsertProductParams{
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
				log.Println(err)
				failure_counter++
				continue
			}
			if product.Inserted {
				products_added++
			} else {
				products_updated++
			}
			option_names := CreateOptionNamesMap(csv_product)
			err = AddProductOptions(dbconfig, product.ID, product.ProductCode, option_names)
			if err != nil {
				log.Println(err)
				failure_counter++
				continue
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
			variant, err := dbconfig.DB.UpsertVariant(c.Request.Context(), database.UpsertVariantParams{
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
			if err == nil {
				variants_added++
			}
			if variant.Inserted {
				variants_added++
			} else {
				variants_updated++
			}
			for _, pricing_value := range csv_product.Pricing {
				// check if the price is acceptable
				if pricing_value.Name == "Selling Price" || pricing_value.Name == "Compare At Price" {
					err = AddPricing(dbconfig, csv_product.SKU, variant.ID, pricing_value.Name, pricing_value.Value)
					if err != nil {
						log.Println(err)
						failure_counter++
						continue
					}
				} else {
					log.Println("invalid price tier " + pricing_value.Name)
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
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ImportResponse{
			ProcessedCounter: processed_counter,
			FailCounter:      failure_counter,
			ProductsAdded:    products_added,
			ProductsUpdated:  products_updated,
			VariantsAdded:    variants_added,
			VariantsUpdated:  variants_updated,
		})
	}
}

/*
Creates and adds a new product to the application

Route: /api/products

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostProductHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params objects.RequestBodyProduct
		err := c.Bind(params)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		validation := ProductValidation(dbconfig, params)
		if validation != nil {
			RespondWithError(c, http.StatusBadRequest, validation.Error())
			return
		}
		err = ValidateDuplicateOption(params)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = ValidateDuplicateSKU(params, dbconfig, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		err = DuplicateOptionValues(params)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		csv_products := ConvertProductToCSV(params)
		for _, csv_product := range csv_products {
			err = ProductValidationDatabase(csv_product, dbconfig, c.Request)
			if err != nil {
				RespondWithError(c, http.StatusBadRequest, err.Error())
				return
			}
		}
		product, err := dbconfig.DB.CreateProduct(c.Request.Context(), database.CreateProductParams{
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
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		for key := range params.ProductOptions {
			_, err := dbconfig.DB.CreateProductOption(c.Request.Context(), database.CreateProductOptionParams{
				ID:        uuid.New(),
				ProductID: product.ID,
				Name:      params.ProductOptions[key].Value,
				Position:  int32(key + 1),
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		for key := range params.Variants {
			variant, err := dbconfig.DB.CreateVariant(c.Request.Context(), database.CreateVariantParams{
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
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			// variant pricing & variant qty
			for key_pricing := range params.Variants[key].VariantPricing {
				// check if the price tier name is acceptable
				if params.Variants[key].VariantPricing[key_pricing].Name == "Selling Price" ||
					params.Variants[key].VariantPricing[key_pricing].Name == "Compare At Price" {
					_, err := dbconfig.DB.CreateVariantPricing(c.Request.Context(), database.CreateVariantPricingParams{
						ID:        uuid.New(),
						VariantID: variant.ID,
						Name:      params.Variants[key].VariantPricing[key_pricing].Name,
						Value:     utils.ConvertStringToSQL(params.Variants[key].VariantPricing[key_pricing].Value),
						Isdefault: params.Variants[key].VariantPricing[key_pricing].IsDefault,
						CreatedAt: time.Now().UTC(),
						UpdatedAt: time.Now().UTC(),
					})
					if err != nil {
						RespondWithError(c, http.StatusInternalServerError, err.Error())
						return
					}
				} else {
					RespondWithError(c, http.StatusInternalServerError, "invalid price tier"+
						params.Variants[key].VariantPricing[key_pricing].Name)
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
						RespondWithError(c, http.StatusInternalServerError, "warehouse "+warehouse_name+" not found")
						return
					}
					RespondWithError(c, http.StatusInternalServerError, err.Error())
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
					RespondWithError(c, http.StatusInternalServerError, err.Error())
					return
				}
			}
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		product_added, err := CompileProductData(dbconfig, product.ID, c.Request.Context(), false)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// queue new products to be added to shopify
		api_key := c.GetString("api_key")
		err = CompileInstructionProduct(dbconfig, product_added, api_key)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		for _, variant := range product_added.Variants {
			err = CompileInstructionVariant(dbconfig, variant, product_added, api_key)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		RespondWithJSON(c, http.StatusCreated, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Filter Searches for certain products based on their vendor, product type and collection

Route: /api/products/filter?page=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ProductFilterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		query_param_type := utils.ConfirmFilters(c.Query("type"))
		query_param_category := utils.ConfirmFilters(c.Query("category"))
		query_param_vendor := utils.ConfirmFilters(c.Query("vendor"))
		response, err := CompileFilterSearch(
			dbconfig,
			c.Request.Context(),
			page,
			utils.ConvertStringToLike(query_param_type),
			utils.ConvertStringToLike(query_param_category),
			utils.ConvertStringToLike(query_param_vendor),
		)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, response)
	}
}

/*
Returns the results of a search query by a product Title and SKU

Route: /api/products/search?q=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ProductSearchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		search_query := c.Query("q")
		if search_query == "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "Invalid search param")
			return
		}
		search, err := dbconfig.DB.GetProductsSearch(c.Request.Context(), utils.ConvertStringToLike(search_query))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		compiled, err := CompileSearchResult(dbconfig, c.Request.Context(), search)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, compiled)
	}
}

/*
Returns the respective page of product data from the database

Route: /api/products?page=

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 404, 401, 500
*/
func (dbconfig *DbConfig) ProductsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}
		dbProducts, err := dbconfig.DB.GetProducts(c.Request.Context(), database.GetProductsParams{
			Limit:  10,
			Offset: int32((page - 1) * 10),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		products := []objects.Product{}
		for _, value := range dbProducts {
			prod, err := CompileProductData(dbconfig, value.ID, c.Request.Context(), false)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			products = append(products, prod)
		}
		RespondWithJSON(c, http.StatusOK, products)
	}
}

/*
Returns the product data having the specific id

Route: /api/product/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 404, 401, 500
*/
func (dbconfig *DbConfig) ProductIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id '"+product_id+"'")
			return
		}
		product_data, err := CompileProductData(dbconfig, product_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, product_data)
	}
}

/*
Logs a user into the application. This does not set any cookies

Route: /api/login

Authorization: Required

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) LoginHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := DecodeLoginRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = UserValidation(body.Username, body.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		db_user, exists, err := dbconfig.CheckUserCredentials(body, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			RespondWithError(c, http.StatusNotFound, "invalid username and password combination")
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseLogin{
			Username: db_user.Name,
			ApiKey:   db_user.ApiKey,
		})
	}
}

/*
Logs a user out of the application. If cookies are set, they will be  set to be expired

Route: /api/logout

Authorization: Required

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) LogoutHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie(cookie_name); err == nil {
			value := make(map[string]string)
			if err = s.Decode(cookie_name, cookie.Value, &value); err == nil {
				// removes the cookie
				cookie := &http.Cookie{
					Name:   cookie_name,
					Value:  "",
					Secure: false,
					Path:   "/",
					MaxAge: -1,
				}
				http.SetCookie(c.Writer, cookie)
			}
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Preregisters a new user. A token is sent to the email that the user provides
Which is then used in the registration.

Route: /api/preregister

Authorization: None

Response-Type: application/json

Possible HTTP Codes:  200, 400, 401, 409, 500
*/
func (dbconfig *DbConfig) PreRegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := utils.LoadEnv("email")
		email_psw := utils.LoadEnv("email_psw")
		if email == "" || email_psw == "" {
			RespondWithError(c, http.StatusInternalServerError, "invalid email or email password")
			return
		}
		request_body, err := DecodePreRegisterRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = PreRegisterValidation(request_body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		// user validation
		exists, err := dbconfig.CheckUserEmailType(email, "app")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if exists {
			RespondWithError(c, http.StatusConflict, "email '"+email+"' already exists")
			return
		}
		token_value := uuid.UUID{}
		token_value, exists, err = dbconfig.CheckTokenExists(email, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			token, err := dbconfig.DB.CreateToken(c.Request.Context(), database.CreateTokenParams{
				ID:        uuid.New(),
				Name:      request_body.Name,
				Email:     request_body.Email,
				Token:     uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			token_value = token.Token
		}
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		err = SendEmail(token_value, request_body.Email, request_body.Name)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, objects.ResponseString{
			Message: "email sent",
		})

	}
}

/*
Registers a new user. It expects an email and a token to be passed
into the body of the request. The token will be verified to confirm if it exists internally

Route: /api/register

Authorization: None

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 409, 500
*/
func (dbconfig *DbConfig) RegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := DecodeUserRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = ValidateTokenValidation(body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := dbconfig.DB.GetTokenValidation(c.Request.Context(), body.Email)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		request_token, err := uuid.Parse(body.Token)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode token: "+body.Token)
			return
		}
		if token.Token != request_token {
			RespondWithError(c, http.StatusNotFound, "invalid token for user")
			return
		}
		err = UserValidation(body.Name, body.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		exists, err := dbconfig.CheckUserExist(body.Name, c.Request)
		if exists {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		user, err := dbconfig.DB.CreateUser(c.Request.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      body.Name,
			UserType:  "app",
			Email:     body.Email,
			Password:  body.Password,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, ConvertDatabaseToRegister(user))
	}
}

/*
Confirms if the API is ready to start accepting requests.

Route: /api/ready

Authorization: None

Response-Type: application/json

Possible HTTP Codes: 200, 503
*/
func (dbconfig *DbConfig) ReadyHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if dbconfig.Valid {
			RespondWithJSON(c, http.StatusOK, gin.H{"message": "OK"})
		} else {
			RespondWithError(c, http.StatusServiceUnavailable, "Unavailable")
		}
	}
}

// Helper function
// logs all error messages in current context to stdout
// the error message in the parameters is returned over the API
// after the chain has been aborted.
func RespondWithError(c *gin.Context, http_code int, err_message string) {
	for _, err := range c.Errors {
		// TODO log previous errors from the authentication middlewares inside database table
		log.Println(err.Err.Error())
		break
	}
	c.AbortWithStatusJSON(http_code, gin.H{
		"message": err_message,
	})
}

// Helper function
// responds with a payload and http code over the API
// after sucessfully processing the request.
func RespondWithJSON(c *gin.Context, http_code int, payload any) {
	c.JSON(http_code, payload)
}

// Helper function
// appends the error to the current context and passes on to the next Middleware
// only used in the Authentication phase when one middleware auth is not met
// it uses another in the format below
// query_param -> api_key as header -> basic
// in the case that the last auth fails, it raises and error with RespondWithError
func AppendErrorNext(c *gin.Context, http_code int, err_message string) {
	c.Error(errors.New(err_message))
	c.Next()
}
