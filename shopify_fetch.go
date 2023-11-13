package main

import (
	"context"
	"fmt"
	"integrator/internal/database"
	"log"
	"shopify"
	"strconv"
	"time"
	"utils"

	"github.com/google/uuid"
)

// loop function that uses Goroutine to run
// the fetch function each interval
func LoopJSONShopify(
	dbconfig *DbConfig,
	shopifyConfig shopify.ConfigShopify) {
	fetch_url := ""
	fetch_time := 5
	fetch_time_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_shopify_fetch_time")
	if err != nil {
		fetch_time = 5
	}
	fetch_time, err = strconv.Atoi(fetch_time_db.Value)
	if err != nil {
		fetch_time = 5
	}
	// do not allow fetch time lower than 5 minutes
	if fetch_time < 5 {
		fetch_time = 5
	}
	ticker := time.NewTicker(time.Duration(fetch_time) * time.Minute)
	for ; ; <-ticker.C {
		fmt.Println("running fetch worker...")
		fetch_enabled := false
		fetch_enabled_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_enable_shopify_fetch")
		if err != nil {
			fetch_enabled = false
		}
		fetch_enabled, err = strconv.ParseBool(fetch_enabled_db.Value)
		if err != nil {
			fetch_enabled = false
		}
		// if the `app_enable_shopify_fetch` setting is enabled
		// then we attempt to fetch the products from Shopify to store locally
		if fetch_enabled {
			shopifyProds, next, err := shopifyConfig.FetchProducts(fetch_url)
			if err != nil {
				log.Println("Shopify > Error fetching next products to process:", err)
				continue
			}
			created_db_product := database.Product{}
			for _, product := range shopifyProds.Products {
				for variant_key, product_variant := range product.Variants {
					internal_product, err := dbconfig.DB.GetVariantBySKU(context.Background(), product_variant.Sku)
					if err != nil {
						if err.Error() != "sql: no rows in result set" {
							log.Println(err)
							break
						}
					}
					// if product is the same as the internal variant
					// then we will UPDATE the product
					// if the `app_fetch_overwrite_products` setting is enabled
					if internal_product.Sku == product_variant.Sku {
						overwrite := false
						overwrite_db, err := dbconfig.DB.GetAppSettingByKey(
							context.Background(),
							"app_fetch_overwrite_products",
						)
						if err != nil {
							if err.Error() != "sql: no rows in result set" {
								log.Println(err)
								break
							}
							overwrite = false
						}
						overwrite, err = strconv.ParseBool(overwrite_db.Value)
						if err != nil {
							overwrite = false
						}
						if overwrite {
							// overwrite the product data once and then the variant data
							if variant_key == 0 {
								// updates product
								// retrieves the existant product category from Shopify
								category := ""
								categories, err := shopifyConfig.GetShopifyCategoryByProductID(fmt.Sprint(product.ID))
								if err != nil {
									log.Println(err)
									break
								}
								if len(categories.CustomCollections) > 0 {
									category = categories.CustomCollections[0].Title
								}
								err = dbconfig.DB.UpdateProductBySKU(context.Background(), database.UpdateProductBySKUParams{
									Active:      "1",
									Title:       utils.ConvertStringToSQL(product.Title),
									BodyHtml:    utils.ConvertStringToSQL(product.BodyHTML),
									Category:    utils.ConvertStringToSQL(category),
									Vendor:      utils.ConvertStringToSQL(product.Vendor),
									ProductType: utils.ConvertStringToSQL(product.ProductType),
									UpdatedAt:   time.Now().UTC(),
									Sku:         product_variant.Sku,
								})
								if err != nil {
									log.Println(err)
									break
								}
								// update product options
								if product.Options[0].Name != "Title" {
									for _, option_value := range product.Options {
										_, err = dbconfig.DB.UpdateProductOption(
											context.Background(),
											database.UpdateProductOptionParams{
												Name:       option_value.Name,
												Position:   int32(option_value.Position),
												ProductID:  internal_product.ProductID,
												Position_2: int32(option_value.Position),
											},
										)
										if err != nil {
											log.Println(err)
											break
										}
									}
								}
							}
							// update variant
							err = dbconfig.DB.UpdateVariant(context.Background(), database.UpdateVariantParams{
								Option1:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option1)),
								Option2:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option2)),
								Option3:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option3)),
								Barcode:   utils.ConvertStringToSQL(product_variant.Barcode),
								UpdatedAt: time.Now().UTC(),
								Sku:       internal_product.Sku,
							})
							if err != nil {
								log.Println(err)
								break
							}
							// update variant pricing
							create_price_tier_enabled := false
							create_price_tier_enabled_db, err := dbconfig.DB.GetAppSettingByKey(
								context.Background(),
								"app_fetch_create_price_tier_enabled",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
								create_price_tier_enabled = false
							}
							create_price_tier_enabled, err = strconv.ParseBool(create_price_tier_enabled_db.Value)
							if err != nil {
								create_price_tier_enabled = false
							}
							pricing_name, err := dbconfig.DB.GetShopifySettingByKey(
								context.Background(),
								"shopify_default_price_tier",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
							}
							// update only the price that is syncing to Shopify
							if pricing_name.Value != "" {
								err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, pricing_name.Value, product_variant.Price)
								if err != nil {
									log.Println(err)
									break
								}
							} else {
								if create_price_tier_enabled {
									// price tier is not set
									// use the default value of `fetch_price`
									err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, "fetch_price", product_variant.Price)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							pricing_compare_name, err := dbconfig.DB.GetShopifySettingByKey(
								context.Background(),
								"shopify_default_compare_at_price_tier",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
							}
							// update only the compare price that is syncing to Shopify
							if pricing_compare_name.Value != "" {
								err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, pricing_compare_name.Value, product_variant.CompareAtPrice)
								if err != nil {
									log.Println(err)
									break
								}
							} else {
								if create_price_tier_enabled {
									// price tier is not set
									// use the default value of `fetch_compare_price`
									err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, "fetch_compare_price", product_variant.CompareAtPrice)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							// check if the product's inventory should be tracked
							if product_variant.InventoryManagement == "shopify" {
								shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
									"",
									fmt.Sprint(product_variant.InventoryItemID),
								)
								if err != nil {
									log.Println(err)
									break
								}
								// create map for warehouse quantity
								total_quantity := make(map[string]int)
								if len(shopify_inventory_levels.InventoryLevels) != 0 {
									for _, inventory_level := range shopify_inventory_levels.InventoryLevels {
										warehouse, err := dbconfig.DB.GetShopifyLocationByLocationID(
											context.Background(),
											fmt.Sprint(inventory_level.LocationID),
										)
										if err != nil {
											if err.Error() == "sql: no rows in result set" {
												continue
											}
											log.Println(err)
											break
										}
										total_quantity[warehouse.WarehouseName] = total_quantity[warehouse.WarehouseName] + inventory_level.Available
									}
								}
								// only fetch the ones that exist locally
								// TODO add any missing warehouses and location maps to database
								for warehouse_name, available := range total_quantity {
									err = dbconfig.DB.UpdateVariantQty(context.Background(), database.UpdateVariantQtyParams{
										Name:      warehouse_name,
										Value:     utils.ConvertIntToSQL(available),
										Isdefault: false,
										Sku:       product_variant.Sku,
										Name_2:    warehouse_name,
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							sync_images_enabled := false
							sync_images_enabled_db, err := dbconfig.DB.GetAppSettingByKey(
								context.Background(),
								"app_fetch_sync_images",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
								sync_images_enabled = false
							}
							sync_images_enabled, err = strconv.ParseBool(sync_images_enabled_db.Value)
							if err != nil {
								sync_images_enabled = false
							}
							// update local images
							if sync_images_enabled {
								for _, image := range product.Images {
									err = AddImagery(dbconfig, internal_product.ProductID, image.Src, image.Position)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
						}
					} else {
						// product and variants can be created if setting is enabled
						create_fetched_product := false
						created_fetch_product_db, err := dbconfig.DB.GetAppSettingByKey(
							context.Background(),
							"app_fetch_add_products",
						)
						if err != nil {
							if err.Error() != "sql: no rows in result set" {
								log.Println(err)
								create_fetched_product = false
							}
						}
						create_fetched_product, err = strconv.ParseBool(created_fetch_product_db.Value)
						if err != nil {
							create_fetched_product = false
						}
						// create product only if the setting is enabled
						if create_fetched_product {
							// create product only once during first iteration
							// creates product code to be the sku of the first variant
							if variant_key == 0 {
								category := ""
								categories, err := shopifyConfig.GetShopifyCategoryByProductID(fmt.Sprint(product.ID))
								if err != nil {
									log.Println(err)
									break
								}
								if len(categories.CustomCollections) > 0 {
									category = categories.CustomCollections[0].Title
								}
								db_product, err := dbconfig.DB.CreateProduct(context.Background(), database.CreateProductParams{
									ID:          uuid.New(),
									ProductCode: product_variant.Sku,
									Active:      "1",
									Title:       utils.ConvertStringToSQL(product.Title),
									BodyHtml:    utils.ConvertStringToSQL(product.BodyHTML),
									Category:    utils.ConvertStringToSQL(category),
									Vendor:      utils.ConvertStringToSQL(product.Vendor),
									ProductType: utils.ConvertStringToSQL(product.ProductType),
									CreatedAt:   time.Now().UTC(),
									UpdatedAt:   time.Now().UTC(),
								})
								if err != nil {
									log.Println(err)
									break
								}
								created_db_product = db_product
								// create product options
								if product.Options[0].Name != "Title" {
									for _, product_options := range product.Options {
										_, err := dbconfig.DB.CreateProductOption(
											context.Background(),
											database.CreateProductOptionParams{
												ID:        uuid.New(),
												ProductID: db_product.ID,
												Name:      product_options.Name,
												Position:  int32(product_options.Position),
											},
										)
										if err != nil {
											log.Println(err)
											break
										}
									}
								}
							}
							// then create this variant and any other variant to come
							db_variant, err := dbconfig.DB.CreateVariant(
								context.Background(),
								database.CreateVariantParams{
									ID:        uuid.New(),
									ProductID: created_db_product.ID,
									Sku:       product_variant.Sku,
									Option1:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option1)),
									Option2:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option2)),
									Option3:   utils.ConvertStringToSQL(IgnoreDefaultOption(product_variant.Option3)),
									Barcode:   utils.ConvertStringToSQL(product_variant.Barcode),
									CreatedAt: time.Now().UTC(),
									UpdatedAt: time.Now().UTC(),
								},
							)
							if err != nil {
								log.Println(err)
								break
							}
							// create variant pricing
							create_price_tier_enabled := false
							create_price_tier_enabled_db, err := dbconfig.DB.GetAppSettingByKey(
								context.Background(),
								"app_fetch_create_price_tier_enabled",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
								create_price_tier_enabled = false
							}
							create_price_tier_enabled, err = strconv.ParseBool(create_price_tier_enabled_db.Value)
							if err != nil {
								create_price_tier_enabled = false
							}
							pricing_name, err := dbconfig.DB.GetShopifySettingByKey(
								context.Background(),
								"shopify_default_price_tier",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
							}
							// update only the price that is syncing to Shopify
							if pricing_name.Value != "" {
								err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, pricing_name.Value, product_variant.Price)
								if err != nil {
									log.Println(err)
									break
								}
							} else {
								if create_price_tier_enabled {
									// price tier is not set
									// use the default value of `fetch_price`
									err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, "fetch_price", product_variant.Price)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							pricing_compare_name, err := dbconfig.DB.GetShopifySettingByKey(
								context.Background(),
								"shopify_default_compare_at_price_tier",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
							}
							// update only the compare price that is syncing to Shopify
							if pricing_compare_name.Value != "" {
								err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, pricing_compare_name.Value, product_variant.CompareAtPrice)
								if err != nil {
									log.Println(err)
									break
								}
							} else {
								if create_price_tier_enabled {
									// price tier is not set
									// use the default value of `fetch_compare_price`
									err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, "fetch_compare_price", product_variant.CompareAtPrice)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							// check if the product's inventory should be tracked
							if product_variant.InventoryManagement == "shopify" {
								shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
									"",
									fmt.Sprint(product_variant.InventoryItemID),
								)
								if err != nil {
									log.Println(err)
									break
								}
								// create map for warehouse quantity
								total_quantity := make(map[string]int)
								if len(shopify_inventory_levels.InventoryLevels) != 0 {
									for _, inventory_level := range shopify_inventory_levels.InventoryLevels {
										warehouse, err := dbconfig.DB.GetShopifyLocationByLocationID(
											context.Background(),
											fmt.Sprint(inventory_level.LocationID),
										)
										if err != nil {
											if err.Error() == "sql: no rows in result set" {
												continue
											}
											log.Println(err)
											break
										}
										total_quantity[warehouse.WarehouseName] = total_quantity[warehouse.WarehouseName] + inventory_level.Available
									}
								}
								// only fetch the ones that exist locally
								// TODO add any missing warehouses and location maps to database
								for warehouse_name, available := range total_quantity {
									_, err = dbconfig.DB.CreateVariantQty(context.Background(), database.CreateVariantQtyParams{
										ID:        uuid.New(),
										VariantID: db_variant.ID,
										Name:      warehouse_name,
										Isdefault: false,
										Value:     utils.ConvertIntToSQL(available),
										CreatedAt: time.Now().UTC(),
										UpdatedAt: time.Now().UTC(),
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
							sync_images_enabled := false
							sync_images_enabled_db, err := dbconfig.DB.GetAppSettingByKey(
								context.Background(),
								"app_fetch_sync_images",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									break
								}
								sync_images_enabled = false
							}
							sync_images_enabled, err = strconv.ParseBool(sync_images_enabled_db.Value)
							if err != nil {
								sync_images_enabled = false
							}
							// add shopify images to database
							if sync_images_enabled {
								for _, image := range product.Images {
									err = AddImagery(dbconfig, created_db_product.ID, image.Src, image.Position)
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
						}
					}
				}
			}
			log.Printf("From Shopify %d products were collected", len(shopifyProds.Products))
			fetch_url = utils.GetNextURL(next)
		}
	}
}

// Updates/Creates the specific price tier for
// a certain SKU
func AddPricing(
	dbconfig *DbConfig,
	sku string,
	variant_id uuid.UUID,
	pricing_name string,
	price string) error {
	exists, err := CheckExistsPriceTier(
		dbconfig,
		context.Background(),
		sku,
		pricing_name,
		false,
	)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateVariantPricing(context.Background(), database.UpdateVariantPricingParams{
			Name:      pricing_name,
			Value:     utils.ConvertStringToSQL(price),
			Isdefault: false,
			Sku:       sku,
			Name_2:    pricing_name,
		})
		if err != nil {
			return err
		}
	} else {
		_, err = dbconfig.DB.CreateVariantPricing(
			context.Background(),
			database.CreateVariantPricingParams{
				ID:        uuid.New(),
				VariantID: variant_id,
				Name:      pricing_name,
				Value:     utils.ConvertStringToSQL(price),
				Isdefault: false,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Updates/Creates an image for a certain product
func AddImagery(
	dbconfig *DbConfig,
	product_id uuid.UUID,
	image_url string,
	position int) error {
	exists, err := CheckExistsProductImage(
		dbconfig,
		context.Background(),
		product_id,
		image_url,
		position,
	)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateProductImage(context.Background(), database.UpdateProductImageParams{
			ImageUrl:  image_url,
			UpdatedAt: time.Now().UTC(),
			ProductID: product_id,
			Position:  int32(position),
		})
		if err != nil {
			return err
		}
	} else {
		err = dbconfig.DB.CreateProductImage(context.Background(), database.CreateProductImageParams{
			ID:        uuid.New(),
			ProductID: product_id,
			ImageUrl:  image_url,
			Position:  int32(position),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Updates/Creates a warehouse for a certain variant
func AddWarehouse(
	dbconfig *DbConfig,
	sku string,
	variant_id uuid.UUID,
	warehouse_name string,
	qty int) error {
	exists, err := CheckExistsWarehouse(dbconfig, context.Background(), sku, warehouse_name)
	if err != nil {
		return err
	}
	if exists {
		err = dbconfig.DB.UpdateVariantQty(context.Background(), database.UpdateVariantQtyParams{
			Name:      warehouse_name,
			Value:     utils.ConvertIntToSQL(qty),
			Isdefault: false,
			Sku:       sku,
			Name_2:    warehouse_name,
		})
		if err != nil {
			return err
		}
	} else {
		_, err = dbconfig.DB.CreateVariantQty(context.Background(), database.CreateVariantQtyParams{
			ID:        uuid.New(),
			VariantID: variant_id,
			Name:      warehouse_name,
			Isdefault: false,
			Value:     utils.ConvertIntToSQL(qty),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Checks if the Option1 has a default option
// and ignores it
func IgnoreDefaultOption(option string) string {
	if option == "Default Title" {
		return ""
	}
	return option
}
