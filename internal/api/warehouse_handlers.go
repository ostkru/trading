package api

import (
	"net/http"
	"strconv"

	"portaldata-api/internal/models"
	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

type WarehouseHandlers struct {
	service *services.WarehouseService
}

func NewWarehouseHandlers(service *services.WarehouseService) *WarehouseHandlers {
	return &WarehouseHandlers{service: service}
}

func (h *WarehouseHandlers) CreateWarehouse(c *gin.Context) {
	var req models.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса: " + err.Error()})
		return
	}
	userID := int64(1)
	warehouse, err := h.service.CreateWarehouse(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, warehouse)
}

func (h *WarehouseHandlers) UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID склада"})
		return
	}
	var req models.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := int64(5) // временно для теста
	warehouse, err := h.service.UpdateWarehouse(id, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

func (h *WarehouseHandlers) DeleteWarehouse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID склада"})
		return
	}
	userID := int64(5) // временно для теста
	if err := h.service.DeleteWarehouse(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *WarehouseHandlers) ListWarehouses(c *gin.Context) {
	userID := int64(5) // временно для теста, чтобы увидеть реальные склады
	warehouses, err := h.service.ListWarehouses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
} 