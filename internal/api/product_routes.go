package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, h *ProductHandlers) {
	r.POST("/meta", h.CreateProduct)
	r.POST("/meta/batch", h.CreateProducts)
	r.GET("/meta", h.ListProducts)
	r.GET("/meta/:id", h.GetProduct)
	r.PUT("/meta/:id", h.UpdateProduct)
	r.DELETE("/meta/:id", h.DeleteProduct)
} 