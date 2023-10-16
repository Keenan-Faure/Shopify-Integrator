package main

import (
	"integrator/internal/database"
	"net/http"
	"objects"
)

// TODO make a better way of storing these keys (struct in objects?)

// POST /api/shopify/settings
func (dbconfig *DbConfig) AddShopifySetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	setting_keys := []objects.ShopifySettings{
		{
			Key: "default_price_tier",
		},
		{
			Key: "enable_push",
		},
	}
	for _, value := range setting_keys {
		if(value.Key == )
	}
}

// DELETE /api/shopify/ettings
func (dbconfig *DbConfig) RemoveShopifySettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {

}