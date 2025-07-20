package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PortalData API is running on port 8095",
			"status":  "ok",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "API is working correctly",
			"port": "8095",
		})
	})

	fmt.Println("Starting server on port 8095...")
	if err := r.Run(":8095"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
