package main

import (
	"context"
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

// loop function that uses Goroutine to run
// a function each interval
func LoopJSONShopify(
	dbconfig *DbConfig,
	shopifyConfig shopify.ConfigShopify) {
	fetch_url := ""
	fetch_time := 1
	fetch_time_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_shopify_fetch_time")
	if err != nil {
		fetch_time = 1
	}
	fetch_time, err = strconv.Atoi(fetch_time_db.Value)
	if err != nil {
		fetch_time = 1
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
							// Breaking out of loop for this iteration incase the database issue is temporary?
							log.Println(err)
							break
						}
						// product is the same as the internal variant
						if internal_product.Sku == product_variant.Sku {
							overwrite := false
							overwrite_db, err := dbconfig.DB.GetShopifySettingByKey(
								context.Background(),
								"app_fetch_overwrite_products",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									overwrite = false
								}
							}
							overwrite, err = strconv.ParseBool(overwrite_db.Value)
							if err != nil {
								log.Println(err)
								overwrite = false
							}
							if overwrite {
								// overwrite the product data once and then the variant data
								if variant_key == 0 {
									// update product here
									err = dbconfig.DB.UpdateProductBySKU(context.Background(), database.UpdateProductBySKUParams{
										Active:      "1",
										Title:       utils.ConvertStringToSQL(product.Title),
										BodyHtml:    utils.ConvertStringToSQL(product.BodyHTML),
										Vendor:      utils.ConvertStringToSQL(product.Vendor),
										ProductType: utils.ConvertStringToSQL(product.Type),
										UpdatedAt:   time.Now().UTC(),
										Sku:         product_variant.Sku,
									})
									// TODO Update product category as well?
									if err != nil {
										log.Println(err)
										break
									}
									// update product options
									if product.Options[0].Values[0] != "Default Title" {
										for _, option_value := range product.Options {
											_, err = dbconfig.DB.UpdateProductOption(
												context.Background(),
												database.UpdateProductOptionParams{
													Name:      option_value.Name,
													Position:  int32(option_value.Position),
													ProductID: internal_product.ProductID,
												},
											)
										}
									}
								}
								// update variant
								err = dbconfig.DB.UpdateVariant(context.Background(), database.UpdateVariantParams{
									Option1:   utils.ConvertStringToSQL(product_variant.Option1),
									Option2:   utils.ConvertStringToSQL(product_variant.Option2),
									Option3:   utils.ConvertStringToSQL(product_variant.Option3),
									Barcode:   utils.ConvertStringToSQL(product_variant.Barcode),
									UpdatedAt: time.Now().UTC(),
									Sku:       product_variant.Sku,
								})
								if err != nil {
									log.Println(err)
									break
								}
								// update variant pricing
								pricing_name, err := dbconfig.DB.GetShopifySettingByKey(
									context.Background(),
									"shopify_default_price_tier",
								)
								if err != nil {
									log.Println(err)
									break
								}
								if pricing_name.Value != "" {
									// update only the price that is syncing to Shopify
									err = dbconfig.DB.UpdateVariantPricing(context.Background(), database.UpdateVariantPricingParams{
										Name:      pricing_name.Value,
										Value:     utils.ConvertStringToSQL(product_variant.Price),
										Isdefault: false,
										Sku:       product_variant.Sku,
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
								pricing_cost_name, err := dbconfig.DB.GetShopifySettingByKey(
									context.Background(),
									"shopify_default_compare_at_price_tier",
								)
								if pricing_cost_name.Value != "" {
									// update only the cost price that is syncing to Shopify
									err = dbconfig.DB.UpdateVariantPricing(context.Background(), database.UpdateVariantPricingParams{
										Name:      pricing_cost_name.Value,
										Value:     utils.ConvertStringToSQL(product_variant.CompareAtPrice),
										Isdefault: false,
										Sku:       product_variant.Sku,
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
								if err != nil {
									log.Println(err)
									break
								}
								shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
									"",
									fmt.Sprint(product_variant.InventoryItemID),
								)
								if err != nil {
									log.Println(err)
									break
								}
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
										// create map for warehouse quantity
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
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
							}
						} else {
							// product variant can be created if setting is enabled
							create_fetched_product := false
							created_fetch_product_db, err := dbconfig.DB.GetAppSettingByKey(
								context.Background(),
								"app_fetch_add_product",
							)
							if err != nil {
								if err.Error() != "sql: no rows in result set" {
									log.Println(err)
									create_fetched_product = false
								}
							}
							create_fetched_product, err = strconv.ParseBool(created_fetch_product_db.Value)
							if err != nil {
								log.Println(err)
								create_fetched_product = false
							}
							if create_fetched_product {
								// create product only once during first iteration
								// creates product code to be the sku of the first variant
								if variant_key == 0 {
									db_product, err := dbconfig.DB.CreateProduct(context.Background(), database.CreateProductParams{
										ID:          uuid.New(),
										ProductCode: product_variant.Sku,
										Active:      "1",
										Title:       utils.ConvertStringToSQL(product.Title),
										BodyHtml:    utils.ConvertStringToSQL(product.BodyHTML),
										Category:    utils.ConvertStringToSQL(""),
										Vendor:      utils.ConvertStringToSQL(product.Vendor),
										ProductType: utils.ConvertStringToSQL(product.Type),
										CreatedAt:   time.Now().UTC(),
										UpdatedAt:   time.Now().UTC(),
									})
									if err != nil {
										log.Println(err)
										break
									}
									created_db_product = db_product
									// create product options
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
								// then create this variant and any other variant to come
								db_variant, err := dbconfig.DB.CreateVariant(
									context.Background(),
									database.CreateVariantParams{
										ID:        uuid.New(),
										ProductID: created_db_product.ID,
										Sku:       product_variant.Sku,
										Option1:   utils.ConvertStringToSQL(product_variant.Option1),
										Option2:   utils.ConvertStringToSQL(product_variant.Option2),
										Option3:   utils.ConvertStringToSQL(product_variant.Option3),
										Barcode:   utils.ConvertStringToSQL(product_variant.Barcode),
										CreatedAt: time.Now().UTC(),
										UpdatedAt: time.Now().UTC(),
									},
								)
								if err != nil {
									log.Println(err)
									break
								}
								// update variant pricing
								pricing_name, err := dbconfig.DB.GetShopifySettingByKey(
									context.Background(),
									"shopify_default_price_tier",
								)
								if err != nil {
									log.Println(err)
									break
								}
								if pricing_name.Value != "" {
									// update only the price that is syncing to Shopify
									_, err = dbconfig.DB.CreateVariantPricing(context.Background(), database.CreateVariantPricingParams{
										ID:        uuid.New(),
										VariantID: db_variant.ID,
										Name:      pricing_name.Value,
										Value:     utils.ConvertStringToSQL(product_variant.Price),
										Isdefault: false,
										CreatedAt: time.Now().UTC(),
										UpdatedAt: time.Now().UTC(),
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
								pricing_cost_name, err := dbconfig.DB.GetShopifySettingByKey(
									context.Background(),
									"shopify_default_compare_at_price_tier",
								)
								if pricing_cost_name.Value != "" {
									// update only the cost price that is syncing to Shopify
									_, err = dbconfig.DB.CreateVariantPricing(context.Background(), database.CreateVariantPricingParams{
										ID:        uuid.New(),
										VariantID: db_variant.ID,
										Name:      pricing_cost_name.Value,
										Value:     utils.ConvertStringToSQL(product_variant.CompareAtPrice),
										Isdefault: false,
										CreatedAt: time.Now().UTC(),
										UpdatedAt: time.Now().UTC(),
									})
									if err != nil {
										log.Println(err)
										break
									}
								}
								if err != nil {
									log.Println(err)
									break
								}
								//////// Create Warehouses
								shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
									"",
									fmt.Sprint(product_variant.InventoryItemID),
								)
								if err != nil {
									log.Println(err)
									break
								}
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
										// create map for warehouse quantity
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
						}
					}
				}
				// steps
				// 1. Check if the SKU exist on the system
				// if it does then check settings if overwrite is enabled
				// if it does not then check settins if create_product_on_fetch is enabled
				// if it is then create product
				// 2. Update product
				log.Printf("From Shopify %d products were collected", len(shopifyProds.Products))
				log.Println(fetch_url)
				fetch_url = utils.GetNextURL(next)
			}
		}
	}
}

// adds shopify products to the database
func ProcessShopifyProducts(dbconfig *DbConfig, products objects.ShopifyProducts) {
	// for _, value := range products.Products {
	// 	for _, sub_value := range value.Variants {
	// 		//create product in database
	// 		// if error = unique value not allowed etc
	// 		// override
	// 	}
	// }
	log.Printf("From Shopify %d products were collected", len(products.Products))
}
