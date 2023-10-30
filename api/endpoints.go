package api

import (
	"objects"
	"time"
)

// GET /api/endpoints
func Endpoints() objects.Endpoints {
	return objects.Endpoints{
		Status:      true,
		Description: "Integrator-Shopify API Documentation",
		Routes:      createRoutes(),
		Version:     "v1",
		Time:        time.Now().UTC(),
	}
}

func createRoutes() map[string]objects.Route {
	routes := map[string]objects.Route{
		"GET /api/endpoints": {
			Description:   "Returns the available endpoints",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"GET /api/ready": {
			Description:   "Returns the status of the API",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"GET /api/products": {
			Description: "Returns a list of products (uses pagination)",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"page": {
					Key:   "page",
					Value: "number",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/products/{id}": {
			Description:   "Returns the specific product",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/products/search": {
			Description: "Product search on specific value",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"query": {
					Key:   "q",
					Value: "value",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/products/filter": {
			Description: "Product search on specific value",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"data": {
					Key:   "product_type | vendor | category",
					Value: "value",
				},
				"page": {
					Key:   "page",
					Value: "number",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/orders": {
			Description: "Returns a list of orders (uses pagination)",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"page": {
					Key:   "page",
					Value: "number",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/orders/{id}": {
			Description:   "Returns the specific order",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/orders/search": {
			Description: "Returns the specific order",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"query": {
					Key:   "q",
					Value: "value",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/customers": {
			Description: "Returns a list of customers (uses pagination)",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"page": {
					Key:   "page",
					Value: "number",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/customers/{id}": {
			Description:   "Returns the specific order",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/customers/search": {
			Description: "Returns the specific order",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"query": {
					Key:   "q",
					Value: "query value",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/shopify/settings": {
			Description: "Returns the value of the shopify setting, or the entire list if the key is blank",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"key": {
					Key:   "key",
					Value: "shopify setting key",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/queue": {
			Description: "Returns the current items in the queue",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"page": {
					Key:   "page",
					Value: "integer number",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/queue/filter": {
			Description: "Filter searches through the queue to return specific results",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"page": {
					Key:   "page",
					Value: "integer number",
				},
				"type": {
					Key:   "type",
					Value: "queue item type",
				},
				"instruction": {
					Key:   "instruction",
					Value: "queue item instruction",
				},
				"status": {
					Key:   "status",
					Value: "queue item status",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/queue/view": {
			Description:   "Returns the current queue count in a structured map",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"GET /api/settings": {
			Description: "Returns the value of the app setting, or all of them if the key is blank",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"key": {
					Key:   "key",
					Value: "app setting key",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/login": {
			Description:   "Login with a new user",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/register": {
			Description:   "Registers a new user",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyUser{},
			Authorization: "None",
		},
		"POST /api/validatetoken": {
			Description:   "Validates a token (user registration)",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyValidateToken{},
			Authorization: "None",
		},
		"POST /api/preregister": {
			Description:   "Validates a token (user registration)",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyPreRegister{},
			Authorization: "None",
		},
		"POST /api/products": {
			Description:   "Adds a product",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyProduct{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/orders": {
			Description: "Adds an order",
			Supports:    []string{"POST"},
			Params: map[string]objects.Params{
				"token": {
					Key:   "token",
					Value: "token_value",
				},
				"api_key": {
					Key:   "api_key",
					Value: "key",
				},
			},
			AcceptsData:   true,
			Format:        objects.RequestBodyOrder{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/customers": {
			Description:   "Adds a customer",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyCustomer{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/inventory": {
			Description:   "Adds location-warehouse map",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestWarehouseLocation{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/shopify/settings": {
			Description:   "Creates a new shopify setting",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestSettings{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/settings": {
			Description:   "Creates a new app setting",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestSettings{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/queue": {
			Description:   "Adds a new item to the queue",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestQueueItem{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"POST /api/queue/worker": {
			Description:   "Pops and processes the next queue item",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"DELETE /api/inventory/{id}": {
			Description:   "Removes a location-warehouse map",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"DELETE /api/shopify/settings": {
			Description:   "Removes a shopify setting",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestSettings{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"DELETE /api/settings": {
			Description:   "Removes an app setting",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestSettings{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"DELETE /api/queue/{id}": {
			Description:   "Removes a queue item",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
		"DELETE /api/queue": {
			Description: "Removes a queue item by filters",
			Supports:    []string{"DELETE"},
			Params: map[string]objects.Params{
				"type": {
					Key:   "type",
					Value: "queue item type",
				},
				"instruction": {
					Key:   "instruction",
					Value: "queue item instruction",
				},
				"status": {
					Key:   "status",
					Value: "queue item status",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
	}
	return routes
}
