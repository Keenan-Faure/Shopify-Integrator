package main

import (
	"context"
	"errors"
	"fmt"
	"integrator/internal/database"
	"log"
	"objects"
	"shopify"
	"strconv"
	"time"
	"utils"

	"github.com/google/uuid"
)

// TODO create a feed that fetches all locations from shopify
// pops up a list of locations and it asks the user which will be used, and which warehouse should be
// mapped to the respective location
func (dbconfig *DbConfig) FetchShopifyLocations(configShopify *shopify.ConfigShopify) {
	// call GetLocationsShopify() functiion to retrieve all of them
	// respond with that on the API
	// let javascript create the respective element after fetching
	// once the user has been taken through the form post it to an endpoint TBC which one
	// done
	response, err := configShopify.GetLocationsShopify()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(response)
}

// Removes mapping between Location and warehouses
func (dbconfig *DbConfig) RemoveLocationMap(id string) error {
	err := IDValidation(id)
	if err != nil {
		return err
	}
	delete_id, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = dbconfig.DB.RemoveShopifyLocationMap(context.Background(), delete_id)
	if err != nil {
		return err
	}
	return err
}

// Pushes an Inventory update to Shopify for a specific SKU
func (dbconfig *DbConfig) PushProductInventory(configShopify *shopify.ConfigShopify, product objects.Product) {
	for _, variant := range product.Variants {
		shopify_inventory, err := dbconfig.DB.GetInventoryIDBySKU(context.Background(), variant.Sku)
		if err != nil {
			log.Println(err)
			return
		}
		int_inventory_id, err := strconv.Atoi(shopify_inventory.ShopifyInventoryID)
		if err != nil {
			log.Println(err)
			return
		}
		for _, variant_qty := range variant.VariantQuantity {
			// checks if the location -> warehouse map has been completed
			data, err := dbconfig.DB.GetShopifyLocationByWarehouse(context.Background(), variant_qty.Name)
			if err != nil {
				log.Println(err)
				return
			}
			int_location_id, err := strconv.Atoi(data.ShopifyLocationID)
			if err != nil {
				log.Println(err)
				return
			}
			// if invalid map
			if int_location_id == 0 {
				// reconfigure the location_id map inside settings
				log.Println(errors.New("invalid location_id, please reconfigure map"))
				return
			}
			// valid map
			// checks if a variant (inventory_item_id) is linked to a location already
			link, err := dbconfig.DB.GetInventoryLocationLink(context.Background(), database.GetInventoryLocationLinkParams{
				InventoryItemID: shopify_inventory.ShopifyInventoryID,
				WarehouseName:   variant_qty.Name,
			})
			if err != nil {
				log.Println(err)
				return
			}
			// check if an item is linked to a Location already
			// if it's not linked
			if link.ShopifyLocationID == "" || len(link.ShopifyLocationID) == 0 {
				// item is not linked to warehouse
				// link it
				linked_location_id, err := strconv.Atoi(data.ShopifyLocationID)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = configShopify.AddInventoryItemToLocation(linked_location_id, int_inventory_id)
				if err != nil {
					log.Println(err)
					return
				}
				err = dbconfig.DB.CreateShopifyLocation(context.Background(), database.CreateShopifyLocationParams{
					ID:                   uuid.New(),
					ShopifyWarehouseName: data.ShopifyWarehouseName,
					ShopifyLocationID:    data.ShopifyLocationID,
					WarehouseName:        variant_qty.Name,
					CreatedAt:            time.Now().UTC(),
					UpdatedAt:            time.Now().UTC(),
				})
				if err != nil {
					log.Println(err)
					return
				}
				_, err = configShopify.AddLocationQtyShopify(int_location_id, int_inventory_id, variant_qty.Value)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				_, err = configShopify.AddLocationQtyShopify(int_location_id, int_inventory_id, variant_qty.Value)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
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
