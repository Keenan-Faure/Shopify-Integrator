package main

import (
	"net/http"
	"objects"
	"utils"

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
				RespondWithError(c, http.StatusUnauthorized, err.Error())
				return
			}
			if !exists {
				RespondWithError(c, http.StatusUnauthorized, "invalid username or password combination")
				return
			}
		} else {
			RespondWithError(c, http.StatusBadRequest, "no authentication found in request")
			return
		}
	}
}

/*
Middleware that checks if the request is authenticating sending the api_key as a query param

Format: {{base_url}}/{{resource}}?api_key={{api_key}}
*/
func QueryParams(dbconfig *DbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.Query("api_key")
		_, err := dbconfig.DB.GetUserByApiKey(c.Request.Context(), api_key)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "user not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

/*
Middleware that checks if the request is sending the ApiKey
inside the headers

Format: ApiKey {{api_key}}
*/
func ApiKeyHeader(dbconfig *DbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key, err := utils.ExtractAPIKey(c.Request.Header["Authorization"][0]) // uses the first Authorization Header in request
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, err.Error())
		}
		_, err = dbconfig.DB.GetUserByApiKey(c.Request.Context(), api_key)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				RespondWithError(c, http.StatusNotFound, "user not found")
				return
			}
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
