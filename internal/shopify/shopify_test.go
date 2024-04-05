package shopify

import (
	"log"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const SHOPIFY_API_URL = "http://localhost:4711"

const MOCK_SHOPIFY_WEBHOOK_ID = "47593067"

func TestDeleteShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(SHOPIFY_API_URL)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response, err := shopifyConfig.DeleteShopifyWebhook(MOCK_SHOPIFY_WEBHOOK_ID)
	log.Println(err)
	httpmock.RegisterResponder(http.MethodGet, SHOPIFY_API_URL+"/api/tunnels",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}
