package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"objects"
	"shopify"
	"testing"
	"utils"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	dbconfig := SetUpDatabase()
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
}

func SetUpDatabase() DbConfig {
	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	dbCon, err := InitConn(connection_string + "@127.0.0.1:5432/" + utils.LoadEnv("db_name") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}
	return dbCon
}
