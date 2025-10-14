package order

import (
	"log"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
		return
	}

	order, err := h.service.CreateOrder(userID.(int64), req)
	if err != nil {
		if err.Error() == "offer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено"})
			return
		}
		if err.Error() == "Нельзя создать заказ на собственное предложение" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error()[:len("not enough lots available")] == "not enough lots available" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("CreateOrder error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handlers) GetOrder(c *gin.Context) {
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
	order, err := h.service.GetOrderByID(id, userID.(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "заказ не найден"})
			return
		}
		log.Printf("GetOrder error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении заказа"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *Handlers) ListOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	role := c.DefaultQuery("role", "all")
	orders, total, err := h.service.ListOrders(userID.(int64), &status, role, page, perPage)
	if err != nil {
		log.Printf("ListOrders error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении заказов"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"orders":   orders,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *Handlers) UpdateOrderStatus(c *gin.Context) {
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

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
		return
	}

    order, err := h.service.UpdateOrderStatus(id, userID.(int64), req.Status)
	if err != nil {
        switch err.Error() {
        case "Заказ не найден", "Order not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
			return
        case "Доступ запрещен", "Access denied":
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
			return
        case "Invalid status transition for role":
            c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый переход статуса для вашей роли"})
            return
		}
		log.Printf("UpdateOrderStatus error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
