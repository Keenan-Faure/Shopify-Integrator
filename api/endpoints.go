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
		"GET /api/ready": {
			Description:   "Returns the status of the API",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
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
			Authorization: "Authorization: ApiKey <key>",
		},
	}
	return routes
}
