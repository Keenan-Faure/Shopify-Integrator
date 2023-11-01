package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"objects"
	"strconv"
	"testing"
)

func TestShopifyVariantPricing(t *testing.T) {
	fmt.Println("Test 1 - Creating product and fetching existing price for Shopify, price tier not set")
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
	productData := objects.Product{}
	err = json.Unmarshal(respBody, &productData)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	price, err := dbconfig.ShopifyVariantPricing(productData.Variants[0], "default_price_tier")
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if price != "0.00" {
		t.Errorf("expected '0.00' but found: " + price)
	}
	fmt.Println("Test 2 - Fetching invalid - non existant - price for Shopify")
	price, err = dbconfig.ShopifyVariantPricing(productData.Variants[0], "Test")
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if price != "0.00" {
		t.Errorf("expected '0.00' but found: " + price)
	}
	UFetchHelperPost("products/"+productData.ID.String(), "DELETE", user.ApiKey, nil)
}
