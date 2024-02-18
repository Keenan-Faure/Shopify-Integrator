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
	dbconfig := setup_database("", "", "", false)
	shopifyConfig := shopify.InitConfigShopify()
	router := setUpAPI(&dbconfig, &shopifyConfig)

	registration_data := payload("registration")
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(registration_data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/register", &buffer)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	response := objects.ResponseRegister{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, "test", response.Name)
	dbconfig.DB.RemoveUser(context.Background(), response.ApiKey)

	/* Test 2 - Invalid request body */

	/* Test 3 - Invalid token */

	/* Test 4 - User already exist */

	/* Test 5 - Email already existing */

}

func TestReadyRoute(t *testing.T) {
	/* Test 1 - Valid database credentials */
	dbconfig := setup_database("", "", "", false)
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
	dbconfig = setup_database("test_user", "test_psw", "database_test", true)
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
func create_queue_item(queue_type string) objects.RequestQueueItem {
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

/*
Returns a struct of the respective object type.

Data is retrived from the project directory `test_payloads`
*/
func payload(object_type string) objects.RequestBodyProduct {
	file, err := os.Open("./test_payloads/" + object_type + ".json")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	productData := objects.RequestBodyProduct{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		log.Println(err)
	}
	return productData
}

/*
Creates a demo user in the database
*/
func create_database_user(dbconfig *DbConfig) database.User {
	user, err := dbconfig.DB.GetUserByEmail(context.Background(), "demo@test.com")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.User{}
		}
	}
	if user.ApiKey == "" {
		user, err := dbconfig.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      "demo",
			UserType:  "app",
			Email:     "demo@test.com",
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
Helper function to setup a database connection for tests.

Will only overwrite with params if `overwrite` is set to true
*/
func setup_database(param_db_user, param_db_psw, param_db_name string, overwrite bool) DbConfig {
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
