package main

import (
	"context"
	"log"
	"objects"
	"shopify"
)

// Pushes a product to Shopify
func (dbconfig *DbConfig) PushProduct(configShopify *shopify.ConfigShopify, product objects.Product) {
	_, err := GetProductID(dbconfig, product.ProductCode)
	if err != nil {
		// TODO log error to something
		log.Println(err)
	}
	// if product_id != "" {
	// 	// If yes, then update (UpdateProductShopify)
	// 	// done
	// 	configShopify.UpdateProductShopify()
	// }

	// If product does not exist
	// Create the product on website (save IDs)
	// Add variants to website as well (save IDs)
	// Add Collection to website
	// done
}

// Pushes a variant to Shopify
func PushVariant() {
	// Check if variant ids exist internally
	// If yes, then update (UpdateVariantShopify)
	// done

	// check if the product exists on the website (getProductBySKU)
	// If it does then retrieve the IDs (and save them)
	// Then update variant using those IDs

	// If no, then create new variant on website under the respective product
	// retrieve the IDs to use in future updates and save id's
}

// Pushes all products in database to Shopify
func Syncronize() {
	// Retrieve all products from database in batches and process them
	// by (loop) products -> (loop) variants

	// TODO errors are logged?
}

// Returns a product id if found in the database
// otherwise an empty string
func GetProductID(dbconfig *DbConfig, product_code string) (string, error) {
	product_id, err := dbconfig.DB.GetPIDByProductCode(context.Background(), product_code)
	if err != nil {
		return "", err
	}
	if product_id.ProductCode == product_code &&
		len(product_id.ShopifyProductID) > 0 &&
		product_id.ShopifyProductID != "" {
		return product_id.ShopifyProductID, nil
	}
	return "", nil
}

// Returns a variant id if found in the database
// otherwise an empty string
func GetVariantID(dbconfig *DbConfig, sku string) (string, error) {
	variant_id, err := dbconfig.DB.GetVIDBySKU(context.Background(), sku)
	if err != nil {
		return "", err
	}
	if variant_id.Sku == sku &&
		len(variant_id.ShopifyVariantID) > 0 &&
		variant_id.ShopifyVariantID != "" {
		return variant_id.ShopifyVariantID, nil
	}
	return "", nil
}
