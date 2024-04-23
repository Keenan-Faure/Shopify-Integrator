package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integrator/internal/database"
	"log"
	"objects"
	"testing"
	"time"
	"utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCompileShopifyToSystemProduct(t *testing.T) {
	// Test 1 - invalid function params
	result := CompileShopifyToSystemProduct(
		objects.ShopifySingleProduct{},
		objects.ShopifyProductVariant{},
		make(map[string]string),
	)

	assert.Equal(t, result.ProductCode, "")
	assert.Equal(t, result.Category, "")
	assert.Equal(t, result.Title, "")
	assert.Equal(t, result.Vendor, "")

	// Test 2 - valid function params
	restrictions := make(map[string]string)
	restrictions["title"] = "app"
	result = CompileShopifyToSystemProduct(
		CreateShopifySingleProduct("test-case-valid-single-product.json"),
		CreateShopifySingleProduct("test-case-valid-single-product.json").Variants[0],
		restrictions,
	)

	assert.Equal(t, result.ProductCode, "")
	assert.Equal(t, result.Category, "")
	assert.Equal(t, result.Title, "")
	assert.Equal(t, result.Vendor, "Burton")
}

func TestCompileRemoveQueueFilter(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileRemoveQueueFilter(&dbconfig, context.Background(), "", "", "")
	assert.Equal(t, result, "success")
	assert.Equal(t, err, nil)

	// Test 2 - valid function params
	result, err = CompileRemoveQueueFilter(&dbconfig, context.Background(), "", "in-queue", "products")
	assert.Equal(t, result, "success")
	assert.Equal(t, err, nil)
}

func TestConvertDatabaseToRegister(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result := ConvertDatabaseToRegister(database.User{})
	assert.Equal(t, result.ApiKey, "")
	assert.Equal(t, result.Email, "")
	assert.Equal(t, result.Name, "")

	// Test 2 - valid function params
	result = ConvertDatabaseToRegister(createDatabaseUser(&dbconfig))
	defer dbconfig.DB.RemoveUser(context.Background(), result.ApiKey)
	assert.Equal(t, result.Email, "test@test.com")
	assert.Equal(t, result.Name, "test")
}

func TestConvertDatabaseToWarehouse(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result := ConvertDatabaseToWarehouse([]database.GetWarehousesRow{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	createDatabaseGlobalWarehouse(&dbconfig)
	dbWarehouse, _ := dbconfig.DB.GetWarehouses(context.Background(), database.GetWarehousesParams{
		Limit:  5,
		Offset: 0,
	})
	result = ConvertDatabaseToWarehouse(dbWarehouse)
	assert.NotEqual(t, len(result), 0)
	RemoveGlobalWarehouse(&dbconfig, context.Background(), MOCK_WAREHOUSE_NAME)
}

func TestCompileQueueFilterSearch(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileQueueFilterSearch(&dbconfig, context.Background(), 0, "", "", "")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	CreateDatabaseQueueItem(&dbconfig, "in-queue")
	result, err = CompileQueueFilterSearch(&dbconfig, context.Background(), 1, "", "in-queue", "")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, len(result), 0)
	ClearQueueMockData(&dbconfig)
}

func TestConvertProductToShopify(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertProductToShopify(objects.Product{})
	assert.Equal(t, result.ShopifyProd.Title, "")
	assert.Equal(t, result.ShopifyProd.Type, "")

	// Test 2 - valid function params
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	result = ConvertProductToShopify(productPayload)
	assert.Equal(t, result.ShopifyProd.Title, "product_title")
	assert.Equal(t, result.ShopifyProd.Type, "product_product_type")
}

func TestConvertVariantToShopifyProdVariant(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertVariantToShopifyProdVariant(objects.Product{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	result = ConvertVariantToShopifyProdVariant(productPayload)
	assert.NotEqual(t, len(result), 0)
}

func TestConvertToShopifyIDs(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertToShopifyIDs(objects.ShopifyProductResponse{})
	assert.Equal(t, result.ProductID, "0")
	assert.Equal(t, len(result.Variants), 0)

	// Test 2 - valid function params
	result = ConvertToShopifyIDs(CreateShopifyProductResponse("test-case-valid-product.json"))
	assert.Equal(t, result.ProductID, fmt.Sprint(MOCK_SHOPIFY_PRODUCT_ID))
	assert.NotEqual(t, len(result.Variants), 0)
}

func TestConvertVariantToShopify(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertVariantToShopify(objects.ProductVariant{})
	assert.Equal(t, result.ShopifyVar.Sku, "")
	assert.Equal(t, result.ShopifyVar.Price, "0")
	assert.Equal(t, result.ShopifyVar.CompareAtPrice, "0")

	// Test 2 - valid function params
	result = ConvertVariantToShopify(InitMockProduct("test-case-valid-product-variable.json").Variants[0])
	assert.Equal(t, result.ShopifyVar.Sku, MOCK_PRODUCT_SKU)
	assert.Equal(t, result.ShopifyVar.Price, "0")
	assert.Equal(t, result.ShopifyVar.CompareAtPrice, "0")
}

func TestConvertVariantToShopifyVariant(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertVariantToShopifyVariant(objects.ProductVariant{})
	assert.Equal(t, result.Sku, "")
	assert.Equal(t, result.Price, "0")
	assert.Equal(t, result.ID, int64(0))

	// Test 2 - valid function params
	result = ConvertVariantToShopifyVariant(InitMockProduct("test-case-valid-product-variable.json").Variants[0])
	assert.Equal(t, result.Sku, MOCK_PRODUCT_SKU)
	assert.Equal(t, result.Price, "0")
	assert.Equal(t, result.ID, int64(0))
	assert.Equal(t, result.Barcode, "2347234-9824")
	assert.Equal(t, result.InventoryManagement, "")
}

func TestCompileShopifyOptions(t *testing.T) {
	// Test 1 - invalid function params
	result := CompileShopifyOptions(objects.Product{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	result = CompileShopifyOptions(InitMockProduct("test-case-valid-product-variable.json"))
	assert.NotEqual(t, len(result), 0)
}

func TestCompileCustomerData(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileCustomerData(&dbconfig, uuid.Nil, context.Background(), true)
	assert.Equal(t, result.ID, uuid.Nil)
	assert.NotEqual(t, err, nil)

	// Test 2 - valid function params | ignore_address
	customerUUID := createDatabaseCustomer(&dbconfig)
	defer dbconfig.DB.RemoveCustomer(context.Background(), customerUUID)
	result, err = CompileCustomerData(&dbconfig, customerUUID, context.Background(), true)
	assert.NotEqual(t, result.ID, uuid.Nil)
	assert.Equal(t, result.FirstName, "TestFirstName")
	assert.Equal(t, err, nil)

	// Test 3 - valid function params
	result, err = CompileCustomerData(&dbconfig, customerUUID, context.Background(), false)
	assert.NotEqual(t, result.ID, uuid.Nil)
	assert.Equal(t, result.FirstName, "TestFirstName")
	assert.NotEqual(t, len(result.Address), 0)
	assert.Equal(t, err, nil)
}

func TestCompileOrderSearchResult(t *testing.T) {
	// Test 1 - invalid function params
	result := CompileOrderSearchResult([]database.GetOrdersSearchByCustomerRow{}, []database.GetOrdersSearchWebCodeRow{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	result = CompileOrderSearchResult([]database.GetOrdersSearchByCustomerRow{
		{
			ID:            uuid.New(),
			Notes:         utils.ConvertStringToSQL("note"),
			Status:        "TestStatus",
			WebCode:       "Test",
			TaxTotal:      utils.ConvertStringToSQL("0"),
			OrderTotal:    utils.ConvertStringToSQL("123"),
			ShippingTotal: utils.ConvertStringToSQL("12"),
			DiscountTotal: utils.ConvertStringToSQL("0"),
			UpdatedAt:     time.Now().UTC(),
		},
	}, []database.GetOrdersSearchWebCodeRow{})
	assert.NotEqual(t, len(result), 0)
}

func TestCompileOrderData(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	// Test 1 - invalid function params
	result, err := CompileOrderData(&dbconfig, uuid.Nil, context.Background(), true)
	assert.Equal(t, result.Notes, "")
	assert.NotEqual(t, err, nil)

	// Test 2 - valid function params
	orderUUID := createDatabaseOrder(&dbconfig)
	defer dbconfig.DB.RemoveOrder(context.Background(), orderUUID)
	result, err = CompileOrderData(&dbconfig, orderUUID, context.Background(), true)
	assert.Equal(t, result.Notes, "Notes not taken")
	assert.Equal(t, err, nil)
}

func TestCompileFilterSearch(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileFilterSearch(&dbconfig, context.Background(), 0, "", "", "")
	assert.Equal(t, len(result), 0)
	assert.NotEqual(t, err, nil)

	// Test 2 - valid function params
	createDatabaseProduct(&dbconfig)
	defer ClearProductTestData(&dbconfig)
	result, err = CompileFilterSearch(&dbconfig, context.Background(), 1, "", "", "product_vendor")
	assert.NotEqual(t, len(result), 0)
	assert.Equal(t, err, nil)
}

func TestCompileProductImages(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileProductImages(uuid.Nil, context.Background(), &dbconfig)
	assert.Equal(t, len(result), 0)
	assert.Equal(t, err, nil)

	// Test 2 - valid function params
	productUUID := createDatabaseProduct(&dbconfig)
	result, err = CompileProductImages(productUUID, context.Background(), &dbconfig)
	assert.Equal(t, len(result), 0)
	assert.Equal(t, err, nil)
}

func TestCompileSearchResult(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid function params
	result, err := CompileSearchResult(&dbconfig, context.Background(), []database.GetProductsSearchRow{})
	assert.Equal(t, len(result), 0)
	assert.Equal(t, err, nil)

	// Test 2 - valid function params
	productUUID := createDatabaseProduct(&dbconfig)
	result, err = CompileSearchResult(&dbconfig, context.Background(), []database.GetProductsSearchRow{
		{
			ID:          productUUID,
			Active:      "",
			ProductCode: "ProductCode",
			Title:       utils.ConvertStringToSQL("product_title"),
			Category:    utils.ConvertStringToSQL("product_category"),
			Vendor:      utils.ConvertStringToSQL("product_vendor"),
			ProductType: utils.ConvertStringToSQL("product_product_type"),
			UpdatedAt:   time.Now().UTC(),
		},
	})
	assert.Equal(t, len(result), 1)
	assert.Equal(t, err, nil)
}

func TestConvertProductToCSVProduct(t *testing.T) {
	// Test 1 - invalid function params
	result := ConvertProductToCSVProduct(objects.RequestBodyProduct{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	result = ConvertProductToCSVProduct(ProductPayload("test-case-valid-product-variable.json"))
	assert.NotEqual(t, len(result), 0)
	assert.Equal(t, result[0].Title, "product_title")
}

func TestCompileProduct(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	// Test 1 - invalid function params
	result, err := CompileProduct(&dbconfig, uuid.Nil, context.Background(), false)
	assert.Equal(t, result.Title, "")
	assert.NotEqual(t, err, nil)

	// Test 2 - valid function params
	productUUID := createDatabaseProduct(&dbconfig)
	defer ClearProductTestData(&dbconfig)
	result, err = CompileProduct(&dbconfig, productUUID, context.Background(), false)
	assert.Equal(t, result.Title, "product_title")
	assert.Equal(t, err, nil)
}

func TestCompileVariants(t *testing.T) {
	// Test 1 - invalid function params

	// Test 2 - valid function params
}

func TestCompileVariantByID(t *testing.T) {
	// Test 1 - invalid function params

	// Test 2 - valid function params
}

/* Returns a test shopify single product response struct */
func CreateShopifySingleProduct(fileName string) objects.ShopifySingleProduct {
	/* Returns a test shopify product response struct */
	fileBytes := payload("./test_payloads/tests/shopify-product/" + fileName)
	shopifyProduct := objects.ShopifySingleProduct{}
	err := json.Unmarshal(fileBytes, &shopifyProduct)
	if err != nil {
		log.Println(err)
	}
	return shopifyProduct
}
