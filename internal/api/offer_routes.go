package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterOfferRoutes(r *gin.RouterGroup, h *OfferHandlers) {
	r.POST("/offers", h.CreateOffer)
	r.PUT("/offers/:id", h.UpdateOffer)
	r.DELETE("/offers/:id", h.DeleteOffer)
	r.GET("/offers", h.ListOffers)
	r.GET("/offers/:id", h.GetOffer)
	r.GET("/offers/public", h.PublicListOffers)
	r.GET("/offers/wb_stock", h.WBStock)
} 