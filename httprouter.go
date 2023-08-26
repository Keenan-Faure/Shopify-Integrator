package main

import (
	"encoding/json"
	"net/http"
	"objects"

	"github.com/go-chi/cors"
)

// Determines the readiness of the API
func (dbconfig *DbConfig) readyHandle(w http.ResponseWriter, r *http.Response) {
	if dbconfig.Valid {
		RespondWithJSON(w, 200, objects.ResponseString{
			Status: "OK",
		})
	} else {
		RespondWithJSON(w, 503, objects.ResponseString{
			Status: "Error",
		})
	}
}

// JSON helper functions
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}

// Middleware that determines which headers, http methods and orgins are allowed
func MiddleWare() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}
}
