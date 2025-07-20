package api

import (
	"net/http"
	"strconv"

	"portaldata-api/internal/models"
	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

type OfferHandlers struct {
	service *services.OfferService
}

func NewOfferHandlers(service *services.OfferService) *OfferHandlers {
	return &OfferHandlers{service: service}
}

func (h *OfferHandlers) CreateOffer(c *gin.Context) {
	var req models.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса: " + err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	offer, err := h.service.CreateOffer(req, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, offer)
}

func (h *OfferHandlers) UpdateOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID оффера"})
		return
	}
	var req models.UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	offer, err := h.service.UpdateOffer(id, req, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offer)
}

func (h *OfferHandlers) DeleteOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID оффера"})
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.DeleteOffer(id, userID.(int64))
	if err != nil {
		if err.Error() == "Нельзя удалить оффер: есть связанные активные заказы" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *OfferHandlers) ListOffers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	offers, err := h.service.ListOffers(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func (h *OfferHandlers) PublicListOffers(c *gin.Context) {
	offers, err := h.service.PublicListOffers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func (h *OfferHandlers) WBStock(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	warehouseID, _ := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)
	supplierID, _ := strconv.ParseInt(c.Query("supplier_id"), 10, 64)
	if productID == 0 || warehouseID == 0 || supplierID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id, warehouse_id, supplier_id required"})
		return
	}
	stock, err := h.service.WBStock(productID, warehouseID, supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stock": stock})
} 