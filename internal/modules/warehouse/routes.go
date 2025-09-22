package warehouse

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	r.POST("/warehouses", h.CreateWarehouse)
	r.POST("/warehouses/batch", h.CreateBatchWarehouses)
	r.GET("/warehouses/:id", h.GetWarehouse)
	r.PUT("/warehouses/:id", h.UpdateWarehouse)
	r.DELETE("/warehouses/:id", h.DeleteWarehouse)
	r.GET("/warehouses", h.ListWarehouses)
}
