package main

import (
	"context"
	"errors"
	"fmt"
	"integrator/internal/database"
	"log"
	"net/http"
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
	price_tier string) (string, error) {
	price, err := dbconfig.DB.GetPriceTierBySKU(context.Background(), database.GetPriceTierBySKUParams{
		Sku:  variant.Sku,
		Name: price_tier,
	})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "0.00", nil
		}
		return "0.00", err
	}
	return price.Value.String, nil
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
		shopify_inventory_level, err := configShopify.GetShopifyInventoryLevel(location_id, inventory_item_id)
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
func (dbconfig *DbConfig) PushProductInventory(configShopify *shopify.ConfigShopify, variant objects.ProductVariant) error {
	shopify_inventory, err := dbconfig.DB.GetInventoryIDBySKU(context.Background(), variant.Sku)
	if err != nil {
		return err
	}
	if shopify_inventory.ShopifyInventoryID == "" || len(shopify_inventory.ShopifyInventoryID) == 0 {
		return errors.New("invalid inventory_item_id for sku: " + variant.Sku)
	}
	int_inventory_id, err := strconv.Atoi(shopify_inventory.ShopifyInventoryID)
	if err != nil {
		return err
	}
	if int_inventory_id == 0 {
		return err
	}
	for _, variant_qty := range variant.VariantQuantity {
		data, err := dbconfig.DB.GetShopifyLocationByWarehouse(context.Background(), variant_qty.Name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return errors.New("invalid location_id, please reconfigure map")
			} else {
				return err
			}
		}
		int_location_id, err := strconv.Atoi(data.ShopifyLocationID)
		if err != nil {
			return err
		}
		if int_location_id == 0 {
			return errors.New("invalid location_id, please reconfigure map")
		}
		link, err := dbconfig.DB.GetInventoryLocationLink(context.Background(), database.GetInventoryLocationLinkParams{
			InventoryItemID: shopify_inventory.ShopifyInventoryID,
			WarehouseName:   variant_qty.Name,
		})
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return err
			}
		}
		if link.ShopifyLocationID == "" || len(link.ShopifyLocationID) == 0 {
			linked_location_id, err := strconv.Atoi(data.ShopifyLocationID)
			if err != nil {
				return err
			}
			_, err = configShopify.AddInventoryItemToLocation(linked_location_id, int_inventory_id)
			if err != nil {
				return err
			}
			available := dbconfig.CalculateAvailableQuantity(
				configShopify,
				int32(variant_qty.Value),
				fmt.Sprint(linked_location_id),
				fmt.Sprint(int_inventory_id),
			)
			_, err = configShopify.AddLocationQtyShopify(int_location_id, int_inventory_id, int(available))
			if err != nil {
				return err
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
				return err
			}
		}
	}
	return nil
}

// Pushes a product to Shopify
func (dbconfig *DbConfig) PushProduct(configShopify *shopify.ConfigShopify, product objects.Product) error {
	product_id, err := GetProductID(dbconfig, product.ProductCode)
	if err != nil {
		return err
	}
	restrictions, err := dbconfig.DB.GetPushRestriction(context.Background())
	if err != nil {
		return err
	}
	shopifyProduct := ConvertProductToShopify(product)
	push_restrictions := PushRestrictionsToMap(restrictions)
	update_shopify_product := ApplyPushRestrictionProduct(push_restrictions, shopifyProduct)
	if product_id != "" && len(product_id) > 0 {
		_, err := configShopify.UpdateProductShopify(update_shopify_product, product_id)
		return err
	}
	dynamic_search_enabled := false
	dynamic_search, err := dbconfig.DB.GetShopifySettingByKey(
		context.Background(),
		"shopify_enable_dynamic_sku_search",
	)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
		dynamic_search.Value = "false"
	}
	dynamic_search_enabled, err = strconv.ParseBool(dynamic_search.Value)
	if err != nil {
		return err
	}
	if !dynamic_search_enabled {
		ids, err := configShopify.GetProductBySKU(product.Variants[0].Sku)
		if err != nil {
			return err
		}
		err = PushAddShopify(configShopify, dbconfig, ids, product, shopifyProduct, update_shopify_product)
		if err != nil {
			return err
		}
		return nil
	} else {
		ids := objects.ResponseIDs{}
		// searches if the product variants exists on shpify already
		for _, variant := range product.Variants {
			ids_search, err := configShopify.GetProductBySKU(variant.Sku)
			if err != nil {
				return err
			}
			if ids_search.ProductID != "" || len(ids_search.ProductID) != 0 {
				ids = ids_search
				break
			}
		}
		err = PushAddShopify(configShopify, dbconfig, ids, product, shopifyProduct, update_shopify_product)
		if err != nil {
			return err
		}

		return nil
	}
}

func PushAddShopify(
	configShopify *shopify.ConfigShopify,
	dbconfig *DbConfig,
	ids objects.ResponseIDs,
	product objects.Product,
	shopifyProduct objects.ShopifyProduct,
	update_shopify_product objects.ShopifyProduct,
) error {
	if ids.ProductID != "" && len(ids.ProductID) > 0 {
		// update existing product on the website
		product_data, err := configShopify.UpdateProductShopify(update_shopify_product, ids.ProductID)
		if err != nil {
			return err
		}
		// add local category to the product (first check if it's set to be sent to Shopify)
		// in the restrictions
		restrictions, err := dbconfig.DB.GetPushRestriction(context.Background())
		if err != nil {
			return err
		}
		if DeterPushRestriction(PushRestrictionsToMap(restrictions), "category") {
			if product.Category != "" || len(product.Category) > 0 {
				err = dbconfig.CollectionShopfy(configShopify, product, int(product_data.Product.ID))
				if err != nil {
					return err
				}
			}
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
			if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
				err = dbconfig.DB.UpdatePID(context.Background(), database.UpdatePIDParams{
					ShopifyProductID: fmt.Sprint(product_data.Product.ID),
					UpdatedAt:        time.Now().UTC(),
					ProductCode:      product.ProductCode,
				})
			}
			return err
		}
		return nil
	} else {
		// add new product to website
		product_data, err := configShopify.AddProductShopify(shopifyProduct)
		if err != nil {
			return err
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
				return err
			}
		}
		if product_data.Product.ID != 0 {
			err = dbconfig.CollectionShopfy(configShopify, product, int(product_data.Product.ID))
			if err != nil {
				return err
			}
		}
		// add all variant skus to the database
		for _, variant := range product.Variants {
			err = SaveVariantIds(dbconfig, product_data, variant.Sku)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (dbconfig *DbConfig) CollectionShopfy(
	configShopify *shopify.ConfigShopify,
	product objects.Product,
	shopify_product_id int) error {
	// check if the product has the associated category linked to it already...
	categories, err := configShopify.GetShopifyCategoryByProductID(fmt.Sprint(shopify_product_id))
	if err != nil {
		return err
	}
	if len(categories.CustomCollections) > 0 {
		if product.Category == categories.CustomCollections[0].Title {
			// the product already has the category added to it
			return nil
		}
	}
	db_category, err := dbconfig.DB.GetShopifyCollection(context.Background(), utils.ConvertStringToSQL(product.Category))
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	// if the database does not contain any ids for the category
	// then it means that it is not linked to shopify
	if db_category.ShopifyCollectionID == "" || db_category.ProductCollection.String == "" {
		shopify_categories, err := configShopify.GetShopifyCategories()
		if err != nil {
			return err
		}
		exists, collection_id := configShopify.CategoryExists(product, shopify_categories)
		if exists {
			err = dbconfig.DB.CreateShopifyCollection(context.Background(), database.CreateShopifyCollectionParams{
				ID:                  uuid.New(),
				ProductCollection:   utils.ConvertStringToSQL(product.Category),
				ShopifyCollectionID: fmt.Sprint(collection_id),
				CreatedAt:           time.Now().UTC(),
				UpdatedAt:           time.Now().UTC(),
			})
			if err != nil {
				return err
			}
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
		err = dbconfig.DB.CreateShopifyCollection(context.Background(), database.CreateShopifyCollectionParams{
			ID:                  uuid.New(),
			ProductCollection:   utils.ConvertStringToSQL(product.Category),
			ShopifyCollectionID: fmt.Sprint(custom_collection_id),
			CreatedAt:           time.Now().UTC(),
			UpdatedAt:           time.Now().UTC(),
		})
		if err != nil {
			return err
		}
		_, err = configShopify.AddProductToCollectionShopify(shopify_product_id, custom_collection_id)
		if err != nil {
			return err
		}
	}

	// if the database does contain the links then we should update it on shopify
	int_shopify_collection, err := strconv.Atoi(db_category.ShopifyCollectionID)
	if err != nil {
		return err
	}
	_, err = configShopify.AddProductToCollectionShopify(shopify_product_id, int_shopify_collection)
	if err != nil {
		return err
	}
	return nil
}

// Pushes a variant to Shopify
func (dbconfig *DbConfig) PushVariant(
	configShopify *shopify.ConfigShopify,
	variant objects.ProductVariant,
	product_variant objects.ShopifyVariant,
	restrictions map[string]string,
	shopify_product_id string,
	shopify_variant_id string,
) error {
	product_variant_adding := ConvertVariantToShopify(variant)
	price, err := dbconfig.ShopifyVariantPricing(variant, "Selling Price")
	if err != nil {
		return err
	}
	product_variant_adding.Price = price
	compare_to_price, err := dbconfig.ShopifyVariantPricing(variant, "Compare At Price")
	if err != nil {
		return err
	}
	product_variant_adding.CompareAtPrice = compare_to_price
	// only update price if the restriction says so
	if DeterPushRestriction(restrictions, "pricing") {
		product_variant.CompareAtPrice = compare_to_price
		product_variant.Price = price
	}
	if shopify_variant_id == "" && len(shopify_variant_id) == 0 {
		db_shopify_variant_id, err := dbconfig.DB.GetVIDBySKU(context.Background(), variant.Sku)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return err
			}
		} else {
			shopify_variant_id = db_shopify_variant_id.ShopifyVariantID
		}
	}
	if shopify_variant_id == "" && len(shopify_variant_id) == 0 {
		// determines if the SKU is found on the website
		// if found we use that product
		ids, err := configShopify.GetProductBySKU(variant.Sku)
		if err != nil {
			return err
		}
		shopify_variant_id = ids.VariantID
	}
	if shopify_variant_id == "" && len(shopify_variant_id) == 0 {
		// if the shopify_variant_id is empty
		// we create a new variant on the website for this product
		if (shopify_product_id == "") || len(shopify_product_id) == 0 {
			db_shopify_product_id, err := dbconfig.DB.GetPIDBySKU(context.Background(), variant.Sku)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return err
				}
				db_shopify_product_id = ""
			}
			shopify_product_id = db_shopify_product_id
			if (shopify_product_id == "") || len(shopify_product_id) == 0 {
				return errors.New("unable to process product with SKU " + variant.Sku + " to shopify")
			}
		}
		variants, err := configShopify.AddVariantShopify(product_variant, shopify_product_id)
		if err != nil {
			return err
		}
		shopify_variant_id = fmt.Sprint(variants.Variant.ID)
		// after adding the new variant, we should update it
		// hence we do not return
	}
	// update variant
	variant_data, err := configShopify.UpdateVariantShopify(product_variant, shopify_variant_id)
	if err != nil {
		return err
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
		if err.Error()[0:50] == "pq: duplicate key value violates unique constraint" {
			err = dbconfig.DB.UpdateVID(context.Background(), database.UpdateVIDParams{
				ShopifyVariantID:   fmt.Sprint(variant_data.Variant.ID),
				ShopifyInventoryID: fmt.Sprint(variant_data.Variant.InventoryItemID),
				UpdatedAt:          time.Now().UTC(),
				Sku:                variant.Sku,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// determine if warehousing should be updated
	if DeterPushRestriction(restrictions, "warehousing") {
		err = dbconfig.PushProductInventory(configShopify, variant)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompileInstructionProduct(dbconfig *DbConfig, product objects.Product, dbUser database.User) error {
	queue_item := objects.RequestQueueHelper{
		Type:        "product",
		Status:      "in-queue",
		Instruction: "add_product",
		Endpoint:    "queue",
		ApiKey:      dbUser.ApiKey,
		Method:      http.MethodPost,
		Object:      nil,
	}
	product_id, err := dbconfig.DB.GetPIDByProductCode(context.Background(), product.ProductCode)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	queue_item_object := objects.RequestQueueItemProducts{
		SystemProductID: product.ID.String(),
		SystemVariantID: "",
		Shopify: struct {
			ProductID string "json:\"product_id\""
			VariantID string "json:\"variant_id\""
		}{
			ProductID: product_id.ShopifyProductID,
			VariantID: "",
		},
	}
	queue_item.Object = queue_item_object
	if product_id.ShopifyProductID == "" {
		// add product
		_, err := dbconfig.QueueHelper(queue_item)
		if err != nil {
			return err
		}
	} else {
		// update product
		queue_item.Instruction = "update_product"
		_, err := dbconfig.QueueHelper(queue_item)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompileInstructionVariant(dbconfig *DbConfig, variant objects.ProductVariant, product objects.Product, dbUser database.User) error {
	queue_item := objects.RequestQueueHelper{
		Type:        "product_variant",
		Status:      "in-queue",
		Instruction: "add_variant",
		Endpoint:    "queue",
		ApiKey:      dbUser.ApiKey,
		Method:      http.MethodPost,
		Object:      nil,
	}
	variant_id, err := dbconfig.DB.GetVIDBySKU(context.Background(), variant.Sku)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	product_id, err := dbconfig.DB.GetPIDByProductCode(context.Background(), product.ProductCode)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return err
		}
	}
	queue_item_object := objects.RequestQueueItemProducts{
		SystemProductID: product.ID.String(),
		SystemVariantID: variant.ID.String(),
		Shopify: struct {
			ProductID string "json:\"product_id\""
			VariantID string "json:\"variant_id\""
		}{
			ProductID: product_id.ShopifyProductID,
			VariantID: variant_id.ShopifyVariantID,
		},
	}
	queue_item.Object = queue_item_object
	if variant_id.ShopifyVariantID == "" {
		// since its blank we should do an add instruction
		_, err := dbconfig.QueueHelper(queue_item)
		if err != nil {
			return err
		}
	} else {
		// update instruction
		// get the object body
		queue_item.Instruction = "update_variant"
		_, err := dbconfig.QueueHelper(queue_item)
		if err != nil {
			return err
		}
	}
	return nil
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

func SaveVariantIds(dbconfig *DbConfig, shopify_product objects.ShopifyProductResponse, sku string) error {
	for _, variant := range shopify_product.Product.Variants {
		if variant.Sku == sku {
			// get VariantID by SKU
			variant_id, err := dbconfig.DB.GetVariantBySKU(context.Background(), sku)
			if err != nil {
				return err
			}
			// save variant_ids locally
			err = dbconfig.DB.CreateVID(context.Background(), database.CreateVIDParams{
				ID:                 uuid.New(),
				Sku:                sku,
				ShopifyVariantID:   fmt.Sprint(variant.ID),
				ShopifyInventoryID: fmt.Sprint(variant.InventoryItemID),
				VariantID:          variant_id.ID,
				CreatedAt:          time.Now().UTC(),
				UpdatedAt:          time.Now().UTC(),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
