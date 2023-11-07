package main

import (
	"context"
	"fmt"
	"log"
	"objects"
	"shopify"
	"strconv"
	"time"
	"utils"
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
			log.Printf("From Shopify %d products were collected", len(shopifyProds.Products))
			log.Println(fetch_url)
			fetch_url = utils.GetNextURL(next)
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
