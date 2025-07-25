package metaproduct

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	// Основные маршруты metaproducts
	r.POST("/metaproducts", h.CreateMetaproduct)
	r.GET("/metaproducts", h.ListMetaproducts)
	r.GET("/metaproducts/:id", h.GetMetaproduct)
	r.PUT("/metaproducts/:id", h.UpdateMetaproduct)
	r.DELETE("/metaproducts/:id", h.DeleteMetaproduct)
	r.POST("/metaproducts/batch", h.CreateMetaproducts)

	// Алиасы для совместимости с тестами
	r.POST("/products", h.CreateMetaproduct)
	r.GET("/products", h.ListMetaproducts)
	r.GET("/products/:id", h.GetMetaproduct)
	r.PUT("/products/:id", h.UpdateMetaproduct)
	r.DELETE("/products/:id", h.DeleteMetaproduct)
	r.POST("/products/batch", h.CreateMetaproducts)
}
