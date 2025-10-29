package characteristics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"portaldata-api/internal/pkg/response"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// ListCharacteristics @Summary Получение списка характеристик
// @Description Получает список всех характеристик с фильтрацией и пагинацией
// @Tags characteristics
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Param search query string false "Поиск по названию характеристики"
// @Param category_id query int false "ID категории"
// @Param type query string false "Тип характеристики"
// @Param sort query string false "Сортировка (name_asc, name_desc, count_asc, count_desc)" default(name_asc)
// @Success 200 {object} response.SuccessResponse{data=ListCharacteristicsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/characteristics [get]
func (h *Handlers) ListCharacteristics(c *gin.Context) {
	var req ListCharacteristicsRequest
	
	// Парсим параметры запроса
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные параметры запроса: "+err.Error())
		return
	}

	// Получаем список характеристик
	result, err := h.service.ListCharacteristics(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения характеристик: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCharacteristicValues @Summary Получение значений характеристики
// @Description Получает список всех значений для указанной характеристики
// @Tags characteristics
// @Accept json
// @Produce json
// @Param name query string true "Название характеристики"
// @Param category_id query int false "ID категории"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Param search query string false "Поиск по значению"
// @Success 200 {object} response.SuccessResponse{data=GetCharacteristicValuesResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/characteristics/values [get]
func (h *Handlers) GetCharacteristicValues(c *gin.Context) {
	var req GetCharacteristicValuesRequest
	
	// Парсим параметры запроса
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные параметры запроса: "+err.Error())
		return
	}

	if req.Name == "" {
		response.Error(c, http.StatusBadRequest, "Название характеристики обязательно")
		return
	}

	// Получаем значения характеристики
	result, err := h.service.GetCharacteristicValues(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения значений: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCategoryCharacteristics @Summary Получение характеристик категории
// @Description Получает список характеристик для указанной категории
// @Tags characteristics
// @Accept json
// @Produce json
// @Param category_name query string true "Название категории"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} response.SuccessResponse{data=GetCategoryCharacteristicsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/characteristics/category [get]
func (h *Handlers) GetCategoryCharacteristics(c *gin.Context) {
	var req GetCategoryCharacteristicsRequest
	
	// Парсим параметры запроса
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные параметры запроса: "+err.Error())
		return
	}

	if req.CategoryName == "" {
		response.Error(c, http.StatusBadRequest, "Название категории обязательно")
		return
	}

	// Получаем характеристики категории
	result, err := h.service.GetCategoryCharacteristics(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения характеристик категории: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCharacteristicStats @Summary Получение статистики характеристики
// @Description Получает статистику по указанной характеристике
// @Tags characteristics
// @Accept json
// @Produce json
// @Param name path string true "Название характеристики"
// @Success 200 {object} response.SuccessResponse{data=CharacteristicStats}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/characteristics/{name}/stats [get]
func (h *Handlers) GetCharacteristicStats(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, "Название характеристики обязательно")
		return
	}

	// Получаем статистику характеристики
	result, err := h.service.GetCharacteristicStats(c.Request.Context(), name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения статистики: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// RegisterRoutes регистрирует маршруты для характеристик
func RegisterRoutes(router *gin.RouterGroup, handlers *Handlers) {
	characteristicsGroup := router.Group("/characteristics")
	{
		// Публичные маршруты
		characteristicsGroup.GET("", handlers.ListCharacteristics)
		characteristicsGroup.GET("/values", handlers.GetCharacteristicValues)
		characteristicsGroup.GET("/category", handlers.GetCategoryCharacteristics)
		characteristicsGroup.GET("/:name/stats", handlers.GetCharacteristicStats)
	}
}
