package main

import (
	"fmt"
	"integrator/internal/database"
	"objects"
)

// Converts the []database.FetchRestriction slice into a map
func FetchRestrictionsToMap(fetch_map []database.FetchRestriction) map[string]string {
	mapp := make(map[string]string)
	for _, value := range fetch_map {
		mapp[value.Field] = value.Flag
	}
	return mapp
}

// Converts the []database.PushRestriction slice into a map
func PushRestrictionsToMap(push_map []database.PushRestriction) map[string]string {
	mapp := make(map[string]string)
	for _, value := range push_map {
		mapp[value.Field] = value.Flag
	}
	return mapp
}

// Determines whether a field should be updated or not
// by default if the `restriction_type“ value inside the `restrictions“
// map is set to `app` or not
// NOTE: the field will NOT be updated if it is set to app
func ApplyFetchRestriction(
	restrictions map[string]string,
	value,
	restriction_type string,
) string {
	// check which restriction is being referenced
	// fetch all restrictions from the database for that type
	// update product to contain restrictions
	// only if the restriction is set to 'app' do we
	// not update the 'app' value with the shopify value

	if map_value, ok := restrictions[restriction_type]; ok {
		if map_value == "app" {
			// do not update
			return ""
		} else {
			return value
		}
	} else {
		// if the key does not exist
		// we should update it
		return value
	}
}

// Determines whether a field should be updated or not
// by default if the `restriction_type` value inside the `restrictions`
// map is set to `app` or not
// NOTE: the field will be updated if it is set to app
func DeterPushRestriction(
	restrictions map[string]string,
	restriction_type string,
) bool {
	// check which restriction is being referenced
	// fetch all restrictions from the database for that type
	// update product to contain restrictions
	// only if the restriction is set to 'app' do we
	// include it in the object as a property value
	// but we return a boolean to indicate this

	fmt.Println(restrictions)
	if map_value, ok := restrictions[restriction_type]; ok {
		if map_value == "app" {
			return true
		}
	}
	return false
}

// Applies the restriction to a product
func ApplyPushRestrictionProduct(
	restrictions map[string]string,
	shopify_product objects.ShopifyProduct,
) objects.ShopifyProduct {
	shopify_product_new := objects.ShopifyProduct{}

	if DeterPushRestriction(restrictions, "title") {
		shopify_product_new.Title = shopify_product.Title
	}
	if DeterPushRestriction(restrictions, "body_html") {
		shopify_product_new.BodyHTML = shopify_product.BodyHTML
	}
	if DeterPushRestriction(restrictions, "product_type") {
		shopify_product_new.Type = shopify_product.Type
	}
	if DeterPushRestriction(restrictions, "vendor") {
		shopify_product_new.Vendor = shopify_product.Vendor
	}
	if DeterPushRestriction(restrictions, "options") {
		shopify_product_new.Options = shopify_product.Options
	}

	// TODO check product category (PUT not supported in the app)

	return shopify_product_new
}

func ApplyPushRestrictionV(
	restrictions map[string]string,
	shopify_variant objects.ShopifyProdVariant,
) objects.ShopifyProdVariant {
	shopify_variant_new := objects.ShopifyProdVariant{}

	if DeterPushRestriction(restrictions, "barcode") {
		shopify_variant_new.Barcode = shopify_variant.Barcode
	}
	if DeterPushRestriction(restrictions, "options") {
		shopify_variant_new.Option1 = shopify_variant.Option1
	}
	if DeterPushRestriction(restrictions, "options") {
		shopify_variant_new.Option2 = shopify_variant.Option2
	}
	if DeterPushRestriction(restrictions, "options") {
		shopify_variant_new.Option3 = shopify_variant.Option3
	}
	if DeterPushRestriction(restrictions, "pricing") {
		shopify_variant_new.Price = shopify_variant.Price
	}
	if DeterPushRestriction(restrictions, "pricing") {
		shopify_variant_new.CompareAtPrice = shopify_variant.CompareAtPrice
	}

	return shopify_variant_new
}

// Determines whether a field should be updated or not
// by default if a key does not exist in the database
// it updates the app value with the shopify value during a fetch
func DeterFetchRestriction(
	restriction_map map[string]string,
	restriction_type string,
) bool {
	if map_value, ok := restriction_map[restriction_type]; ok {
		if map_value == "app" {
			// do not update
			return false
		} else {
			return true
		}
	} else {
		// if the key does not exist
		// we should update it
		return true
	}
}
