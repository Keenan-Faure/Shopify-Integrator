package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"objects"
	"testing"

	"github.com/google/uuid"
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

func TestProductSKUValidation(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - empty request body
	result := ProductSKUValidation("", &dbconfig)
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request | duplicate SKU
	createDatabaseProduct(&dbconfig)
	defer ClearProductTestData(&dbconfig)
	result = ProductSKUValidation("product_sku", &dbconfig)
	assert.NotEqual(t, result, nil)

	// Test 3  - valid request body | sku not found
	result = ProductSKUValidation("mock-product-sku", &dbconfig)
	assert.Equal(t, result, nil)
}

func TestProductOptionNameValidation(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - empty request body
	result := ProductOptionNameValidation("", "", &dbconfig)
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request | duplicate option name
	createDatabaseProduct(&dbconfig)
	defer ClearProductTestData(&dbconfig)
	result = ProductOptionNameValidation("product_code", "Size", &dbconfig)
	assert.NotEqual(t, result, nil)
}

func TestValidateTokenValidation(t *testing.T) {
	requestBody := UserPayload("test-case-valid-user.json")

	// Test 1 - empty request body
	result := ValidateTokenValidation(objects.RequestBodyRegister{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = ValidateTokenValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodeValidateTokenRequestBody(t *testing.T) {
	requestBody := UserPayload("test-case-valid-user.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeValidateTokenRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Password, "")
	assert.Equal(t, result.Token, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeValidateTokenRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Password, "")
	assert.Equal(t, result.Token, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeValidateTokenRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "test@test.com")
	assert.Equal(t, result.Name, "test")
	assert.Equal(t, result.Password, "abc12345678910")
	assert.Equal(t, result.Token, "c266e9f6-1ca6-4e27-8dd8-cce2bf5fdba5")
}

func TestPreRegisterValidation(t *testing.T) {
	requestBody := CreateTokenPayload("test-case-valid-preregister.json")

	// Test 1 - empty request body
	result := PreRegisterValidation(objects.RequestBodyPreRegister{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = PreRegisterValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodePreRegisterRequestBody(t *testing.T) {
	requestBody := CreateTokenPayload("test-case-valid-preregister.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodePreRegisterRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodePreRegisterRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodePreRegisterRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "test@test.com")
	assert.Equal(t, result.Name, "test")
}

func TestCustomerValidation(t *testing.T) {
	requestBody := CustomerPayload("test-case-valid-customer.json")

	// Test 1 - empty request body
	result := CustomerValidation(objects.RequestBodyCustomer{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = CustomerValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodeCustomerRequestBody(t *testing.T) {
	requestBody := CustomerPayload("test-case-valid-customer.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeCustomerRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.FirstName, "")
	assert.Equal(t, result.LastName, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeCustomerRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.FirstName, "")
	assert.Equal(t, result.LastName, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeCustomerRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "keenan@test.com")
	assert.Equal(t, result.FirstName, "TestFirstName")
	assert.Equal(t, result.LastName, "TestLastName")
}

func TestOrderValidation(t *testing.T) {
	requestBody := OrderPayload("test-case-valid-order.json")

	// Test 1 - empty request body
	result := OrderValidation(objects.RequestBodyOrder{})
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = OrderValidation(requestBody)
	assert.Equal(t, result, nil)
}

func TestDecodeOrderRequestBody(t *testing.T) {
	requestBody := OrderPayload("test-case-valid-order.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-simple.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeOrderRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Note, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeOrderRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Note, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeOrderRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "keenan@stock2shop.com")
	assert.Equal(t, result.Name, "#999999")
	assert.Equal(t, result.Note, "Notes not taken")
}

func TestTokenValidation(t *testing.T) {
	// Test 1 - empty request body
	result := TokenValidation("")
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = TokenValidation("c266e9f6-1ca6-4e27-8dd8-cce2bf5fdba5")
	assert.Equal(t, result, nil)
}

func TestUserValidation(t *testing.T) {
	// Test 1 - empty request body
	result := UserValidation("", "")
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = UserValidation("mock-user-name", "mock-user-password")
	assert.Equal(t, result, nil)
}

func TestIDValidation(t *testing.T) {
	// Test 1 - empty request body
	result := IDValidation("")
	assert.NotEqual(t, result, nil)

	// Test 2  - valid request body
	result = IDValidation("c266e9f6-1ca6-4e27-8dd8-cce2bf5fdba5")
	assert.Equal(t, result, nil)
}

func TestProductValidation(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := ProductPayload("test-case-valid-product-variable.json")

	// Test 1 - empty request body
	result := ProductValidation(&dbconfig, objects.RequestBodyProduct{})
	assert.NotEqual(t, result, nil)

	// Test 2 - empty request body | empty title
	requestBody.Title = ""
	result = ProductValidation(&dbconfig, requestBody)
	assert.NotEqual(t, result, nil)
	requestBody.Title = "mock-title"

	// Test 3 - empty request body | invalid SKU
	requestBody.Variants[0].Sku = ""
	result = ProductValidation(&dbconfig, requestBody)
	assert.NotEqual(t, result, nil)
	requestBody.Variants[0].Sku = "mock-variant-0-sku"

	// Test 4 - empty request body | empty variants
	requestBody.Variants = []objects.RequestBodyVariant{}
	result = ProductValidation(&dbconfig, requestBody)
	assert.NotEqual(t, result, nil)

	// Test 5  - valid request body
	requestBody = ProductPayload("test-case-valid-product-variable.json")
	result = ProductValidation(&dbconfig, requestBody)
	assert.Equal(t, result, nil)
}

func TestValidateDuplicateOption(t *testing.T) {
	requestBody := ProductPayload("test-case-valid-product-variable.json")
	invalidRequestBody := ProductPayload("test-case-invalid-product-variable.json")

	// Test 1 - empty request body
	result := ValidateDuplicateOption(objects.RequestBodyProduct{})
	assert.Equal(t, result, nil)

	// Test 2 - invalid request body
	result = ValidateDuplicateOption(invalidRequestBody)
	assert.NotEqual(t, result, nil)

	// Test 3  - valid request body
	result = ValidateDuplicateOption(requestBody)
	assert.Equal(t, result, nil)
}

func TestValidateDuplicateSKU(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := ProductPayload("test-case-valid-product-variable.json")
	invalidRequestBody := ProductPayload("test-case-invalid-product-variable.json")

	// Test 1 - empty request body
	result := ValidateDuplicateSKU(objects.RequestBodyProduct{}, &dbconfig)
	assert.Equal(t, result, nil)

	// Test 2 - invalid request body | duplicate SKU
	createDatabaseProduct(&dbconfig)
	result = ValidateDuplicateSKU(invalidRequestBody, &dbconfig)
	assert.NotEqual(t, result, nil)
	ClearProductTestData(&dbconfig)

	// Test 3  - valid request body
	result = ValidateDuplicateSKU(requestBody, &dbconfig)
	assert.Equal(t, result, nil)
}

func TestDuplicateOptionValues(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	requestBody := ProductPayload("test-case-valid-product-variable.json")

	// Test 1 - empty request body
	result := DuplicateOptionValues(&dbconfig, objects.RequestBodyVariant{}, uuid.Nil)
	assert.Equal(t, result, nil)

	// Test 2 - invalid request body | duplicate SKU
	productID := createDatabaseProduct(&dbconfig)
	result = DuplicateOptionValues(&dbconfig, requestBody.Variants[0], productID)
	assert.NotEqual(t, result, nil)
	ClearProductTestData(&dbconfig)

	// Test 3 - invalid request body | duplicate option values
	productID = createDatabaseProduct(&dbconfig)
	requestBody.Variants[0].Sku = "product_sku1"
	requestBody.Variants[0].Option1 = "option1"
	requestBody.Variants[0].Option2 = "option2"
	requestBody.Variants[0].Option3 = "option3"
	result = DuplicateOptionValues(&dbconfig, requestBody.Variants[0], productID)
	assert.NotEqual(t, result, nil)
	ClearProductTestData(&dbconfig)

	// Test 4 - invalid request body | invalid option value
	productID = createDatabaseProduct(&dbconfig)
	requestBody.Variants[0].Sku = "product_sku1"
	requestBody.Variants[0].Option1 = "option4"
	requestBody.Variants[0].Option2 = "option5"
	requestBody.Variants[0].Option3 = ""
	result = DuplicateOptionValues(&dbconfig, requestBody.Variants[0], productID)
	assert.NotEqual(t, result, nil)
	ClearProductTestData(&dbconfig)

	// Test 5  - valid request body
	requestBody.Variants[0].Sku = "product_sku1"
	requestBody.Variants[0].Option1 = "option4"
	requestBody.Variants[0].Option2 = "option5"
	requestBody.Variants[0].Option3 = "option6"
	result = DuplicateOptionValues(&dbconfig, requestBody.Variants[0], uuid.Nil)
	assert.Equal(t, result, nil)
}

func TestCreateProductOptionSlice(t *testing.T) {
	// Test 1 - invalid (empty option values)
	assert.Equal(t, len(CreateProductOptionSlice("", "", "")), 0)

	// Test 2 - invalid (partial empty option values)
	assert.Equal(t, len(CreateProductOptionSlice("", "option2", "option3")), 0)

	// Test 3 - valid option values
	assert.Equal(t, len(CreateProductOptionSlice("option1", "option2", "option3")), 3)
}

func TestDecodeProductRequestBody(t *testing.T) {
	requestBody := ProductPayload("test-case-valid-product-variable.json")
	invalidRequestBody := CustomerPayload("test-case-valid-customer.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeProductRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Active, "")
	assert.Equal(t, result.Title, "")
	assert.Equal(t, result.Category, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeProductRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Active, "")
	assert.Equal(t, result.Title, "")
	assert.Equal(t, result.Category, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeProductRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.BodyHTML, "<p>I am a body_html</p>")
	assert.Equal(t, result.Title, "product_title")
	assert.Equal(t, result.Category, "product_category")
}

func TestDecodeUserRequestBody(t *testing.T) {
	requestBody := UserPayload("test-case-valid-user.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-variable.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeUserRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Password, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeUserRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")
	assert.Equal(t, result.Password, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeUserRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Email, "test@test.com")
	assert.Equal(t, result.Name, "test")
	assert.Equal(t, result.Password, "abc12345678910")
	assert.Equal(t, result.Token, "c266e9f6-1ca6-4e27-8dd8-cce2bf5fdba5")
}

func TestDecodeLoginRequestBody(t *testing.T) {
	requestBody := LoginPayload("test-case-valid-login.json")
	invalidRequestBody := ProductPayload("test-case-valid-product-variable.json")

	// Test 1 - empty request body
	request := InitMockHttpRequest(nil, "", "")
	result, err := DecodeLoginRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Username, "")
	assert.Equal(t, result.Password, "")

	// Test 2 - invalid request body
	request = InitMockHttpRequest(invalidRequestBody, "", "")
	result, err = DecodeLoginRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, result.Username, "")
	assert.Equal(t, result.Password, "")

	// Test 3  - valid request body
	request = InitMockHttpRequest(requestBody, "", "")
	result, err = DecodeLoginRequestBody(request)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, result.Username, "test")
	assert.Equal(t, result.Password, "testPassword")
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
