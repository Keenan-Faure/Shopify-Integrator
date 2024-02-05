package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Middleware that checks if the request is using Basic Authentication.
The username and password values needs to be passed in the headers of the request
*/
func Basic() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()
		if hasAuth {
			// validation
		} else {
			err := errors.New("request did not meet authentication standards")
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		}
	}
}

/*
Middleware that checks if the request is authenticating sending the api_key as a query param

Format: {{base_url}}/{{resource}}?api_key={{api_key}}
*/
func QueryParams() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

/*
Middleware that checks if the request is sending the ApiKey
inside the headers

Format: ApiKey {{api_key}}
*/
func ApiKeyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
