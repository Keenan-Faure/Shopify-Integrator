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
	"objects"
	"os"
	"strconv"
	"testing"
	"time"
	"utils"

	"github.com/google/uuid"
)

func SetUpDatabase() DbConfig {
	dbCon, err := InitConn(utils.LoadEnv("docker_db_url") + utils.LoadEnv("database") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}
	return dbCon
}

func UFetchHelper(endpoint, method, auth string) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, "http://localhost:"+utils.LoadEnv("port")+"/api/"+endpoint, nil)
	if auth != "" {
		req.Header.Add("Authorization", "ApiKey "+auth)
	}
	if err != nil {
		log.Println(err)
		return &http.Response{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return &http.Response{}, err
	}
	return res, nil
}

func UFetchHelperPost(endpoint, method, auth string, body io.Reader) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, "http://localhost:"+utils.LoadEnv("port")+"/api/"+endpoint, body)
	if auth != "" {
		req.Header.Add("Authorization", "ApiKey "+auth)
	}
	if err != nil {
		log.Println(err)
		return &http.Response{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return &http.Response{}, err
	}
	return res, nil
}

func CreateOrdr() objects.RequestBodyOrder {
	file, err := os.Open("./test_payloads/order.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	orderData := objects.RequestBodyOrder{}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		fmt.Println(err)
	}
	return orderData
}

func CreateProd() objects.RequestBodyProduct {
	return objects.RequestBodyProduct{
		Title:          "TestProduct",
		BodyHTML:       "",
		Category:       "",
		Vendor:         "",
		ProductType:    "",
		Variants:       []objects.ProductVariant{{Sku: "Test", Option1: "", Option2: "", Option3: "", Barcode: "", VariantPricing: []objects.VariantPrice{{Name: "Test", Value: "0.00"}}, VariantQuantity: []objects.VariantQty{{Name: "Test", Value: 0}}, UpdatedAt: time.Time{}}},
		ProductOptions: []objects.ProductOptions{{Value: ""}},
	}
}

func CreateDemoUser(dbconfig *DbConfig) database.User {
	user, err := dbconfig.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "Demo",
		Email:     "Demo@test.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.User{}
	}
	return user
}

func TestDatabaseConnection(t *testing.T) {
	fmt.Println("Test Case 1 - Invalid database url string")
	dbconfig, err := InitConn("abc123")
	if err != nil && dbconfig.Valid {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	_, err = dbconfig.DB.GetOrders(context.Background(), database.GetOrdersParams{
		Limit:  1,
		Offset: 0,
	})
	if err == nil {
		t.Errorf("Expected 'error' but found 'nil'")
	}
	fmt.Println("Test Case 2 - Invalid database")
	dbconfig, err = InitConn(utils.LoadEnv("db_url") + "fake_abc123" + "?sslmode=disable")
	if err != nil && dbconfig.Valid {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	_, err = dbconfig.DB.GetOrders(context.Background(), database.GetOrdersParams{
		Limit:  1,
		Offset: 0,
	})
	if err == nil {
		t.Errorf("Expected 'error' but found 'nil'")
	}
	fmt.Println("Test Case 3 - Valid connection url")
	dbconfig, err = InitConn(utils.LoadEnv("db_url") + utils.LoadEnv("database") + "?sslmode=disable")
	if err != nil && !dbconfig.Valid {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	_, err = dbconfig.DB.GetOrders(context.Background(), database.GetOrdersParams{
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		t.Errorf("Expected 'nil' but found 'error'")
	}
}

func TestProductCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating product")
	dbconfig := SetUpDatabase()
	body := CreateProd()
	user := CreateDemoUser(&dbconfig)
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err := UFetchHelperPost("products", "POST", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 201 {
		t.Errorf("Expected '201' but found: " + strconv.Itoa(res.StatusCode))
	}
	productData := objects.Product{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if productData.Title != "TestProduct" {
		t.Errorf("Expected 'TestProduct' but found: " + productData.Title)
	}
	fmt.Println("Test 2 - Fetching product")
	res, err = UFetchHelper("products/"+productData.ID.String(), "GET", user.ApiKey)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	productData = objects.Product{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if productData.Title != "TestProduct" {
		t.Errorf("Expected 'TestProduct' but found: " + productData.Title)
	}

	fmt.Println("Test 3 - Deleting product & recheck")
	dbconfig.DB.RemoveProduct(context.Background(), productData.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("products/"+productData.ID.String(), "GET", user.ApiKey)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 404 {
		t.Errorf("Expected '404' but found: " + strconv.Itoa(res.StatusCode))
	}
	data := ErrorStruct{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if data.Error != "not found" {
		t.Errorf("Expected 'not found' but found: " + data.Error)
	}
	dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
}

func TestOrderCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating order")
	dbconfig := SetUpDatabase()
	body := CreateOrdr()
	user := CreateDemoUser(&dbconfig)
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err := UFetchHelperPost("orders?token="+user.WebhookToken, "POST", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 201 {
		t.Errorf("Expected '201' but found: " + strconv.Itoa(res.StatusCode))
	}
	orderData := objects.RequestString{}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	fmt.Println("Test 2 - Fetching order")
	res, err = UFetchHelper("orders/"+orderData.Message, "GET", user.ApiKey)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	order_id, err := uuid.Parse(orderData.Message)
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	orderData_fetched, err := CompileOrderData(&dbconfig, order_id, res.Request, true)
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if orderData_fetched.ID.String() != orderData.Message {
		t.Errorf("Expected '" + orderData_fetched.ID.String() + "' but found: " + orderData.Message)
	}

	fmt.Println("Test 3 - Deleting order & recheck")
	dbconfig.DB.RemoveOrder(context.Background(), orderData_fetched.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("orders/"+orderData.Message, "GET", user.ApiKey)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 404 {
		t.Errorf("Expected '404' but found: " + strconv.Itoa(res.StatusCode))
	}
	data := ErrorStruct{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if data.Error != "not found" {
		t.Errorf("Expected 'not found' but found: " + data.Error)
	}
	dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
}

func TestCustomerCRUD(t *testing.T) {

}

// import / export should also appear here
