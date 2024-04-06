package shopify

import (
	"net/http"
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

func TestDeleteShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodDelete, MOCK_SHOPIFY_API_URL+"webhooks/"+MOCK_SHOPIFY_WEBHOOK_ID+".json",
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
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, response, "")

	// Test Case 2 - "valid" webhook ID
	response, err = shopifyConfig.DeleteShopifyWebhook(MOCK_SHOPIFY_WEBHOOK_ID)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, response, "")
}
