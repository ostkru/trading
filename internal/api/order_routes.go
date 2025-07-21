package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.RouterGroup, handler *OrderHandlers) {
	router.GET("/orders/:id", handler.GetOrder)
	router.GET("/orders", handler.ListOrders)
	router.POST("/orders", handler.CreateOrder)
} 