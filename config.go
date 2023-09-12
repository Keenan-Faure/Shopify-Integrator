package main

import (
	"context"
	"database/sql"
	"integrator/internal/database"
	"log"
)

// Initiates a connection to the database and
// if successful returns the connection
func InitConn(dbURL string) (DbConfig, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return DbConfig{}, err
	}
	return storeConfig(db), nil
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
		config := DbConfig{
			DB:    database.New(conn),
			Valid: false,
		}
		return config
	}
}
