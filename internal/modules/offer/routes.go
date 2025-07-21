package offer

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handlers) {
	r.POST("/offers", h.CreateOffer)
	r.PUT("/offers/:id", h.UpdateOffer)
	r.DELETE("/offers/:id", h.DeleteOffer)
	r.GET("/offers", h.ListOffers)
	r.GET("/offers/:id", h.GetOffer)
	r.GET("/offers/wb_stock", h.WBStock)
}

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handlers) {
	r.GET("/offers/public", h.PublicListOffers)
} 