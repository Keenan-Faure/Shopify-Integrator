package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"objects"
	"testing"

	"github.com/stretchr/testify/assert"
)

const MOCK_HTTP_REQUEST_URL = "http://mock.url:8212"
const MOCK_HTTP_REQUEST_METHOD = "POST"

func TestRestrictionValidation(t *testing.T) {
	// Test Case 1 - invalid push restrictions
	pushRestrictions := PushRestrictionPayload("test-case-invalid-request.json")
	valid := RestrictionValidation(pushRestrictions)
	assert.NotEqual(t, nil, valid)

	// Test Case 2 - valid push restrictions
	pushRestrictions = PushRestrictionPayload("test-case-valid-request.json")
	valid = RestrictionValidation(pushRestrictions)
	assert.Equal(t, nil, valid)

	// Test Case 3 - invalid push restrictions
	fetchRestrictions := FetchRestrictionPayload("test-case-invalid-request.json")
	valid = RestrictionValidation(fetchRestrictions)
	assert.NotEqual(t, nil, valid)

	// Test Case 4 - valid push restrictions
	fetchRestrictions = FetchRestrictionPayload("test-case-valid-request.json")
	valid = RestrictionValidation(fetchRestrictions)
	assert.Equal(t, nil, valid)
}

func TestDecodeRestriction(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := FetchRestrictionPayload("test-case-valid-request.json")
	invalidRequestBody := AppSettingsPayload("test-case-valid-request.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeRestriction(&dbconfig, request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(result), 0)

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeRestriction(&dbconfig, request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(result), 8)
	assert.Equal(t, result[0].Field, "")
	assert.Equal(t, result[0].Flag, "")
	assert.Equal(t, result[1].Field, "")
	assert.Equal(t, result[1].Field, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeRestriction(&dbconfig, request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result[0].Field, "title")
	assert.Equal(t, result[0].Flag, "shopify")
	assert.Equal(t, result[1].Field, "body_html")
	assert.Equal(t, result[1].Flag, "app")
}

func TestGlobalWarehouseValidation(t *testing.T) {
	// Test Case 1 - invalid push restrictions
	payload := WarehousePayload("test-case-invalid-warehouse.json")
	valid := GlobalWarehouseValidation(payload)
	assert.NotEqual(t, nil, valid)

	// Test Case 2 - valid push restrictions
	payload = WarehousePayload("test-case-valid-warehouse.json")
	valid = GlobalWarehouseValidation(payload)
	assert.Equal(t, nil, valid)
}

func TestDecodeGlobalWarehouse(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := WarehousePayload("test-case-valid-warehouse.json")
	invalidRequestBody := AppSettingsPayload("test-case-valid-request.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeGlobalWarehouse(&dbconfig, request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Name, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeGlobalWarehouse(&dbconfig, request)
	if err == nil {
		t.Errorf("Expected 'json: cannot unmarshal array into Go value of type objects.RequestGlobalWarehouse' but found: 'nil'")
	}
	assert.Equal(t, result.Name, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeGlobalWarehouse(&dbconfig, request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Name, "TestHouse")
}

func TestDecodeQueueItemOrder(t *testing.T) {
	requestBody := QueueItemPayload("test-case-valid-order-queue-item.json")
	invalidRequestBody := QueueItemPayload("test-case-invalid-order-queue-item.json")

	// Test 1 - empty request body
	result, err := DecodeQueueItemOrder(json.RawMessage{})
	if err == nil {
		t.Errorf("Expected 'unexpected end of JSON input' but found: 'nil'")
	}
	assert.Equal(t, result.Name, "")
	assert.Equal(t, int(result.ID), 0)

	// Test 2 - invalid request body
	result, err = DecodeQueueItemOrder(invalidRequestBody.Object)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Name, "")
	assert.Equal(t, int(result.ID), 0)

	// Test 3  - valid request body
	result, err = DecodeQueueItemOrder(requestBody.Object)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Name, "#999999")
	assert.Equal(t, int(result.ID), 5010223726653)
	assert.Equal(t, int(result.Number), 999999)
}

func TestDecodeQueueItemProduct(t *testing.T) {
	requestBody := QueueItemPayload("test-case-valid-product-queue-item.json")
	invalidRequestBody := QueueItemPayload("test-case-invalid-product-queue-item.json")

	// Test 1 - empty request body
	result, err := DecodeQueueItemProduct(json.RawMessage{})
	if err == nil {
		t.Errorf("Expected 'unexpected end of JSON input' but found: 'nil'")
	}
	assert.Equal(t, result.Shopify.ProductID, "")
	assert.Equal(t, result.Shopify.VariantID, "")
	assert.Equal(t, result.SystemProductID, "")
	assert.Equal(t, result.SystemVariantID, "")

	// Test 2 - invalid request body
	result, err = DecodeQueueItemProduct(invalidRequestBody.Object)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Shopify.ProductID, "")
	assert.Equal(t, result.Shopify.VariantID, "")
	assert.Equal(t, result.SystemProductID, "")
	assert.Equal(t, result.SystemVariantID, "")

	// Test 3  - valid request body
	result, err = DecodeQueueItemProduct(requestBody.Object)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.SystemProductID, "f6b1c96b-6079-41d5-9ec5-1203fb157206")
	assert.Equal(t, result.SystemVariantID, "8ca56b49-6655-4c98-9fa6-45804f61425d")
	assert.Equal(t, result.Shopify.ProductID, "7127845371965")
	assert.Equal(t, result.Shopify.VariantID, "40807533772861")
}

func TestQueueItemProductValidation(t *testing.T) {
	// Test Case 1 - valid push restrictions
	payload := QueueItemPayload("test-case-valid-product-queue-item.json").Object
	productPayload, err := DecodeQueueItemProduct(payload)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	valid := QueueItemProductValidation(productPayload)
	assert.Equal(t, nil, valid)

	// Test Case 2 - invalid push restrictions
	payload = QueueItemPayload("test-case-invalid-product-queue-item.json").Object
	productPayload, err = DecodeQueueItemProduct(payload)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	valid = QueueItemProductValidation(productPayload)
	assert.NotEqual(t, nil, valid)
}

func TestQueueItemValidation(t *testing.T) {
	requestBody := QueuePayload("queue_add_product.json")
	// Test 1 - empty request body
	result := QueueItemValidation(objects.RequestQueueItem{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = QueueItemValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodeQueueItem(t *testing.T) {
	requestBody := QueuePayload("queue_add_product.json")
	invalidRequestBody := AppSettingsPayload("test-case-valid-request.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeQueueItem(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Instruction, "")
	assert.Equal(t, result.Status, "")
	assert.Equal(t, result.Type, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeQueueItem(request)
	if err == nil {
		t.Errorf("Expected 'json: cannot unmarshal array into Go value of type objects.RequestQueueItem' but found: 'nil'")
	}
	assert.Equal(t, result.Instruction, "")
	assert.Equal(t, result.Status, "")
	assert.Equal(t, result.Type, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeQueueItem(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Instruction, "add_product")
	assert.Equal(t, result.Status, "in-queue")
	assert.Equal(t, result.Type, "product")
}

func TestSettingsValidation(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := AppSettingsPayload("test-case-valid-request")
	setting_keys, err := dbconfig.DB.GetAppSettingsList(context.Background())
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	// Test 1 - empty request body
	result := SettingsValidation([]objects.RequestSettings{}, setting_keys)
	assert.Equal(t, result, nil)

	// Test 2 - invalid setting key
	result = SettingsValidation([]objects.RequestSettings{
		{
			Key:   "mock_key",
			Value: "mock_value",
		},
	}, setting_keys)
	assert.NotEqual(t, result, nil)

	// Test 3  - valid request body
	result = SettingsValidation(requestBody, setting_keys)
	assert.Equal(t, result, nil)
}

func TestSettingValidation(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	setting_keys, err := dbconfig.DB.GetAppSettingsList(context.Background())
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	// Test 1 - empty request body
	result := SettingValidation(objects.RequestSettings{}, setting_keys)
	assert.NotEqual(t, result, nil)

	// Test 2 - invalid setting key
	result = SettingValidation(objects.RequestSettings{
		Key:   "mock_key",
		Value: "mock_value",
	},
		setting_keys)
	assert.NotEqual(t, result, nil)

	// Test 3  - valid request body
	result = SettingValidation(objects.RequestSettings{
		Key:   "app_enable_shopify_push",
		Value: "false",
	}, setting_keys)
	assert.Equal(t, result, nil)
}

func TestDecodeSettings(t *testing.T) {
	requestBody := AppSettingsPayload("test-case-valid-request.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeSettings(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(result), 0)

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeSettings(request)
	if err == nil {
		t.Errorf("Expected 'json: cannot unmarshal array into Go value of type objects.RequestBodyProduct' but found: 'nil'")
	}
	assert.Equal(t, len(result), 0)

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeSettings(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(result), 8)
}

func TestDecodeSetting(t *testing.T) {
	requestBody := AppSettingsPayload("test-case-valid-request.json")[0]
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeSetting(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Key, "")
	assert.Equal(t, result.Value, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeSetting(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Key, "")
	assert.Equal(t, result.Value, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeSetting(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Key, "app_queue_process_limit")
	assert.Equal(t, result.Value, "20")
}

func TestInventoryMapValidation(t *testing.T) {
	requestBody := WarehouseLocationPayload("test-case-valid-warehouse-location.json")

	// Test 1 - empty request body
	result := InventoryMapValidation(objects.RequestWarehouseLocation{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = InventoryMapValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodeInventoryMap(t *testing.T) {
	requestBody := WarehouseLocationPayload("test-case-valid-warehouse-location.json")
	invalidRequestBody := AppSettingsPayload("test-case-valid-request.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeInventoryMap(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.LocationID, "")
	assert.Equal(t, result.ShopifyWarehouseName, "")
	assert.Equal(t, result.WarehouseName, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeInventoryMap(request)
	if err == nil {
		t.Errorf("Expected 'json: cannot unmarshal array into Go value of type objects.RequestQueueItem' but found: 'nil'")
	}
	assert.Equal(t, result.LocationID, "")
	assert.Equal(t, result.ShopifyWarehouseName, "")
	assert.Equal(t, result.WarehouseName, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeInventoryMap(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.LocationID, "62274240573")
	assert.Equal(t, result.ShopifyWarehouseName, "Cape Town warehouse")
	assert.Equal(t, result.WarehouseName, "TestHouse")
}

func TestProductValidationDatabase(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := objects.CSVProduct{
		ProductCode:  "MOCK-PRODUCT-CODE",
		Active:       "1",
		Title:        "MOCK-TITLE",
		BodyHTML:     "",
		Category:     "MOCK-CATEGORY",
		Vendor:       "MOCK-VENDOR",
		ProductType:  "MOCK-PRODUCT-TYPE",
		SKU:          "product_sku",
		Option1Name:  "",
		Option1Value: "",
		Option2Name:  "",
		Option2Value: "",
		Option3Name:  "",
		Option3Value: "",
		Barcode:      "",
		Image1:       "",
		Image2:       "",
		Image3:       "",
		Warehouses:   []objects.CSVQuantity{},
		Pricing:      []objects.CSVPricing{},
	}
	// Test 1 - empty request body
	result := ProductValidationDatabase(objects.CSVProduct{}, &dbconfig)
	assert.Equal(t, result, nil)

	// Test 2  - valid request body | duplicate SKU
	createDatabaseProduct(&dbconfig)
	defer ClearProductTestData(&dbconfig)
	result = ProductValidationDatabase(requestBody, &dbconfig)
	assert.NotEqual(t, result, nil)

	// this should never be a case because the struct only allows
	// for three product options
	// hence it is excluded from tests

	// // Test 3  - valid request body | 3 product options
	// result = ProductValidationDatabase(requestBody, &dbconfig)
	// assert.Equal(t, result, nil)

	// this will never be a case either
	// unless an invalid product with 4 options
	// is added directly to the database

	// Test 4  - valid request body | invalid option position
	// result = ProductValidationDatabase(requestBody, &dbconfig)
	// assert.Equal(t, result, nil)

	// Test 5  - valid request body | option name already exists
	requestBody.SKU = "MOCK-PRODUCT-SKU"
	result = ProductValidationDatabase(requestBody, &dbconfig)
	assert.Equal(t, result, nil)
}

func InitMockHttpRequest(requestBody interface{}, requestMethod, requestURL string) *http.Request {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	if requestMethod == "" {
		requestMethod = MOCK_HTTP_REQUEST_METHOD
	}
	if requestURL == "" {
		requestURL = MOCK_HTTP_REQUEST_URL
	}
	request, err := http.NewRequest(requestMethod, requestURL, &buffer)
	if err != nil {
		log.Fatal("unable to create mock http request: " + err.Error())
	}
	return request
}
