package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"objects"
	"os"
	"strings"
	"testing"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const PRODUCT_CODE = "product_code"
const PRODUCT_CODE_SIMPLE = "product_code_simple"
const WEB_CUSTOMER_CODE = "TestFirstName TestLastName"
const ORDER_WEB_CODE = "#999999"

func TestPostOrderHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	orderPayload := OrderPayload("test-case-valid-order.json")
	w := Init("/api/orders", http.MethodPost, map[string][]string{}, orderPayload, &dbconfig, router)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid token param */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init("/api/orders?token=&api_key="+dbUser.ApiKey, http.MethodPost, map[string][]string{}, orderPayload, &dbconfig, router)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "invalid token", response.Message)

	/* Test 4 - valid token | invalid user referenced */
	w = Init(
		"/api/orders?token=b23a8af2f57870d8afd88fec713c5e59eb84ce8657321aeace26222536fa1565&api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 404, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "invalid token for user", response.Message)

	/* Test 5 - valid request | duplicate order */
	createDatabaseOrder(&dbconfig)
	requestHeaders := make(map[string][]string)
	requestHeaders["Mocker"] = []string{"true"}
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)
	ClearOrderTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	responsePostOrder := objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 6 - valid request | one line item | one shipping item taxes */
	ClearOrderTestData(&dbconfig)
	requestHeaders = make(map[string][]string)
	requestHeaders["Mocker"] = []string{"true"}
	orderPayload = OrderPayload("test-case-valid-order-one-shipping-product-line.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 7 - valid request | one line item | one shipping item no taxes */
	ClearOrderTestData(&dbconfig)
	orderPayload = OrderPayload("test-case-valid-order-one-product-line-no-tax.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 10 - valid request | duplicated customer web_code */
	createDatabaseCustomer(&dbconfig)
	ClearOrderTestData(&dbconfig)
	orderPayload = OrderPayload("test-case-valid-order.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 11 - valid request | no customer addresses */
	ClearOrderTestData(&dbconfig)
	orderPayload = OrderPayload("test-case-valid-order-no-customer-address.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 12 - valid request | order total of zero */
	ClearOrderTestData(&dbconfig)
	orderPayload = OrderPayload("test-case-valid-order-zero-value-totals.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)

	/* Test 13 - valid request | non-zero order total */
	ClearOrderTestData(&dbconfig)
	orderPayload = OrderPayload("test-case-valid-order.json")
	w = Init(
		"/api/orders?token="+dbUser.WebhookToken+"&api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, orderPayload, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	responsePostOrder = objects.RequestQueueHelper{}
	err = json.Unmarshal(w.Body.Bytes(), &responsePostOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "add_order", responsePostOrder.Instruction)
	assert.Equal(t, "in-queue", responsePostOrder.Status)
	assert.Equal(t, "order", responsePostOrder.Type)
}

func TestOrdersHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/orders?page=1",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid page number */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/orders?page=-16&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	responseOrders := []objects.Order{}
	err := json.Unmarshal(w.Body.Bytes(), &responseOrders)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(responseOrders))

	/* Test 3 - valid request | no results */
	w = Init(
		"/api/orders?page=1&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	responseOrders = []objects.Order{}
	err = json.Unmarshal(w.Body.Bytes(), &responseOrders)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(responseOrders))

	/* Test 4 - valid request | with results */
	createDatabaseOrder(&dbconfig)
	w = Init(
		"/api/orders?page=1&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearOrderTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	responseOrders = []objects.Order{}
	err = json.Unmarshal(w.Body.Bytes(), &responseOrders)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 1, len(responseOrders))
}

func TestOrderIDHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/orders/id",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid order ID */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/orders/id?api_key="+dbUser.ApiKey,
		http.MethodGet, make(map[string][]string), nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "could not decode order id: id", response.Message)

	/* Test 3 - valid request | do not exist */
	w = Init(
		"/api/orders/c2d29867-3d0b-d497-9191-18a9d8ee7830?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 404, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "not found", response.Message)

	/* Test 4 - valid request | exists */
	orderUUID := createDatabaseOrder(&dbconfig)
	w = Init(
		"/api/orders/"+orderUUID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearOrderTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	responseOrder := objects.Order{}
	err = json.Unmarshal(w.Body.Bytes(), &responseOrder)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, orderUUID, responseOrder.ID)
}

func TestOrderSearchHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/orders/search?q=test",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid search query param */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/orders/search?q=&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "invalid search param", response.Message)

	/* Test 3 - valid search query param | no results */
	w = Init(
		"/api/orders/search?q=test_param&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	orderSearchResponse := []objects.SearchOrder{}
	err = json.Unmarshal(w.Body.Bytes(), &orderSearchResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(orderSearchResponse))

	/* Test 4 - valid request | results */
	createDatabaseOrder(&dbconfig)
	w = Init(
		"/api/orders/search?q=%23"+ORDER_WEB_CODE[1:]+"&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	orderSearchResponse = []objects.SearchOrder{}
	err = json.Unmarshal(w.Body.Bytes(), &orderSearchResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 1, len(orderSearchResponse))
}

func TestProductVariantRemoveIDHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/products/id/variants/variant_id",
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid variant ID*/
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/id/variants/variant_id?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "could not decode variant id: id", response.Message)

	/* Test 3 - valid product ID but do not exist */
	w = Init(
		"/api/products/c2d29867-3d0b-d497-9191-18a9d8ee7830/variants/variant_id?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "could not decode variant id: variant_id", response.Message)

	/* Test 4 - valid product ID and variant ID but do not exist */
	w = Init(
		"/api/products/c2d29867-3d0b-d497-9191-18a9d8ee7830/variants/c2d29867-3d0b-d497-9191-18a9d8ee7830?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)

	/* Test 5 - valid request */
	productUUID := createDatabaseProduct(&dbconfig)
	dbProduct, _ := CompileProduct(&dbconfig, productUUID, context.Background(), false)
	w = Init(
		"/api/products/"+productUUID.String()+"/variants/"+dbProduct.Variants[0].ID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
}

func TestProductRemoveIDHandle(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/products/id",
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - invalid product ID */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/id?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "could not decode product id: id", response.Message)

	/* Test 3 - invalid product ID */
	w = Init(
		"/api/products/id?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "could not decode product id: id", response.Message)

	/* Test 4 - valid product ID but do not exist */
	w = Init(
		"/api/products/c2d29867-3d0b-d497-9191-18a9d8ee7830?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)

	/* Test 5 - valid request */
	productUUID := createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products/"+productUUID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodDelete, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
}

func TestProductExportRoute(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/products/export",
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - valid request | no products */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/export?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, true, strings.Contains(response.Message, "product_export-"))

	file1, err := os.Stat("." + response.Message[21:])
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, 0, file1.Size())

	/* Test 3 - valid request | products */
	createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products/export?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, true, strings.Contains(response.Message, "product_export-"))

	file2, err := os.Stat("." + response.Message[21:])
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, 0, file2.Size())
	assert.NotEqual(t, file2.Size(), file1.Size())
}

func TestProductImportRoute(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	mpw, multiPartFormData := CreateMultiPartFormData("test-case-invalid-request-no-auth.csv", "file")
	requestHeaders := make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType()}
	w := Init(
		"/api/products/import",
		http.MethodPost, requestHeaders, &multiPartFormData, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid request - no file */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	mpw, multiPartFormData = CreateMultiPartFormData("", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType()}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, &multiPartFormData, &dbconfig, router,
	)

	assert.Equal(t, 500, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "request Content-Type isn't multipart/form-data", response.Message)

	/* Test 3 - Invalid request headers */
	_, multiPartFormData = CreateMultiPartFormData("test-case-invalid-request-headers.csv", "file")
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, make(map[string][]string), &multiPartFormData, &dbconfig, router,
	)

	assert.Equal(t, 500, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "request Content-Type isn't multipart/form-data", response.Message)

	/* Test 4 - Invalid request file too large */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-invalid-request-file-size.csv", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType()}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)

	assert.Equal(t, 500, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "http: request body too large", response.Message)

	/* Test 5 - invalid request - not CSV file (not specified) */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-request.csv", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType()}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)

	assert.Equal(t, 500, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "only CSV extensions are supported", response.Message)

	/* Test 6 - invalid request - incorrect form key */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-request.csv", "test_file_key")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 500, w.Code)
	successResponse := objects.ImportResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &successResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, successResponse.FailCounter)
	assert.Equal(t, 0, successResponse.ProcessedCounter)
	assert.Equal(t, 0, successResponse.ProductsAdded)
	assert.Equal(t, 0, successResponse.ProductsUpdated)
	assert.Equal(t, 0, successResponse.VariantsAdded)
	assert.Equal(t, 0, successResponse.VariantsUpdated)

	/* Test 7 - Valid request - products failed to import (attempted duplicate SKU) */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-request-duplicate-sku.csv", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	successResponse = objects.ImportResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &successResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, successResponse.FailCounter)
	assert.Equal(t, 5, successResponse.ProcessedCounter)
	assert.Equal(t, 1, successResponse.ProductsAdded)
	assert.Equal(t, 4, successResponse.ProductsUpdated)
	assert.Equal(t, 1, successResponse.VariantsAdded)
	assert.Equal(t, 4, successResponse.VariantsUpdated)

	/* Test 8 - Valid request - Products/variants created */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-request-variants-products-added.csv", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	successResponse = objects.ImportResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &successResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, successResponse.FailCounter)
	assert.Equal(t, 1, successResponse.ProcessedCounter)
	assert.Equal(t, 1, successResponse.ProductsAdded)
	assert.Equal(t, 0, successResponse.ProductsUpdated)
	assert.Equal(t, 1, successResponse.VariantsAdded)
	assert.Equal(t, 0, successResponse.VariantsUpdated)

	/* Test 9 - Valid request - Products/variants updated (should be zero created) */
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-request-variants-products-updated.csv", "file")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	w = Init(
		"/api/products/import?api_key="+dbUser.ApiKey,
		http.MethodPost, requestHeaders, multiPartFormData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	successResponse = objects.ImportResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &successResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, successResponse.FailCounter)
	assert.Equal(t, 2, successResponse.ProcessedCounter)
	assert.Equal(t, 1, successResponse.ProductsAdded)
	assert.Equal(t, 1, successResponse.ProductsUpdated)
	assert.Equal(t, 1, successResponse.VariantsAdded)
	assert.Equal(t, 1, successResponse.VariantsUpdated)
}

func TestProductCreationRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/products",
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid simple product data */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	productData := ProductPayload("test-case-invalid-product-title-simple.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "empty title not allowed", response.Message)

	/* Test 3 - Invalid variable product data */
	productData = ProductPayload("test-case-invalid-product-title-variable.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "empty title not allowed", response.Message)

	/* Test 3 - Invalid product options */
	productData = ProductPayload("test-case-invalid-product-variable.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "duplicate product option names not allowed: Size", response.Message)

	/* Test 4 - Invalid product sku | duplicated SKU */
	createDatabaseProduct(&dbconfig)
	productData = ProductPayload("test-case-valid-product-variable.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 409, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "SKU with code product_sku already exists", response.Message)

	/* Test 5 - Invalid product | duplicate options */
	createDatabaseProduct(&dbconfig)
	productData = ProductPayload("test-case-invalid-product-variable-duplicate-options.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "duplicate product option names not allowed: Size", response.Message)

	/* Test 6 - Valid variable product request | not added to shopify */
	productData = ProductPayload("test-case-valid-product-variable.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 201, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)

	/* Test 6 - Valid simple product request | not added to shopify */
	productData = ProductPayload("test-case-valid-product-simple.json")
	w = Init(
		"/api/products?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, productData, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 201, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
}

func TestProductFilterRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/products/filter?page=1",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid filter params (empty) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/filter?page=&type=&category=&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 4 - No filter results */
	w = Init(
		"/api/products/filter?type=simple&category=test&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response := []objects.SearchProduct{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(response))

	/* Test 5 - Valid filter request */
	createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products/filter?type=product_product_type&vendor=product_vendor&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response = []objects.SearchProduct{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 1, len(response))
}

func TestProductSearchRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/products/search?q=product_title",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid search param (empty) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/search?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 4 - No search results */
	w = Init(
		"/api/products/search?q=simple&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response := []objects.SearchProduct{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(response))

	/* Test 5 - Valid search request */
	createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products/search?q=product_title&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response = []objects.SearchProduct{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 1, len(response))
}

func TestProductsRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/products?page=1",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid page number */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products?page=-1&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 4 - Invalid page number (string) */
	w = Init(
		"/api/products?page=two&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)

	/* Test 5 - Valid request */
	createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products?page=1&api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)

	assert.Equal(t, 200, w.Code)
	response := []objects.Product{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.NotEqual(t, 0, len(response))
}

func TestProductIDRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)
	w := Init(
		"/api/products/abctest123",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid product_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = Init(
		"/api/products/abctest123?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 4 - Invalid product_id (404) */
	w = Init(
		"/api/products/"+uuid.New().String()+"?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 404, w.Code)

	/* Test 5 - Valid request */
	productUUID := createDatabaseProduct(&dbconfig)
	w = Init(
		"/api/products/"+productUUID.String()+"?api_key="+dbUser.ApiKey,
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)
	ClearProductTestData(&dbconfig)
	assert.Equal(t, 200, w.Code)
}

func TestLoginRoute(t *testing.T) {
	/* Test 1 - Invalid request - empty username/password */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	loginData := LoginPayload("test-case-invalid-login.json")
	w := Init(
		"/api/login",
		http.MethodPost, map[string][]string{}, loginData, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "empty username not allowed", response.Message)

	/* Test 2 - Invalid request - non empty username/password but invalid credentials) */
	loginData = LoginPayload("test-case-valid-login.json")
	w = Init(
		"/api/login",
		http.MethodPost, map[string][]string{}, loginData, &dbconfig, router,
	)
	assert.Equal(t, 404, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "invalid username and password combination", response.Message)

	/* Test 3 - Valid request */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	loginData = LoginPayload("test-case-valid-login.json")
	loginData.Username = dbUser.Name
	loginData.Password = dbUser.Password
	w = Init(
		"/api/login",
		http.MethodPost, map[string][]string{}, loginData, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	responseLogin := objects.ResponseLogin{}
	err = json.Unmarshal(w.Body.Bytes(), &responseLogin)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "test", responseLogin.Username)
}

func TestLogoutHandle(t *testing.T) {
	/* Test 1 - Invalid request - no cookies and no authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/logout",
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid request - no cookies sent with request */
	dbUser := createDatabaseUser(&dbconfig)
	router = setUpAPI(&dbconfig)

	w = Init(
		"/api/logout?api_key="+dbUser.ApiKey,
		http.MethodPost, map[string][]string{}, nil, &dbconfig, router,
	)
	assert.Equal(t, 200, w.Code)
}

func TestPreregisterRoute(t *testing.T) {
	/* Test 1 - Invalid request (empty email) */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	preregisterData := PreRegisterPayload("test-case-invalid-preregister.json")
	w := Init(
		"/api/preregister",
		http.MethodPost, map[string][]string{}, preregisterData, &dbconfig, router,
	)
	assert.Equal(t, 400, w.Code)

	/* Test 2 - Email already exists */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	preregisterData = PreRegisterPayload("test-case-valid-preregister.json")
	w = Init(
		"/api/preregister",
		http.MethodPost, map[string][]string{}, preregisterData, &dbconfig, router,
	)

	assert.Equal(t, 409, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "email '"+preregisterData.Email+"' already exists", response.Message)
	dbconfig.DB.DeleteTokenByEmail(context.Background(), preregisterData.Email)

	/* Test 3 - Valid request */
	preregisterData = PreRegisterPayload("test-case-valid-preregister.json")
	requestHeaders := make(map[string][]string)
	requestHeaders["Mocker"] = []string{"true"}
	w = Init(
		"/api/preregister",
		http.MethodPost, requestHeaders, preregisterData, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "email sent", response.Message)

	dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	dbconfig.DB.DeleteTokenByEmail(context.Background(), preregisterData.Email)
}

func TestRegisterRoute(t *testing.T) {
	/* Test 1 - Valid request */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	registrationData := UserPayload("test-case-valid-user.json")
	register_data_token := createDatabasePreregister(registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	w := Init(
		"/api/register",
		http.MethodPost, map[string][]string{}, registrationData, &dbconfig, router,
	)

	assert.Equal(t, 201, w.Code)
	response := objects.ResponseRegister{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "test", response.Name)
	assert.Equal(t, "test@test.com", response.Email)

	dbconfig.DB.RemoveUser(context.Background(), response.ApiKey)
	dbconfig.DB.DeleteToken(context.Background(), database.DeleteTokenParams{
		Token: register_data_token.Token,
		Email: register_data_token.Email,
	})

	/* Test 2 - Invalid request body */
	new_registration_data := ProductPayload("test-case-valid-product-simple.json")
	w = Init(
		"/api/register",
		http.MethodPost, map[string][]string{}, new_registration_data, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Email)

	/* Test 3 - Invalid token */
	registrationData = UserPayload("test-case-invalid-user.json")
	w = Init(
		"/api/register",
		http.MethodPost, map[string][]string{}, registrationData, &dbconfig, router,
	)

	assert.Equal(t, 404, w.Code)
	response = objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Email)

	/* Test 4 - User already exist */
	db_user := createDatabaseUser(&dbconfig)
	registrationData = UserPayload("test-case-valid-user.json")
	register_data_token = createDatabasePreregister(registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	w = Init(
		"/api/register",
		http.MethodPost, map[string][]string{}, registrationData, &dbconfig, router,
	)

	assert.Equal(t, 409, w.Code)
	response = objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Email)
	dbconfig.DB.RemoveUser(context.Background(), db_user.ApiKey)
	dbconfig.DB.DeleteToken(context.Background(), database.DeleteTokenParams{
		Token: register_data_token.Token,
		Email: register_data_token.Email,
	})

	/* Test 5 - Empty username/password in request */
	registrationData = UserPayload("test-case-invalid-user-empty-username-password.json")
	register_data_token = createDatabasePreregister(registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	w = Init(
		"/api/register",
		http.MethodPost, map[string][]string{}, registrationData, &dbconfig, router,
	)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Email)
	dbconfig.DB.RemoveUser(context.Background(), db_user.ApiKey)
	dbconfig.DB.DeleteToken(context.Background(), database.DeleteTokenParams{
		Token: register_data_token.Token,
		Email: register_data_token.Email,
	})
}

func TestReadyRoute(t *testing.T) {
	/* Test 1 - Valid database credentials */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := Init(
		"/api/ready",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 200, w.Code)
	response_string := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response_string)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "OK", response_string.Message)

	/* Test 2 - Invalid database credentials */
	dbconfig = setupDatabase("test_user", "test_psw", "database_test", true)
	w = Init(
		"/api/ready",
		http.MethodGet, map[string][]string{}, nil, &dbconfig, router,
	)

	assert.Equal(t, 503, w.Code)
	response_string = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response_string)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "Unavailable", response_string.Message)
}

/* Function that creates a multipart/form request to be used in the import handle */
func CreateMultiPartFormData(fileName, formKey string) (*multipart.Writer, bytes.Buffer) {
	// Create a buffer to store the request body
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	// Add a file to the request
	file, err := os.Open("./test_payloads/tests/import/" + fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	defer file.Close()

	// Create a new form field
	fw, err := w.CreateFormFile(formKey, "./test_payloads/tests/import/"+fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	w.Close()
	return w, buf
}

/*
Creates a queue item inside the queue of the respective queue_type

Data is retrived from the project directory `test_payloads`
*/
func createQueueItem(queue_type string) objects.RequestQueueItem {
	file, err := os.Open("./test_payloads/queue/queue_" + queue_type + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	orderData := objects.RequestQueueItem{}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		fmt.Println(err)
	}
	return orderData
}

/* Returns an order request body struct */
func OrderPayload(fileName string) objects.RequestBodyOrder {
	fileBytes := payload("./test_payloads/tests/orders/" + fileName)
	orderData := objects.RequestBodyOrder{}
	err := json.Unmarshal(fileBytes, &orderData)
	if err != nil {
		log.Println(err)
	}
	return orderData
}

/* Returns an order request body struct */
func CustomerPayload(fileName string) objects.RequestBodyCustomer {
	fileBytes := payload("./test_payloads/tests/customers/" + fileName)
	customerData := objects.RequestBodyCustomer{}
	err := json.Unmarshal(fileBytes, &customerData)
	if err != nil {
		log.Println(err)
	}
	return customerData
}

/* Returns a product request body struct */
func ProductPayload(fileName string) objects.RequestBodyProduct {
	fileBytes := payload("./test_payloads/tests/products/" + fileName)
	productData := objects.RequestBodyProduct{}
	err := json.Unmarshal(fileBytes, &productData)
	if err != nil {
		log.Println(err)
	}
	return productData
}

/* Returns a register request body struct */
func RegisterPayload(fileName string) objects.RequestBodyRegister {
	fileBytes := payload("./test_payloads/tests/register/" + fileName)
	registerData := objects.RequestBodyRegister{}
	err := json.Unmarshal(fileBytes, &registerData)
	if err != nil {
		log.Println(err)
	}
	return registerData
}

/* Returns a pre-registrater request body struct */
func PreRegisterPayload(fileName string) objects.RequestBodyPreRegister {
	fileBytes := payload("./test_payloads/tests/preregister/" + fileName)
	preregData := objects.RequestBodyPreRegister{}
	err := json.Unmarshal(fileBytes, &preregData)
	if err != nil {
		log.Println(err)
	}
	return preregData
}

/* Returns a login request body struct */
func LoginPayload(fileName string) objects.RequestBodyLogin {
	fileBytes := payload("./test_payloads/tests/login/" + fileName)
	loginData := objects.RequestBodyLogin{}
	err := json.Unmarshal(fileBytes, &loginData)
	if err != nil {
		log.Println(err)
	}
	return loginData
}

/* Returns a test user RequestBodyRegister struct */
func UserPayload(fileName string) objects.RequestBodyRegister {
	fileBytes := payload("./test_payloads/tests/users/" + fileName)
	userRegistrationData := objects.RequestBodyRegister{}
	err := json.Unmarshal(fileBytes, &userRegistrationData)
	if err != nil {
		log.Println(err)
	}
	return userRegistrationData
}

/* Returns a test user RequestBodyRegister struct */
func CreateTokenPayload(fileName string) objects.RequestBodyPreRegister {
	fileBytes := payload("./test_payloads/tests/preregister/" + fileName)
	userRegistrationData := objects.RequestBodyPreRegister{}
	err := json.Unmarshal(fileBytes, &userRegistrationData)
	if err != nil {
		log.Println(err)
	}
	return userRegistrationData
}

/*
Returns a byte array representing the file data that was read

Data is retrived from the project directory `test_payloads`
*/
func payload(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	return respBody
}

/*
Creates a test product in the database
*/
func createDatabaseProduct(dbconfig *DbConfig) uuid.UUID {
	product, err := dbconfig.DB.GetProductByProductCode(context.Background(), "product_code")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return uuid.Nil
		}
	}
	if product.ProductCode == "" {
		product := ProductPayload("test-case-valid-product-variable.json")
		productUUID, _, err := AddProduct(dbconfig, product)
		if err != nil {
			log.Println(err)
			return uuid.Nil
		}
		return productUUID
	}
	return product.ID
}

/*
Creates a test user in the database
*/
func createDatabaseCustomer(dbconfig *DbConfig) uuid.UUID {
	customer, err := dbconfig.DB.GetCustomerByWebCode(context.Background(), WEB_CUSTOMER_CODE)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return uuid.Nil
		}
	}
	if customer.WebCustomerCode == "" {
		customer := CustomerPayload("test-case-valid-customer.json")
		dbCustomer, err := AddCustomer(dbconfig, customer, customer.FirstName+" "+customer.LastName)
		if err != nil {
			log.Println(err)
			return uuid.Nil
		}
		return dbCustomer
	}
	return customer.ID
}

/*
Creates a test order in the database
*/
func createDatabaseOrder(dbconfig *DbConfig) uuid.UUID {
	order, err := dbconfig.DB.GetOrderByWebCode(context.Background(), "1000")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return uuid.Nil
		}
	}
	if order.WebCode == "" {
		order := OrderPayload("test-case-valid-order.json")
		orderUUID, err := AddOrder(dbconfig, order)
		if err != nil {
			log.Println(err)
			return uuid.Nil
		}
		return orderUUID
	}
	return order.ID
}

/*
Creates a test user in the database
*/
func createDatabaseUser(dbconfig *DbConfig) database.User {
	user, err := dbconfig.DB.GetUserByEmail(context.Background(), "test@test.com")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.User{}
		}
	}
	if user.ApiKey == "" {
		user := UserPayload("test-case-valid-user.json")
		dbUser, err := AddUser(dbconfig, user)
		if err != nil {
			log.Println(err)
			return database.User{}
		}
		return dbUser
	}
	return user
}

/*
Creates a demo token in the database for registration
*/
func createDatabasePreregister(email string, dbconfig *DbConfig) database.RegisterToken {
	token, err := dbconfig.DB.GetTokenValidation(context.Background(), email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.RegisterToken{}
		}
	}
	if token.Token == uuid.Nil {
		token := CreateTokenPayload("test-case-valid-preregister.json")
		dbToken, err := AddUserRegistration(dbconfig, token)
		if err != nil {
			log.Println(err)
			return database.RegisterToken{}
		}
		return dbToken
	}
	return database.RegisterToken{}
}

/*
Helper function to setup a database connection for tests.

Will only overwrite with params if `overwrite` is set to true
*/
func setupDatabase(param_db_user, param_db_psw, param_db_name string, overwrite bool) DbConfig {
	db_user := utils.LoadEnv("db_user")
	db_psw := utils.LoadEnv("db_psw")
	db_name := utils.LoadEnv("db_name")
	if overwrite {
		db_user = param_db_user
		db_psw = param_db_psw
		db_name = param_db_name
	}
	connection_string := "postgres://" + db_user + ":" + db_psw
	dbCon, err := InitConn(connection_string + "@127.0.0.1:5432/" + db_name + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}
	return dbCon
}

func ClearTestData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveCustomerByWebCustomerCode(context.Background(), WEB_CUSTOMER_CODE)
	dbconfig.DB.RemoveProductByCode(context.Background(), PRODUCT_CODE_SIMPLE)
	dbconfig.DB.RemoveProductByCode(context.Background(), PRODUCT_CODE)
	dbconfig.DB.RemoveOrderByWebCode(context.Background(), ORDER_WEB_CODE)
}

func ClearOrderTestData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveOrderByWebCode(context.Background(), ORDER_WEB_CODE)
	dbconfig.DB.RemoveCustomerByWebCustomerCode(context.Background(), WEB_CUSTOMER_CODE)
}

func ClearProductTestData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveProductByCode(context.Background(), PRODUCT_CODE_SIMPLE)
	dbconfig.DB.RemoveProductByCode(context.Background(), PRODUCT_CODE)
}

func ClearCustomerTestData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveCustomerByWebCustomerCode(context.Background(), WEB_CUSTOMER_CODE)
}

func Init(
	requestURL,
	requestMethod string,
	additionalRequestHeaders map[string][]string,
	payload interface{},
	dbconfig *DbConfig,
	router *gin.Engine,
) *httptest.ResponseRecorder {
	var buffer bytes.Buffer
	switch dataType := payload.(type) {
	case bytes.Buffer:
		buffer = dataType
	default:
		err := json.NewEncoder(&buffer).Encode(payload)
		if err != nil {
			log.Fatal(err)
		}
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(requestMethod, requestURL, &buffer)
	for key, value := range additionalRequestHeaders {
		for _, sub_value := range value {
			req.Header.Add(key, sub_value)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}
