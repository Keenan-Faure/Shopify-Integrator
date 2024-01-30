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
	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	dbCon, err := InitConn(connection_string + "@127.0.0.1:5432/" + utils.LoadEnv("db_name") + "?sslmode=disable")
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

func CreateCustmr() objects.RequestBodyCustomer {
	file, err := os.Open("./test_payloads/customer.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	customerData := objects.RequestBodyCustomer{}
	err = json.Unmarshal(respBody, &customerData)
	if err != nil {
		fmt.Println(err)
	}
	return customerData
}

func CreateQueueItemProduct(dbconfig *DbConfig, user database.User, product objects.RequestBodyProduct) objects.RequestQueueItem {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(product)
	if err != nil {
		fmt.Println(err)
	}
	res, err := UFetchHelperPost("products", "POST", user.ApiKey, &buffer)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 201 {
		fmt.Println("invalid response code")
	}
	productData := objects.Product{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		fmt.Println(err)
	}
	dbconfig.DB.RemoveProduct(context.Background(), productData.ID)
	queue_item := objects.RequestQueueItem{
		Type:        "product",
		Status:      "in-queue",
		Instruction: "add_product",
		Object: objects.RequestQueueItemProducts{
			SystemProductID: productData.ID.String(),
			SystemVariantID: "",
			Shopify: struct {
				ProductID string "json:\"product_id\""
				VariantID string "json:\"variant_id\""
			}{
				ProductID: "",
				VariantID: "",
			},
		},
	}
	return queue_item
}

func CreateQueueItemOrder(queue_type string) objects.RequestQueueItem {
	file, err := os.Open("./test_payloads//queue/queue_" + queue_type + ".json")
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
		Active:         "1",
		ProductCode:    "ABC123",
		Title:          "TestProduct",
		BodyHTML:       "",
		Category:       "",
		Vendor:         "",
		ProductType:    "",
		Variants:       []objects.RequestBodyVariant{{Sku: "Test", Option1: "", Option2: "", Option3: "", Barcode: "", VariantPricing: []objects.VariantPrice{{Name: "Selling Price", Value: "0.00"}}, VariantQuantity: []objects.VariantQty{}, UpdatedAt: time.Time{}}},
		ProductOptions: []objects.ProductOptions{{Value: ""}},
	}
}

func CreateDemoUser(dbconfig *DbConfig) database.User {
	user, err := dbconfig.DB.GetUserByEmail(context.Background(), "Demo@test.com")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return database.User{}
		}
	}
	if user.ApiKey == "" {
		user, err := dbconfig.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      "Demo",
			UserType:  "",
			Email:     "Demo@test.com",
			Password:  "",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return database.User{}
		}
		return user
	} else {
		return user
	}
}

func TestDatabaseConnection(t *testing.T) {
	fmt.Println("Test Case 1 - Invalid database url string")
	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	docker_url := connection_string + "@localhost:5432/"

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
	dbconfig, err = InitConn(docker_url + "fake_abc123" + "?sslmode=disable")
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
	dbconfig, err = InitConn(docker_url + utils.LoadEnv("db_name") + "?sslmode=disable")
	if err != nil && !dbconfig.Valid {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	_, err = dbconfig.DB.GetOrders(context.Background(), database.GetOrdersParams{
		Limit:  1,
		Offset: 0,
	})
	fmt.Println(err)
	if err != nil {
		t.Errorf("Expected 'nil' but found 'error'")
	}
}

func TestProductCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating product")
	dbconfig := SetUpDatabase()
	body := CreateProd()
	user := CreateDemoUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
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
	response_string := objects.ResponseString{}
	err = json.Unmarshal(respBody, &response_string)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if response_string.Message != "success" {
		t.Errorf("Expected 'success' but found: " + response_string.Message)
	}
	fmt.Println("Test 2 - Fetching product by search title param")
	res, err = UFetchHelper("products/search?q="+body.Title, "GET", user.ApiKey)
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
	productData := []objects.SearchProduct{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if len(productData) == 0 {
		t.Errorf("expected '1' but found: 0")
	}
	fmt.Println("Test 3 - Fetching product by ID ")
	res, err = UFetchHelper("products/"+productData[0].ID.String(), "GET", user.ApiKey)
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
	productData_id := objects.Product{}
	err = json.Unmarshal(respBody, &productData_id)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if productData_id.Title != "TestProduct" {
		t.Errorf("Expected 'TestProduct' but found: " + productData_id.Title)
	}

	fmt.Println("Test 4 - Deleting product & recheck")
	dbconfig.DB.RemoveProduct(context.Background(), productData_id.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("products/"+productData_id.ID.String(), "GET", user.ApiKey)
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
}

func TestOrderCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating order")
	dbconfig := SetUpDatabase()
	body := CreateOrdr()
	user := CreateDemoUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err := UFetchHelperPost("orders?token="+user.WebhookToken+"&api_key="+user.ApiKey, "POST", "", &buffer)
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
	queueData := objects.ResponseQueueItem{}
	err = json.Unmarshal(respBody, &queueData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	fmt.Println("Test 2 - Fetching order")
	res, err = UFetchHelper("queue/"+queueData.ID.String(), "GET", user.ApiKey)
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
	_, err = uuid.Parse(queueData.ID.String())
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	queueOrder_fetched := objects.ResponseQueueWorker{}
	err = json.Unmarshal(respBody, &queueOrder_fetched)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if queueOrder_fetched.ID != queueData.ID.String() {
		t.Errorf("Expected '" + queueOrder_fetched.ID + "' but found: " + queueData.ID.String())
	}
	fmt.Println("Test 3 - Deleting order & recheck")
	dbconfig.DB.RemoveOrder(context.Background(), queueData.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("orders/"+queueData.ID.String(), "GET", user.ApiKey)
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
}

func TestCustomerCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating customer")
	dbconfig := SetUpDatabase()
	body := CreateCustmr()
	user := CreateDemoUser(&dbconfig)
	defer dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err := UFetchHelperPost("customers", "POST", user.ApiKey, &buffer)
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
	customerData := objects.Customer{}
	err = json.Unmarshal(respBody, &customerData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	fmt.Println("Test 2 - Fetching customer")
	res, err = UFetchHelper("customers/"+customerData.ID.String(), "GET", user.ApiKey)
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
	customer_id, err := uuid.Parse(customerData.ID.String())
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	customerData_fetched, err := CompileCustomerData(&dbconfig, customer_id, context.Background(), true)
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	err = json.Unmarshal(respBody, &customerData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if customerData_fetched.ID.String() != customerData.ID.String() {
		t.Errorf("Expected '" + customerData_fetched.ID.String() + "' but found: " + customerData.ID.String())
	}

	fmt.Println("Test 3 - Deleting customer & recheck")
	dbconfig.DB.RemoveCustomer(context.Background(), customerData_fetched.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("customers/"+customerData.ID.String(), "GET", user.ApiKey)
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
}

func TestProductIOCRUD(t *testing.T) {
	fmt.Println("Test 1 - Importing products")
	dbconfig := SetUpDatabase()
	user := CreateDemoUser(&dbconfig)
	res, err := UFetchHelper("products/import?test=true", "POST", user.ApiKey)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	importResponse := objects.ImportResponse{}
	err = json.Unmarshal(respBody, &importResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if importResponse.FailCounter != 0 {
		t.Errorf("Expected '0', but found " + fmt.Sprint(importResponse.FailCounter))
	}
	variant, err := dbconfig.DB.GetVariantBySKU(context.Background(), "skubca")
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if variant.Sku != "skubca" {
		t.Errorf("Expected 'skubca', but found " + variant.Sku)
	}
	_, err = os.Open("test_import.csv")
	if err == nil {
		t.Errorf("Expected error but found nil")
	}
	dbconfig.DB.RemoveProductByCode(context.Background(), "grouper")
	fmt.Println("Test 2 - Exporting products")
	res, err = UFetchHelperPost("products/export?test=true", "GET", user.ApiKey, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	exportResponse := objects.ResponseString{}
	err = json.Unmarshal(respBody, &exportResponse)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if exportResponse.Message == "" {
		t.Errorf("Expected a file name, but found " + exportResponse.Message)
	}
}

func TestQueueCRUD(t *testing.T) {
	fmt.Println("Test 1 - Creating new items inside the queue")
	dbconfig := SetUpDatabase()
	user := CreateDemoUser(&dbconfig)
	body := CreateQueueItemOrder("add_order")
	body2 := CreateQueueItemProduct(&dbconfig, user, CreateProd())
	defer dbconfig.DB.RemoveUser(context.Background(), user.ApiKey)
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	var buffer2 bytes.Buffer
	err = json.NewEncoder(&buffer2).Encode(body2)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err := UFetchHelperPost("queue", "POST", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res2, err := UFetchHelperPost("queue", "POST", user.ApiKey, &buffer2)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	defer res2.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 201 {
		t.Errorf("Expected '201' but found: " + strconv.Itoa(res.StatusCode))
	}
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res2.StatusCode != 201 {
		t.Errorf("Expected '201' but found: " + strconv.Itoa(res2.StatusCode))
	}
	fmt.Println("Test 2 - Reading the data of the new queue items")
	res, err = UFetchHelperPost("queue?page=1", "GET", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	queueDataList := []objects.ResponseQueueItemFilter{}
	err = json.Unmarshal(respBody, &queueDataList)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if queueDataList[0].Instruction != "add_product" {
		t.Errorf("expected 'add_product' but found: " + queueDataList[0].Instruction)
	}
	if queueDataList[0].Status != "in-queue" {
		t.Errorf("expected 'in-queue' but found: " + queueDataList[0].Status)
	}
	if queueDataList[1].Instruction != "add_order" {
		t.Errorf("expected 'add_order' but found: " + queueDataList[1].Instruction)
	}
	if queueDataList[1].QueueType != "order" {
		t.Errorf("expected 'order' but found: " + queueDataList[1].QueueType)
	}
	fmt.Println("Test 3 - Processing queue item in the queue and check status")
	// depends on how often the worker runs
	// by default I set time for 10 seconds
	time.Sleep(10 * time.Second)
	fmt.Println("Test 4 - Delete queue item in the queue")
	body = CreateQueueItemOrder("add_order")
	err = json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	res, err = UFetchHelperPost("queue", "POST", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 201 {
		t.Errorf("Expected '201' but found: " + strconv.Itoa(res.StatusCode))
	}
	res, err = UFetchHelperPost("queue?instruction=add_order", "DELETE", user.ApiKey, &buffer)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected '200' but found: " + strconv.Itoa(res.StatusCode))
	}
	fmt.Println("Test 5 - check the queue view")
	res, err = UFetchHelperPost("queue/view", "GET", user.ApiKey, &buffer)
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
	queueCount := objects.ResponseQueueCount{}
	err = json.Unmarshal(respBody, &queueCount)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if queueCount.AddOrder != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.AddOrder))
	}
	if queueCount.AddProduct != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.AddProduct))
	}
	if queueCount.AddVariant != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.AddVariant))
	}
	if queueCount.UpdateOrder != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.UpdateOrder))
	}
	if queueCount.UpdateProduct != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.UpdateProduct))
	}
	if queueCount.UpdateVariant != 0 {
		t.Errorf("Expected '0' but found " + fmt.Sprint(queueCount.UpdateVariant))
	}
	defer UFetchHelperPost("queue?status=completed", "DELETE", user.ApiKey, &buffer)
	defer UFetchHelperPost("queue?instruction=add_product", "DELETE", user.ApiKey, &buffer)
}
