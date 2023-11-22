package main

import (
	"flag"
	"fmt"
	"integrator/internal/database"
	"iocsv"
	"log"
	"net/http"
	"os"
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
	fmt.Println(os.Getwd())
	workers := flag.Bool("workers", false, "Enable server and worker for tests only")
	use_localhost := flag.Bool("localhost", false, "Enable localhost for tests only")
	flag.Parse()

	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	host := "@localhost:5432/"
	if !*use_localhost {
		host = "@postgres:5432/"
	}
	dbCon, err := InitConn(connection_string + host + utils.LoadEnv("db_name") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}

	shopifyConfig := shopify.InitConfigShopify()
	if !*workers {
		fmt.Println("Starting Workers")
		go iocsv.LoopRemoveCSV()
		if shopifyConfig.Valid {
			go LoopJSONShopify(&dbCon, shopifyConfig)
		}
		QueueWorker(&dbCon)
	}
	fmt.Println("Starting API")
	setupAPI(dbCon, shopifyConfig)
}

// starts up the API
func setupAPI(dbconfig DbConfig, shopifyConfig shopify.ConfigShopify) {
	r := chi.NewRouter()
	r.Use(cors.Handler(MiddleWare()))
	api := chi.NewRouter()

	api.Post("/products/import", dbconfig.middlewareAuth(dbconfig.ProductImportHandle))
	api.Post("/products", dbconfig.middlewareAuth(dbconfig.PostProductHandle))
	api.Post("/customers", dbconfig.middlewareAuth(dbconfig.PostCustomerHandle))
	api.Post("/orders", dbconfig.middlewareAuth(dbconfig.PostOrderHandle))
	api.Post("/register", dbconfig.RegisterHandle)
	api.Post("/preregister", dbconfig.PreRegisterHandle)
	api.Post("/login", dbconfig.middlewareAuth(dbconfig.LoginHandle))
	api.Get("/ready", dbconfig.ReadyHandle)
	api.Get("/products", dbconfig.middlewareAuth(dbconfig.ProductsHandle))
	api.Get("/products/{id}", dbconfig.middlewareAuth(dbconfig.ProductHandle))
	api.Get("/products/search", dbconfig.middlewareAuth(dbconfig.ProductSearchHandle))
	api.Get("/products/filter", dbconfig.middlewareAuth(dbconfig.ProductFilterHandle))
	api.Get("/orders", dbconfig.middlewareAuth(dbconfig.OrdersHandle))
	api.Get("/orders/{id}", dbconfig.middlewareAuth(dbconfig.OrderHandle))
	api.Get("/orders/search", dbconfig.middlewareAuth(dbconfig.OrderSearchHandle))
	api.Get("/customers", dbconfig.middlewareAuth(dbconfig.CustomersHandle))
	api.Get("/customers/{id}", dbconfig.middlewareAuth(dbconfig.CustomerHandle))
	api.Get("/customers/search", dbconfig.middlewareAuth(dbconfig.CustomerSearchHandle))
	api.Get("/products/export", dbconfig.middlewareAuth(dbconfig.ExportProductsHandle))
	api.Delete("/products/{id}", dbconfig.middlewareAuth(dbconfig.RemoveProductHandle))
	api.Delete("/products/{variant_id}", dbconfig.middlewareAuth(dbconfig.RemoveProductVariantHandle))

	api.Get("/inventory", dbconfig.middlewareAuth(dbconfig.GetWarehouseLocations))
	api.Post("/inventory", dbconfig.middlewareAuth(dbconfig.AddWarehouseLocationMap))
	api.Delete("/inventory/{id}", dbconfig.middlewareAuth(dbconfig.RemoveWarehouseLocation))

	// shopify settings
	api.Get("/shopify/settings", dbconfig.middlewareAuth(dbconfig.GetShopifySettingValue))
	api.Post("/shopify/settings", dbconfig.middlewareAuth(dbconfig.AddShopifySetting))
	api.Delete("/shopify/settings", dbconfig.middlewareAuth(dbconfig.RemoveShopifySettings))

	// app settings
	api.Get("/settings", dbconfig.middlewareAuth(dbconfig.GetAppSettingValue))
	api.Post("/settings", dbconfig.middlewareAuth(dbconfig.AddAppSetting))
	api.Delete("/settings", dbconfig.middlewareAuth(dbconfig.RemoveAppSettings))

	// queue
	api.Get("/queue/{id}", dbconfig.middlewareAuth(dbconfig.GetQueueItemByID))
	api.Get("/queue", dbconfig.middlewareAuth(dbconfig.QueueViewNextItems))
	api.Get("/queue/filter", dbconfig.middlewareAuth(dbconfig.FilterQueueItems))
	api.Get("/queue/view", dbconfig.middlewareAuth(dbconfig.QueueView))
	api.Get("/queue/processing", dbconfig.middlewareAuth(dbconfig.QueueViewCurrentItem))
	api.Post("/queue", dbconfig.middlewareAuth(dbconfig.QueuePush))
	api.Post("/shopify/sync", dbconfig.middlewareAuth(dbconfig.Synchronize))
	// api.Post("/queue/worker", dbconfig.middlewareAuth(dbconfig.QueuePopAndProcess))
	api.Delete("/queue/{id}", dbconfig.middlewareAuth(dbconfig.ClearQueueByID))
	api.Delete("/queue", dbconfig.middlewareAuth(dbconfig.ClearQueueByFilter))

	r.Mount("/api", api)

	fs := http.FileServer(http.Dir(file_path))
	fsHandle := http.StripPrefix("/app", fs)
	r.Handle("/app", fsHandle)
	r.Handle("/app/*", fsHandle)

	port := utils.LoadEnv("port")
	if port == "" {
		log.Fatal("Port not defined in Environment")
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s and listening on port %s", file_path, port)
	log.Fatal(server.ListenAndServe())
}
