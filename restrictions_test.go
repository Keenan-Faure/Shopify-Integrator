package main

import (
	"context"
	"integrator/internal/database"
	"objects"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchRestrictionsToMap(t *testing.T) {
	// Test 1 - invalid function params
	result := FetchRestrictionsToMap([]database.FetchRestriction{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	dbconfig := setupDatabase("", "", "", false)
	restrictions, _ := dbconfig.DB.GetFetchRestriction(context.Background())
	result = FetchRestrictionsToMap(restrictions)
	assert.NotEqual(t, len(result), 0)
}

func TestPushRestrictionsToMap(t *testing.T) {
	// Test 1 - invalid function params
	result := PushRestrictionsToMap([]database.PushRestriction{})
	assert.Equal(t, len(result), 0)

	// Test 2 - valid function params
	dbconfig := setupDatabase("", "", "", false)
	restrictions, _ := dbconfig.DB.GetPushRestriction(context.Background())
	result = PushRestrictionsToMap(restrictions)
	assert.NotEqual(t, len(result), 0)
}

func TestApplyFetchRestriction(t *testing.T) {
	// Test 1 - invalid function params
	result := ApplyFetchRestriction(make(map[string]string), "", "")
	assert.Equal(t, result, "")

	// Test 2 - valid function params | app
	mockRestrictionMap := make(map[string]string)
	mockRestrictionMap["title"] = "app"
	result = ApplyFetchRestriction(mockRestrictionMap, "MOCK_VALUE", "title")
	assert.Equal(t, result, "")

	// Test 3 - valid function params | shopify
	mockRestrictionMap["title"] = "shopify"
	result = ApplyFetchRestriction(mockRestrictionMap, "MOCK_VALUE", "title")
	assert.Equal(t, result, "MOCK_VALUE")
}

func TestDeterPushRestriction(t *testing.T) {
	// Test 1 - invalid function params
	result := DeterPushRestriction(make(map[string]string), "")
	assert.Equal(t, result, false)

	// Test 2 - valid function params | app
	mockRestrictionMap := make(map[string]string)
	mockRestrictionMap["title"] = "app"
	result = DeterPushRestriction(mockRestrictionMap, "title")
	assert.Equal(t, result, true)

	// Test 3 - valid function params | shopify
	mockRestrictionMap["title"] = "shopify"
	result = DeterPushRestriction(mockRestrictionMap, "title")
	assert.Equal(t, result, false)
}

func TestApplyPushRestrictionProduct(t *testing.T) {
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	shopifyProduct := ConvertProductToShopify(productPayload)

	// Test 1 - invalid function params
	result := ApplyPushRestrictionProduct(make(map[string]string), objects.ShopifyProduct{})
	assert.Equal(t, result.Title, "")

	// Test 2 - valid function params | app
	mockRestrictionMap := make(map[string]string)
	mockRestrictionMap["title"] = "app"
	result = ApplyPushRestrictionProduct(mockRestrictionMap, shopifyProduct)
	assert.Equal(t, result.Title, "product_title")

	// Test 3 - valid function params | shopify
	mockRestrictionMap["title"] = "shopify"
	result = ApplyPushRestrictionProduct(mockRestrictionMap, shopifyProduct)
	assert.Equal(t, result.Title, "")
}

func TestApplyPushRestrictionV(t *testing.T) {
	productPayload := InitMockProduct("test-case-valid-product-variable.json")
	shopifyVariant := ConvertVariantToShopify(productPayload.Variants[0])

	// Test 1 - invalid function params
	result := ApplyPushRestrictionV(make(map[string]string), objects.ShopifyVariant{})
	assert.Equal(t, result.Barcode, "")

	// Test 2 - valid function params | app
	mockRestrictionMap := make(map[string]string)
	mockRestrictionMap["barcode"] = "app"
	result = ApplyPushRestrictionV(mockRestrictionMap, shopifyVariant)
	assert.Equal(t, result.Barcode, "2347234-9824")

	// Test 3 - valid function params | shopify
	mockRestrictionMap["barcode"] = "shopify"
	result = ApplyPushRestrictionV(mockRestrictionMap, shopifyVariant)
	assert.Equal(t, result.Barcode, "")
}

func TestDeterFetchRestriction(t *testing.T) {
	// Test 1 - invalid function params
	result := DeterFetchRestriction(make(map[string]string), "")
	assert.Equal(t, result, true)

	// Test 2 - valid function params | app
	mockRestrictionMap := make(map[string]string)
	mockRestrictionMap["title"] = "app"
	result = DeterFetchRestriction(mockRestrictionMap, "title")
	assert.Equal(t, result, false)

	// Test 3 - valid function params | shopify
	mockRestrictionMap["title"] = "shopify"
	result = DeterFetchRestriction(mockRestrictionMap, "title")
	assert.Equal(t, result, true)
}
