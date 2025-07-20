package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterWarehouseRoutes(r *gin.RouterGroup, h *WarehouseHandlers) {
	r.POST("/warehouses", h.CreateWarehouse)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)
	r.GET("/warehouses", h.ListWarehouses)
} 