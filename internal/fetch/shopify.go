package shopify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"objects"
	"time"
	"utils"

	"github.com/shurcooL/graphql"
)

const PRODUCT_FETCH_LIMIT = "50" // limit on products to fetch

type ConfigShopify struct {
	APIKey      string
	APIPassword string
	Version     string
	Url         string
	Valid       bool
}

// TODO log the fetch errors?
// Adds a product to Shopify
func (configShopify *ConfigShopify) AddProductShopify(shopifyProduct objects.ShopifyProduct) (objects.ShopifyProductResponse, error) {
	fmt.Println(shopifyProduct)
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
	return products, err
}

// Updates a product on Shopify
func (configShopify *ConfigShopify) UpdateProductShopify(shopifyProduct objects.ShopifyProduct, id string) error {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(shopifyProduct)
	if err != nil {
		return err
	}
	res, err := configShopify.FetchHelper("products/"+id+".json", http.MethodPut, &buffer)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(string(respBody))
	}
	products := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Adds a product variant on Shopify
func (configShopify *ConfigShopify) AddVariantShopify(variant objects.ShopifyVariant, product_id string) (string, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(variant)
	if err != nil {
		return "", err
	}
	fmt.Println(product_id)
	res, err := configShopify.FetchHelper("products/"+product_id+"/variants.json", http.MethodPost, &buffer)
	if err != nil {
		return "", err
	}
	fmt.Println(res)
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if res.StatusCode != 201 {
		return "", errors.New(string(respBody))
	}
	fmt.Println(string(respBody))
	variant_data := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &variant_data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprint(variant_data.Variant.ID), err
}

// Updates a product variant on Shopify
func (configShopify *ConfigShopify) UpdateVariantShopify(variant objects.ShopifyVariant, variant_id string) error {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(variant)
	if err != nil {
		return err
	}
	res, err := configShopify.FetchHelper("variants/"+variant_id+".json", http.MethodPut, &buffer)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(string(respBody))
	}
	products := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Adds a product to an existing collection in Shopify. Requires the Shopify product_id and the collection_id
// https://shopify.dev/docs/api/admin-rest/2023-10/resources/collect#post-collects
func (configShopify *ConfigShopify) AddProductToCollectionShopify(
	product_id,
	collection_id int) (objects.ResponseAddProductToShopifyCollection, error) {
	collection := objects.AddProducToShopifyCollection{
		Collect: struct {
			ProductID    int
			CollectionID int
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
		log.Println(err)
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
	shopify_collection := objects.AddShopifyCustomCollection{
		CustomCollection: struct{ Title string }{
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
	if res.StatusCode != 201 {
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
	client := graphql.NewClient(configShopify.Url+"/graphql.json", nil)
	variables := map[string]any{
		"sku": graphql.String(sku),
	}
	var respData struct {
		ProductVariants struct {
			Edges []struct {
				Node struct {
					Sku     string
					Id      string
					Product struct {
						Id string
					}
				}
			}
		} `graphql:"productVariants(query: $sku, first: 1)"`
	}
	err := client.Query(context.Background(), &respData, variables)
	if err != nil {
		return objects.ResponseIDs{}, err
	}
	for _, value := range respData.ProductVariants.Edges {
		if value.Node.Sku == sku {
			return objects.ResponseIDs{
				ProductID: utils.ExtractVID(respData.ProductVariants.Edges[0].Node.Id),
				VariantID: utils.ExtractPID(respData.ProductVariants.Edges[0].Node.Product.Id),
			}, nil
		}
	}
	return objects.ResponseIDs{}, nil
}

// Initiates the connection string for shopify
func InitConfigShopify(store_name, api_key, api_password, version string) ConfigShopify {
	validation := validateConfigShopify(store_name, api_key, api_password)
	if !validation {
		log.Println("Error setting up connection string for Shopify")
	}
	return ConfigShopify{
		APIKey:      api_key,
		APIPassword: api_password,
		Version:     version,
		Url:         "https://" + api_key + ":" + api_password + "@" + store_name + ".myshopify.com/admin/api/" + version,
		Valid:       validation,
	}
}

// Validates the config settings for Shopify
func validateConfigShopify(store_name, api_key, api_password string) bool {
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

func (configShopify *ConfigShopify) FetchProducts() (objects.ShopifyProducts, error) {
	res, err := configShopify.FetchHelper("products.json?limit="+PRODUCT_FETCH_LIMIT, http.MethodGet, nil)
	if err != nil {
		return objects.ShopifyProducts{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProducts{}, err
	}
	products := objects.ShopifyProducts{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err) // TODO Log these errors?
		return objects.ShopifyProducts{}, err
	}
	return products, nil
}

func (shopifyConfig *ConfigShopify) FetchHelper(endpoint, method string, body io.Reader) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, shopifyConfig.Url+"/"+endpoint, body)
	fmt.Println(shopifyConfig.Url + "/" + endpoint)
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
