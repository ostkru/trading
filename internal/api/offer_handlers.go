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
	newOffer, err := h.service.CreateOffer(userID.(int64), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newOffer)
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Ограничение максимального лимита
	}
	response, err := h.service.ListOffers(userID.(int64), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *OfferHandlers) GetOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID оффера"})
		return
	}

	offer, err := h.service.GetOfferByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Оффер не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверка доступа (опционально, если нужно)
	// userID, _ := c.Get("user_id")
	// if offer.UserID != userID.(int64) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
	// 	return
	// }

	c.JSON(http.StatusOK, offer)
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