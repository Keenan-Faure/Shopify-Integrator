package main

import (
	"net/http"
	"objects"
	"utils"

	"github.com/gin-gonic/gin"
)

/*
Middleware that checks if the request is using Basic Authentication.
*/
func Basic(dbconfig *DbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetString("api_key")
		if api_key != "" {
			c.Next()
		} else {
			user, password, hasAuth := c.Request.BasicAuth()
			if hasAuth {
				_, exists, err := CheckUserCredentials(dbconfig, objects.RequestBodyLogin{
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
				RespondWithError(c, http.StatusUnauthorized, "no authentication found in request")
				return
			}
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
		db_user, err := dbconfig.DB.GetUserByApiKey(c.Request.Context(), api_key)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.Next()
				return
			}
			AppendErrorNext(c, http.StatusInternalServerError, err.Error())
			return
		} else {
			c.Set("api_key", db_user.ApiKey)
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
		api_key := c.GetString("api_key")
		if api_key != "" {
			c.Next()
		}
		auth_headers := c.Request.Header["Authorization"]
		if len(auth_headers) > 0 {
			api_key, err := utils.ExtractAPIKey(auth_headers[0]) // uses the first Authorization Header in request
			if err != nil {
				c.Next()
				return
			}
			db_user, err := dbconfig.DB.GetUserByApiKey(c.Request.Context(), api_key)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					c.Next()
					return
				}
				AppendErrorNext(c, http.StatusInternalServerError, err.Error())
				return
			} else {
				c.Set("api_key", db_user.ApiKey)
			}
		}
	}
}
