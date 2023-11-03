package main

import (
	"context"
	"log"
	"objects"
	"shopify"
	"strconv"
	"time"
)

// loop function that uses Goroutine to run
// a function each interval
func LoopJSONShopify(
	dbconfig *DbConfig,
	shopifyConfig shopify.ConfigShopify) {
	fetch_time, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_shopify_fetch_time")
	if err != nil {
		log.Println(err)
	}
	timer, err := strconv.Atoi(fetch_time.Value)
	if err != nil {
		log.Println(err)
	}
	ticker := time.NewTicker(time.Duration(timer) * time.Second)
	for ; ; <-ticker.C {
		fetch_enabled, err := dbconfig.DB.GetAppSettingByKey(context.Background(), "app_enable_shopify_fetch")
		if err != nil {
			log.Println(err)
		}
		is_enabled, err := strconv.ParseBool(fetch_enabled.Value)
		if err != nil {
			log.Println(err)
		}
		if is_enabled {
			shopifyProds, err := shopifyConfig.FetchProducts()
			if err != nil {
				log.Println("Shopify > Error fetching next products to process:", err)
				continue
			}
			ProcessShopifyProducts(dbconfig, shopifyProds)
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
