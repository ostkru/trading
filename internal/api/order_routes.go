package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, h *OrderHandlers) {
	r.GET("/orders/:id", h.GetOrder)
	r.GET("/orders", h.ListOrders)
} 