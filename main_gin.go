package main

import (
	"context"
	"flag"
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
		log.Fatalf("error occured when setting up database: %v", err.Error())
	}

	shopifyConfig := shopify.InitConfigShopify()

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
	setUpAPI(&dbCon, &shopifyConfig)
}

func setUpAPI(dbconfig *DbConfig, shopifyconfig *shopify.ConfigShopify) {
	r := gin.Default()

	// authentication methods
	// hover for more details
	// Middleware runs in the format specficied
	// query_params -> api_keys inside header -> Basic authentication
	r.Use(QueryParams(dbconfig))
	r.Use(ApiKeyHeader(dbconfig))
	r.Use(Basic(dbconfig))

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.POST("/preregister", dbconfig.PreRegisterHandle())
	r.POST("/register", dbconfig.RegisterHandle())
	r.POST("/logout", dbconfig.LogoutHandle())
	r.POST("/login", dbconfig.LoginHandle())

	r.GET("/ready", dbconfig.ReadyHandle())

	r.Run(":8080")
}
