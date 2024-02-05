package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TIPS
// URLParams | name := c.Param("name")
// Query Params | c.Query("lastname")

/*
Returns the product data having the specific id

Authorization: Basic, QueryParams, Headers

Response-Type: application/json

HTTP Codes: 200, 503
*/
func (dbconfig *DbConfig) ProductIDHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Param("id")
		err := IDValidation(product_id)
		if err != nil {
			RespondWithError(c, err, http.StatusInternalServerError)
			return
		}
		product_uuid, err := uuid.Parse(product_id)
		if err != nil {
			RespondWithError(c, errors.New("could not decode product id: "+product_id), http.StatusBadRequest)
			return
		}
		product_data, err := CompileProductData(dbconfig, product_uuid, c.Request.Context(), false)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, errors.New("not found"), http.StatusNotFound)
				return
			}
			RespondWithError(c, err, http.StatusInternalServerError)
			return
		}
		RespondWithJSON(c, product_data, http.StatusOK)
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
			RespondWithJSON(c, gin.H{"message": "OK"}, http.StatusOK)
		} else {
			RespondWithError(c, errors.New("Unavailable"), http.StatusServiceUnavailable)
		}
	}
}

func RespondWithError(c *gin.Context, err error, http_code int) {
	c.Error(err)
	c.AbortWithStatusJSON(http_code, gin.H{
		"message": err.Error(),
	})
}

func RespondWithJSON(c *gin.Context, payload any, http_code int) {
	c.JSON(http_code, payload)
}
