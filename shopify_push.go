package main

import (
	"context"
	"fmt"
	"integrator/internal/database"
	"log"
	"objects"
	"shopify"
	"time"
	"utils"

	"github.com/google/uuid"
)

// TODO create a feed that fetches all locations from shopify
// pops up a list of locations and it asks the user which will be used, and which warehouse should be
// mapped to the respective location
func (dbconfig *DbConfig) FetchShopifyLocations(configShopify *shopify.ConfigShopify, product objects.Product) {
	// call GetLocationsShopify() functiion to retrieve all of them
	// respond with that on the API
	// let javascript create the respective element after fetching
	// once the user has been taken through the form post it to an endpoint TBC which one
	// done
}

// Pushes an Inventory update to Shopify for a specific SKU
func (dbconfig *DbConfig) PushProductInventory(configShopify *shopify.ConfigShopify, product objects.Product) {
	// field that links warehouse to location id should be added to the variant_qty table
	for _, variant := range product.Variants {
		for _, variant_qty := range variant.VariantQuantity {
			data, err := dbconfig.DB.GetShopifyLocationByWarehouse(context.Background(), variant_qty.Name)
			if err != nil {
				log.Println(err)
				return
			}
			if data.ShopifyLocationID != "" && len(data.ShopifyLocationID) > 0 {
				// create the link between the item and the respective location

			}
		}
	}

	// loop through all variants of the product
	// get the location_id for the respective variant from the database variants_qty table
	// dbconfig.DB.GetShopifyLocationByWarehouse()
	// if the location_id exists and is non-zero
	// create the link between the item and the respective location
	// this will use the
	// if its a valid location_id, then proceed with adjusting the qty of the item
	// otherwise if its invalid then error message should pop up (404)
	// advise user to reconfigure the location_id map

}

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
	if product_response.Product.ID != 0 {
		err = dbconfig.DB.CreatePID(context.Background(), database.CreatePIDParams{
			ID:               uuid.New(),
			ProductCode:      product.ProductCode,
			ProductID:        product.ID,
			ShopifyProductID: fmt.Sprint(product_response.Product.ID),
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
	for key := range product_response.Product.Variants {
		if product_response.Product.Variants[key].ID != 0 {
			err = dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
				ID:                 uuid.New(),
				Sku:                product.Variants[key].Sku,
				ShopifyVariantID:   fmt.Sprint(product_response.Product.Variants[key].ID),
				ShopifyInventoryID: fmt.Sprint(product_response.Product.Variants[key].InventoryItemID),
				VariantID:          product.Variants[key].ID,
				CreatedAt:          time.Now().UTC(),
				UpdatedAt:          time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	if product_response.Product.ID != 0 {
		err = dbconfig.CollectionShopfy(configShopify, product, product_response.Product.ID)
		if err != nil {
			log.Println(err)
			return
		}
	}
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
	if db_category.ShopifyCollectionID == 0 || db_category.ProductCollection.String == "" {
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
		variant_data, err := configShopify.AddVariantShopify(ConvertVariantToShopify(variant), product_id)
		if err != nil {
			log.Println(err)
			return
		}
		err = dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
			ID:                 uuid.New(),
			Sku:                variant.Sku,
			ShopifyVariantID:   fmt.Sprint(variant_data.Variant.ID),
			ShopifyInventoryID: fmt.Sprint(variant_data.Variant.InventoryItemID),
			VariantID:          variant.ID,
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
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
