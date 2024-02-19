package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"objects"
	"os"
	"shopify"
	"testing"
	"time"
	"utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProductsRoute(t *testing.T) {
	/* Test 1 - Invalid authentication */
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
	productData := createDatabaseProduct(&dbconfig)
	defer dbconfig.DB.RemoveProductByCode(context.Background(), productData.ProductCode)
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
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
	productData := createDatabaseProduct(&dbconfig)
	defer dbconfig.DB.RemoveProductByCode(context.Background(), productData.ProductCode)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+productData.ID.String()+"?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLoginRoute(t *testing.T) {
	/* Test 1 - Invalid request - empty username/password */
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	/* Test 2 - Invalid request - no cookies sent with request */
	dbUser := createDatabaseUser(&dbconfig)
	shopifyConfig = shopify.InitConfigShopify()
	router = setUpAPI(&dbconfig, &shopifyConfig)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/logout?api_key="+dbUser.ApiKey, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPreregisterRoute(t *testing.T) {
	/* Test 1 - Invalid request (empty email) */
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

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
func createDatabaseProduct(dbconfig *DbConfig) database.Product {
	product, err := dbconfig.DB.GetProductByProductCode(context.Background(), "product_code")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.Product{}
		}
	}
	if product.ProductCode == "" {
		product, err := dbconfig.DB.CreateProduct(context.Background(), database.CreateProductParams{
			ID:          uuid.New(),
			ProductCode: "product_code",
			Active:      "1",
			Title:       utils.ConvertStringToSQL("test_title"),
			BodyHtml:    utils.ConvertStringToSQL("test_body_html"),
			Category:    utils.ConvertStringToSQL("test_category"),
			Vendor:      utils.ConvertStringToSQL("test_vendor"),
			ProductType: utils.ConvertStringToSQL("test_product_type"),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return database.Product{}
		}
		return product
	}
	return database.Product{}
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
		db_user = param_db_name
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
