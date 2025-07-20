package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productService interface{}, parserService interface{}) {
	// Тестовый эндпоинт
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PortalData API is running on port 8095",
			"status":  "ok",
		})
	})

	// Эндпоинт для продуктов
	router.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"products": []gin.H{
				{"id": 1, "name": "Product 1", "price": 100},
				{"id": 2, "name": "Product 2", "price": 200},
			},
		})
	})

	// Эндпоинт для тестирования
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "API is working correctly",
			"port": "8095",
		})
	})
}
