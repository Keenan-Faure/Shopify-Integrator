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
	"testing"
	"time"
	"utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductImportRoute(t *testing.T) {
	/* Test 1 - invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	mpw, multiPartFormData := CreateMultiPartFormData("test-case-invalid-request-no-auth.csv", "file")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products/import", &multiPartFormData)
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid request - no file */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	mpw, multiPartFormData = CreateMultiPartFormData("test-case-invalid-request-no-auth.csv", "file")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products/import?api_key="+dbUser.ApiKey, &multiPartFormData)
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	response := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "http: no such file", response.Message)

	/* Test 3 - Invalid request headers */
	_, multiPartFormData = CreateMultiPartFormData("test-case-invalid-request-headers.csv", "file")

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products/import?api_key="+dbUser.ApiKey, &multiPartFormData)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "request Content-Type isn't multipart/form-data", response.Message)

	/* Test 4 - Invalid request file too large */

	/* Test 5 - invalid request - incorrect form key */

	/* Test 6 - Valid request - products failed to import (duplicate SKU) */

	/* Test 7 - Valid request - Products/variants created */

	/* Test 8 - Valid request - Products/variants updated (should be zero created) */
}

func TestProductCreationRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", nil)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid product data */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	productData := ProductPayload()
	productData.Title = ""
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products?api_key="+dbUser.ApiKey, &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	response := objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "empty title not allowed", response.Message)

	/* Test 3 - Invalid product options */
	productData.Title = "product_title"
	productData.ProductOptions[0].Value = "Colour"
	err = json.NewEncoder(&buffer).Encode(productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products?api_key="+dbUser.ApiKey, &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "duplicate product option names not allowed: Colour", response.Message)

	/* Test 4 - Invalid product sku | duplicated SKU */
	createDatabaseProduct(&dbconfig)
	productData.ProductOptions[0].Value = "Size"

	err = json.NewEncoder(&buffer).Encode(productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products?api_key="+dbUser.ApiKey, &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "SKU with code product_sku already exists", response.Message)
	dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")

	/* Test 5 - Invalid product | duplicate options */
	createDatabaseProduct(&dbconfig)

	productData.Variants[0].Option1 = "option1"
	productData.Variants[0].Option2 = "option2"
	productData.Variants[0].Option3 = "option3"
	productData.Variants[0].Sku = "product_sku1"

	err = json.NewEncoder(&buffer).Encode(productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products?api_key="+dbUser.ApiKey, &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "duplicate option values not allowed", response.Message)
	dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")

	/* Test 6 - Valid product request | not added to shopify */

	err = json.NewEncoder(&buffer).Encode(productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/products?api_key="+dbUser.ApiKey, &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	response = objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "success", response.Message)
	dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")
}

func TestProductFilterRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/filter?page=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid filter params (empty) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/filter?page=&type=&category=&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	/* Test 4 - No filter results */
	defer dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/filter?type=simple&category=test&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	response := []objects.SearchProduct{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(response))

	/* Test 5 - Valid filter request */
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/filter?type=product_product_type&vendor=product_vendor&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/search?q=product_title", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid search param (empty) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/search?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 4 - No search results */
	createDatabaseProduct(&dbconfig)
	defer dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/search?q=simple&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	response := []objects.SearchProduct{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, 0, len(response))

	/* Test 5 - Valid search request */
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/search?q=product_title&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products?page=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid page number */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products?page=-1&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	/* Test 4 - Invalid page number (string) */
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products?page=two&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	/* Test 5 - Valid request */
	createDatabaseProduct(&dbconfig)
	defer dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products?page=1&api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/abctest123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	/* Test 2 - Invalid product_id (malformed) */
	dbUser := createDatabaseUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/abctest123?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 4 - Invalid product_id (404) */
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+uuid.New().String()+"?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	/* Test 5 - Valid request */
	productUUID := createDatabaseProduct(&dbconfig)
	defer dbconfig.DB.RemoveProductByCode(context.Background(), "product_code")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+productUUID.String()+"?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLoginRoute(t *testing.T) {
	/* Test 1 - Invalid request - empty username/password */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	loginData := LoginPayload()
	loginData.Username = ""
	loginData.Password = ""
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(loginData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 2 - Invalid request - non empty username/password but invalid credentials) */
	loginData = LoginPayload()
	err = json.NewEncoder(&buffer).Encode(loginData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/login", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	response := objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "invalid username and password combination", response.Message)

	/* Test 3 - Valid request */
	dbUser := createDatabaseUser(&dbconfig)

	loginData = LoginPayload()
	loginData.Username = dbUser.Name
	loginData.Password = dbUser.Password
	err = json.NewEncoder(&buffer).Encode(loginData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/login", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	responseLogin := objects.ResponseLogin{}
	err = json.Unmarshal(w.Body.Bytes(), &responseLogin)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "test", responseLogin.Username)
	dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)
}

func TestLogoutHandle(t *testing.T) {
	/* Test 1 - Invalid request - no cookies and no authentication */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/logout", nil)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 2 - Invalid request - no cookies sent with request */
	dbUser := createDatabaseUser(&dbconfig)
	router = setUpAPI(&dbconfig)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/logout?api_key="+dbUser.ApiKey, nil)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPreregisterRoute(t *testing.T) {
	/* Test 1 - Invalid request (empty email) */
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	preregisterData := PreRegisterPayload()
	preregisterData.Email = ""
	preregisterData.Name = ""
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(preregisterData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 2 - Email already exists */
	dbUser := createDatabaseUser(&dbconfig)

	preregisterData = PreRegisterPayload()
	err = json.NewEncoder(&buffer).Encode(preregisterData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/preregister", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
	response := objects.ResponseString{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "email '"+preregisterData.Email+"' already exists", response.Message)
	dbconfig.DB.RemoveUser(context.Background(), dbUser.ApiKey)

	/* Test 3 - Valid request */
	preregisterData = PreRegisterPayload()
	err = json.NewEncoder(&buffer).Encode(preregisterData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/preregister?test=true", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

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
	/* Test 1 - Valid request*/
	dbconfig := setupDatabase("", "", "", false)
	router := setUpAPI(&dbconfig)

	registrationData := RegisterPayload()
	register_data_token := createDatabasePreregister(registrationData.Name, registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(registrationData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	response := objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
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
	new_registration_data := ProductPayload()
	err = json.NewEncoder(&buffer).Encode(new_registration_data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/register", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	response = objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "", response.Name)
	assert.Equal(t, "", response.Email)

	/* Test 3 - Invalid token */
	registrationData = RegisterPayload()
	register_data_token = createDatabasePreregister(registrationData.Name, registrationData.Email, &dbconfig)
	err = json.NewEncoder(&buffer).Encode(registrationData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/register", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

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

	registrationData = RegisterPayload()
	register_data_token = createDatabasePreregister(registrationData.Name, registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	err = json.NewEncoder(&buffer).Encode(registrationData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/register", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

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
	registrationData = RegisterPayload()
	register_data_token = createDatabasePreregister(registrationData.Name, registrationData.Email, &dbconfig)
	registrationData.Token = register_data_token.Token.String()
	registrationData.Name = ""
	registrationData.Email = ""
	err = json.NewEncoder(&buffer).Encode(registrationData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/register", &buffer)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ready", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	response_string := objects.ResponseString{}
	err := json.Unmarshal(w.Body.Bytes(), &response_string)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "OK", response_string.Message)

	/* Test 2 - Invalid database credentials */
	dbconfig = setupDatabase("test_user", "test_psw", "database_test", true)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/ready", nil)
	router.ServeHTTP(w, req)

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
	var buf bytes.Buffer
	mpw := multipart.NewWriter(&buf)
	mpw.Close()

	file, err := os.Open("./test_payloads/import/" + fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	defer file.Close()

	w, err := mpw.CreateFormFile(formKey, "./test_payloads/import/"+fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	if err := mpw.Close(); err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	if _, err := io.Copy(w, file); err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}

	return mpw, buf
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

/* Returns a product request body struct */
func ProductPayload() objects.RequestBodyProduct {
	fileBytes := payload("products")
	productData := objects.RequestBodyProduct{}
	err := json.Unmarshal(fileBytes, &productData)
	if err != nil {
		log.Println(err)
	}
	return productData
}

/* Returns a register request body struct */
func RegisterPayload() objects.RequestBodyRegister {
	fileBytes := payload("registration")
	registerData := objects.RequestBodyRegister{}
	err := json.Unmarshal(fileBytes, &registerData)
	if err != nil {
		log.Println(err)
	}
	return registerData
}

/* Returns a pre-registrater request body struct */
func PreRegisterPayload() objects.RequestBodyPreRegister {
	fileBytes := payload("preregister")
	preregData := objects.RequestBodyPreRegister{}
	err := json.Unmarshal(fileBytes, &preregData)
	if err != nil {
		log.Println(err)
	}
	return preregData
}

/* Returns a login request body struct */
func LoginPayload() objects.RequestBodyLogin {
	fileBytes := payload("login")
	loginData := objects.RequestBodyLogin{}
	err := json.Unmarshal(fileBytes, &loginData)
	if err != nil {
		log.Println(err)
	}
	return loginData
}

/*
Returns a byte array representing the file data that was read

Data is retrived from the project directory `test_payloads`
*/
func payload(object_type string) []byte {
	file, err := os.Open("./test_payloads/" + object_type + ".json")
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
		product := ProductPayload()
		productUUID, _, err := AddProduct(dbconfig, product)
		if err != nil {
			log.Println(err)
			return uuid.Nil
		}
		return productUUID
	}
	return uuid.Nil
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
		user, err := dbconfig.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      "test",
			UserType:  "app",
			Email:     "test@test.com",
			Password:  utils.RandStringBytes(20),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return database.User{}
		}
		return user
	}
	return user
}

/*
Creates a demo token in the database for registration
*/
func createDatabasePreregister(name, email string, dbconfig *DbConfig) database.RegisterToken {
	token, err := dbconfig.DB.GetTokenValidation(context.Background(), email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.RegisterToken{}
		}
	}
	if token.Token == uuid.Nil {
		token, err := dbconfig.DB.CreateToken(context.Background(), database.CreateTokenParams{
			ID:        uuid.New(),
			Name:      name,
			Email:     email,
			Token:     uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return database.RegisterToken{}
		}
		return token
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
