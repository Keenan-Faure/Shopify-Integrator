package main

import (
	"github.com/gin-gonic/gin"
)

/*
Confirms if the API is ready to start accepting requests.

Authorization: None

Response-Type: application/json

HTTP Codes: 200, 503
*/
func (dbconfig *DbConfig) ReadyHandle() {
	r := gin.Default()
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	r.Run()
}
