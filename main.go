package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin",
		})
	})
	router.Run("localhost:3000")
}