package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"log"
	"net/http"
	"objects"
	"shopify"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const MOCK_SHOPIFY_API_URL = "http://localhost:4711"
const MOCK_SHOPIFY_API_KEY = "9812Y3N13981UO1NWD"
const MOCK_SHOPIFY_API_PSW = "shpat_92UHEYF927YR2"
const MOCK_SHOPIFY_API_VERSION = "2021-07"
const MOCK_SHOPIFY_STORE_NAME = "test-test"

const MOCK_SHOPIFY_WEBHOOK_ID = "47593067"

const MOCK_NGROK_WEBHOOK_URL = "https://f5fa-102-135-246-72.ngrok-free.app"

const MOCK_SHOPIFY_LOCATION_ID = 10293810823
const MOCK_SHOPIFY_INVENTORY_LEVEL_ID = 23087120381
const MOCK_INVENTORY_ITEM_ID = 1023781023

const MOCK_SHOPIFY_PRODUCT_ID = 1072481085
const MOCK_SHOPIFY_VARIANT_ID = 1070325083
const MOCK_SHOPIFY_COLLECTION_ID = 2039482049
const MOCK_SHOPIFY_CUSTOM_COLLECTION_ID = 1063001407
const MOCK_PRODUCT_SKU = "product_sku"

const MOCK_SHOPIFY_LOCATION_MAP_UUID = "c266e9f6-1ca6-4e27-8dd8-cce2bf5fdba5"
const MOCK_SHOPIFY_WAREHOUSE_NAME = "MOCK-SHOPIFY-WAREHOUSE-NAME"

func TestShopifyVariantPricing(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - Invalid function params
	result, err := dbconfig.ShopifyVariantPricing(objects.ProductVariant{}, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, "0.00", result)

	// Test 2 - Price Tier not found, 0.00 returned
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	result, err = dbconfig.ShopifyVariantPricing(productPayload.Variants[0], "Price Tier 1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "0.00", result)

	// Test 3 - Valid price returned
	result, err = dbconfig.ShopifyVariantPricing(productPayload.Variants[0], "Selling Price")
	assert.Equal(t, nil, err)
	assert.Equal(t, "1500.99", result)
}

func TestCalculateAvailableQuantity(t *testing.T) {
	// Test 1 - Invalid function params
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)
	result := dbconfig.CalculateAvailableQuantity(&shopifyConfig, 0, "", "")
	assert.Equal(t, int32(0), result)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	CreateDatabaseShopifyInventory(&dbconfig)
	defer ClearShopifyInventoryData(&dbconfig)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/inventory_levels.json?location_ids="+
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+"&inventory_item_ids="+fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				CreateShopifyInventoryLevelsResponse("test-case-valid-shopify-inventory-levels.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test 2 - valid params - positive integer
	result = dbconfig.CalculateAvailableQuantity(
		&shopifyConfig,
		5,
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID),
		fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
	)
	assert.Equal(t, int32(13), result)

	// Test 3 - valid params - zero value
	result = dbconfig.CalculateAvailableQuantity(
		&shopifyConfig,
		0,
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID),
		fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
	)
	assert.Equal(t, int32(8), result)

	// Test 4 - valid params - negative number
	result = dbconfig.CalculateAvailableQuantity(
		&shopifyConfig,
		-2,
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID),
		fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
	)
	assert.Equal(t, int32(6), result)
}

func TestRemoveLocationMap(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	err := dbconfig.RemoveLocationMap("")
	assert.NotEqual(t, nil, err)

	// Test 2 - valid function params
	CreateDatabaseShopifyLocationMap(&dbconfig)
	defer ClearShopifyLocationData(&dbconfig)
	err = dbconfig.RemoveLocationMap(MOCK_SHOPIFY_LOCATION_MAP_UUID)
	assert.Equal(t, nil, err)
}

func TestPushProductInventory(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")

	// Test 1 - invalid function param
	err := dbconfig.PushProductInventory(&shopifyConfig, productPayload.Variants[0])
	assert.NotEqual(t, nil, err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	createDatabaseProduct(&dbconfig)
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/inventory_levels.json?location_ids="+
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+"&inventory_item_ids="+fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				CreateShopifyInventoryLevelsResponse("test-case-valid-shopify-inventory-levels.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/inventory_levels/connect.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				201,
				CreateShopifyInventoryItemConnectResponse("test-case-valid-inventory-item-connect.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/inventory_levels/adjust.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				CreateShopifyInventoryItemAdjustResponse("test-case-valid-level-item-adjust.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	// Test 2 - valid variant data | SKU not found
	productPayload.Variants[0].Sku = "MOCK-SKU-NOT-FOUND"
	err = dbconfig.PushProductInventory(&shopifyConfig, productPayload.Variants[0])
	assert.NotEqual(t, nil, err)

	// Test 3 - valid variant data | invalid warehouse name
	productPayload.Variants[0].Sku = "product_sku"
	productPayload.Variants[0].VariantQuantity = append(
		productPayload.Variants[0].VariantQuantity,
		objects.VariantQty{
			IsDefault: false,
			Name:      "MOCK-WAREHOUSE-NOT-FOUND",
			Value:     0,
		},
	)
	err = dbconfig.PushProductInventory(&shopifyConfig, productPayload.Variants[0])
	assert.NotEqual(t, nil, err)

	// Test 3 - valid variant data
	CreateDatabaseShopifyLocationMap(&dbconfig)
	defer ClearShopifyLocationData(&dbconfig)
	productPayload.Variants[0].VariantQuantity[0].Name = "TestHouse"
	err = dbconfig.PushProductInventory(&shopifyConfig, productPayload.Variants[0])
	assert.Equal(t, nil, err)
}

/* Returns a test shopify inventory item adjust response struct */
func CreateShopifyInventoryItemAdjustResponse(fileName string) objects.ResponseAddInventoryItem {
	fileBytes := payload("./test_payloads/tests/inventory-item-adjust/" + fileName)
	shopifyInventoryLevel := objects.ResponseAddInventoryItem{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a test shopify inventory item connect response struct */
func CreateShopifyInventoryItemConnectResponse(fileName string) objects.ResponseAddInventoryItemLocation {
	fileBytes := payload("./test_payloads/tests/inventory-item-connect/" + fileName)
	shopifyInventoryLevel := objects.ResponseAddInventoryItemLocation{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a test shopify inventory level response struct */
func CreateShopifyInventoryLevelsResponse(fileName string) objects.GetShopifyInventoryLevelsList {
	fileBytes := payload("./test_payloads/tests/shopify-inventory-level/" + fileName)
	shopifyInventoryLevel := objects.GetShopifyInventoryLevelsList{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a database.CreateShopifyInventoryRecordParams struct */
func CreateShopifyInventoryRecordDatabaseStruct(fileName string) database.CreateShopifyInventoryRecordParams {
	fileBytes := payload("./test_payloads/tests/shopify-inventory/" + fileName)
	shopifyInventoryLevel := database.CreateShopifyInventoryRecordParams{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a database.CreateShopifyLocationParams struct */
func CreateShopifyLocationRecordDatabaseStruct(fileName string) database.CreateShopifyLocationParams {
	fileBytes := payload("./test_payloads/tests/shopify-location/" + fileName)
	shopifyInventoryLevel := database.CreateShopifyLocationParams{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a database.CreateVIDParams struct */
func CreateCreateVIDStruct(fileName string) database.CreateVIDParams {
	fileBytes := payload("./test_payloads/tests/shopify-vid/" + fileName)
	shopifyVID := database.CreateVIDParams{}
	err := json.Unmarshal(fileBytes, &shopifyVID)
	if err != nil {
		log.Println(err)
	}
	return shopifyVID
}

/*
Creates an internal Shopify Inventory row in the database
*/
func CreateDatabaseShopifyInventory(dbconfig *DbConfig) {
	err := dbconfig.DB.CreateShopifyInventoryRecord(context.Background(),
		CreateShopifyInventoryRecordDatabaseStruct("test-case-valid-shopify-inventory.json"),
	)
	if err != nil {
		log.Println(err)
	}
}

/*
Creates an internal Shopify VID link row in the database
*/
func CreateDatabaseShopifyVID(dbconfig *DbConfig, variantID uuid.UUID) {
	databaseParams := CreateCreateVIDStruct("test-case-valid-shopify-vid.json")
	databaseParams.VariantID = variantID
	err := dbconfig.DB.CreateVID(context.Background(), databaseParams)
	if err != nil {
		log.Println(err)
	}
}

/*
Creates an internal Shopify Location Warehouse Map
*/
func CreateDatabaseShopifyLocationMap(dbconfig *DbConfig) {
	_, err := dbconfig.DB.CreateShopifyLocation(context.Background(),
		CreateShopifyLocationRecordDatabaseStruct("test-case-valid-shopify-location-map.json"),
	)
	if err != nil {
		log.Println(err)
	}
}

/* Clears Internal Shopify Inventory IDs */
func ClearShopifyInventoryData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveShopifyInventoryRecord(context.Background(), database.RemoveShopifyInventoryRecordParams{
		ShopifyLocationID: fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID),
		InventoryItemID:   fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
	})
	dbconfig.DB.RemoveShopifyVIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
}

/* Clears Internal Shopify Locations */
func ClearShopifyLocationData(dbconfig *DbConfig) {
	uUID, err := uuid.Parse(MOCK_SHOPIFY_LOCATION_MAP_UUID)
	if err != nil {
		log.Println(err)
	}
	err = dbconfig.DB.RemoveShopifyLocationMap(context.Background(), uUID)
	if err != nil {
		log.Println(err)
	}
}

/* Returns a product struct */
func InitMockProduct(fileName string) objects.Product {
	/* Returns a product request body struct */
	fileBytes := payload("./test_payloads/tests/products/" + fileName)
	productData := objects.Product{}
	err := json.Unmarshal(fileBytes, &productData)
	if err != nil {
		log.Println(err)
	}
	return productData
}
