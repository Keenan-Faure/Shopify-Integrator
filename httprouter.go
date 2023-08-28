package main

import (
	"api"
	"encoding/json"
	"integrator/internal/database"
	"net/http"
	"objects"
	"time"

	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

// POST /api/login
func (dbconfig *DbConfig) LoginHandle(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/register
func (dbconfig *DbConfig) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeUserRequestBody(r)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if UserValidation(body) != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	exists, err := dbconfig.CheckUserExist(body.Name, r)
	if exists {
		RespondWithError(w, http.StatusConflict, err.Error())
		return
	}
	user, err := dbconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New().String(),
		Name:      body.Name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		ApiKey:    uuid.New().String(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, user)
}

// GET /api/endpoints
func (dbconfig *DbConfig) EndpointsHandle(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, api.Endpoints())
}

// puppy function for a handler
//func Handle{{Name}} {
// 1.) decode params
// 2.) validate requestBody
// 3.) Check for duplicates in db
// 4.) (fi) create/insert new record
// 4.) (if) return record
// 5.) (else) return error
//}

// GET /api/ready
func (dbconfig *DbConfig) ReadyHandle(w http.ResponseWriter, r *http.Response) {
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
