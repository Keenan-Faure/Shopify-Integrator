package main

import (
	"context"
	"flag"
	"integrator/internal/database"
	"iocsv"
	"log"
	"shopify"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	DB    *database.Queries
	Valid bool
}

func main() {
	workers := flag.Bool("workers", false, "Enable server and worker for tests only")
	connectionString := InitConnectionString(false)
	dbCon, err := InitConn(connectionString)
	if err != nil {
		log.Fatalf("error occured when setting up database: %v", err.Error())
	}

	shopifyConfig := shopify.InitConfigShopify("")

	if !*workers {
		go iocsv.LoopRemoveCSV()
		if shopifyConfig.Valid {
			go LoopJSONShopify(&dbCon, shopifyConfig)
		}
		QueueWorker(&dbCon)
		err = dbCon.DB.ResetFetchWorker(context.Background(), "0")
		if err != nil {
			if err.Error()[0:12] != "pq: relation" {
				log.Fatalf("error occured %v", err.Error())
			}
		}
	}
	r := setUpAPI(&dbCon)
	r.Run(":8080")
}

func setUpAPI(dbconfig *DbConfig) *gin.Engine {
	r := gin.Default()

	// authentication methods
	// hover for more details
	// Middleware runs in the format specficied
	// query_params -> api_keys inside header -> Basic authentication

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	/* --------- N/A Auth routes --------- */

	nauth := r.Group("/api")

	nauth.GET("/ready", dbconfig.ReadyHandle())
	nauth.POST("/preregister", dbconfig.PreRegisterHandle())
	nauth.POST("/register", dbconfig.RegisterHandle())
	nauth.POST("/login", dbconfig.LoginHandle())

	/* --------- Auth routes --------- */
	auth := r.Group("/api")

	auth.Use(QueryParams(dbconfig))
	auth.Use(ApiKeyHeader(dbconfig))
	auth.Use(Basic(dbconfig))

	auth.POST("/logout", dbconfig.LogoutHandle())

	/* Products */
	auth.GET("/products", dbconfig.ProductsHandle())
	auth.GET("/products/:id", dbconfig.ProductIDHandle())
	auth.GET("/products/search", dbconfig.ProductSearchHandle())
	auth.GET("/products/filter", dbconfig.ProductFilterHandle())

	auth.PUT("/products/:id", dbconfig.UpdateProductHandle())

	auth.POST("/products", dbconfig.PostProductHandle())
	auth.POST("/products/import", dbconfig.ProductImportHandle())
	auth.POST("/products/export", dbconfig.ProductExportHandle())

	auth.DELETE("/products/:id/variants/:variant_id", dbconfig.RemoveProductVariantHandle())
	auth.DELETE("/products/:id", dbconfig.RemoveProductHandle())

	/* Orders */
	auth.GET("/orders", dbconfig.OrdersHandle())
	auth.GET("/orders/:id", dbconfig.OrderIDHandle())
	auth.GET("/orders/search", dbconfig.OrderSearchHandle())

	auth.POST("/orders", dbconfig.PostOrderHandle())

	/* Customers */
	auth.GET("/customers", dbconfig.CustomersHandle())
	auth.GET("/customers/:id", dbconfig.CustomerIDHandle())
	auth.GET("/customers/search", dbconfig.CustomerSearchHandle())

	auth.POST("/customers", dbconfig.PostCustomerHandle())

	/* Inventory Config Handle */
	auth.GET("/inventory/config", dbconfig.ConfigLocationWarehouseHandle())

	/* Inventory Map */
	auth.GET("/inventory/map", dbconfig.LocationWarehouseHandle())
	auth.POST("/inventory/map", dbconfig.AddWarehouseLocationMap())
	auth.DELETE("/inventory/map/:id", dbconfig.RemoveWarehouseLocation())

	/* Statistics */
	auth.GET("/stats/orders", dbconfig.GetOrderStats())
	auth.GET("/stats/fetch", dbconfig.GetFetchStats())

	/* Inventory Warehouses */
	auth.GET("/inventory/warehouse/:id", dbconfig.GetInventoryWarehouse())
	auth.GET("/inventory/warehouse", dbconfig.GetInventoryWarehouses())
	auth.POST("/inventory/warehouse", dbconfig.AddInventoryWarehouseHandle())
	auth.DELETE("/inventory/warehouse/:id", dbconfig.DeleteInventoryWarehouse())

	/* Fetch Worker */
	auth.GET("/worker/fetch", dbconfig.WorkerFetchProductsHandle())

	/* Restrictions */
	auth.GET("/fetch/restriction", dbconfig.GetFetchRestrictionHandle())
	auth.PUT("/fetch/restriction", dbconfig.FetchRestrictionHandle())

	auth.GET("/push/restriction", dbconfig.GetPushRestrictionHandle())
	auth.PUT("/push/restriction", dbconfig.PushRestrictionHandle())

	/* Webhooks */
	auth.GET("/shopify/webhook", dbconfig.AddWebhookHandle())
	auth.DELETE("/shopify/webhook", dbconfig.DeleteWebhookHandle())

	/* Settings */
	auth.GET("/shopify/settings", dbconfig.GetShopifySettingValue())
	auth.PUT("/shopify/settings", dbconfig.AddShopifySetting())

	auth.GET("/settings", dbconfig.GetAppSettingValue())
	auth.PUT("/settings", dbconfig.AddAppSetting())

	/* Queue */
	auth.GET("/queue/:id", dbconfig.GetQueueItemByID())
	auth.GET("/queue", dbconfig.QueueViewNextItems())
	auth.GET("/queue/filter", dbconfig.FilterQueueItems())
	auth.GET("/queue/view", dbconfig.QueueView())
	auth.GET("/queue/processing", dbconfig.QueueViewCurrentItem())
	auth.POST("/queue", dbconfig.QueuePush())
	auth.DELETE("/queue", dbconfig.ClearQueueByFilter())
	auth.DELETE("/queue/:id", dbconfig.ClearQueueByID())

	/* OAuth2.0 */
	auth.GET("/google/login", dbconfig.OAuthGoogleLogin())
	auth.GET("/google/callback", dbconfig.OAuthGoogleCallback())
	auth.GET("/google/oauth2/login", dbconfig.OAuthGoogleOAuth())

	/* setup file server */
	r.StaticFS("/static", gin.Dir("./app/export", true))

	return r
}
