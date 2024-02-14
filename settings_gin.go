package main

import (
	"context"
	"integrator/internal/database"
	"net/http"
	"objects"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Returns a list of the internal shopify settings. If the key query param is used, it returns the data for the specific setting.

Route: /api/shopify/settings

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetShopifySettingValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key != "" {
			setting_value, err := dbconfig.DB.GetShopifySettingByKey(context.Background(), key)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusInternalServerError, "no setting value found for "+key)
					return
				}
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, setting_value)
		} else {
			setting_value, err := dbconfig.DB.GetShopifySettings(context.Background())
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusInternalServerError, "no setting value found for "+key)
					return
				}
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, setting_value)
		}
	}
}

/*
Updates an existing shopify setting inside the database.

Route: /api/shopify/settings

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddShopifySetting() gin.HandlerFunc {
	return func(c *gin.Context) {
		setting_keys, err := dbconfig.DB.GetShopifySettingsList(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		shopify_settings_map, err := DecodeSettings(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = SettingsValidation(shopify_settings_map, setting_keys)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		for _, setting := range shopify_settings_map {
			err = dbconfig.DB.UpdateShopifySetting(c.Request.Context(), database.UpdateShopifySettingParams{
				Value:     setting.Value,
				UpdatedAt: time.Now().UTC(),
				Key:       setting.Key,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Updates an existing shopify setting inside the database.

Route: /api/shopify/settings

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) AddAppSetting() gin.HandlerFunc {
	return func(c *gin.Context) {
		setting_keys, err := dbconfig.DB.GetAppSettingsList(c.Request.Context())
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		app_settings_map, err := DecodeSettings(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = SettingsValidation(app_settings_map, setting_keys)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		for _, setting := range app_settings_map {
			err = dbconfig.DB.UpdateAppSetting(c.Request.Context(), database.UpdateAppSettingParams{
				Value:     setting.Value,
				UpdatedAt: time.Now().UTC(),
				Key:       setting.Key,
			})
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		RespondWithJSON(c, http.StatusOK, objects.ResponseString{
			Message: "success",
		})
	}
}

/*
Returns a list of the internal app settings. If the key query param is used, it returns the data for the specific setting.

Route: /api/settings

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

Possible HTTP Codes: 200, 400, 401, 404, 500
*/
func (dbconfig *DbConfig) GetAppSettingValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.ToLower(c.Query("key"))
		if key != "" {
			setting_value, err := dbconfig.DB.GetAppSettingByKey(context.Background(), key)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusInternalServerError, "no setting value found for "+key)
					return
				}
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, setting_value)
		} else {
			setting_value, err := dbconfig.DB.GetAppSettings(context.Background())
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					RespondWithError(c, http.StatusInternalServerError, "no setting value found for "+key)
					return
				}
				RespondWithError(c, http.StatusInternalServerError, err.Error())
				return
			}
			RespondWithJSON(c, http.StatusOK, setting_value)
		}
	}
}
