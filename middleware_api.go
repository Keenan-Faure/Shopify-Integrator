package main

import (
	"integrator/internal/database"
	"net/http"
	"utils"
)

// custom Authhandler
type authHandler func(w http.ResponseWriter, r *http.Request, dbuser database.User)

// Authentication middleware
func (dbconfig *DbConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := utils.ExtractAPIKey(r.Header.Get("Authorization"))
		if apiKey == "" {
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		dbUser, err := dbconfig.DB.GetUserByApiKey(r.Context(), []byte(apiKey))
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(w, http.StatusNotFound, "record not found")
				return
			}
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		handler(w, r, dbUser)
	}
}
