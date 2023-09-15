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
	return objects.RequestBodyOrder{
		ID:                    0,
		AdminGraphqlAPIID:     "",
		AppID:                 0,
		BrowserIP:             "",
		BuyerAcceptsMarketing: false,
		CancelReason:          nil,
		CancelledAt:           nil,
		CartToken:             nil,
		CheckoutID:            0,
		CheckoutToken:         "",
		ClientDetails: struct {
			AcceptLanguage any    "json:\"accept_language\""
			BrowserHeight  any    "json:\"browser_height\""
			BrowserIP      string "json:\"browser_ip\""
			BrowserWidth   any    "json:\"browser_width\""
			SessionHash    any    "json:\"session_hash\""
			UserAgent      string "json:\"user_agent\""
		}{},
		ClosedAt:             nil,
		Confirmed:            false,
		ContactEmail:         "",
		CreatedAt:            time.Time{},
		Currency:             "",
		CurrentSubtotalPrice: "",
		CurrentSubtotalPriceSet: struct {
			ShopMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"shop_money\""
			PresentmentMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"presentment_money\""
		}{},
		CurrentTotalDiscounts: "",
		CurrentTotalDiscountsSet: struct {
			ShopMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"shop_money\""
			PresentmentMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"presentment_money\""
		}{},
		CurrentTotalDutiesSet: nil,
		CurrentTotalPrice:     "",
		CurrentTotalPriceSet: struct {
			ShopMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"shop_money\""
			PresentmentMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"presentment_money\""
		}{},
		CurrentTotalTax: "",
		CurrentTotalTaxSet: struct {
			ShopMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"shop_money\""
			PresentmentMoney struct {
				Amount       string "json:\"amount\""
				CurrencyCode string "json:\"currency_code\""
			} "json:\"presentment_money\""
		}{},
		CustomerLocale: "",
		DeviceID:       nil,
		DiscountCodes: []struct {
			Code   string "json:\"code\""
			Amount string "json:\"amount\""
			Type   string "json:\"type\""
		}{},
		Email:                  "",
		EstimatedTaxes:         false,
		FinancialStatus:        "",
		FulfillmentStatus:      nil,
		Gateway:                "",
		LandingSite:            nil,
		LandingSiteRef:         nil,
		LocationID:             nil,
		MerchantOfRecordAppID:  nil,
		Name:                   "0123",
		Note:                   nil,
		NoteAttributes:         []any{},
		Number:                 0,
		OrderNumber:            0,
		OrderStatusURL:         "",
		OriginalTotalDutiesSet: nil,
		PaymentGatewayNames:    []string{},
		Phone:                  "",
		PresentmentCurrency:    "",
		ProcessedAt:            time.Time{},
		ProcessingMethod:       "",
		Reference:              "",
		ReferringSite:          nil,
		SourceIdentifier:       "",
		SourceName:             "",
		SourceURL:              nil,
		SubtotalPrice:          "",
		LineItems: []struct {
			ID                  int64  "json:\"id\""
			AdminGraphqlAPIID   string "json:\"admin_graphql_api_id\""
			FulfillableQuantity int    "json:\"fulfillable_quantity\""
			FulfillmentService  string "json:\"fulfillment_service\""
			FulfillmentStatus   any    "json:\"fulfillment_status\""
			GiftCard            bool   "json:\"gift_card\""
			Grams               int    "json:\"grams\""
			Name                string "json:\"name\""
			Price               string "json:\"price\""
			PriceSet            struct {
				ShopMoney struct {
					Amount       string "json:\"amount\""
					CurrencyCode string "json:\"currency_code\""
				} "json:\"shop_money\""
				PresentmentMoney struct {
					Amount       string "json:\"amount\""
					CurrencyCode string "json:\"currency_code\""
				} "json:\"presentment_money\""
			} "json:\"price_set\""
			ProductExists    bool   "json:\"product_exists\""
			ProductID        any    "json:\"product_id\""
			Properties       []any  "json:\"properties\""
			Quantity         int    "json:\"quantity\""
			RequiresShipping bool   "json:\"requires_shipping\""
			Sku              string "json:\"sku\""
			Taxable          bool   "json:\"taxable\""
			Title            string "json:\"title\""
			TotalDiscount    string "json:\"total_discount\""
			TotalDiscountSet struct {
				ShopMoney struct {
					Amount       string "json:\"amount\""
					CurrencyCode string "json:\"currency_code\""
				} "json:\"shop_money\""
				PresentmentMoney struct {
					Amount       string "json:\"amount\""
					CurrencyCode string "json:\"currency_code\""
				} "json:\"presentment_money\""
			} "json:\"total_discount_set\""
			VariantID                  any    "json:\"variant_id\""
			VariantInventoryManagement any    "json:\"variant_inventory_management\""
			VariantTitle               string "json:\"variant_title\""
			Vendor                     string "json:\"vendor\""
			TaxLines                   []struct {
				ChannelLiable bool   "json:\"channel_liable\""
				Price         string "json:\"price\""
				PriceSet      struct {
					ShopMoney struct {
						Amount       string "json:\"amount\""
						CurrencyCode string "json:\"currency_code\""
					} "json:\"shop_money\""
					PresentmentMoney struct {
						Amount       string "json:\"amount\""
						CurrencyCode string "json:\"currency_code\""
					} "json:\"presentment_money\""
				} "json:\"price_set\""
				Rate  float64 "json:\"rate\""
				Title string  "json:\"title\""
			} "json:\"tax_lines\""
			Duties              []any "json:\"duties\""
			DiscountAllocations []struct {
				Amount    string "json:\"amount\""
				AmountSet struct {
					ShopMoney struct {
						Amount       string "json:\"amount\""
						CurrencyCode string "json:\"currency_code\""
					} "json:\"shop_money\""
					PresentmentMoney struct {
						Amount       string "json:\"amount\""
						CurrencyCode string "json:\"currency_code\""
					} "json:\"presentment_money\""
				} "json:\"amount_set\""
				DiscountApplicationIndex int "json:\"discount_application_index\""
			} "json:\"discount_allocations\""
		}{},
	}
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
	orderData := objects.Order{}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if orderData.WebCode != "0123" {
		t.Errorf("Expected '0123' but found: " + orderData.WebCode)
	}
	fmt.Println("Test 2 - Fetching order")
	res, err = UFetchHelper("orders/"+orderData.ID.String(), "GET", user.ApiKey)
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
	orderData = objects.Order{}
	err = json.Unmarshal(respBody, &orderData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if orderData.WebCode != "0123" {
		t.Errorf("Expected '0123' but found: " + orderData.WebCode)
	}

	fmt.Println("Test 3 - Deleting order & recheck")
	dbconfig.DB.RemoveOrder(context.Background(), orderData.ID)
	type ErrorStruct struct {
		Error string `json:"error"`
	}
	res, err = UFetchHelper("orders/"+orderData.ID.String(), "GET", user.ApiKey)
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
