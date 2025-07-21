package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterMetaproductRoutes(r *gin.RouterGroup, h *MetaproductHandlers) {
	r.POST("/metaproduct", h.CreateMetaproduct)
	r.POST("/metaproduct/batch", h.CreateMetaproducts)
	r.GET("/metaproduct", h.ListMetaproducts)
	r.GET("/metaproduct/:id", h.GetMetaproduct)
	r.PUT("/metaproduct/:id", h.UpdateMetaproduct)
	r.DELETE("/metaproduct/:id", h.DeleteMetaproduct)
} 