package main

import (
	"context"
	"flag"
	"fmt"
	"integrator/internal/database"
	"iocsv"
	"log"
	"shopify"
	"utils"

	"github.com/gin-gonic/gin"
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
	setUpAPI(&dbCon, &shopifyConfig)
}

func setUpAPI(dbconfig *DbConfig, shopifyconfig *shopify.ConfigShopify) {
	r := gin.Default()

	// use basic authentication
	r.Use(ApiKeyHeader(dbconfig))

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/ready", dbconfig.ReadyHandle())
	r.GET("/products/:id", dbconfig.ProductIDHandle())

	r.Run(":8080")
}
