package main

import (
	"context"
	"integrator/internal/database"
	"net/http"
	"objects"
	"time"

	"github.com/google/uuid"
)

// POST /api/shopify/settings
func (dbconfig *DbConfig) AddShopifySetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	// TODO make a better way of storing these keys (struct in objects?)
	setting_keys := []objects.ShopifySettings{
		{
			Key: "default_price_tier",
		},
		{
			Key: "enable_push",
		},
		{
			Key: "enable_dynamic_inventory_management",
		},
		{
			Key: "default_cost_price_tier",
		},
	}
	shopify_settings_map, err := DecodeShopifySettings(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = ShopifySettingsValidation(shopify_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for _, setting := range shopify_settings_map {
		err = dbconfig.DB.AddShopifySetting(r.Context(), database.AddShopifySettingParams{
			ID:        uuid.New(),
			Key:       setting.Key,
			Value:     setting.Value,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				err = dbconfig.DB.UpdateShopifySetting(r.Context(), database.UpdateShopifySettingParams{
					Value: setting.Value,
					Key:   setting.Key,
				})
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
			} else {
				RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	RespondWithJSON(w, http.StatusOK, []string{"success"})
}

// DELETE /api/shopify/settings
func (dbconfig *DbConfig) RemoveShopifySettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	// TODO make a better way of storing these keys
	// create global struct to use for these functions
	// DRY
	setting_keys := []objects.ShopifySettings{
		{
			Key: "default_price_tier",
		},
		{
			Key: "enable_push",
		},
		{
			Key: "enable_dynamic_inventory_management",
		},
		{
			Key: "default_cost_price_tier",
		},
	}
	shopify_settings_map, err := DecodeShopifySetting(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = ShopifySettingValidation(shopify_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = dbconfig.DB.RemoveShopifySetting(r.Context(), shopify_settings_map.Key)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, []string{"success"})
}

// Returns the value of the setting respective to the Key
func (dbconfig *DbConfig) GetSettingValue(key string) (string, error) {
	setting_value, err := dbconfig.DB.GetShopifySettingByKey(context.Background(), key)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", nil
		}
		return "", err
	}
	return setting_value.Value, nil
}
