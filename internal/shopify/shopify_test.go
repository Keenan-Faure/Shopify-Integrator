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
const MOCK_SHOPIFY_API_PSW = "shpat_92UHEYF927YR2"
const MOCK_SHOPIFY_API_VERSION = "2021-07"
const MOCK_SHOPIFY_STORE_NAME = "test-test"

const MOCK_SHOPIFY_WEBHOOK_ID = "47593067"

const MOCK_NGROK_WEBHOOK_URL = "https://f5fa-102-135-246-72.ngrok-free.app"

const MOCK_SHOPIFY_LOCATION_ID = 10293810823
const MOCK_SHOPIFY_INVENTORY_LEVEL_ID = 23087120381
const MOCK_INVENTORY_ITEM_ID = 1023781023

const MOCK_SHOPIFY_PRODUCT_ID = 1072481085
const MOCK_SHOPIFY_VARIANT_ID = 1070325083
const MOCK_SHOPIFY_COLLECTION_ID = 2039482049
const MOCK_SHOPIFY_CUSTOM_COLLECTION_ID = 1063001407
const MOCK_PRODUCT_SKU = "MOCK_PRODUCT_SKU"

func TestDeleteShopifyWebhook(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyLocations()
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, len(response.Locations), 5)
}

func TestGetShopifyLocation(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyLocationByID(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, response.Location.ID, 655441491)
	assert.Equal(t, response.Location.Name, "50 Rideau Street")
}

func TestGetShopifyInventoryLevel(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	response, err = shopifyConfig.GetShopifyInventoryLevel(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), fmt.Sprint(MOCK_INVENTORY_ITEM_ID))
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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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
	response, err = shopifyConfig.GetShopifyInventoryLevels(fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID), fmt.Sprint(MOCK_INVENTORY_ITEM_ID))
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
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - empty parameters
	response, err := shopifyConfig.AddLocationQtyShopify(0, 0, 0)
	if err == nil {
		t.Errorf("expected 'invalid location id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 0)
	assert.Equal(t, int(response.InventoryLevel.Available), 0)

	// Test Case 2 - 1 empty parameter
	response, _ = shopifyConfig.AddLocationQtyShopify(MOCK_SHOPIFY_LOCATION_ID, 0, 0)
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 23087120381)
	assert.Equal(t, int(response.InventoryLevel.Available), 5)

	// Test Case 3 - valid parameters
	response, err = shopifyConfig.AddLocationQtyShopify(MOCK_SHOPIFY_LOCATION_ID, MOCK_SHOPIFY_INVENTORY_LEVEL_ID, 2)
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.InventoryLevel.InventoryItemID), 23087120381)
	assert.Equal(t, int(response.InventoryLevel.LocationID), 10293810823)
	assert.Equal(t, int(response.InventoryLevel.Available), 5)
	assert.Equal(t, response.InventoryLevel.UpdatedAt, "2024-04-01T13:24:55-04:00")
}

func TestAddProductShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

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

func TestUpdateProductShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.UpdateProductShopify(objects.ShopifyProduct{}, "")
	if err == nil {
		t.Errorf("expected 'invalid product id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.Product.ID), 0)
	assert.Equal(t, response.Product.Vendor, "")
	assert.Equal(t, response.Product.ProductType, "")
	assert.Equal(t, len(response.Product.Variants), 0)

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.UpdateProductShopify(objects.ShopifyProduct{
		ShopifyProd: objects.ShopifyProd{
			Title:    "Burton Custom Freestyle 151",
			BodyHTML: "<strong>Good snowboard!</strong>",
			Vendor:   "Burton",
			Type:     "Snowboard",
			Status:   "draft",
		},
	}, fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Product.ID), 1072481085)
	assert.Equal(t, response.Product.ProductType, "Snowboard")
	assert.Equal(t, response.Product.BodyHTML, "<strong>Good snowboard!</strong>")
	assert.Equal(t, int(response.Product.Variants[0].ID), 1070325083)
}

func TestAddVariantShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.AddVariantShopify(objects.ShopifyVariant{
		ShopifyVar: objects.ShopifyVar{
			Price:   "1.00",
			Option1: "Yellow",
		},
	}, "")
	if err == nil {
		t.Errorf("expected 'invalid product id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.Variant.ID), 0)
	assert.Equal(t, int(response.Variant.ProductID), 0)

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.AddVariantShopify(objects.ShopifyVariant{
		ShopifyVar: objects.ShopifyVar{
			Price:   "1.00",
			Option1: "Yellow",
		},
	}, fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Variant.ID), 1070325074)
	assert.Equal(t, int(response.Variant.ProductID), 632910392)
	assert.Equal(t, response.Variant.Title, "Yellow")
}

func TestUpdateVariantShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.UpdateVariantShopify(objects.ShopifyVariant{
		ShopifyVar: objects.ShopifyVar{
			Price:   "1.00",
			Option1: "Yellow",
		},
	}, "")
	if err == nil {
		t.Errorf("expected 'invalid variant id not allowed' but found: 'nil'")
	}
	assert.Equal(t, response.Variant.Title, "")
	assert.Equal(t, response.Variant.InventoryManagement, "")

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.UpdateVariantShopify(objects.ShopifyVariant{
		ShopifyVar: objects.ShopifyVar{
			Price:   "1.00",
			Option1: "Yellow",
		},
	}, fmt.Sprint(MOCK_SHOPIFY_VARIANT_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Variant.ID), 1070325074)
	assert.Equal(t, int(response.Variant.ProductID), 632910392)
	assert.Equal(t, response.Variant.Title, "Yellow")
	assert.Equal(t, response.Variant.Sku, "")
	assert.Equal(t, response.Variant.InventoryManagement, "shopify")
}

func TestAddProductToCollectionShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.AddProductToCollectionShopify(0, MOCK_SHOPIFY_COLLECTION_ID)
	if err == nil {
		t.Errorf("expected 'invalid product id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response.Collect.CollectionID), 0)
	assert.Equal(t, int(response.Collect.ID), 0)

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.AddProductToCollectionShopify(MOCK_SHOPIFY_PRODUCT_ID, MOCK_SHOPIFY_COLLECTION_ID)
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response.Collect.CollectionID), 841564295)
	assert.Equal(t, int(response.Collect.ID), 1071559588)
}

func TestAddCustomCollectionShopify(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.AddCustomCollectionShopify("")
	if err == nil {
		t.Errorf("expected 'invalid product id not allowed' but found: 'nil'")
	}
	assert.Equal(t, int(response), 0)

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.AddCustomCollectionShopify("IPods")
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, int(response), MOCK_SHOPIFY_CUSTOM_COLLECTION_ID)
}

func TestGetShopifyCategories(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - valid request
	response, err := shopifyConfig.GetShopifyCategories()
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(response.CustomCollections), 3)
	assert.Equal(t, response.CustomCollections[0].Title, "IPods")
	assert.Equal(t, response.CustomCollections[1].Title, "IPods Two")
	assert.Equal(t, response.CustomCollections[2].Title, "Non Ipods")
}

func TestGetShopifyCategoryByProductID(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.GetShopifyCategoryByProductID("")
	if err == nil {
		t.Errorf("expected 'invalid product id not allowed' but found: 'nil'")
	}
	assert.Equal(t, len(response.CustomCollections), 0)

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.GetShopifyCategoryByProductID(fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID))
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, len(response.CustomCollections), 3)
	assert.Equal(t, response.CustomCollections[0].Title, "IPods")
	assert.Equal(t, response.CustomCollections[1].Title, "IPods Two")
	assert.Equal(t, response.CustomCollections[2].Title, "Non Ipods")
}

func TestFetchProducts(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	httpmock.Activate()
	InitMockShopifyAPI()
	defer httpmock.DeactivateAndReset()

	// Test Case 1 - valid request
	shopifyProducts, _, err := shopifyConfig.FetchProducts("")
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	assert.Equal(t, len(shopifyProducts.Products), 2)
	assert.Equal(t, shopifyProducts.Products[0].Title, "IPod Nano - 8GB")
	assert.Equal(t, shopifyProducts.Products[1].Variants[0].Title, "Black")
}

func TestCategoryExists(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)

	// Test Case 1 - category does not exist
	categories := CreateShopifCollectionsResponse("test-case-valid-custom-collections.json")
	product := objects.Product{
		Category: "Mock Category",
	}
	exists, category_id := shopifyConfig.CategoryExists(product, categories)
	assert.Equal(t, exists, false)
	assert.Equal(t, category_id, 0)

	// Test Case 2 - invalid category
	product.Category = ""
	exists, category_id = shopifyConfig.CategoryExists(product, categories)
	assert.Equal(t, exists, false)
	assert.Equal(t, category_id, 0)

	// Test Case 3 - category does exist
	product.Category = "IPods Two"
	exists, category_id = shopifyConfig.CategoryExists(product, categories)
	assert.Equal(t, exists, true)
	assert.Equal(t, category_id, 395646240)
}

func TestGetProductBySKU(t *testing.T) {
	shopifyConfig := InitConfigShopify(MOCK_SHOPIFY_API_URL)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost,
		MOCK_SHOPIFY_API_URL+"/graphql.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyGraphQLResponse("test-case-valid-graph-ql.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	// Test Case 1 - invalid parameter
	response, err := shopifyConfig.GetProductBySKU("")
	if err == nil {
		t.Errorf("expected 'invalid sku not allowed' but found: 'nil'")
	}
	assert.Equal(t, response.ProductID, "")
	assert.Equal(t, response.VariantID, "")

	// Test Case 2 - valid parameter
	response, err = shopifyConfig.GetProductBySKU(MOCK_PRODUCT_SKU)
	if err != nil {
		t.Errorf("expected 'nil' but found: :" + err.Error())
	}
	assert.Equal(t, response.ProductID, "1072481085")
	assert.Equal(t, response.VariantID, "1070325083")
}

func TestValidateConfigShopify(t *testing.T) {
	// Test Case 1 - invalid parameter (x1)
	valid := ValidateConfigShopify("", MOCK_SHOPIFY_API_KEY, MOCK_SHOPIFY_API_PSW)
	assert.Equal(t, valid, false)

	// Test Case 2 - invalid parameters (x2)
	valid = ValidateConfigShopify("", "", MOCK_SHOPIFY_API_PSW)
	assert.Equal(t, valid, false)

	// Test Case 3 - invalid parameter (x3)
	valid = ValidateConfigShopify("", "", "")
	assert.Equal(t, valid, false)

	// Test Case 4 - valid parameter
	valid = ValidateConfigShopify(MOCK_SHOPIFY_STORE_NAME, MOCK_SHOPIFY_API_KEY, MOCK_SHOPIFY_API_PSW)
	assert.Equal(t, valid, true)
}

/* Sets up a mock shopify API */
func InitMockShopifyAPI() {

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/variants/"+fmt.Sprint(MOCK_SHOPIFY_VARIANT_ID)+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifVariantResponse("test-case-valid-variant.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/products/count.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyProductCountResponse("test-case-valid-product-count.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyWebhooksResponse("test-case-valid-webhooks.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/webhooks.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyWebhookResponse("test-case-valid-webhook.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/webhooks.json?topic=orders/updated",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyWebhooksResponse("test-case-valid-webhooks.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/webhooks/"+MOCK_SHOPIFY_WEBHOOK_ID+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyWebhookResponse("test-case-valid-webhook.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodDelete, MOCK_SHOPIFY_API_URL+"/webhooks/"+MOCK_SHOPIFY_WEBHOOK_ID+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, "")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/collects.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyCollectionResponse("test-case-valid-collection.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/custom_collections.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifCustomCollectionResponse("test-case-valid-custom-collection.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/custom_collections.json?fields=title,id",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifCollectionsResponse("test-case-valid-custom-collections.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet,
		MOCK_SHOPIFY_API_URL+"/custom_collections.json?fields=title,id&product_id="+fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifCollectionsResponse("test-case-valid-custom-collections.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/products.json?limit="+PRODUCT_FETCH_LIMIT,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyProductsResponse("test-case-valid-products.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/products.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyProductResponse("test-case-valid-product.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/locations/"+fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyLocationResponse("test-case-valid-location.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/locations.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyLocationsResponse("test-case-valid-locations.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPut, MOCK_SHOPIFY_API_URL+"/products/"+fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID)+".json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyProductResponse("test-case-valid-product.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/products/"+fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID)+"/variants.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, CreateShopifyVariantResponse("test-case-valid-variant.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost,
		MOCK_SHOPIFY_API_URL+"/graphql.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, CreateShopifyGraphQLResponse("test-case-valid-graph-ql.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodGet, MOCK_SHOPIFY_API_URL+"/inventory_levels.json?location_ids="+
		fmt.Sprint(MOCK_SHOPIFY_LOCATION_ID)+"&inventory_item_ids="+fmt.Sprint(MOCK_INVENTORY_ITEM_ID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				CreateShopifyInventoryLevelsResponse("test-case-valid-inventory-levels.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/inventory_levels/connect.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				201,
				CreateShopifyInventoryItemConnectResponse("test-case-valid-inventory-item-connect.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder(http.MethodPost, MOCK_SHOPIFY_API_URL+"/inventory_levels/adjust.json",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(
				200,
				CreateShopifyInventoryItemAdjustResponse("test-case-valid-level-item-adjust.json"),
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

/* Returns a test shopify inventory item adjust response struct */
func CreateShopifyInventoryItemAdjustResponse(fileName string) objects.ResponseAddInventoryItem {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyInventoryLevel := objects.ResponseAddInventoryItem{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
}

/* Returns a test shopify graph ql response response struct */
func CreateShopifyGraphQLResponse(fileName string) objects.JSONResponseShopifyGraphQL {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyGraphQL := objects.JSONResponseShopifyGraphQL{}
	err := json.Unmarshal(fileBytes, &shopifyGraphQL)
	if err != nil {
		log.Println(err)
	}
	return shopifyGraphQL
}

/* Returns a test shopify product response struct */
func CreateShopifProductsResponse(fileName string) objects.ShopifyProducts {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyProducts := objects.ShopifyProducts{}
	err := json.Unmarshal(fileBytes, &shopifyProducts)
	if err != nil {
		log.Println(err)
	}
	return shopifyProducts
}

/* Returns a test shopify custom collections response struct */
func CreateShopifCollectionsResponse(fileName string) objects.ResponseGetCustomCollections {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyCustomCollections := objects.ResponseGetCustomCollections{}
	err := json.Unmarshal(fileBytes, &shopifyCustomCollections)
	if err != nil {
		log.Println(err)
	}
	return shopifyCustomCollections
}

/* Returns a test shopify custom collection response struct */
func CreateShopifCustomCollectionResponse(fileName string) objects.ResponseShopifyCustomCollection {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyCustomCollection := objects.ResponseShopifyCustomCollection{}
	err := json.Unmarshal(fileBytes, &shopifyCustomCollection)
	if err != nil {
		log.Println(err)
	}
	return shopifyCustomCollection
}

/* Returns a test shopify collection response struct */
func CreateShopifyCollectionResponse(fileName string) objects.ResponseAddProductToShopifyCollection {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyCollection := objects.ResponseAddProductToShopifyCollection{}
	err := json.Unmarshal(fileBytes, &shopifyCollection)
	if err != nil {
		log.Println(err)
	}
	return shopifyCollection
}

/* Returns a test shopify variant response struct */
func CreateShopifVariantResponse(fileName string) objects.ShopifyVariantResponse {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyVariant := objects.ShopifyVariantResponse{}
	err := json.Unmarshal(fileBytes, &shopifyVariant)
	if err != nil {
		log.Println(err)
	}
	return shopifyVariant
}

/* Returns a test shopify inventory item connect response struct */
func CreateShopifyInventoryItemConnectResponse(fileName string) objects.ResponseAddInventoryItemLocation {
	fileBytes := payload("./test_payloads/tests/inventory-item-connect/" + fileName)
	shopifyInventoryLevel := objects.ResponseAddInventoryItemLocation{}
	err := json.Unmarshal(fileBytes, &shopifyInventoryLevel)
	if err != nil {
		log.Println(err)
	}
	return shopifyInventoryLevel
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

/* Returns a test shopify product response struct */
func CreateShopifyProductsResponse(fileName string) objects.ShopifyProducts {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyProducts := objects.ShopifyProducts{}
	err := json.Unmarshal(fileBytes, &shopifyProducts)
	if err != nil {
		log.Println(err)
	}
	return shopifyProducts
}

/* Returns a test shopify product response struct */
func CreateShopifyVariantResponse(fileName string) objects.ShopifyVariantResponse {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyVariant := objects.ShopifyVariantResponse{}
	err := json.Unmarshal(fileBytes, &shopifyVariant)
	if err != nil {
		log.Println(err)
	}
	return shopifyVariant
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
func CreateShopifyLocationsResponse(fileName string) objects.ShopifyLocations {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyLocations := objects.ShopifyLocations{}
	err := json.Unmarshal(fileBytes, &shopifyLocations)
	if err != nil {
		log.Println(err)
	}
	return shopifyLocations
}

/* Returns a test shopify location response struct */
func CreateShopifyLocationResponse(fileName string) objects.ShopifyLocation {
	fileBytes := payload("./test_payloads/" + fileName)
	shopifyLocation := objects.ShopifyLocation{}
	err := json.Unmarshal(fileBytes, &shopifyLocation)
	if err != nil {
		log.Println(err)
	}
	return shopifyLocation
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
