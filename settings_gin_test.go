package main

import (
	"context"
	"encoding/json"
	"integrator/internal/database"
	"net/http"
	"objects"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppSettingValue(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/settings",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - valid request */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	createDatabaseAppSettings(&dbconfig)
	w = Init(
		"/api/settings?api_key="+dbUser.ApiKey,
		http.MethodGet, make(map[string][]string), nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	responseAppSettings := []database.GetAppSettingsRow{}
	err := json.Unmarshal(w.Body.Bytes(), &responseAppSettings)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	assert.Equal(t, "app_enable_shopify_push", responseAppSettings[0].Key)
	assert.Equal(t, "Enables products to be pushed to Shopify.", responseAppSettings[0].Description)
	assert.Equal(t, "Enable Shopify Push", responseAppSettings[0].FieldName)
	assert.Equal(t, "false", responseAppSettings[0].Value)

	assert.Equal(t, "app_queue_cron_time", responseAppSettings[3].Key)
	assert.Equal(t, "Interval between each run of the queue worker.", responseAppSettings[3].Description)
	assert.Equal(t, "Queue Cron Time", responseAppSettings[3].FieldName)
	assert.Equal(t, "5", responseAppSettings[3].Value)

	assert.Equal(t, "app_fetch_sync_images", responseAppSettings[6].Key)
	assert.Equal(t, "Enabled products to be pulled from Shopify when fetching data.", responseAppSettings[6].Description)
	assert.Equal(t, "Add Shopify Images", responseAppSettings[6].FieldName)
	assert.Equal(t, "false", responseAppSettings[6].Value)
}

func TestAddAppSetting(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	appSettingsPayload := AppSettingsPayload("test-case-valid-request.json")
	w := Init("/api/settings", http.MethodPut, map[string][]string{}, appSettingsPayload, &dbconfig, router)

	assert.Equal(t, 401, w.Code)

	/* Test Case 2 - invalid request body | invalid fields */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	appSettingsPayload = AppSettingsPayload("test-case-invalid-request.json")
	w = Init(
		"/api/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, appSettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "setting 'app_queue_process_limiters' not allowed", response.Message)

	/* Test Case 3 - invalid request body | blank fields */
	appSettingsPayload = AppSettingsPayload("test-case-invalid-request-blank-field.json")
	w = Init(
		"/api/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, appSettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "settings key cannot be blank", response.Message)

	/* Test Case 4 - valid request */
	appSettingsPayload = AppSettingsPayload("test-case-valid-request.json")
	w = Init(
		"/api/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, appSettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
}

func TestAddShopifySetting(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	shopifySettingsPayload := ShopifySettingsPayload("test-case-valid-request.json")
	w := Init("/api/shopify/settings", http.MethodPut, map[string][]string{}, shopifySettingsPayload, &dbconfig, router)

	assert.Equal(t, 401, w.Code)

	/* Test Case 2 - invalid request body | invalid fields */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	shopifySettingsPayload = ShopifySettingsPayload("test-case-invalid-request.json")
	w = Init(
		"/api/shopify/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, shopifySettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "setting 'abc123_test' not allowed", response.Message)

	/* Test Case 3 - invalid request body | blank fields */
	shopifySettingsPayload = ShopifySettingsPayload("test-case-invalid-request-blank-field.json")
	w = Init(
		"/api/shopify/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, shopifySettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "settings key cannot be blank", response.Message)

	/* Test Case 4 - valid request */
	shopifySettingsPayload = ShopifySettingsPayload("test-case-valid-request.json")
	w = Init(
		"/api/shopify/settings?api_key="+dbUser.ApiKey,
		http.MethodPut, map[string][]string{}, shopifySettingsPayload, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
}

func TestGetShopifySettingValue(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/shopify/settings",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - valid request */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	createDatabaseShopifySettings(&dbconfig)
	w = Init(
		"/api/shopify/settings?api_key="+dbUser.ApiKey,
		http.MethodGet, make(map[string][]string), nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	responseShopifySettings := []database.GetShopifySettingsRow{}
	err := json.Unmarshal(w.Body.Bytes(), &responseShopifySettings)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "shopify_enable_dynamic_sku_search", responseShopifySettings[0].Key)
	assert.Equal(t, "Enables the dynamic searching of SKUs on Shopify when adding new products. If disabled, only first product SKU will be considered.", responseShopifySettings[0].Description)
	assert.Equal(t, "Dynamic SKU Search", responseShopifySettings[0].FieldName)
	assert.Equal(t, "true", responseShopifySettings[0].Value)
}
