package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handlers) {
	router.GET("/orders/:id", handler.GetOrder)
	router.GET("/orders", handler.ListOrders)
	router.POST("/orders", handler.CreateOrder)
}
