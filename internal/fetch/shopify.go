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
	"time"

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
func (configShopify *ConfigShopify) AddProductShopify(shopifyProduct objects.ShopifyProduct) error {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(shopifyProduct)
	if err != nil {
		return err
	}
	res, err := configShopify.FetchHelper("products.json", http.MethodPost, &buffer)
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		return errors.New("unexpected http status code")
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	products := objects.ShopifyProductResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return err
	}
	// FIXME add IDs in response to database
	return nil
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
	if res.StatusCode != 200 {
		return errors.New("unexpected http status code")
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
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
func (configShopify *ConfigShopify) AddVariantShopify(variant objects.ShopifyVariant, id string) error {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(variant)
	if err != nil {
		return err
	}
	res, err := configShopify.FetchHelper("products/"+id+"/variants.json", http.MethodPost, &buffer)
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		return errors.New("unexpected http status code")
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	products := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return err
	}
	// FIXME add IDs in response to database
	return nil
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
	if res.StatusCode != 200 {
		return errors.New("unexpected http status code")
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	products := objects.ShopifyVariantResponse{}
	err = json.Unmarshal(respBody, &products)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Adds a product collection to Shopify
func (configShopify *ConfigShopify) AddCollectionShopify(collection string) error {
	return nil
}

// Checks if the product SKU exists on the website
func (configShopify *ConfigShopify) GetProductBySKU(sku string) (bool, error) {
	client := graphql.NewClient(configShopify.Url+"/graphql.json", nil)
	variables := map[string]any{
		"sku": graphql.String(sku),
	}
	var respData struct {
		ProductVariants struct {
			Edges []struct {
				Node struct {
					Sku string
				}
			}
		} `graphql:"productVariants(query: $sku, first: 1)"`
	}

	err := client.Query(context.Background(), &respData, variables)
	if err != nil {
		return false, err
	}
	for _, value := range respData.ProductVariants.Edges {
		if value.Node.Sku == sku {
			return true, nil
		}
	}
	return false, nil
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
	if api_key == "" || api_key[0:3] != "ck_" {
		return false
	}
	if api_password == "" || api_password[0:3] != "cs_" {
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
