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
					Key:   "product_type",
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
					Value: "value",
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
			AcceptsData:   true,
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
		"DELETE /api/inventory/{id}": {
			Description:   "Removes a location-warehouse map",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        []string{},
			Authorization: "Authorization: ApiKey <key>",
		},
	}
	return routes
}
