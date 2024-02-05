package main

// import (
// 	"context"
// 	"integrator/internal/database"
// 	"net/http"
// 	"objects"
// 	"strings"
// 	"time"
// 	"utils"
// )

// // GET /api/settings
// func (dbconfig *DbConfig) GetAppSettingValue(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	key := strings.ToLower(r.URL.Query().Get("key"))
// 	if key != "" {
// 		setting_value, err := dbconfig.DB.GetAppSettingByKey(context.Background(), key)
// 		if err != nil {
// 			if err.Error() == "sql: no rows in result set" {
// 				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
// 				return
// 			}
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		RespondWithJSON(w, http.StatusOK, setting_value)
// 	} else {
// 		setting_value, err := dbconfig.DB.GetAppSettings(context.Background())
// 		if err != nil {
// 			if err.Error() == "sql: no rows in result set" {
// 				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
// 				return
// 			}
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		RespondWithJSON(w, http.StatusOK, setting_value)
// 	}
// }

// // PUT /api/settings
// func (dbconfig *DbConfig) AddAppSetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
// 	setting_keys := utils.GetAppSettings("app")
// 	app_settings_map, err := DecodeSettings(r)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = SettingsValidation(app_settings_map, setting_keys)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	for _, setting := range app_settings_map {
// 		err = dbconfig.DB.UpdateAppSetting(r.Context(), database.UpdateAppSettingParams{
// 			Value:     setting.Value,
// 			UpdatedAt: time.Now().UTC(),
// 			Key:       setting.Key,
// 		})
// 		if err != nil {
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "success",
// 	})
// }

// // DELETE /api/settings
// func (dbconfig *DbConfig) RemoveAppSettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {
// 	setting_keys := utils.GetAppSettings("app")
// 	key := r.URL.Query().Get("key")
// 	err := SettingValidation(
// 		objects.RequestSettings{
// 			Key:   key,
// 			Value: "",
// 		},
// 		setting_keys,
// 	)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = dbconfig.DB.RemoveAppSetting(r.Context(), key)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "success",
// 	})
// }

// // GET /api/shopify/settings
// func (dbconfig *DbConfig) GetShopifySettingValue(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	user database.User) {
// 	key := r.URL.Query().Get("key")
// 	if key != "" {
// 		setting_value, err := dbconfig.DB.GetShopifySettingByKey(context.Background(), key)
// 		if err != nil {
// 			if err.Error() == "sql: no rows in result set" {
// 				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
// 				return
// 			}
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		RespondWithJSON(w, http.StatusOK, setting_value)
// 	} else {
// 		setting_value, err := dbconfig.DB.GetShopifySettings(context.Background())
// 		if err != nil {
// 			if err.Error() == "sql: no rows in result set" {
// 				RespondWithError(w, http.StatusInternalServerError, "no setting value found for "+key)
// 				return
// 			}
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		RespondWithJSON(w, http.StatusOK, setting_value)
// 	}
// }

// // PUT /api/shopify/settings
// func (dbconfig *DbConfig) AddShopifySetting(w http.ResponseWriter, r *http.Request, dbUser database.User) {
// 	setting_keys := utils.GetAppSettings("shopify")
// 	shopify_settings_map, err := DecodeSettings(r)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = SettingsValidation(shopify_settings_map, setting_keys)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	for _, setting := range shopify_settings_map {
// 		err = dbconfig.DB.UpdateShopifySetting(r.Context(), database.UpdateShopifySettingParams{
// 			Value:     setting.Value,
// 			UpdatedAt: time.Now().UTC(),
// 			Key:       setting.Key,
// 		})
// 		if err != nil {
// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "success",
// 	})
// }

// // DELETE /api/shopify/settings
// func (dbconfig *DbConfig) RemoveShopifySettings(w http.ResponseWriter, r *http.Request, dbUser database.User) {
// 	setting_keys := utils.GetAppSettings("shopify")
// 	key := r.URL.Query().Get("key")
// 	err := SettingValidation(
// 		objects.RequestSettings{
// 			Key:   key,
// 			Value: "",
// 		},
// 		setting_keys,
// 	)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	err = dbconfig.DB.RemoveShopifySetting(r.Context(), key)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	RespondWithJSON(w, http.StatusOK, objects.ResponseString{
// 		Message: "success",
// 	})
// }
