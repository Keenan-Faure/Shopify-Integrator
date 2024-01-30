package main

import (
	"context"
	"flag"
	"fmt"
	"integrator/internal/database"
	"iocsv"
	"log"
	"net/http"
	"shopify"
	"time"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	DB    *database.Queries
	Valid bool
}

const file_path = "./app"

func main() {
	fmt.Println("starting up app")
	// flags
	workers := flag.Bool("workers", false, "Enable server and worker for tests only")
	use_localhost := flag.Bool("localhost", false, "Enable localhost for tests only")
	flag.Parse()

	// db connnection config
	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	host := "@localhost:5432/"
	if !*use_localhost {
		host = "@postgres:5432/"
	}
	dbCon, err := InitConn(connection_string + host + utils.LoadEnv("db_name") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("error occured when setting up database: %v", err.Error())
	}
	// shopify connection config
	shopifyConfig := shopify.InitConfigShopify()

	// config workers only if flags are set
	if !*workers {
		go iocsv.LoopRemoveCSV()
		if shopifyConfig.Valid {
			go LoopJSONShopify(&dbCon, shopifyConfig)
		}
		QueueWorker(&dbCon)
		fmt.Println("resetting broken workers")
		err = dbCon.DB.ResetFetchWorker(context.Background(), "0")
		if err != nil {
			if err.Error()[0:12] != "pq: relation" {
				log.Fatalf("error occured %v", err.Error())
			}
		}
	}
	fmt.Println("starting API")
	setupAPI(dbCon, shopifyConfig)
}

// starts up the API
func setupAPI(dbconfig DbConfig, shopifyConfig shopify.ConfigShopify) {
	r := chi.NewRouter()
	r.Use(cors.Handler(MiddleWare()))
	api := chi.NewRouter()

	// utils
	api.Get("/ready", dbconfig.ReadyHandle)

	// customers
	api.Post("/customers", dbconfig.middlewareAuth(dbconfig.PostCustomerHandle))
	api.Get("/customers", dbconfig.middlewareAuth(dbconfig.CustomersHandle))
	api.Get("/customers/{id}", dbconfig.middlewareAuth(dbconfig.CustomerHandle))
	api.Get("/customers/search", dbconfig.middlewareAuth(dbconfig.CustomerSearchHandle))

	// orders
	api.Get("/orders", dbconfig.middlewareAuth(dbconfig.OrdersHandle))
	api.Get("/orders/{id}", dbconfig.middlewareAuth(dbconfig.OrderHandle))
	api.Get("/orders/search", dbconfig.middlewareAuth(dbconfig.OrderSearchHandle))
	api.Post("/orders", dbconfig.PostOrderHandle)

	// registration
	api.Post("/register", dbconfig.RegisterHandle)
	api.Post("/preregister", dbconfig.PreRegisterHandle)
	api.Post("/login", dbconfig.LoginHandle)
	api.Post("/logout", dbconfig.middlewareAuth(dbconfig.LogoutHandle))

	// products
	api.Post("/products/import", dbconfig.middlewareAuth(dbconfig.ProductImportHandle))
	api.Post("/products", dbconfig.middlewareAuth(dbconfig.PostProductHandle))
	api.Get("/products", dbconfig.middlewareAuth(dbconfig.ProductsHandle))
	api.Get("/products/{id}", dbconfig.middlewareAuth(dbconfig.ProductHandle))
	api.Get("/products/search", dbconfig.middlewareAuth(dbconfig.ProductSearchHandle))
	api.Get("/products/filter", dbconfig.middlewareAuth(dbconfig.ProductFilterHandle))
	api.Get("/products/export", dbconfig.middlewareAuth(dbconfig.ExportProductsHandle))
	api.Delete("/products/{id}", dbconfig.middlewareAuth(dbconfig.RemoveProductHandle))
	// api.Delete("/products/{variant_id}", dbconfig.middlewareAuth(dbconfig.RemoveProductVariantHandle))
	api.Put("/products/{id}", dbconfig.middlewareAuth(dbconfig.UpdateProductHandle))

	// general endpoint that returns the shopify_locations & internal warehouses
	api.Get("/inventory/config", dbconfig.middlewareAuth(dbconfig.ConfigLocationWarehouse))

	// shopify_location-internal warehouses map
	api.Get("/inventory/map", dbconfig.middlewareAuth(dbconfig.GetWarehouseLocations))
	api.Post("/inventory/map", dbconfig.middlewareAuth(dbconfig.AddWarehouseLocationMap))
	api.Delete("/inventory/map/{id}", dbconfig.middlewareAuth(dbconfig.RemoveWarehouseLocation))

	// internal warehouses
	api.Get("/inventory/warehouse", dbconfig.middlewareAuth(dbconfig.GetInventoryWarehouses))
	api.Get("/inventory/warehouse/{id}", dbconfig.middlewareAuth(dbconfig.GetInventoryWarehouse))
	api.Post("/inventory/warehouse", dbconfig.middlewareAuth(dbconfig.AddInventoryWarehouse))
	api.Delete("/inventory/warehouse/{id}", dbconfig.middlewareAuth(dbconfig.DeleteInventoryWarehouse))

	// shopify settings
	api.Get("/shopify/settings", dbconfig.middlewareAuth(dbconfig.GetShopifySettingValue))
	api.Put("/shopify/settings", dbconfig.middlewareAuth(dbconfig.AddShopifySetting))
	api.Delete("/shopify/settings", dbconfig.middlewareAuth(dbconfig.RemoveShopifySettings))

	// app settings
	api.Get("/settings", dbconfig.middlewareAuth(dbconfig.GetAppSettingValue))
	api.Put("/settings", dbconfig.middlewareAuth(dbconfig.AddAppSetting))
	api.Delete("/settings", dbconfig.middlewareAuth(dbconfig.RemoveAppSettings))

	// queue
	api.Get("/queue/{id}", dbconfig.middlewareAuth(dbconfig.GetQueueItemByID))
	api.Get("/queue", dbconfig.middlewareAuth(dbconfig.QueueViewNextItems))
	api.Get("/queue/filter", dbconfig.middlewareAuth(dbconfig.FilterQueueItems))
	api.Get("/queue/view", dbconfig.middlewareAuth(dbconfig.QueueView))
	api.Get("/queue/processing", dbconfig.middlewareAuth(dbconfig.QueueViewCurrentItem))
	api.Post("/queue", dbconfig.middlewareAuth(dbconfig.QueuePush))
	// api.Post("/queue/worker", dbconfig.middlewareAuth(dbconfig.QueuePopAndProcess))
	api.Delete("/queue/{id}", dbconfig.middlewareAuth(dbconfig.ClearQueueByID))
	api.Delete("/queue", dbconfig.middlewareAuth(dbconfig.ClearQueueByFilter))

	// not visible endpoints
	api.Get("/stats/fetch", dbconfig.middlewareAuth(dbconfig.GetFetchStats))
	api.Get("/stats/orders", dbconfig.middlewareAuth(dbconfig.GetOrderStats))

	// fetch handle
	api.Get("/shopify/fetch", dbconfig.middlewareAuth(dbconfig.WorkerFetchProductsHandle))

	// restrictions
	api.Put("/push/restriction", dbconfig.middlewareAuth(dbconfig.PushRestrictionHandle))
	api.Get("/push/restriction", dbconfig.middlewareAuth(dbconfig.GetPushRestrictionHandle))

	api.Put("/fetch/restriction", dbconfig.middlewareAuth(dbconfig.FetchRestrictionHandle))
	api.Get("/fetch/restriction", dbconfig.middlewareAuth(dbconfig.GetFetchRestrictionHandle))

	// Shopify
	api.Post("/shopify/sync", dbconfig.middlewareAuth(dbconfig.Synchronize))

	// webhooks
	api.Post("/shopify/webhook", dbconfig.middlewareAuth(dbconfig.PostWebhookHandle))
	api.Delete("/shopify/webhook", dbconfig.middlewareAuth(dbconfig.DeleteWebhookHandle))

	// OAuth2.0
	api.Get("/google/login", dbconfig.OAuthGoogleLogin)
	api.Get("/google/callback", dbconfig.OAuthGoogleCallback)
	api.Get("/google/oauth2/login", dbconfig.OAuthGoogleOAuth)

	r.Mount("/api", api)

	fs := http.FileServer(http.Dir(file_path))
	fsHandle := http.StripPrefix("/app", fs)
	r.Handle("/app", fsHandle)
	r.Handle("/app/*", fsHandle)

	port := utils.LoadEnv("port")
	if port == "" {
		log.Fatal("error port not defined in environment")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		// 10 second timeout
		ReadHeaderTimeout: 10 * time.Second,
	}
	fmt.Println("serving files from " + file_path + " and listening on port " + port)
	log.Fatal(server.ListenAndServe())
}
