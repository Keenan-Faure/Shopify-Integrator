package main

import (
	"context"
	"errors"
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
	fetch_time := 3
	fetch_time_db, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_shopify_fetch_time")
	if err != nil {
		fetch_time = 3
	}
	fetch_time, err = strconv.Atoi(fetch_time_db.Value)
	if err != nil {
		fetch_time = 3
	}
	// do not allow fetch time lower than 3 minutes
	if fetch_time < 3 {
		fetch_time = 3
	}
	ticker := time.NewTicker(time.Duration(fetch_time) * time.Minute)
	for ; ; <-ticker.C {
		err = FetchShopifyProducts(dbconfig, shopifyConfig)
		if err != nil {
			log.Println(err)
		}
	}
}

// fetch url to be stored inside the database
// this way the next fetch items can be paginated

// the fetch_url should be posted to the database
// along with the type of status that is currently used...

// the FectshopifyProducts should be used in the LoopJSONShopify function
// inside the loop...
func FetchShopifyProducts(dbconfig *DbConfig,
	shopifyConfig shopify.ConfigShopify) error {
	db_fetch_worker, err := dbconfig.DB.GetFetchWorker(context.Background())
	if err != nil {
		return err
	}
	fetch_url := db_fetch_worker.FetchUrl
	fetch_shopify_product_count := db_fetch_worker.ShopifyProductCount
	local_product_fetch_count := db_fetch_worker.LocalCount
	// during the first iteration it should fetch the count from shopify and update the counter
	if fetch_shopify_product_count == 0 {
		fetch_shopify_product_count_object, err := shopifyConfig.GetShopifyProductCount()
		if err != nil {
			log.Fatal("Shopify > Error fetching next products to process:", err)
		}
		// sets it to the counter
		fetch_shopify_product_count = int32(fetch_shopify_product_count_object.Count)
	} else {
		// otherwise if the value is non-zero (it has been set already)
		// then we check if the local counter equals the shopify count
		if fetch_shopify_product_count == local_product_fetch_count {
			// resets the url and the counters
			fetch_url = ""
			local_product_fetch_count = 0
			fetch_shopify_product_count = 0
		}
	}
	log.Println("running fetch worker...")
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
		// add check here to check if the worker is running or not
		status, err := dbconfig.DB.GetFetchWorker(context.Background())
		if err != nil {
			return err
		}
		if status.Status == "1" {
			return errors.New("worker is currently running")
		} else {
			err = dbconfig.DB.UpdateFetchWorker(context.Background(), database.UpdateFetchWorkerParams{
				Status:              "1",
				FetchUrl:            fetch_url,
				LocalCount:          local_product_fetch_count,
				ShopifyProductCount: fetch_shopify_product_count,
				UpdatedAt:           time.Now().UTC(),
				ID:                  db_fetch_worker.ID,
			})
			if err != nil {
				return err
			}
			shopifyProds, next, err := shopifyConfig.FetchProducts(fetch_url)
			if err != nil {
				return errors.New("Shopify > Error fetching next products to process: " + err.Error())
			}
			created_db_product := database.Product{}
			for _, product := range shopifyProds.Products {
				for variant_key, product_variant := range product.Variants {
					internal_product, err := dbconfig.DB.GetVariantBySKU(context.Background(), product_variant.Sku)
					if err != nil {
						if err.Error() != "sql: no rows in result set" {
							return err
						}
					}
					// if product is the same as the internal variant
					// then we will UPDATE the product
					// if the `app_fetch_overwrite_products` setting is enabled
					if internal_product.Sku == product_variant.Sku {
						// UPDATE PRODUCT: NOTE restrictions will only have effect
						// on product updates

						// fetch internal database restrictions for current iteration
						restrictions, err := dbconfig.DB.GetFetchRestriction(context.Background())
						if err != nil {
							return err
						}
						restrictions_map := FetchRestrictionsToMap(restrictions)

						overwrite := false
						overwrite_db, err := dbconfig.DB.GetAppSettingByKey(
							context.Background(),
							"app_fetch_overwrite_products",
						)
						if err != nil {
							if err.Error() != "sql: no rows in result set" {
								return err
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
									return err
								}
								if len(categories.CustomCollections) > 0 {
									category = categories.CustomCollections[0].Title
								}
								err = dbconfig.DB.UpdateProductBySKU(context.Background(), database.UpdateProductBySKUParams{
									Title:       utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, product.Title, "title")),
									BodyHtml:    utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, product.BodyHTML, "body_html")),
									Category:    utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, category, "category")),
									Vendor:      utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, product.Vendor, "vendor")),
									ProductType: utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, product.ProductType, "product_type")),
									UpdatedAt:   time.Now().UTC(),
									Sku:         product_variant.Sku,
								})
								if err != nil {
									return err
								}
								// update product options
								// check if product options should be updated
								// only if true should it be updated
								if DeterFetchRestriction(restrictions_map, "options") {
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
												return err
											}
										}
									}
								}
							}
							// update variant
							err = dbconfig.DB.UpdateVariant(context.Background(), database.UpdateVariantParams{
								Option1: utils.ConvertStringToSQL(
									ApplyFetchRestriction(
										restrictions_map,
										IgnoreDefaultOption(product_variant.Option1),
										"options",
									),
								),
								Option2: utils.ConvertStringToSQL(ApplyFetchRestriction(
									restrictions_map,
									IgnoreDefaultOption(product_variant.Option2),
									"options",
								)),
								Option3: utils.ConvertStringToSQL(ApplyFetchRestriction(
									restrictions_map,
									IgnoreDefaultOption(product_variant.Option3),
									"options",
								)),
								Barcode:   utils.ConvertStringToSQL(ApplyFetchRestriction(restrictions_map, product_variant.Barcode, "barcode")),
								UpdatedAt: time.Now().UTC(),
								Sku:       internal_product.Sku,
							})
							if err != nil {
								return err
							}
							// update variant pricing
							// check if pricing should be updated
							if DeterFetchRestriction(restrictions_map, "pricing") {
								if err != nil {
									if err.Error() != "sql: no rows in result set" {
										return err
									}
								}
								err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, "Selling Price", product_variant.Price)
								if err != nil {
									return err
								}
								err = AddPricing(dbconfig, internal_product.Sku, internal_product.ID, "Compare At Price", product_variant.Price)
								if err != nil {
									return err
								}
							}
							// check if the product's inventory should be tracked
							// check if the product's quantities should be updated
							if DeterFetchRestriction(restrictions_map, "warehousing") {
								if product_variant.InventoryManagement == "shopify" {
									shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
										"",
										fmt.Sprint(product_variant.InventoryItemID),
									)
									if err != nil {
										return err
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
												return err
											}
											total_quantity[warehouse.WarehouseName] = total_quantity[warehouse.WarehouseName] + inventory_level.Available
										}
									}
									// only update the warehouses that exist locally
									for warehouse_name, available := range total_quantity {
										err = dbconfig.DB.UpdateVariantQty(context.Background(), database.UpdateVariantQtyParams{
											Name:      warehouse_name,
											Value:     utils.ConvertIntToSQL(available),
											Isdefault: false,
											UpdatedAt: time.Now().UTC(),
											Sku:       product_variant.Sku,
											Name_2:    warehouse_name,
										})
										if err != nil {
											return err
										}
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
									return err
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
										return err
									}
								}
							}
						}
					} else {
						// PRODUCT CREATION: Note that the restrictions will have no impact
						// on a new product at all

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
									return err
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
									return err
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
											return err
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
								return err
							}
							err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, "Selling Price", product_variant.Price)
							if err != nil {
								return err
							}
							err = AddPricing(dbconfig, db_variant.Sku, db_variant.ID, "Compare At Price", product_variant.CompareAtPrice)
							if err != nil {
								return err
							}
							// check if the product's inventory should be tracked
							if product_variant.InventoryManagement == "shopify" {
								shopify_inventory_levels, err := shopifyConfig.GetShopifyInventoryLevels(
									"",
									fmt.Sprint(product_variant.InventoryItemID),
								)
								if err != nil {
									return err
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
											return err
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
										return err
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
									return err
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
										return err
									}
								}
							}
						}
					}
				}
			}
			err = dbconfig.DB.CreateFetchStat(context.Background(), database.CreateFetchStatParams{
				ID:               uuid.New(),
				AmountOfProducts: int32(len(shopifyProds.Products)),
				CreatedAt:        time.Now().UTC(),
				UpdatedAt:        time.Now().UTC(),
			})
			if err != nil {
				return err
			}
			log.Printf("From Shopify %d products were collected", len(shopifyProds.Products))
			local_product_fetch_count = int32(local_product_fetch_count) + int32(len(shopifyProds.Products))
			fetch_url = utils.GetNextURL(next)
			err = dbconfig.DB.UpdateFetchWorker(context.Background(), database.UpdateFetchWorkerParams{
				Status:              "0",
				FetchUrl:            fetch_url,
				LocalCount:          local_product_fetch_count,
				ShopifyProductCount: fetch_shopify_product_count,
				UpdatedAt:           time.Now().UTC(),
				ID:                  db_fetch_worker.ID,
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
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
			UpdatedAt: time.Now().UTC(),
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
	_, err := dbconfig.DB.GetWarehouseByName(context.Background(), warehouse_name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("warehouse " + warehouse_name + " not found")
		}
		return err
	}
	// if warehouse is found, we update the qty, we cannot create a new one
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
	return nil
}

// Creates/Updates product options for a certain product
func AddProductOptions(
	dbconfig *DbConfig,
	product_id uuid.UUID,
	product_code string,
	option_names []string,
) error {
	product_options, err := dbconfig.DB.GetProductOptions(context.Background(), product_id)
	if err != nil {
		return err
	}
	// product does not have any options
	if len(product_options) == 0 {
		for key, option_name := range option_names {
			if option_name != "" {
				_, err := dbconfig.DB.CreateProductOption(context.Background(), database.CreateProductOptionParams{
					ID:        uuid.New(),
					ProductID: product_id,
					Name:      option_name,
					Position:  int32(key + 1),
				})
				if err != nil {
					return err
				}
			}
		}
	} else {
		// product has options, we should update
		for key, option_name := range option_names {
			if option_name != "" {
				_, err := dbconfig.DB.UpdateProductOption(context.Background(), database.UpdateProductOptionParams{
					Name:       option_name,
					Position:   int32(key + 1),
					ProductID:  product_id,
					Position_2: int32(key + 1),
				})
				if err != nil {
					return err
				}
			}
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
