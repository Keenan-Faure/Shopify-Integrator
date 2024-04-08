package shopify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"objects"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/jarcoal/httpmock"
)

const MOCK_SHOPIFY_API_URL = "http://localhost:4711"
const MOCK_SHOPIFY_API_KEY = "9812Y3N13981UO1NWD"
const MOCK_SHOPIFY_API_PSW = "92UHEYF927YR2"
const MOCK_SHOPIFY_API_VERSION = "2021-07"
const MOCK_SHOPIFY_STORE_NAME = "test-test"

const MOCK_SHOPIFY_WEBHOOK_ID = "47593067"

const MOCK_NGROK_WEBHOOK_URL = "https://f5fa-102-135-246-72.ngrok-free.app"

const MOCK_SHOPIFY_LOCATION_ID = 10293810823
const MOCK_SHOPIFY_INVENTORY_LEVEL_ID = 23087120381

func TestDeleteShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodDelete, MOCK_SHOPIFY_API_URL+"/webhooks/"+MOCK_SHOPIFY_WEBHOOK_ID+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, "")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty webhook ID
	response, err := shopifyConfig.DeleteShopifyWebhook("")
	if err == nil {
		t.Errorf("expected 'Delete http://localhost:4711/webhooks/.json: no responder found' but found: 'nil'")
	}
	assert.Equal(t, response, "")

	// Test Case 2 - "valid" webhook ID
	response, err = shopifyConfig.DeleteShopifyWebhook(MOCK_SHOPIFY_WEBHOOK_ID)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, response, "")
}

func TestUpdateShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	webhookResponse := CreateShopifyWebhookResponse("test-case-valid-webhook.json")

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, webhookResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty webhook ID
	response, err := shopifyConfig.UpdateShopifyWebhook("", MOCK_NGROK_WEBHOOK_URL)
	if err == nil {
		t.Errorf("expected 'strconv.Atoi: parsing '': invalid syntax' but found: 'nil'")
	}
	assert.Equal(t, int(response.ShopifyWebhook.ID), 0)
	assert.Equal(t, response.ShopifyWebhook.Address, "")
	assert.Equal(t, response.ShopifyWebhook.Topic, "")
	assert.Equal(t, response.ShopifyWebhook.APIVersion, "")

	// Test Case 2 - "valid" webhook ID | invalid webhook URL
	response, err = shopifyConfig.UpdateShopifyWebhook(MOCK_SHOPIFY_WEBHOOK_ID, "")
	if err == nil {
		t.Errorf("expected 'invalid webhook url not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.ShopifyWebhook.ID), 0)
	assert.Equal(t, response.ShopifyWebhook.Address, "")
	assert.Equal(t, response.ShopifyWebhook.Topic, "")
	assert.Equal(t, response.ShopifyWebhook.APIVersion, "")

	// Test Case 3 - "valid" webhook ID
	response, err = shopifyConfig.UpdateShopifyWebhook(MOCK_SHOPIFY_WEBHOOK_ID, MOCK_NGROK_WEBHOOK_URL)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, int(response.ShopifyWebhook.ID), 4759306)
	assert.Equal(t, response.ShopifyWebhook.Address, "https://somewhere-else.com/")
	assert.Equal(t, response.ShopifyWebhook.Topic, "orders/create")
	assert.Equal(t, response.ShopifyWebhook.APIVersion, "unstable")
}

func TestCreateShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	webhookResponse := CreateShopifyWebhookResponse("test-case-valid-webhook.json")

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, webhookResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty NGROK webhook URL
	response, err := shopifyConfig.CreateShopifyWebhook("")
	if err == nil {
		t.Errorf("expected 'invalid webhook url not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.ShopifyWebhook.ID), 0)
	assert.Equal(t, response.ShopifyWebhook.Address, "")
	assert.Equal(t, response.ShopifyWebhook.Topic, "")
	assert.Equal(t, response.ShopifyWebhook.APIVersion, "")

	// Test Case 2 - "valid" webhook URL
	response, err = shopifyConfig.CreateShopifyWebhook(MOCK_NGROK_WEBHOOK_URL)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, int(response.ShopifyWebhook.ID), 4759306)
	assert.Equal(t, response.ShopifyWebhook.Address, "https://somewhere-else.com/")
	assert.Equal(t, response.ShopifyWebhook.Topic, "orders/create")
	assert.Equal(t, response.ShopifyWebhook.APIVersion, "unstable")
}

func TestGetShopifyWebhooks(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	webhookResponse := CreateShopifyWebhooksResponse("test-case-valid-webhooks.json")

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, webhookResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyWebhooks()
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, len(response.Webhooks), 4)
	assert.Equal(t, int(response.Webhooks[0].ID), 4759306)
	assert.Equal(t, response.Webhooks[0].Address, "https://apple.com")
	assert.Equal(t, response.Webhooks[0].Topic, "orders/create")
	assert.Equal(t, response.Webhooks[0].APIVersion, "unstable")
}

func TestGetShopifyProductCount(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	productCountResponse := CreateShopifyProductCountResponse("test-case-valid-product-count.json")

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/products/count.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, productCountResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyProductCount()
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Count), 2)
}

func TestGetShopifyLocations(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	shopifyLocationResponse := CreateShopifyLocationResponse("test-case-valid-locations.json")

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/locations.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, shopifyLocationResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyLocations()
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, len(response.Locations), 5)
}

func TestGetShopifyInventoryLevel(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	shopifyInventoryLevelResponse := CreateShopifyInventoryLevelsResponse("test-case-valid-inventory-levels.json")

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/inventory_levels.json?location_ids="+
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+"&inventory_item_ids="+fmt.Sprint(MOCK_SHOPIFY_INVENTORY_LEVEL_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, shopifyInventoryLevelResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty parameters
	response, err := shopifyConfig.GetShopifyInventoryLevel("", "")
	if err == nil {
		t.Errorf("expected 'invalid location id not allowed' but found: 'nil'")
	}
	assert.Equal(t, response.Available, 0)
	assert.Equal(t, response.InventoryItemID, 0)
	assert.Equal(t, response.LocationID, 0)

	// Test Case 2 - 1 empty parameter
	response, err = shopifyConfig.GetShopifyInventoryLevel(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), "")
	if err == nil {
		t.Errorf("expected 'invalid inventory item id not allowed' but found: 'nil'")
	}
	assert.Equal(t, response.Available, 0)
	assert.Equal(t, response.InventoryItemID, 0)
	assert.Equal(t, response.LocationID, 0)

	// Test Case 3 - valid parameters
	response, err = shopifyConfig.GetShopifyInventoryLevel(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), fmt.Sprint(MOCK_SHOPIFY_INVENTORY_LEVEL_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Available), 2)
	assert.Equal(t, int(response.InventoryItemID), 23087120381)
	assert.Equal(t, int(response.LocationID), 10293810823)
}

func TestGetShopifyInventoryLevels(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	shopifyInventoryLevelResponse := CreateShopifyInventoryLevelsResponse("test-case-valid-inventory-levels.json")

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/inventory_levels.json?location_ids="+
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+"&inventory_item_ids="+fmt.Sprint(MOCK_SHOPIFY_INVENTORY_LEVEL_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, shopifyInventoryLevelResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty parameters
	response, err := shopifyConfig.GetShopifyInventoryLevels("", "")
	if err == nil {
		t.Errorf("expected 'invalid location id not allowed' but found: 'nil'")
	}
	assert.Equal(t, len(response.InventoryLevels), 0)

	// Test Case 2 - 1 empty parameter
	response, err = shopifyConfig.GetShopifyInventoryLevels(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), "")
	if err == nil {
		t.Errorf("expected 'invalid inventory item id not allowed' but found: 'nil'")
	}
	assert.Equal(t, len(response.InventoryLevels), 0)

	// Test Case 3 - valid parameters
	response, err = shopifyConfig.GetShopifyInventoryLevels(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), fmt.Sprint(MOCK_SHOPIFY_INVENTORY_LEVEL_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, len(response.InventoryLevels), 4)
	assert.Equal(t, int(response.InventoryLevels[1].Available), 1)
	assert.Equal(t, int(response.InventoryLevels[1].InventoryItemID), 808950810)
	assert.Equal(t, int(response.InventoryLevels[1].LocationID), 655441491)
}

func TestAddLocationQtyShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	shopifyInventoryLevelResponse := CreateShopifyInventoryLevelResponse("test-case-valid-inventory-level.json")

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/inventory_levels/adjust.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, shopifyInventoryLevelResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - empty parameters
	response, err := shopifyConfig.AddLocationQtyShopify(0, 0, 0)
	if err == nil {
		t.Errorf("expected 'invalid location id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 0)
	assert.Equal(t, int(response.InventoryLevel.Available), 0)

	// Test Case 2 - 1 empty parameter
	response, err = shopifyConfig.AddLocationQtyShopify(MOCK_SHOPIFY_LOCATION_ID, 0, 0)
	if err == nil {
		t.Errorf("expected 'invalid inventory item id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 0)
	assert.Equal(t, int(response.InventoryLevel.Available), 0)

	// Test Case 3 - valid parameters
	response, err = shopifyConfig.AddLocationQtyShopify(MOCK_SHOPIFY_LOCATION_ID, MOCK_SHOPIFY_INVENTORY_LEVEL_ID, 2)
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 23087120381)
	assert.Equal(t, int(response.InventoryLevel.LocationID), 10293810823)
	assert.Equal(t, int(response.InventoryLevel.Available), 2)
	assert.Equal(t, response.InventoryLevel.UpdatedAt, "2024-04-01T13:24:55-04:00")
}

func TestAddProductShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	shopifyResponse := CreateShopifyProductResponse("test-case-valid-product.json")

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/products.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, shopifyResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	// Test Case 1 - valid parameter
	response, err := shopifyConfig.AddProductShopify(objects.ShopifyProduct{
		ShopifyProd: objects.ShopifyProd{
			Title:    "Burton Custom Freestyle 151",
			BodyHTML: "<strong>Good snowboard!</strong>",
			Vendor:   "Burton",
			Type:     "Snowboard",
			Status:   "draft",
		},
	})
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Product.ID), 1072481085)
	assert.Equal(t, response.Product.ProductType, "Snowboard")
	assert.Equal(t, response.Product.BodyHTML, "<strong>Good snowboard!</strong>")
	assert.Equal(t, int(response.Product.Variants[0].ID), 1070325083)
}

/* Returns a test shopify product response struct */
func CreateShopifyProductResponse(fileName string) objects.ShopifyProductResponse {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyProduct := objects.ShopifyProductResponse{}
	err := json.Unmarshal(fileBytes, &shopifyProduct)
	if err != nil {
		log.Println(err)
	}
	return shopifyProduct
}

/* Returns a test shopify inventory level response struct */
func CreateShopifyInventoryLevelResponse(fileName string) objects.ResponseAddInventoryItem {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyInventoryLevel := objects.ResponseAddInventoryItem{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a test shopify inventory levels response struct */
func CreateShopifyInventoryLevelsResponse(fileName string) objects.GetShopifyInventoryLevelsList {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyInventoryLevel := objects.GetShopifyInventoryLevelsList{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a test shopify location response struct */
func CreateShopifyLocationResponse(fileName string) objects.ShopifyLocations {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyLocations := objects.ShopifyLocations{}
	err := json.Unmarshal(fileBytes, &shopifyLocations)
	if err != nil {
		log.Println(err)
	}
	return shopifyLocations
}

/* Returns a test shopify product count struct */
func CreateShopifyProductCountResponse(fileName string) objects.ShopifyProductCount {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyProductCount := objects.ShopifyProductCount{}
	err := json.Unmarshal(fileBytes, &shopifyProductCount)
	if err != nil {
		log.Println(err)
	}
	return shopifyProductCount
}

/* Returns a test shopify webhook response struct */
func CreateShopifyWebhookResponse(fileName string) objects.ShopifyWebhookRequest {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyWebhookResponse := objects.ShopifyWebhookRequest{}
	err := json.Unmarshal(fileBytes, &shopifyWebhookResponse)
	if err != nil {
		log.Println(err)
	}
	return shopifyWebhookResponse
}

/* Returns a test shopify webhook response struct */
func CreateShopifyWebhooksResponse(fileName string) objects.ShopifyWebhookResponse {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyWebhookResponse := objects.ShopifyWebhookResponse{}
	err := json.Unmarshal(fileBytes, &shopifyWebhookResponse)
	if err != nil {
		log.Println(err)
	}
	return shopifyWebhookResponse
}

/*
Returns a byte array representing the file data that was read

Data is retrived from the project directory `test_payloads`
*/
func payload(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	return respBody
}
