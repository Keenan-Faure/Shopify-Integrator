package main

import (
	"context"
	"fetch"
	"integrator/internal/database"
	"log"
	"objects"
	"time"

	"github.com/google/uuid"
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
	for _, value := range products.Products {
		for _, sub_value := range value.Variants {
			_, err := dbconfig.DB.CreateShopifyProduct(context.Background(), database.CreateShopifyProductParams{
				ID:        uuid.New(),
				Title:     value.Title,
				Sku:       sub_value.Sku,
				Price:     sub_value.Price,
				Qty:       int32(sub_value.InventoryQuantity),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})

		}
	}
	log.Printf("From Shopify %d products were collected", len(products.Products))
}
