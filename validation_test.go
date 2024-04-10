package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
