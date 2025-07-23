package offer

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"log"
	"database/sql"
)
type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateOffer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	offer, err := h.service.CreateOffer(req, userID.(int64))
	if err != nil {
		log.Printf("CreateOffer error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

func (h *Handlers) UpdateOffer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offer, err := h.service.UpdateOffer(id, req, userID.(int64))
	if err != nil {
		log.Printf("UpdateOffer error: %v", err)
		if err.Error() == "Доступ запрещён" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offer)
}

func (h *Handlers) DeleteOffer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteOffer(id, userID.(int64)); err != nil {
		log.Printf("DeleteOffer error: %v", err)
		if err.Error() == "Нельзя удалить оффер: есть связанные активные заказы" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "Доступ запрещён" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handlers) ListOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	filter := c.DefaultQuery("filter", "my") // my, others, all
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	response, err := h.service.ListOffers(userID.(int64), page, limit, filter)
	if err != nil {
		log.Printf("ListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	offer, err := h.service.GetOfferByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Оффер не найден"})
			return
		}
		log.Printf("GetOffer error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offer)
}

func (h *Handlers) PublicListOffers(c *gin.Context) {
	offers, err := h.service.PublicListOffers()
	if err != nil {
		log.Printf("PublicListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func (h *Handlers) WBStock(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
		return
	}
	warehouseID, err := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse_id"})
		return
	}
	supplierID, err := strconv.ParseInt(c.Query("supplier_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier_id"})
		return
	}

	stock, err := h.service.WBStock(productID, warehouseID, supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *Handlers) CreateOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	
	var req CreateOffersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Проверяем количество офферов (максимум 100)
	if len(req.Offers) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Максимальное количество офферов за один запрос: 100"})
		return
	}
	
	offers, err := h.service.CreateOffers(req, userID.(int64))
	if err != nil {
		log.Printf("CreateOffers error: %v", err)
		if strings.Contains(err.Error(), "уже существует") || strings.Contains(err.Error(), "не может быть пустым") || strings.Contains(err.Error(), "Максимальное количество") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusCreated, offers)
} 