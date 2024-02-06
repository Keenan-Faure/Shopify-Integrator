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
func RespondWithError(c *gin.Context, http_code int, err_message string) {
	c.Error(errors.New(err_message))
	c.AbortWithStatusJSON(http_code, gin.H{
		"message": err_message,
	})
}

func RespondWithJSON(c *gin.Context, http_code int, payload any) {
	c.JSON(http_code, payload)
}
