package main

import (
	"flag"
	"fmt"
	"integrator/internal/database"
	"log"
	"net/http"
	"utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	DB    *database.Queries
	Valid bool
}

const file_path = "./app"

func main() {
	dbCon, err := InitConn(utils.LoadEnv("db_url"))
	if err != nil {
		log.Fatalf("Error occured %v", err.Error())
	}
	flags := flag.Bool("test", false, "Enable server for tests only")
	flag.Parse()

	if !*flags {
		fmt.Println("Starting Worker")
	}
	fmt.Println("Starting API")
	setupAPI(&dbCon)
}

// starts up the API
func setupAPI(dbconfig *DbConfig) {
	r := chi.NewRouter()
	r.Use(cors.Handler(MiddleWare()))

	api := chi.NewRouter()
	api.Mount("/api", api)

	api.Post("/register", dbconfig.RegisterHandle)
	api.Post("/login", dbconfig.middlewareAuth(dbconfig.LoginHandle))
	api.Get("/endpoints", dbconfig.EndpointsHandle)

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
