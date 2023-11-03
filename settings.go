package main

import (
	"context"
	"integrator/internal/database"
	"net/http"
	"strings"
	"time"
	"utils"

	"github.com/google/uuid"
)

// GET /api/settings
func (dbconfig *DbConfig) GetAppSettingValue(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	key := strings.ToLower(r.URL.Query().Get("key"))
	if key != "" {
		setting_value, err := dbconfig.DB.GetAppSettingByKey(context.Background(), key)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, setting_value)
	} else {
		setting_value, err := dbconfig.DB.GetAppSettings(context.Background())
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, setting_value)
	}
}

// POST /api/settings
func (dbconfig *DbConfig) AddAppSetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	setting_keys := utils.GetAppSettings("app")
	app_settings_map, err := DecodeSettings(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = SettingsValidation(app_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for _, setting := range app_settings_map {
		err = dbconfig.DB.AddAppSetting(r.Context(), database.AddAppSettingParams{
			ID:          uuid.New(),
			Key:         setting.Key,
			Description: setting_keys[strings.ToUpper(setting.Key)],
			Value:       setting.Value,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})
		if err != nil {
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				err = dbconfig.DB.UpdateAppSetting(r.Context(), database.UpdateAppSettingParams{
					Value:     setting.Value,
					UpdatedAt: time.Now().UTC(),
					Key:       setting.Key,
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
	RespondWithJSON(w, http.StatusOK, "success")
}

// DELETE /api/settings
func (dbconfig *DbConfig) RemoveAppSettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	setting_keys := utils.GetAppSettings("app")
	app_settings_map, err := DecodeSetting(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = SettingValidation(app_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = dbconfig.DB.RemoveAppSetting(r.Context(), app_settings_map.Key)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, "success")
}

// GET /api/shopify/settings
func (dbconfig *DbConfig) GetShopifySettingValue(
	w http.ResponseWriter,
	r *http.Request,
	user database.User) {
	key := r.URL.Query().Get("key")
	if key != "" {
		setting_value, err := dbconfig.DB.GetShopifySettingByKey(context.Background(), key)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, setting_value)
	} else {
		setting_value, err := dbconfig.DB.GetShopifySettings(context.Background())
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, setting_value)
	}
}

// POST /api/shopify/settings
func (dbconfig *DbConfig) AddShopifySetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	setting_keys := utils.GetAppSettings("shopify")
	shopify_settings_map, err := DecodeSettings(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = SettingsValidation(shopify_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for _, setting := range shopify_settings_map {
		err = dbconfig.DB.AddShopifySetting(r.Context(), database.AddShopifySettingParams{
			ID:          uuid.New(),
			Key:         setting.Key,
			Description: setting_keys[strings.ToUpper(setting.Key)],
			Value:       setting.Value,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})
		if err != nil {
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				err = dbconfig.DB.UpdateShopifySetting(r.Context(), database.UpdateShopifySettingParams{
					Value:     setting.Value,
					UpdatedAt: time.Now().UTC(),
					Key:       setting.Key,
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
	RespondWithJSON(w, http.StatusOK, "success")
}

// DELETE /api/shopify/settings
func (dbconfig *DbConfig) RemoveShopifySettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	setting_keys := utils.GetAppSettings("shopify")
	shopify_settings_map, err := DecodeSetting(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = SettingValidation(shopify_settings_map, setting_keys)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = dbconfig.DB.RemoveShopifySetting(r.Context(), shopify_settings_map.Key)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, "success")
}
