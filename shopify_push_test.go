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

const MOCK_APP_API_URL = "http://localhost:8080"

const MOCK_QUEUE_ITEM_ID = "66608bb9-6bef-4424-b48f-59f487ec2933"

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	CreateDatabaseShopifyInventory(&dbconfig)
	defer ClearShopifyInventoryData(&dbconfig)

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	createDatabaseProduct(&dbconfig)
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

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

func TestPushProduct(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")

	// Test 1 - invalid function param
	err := dbconfig.PushProduct(&shopifyConfig, objects.Product{})
	assert.NotEqual(t, nil, err)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 2 - valid data | valid productID not empty
	CreateDatabaseShopifyLocationMap(&dbconfig)
	defer ClearShopifyLocationData(&dbconfig)
	err = dbconfig.PushProduct(&shopifyConfig, productPayload)
	assert.Equal(t, nil, err)

	// Test 3 - valid data | productID empty
	dbconfig.DB.RemovePIDByProductCode(context.Background(), MOCK_PRODUCT_CODE)
	err = dbconfig.PushProduct(&shopifyConfig, productPayload)
	assert.Equal(t, nil, err)
}

func TestPushAddShopify(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)

	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	ids, _ := shopifyConfig.GetProductBySKU(MOCK_PRODUCT_SKU)
	shopifyProduct := ConvertProductToShopify(productPayload)

	restrictions, _ := dbconfig.DB.GetPushRestriction(context.Background())
	push_restrictions := PushRestrictionsToMap(restrictions)

	updateShopifyProduct := ApplyPushRestrictionProduct(push_restrictions, shopifyProduct)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test 1 - invalid function param
	err := PushAddShopify(
		&shopifyConfig,
		&dbconfig,
		objects.ResponseIDs{},
		objects.Product{},
		objects.ShopifyProduct{},
		objects.ShopifyProduct{},
	)
	assert.NotEqual(t, nil, err)

	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	CreateDatabaseShopifyLocationMap(&dbconfig)
	defer ClearShopifyLocationData(&dbconfig)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 2 - valid data | valid productID not empty
	err = PushAddShopify(
		&shopifyConfig,
		&dbconfig,
		objects.ResponseIDs{},
		productPayload,
		shopifyProduct,
		updateShopifyProduct,
	)
	assert.Equal(t, nil, err)

	// Test 3 - valid data
	dbconfig.DB.RemovePIDByProductCode(context.Background(), MOCK_PRODUCT_CODE)
	err = PushAddShopify(
		&shopifyConfig,
		&dbconfig,
		ids,
		productPayload,
		shopifyProduct,
		updateShopifyProduct,
	)
	assert.Equal(t, nil, err)
}

func TestCollectionShopfy(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)

	productPayload := InitMockProduct("test-case-valid-product-variable.json")

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test 1 - invalid function params
	err := dbconfig.CollectionShopfy(&shopifyConfig, objects.Product{}, 0)
	assert.NotEqual(t, nil, err)

	// Test 2 - valid product | invalid shopify product id
	err = dbconfig.CollectionShopfy(&shopifyConfig, productPayload, 0)
	assert.NotEqual(t, nil, err)

	// Test 3 - valid product | invalid product category
	productPayload.Category = ""
	err = dbconfig.CollectionShopfy(&shopifyConfig, productPayload, MOCK_SHOPIFY_PRODUCT_ID)
	assert.NotEqual(t, nil, err)

	// Test 4 - valid function params
	productPayload.Category = "product_category"
	err = dbconfig.CollectionShopfy(&shopifyConfig, productPayload, MOCK_SHOPIFY_PRODUCT_ID)
	assert.Equal(t, nil, err)
}

func TestPushVariant(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify(MOCK_SHOPIFY_API_URL)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")

	shopifyVariant := ConvertVariantToShopify(productPayload.Variants[0])

	restrictions, _ := dbconfig.DB.GetPushRestriction(context.Background())
	push_restrictions := PushRestrictionsToMap(restrictions)

	// Test 1 - invalid function param
	err := dbconfig.PushVariant(
		&shopifyConfig,
		objects.ProductVariant{},
		objects.ShopifyVariant{},
		map[string]string{},
		"",
		"",
	)
	assert.NotEqual(t, nil, err)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 2 - valid data | invalid shopify variant ID
	CreateDatabaseShopifyLocationMap(&dbconfig)
	defer ClearShopifyLocationData(&dbconfig)
	err = dbconfig.PushVariant(
		&shopifyConfig,
		productPayload.Variants[0],
		shopifyVariant,
		push_restrictions,
		fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID),
		fmt.Sprint(0),
	)
	assert.NotEqual(t, nil, err)

	// Test 3 - invalid data | invalid shopify product ID
	err = dbconfig.PushVariant(
		&shopifyConfig,
		productPayload.Variants[0],
		shopifyVariant,
		push_restrictions,
		fmt.Sprint(0),
		fmt.Sprint(MOCK_SHOPIFY_VARIANT_ID),
	)
	assert.NotEqual(t, nil, err)

	// Test 3 - valid data | productID empty
	dbconfig.DB.RemovePIDByProductCode(context.Background(), MOCK_PRODUCT_CODE)
	err = dbconfig.PushVariant(
		&shopifyConfig,
		productPayload.Variants[0],
		shopifyVariant,
		push_restrictions,
		fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID),
		fmt.Sprint(MOCK_SHOPIFY_VARIANT_ID),
	)
	assert.Equal(t, nil, err)
}

func TestCompileInstructionProduct(t *testing.T) {
	httpmock.Activate()
	InitMockQueue()
	defer httpmock.DeactivateAndReset()

	dbconfig := setupDatabase("", "", "", false)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	productUUID := createDatabaseProduct(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	// Test 1 - invalid params
	err := CompileInstructionProduct(&dbconfig, objects.Product{}, "")
	assert.NotEqual(t, nil, err)

	// Test 2 - valid params
	err = CompileInstructionProduct(&dbconfig, productPayload, dbUser.ApiKey)
	assert.Equal(t, nil, err)
}

func TestCompileInstructionVariant(t *testing.T) {
	httpmock.Activate()
	InitMockQueue()
	defer httpmock.DeactivateAndReset()

	dbconfig := setupDatabase("", "", "", false)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	// Test 1 - invalid params
	err := CompileInstructionVariant(&dbconfig, objects.ProductVariant{}, objects.Product{}, "")
	assert.NotEqual(t, nil, err)

	// Test 2 - valid params
	err = CompileInstructionVariant(&dbconfig, productPayload.Variants[0], productPayload, dbUser.ApiKey)
	assert.Equal(t, nil, err)
}

func TestGetShopifyProductID(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 1 - invalid param
	result, err := GetShopifyProductID(&dbconfig, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, "", result)

	// Test 2 - valid params
	result, err = GetShopifyProductID(&dbconfig, productPayload.ProductCode)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", result)
}

func TestGetShopifyVariantID(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 1 - invalid param
	result, err := GetShopifyVariantID(&dbconfig, "")
	assert.Equal(t, nil, err)
	assert.Equal(t, "", result)

	// Test 2 - valid params
	result, err = GetShopifyVariantID(&dbconfig, productPayload.Variants[0].Sku)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", result)
}

func TestSaveVariantIds(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	productUUID := createDatabaseProduct(&dbconfig)
	productPayload.ID = productUUID
	variantUUID, _ := dbconfig.DB.GetVariantIDBySKU(context.Background(), MOCK_PRODUCT_SKU)
	CreateDatabaseShopifyInventory(&dbconfig)
	CreateDatabaseShopifyPID(&dbconfig, productUUID)
	CreateDatabaseShopifyVID(&dbconfig, variantUUID)
	defer ClearShopifyInventoryData(&dbconfig)

	// Test 1 - invalid param
	err := SaveVariantIds(&dbconfig, objects.ShopifyProductResponse{}, productPayload.Variants[0].Sku)
	assert.Equal(t, nil, err)

	// Test 2 - valid params
	err = SaveVariantIds(
		&dbconfig,
		CreateShopifyProductResponse("test-case-valid-product.json"),
		productPayload.Variants[0].Sku,
	)
	assert.Equal(t, nil, err)
}

/* Returns a test shopify collection response struct */
func CreateShopifyCollectionResponse(fileName string) objects.ResponseAddProductToShopifyCollection {
	fileBytes := payload("./test_payloads/tests/shopify-collection/" + fileName)
	shopifyCollection := objects.ResponseAddProductToShopifyCollection{}
	err := json.Unmarshal(fileBytes, &shopifyCollection)
	if err != nil {
		log.Println(err)
	}
	return shopifyCollection
}

/* Returns a test shopify custom collection response struct */
func CreateShopifCustomCollectionResponse(fileName string) objects.ResponseShopifyCustomCollection {
	fileBytes := payload("./test_payloads/tests/shopify-custom-collections/" + fileName)
	shopifyCustomCollection := objects.ResponseShopifyCustomCollection{}
	err := json.Unmarshal(fileBytes, &shopifyCustomCollection)
	if err != nil {
		log.Println(err)
	}
	return shopifyCustomCollection
}

/* Returns a test shopify custom collections response struct */
func CreateShopifCollectionsResponse(fileName string) objects.ResponseGetCustomCollections {
	fileBytes := payload("./test_payloads/tests/shopify-collections/" + fileName)
	shopifyCustomCollections := objects.ResponseGetCustomCollections{}
	err := json.Unmarshal(fileBytes, &shopifyCustomCollections)
	if err != nil {
		log.Println(err)
	}
	return shopifyCustomCollections
}

/* Returns a test shopify variant response struct */
func CreateShopifVariantResponse(fileName string) objects.ShopifyVariantResponse {
	fileBytes := payload("./test_payloads/tests/shopify-variant/" + fileName)
	shopifyVariant := objects.ShopifyVariantResponse{}
	err := json.Unmarshal(fileBytes, &shopifyVariant)
	if err != nil {
		log.Println(err)
	}
	return shopifyVariant
}

/* Returns a test shopify location response struct */
func CreateShopifyLocationResponse(fileName string) objects.ShopifyLocations {
	fileBytes := payload("./test_payloads/tests/shopify-location/" + fileName)
	shopifyLocations := objects.ShopifyLocations{}
	err := json.Unmarshal(fileBytes, &shopifyLocations)
	if err != nil {
		log.Println(err)
	}
	return shopifyLocations
}

/* Returns a test shopify product response struct */
func CreateShopifyProductResponse(fileName string) objects.ShopifyProductResponse {
	fileBytes := payload("./test_payloads/tests/shopify-product/" + fileName)
	shopifyProduct := objects.ShopifyProductResponse{}
	err := json.Unmarshal(fileBytes, &shopifyProduct)
	if err != nil {
		log.Println(err)
	}
	return shopifyProduct
}

/* Returns a test shopify graph ql response response struct */
func CreateShopifyGraphQLResponse(fileName string) objects.JSONResponseShopifyGraphQL {
	fileBytes := payload("./test_payloads/tests/shopify-graph-ql/" + fileName)
	shopifyGraphQL := objects.JSONResponseShopifyGraphQL{}
	err := json.Unmarshal(fileBytes, &shopifyGraphQL)
	if err != nil {
		log.Println(err)
	}
	return shopifyGraphQL
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

/* Returns a test shopify webhook response struct */
func CreateShopifyWebhookResponse(fileName string) objects.ShopifyWebhookRequest {
	fileBytes := payload("./test_payloads/tests/shopify-webhook/" + fileName)
	shopifyWebhookResponse := objects.ShopifyWebhookRequest{}
	err := json.Unmarshal(fileBytes, &shopifyWebhookResponse)
	if err != nil {
		log.Println(err)
	}
	return shopifyWebhookResponse
}

/* Returns a test shopify product count struct */
func CreateShopifyProductCountResponse(fileName string) objects.ShopifyProductCount {
	fileBytes := payload("./test_payloads/tests/shopify-product-count/" + fileName)
	shopifyProductCount := objects.ShopifyProductCount{}
	err := json.Unmarshal(fileBytes, &shopifyProductCount)
	if err != nil {
		log.Println(err)
	}
	return shopifyProductCount
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

/* Returns a database.CreatePIDParams struct */
func CreateCreatePIDStruct(fileName string) database.CreatePIDParams {
	fileBytes := payload("./test_payloads/tests/shopify-pid/" + fileName)
	shopifyPID := database.CreatePIDParams{}
	err := json.Unmarshal(fileBytes, &shopifyPID)
	if err != nil {
		log.Println(err)
	}
	return shopifyPID
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
Creates an internal Shopify PID link row in the database
*/
func CreateDatabaseShopifyPID(dbconfig *DbConfig, productID uuid.UUID) {
	databaseParams := CreateCreatePIDStruct("test-case-valid-shopify-pid.json")
	databaseParams.ProductID = productID
	err := dbconfig.DB.CreatePID(context.Background(), databaseParams)
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
	dbconfig.DB.RemovePIDByProductCode(context.Background(), MOCK_PRODUCT_CODE)
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

/* Returns a mock queue_item response */
func QueueItemResponse(fileName string) objects.ResponseQueueItem {
	fileBytes := payload("./test_payloads/tests/queue-reply/" + fileName)
	queueItem := objects.ResponseQueueItem{}
	err := json.Unmarshal(fileBytes, &queueItem)
	if err != nil {
		log.Println(err)
	}
	return queueItem
}

/* Returns a mock queue_item response */
func QueueItemFilterResponse(fileName string) []objects.ResponseQueueItemFilter {
	fileBytes := payload("./test_payloads/tests/queue-filter/" + fileName)
	queueItem := []objects.ResponseQueueItemFilter{}
	err := json.Unmarshal(fileBytes, &queueItem)
	if err != nil {
		log.Println(err)
	}
	return queueItem
}

/* Sets up mock queue endpoints */
func InitMockQueue() {
	httpmock.RegisterResponder(http.MethodPost, MOCK_APP_API_URL+"/api/queue",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, QueueItemResponse("test-case-valid-product.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_APP_API_URL+"/api/shopify/sync",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, objects.ResponseString{
				Message: "synconizing started",
			})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/queue/"+MOCK_QUEUE_ITEM_ID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, QueueItemPayload("test-case-valid-product-queue-item.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/queue/filter",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, QueueItemFilterResponse("test-case-valid-queue-filter.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/queue/processing",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				QueueItemPayload("test-case-valid-product-queue-processing-item.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/queue",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, QueueItemPayload("test-case-valid-product-queue-items.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/queue/view",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, QueueItemPayload("test-case-valid-queue-count.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodDelete, MOCK_APP_API_URL+"/api/queue/"+MOCK_QUEUE_ITEM_ID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, objects.ResponseString{
				Message: "success",
			})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodDelete, MOCK_APP_API_URL+"/api/queue/filter",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, objects.ResponseString{
				Message: "success",
			})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

/* Sets up a mock shopify API */
func InitMockShopifyAPI() {

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/variants/"+fmt.Sprint(MOCK_SHOPIFY_VARIANT_ID)+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifVariantResponse("test-case-valid-variant.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/products/count.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyProductCountResponse("test-case-valid-product-count.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyWebhookResponse("test-case-valid-webhooks.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyWebhookResponse("test-case-valid-webhook.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodDelete, MOCK_SHOPIFY_API_URL+"/webhooks/"+MOCK_SHOPIFY_WEBHOOK_ID+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, "")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyWebhookResponse("test-case-valid-webhook.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/collects.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyCollectionResponse("test-case-valid-collection.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/custom_collections.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifCustomCollectionResponse("test-case-valid-custom-collection.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/custom_collections.json?fields=title,id",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifCollectionsResponse("test-case-valid-custom-collections.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet,
		MOCK_SHOPIFY_API_URL+"/custom_collections.json?fields=title,id&product_id="+fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifCollectionsResponse("test-case-valid-custom-collections.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/products.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyProductResponse("test-case-valid-product.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/locations.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyLocationResponse("test-case-valid-locations.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/products/"+fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID)+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyProductResponse("test-case-valid-product.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost,
		MOCK_SHOPIFY_API_URL+"/graphql.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyGraphQLResponse("test-case-valid-graph-ql.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

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
}
