package main

import (
	"flag"
	"fmt"
	"integrator/internal/database"
	"log"
	"net/http"
	"shopify"
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
	dbCon, err := InitConn(utils.LoadEnv("docker_db_url") + utils.LoadEnv("database") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}
	flags := flag.Bool("test", false, "Enable server for tests only")
	flag.Parse()

	shopifyConfig := shopify.InitConfigShopify(
		utils.LoadEnv("store_name"),
		utils.LoadEnv("api_key"),
		utils.LoadEnv("api_password"),
		utils.LoadEnv("api_version"),
	)
	if !*flags {
		fmt.Println("Starting Workers")
		// go iocsv.LoopRemoveCSV()
		// go shopify.LoopJSONShopify()
	}
	fmt.Println("Starting API")
	fmt.Println(shopifyConfig.GetProductBySKU("GenImp-R-EC"))
	setupAPI(dbCon)
}

// starts up the API
func setupAPI(dbconfig DbConfig) {
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
	api.Get("/endpoints", dbconfig.EndpointsHandle)
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
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Serving files from %s and listening on port %s", file_path, port)
	log.Fatal(server.ListenAndServe())
}
