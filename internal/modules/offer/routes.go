package offer

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	r.POST("/offers", h.CreateOffer)
	r.PUT("/offers/:id", h.UpdateOffer)
	r.DELETE("/offers/:id", h.DeleteOffer)
	r.GET("/offers", h.ListOffers)
	r.POST("/offers/filter", h.ListOffersWithFilters) // Новый маршрут для фильтрованных запросов
	r.GET("/offers/:id", h.GetOffer)
	r.GET("/offers/wb_stock", h.WBStock)
	r.POST("/offers/batch", h.CreateOffers)
}

// RegisterPublicRoutes регистрирует публичные маршруты (без авторизации)
func RegisterPublicRoutes(r *gin.RouterGroup, h *Handlers) {
	r.GET("/offers/public", h.PublicListOffers)
	r.POST("/offers/public/filter", h.PublicListOffersWithFilters) // Новый маршрут для фильтрованных публичных запросов
}
