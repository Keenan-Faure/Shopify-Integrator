package main

import (
	"context"
	"database/sql"
	"integrator/internal/database"
	"log"
	"utils"

	"github.com/jackc/pgx/v5"
)

// Initiates a connection to the database and
// if successful returns the connection
func InitConn(dbURL string) (DbConfig, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
		return DbConfig{}, err
	}
	return storeConfig(db), nil
}

// Initiates a connection to the database
// used only for custom queries that do not
// exist inside the /internal/database
func InitCustomConnection(dbURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return &pgx.Conn{}, err
	}
	return conn, err
}

// Creates a connection string
func InitConnectionString(useLocalhost, mock bool) string {
	connection_string := "postgres://" + utils.LoadEnv("db_user") + ":" + utils.LoadEnv("db_psw")
	host := "@localhost:5432/"
	if !useLocalhost {
		host = "@postgres:5432/"
	}
	if mock {
		host = "@localhost:5432/"
	}
	log.Println("USR: " + utils.LoadEnv("db_user"))
	log.Println("PSS: " + utils.LoadEnv("db_psw"))
	return connection_string + host + utils.LoadEnv("db_name")
}

// Stores the database connection inside a config struct
func storeConfig(conn *sql.DB) DbConfig {
	_, err := database.New(conn).GetUsers(context.Background())
	if err == nil {
		config := DbConfig{
			DB:    database.New(conn),
			Valid: true,
		}
		return config
	} else {
		if err.Error() == "sql: no rows in result set" {
			config := DbConfig{
				DB:    database.New(conn),
				Valid: true,
			}
			return config
		}
		config := DbConfig{
			DB:    database.New(conn),
			Valid: false,
		}
		return config
	}
}
