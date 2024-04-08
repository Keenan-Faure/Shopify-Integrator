package shopify

import (
	"encoding/json"
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
	// Test Case 1 - empty NGROK webhook URL
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

/* Returns a test shopify webhook response struct */
func CreateShopifyWebhookResponse(fileName string) objects.ShopifyWebhookRequest {
	fileBytes := payload("./test_payloads/webhooks/" + fileName)
	shopifyWebhookResponse := objects.ShopifyWebhookRequest{}
	err := json.Unmarshal(fileBytes, &shopifyWebhookResponse)
	if err != nil {
		log.Println(err)
	}
	return shopifyWebhookResponse
}

/* Returns a test shopify webhook response struct */
func CreateShopifyWebhooksResponse(fileName string) objects.ShopifyWebhookResponse {
	fileBytes := payload("./test_payloads/webhooks/" + fileName)
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
