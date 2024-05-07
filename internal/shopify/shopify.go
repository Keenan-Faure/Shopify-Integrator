package shopify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"objects"
	"strconv"
	"time"
	"utils"

	"github.com/shurcooL/graphql"
)

const PRODUCT_FETCH_LIMIT = "20" // limit on products to fetch

type ConfigShopify struct {
	APIKey      string
	APIPassword string
	Version     string
	Url         string
	Valid       bool
}

// Deletes a webhook on Shopify
func (configShopify *ConfigShopify) DeleteShopifyWebhook(shopify_webhook_id string) (string, error) {
	res, err := configShopify.FetchHelper(
		"webhooks/"+shopify_webhook_id+".json",
		http.MethodDelete,
		nil,
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(string(respBody))
	}
	return "", nil
}

// Updates a webhook on Shopify
func (configShopify *ConfigShopify) UpdateShopifyWebhook(
	shopify_webhook_id,
	webhook_url string,
) (objects.ShopifyWebhookRequest, error) {
	int_webhook_id, err := strconv.Atoi(shopify_webhook_id)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	if webhook_url == "" || len(webhook_url) == 0 {
		return objects.ShopifyWebhookRequest{}, errors.New("invalid webhook url not allowed")
	}
	webhook := objects.ShopifyWebhookRequest{
		ShopifyWebhook: objects.ShopifyWebhook{
			ID:      int64(int_webhook_id),
			Address: webhook_url,
		},
	}
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(webhook)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	res, err := configShopify.FetchHelper(
		"webhooks.json",
		http.MethodPut,
		&buffer,
	)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyWebhookRequest{}, errors.New(string(respBody))
	}
	response := objects.ShopifyWebhookRequest{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	return response, nil
}

// Creates a webhook on Shopify
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/webhook#post-webhooks
func (configShopify *ConfigShopify) CreateShopifyWebhook(webhook_url string) (objects.ShopifyWebhookRequest, error) {
	if webhook_url == "" || len(webhook_url) == 0 {
		return objects.ShopifyWebhookRequest{}, errors.New("invalid webhook url not allowed")
	}
	webhook := objects.ShopifyWebhookRequest{
		ShopifyWebhook: objects.ShopifyWebhook{
			Address: webhook_url,
			Topic:   "orders/updated",
			Format:  "json",
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(webhook)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	res, err := configShopify.FetchHelper(
		"webhooks.json",
		http.MethodPost,
		&buffer,
	)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	if res.StatusCode != 201 {
		return objects.ShopifyWebhookRequest{}, errors.New(string(respBody))
	}
	response := objects.ShopifyWebhookRequest{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ShopifyWebhookRequest{}, err
	}
	return response, nil
}

// Retrieves a list of webhooks on Shopify
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/webhook#get-webhooks
func (configShopify *ConfigShopify) GetShopifyWebhooks() (objects.ShopifyWebhookResponse, error) {
	res, err := configShopify.FetchHelper(
		"webhooks.json?topic=orders/updated",
		http.MethodGet,
		nil,
	)
	if err != nil {
		return objects.ShopifyWebhookResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyWebhookResponse{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyWebhookResponse{}, errors.New(string(respBody))
	}
	response := objects.ShopifyWebhookResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ShopifyWebhookResponse{}, err
	}
	return response, nil
}

// Retrieves a count of products on Shopify
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/product#get-products-count
func (configShopify *ConfigShopify) GetShopifyProductCount() (objects.ShopifyProductCount, error) {
	res, err := configShopify.FetchHelper(
		"products/count.json",
		http.MethodGet,
		nil,
	)
	if err != nil {
		return objects.ShopifyProductCount{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyProductCount{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyProductCount{}, errors.New(string(respBody))
	}
	response := objects.ShopifyProductCount{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ShopifyProductCount{}, err
	}
	return response, nil
}

// Retrieves a list of locations
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/location#get-locations
func (configShopify *ConfigShopify) GetShopifyLocations() (objects.ShopifyLocations, error) {
	res, err := configShopify.FetchHelper(
		"locations.json",
		http.MethodGet,
		nil,
	)
	if err != nil {
		return objects.ShopifyLocations{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyLocations{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyLocations{}, errors.New(string(respBody))
	}
	response := objects.ShopifyLocations{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ShopifyLocations{}, err
	}
	return response, nil
}

// Retrieves a list of inventory levels
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/inventorylevel#get-inventory-levels
func (configShopify *ConfigShopify) GetShopifyInventoryLevel(
	location_id,
	inventory_item_id string) (objects.GetShopifyInventoryLevels, error) {
	if location_id == "" {
		return objects.GetShopifyInventoryLevels{}, errors.New("invalid location id not allowed")
	}
	if inventory_item_id == "" {
		return objects.GetShopifyInventoryLevels{}, errors.New("invalid inventory item id not allowed")
	}
	res, err := configShopify.FetchHelper(
		"inventory_levels.json?location_ids="+location_id+"&inventory_item_ids="+inventory_item_id,
		http.MethodGet,
		nil,
	)
	if err != nil {
		return objects.GetShopifyInventoryLevels{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.GetShopifyInventoryLevels{}, err
	}
	if res.StatusCode != 200 {
		return objects.GetShopifyInventoryLevels{}, errors.New(string(respBody))
	}
	response := objects.GetShopifyInventoryLevelsList{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.GetShopifyInventoryLevels{}, err
	}
	return response.InventoryLevels[0], nil
}

// Retrieves a list of inventory levels
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/inventorylevel#get-inventory-levels
func (configShopify *ConfigShopify) GetShopifyInventoryLevels(
	location_id,
	inventory_item_id string) (objects.GetShopifyInventoryLevelsList, error) {
	if location_id == "" {
		return objects.GetShopifyInventoryLevelsList{}, errors.New("invalid location id not allowed")
	}
	if inventory_item_id == "" {
		return objects.GetShopifyInventoryLevelsList{}, errors.New("invalid inventory item id not allowed")
	}
	res, err := configShopify.FetchHelper(
		"inventory_levels.json?location_ids="+location_id+"&inventory_item_ids="+inventory_item_id,
		http.MethodGet,
		nil,
	)
	if err != nil {
		return objects.GetShopifyInventoryLevelsList{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.GetShopifyInventoryLevelsList{}, err
	}
	if res.StatusCode != 200 {
		return objects.GetShopifyInventoryLevelsList{}, errors.New(string(respBody))
	}
	response := objects.GetShopifyInventoryLevelsList{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.GetShopifyInventoryLevelsList{}, err
	}
	return response, nil
}

// Adjusts the inventory level of an inventory item at a location
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/inventorylevel#post-inventory-levels-adjust
func (configShopify *ConfigShopify) AddLocationQtyShopify(
	location_id, inventory_item_id, qty int) (objects.ResponseAddInventoryItem, error) {
	inventory_adjustment := objects.AddInventoryItem{
		LocationID:          location_id,
		InventoryItemID:     inventory_item_id,
		AvailableAdjustment: qty,
	}
	if location_id == 0 {
		return objects.ResponseAddInventoryItem{}, errors.New("invalid location id not allowed")
	}
	if inventory_item_id == 0 {
		return objects.ResponseAddInventoryItem{}, errors.New("invalid inventory item id not allowed")
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(inventory_adjustment)
	if err != nil {
		return objects.ResponseAddInventoryItem{}, err
	}
	res, err := configShopify.FetchHelper("inventory_levels/adjust.json", http.MethodPost, &buffer)
	if err != nil {
		return objects.ResponseAddInventoryItem{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ResponseAddInventoryItem{}, err
	}
	if res.StatusCode != 200 {
		return objects.ResponseAddInventoryItem{}, errors.New(string(respBody))
	}
	response := objects.ResponseAddInventoryItem{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ResponseAddInventoryItem{}, err
	}
	return response, nil
}

// Connects an inventory item to a location:
// https://shopify.dev/docs/api/admin-rest/2023-04/resources/inventorylevel#post-inventory-levels-connect
func (configShopify *ConfigShopify) AddInventoryItemToLocation(
	location_id, inventory_item_id int) (objects.ResponseAddInventoryItemLocation, error) {
	inventory_level := objects.AddInventoryItemToLocation{
		LocationID:      location_id,
		InventoryItemID: inventory_item_id,
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(inventory_level)
	if err != nil {
		return objects.ResponseAddInventoryItemLocation{}, err
	}
	res, err := configShopify.FetchHelper("inventory_levels/connect.json", http.MethodPost, &buffer)
	if err != nil {
		return objects.ResponseAddInventoryItemLocation{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ResponseAddInventoryItemLocation{}, err
	}
	if res.StatusCode != 201 {
		return objects.ResponseAddInventoryItemLocation{}, errors.New(string(respBody))
	}
	response := objects.ResponseAddInventoryItemLocation{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return objects.ResponseAddInventoryItemLocation{}, err
	}
	return response, nil
}

// Adds a product to Shopify:
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/product#post-products
func (configShopify *ConfigShopify) AddProductShopify(shopifyProduct objects.ShopifyProduct) (objects.ShopifyProductResponse, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(shopifyProduct)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	res, err := configShopify.FetchHelper("products.json", http.MethodPost, &buffer)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	if res.StatusCode != 201 {
		return objects.ShopifyProductResponse{}, errors.New(string(respBody))
	}
	products := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	return products, nil
}

// Updates a product on Shopify:
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/product#put-products-product-id
func (configShopify *ConfigShopify) UpdateProductShopify(shopifyProduct objects.ShopifyProduct, product_id string) (objects.ShopifyProductResponse, error) {
	if product_id == "" || len(product_id) == 0 {
		return objects.ShopifyProductResponse{}, errors.New("invalid product id not allowed")
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(shopifyProduct)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	res, err := configShopify.FetchHelper("products/"+product_id+".json", http.MethodPut, &buffer)
	if err != nil {
		return objects.ShopifyProductResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProductResponse{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyProductResponse{}, errors.New(string(respBody))
	}
	products := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProductResponse{}, err
	}
	return products, nil
}

// Adds a product variant on Shopify:
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/product-variant#post-products-product-id-variants
func (configShopify *ConfigShopify) AddVariantShopify(
	variant objects.ShopifyVariant,
	product_id string) (objects.ShopifyVariantResponse, error) {
	if product_id == "" || len(product_id) == 0 {
		return objects.ShopifyVariantResponse{}, errors.New("invalid product id not allowed")
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(variant)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	res, err := configShopify.FetchHelper("products/"+product_id+"/variants.json", http.MethodPost, &buffer)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	if res.StatusCode != 201 {
		return objects.ShopifyVariantResponse{}, errors.New(string(respBody))
	}
	variant_data := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &variant_data)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	return variant_data, nil
}

// Updates a product variant on Shopify:
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/product-variant#put-variants-variant-id
func (configShopify *ConfigShopify) UpdateVariantShopify(
	variant any,
	variant_id string) (objects.ShopifyVariantResponse, error) {
	if variant_id == "" || len(variant_id) == 0 {
		return objects.ShopifyVariantResponse{}, errors.New("invalid variant id not allowed")
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(variant)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	res, err := configShopify.FetchHelper("variants/"+variant_id+".json", http.MethodPut, &buffer)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ShopifyVariantResponse{}, err
	}
	if res.StatusCode != 200 {
		return objects.ShopifyVariantResponse{}, errors.New(string(respBody))
	}
	variant_data := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &variant_data)
	if err != nil {
		log.Println(err)
		return objects.ShopifyVariantResponse{}, err
	}
	return variant_data, nil
}

// Adds a product to an existing collection in Shopify. Requires the Shopify product_id and the collection_id
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/collect#post-collects
func (configShopify *ConfigShopify) AddProductToCollectionShopify(
	product_id,
	collection_id int) (objects.ResponseAddProductToShopifyCollection, error) {
	if product_id == 0 {
		return objects.ResponseAddProductToShopifyCollection{}, errors.New("invalid product id not allowed")
	}
	if collection_id == 0 {
		return objects.ResponseAddProductToShopifyCollection{}, errors.New("invalid collection id not allowed")
	}
	collection := objects.AddProducToShopifyCollection{
		Collect: struct {
			ProductID    int "json:\"product_id\""
			CollectionID int "json:\"collection_id\""
		}{
			ProductID:    product_id,
			CollectionID: collection_id,
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(collection)
	if err != nil {
		return objects.ResponseAddProductToShopifyCollection{}, err
	}
	res, err := configShopify.FetchHelper("collects.json", http.MethodPost, &buffer)
	if err != nil {
		return objects.ResponseAddProductToShopifyCollection{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return objects.ResponseAddProductToShopifyCollection{}, err
	}
	if res.StatusCode != 201 {
		return objects.ResponseAddProductToShopifyCollection{}, errors.New(string(respBody))
	}
	response := objects.ResponseAddProductToShopifyCollection{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Println(err)
		return objects.ResponseAddProductToShopifyCollection{}, err
	}
	return response, nil
}

// Adds a custom collection to Shopify
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/customcollection#post-custom-collections
func (configShopify *ConfigShopify) AddCustomCollectionShopify(collection string) (int, error) {
	if collection == "" || len(collection) == 0 {
		return 0, errors.New("invalid collection id not allowed")
	}
	shopify_collection := objects.AddShopifyCustomCollection{
		CustomCollection: struct {
			Title string "json:\"title\""
		}{
			Title: collection,
		},
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(shopify_collection)
	if err != nil {
		return 0, err
	}
	res, err := configShopify.FetchHelper("custom_collections.json", http.MethodPost, &buffer)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if res.StatusCode != 201 {
		return 0, errors.New(string(respBody))
	}
	response := objects.ResponseShopifyCustomCollection{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return response.CustomCollection.ID, nil
}

// Retreives all custom categories from Shopify
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/customcollection#get-custom-collections
func (configShopify *ConfigShopify) GetShopifyCategories() (objects.ResponseGetCustomCollections, error) {
	res, err := configShopify.FetchHelper("custom_collections.json?fields=title,id", http.MethodGet, nil)
	if err != nil {
		return objects.ResponseGetCustomCollections{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.ResponseGetCustomCollections{}, err
	}
	if res.StatusCode != 200 {
		return objects.ResponseGetCustomCollections{}, errors.New(string(respBody))
	}
	response := objects.ResponseGetCustomCollections{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Println(err)
		return objects.ResponseGetCustomCollections{}, err
	}
	return response, nil
}

// Retrieves a product's collection from Shopify
// used for shopify_fetch.go
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/customcollection#get-custom-collections
func (configShopify *ConfigShopify) GetShopifyCategoryByProductID(product_id string) (objects.ResponseGetCustomCollections, error) {
	if product_id == "" || len(product_id) == 0 {
		return objects.ResponseGetCustomCollections{}, errors.New("invalid product id not allowed")
	}
	res, err := configShopify.FetchHelper("custom_collections.json?fields=title,id&product_id="+product_id, http.MethodGet, nil)
	if err != nil {
		return objects.ResponseGetCustomCollections{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.ResponseGetCustomCollections{}, err
	}
	if res.StatusCode != 200 {
		return objects.ResponseGetCustomCollections{}, errors.New(string(respBody))
	}
	response := objects.ResponseGetCustomCollections{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Println(err)
		return objects.ResponseGetCustomCollections{}, err
	}
	return response, nil
}

// Determines if a product's category already exists on Shopify
func (configShopify *ConfigShopify) CategoryExists(product objects.Product, categories objects.ResponseGetCustomCollections) (bool, int) {
	for _, value := range categories.CustomCollections {
		if product.Category == value.Title {
			return true, int(value.ID)
		}
	}
	return false, 0
}

// Checks if the product SKU exists on the website
func (configShopify *ConfigShopify) GetProductBySKU(sku string) (objects.ResponseIDs, error) {
	if sku == "" || len(sku) == 0 {
		return objects.ResponseIDs{}, errors.New("invalid sku not allowed")
	}
	client := graphql.NewClient(configShopify.Url+"/graphql.json", nil)
	variables := map[string]any{
		"sku": "sku:" + graphql.String(sku),
	}
	graphQL := objects.ResponseShopifyGraphQL{}
	err := client.Query(context.Background(), &graphQL, variables)
	if err != nil {
		return objects.ResponseIDs{}, err
	}
	for _, value := range graphQL.ProductVariants.Edges {
		if value.Node.Sku == sku {
			return objects.ResponseIDs{
				VariantID: utils.ExtractVID(graphQL.ProductVariants.Edges[0].Node.Id),
				ProductID: utils.ExtractPID(graphQL.ProductVariants.Edges[0].Node.Product.Id),
			}, nil
		}
	}
	return objects.ResponseIDs{}, nil
}

// Initiates the connection string for shopify
func InitConfigShopify(apiURL string) ConfigShopify {
	if apiURL != "" {
		return ConfigShopify{
			APIKey:      "",
			APIPassword: "",
			Version:     "",
			Url:         apiURL,
			Valid:       true,
		}
	}
	store_name := utils.LoadEnv("store_name")
	api_key := utils.LoadEnv("api_key")
	api_password := utils.LoadEnv("api_password")
	api_version := utils.LoadEnv("api_version")
	validation := ValidateConfigShopify(store_name, api_key, api_password)
	if !validation {
		log.Println("Error setting up connection string for Shopify")
	}
	return ConfigShopify{
		APIKey:      api_key,
		APIPassword: api_password,
		Version:     api_version,
		Url:         "https://" + api_key + ":" + api_password + "@" + store_name + ".myshopify.com/admin/api/" + api_version,
		Valid:       validation,
	}
}

// Validates the config settings for Shopify
func ValidateConfigShopify(store_name, api_key, api_password string) bool {
	if store_name == "" {
		return false
	}
	if api_key == "" {
		return false
	}
	if api_password == "" || api_password[0:6] != "shpat_" {
		return false
	}
	return true
}

func (configShopify *ConfigShopify) FetchProducts(fetch_url string) (objects.ShopifyProducts, string, error) {
	if fetch_url == "" {
		fetch_url = "products.json?limit=" + PRODUCT_FETCH_LIMIT
	}
	res, err := configShopify.FetchHelper(fetch_url, http.MethodGet, nil)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProducts{}, "", err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProducts{}, "", err
	}
	products := objects.ShopifyProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProducts{}, "", err
	}
	return products, string(res.Header.Get("Link")), nil
}

func (shopifyConfig *ConfigShopify) FetchHelper(endpoint, method string, body io.Reader) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, shopifyConfig.Url+"/"+endpoint, body)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return res, nil
}
