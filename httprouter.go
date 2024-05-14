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
	"ngrok"
	"objects"
	"os"
	"shopify"
	"strconv"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Creates a new webhook on Shopify

Route: /api/shopify/webhook

Authorization: Basic, QueryParams, Headers

Header: Optional Mocker Header can be sent with request, used just for tests

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddWebhookHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockRequest := c.Request.Header.Get("Mocker")
		if mockRequest == "true" {
			// if the request is a mock request
			// then we will not update the database
			// as this will overwrite
			// the users data (if in production)
			RespondWithJSON(c, http.StatusOK, objects.ResponseString{
				Message: "success",
			})
			return
		}
		// request will take care of the post and put in a single request

		ngrok_tunnels, err := ngrok.FetchNgrokTunnels()
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		domain_url := ngrok.FetchWebsiteTunnel(ngrok_tunnels)
		if domain_url == "" {
			RespondWithError(c, http.StatusInternalServerError, "could not locate ngrok tunnel")
			return
		}
		// checks if there is an internal record of a webhook URL inside the database
		// by default there should always only be one
		db_shopify_webhook, err := dbconfig.DB.GetShopifyWebhooks(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		shopifyConfig := shopify.InitConfigShopify("")
		dbUser, err := dbconfig.DB.GetUserByApiKey(c.Request.Context(), c.GetString("api_key"))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// create new webhook on Shopify
		if db_shopify_webhook.ShopifyWebhookID == "" {
			webhook_response, err := shopifyConfig.CreateShopifyWebhook(
				ngrok.SetUpWebhookURL(
					domain_url,
					dbUser.ApiKey,
					dbUser.WebhookToken,
				),
			)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			err = dbconfig.DB.UpdateShopifyWebhook(c.Request.Context(), database.UpdateShopifyWebhookParams{
				ShopifyWebhookID: fmt.Sprint(webhook_response.ID),
				WebhookUrl:       webhook_response.Address,
				Topic:            webhook_response.Topic,
				ID:               db_shopify_webhook.ID,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusCreated, objects.ResponseString{
				Message: "success",
			})
			return
		} else {
			// update existing webhook on Shopify
			webhook_response, err := shopifyConfig.UpdateShopifyWebhook(
				db_shopify_webhook.ShopifyWebhookID,
				ngrok.SetUpWebhookURL(
					domain_url,
					dbUser.ApiKey,
					dbUser.WebhookToken,
				),
			)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			err = dbconfig.DB.UpdateShopifyWebhook(c.Request.Context(), database.UpdateShopifyWebhookParams{
				ShopifyWebhookID: fmt.Sprint(webhook_response.ID),
				WebhookUrl:       webhook_response.Address,
				Topic:            webhook_response.Topic,
				ID:               db_shopify_webhook.ID,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, objects.ResponseString{
				Message: "success",
			})
		}
	}
}

/*
Deletes the webhook on Shopify

Route: /api/shopify/webhook

Authorization: Basic, QueryParams, Headers

Header: Optional Mocker Header can be sent with request, used just for tests

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) DeleteWebhookHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockRequest := c.Request.Header.Get("Mocker")
		if mockRequest == "true" {
			// if the request is a mock request
			// then we will not update the database
			// as this will overwrite
			// the users data (if in production)
			RespondWithJSON(c, http.StatusOK, objects.ResponseString{
				Message: "success",
			})
			return
		}
		// fetches data of internal shopify webhook and confirms if it's set
		db_shopify_webhook, err := dbconfig.DB.GetShopifyWebhooks(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// delete the webhook on Shopify
		if db_shopify_webhook.ShopifyWebhookID != "" {
			shopifyConfig := shopify.InitConfigShopify("")
			_, err := shopifyConfig.DeleteShopifyWebhook(db_shopify_webhook.ShopifyWebhookID)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			err = dbconfig.DB.UpdateShopifyWebhook(c.Request.Context(), database.UpdateShopifyWebhookParams{
				ShopifyWebhookID: "",
				WebhookUrl:       "",
				Topic:            "",
				ID:               db_shopify_webhook.ID,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, objects.ResponseString{
				Message: "success",
			})
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Updates the push restriction.

Route: /api/push/restriction

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PushRestrictionHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		restrictions, err := DecodeRestriction(dbconfig, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = RestrictionValidation(restrictions)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = UpdatePushRestriction(dbconfig, restrictions)
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
Returns all push restriction values

Route: /api/push/restriction

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetPushRestrictionHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		restrictions, err := dbconfig.DB.GetPushRestriction(context.Background())
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusInternalServerError, "no push restrictions found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, restrictions)
	}
}

/*
Returns all fetch restriction values

Route: /api/fetch/restriction

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetFetchRestrictionHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		restrictions, err := dbconfig.DB.GetFetchRestriction(context.Background())
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusInternalServerError, "no fetch restrictions found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, restrictions)
	}
}

/*
Updates the fetch restriction.

Route: /api/fetch/restriction

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) FetchRestrictionHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		restrictions, err := DecodeRestriction(dbconfig, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = RestrictionValidation(restrictions)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = UpdateFetchRestriction(dbconfig, restrictions)
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
Runs the worker that fetch products from shopify.

Route: /api/worker/fetch

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) WorkerFetchProductsHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create database table containing the status of this
		shopifyConfig := shopify.InitConfigShopify("")
		err := FetchShopifyProducts(dbconfig, shopifyConfig)
		if err != nil {
			if err.Error() == "worker is currently running" {
				RespondWithError(c, http.StatusConflict, err.Error())
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Resets the fetch worker of the application.

Route: /api/shopify/reset_fetch

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) ResetShopifyFetchHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := dbconfig.DB.ResetFetchWorker(context.Background(), "0")
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
Adds a new internal Warehouse to the application. Comes with an additional reindex param
that will set missing warehouses for older products.

Route: /api/inventory/warehouse?reindex=false

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddInventoryWarehouseHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		reindex := c.Query("reindex")
		warehouse, err := DecodeGlobalWarehouse(dbconfig, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = GlobalWarehouseValidation(warehouse)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = nil
		httpStatus := 200
		if reindex == "true" {
			httpStatus, err = AddGlobalWarehouse(dbconfig, c.Request.Context(), warehouse.Name, true)
		} else {
			httpStatus, err = AddGlobalWarehouse(dbconfig, c.Request.Context(), warehouse.Name, false)
		}
		if err != nil {
			RespondWithError(c, httpStatus, err.Error())
		}
		RespondWithJSON(c, httpStatus, objects.ResponseString{
			Message: "success",
		})
	}
}

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
		if err != nil || page < 0 {
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
			RespondWithError(c, http.StatusBadRequest, "could not decode warehouse id: "+warehouse_id)
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
			RespondWithError(c, http.StatusBadRequest, "could not decode warehouse id: "+warehouse_id)
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
		productID := c.Param("id")
		requestData, err := DecodeProductRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = UpdateProduct(dbconfig, requestData, productID, c.GetString("api_key"))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, requestData)
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
		if err != nil || page < 0 {
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

Possible HTTP Codes: 201, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddWarehouseLocationMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		location_map, err := DecodeInventoryMap(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := InventoryMapValidation(location_map); err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		result, err := AddWarehouseLocation(dbconfig, location_map)
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
		locations := objects.ShopifyLocations{}
		mockRequest := c.Request.Header.Get("Mocker")
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 0 {
			page = 1
		}
		if mockRequest == "true" {
			// if the request is a mock request
			// then we will not actually fetch data from Shopify
			locations = objects.ShopifyLocations{}
		} else {
			shopifyConfig := shopify.InitConfigShopify("")
			if !shopifyConfig.Valid {
				RespondWithError(c, http.StatusInternalServerError, "invalid shopify config")
				return
			}
			locations, err = shopifyConfig.GetShopifyLocations()
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
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

Possible HTTP Codes: 201, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostCustomerHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		customer_body, err := DecodeCustomerRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		if CustomerValidation(customer_body) != nil {
			RespondWithError(c, http.StatusBadRequest, "invalid customer first name")
			return
		}
		dbCustomer, err := AddCustomer(dbconfig, customer_body, customer_body.FirstName+" "+customer_body.LastName)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, dbCustomer)
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
		customer, err := CompileCustomerData(dbconfig, customer_uuid, false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "customer with ID '"+customer_id+"' not found")
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
		if err != nil || page < 0 {
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
			cust, err := CompileCustomerData(dbconfig, value.ID, true)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusNotFound, "customer with ID '"+value.ID.String()+"' not found")
					return
				}
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

Header: Optional Mocker Header can be sent with request, used just for tests

Response-Type: application/json

Possible HTTP Codes: 201, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostOrderHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockRequest := c.Request.Header.Get("Mocker")
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
				RespondWithError(c, http.StatusNotFound, "invalid token for user")
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
		orderID, err := CheckExistsOrder(dbconfig, c.Request.Context(), fmt.Sprint(order_body.Name))
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		queueRequest := objects.RequestQueueHelper{
			Type:        "order",
			Status:      "in-queue",
			Instruction: "update_order",
			Endpoint:    "queue",
			ApiKey:      api_key,
			Method:      http.MethodPost,
			Object:      order_body,
		}
		status := http.StatusOK
		if orderID == uuid.Nil {
			queueRequest.Instruction = "add_order"
			status = http.StatusCreated
		}
		if mockRequest != "true" {
			response_payload, err := dbconfig.QueueHelper(queueRequest)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, status, response_payload)
		}
		queueRequest.ApiKey = "***"
		RespondWithJSON(c, status, queueRequest)
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
		if search_query == "" || len(search_query) == 0 {
			RespondWithError(c, http.StatusBadRequest, "invalid search param")
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
		order_data, err := CompileOrderData(dbconfig, order_uuid, false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "order with ID '"+order_id+"' not found")
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
		if err != nil || page < 0 {
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
			ord, err := CompileOrderData(dbconfig, value.ID, true)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusNotFound, "order with ID '"+value.ID.String()+"' not found")
					return
				}
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
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusBadRequest, "not found")
				return
			}
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

Route: /api/products/{id}/variants/{variant_id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) RemoveProductVariantHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode variant id: "+product_id)
			return
		}
		variant_id := c.Param("variant_id")
		err = IDValidation(variant_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		variant_uuid, err := uuid.Parse(variant_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode variant id: "+variant_id)
			return
		}
		err = dbconfig.DB.RemoveVariant(c.Request.Context(), database.RemoveVariantParams{
			ID:        variant_uuid,
			ProductID: product_uuid,
		})
		if err != nil {
			// TODO whether it exists or not is something that the query can decide
			// or should we do a prior check to avoid unnecessary code from running?
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
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
func (dbconfig *DbConfig) ProductExportHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockRequest := c.Request.Header.Get("Mocker")
		if mockRequest == "true" {
			// if the request is a mock request
			// then we will not get the file as it
			// might be that we trying to open the file
			// on github servers which gives us the permission error
			RespondWithJSON(c, http.StatusOK, objects.ResponseString{
				Message: "success",
			})
			return
		}
		product_ids, err := dbconfig.DB.GetProductIDs(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		products := []objects.Product{}
		for _, product_id := range product_ids {
			product, err := CompileProduct(dbconfig, product_id, false)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusNotFound, "product with ID '"+product_id.String()+"' not found")
					return
				}
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
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1024*1024) // 10 Mb
		file_name, err := iocsv.UploadFile(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		CSVProducts, err := iocsv.ReadFile(wd + "/" + file_name)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		importingRecord := objects.ImportResponse{
			ProcessedCounter: 0,
			FailCounter:      0,
			ProductsAdded:    0,
			ProductsUpdated:  0,
			VariantsAdded:    0,
			VariantsUpdated:  0,
		}
		for _, CSVProduct := range CSVProducts {
			importingRecord = UpsertProduct(dbconfig, importingRecord, CSVProduct)
			importingRecord = UpsertImages(dbconfig, importingRecord, CSVProduct)
			importingRecord = UpsertVariant(dbconfig, importingRecord, CSVProduct)
			importingRecord = UpsertPrice(dbconfig, importingRecord, CSVProduct)
			importingRecord = UpsertWarehouse(dbconfig, importingRecord, CSVProduct)
			importingRecord.ProcessedCounter++
		}
		RespondWithJSON(c, http.StatusOK, importingRecord)
	}
}

/*
Creates and adds a new product to the application

Route: /api/products

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 201, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) PostProductHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestProduct objects.RequestBodyProduct
		err := c.Bind(&requestProduct)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		if requestProduct.Active == "" || len(requestProduct.Active) == 0 {
			requestProduct.Active = "0"
		}
		dbProductID, httpCode, err := AddProduct(dbconfig, requestProduct)
		if err != nil {
			RespondWithError(c, httpCode, err.Error())
			return
		}
		product_added, err := CompileProduct(dbconfig, dbProductID, false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "product with ID '"+dbProductID.String()+"' not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Only products that are active
		// should be pushed (automatically) to Shopify
		if product_added.Active == "1" {
			api_key := c.GetString("api_key")
			err = CompileInstructionProduct(dbconfig, product_added, api_key)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusNotFound, "shopify product ID not found for: '"+product_added.ProductCode+"'")
					return
				}
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			for _, variant := range product_added.Variants {
				err = CompileInstructionVariant(dbconfig, variant, product_added, api_key)
				if err != nil {
					if err.Error() == "sql: no rows in result set" {
						RespondWithError(c, http.StatusNotFound, "shopify variant ID not found for: '"+variant.Sku+"'")
						return
					}
					RespondWithError(c, http.StatusInternalServerError, err.Error())
					return
				}
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
		if err != nil || page < 0 {
			page = 1
		}
		query_param_type := utils.ConfirmFilters(c.Query("type"))
		query_param_category := utils.ConfirmFilters(c.Query("category"))
		query_param_vendor := utils.ConfirmFilters(c.Query("vendor"))
		response, err := CompileFilterSearch(
			dbconfig,
			false,
			page,
			query_param_type,
			query_param_category,
			query_param_vendor,
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
		compiled, err := CompileSearchResult(dbconfig, search)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "product ID not found")
			}
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
		if err != nil || page < 0 {
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
			prod, err := CompileProduct(dbconfig, value.ID, false)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusNotFound, "product ID not found for: '"+value.ID.String()+"'")
					return
				}
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

Route: /api/products/{id}

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 404, 401, 500
*/
func (dbconfig *DbConfig) ProductIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id '"+product_id+"'")
			return
		}
		product_data, err := CompileProduct(dbconfig, product_uuid, false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "product with ID '"+product_id+"' not found")
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

Authorization: None

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
		db_user, exists, err := CheckUserCredentials(dbconfig, body, c.Request)
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
Logs a user out of the application. If cookies are set, they will be set to be expired

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

Header: Optional Mocker Header can be sent with request, used just for tests

Response-Type: application/json

Possible HTTP Codes:  200, 400, 401, 409, 500
*/
func (dbconfig *DbConfig) PreRegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		mockRequest := c.Request.Header.Get("Mocker")
		email := utils.LoadEnv("email")
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
		exists, err := CheckUserEmailType(dbconfig, request_body.Email, "app")
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if exists {
			RespondWithError(c, http.StatusConflict, "email '"+request_body.Email+"' already exists")
			return
		}
		token_value, exists, err := CheckExistsToken(dbconfig, email, c.Request)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !exists {
			dbTokenDetails, err := AddUserRegistration(dbconfig, request_body)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			token_value = dbTokenDetails.Token
		}
		if mockRequest != "true" {
			err = Email(token_value, true, request_body.Email, request_body.Name)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
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
		requestUserData, err := DecodeUserRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		exists, err := CheckUExistsUser(dbconfig, requestUserData.Name, c.Request)
		if exists {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		err = ValidateTokenValidation(requestUserData)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := dbconfig.DB.GetTokenValidation(c.Request.Context(), requestUserData.Email)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		request_token, err := uuid.Parse(requestUserData.Token)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode token: "+requestUserData.Token)
			return
		}
		if token.Token != request_token {
			RespondWithError(c, http.StatusNotFound, "invalid token for user")
			return
		}
		err = UserValidation(requestUserData.Name, requestUserData.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		dbUser, err := AddUser(dbconfig, requestUserData)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, ConvertDatabaseToRegister(dbUser))
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
