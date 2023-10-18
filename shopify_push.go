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

// Return the price of the product for a specific tier
func (dbconfig *DbConfig) ShopifyVariantPricing(
	variant objects.ProductVariant,
	price_tier string) string {
	price_name, err := dbconfig.GetSettingValue(price_tier)
	// get the price of the product here
	if err != nil {
		log.Println(err)
		return "0"
	}
	for _, price := range variant.VariantPricing {
		if price.Name == price_name {
			return price.Value
		}
	}
	return "0"
}

// Calculate stock to send as the available_adjustment
func (dbconfig *DbConfig) CalculateAvailableQuantity(
	configShopify *shopify.ConfigShopify,
	db_quantity int32,
	location_id,
	inventory_item_id string) int32 {
	db_inventory_level, err := dbconfig.DB.GetShopifyInventory(context.Background(), database.GetShopifyInventoryParams{
		InventoryItemID:   inventory_item_id,
		ShopifyLocationID: location_id,
	})
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println(err)
			return 0
		}
	}
	if db_inventory_level.CreatedAt.IsZero() {
		shopify_inventory_level, err := configShopify.GetShopifyInventoryLevels(location_id, inventory_item_id)
		if err != nil {
			log.Println(err)
			return 0
		}
		available := int32(db_quantity) - (int32(shopify_inventory_level.Available) - db_inventory_level.Available)
		err = dbconfig.DB.CreateShopifyInventoryRecord(context.Background(), database.CreateShopifyInventoryRecordParams{
			ID:                uuid.New(),
			ShopifyLocationID: fmt.Sprint(location_id),
			InventoryItemID:   fmt.Sprint(inventory_item_id),
			Available:         available,
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		})
		if err != nil {
			log.Println(err)
			return 0
		}
		return available
	} else {
		available := int32(db_quantity) - db_inventory_level.Available
		err = dbconfig.DB.UpdateShopifyInventoryRecord(context.Background(), database.UpdateShopifyInventoryRecordParams{
			Available:         available,
			UpdatedAt:         time.Now().UTC(),
			ShopifyLocationID: fmt.Sprint(location_id),
			InventoryItemID:   fmt.Sprint(inventory_item_id),
		})
		if err != nil {
			log.Println(err)
			return 0
		}
		return available
	}
}

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
func (dbconfig *DbConfig) PushProductInventory(configShopify *shopify.ConfigShopify, variant objects.ProductVariant) {
	shopify_inventory, err := dbconfig.DB.GetInventoryIDBySKU(context.Background(), variant.Sku)
	if err != nil {
		log.Println("1 | " + err.Error())
		return
	}
	int_inventory_id, err := strconv.Atoi(shopify_inventory.ShopifyInventoryID)
	if err != nil {
		log.Println("2 | " + err.Error())
		return
	}
	if int_inventory_id == 0 {
		log.Println(errors.New("invalid variant inventory id"))
		return
	}
	for _, variant_qty := range variant.VariantQuantity {
		data, err := dbconfig.DB.GetShopifyLocationByWarehouse(context.Background(), variant_qty.Name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				log.Println(errors.New("invalid location_id, please reconfigure map"))
				return
			} else {
				log.Println("3 | " + err.Error())
				return
			}
		}
		int_location_id, err := strconv.Atoi(data.ShopifyLocationID)
		if err != nil {
			log.Println("4 | " + err.Error())
			return
		}
		if int_location_id == 0 {
			log.Println(errors.New("invalid location_id, please reconfigure map"))
			return
		}
		link, err := dbconfig.DB.GetInventoryLocationLink(context.Background(), database.GetInventoryLocationLinkParams{
			InventoryItemID: shopify_inventory.ShopifyInventoryID,
			WarehouseName:   variant_qty.Name,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				log.Println("5 | " + err.Error())
				return
			}
		}
		if link.ShopifyLocationID == "" || len(link.ShopifyLocationID) == 0 {
			linked_location_id, err := strconv.Atoi(data.ShopifyLocationID)
			if err != nil {
				log.Println("6 | " + err.Error())
				return
			}
			_, err = configShopify.AddInventoryItemToLocation(linked_location_id, int_inventory_id)
			if err != nil {
				log.Println("7 | " + err.Error())
				return
			}
			available := dbconfig.CalculateAvailableQuantity(
				configShopify,
				int32(variant_qty.Value),
				fmt.Sprint(linked_location_id),
				fmt.Sprint(int_inventory_id),
			)
			_, err = configShopify.AddLocationQtyShopify(int_location_id, int_inventory_id, int(available))
			if err != nil {
				log.Println("8 | " + err.Error())
				return
			}
		} else {
			available := dbconfig.CalculateAvailableQuantity(
				configShopify,
				int32(variant_qty.Value),
				fmt.Sprint(int_location_id),
				fmt.Sprint(int_inventory_id),
			)
			_, err = configShopify.AddLocationQtyShopify(int_location_id, int_inventory_id, int(available))
			if err != nil {
				log.Println("10 | " + err.Error())
				return
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
	if product_id != "" && len(product_id) > 0 {
		configShopify.UpdateProductShopify(shopifyProduct, product_id)
		return
	}
	// TODO should it look until it finds a variant?
	// Make this a setting
	ids, err := configShopify.GetProductBySKU(product.Variants[0].Sku)
	if err != nil {
		log.Println("1 -" + err.Error())
		return
	}
	if ids.ProductID != "" && len(ids.ProductID) > 0 {
		// update existing product on the website
		product_data, err := configShopify.UpdateProductShopify(shopifyProduct, ids.ProductID)
		if err != nil {
			log.Println("2 -" + err.Error())
			return
		}
		err = dbconfig.DB.CreatePID(context.Background(), database.CreatePIDParams{
			ID:               uuid.New(),
			ProductCode:      product.ProductCode,
			ProductID:        product.ID,
			ShopifyProductID: fmt.Sprint(product_data.Product.ID),
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
		})
		if err != nil {
			log.Println("3 -" + err.Error())
			return
		}
		for key := range product.Variants {
			dbconfig.PushVariant(
				configShopify,
				product.Variants[key],
				fmt.Sprint(product_data.Product.ID),
				fmt.Sprint(product_data.Product.Variants[key].ID))
		}
	} else {
		// add new product to website
		product_data, err := configShopify.AddProductShopify(shopifyProduct)
		if err != nil {
			// TODO log error to something
			log.Println(err)
			return
		}
		if product_data.Product.ID != 0 {
			err = dbconfig.DB.CreatePID(context.Background(), database.CreatePIDParams{
				ID:               uuid.New(),
				ProductCode:      product.ProductCode,
				ProductID:        product.ID,
				ShopifyProductID: fmt.Sprint(product_data.Product.ID),
				CreatedAt:        time.Now().UTC(),
				UpdatedAt:        time.Now().UTC(),
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
		if product_data.Product.ID != 0 {
			err = dbconfig.CollectionShopfy(configShopify, product, product_data.Product.ID)
			if err != nil {
				log.Println(err)
				return
			}
		}
		for key := range product.Variants {
			dbconfig.PushVariant(
				configShopify,
				product.Variants[key],
				fmt.Sprint(product_data.Product.ID),
				fmt.Sprint(product_data.Product.Variants[key].ID))
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
	product_id string,
	product_variant_id string) {
	shopifyVariant := ConvertVariantToShopify(variant)

	// Sets the price for the variant here
	shopifyVariant.Price = dbconfig.ShopifyVariantPricing(variant, "default_price_tier")
	shopifyVariant.CompareAtPrice = dbconfig.ShopifyVariantPricing(variant, "default_compare_at_price")
	if product_variant_id != "" && len(product_variant_id) > 0 {
		variant_data, err := configShopify.UpdateVariantShopify(shopifyVariant, product_variant_id)
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
		dbconfig.PushProductInventory(configShopify, variant)
		return
	}
	variant_id, err := GetVariantID(dbconfig, variant.Sku)
	if err != nil {
		// TODO log error to something
		log.Println(err)
		return
	}
	if variant_id != "" && len(variant_id) > 0 {
		configShopify.UpdateVariantShopify(shopifyVariant, variant_id)
		dbconfig.PushProductInventory(configShopify, variant)
		return
	}
	ids, err := configShopify.GetProductBySKU(variant.Sku)
	if err != nil {
		log.Println(err)
		return
	}
	if ids.VariantID != "" && len(ids.VariantID) > 0 {
		variant_data, err := configShopify.UpdateVariantShopify(shopifyVariant, ids.VariantID)
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
		dbconfig.PushProductInventory(configShopify, variant)
	} else {
		variant_data, err := configShopify.AddVariantShopify(shopifyVariant, product_id)
		if err != nil {
			log.Println("variant_error: " + err.Error())
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
// make this a go routine
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
