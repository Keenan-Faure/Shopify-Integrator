package main

import (
	"context"
	"integrator/internal/database"
	"log"
	"objects"
	"shopify"
	"strconv"
	"time"
	"utils"

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
		configShopify.UpdateProductShopify(shopifyProduct, product_id)
		return
	}
	product_response, err := configShopify.AddProductShopify(shopifyProduct)
	if err != nil {
		// TODO log error to something
		log.Println(err)
		return
	}
	ids := ConvertToShopifyIDs(product_response)
	if ids.ProductID != "" && len(ids.ProductID) > 0 {
		err = dbconfig.DB.CreatePID(context.Background(), database.CreatePIDParams{
			ID:               uuid.New(),
			ProductCode:      product.ProductCode,
			ProductID:        product.ID,
			ShopifyProductID: ids.ProductID,
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
	for key := range ids.Variants {
		if ids.Variants[key].VariantID != "" && len(ids.Variants[key].VariantID) > 0 {
			err = dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
				ID:               uuid.New(),
				Sku:              product.Variants[key].Sku,
				ShopifyVariantID: ids.Variants[key].VariantID,
				VariantID:        product.Variants[key].ID,
				CreatedAt:        time.Now().UTC(),
				UpdatedAt:        time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	str_int, err := strconv.Atoi(ids.ProductID)
	if err != nil {
		log.Println(err)
		return
	}
	dbconfig.CollectionShopfy(configShopify, product, str_int)
}

func (dbconfig *DbConfig) CollectionShopfy(
	configShopify *shopify.ConfigShopify,
	product objects.Product,
	shopify_product_id int) error {
	db_category, err := dbconfig.DB.GetShopifyCollection(context.Background(), utils.ConvertStringToSQL(product.Category))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	if db_category.ShopifyCollectionID != 0 || db_category.ProductCollection.String != "" {
		shopify_categories, err := configShopify.GetShopifyCategories()
		if err != nil {
			return err
		}
		exists, collection_id := configShopify.CategoryExists(product, shopify_categories)
		if exists {
			err = dbconfig.DB.CreateShopifyCollection(context.Background(), database.CreateShopifyCollectionParams{
				ID:                  uuid.New(),
				ProductCollection:   utils.ConvertStringToSQL(product.Category),
				ShopifyCollectionID: int32(collection_id),
				CreatedAt:           time.Now().UTC(),
				UpdatedAt:           time.Now().UTC(),
			})
			if err != nil {
				return err
			}
			// TODO might need the response here?
			_, err = configShopify.AddProductToCollectionShopify(shopify_product_id, collection_id)
			if err != nil {
				return err
			}
			return nil
		}
		custom_collection_id, err := configShopify.AddCustomCollectionShopify(product.Category)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		err = dbconfig.DB.CreateShopifyCollection(context.Background(), database.CreateShopifyCollectionParams{
			ID:                  uuid.New(),
			ProductCollection:   utils.ConvertStringToSQL(product.Category),
			ShopifyCollectionID: int32(custom_collection_id),
			CreatedAt:           time.Now().UTC(),
			UpdatedAt:           time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Pushes a variant to Shopify
func (dbconfig *DbConfig) PushVariant(
	configShopify *shopify.ConfigShopify,
	variant objects.ProductVariant,
	product_id string) {
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
	} else {
		// create new variant
		variant_id, err = configShopify.AddVariantShopify(ConvertVariantToShopify(variant), product_id)
		if err != nil {
			log.Println(err)
			return
		}
		err := dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
			ID:               uuid.New(),
			Sku:              variant.Sku,
			ShopifyVariantID: variant_id,
			VariantID:        variant.ID,
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
		if err.Error() == "sql: no rows in result set" {
			return "", nil
		}
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
		if err.Error() == "sql: no rows in result set" {
			return "", nil
		}
		return "", err
	}
	if variant_id.Sku == sku &&
		len(variant_id.ShopifyVariantID) > 0 &&
		variant_id.ShopifyVariantID != "" {
		return variant_id.ShopifyVariantID, nil
	}
	return "", nil
}
