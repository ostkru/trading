package metaproduct

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	r.POST("/products", h.CreateMetaproduct)
	r.GET("/products", h.ListMetaproducts)
	r.GET("/products/:id", h.GetMetaproduct)
	r.PUT("/products/:id", h.UpdateMetaproduct)
	r.DELETE("/products/:id", h.DeleteMetaproduct)
	r.POST("/products/batch", h.CreateMetaproducts)
} 