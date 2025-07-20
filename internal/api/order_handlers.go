package api

import (
	"net/http"
	"strconv"

	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderHandlers struct {
	service *services.OrderService
}

func NewOrderHandlers(service *services.OrderService) *OrderHandlers {
	return &OrderHandlers{service: service}
}

func (h *OrderHandlers) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заказа"})
		return
	}
	userID := int64(1)
	order, err := h.service.GetOrder(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandlers) ListOrders(c *gin.Context) {
	userID := int64(1)
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	orders, total, err := h.service.ListOrders(userID, statusPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders, "total": total, "page": page, "per_page": perPage})
} 