package main

import (
	"github.com/gin-gonic/gin"
)

func initializeRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})
	return r
}
