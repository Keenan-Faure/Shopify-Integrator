package main

import (
	"fetch"
	"log"
	"objects"
	"time"
)

const fetch_time_shopify = 120 * time.Second // 120 seconds

// loop function that uses Goroutine to run
// a function each interval
func LoopJSONShopify(
	dbconfig *DbConfig,
	shopifyConfig fetch.ConfigShopify,
	interval time.Duration) {
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		shopifyProds, err := shopifyConfig.FetchProducts()
		if err != nil {
			log.Println("Shopify > Error fetching next products to process:", err)
			continue
		}
		ProcessShopifyProducts(dbconfig, shopifyProds)
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
