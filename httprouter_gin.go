package main

import (
	"errors"
	"integrator/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
Returns the product data having the specific id

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

HTTP Codes: 200, 400, 401, 500
*/
func (dbconfig *DbConfig) ProductIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode product id '"+product_id+"'")
			return
		}
		product_data, err := CompileProductData(dbconfig, product_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusOK, product_data)
	}
}

/*
Registers a new user

Authorization: None

Response-Type: application/json

HTTP Codes: 200, 503
*/

func (dbconfig *DbConfig) RegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := DecodeUserRequestBody(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		err = ValidateTokenValidation(body)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		token, err := dbconfig.DB.GetTokenValidation(c.Request.Context(), body.Email)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		request_token, err := uuid.Parse(body.Token)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, "could not decode token: "+body.Token)
			return
		}
		if token.Token != request_token {
			RespondWithError(c, http.StatusNotFound, "invalid token for user")
			return
		}
		err = UserValidation(body.Name, body.Password)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		exists, err := dbconfig.CheckUserExist(body.Name, c.Request)
		if exists {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		user, err := dbconfig.DB.CreateUser(c.Request.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      body.Name,
			UserType:  "app",
			Email:     body.Email,
			Password:  body.Password,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(c, http.StatusCreated, ConvertDatabaseToRegister(user))
	}
}

/*
Confirms if the API is ready to start accepting requests.

Authorization: None

Response-Type: application/json

HTTP Codes: 200, 503
*/
func (dbconfig *DbConfig) ReadyHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if dbconfig.Valid {
			RespondWithJSON(c, http.StatusOK, gin.H{"message": "OK"})
		} else {
			RespondWithError(c, http.StatusServiceUnavailable, "Unavailable")
		}
	}
}

// Helper function
// logs all error messages in current context to stdout
// the error message in the parameters is returned over the API
// after the chain has been aborted.
func RespondWithError(c *gin.Context, http_code int, err_message string) {
	for _, err := range c.Errors {
		// TODO log previous errors from the authentication middlewares inside database table
		log.Println(err.Err.Error())
		break
	}
	c.AbortWithStatusJSON(http_code, gin.H{
		"message": err_message,
	})
}

// Helper function
// responds with a payload and http code over the API
// after sucessfully processing the request.
func RespondWithJSON(c *gin.Context, http_code int, payload any) {
	c.JSON(http_code, payload)
}

// Helper function
// appends the error to the current context and passes on to the next Middleware
// only used in the Authentication phase when one middleware auth is not met
// it uses another in the format below
// query_param -> api_key as header -> basic
// in the case that the last auth fails, it raises and error with RespondWithError
func AppendErrorNext(c *gin.Context, http_code int, err_message string) {
	c.Error(errors.New(err_message))
	c.Next()
}
