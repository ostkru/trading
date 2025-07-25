package offer

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		if err.Error() == "Доступ запрещён" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		log.Printf("UpdateOffer error: %v", err)
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
		if err.Error() == "Доступ запрещён" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if err.Error() == "Нельзя удалить оффер: есть связанные активные заказы" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("DeleteOffer error: %v", err)
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

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Параметр фильтрации
	filter := c.DefaultQuery("filter", "my")
	if filter != "my" && filter != "others" && filter != "all" {
		// Используем "my" по умолчанию для неверных значений
		filter = "my"
	}

	// Параметр типа оффера
	offerType := c.Query("offer_type")
	if offerType != "" && offerType != "sale" && offerType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.ListOffers(userID.(int64), page, limit, filter, offerType)
	if err != nil {
		log.Printf("ListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// ListOffersWithFilters обрабатывает запросы с расширенными фильтрами
func (h *Handlers) ListOffersWithFilters(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Парсим JSON с фильтрами
	var filters OfferFilterRequest
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters format"})
		return
	}

	// Валидация фильтров
	if filters.Filter != "" && filters.Filter != "my" && filters.Filter != "others" && filters.Filter != "all" {
		filters.Filter = "my" // По умолчанию
	}
	if filters.OfferType != "" && filters.OfferType != "sale" && filters.OfferType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.ListOffersWithFilters(userID.(int64), page, limit, &filters)
	if err != nil {
		log.Printf("ListOffersWithFilters error: %v", err)
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

func (h *Handlers) PublicListOffers(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	response, err := h.service.PublicListOffers(page, limit)
	if err != nil {
		log.Printf("PublicListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// PublicListOffersWithFilters обрабатывает запросы публичных офферов с фильтрами
func (h *Handlers) PublicListOffersWithFilters(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Парсим JSON с фильтрами
	var filters OfferFilterRequest
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters format"})
		return
	}

	// Валидация фильтров
	if filters.OfferType != "" && filters.OfferType != "sale" && filters.OfferType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.PublicListOffersWithFilters(page, limit, &filters)
	if err != nil {
		log.Printf("PublicListOffersWithFilters error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) CreateOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var req CreateOffersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	offers, err := h.service.CreateOffers(req, userID.(int64))
	if err != nil {
		log.Printf("CreateOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"offers": offers})
}
