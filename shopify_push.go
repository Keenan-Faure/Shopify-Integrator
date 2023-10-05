package main

import (
	"context"
	"integrator/internal/database"
	"log"
	"objects"
	"shopify"
	"time"

	"github.com/google/uuid"
)

// Pushes a product to Shopify
func (dbconfig *DbConfig) PushProduct(configShopify *shopify.ConfigShopify, product objects.Product) {
	product_id, err := GetProductID(dbconfig, product.ProductCode)
	if err != nil {
		// TODO log error to something
		log.Println(err)
		return
	}
	shopifyProduct := ConvertProductToShopify(product)
	if product_id != "" {
		// If yes, then update (UpdateProductShopify)
		// done
		configShopify.UpdateProductShopify(shopifyProduct, product_id)
		return
	}
	// If product does not exist
	// Create the product on website (save IDs)
	product_id, err = configShopify.AddProductShopify(shopifyProduct)
	if err != nil {
		// TODO log error to something
		log.Println(err)
		return
	}
	err = dbconfig.DB.CreatePID(context.Background(), database.CreatePIDParams{
		ID:               uuid.New(),
		ProductCode:      product.ProductCode,
		ShopifyProductID: product_id,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		return
	}
	for _, variant := range product.Variants {
		dbconfig.PushVariant(configShopify, variant)
	}
	// Add variants to website as well (save IDs)
	// Add Collection to website
	// done
}

// Pushes a variant to Shopify
func (dbconfig *DbConfig) PushVariant(configShopify *shopify.ConfigShopify, variant objects.ProductVariant) {
	variant_id, err := GetVariantID(dbconfig, variant.Sku)
	if err != nil {
		// TODO log error to something
		log.Println(err)
		return
	}
	if variant_id != "" {
		// If yes, then update (UpdateProductShopify)
		// done
		configShopify.UpdateVariantShopify(ConvertVariantToShopify(variant), variant_id)
	}

	ids, err := configShopify.GetProductBySKU(variant.Sku)
	if err != nil {
		log.Println(err)
		return
	}
	if ids.VariantID != "" && len(ids.VariantID) > 0 {
		// update existing variant.
		err = configShopify.UpdateVariantShopify(ConvertVariantToShopify(variant), ids.VariantID)
		if err != nil {
			log.Println(err)
			return
		}
		// TODO should we add Ids to the DB when updating?
	} else {
		// create new variant
		variant_id, err = configShopify.AddVariantShopify(ConvertVariantToShopify(variant), ids.ProductID)
		if err != nil {
			log.Println(err)
			return
		}
		err := dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
			ID:               uuid.New(),
			Sku:              variant.Sku,
			ShopifyVariantID: variant_id,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
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
