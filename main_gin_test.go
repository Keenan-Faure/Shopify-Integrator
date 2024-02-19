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

func TestRegisterRoute(t *testing.T) {
	/* Test 1 - Valid request*/
	dbconfig := setupDatabase("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

	registration_data := RegisterPayload()
	register_data_token := createDatabasePreregister(registration_data.Name, registration_data.Email, &dbconfig)
	registration_data.Token = register_data_token.Token.String()
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(registration_data)
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
	registration_data = RegisterPayload()
	register_data_token = createDatabasePreregister(registration_data.Name, registration_data.Email, &dbconfig)
	err = json.NewEncoder(&buffer).Encode(registration_data)
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

	registration_data = RegisterPayload()
	register_data_token = createDatabasePreregister(registration_data.Name, registration_data.Email, &dbconfig)
	registration_data.Token = register_data_token.Token.String()
	err = json.NewEncoder(&buffer).Encode(registration_data)
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
	registration_data = RegisterPayload()
	register_data_token = createDatabasePreregister(registration_data.Name, registration_data.Email, &dbconfig)
	registration_data.Token = register_data_token.Token.String()
	registration_data.Name = ""
	registration_data.Email = ""
	err = json.NewEncoder(&buffer).Encode(registration_data)
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

/* Returns a product request body struct */
func RegisterPayload() objects.RequestBodyRegister {
	fileBytes := payload("registration")
	registerData := objects.RequestBodyRegister{}
	err := json.Unmarshal(fileBytes, &registerData)
	if err != nil {
		log.Println(err)
	}
	return registerData
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
Creates a demo user in the database
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
