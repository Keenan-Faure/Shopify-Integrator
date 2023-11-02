package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"objects"
	"shopify"
	"strconv"
	"testing"
	"time"
)

func SetUpShopify() shopify.ConfigShopify {
	return shopify.InitConfigShopify()
}

func TestPushProduct(t *testing.T) {
	fmt.Println("Test Case 1 - Push new product that does not exist on the website")
	dbconfig := SetUpDatabase()
	shopifyConfig := SetUpShopify()
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
	defer dbconfig.DB.RemoveProduct(context.Background(), productData.ID)
	// create queue item
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
	err = json.NewEncoder(&buffer).Encode(queue_item)
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

	// wait until the queue processes the queue item
	time.Sleep(10 * time.Second)

	// queue processed item, now check shopify if the data is correct of the product
	// we should have saved the values of the product internally
	product_id, err := dbconfig.DB.GetPIDByProductCode(context.Background(), productData.ProductCode)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if product_id.ProductCode != productData.ProductCode {
		t.Errorf("expected '" + productData.ProductCode + "' but found: " + err.Error())
	}
	if product_id.ShopifyProductID == "" || len(product_id.ShopifyProductID) == 0 {
		t.Errorf("unexpected product code found")
	}
	// checks shopify data
	res, err = shopifyConfig.FetchHelper("products/"+product_id.ShopifyProductID+".json", http.MethodGet, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("unexpected status code found: " + fmt.Sprint(res.StatusCode))
	}
	product_data := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &product_data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if product_data.Product.Title != productData.Title {
		t.Errorf("expected '" + productData.Title + "' but found: " + product_data.Product.Title)
	}
	if fmt.Sprint(product_data.Product.ID) != product_id.ShopifyProductID {
		t.Errorf("expected '" + product_id.ShopifyProductID + "' but found: " + fmt.Sprint(product_data.Product.ID))
	}

	fmt.Println("Test Case 2 - Push Product that does not exist on the website")
	// push the same product that was previously pushed...
	// create queue item
	queue_item = objects.RequestQueueItem{
		Type:        "product",
		Status:      "in-queue",
		Instruction: "update_product",
		Object: objects.RequestQueueItemProducts{
			SystemProductID: productData.ID.String(),
			SystemVariantID: "",
			Shopify: struct {
				ProductID string "json:\"product_id\""
				VariantID string "json:\"variant_id\""
			}{
				ProductID: product_id.ShopifyProductID,
				VariantID: "",
			},
		},
	}
	err = json.NewEncoder(&buffer).Encode(queue_item)
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

	// wait until the queue processes the queue item
	time.Sleep(10 * time.Second)

	// after it processed, you can check the data of the
	// ids internally and the data on shopify

	product_id_updated, err := dbconfig.DB.GetPIDByProductCode(context.Background(), productData.ProductCode)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if product_id_updated.ProductCode != productData.ProductCode {
		t.Errorf("expected '" + productData.ProductCode + "' but found: " + product_id_updated.ProductCode)
	}
	if product_id_updated.ShopifyProductID == "" || len(product_id_updated.ShopifyProductID) == 0 {
		t.Errorf("unexpected product code found")
	}
	// checks shopify data
	res, err = shopifyConfig.FetchHelper("products/"+product_id_updated.ShopifyProductID+".json", http.MethodGet, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("unexpected status code found: " + fmt.Sprint(res.StatusCode))
	}
	product_data_updated := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &product_data_updated)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if product_data_updated.Product.Title != productData.Title {
		t.Errorf("expected '" + productData.Title + "' but found: " + product_data_updated.Product.Title)
	}
	if fmt.Sprint(product_data_updated.Product.ID) != product_id_updated.ShopifyProductID {
		t.Errorf("expected '" + product_id_updated.ShopifyProductID + "' but found: " + fmt.Sprint(product_data_updated.Product.ID))
	}
}

func TestPushVariant(t *testing.T) {
	fmt.Println("Test Case 1 - Push new variant that does not exist on the website")
	dbconfig := SetUpDatabase()
	shopifyConfig := SetUpShopify()
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
	defer dbconfig.DB.RemoveProduct(context.Background(), productData.ID)
	// create queue item
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
	err = json.NewEncoder(&buffer).Encode(queue_item)
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

	// wait until the queue processes the queue item
	time.Sleep(10 * time.Second)

	// queue processed item, now check shopify if the data is correct of the product
	// we should have saved the values of the product internally
	product_id, err := dbconfig.DB.GetPIDByProductCode(context.Background(), productData.ProductCode)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	queue_item = objects.RequestQueueItem{
		Type:        "product_variant",
		Status:      "in-queue",
		Instruction: "add_variant",
		Object: objects.RequestQueueItemProducts{
			SystemProductID: productData.ID.String(),
			SystemVariantID: productData.Variants[0].ID.String(),
			Shopify: struct {
				ProductID string "json:\"product_id\""
				VariantID string "json:\"variant_id\""
			}{
				ProductID: product_id.ShopifyProductID,
				VariantID: "",
			},
		},
	}
	err = json.NewEncoder(&buffer).Encode(queue_item)
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
	time.Sleep(10 * time.Second)
	variant_id, err := dbconfig.DB.GetVIDBySKU(context.Background(), productData.Variants[0].Sku)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if variant_id.Sku != productData.Variants[0].Sku {
		t.Errorf("expected '" + productData.Variants[0].Sku + "' but found: " + variant_id.Sku)
	}
	if variant_id.ShopifyVariantID == "" || len(variant_id.ShopifyVariantID) == 0 {
		t.Errorf("unexpected variant code found")
	}
	// checks shopify data
	res, err = shopifyConfig.FetchHelper("variants/"+variant_id.ShopifyVariantID+".json", http.MethodGet, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("unexpected status code found: " + fmt.Sprint(res.StatusCode))
	}
	variant_data := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &variant_data)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if variant_data.Variant.Sku != productData.Variants[0].Sku {
		t.Errorf("expected '" + productData.Variants[0].Sku + "' but found: " + variant_data.Variant.Sku)
	}
	if fmt.Sprint(variant_data.Variant.ID) != variant_id.ShopifyVariantID {
		t.Errorf("expected '" + variant_id.ShopifyVariantID + "' but found: " + fmt.Sprint(variant_data.Variant.ID))
	}

	fmt.Println("Test Case 2 - Push Product that does not exist on the website")
	// push the same product that was previously pushed...
	// create queue item
	queue_item = objects.RequestQueueItem{
		Type:        "product_variant",
		Status:      "in-queue",
		Instruction: "update_variant",
		Object: objects.RequestQueueItemProducts{
			SystemProductID: productData.ID.String(),
			SystemVariantID: productData.Variants[0].ID.String(),
			Shopify: struct {
				ProductID string "json:\"product_id\""
				VariantID string "json:\"variant_id\""
			}{
				ProductID: product_id.ShopifyProductID,
				VariantID: variant_id.ShopifyVariantID,
			},
		},
	}
	err = json.NewEncoder(&buffer).Encode(queue_item)
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

	// wait until the queue processes the queue item
	time.Sleep(10 * time.Second)

	// after it processed, you can check the data of the
	// ids internally and the data on shopify

	variant_id_updated, err := dbconfig.DB.GetVIDBySKU(context.Background(), productData.Variants[0].Sku)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if variant_id_updated.ShopifyVariantID != variant_id.ShopifyVariantID {
		t.Errorf("expected '" + variant_id.ShopifyVariantID + "' but found: " + variant_id_updated.ShopifyVariantID)
	}
	if variant_id_updated.ShopifyVariantID == "" || len(variant_id_updated.ShopifyVariantID) == 0 {
		t.Errorf("unexpected variant code found")
	}
	// checks shopify data
	res, err = shopifyConfig.FetchHelper("variants/"+variant_id.ShopifyVariantID+".json", http.MethodGet, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("unexpected status code found: " + fmt.Sprint(res.StatusCode))
	}
	variant_data_updated := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &variant_data_updated)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if variant_data_updated.Variant.Sku != productData.Variants[0].Sku {
		t.Errorf("expected '" + productData.Variants[0].Sku + "' but found: " + variant_data_updated.Variant.Sku)
	}
	if fmt.Sprint(variant_data_updated.Variant.ID) != variant_id.ShopifyVariantID {
		t.Errorf("expected '" + variant_id.ShopifyVariantID + "' but found: " + fmt.Sprint(variant_data_updated.Variant.ID))
	}

	// check internal inventory ids
	// check internal inventory stored

	fmt.Println("Test Case 2 - Push variant that does not exist on the website")
}

func TestCalculateAvailableQuantity(t *testing.T) {

}

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
