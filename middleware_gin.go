package main

import (
	"errors"
	"net/http"
	"objects"

	"github.com/gin-gonic/gin"
)

/*
Middleware that checks if the request is using Basic Authentication.
The username and password values needs to be passed in the headers of the request
*/
func Basic(dbconfig *DbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()
		if hasAuth {
			_, exists, err := dbconfig.CheckUserCredentials(objects.RequestBodyLogin{
				Username: user,
				Password: password,
			}, c.Request)
			if err != nil {
				RespondWithError(c, err, http.StatusUnauthorized)
				return
			}
			if !exists {
				RespondWithError(c, errors.New("invalid username or password combination"), http.StatusUnauthorized)
				return
			}
		} else {
			RespondWithError(c, errors.New("no authentication found in request"), http.StatusBadRequest)
			return
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
