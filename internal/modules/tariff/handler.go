package tariff

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTariff создает новый тариф
// @POST /tariffs
func (h *Handler) CreateTariff(c *gin.Context) {
	var req CreateTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные: " + err.Error(),
		})
		return
	}

	tariff, err := h.service.CreateTariff(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tariff,
	})
}

// GetTariff получает тариф по ID
// @GET /tariffs/{id}
func (h *Handler) GetTariff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID тарифа",
		})
		return
	}

	tariff, err := h.service.GetTariff(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tariff,
	})
}

// ListTariffs получает список тарифов
// @GET /tariffs
func (h *Handler) ListTariffs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	activeOnly := c.DefaultQuery("active_only", "false") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := h.service.ListTariffs(page, limit, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateTariff обновляет тариф
// @PUT /tariffs/{id}
func (h *Handler) UpdateTariff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID тарифа",
		})
		return
	}

	var req UpdateTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные: " + err.Error(),
		})
		return
	}

	tariff, err := h.service.UpdateTariff(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tariff,
	})
}

// DeleteTariff удаляет тариф
// @DELETE /tariffs/{id}
func (h *Handler) DeleteTariff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID тарифа",
		})
		return
	}

	err = h.service.DeleteTariff(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Тариф успешно удален",
	})
}

// GetUserTariffInfo получает информацию о тарифе пользователя
// @GET /tariffs/user/{user_id}
func (h *Handler) GetUserTariffInfo(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректный ID пользователя",
		})
		return
	}

	info, err := h.service.GetUserTariffInfo(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// ChangeUserTariff изменяет тариф пользователя
// @PUT /tariffs/user/change
func (h *Handler) ChangeUserTariff(c *gin.Context) {
	var req ChangeUserTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные: " + err.Error(),
		})
		return
	}

	err := h.service.ChangeUserTariff(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Тариф пользователя успешно изменен",
	})
}

// GetTariffUsageStats получает статистику использования тарифов
// @GET /tariffs/stats
func (h *Handler) GetTariffUsageStats(c *gin.Context) {
	stats, err := h.service.GetTariffUsageStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetMyTariffInfo получает информацию о тарифе текущего пользователя
// @GET /tariffs/my
func (h *Handler) GetMyTariffInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Пользователь не авторизован",
		})
		return
	}

	info, err := h.service.GetUserTariffInfo(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}
