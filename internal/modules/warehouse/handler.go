package warehouse

import (
	"net/http"
	"strconv"

	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateWarehouse(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса: " + err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	warehouse, err := h.service.CreateWarehouse(req, userID.(int64))
	if err != nil {
		log.Printf("CreateWarehouse error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, warehouse)
}

func (h *Handlers) UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID склада"})
		return
	}
	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	warehouse, err := h.service.UpdateWarehouse(id, req, userID.(int64))
	if err != nil {
		if err.Error() == "Доступ запрещён" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		log.Printf("UpdateWarehouse error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

func (h *Handlers) DeleteWarehouse(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат ID"})
		return
	}

	if err := h.service.DeleteWarehouse(id, userID.(int64)); err != nil {
		if strings.Contains(err.Error(), "принадлежит другому пользователю") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted"})
}

func (h *Handlers) ListWarehouses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	warehouses, err := h.service.ListWarehouses(userID.(int64), page, limit)
	if err != nil {
		log.Printf("ListWarehouses error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouses)
}
