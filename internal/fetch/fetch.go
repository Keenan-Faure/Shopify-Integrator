package fetch

import (
	"context"
	"encoding/json"
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

// Checks if the product SKU exists on the website
func (shopifyConfig *ConfigShopify) GetProductBySKU(sku string) bool {
	client := graphql.NewClient(shopifyConfig.Url+"/graphql.json", nil)
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
		log.Println(err)
	}
	for _, value := range respData.ProductVariants.Edges {
		if value.Node.Sku == sku {
			return true
		}
	}
	return false
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

func (shopifyConfig *ConfigShopify) FetchProducts() (objects.ShopifyProducts, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, shopifyConfig.Url+"?limit="+PRODUCT_FETCH_LIMIT, nil)
	if err != nil {
		log.Println(err)
		return objects.ShopifyProducts{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return objects.ShopifyProducts{}, err
	}
	return products, nil
}
