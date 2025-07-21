package metaproduct

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	r.POST("/metaproducts", h.CreateMetaproduct)
	r.GET("/metaproducts", h.ListMetaproducts)
	r.GET("/metaproducts/:id", h.GetMetaproduct)
	r.PUT("/metaproducts/:id", h.UpdateMetaproduct)
	r.DELETE("/metaproducts/:id", h.DeleteMetaproduct)
	r.POST("/metaproducts/batch", h.CreateMetaproducts)
} 