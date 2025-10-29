package categories

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"portaldata-api/internal/pkg/response"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// ListCategories @Summary Получение списка категорий
// @Description Получает список всех категорий с фильтрацией и пагинацией
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Param search query string false "Поиск по названию категории"
// @Param parent_id query int false "ID родительской категории"
// @Param active query bool false "Фильтр по активности"
// @Param sort query string false "Сортировка (name_asc, name_desc, products_asc, products_desc)" default(created_at)
// @Success 200 {object} response.SuccessResponse{data=ListCategoriesResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories [get]
func (h *Handlers) ListCategories(c *gin.Context) {
	var req ListCategoriesRequest
	
	// Парсим параметры запроса
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные параметры запроса: "+err.Error())
		return
	}

	// Получаем список категорий
	result, err := h.service.ListCategories(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения категорий: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// ListCategoryIDsAndNames @Summary Получить все category_id и category_name
// @Description Возвращает список всех уникальных пар id+name категорий из MySQL (таблица products)
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=ListCategoryIDNameResponse}
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories/list [get]
func (h *Handlers) ListCategoryIDsAndNames(c *gin.Context) {
    res, err := h.service.ListAllCategoryIDsAndNames(c.Request.Context())
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Ошибка получения категорий: "+err.Error())
        return
    }
    response.SuccessWithData(c, http.StatusOK, res)
}

// GetCategoryCharacteristics @Summary Получение характеристик категории
// @Description Получает список характеристик для указанной категории
// @Tags categories
// @Accept json
// @Produce json
// @Param category_name path string true "Название категории"
// @Success 200 {object} response.SuccessResponse{data=GetCategoryCharacteristicsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories/{category_name}/characteristics [get]
func (h *Handlers) GetCategoryCharacteristics(c *gin.Context) {
	categoryName := c.Param("category_name")
	if categoryName == "" {
		response.Error(c, http.StatusBadRequest, "Название категории обязательно")
		return
	}

	// Получаем характеристики категории
	result, err := h.service.GetCategoryCharacteristics(c.Request.Context(), categoryName)
	if err != nil {
		if err.Error() == "категория не найдена" {
			response.Error(c, http.StatusNotFound, "Категория не найдена")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Ошибка получения характеристик: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCategoryCharacteristicsByID @Summary Получение характеристик категории по ID
// @Description Получает список характеристик для указанной категории по ID из OpenSearch
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id path int true "ID категории"
// @Success 200 {object} response.SuccessResponse{data=GetCategoryCharacteristicsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories/by-id/:category_id/characteristics [get]
func (h *Handlers) GetCategoryCharacteristicsByID(c *gin.Context) {
	fmt.Printf("DEBUG: GetCategoryCharacteristicsByID handler called\n")
	
	categoryIDStr := c.Param("category_id")
	fmt.Printf("DEBUG: category_id param: %s\n", categoryIDStr)
	
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		fmt.Printf("DEBUG: Error parsing category_id: %v\n", err)
		response.Error(c, http.StatusBadRequest, "Неверный ID категории")
		return
	}

	fmt.Printf("DEBUG: Parsed categoryID: %d\n", categoryID)

	// Получаем характеристики категории по ID из OpenSearch
	result, err := h.service.GetCategoryCharacteristicsByID(c.Request.Context(), categoryID)
	if err != nil {
		fmt.Printf("DEBUG: Service error: %v\n", err)
		if err.Error() == "категория не найдена" {
			response.Error(c, http.StatusNotFound, "Категория не найдена")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Ошибка получения характеристик: "+err.Error())
		return
	}

	fmt.Printf("DEBUG: Service success, returning result\n")
	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCategoryStats @Summary Получение статистики категории
// @Description Получает статистику по категории (количество продуктов, офферов, цены)
// @Tags categories
// @Accept json
// @Produce json
// @Param category_name path string true "Название категории"
// @Success 200 {object} response.SuccessResponse{data=GetCategoryStatsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories/{category_name}/stats [get]
func (h *Handlers) GetCategoryStats(c *gin.Context) {
	categoryName := c.Param("category_name")
	if categoryName == "" {
		response.Error(c, http.StatusBadRequest, "Название категории обязательно")
		return
	}

	// Получаем статистику категории
	result, err := h.service.GetCategoryStats(c.Request.Context(), categoryName)
	if err != nil {
		if err.Error() == "категория не найдена" {
			response.Error(c, http.StatusNotFound, "Категория не найдена")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Ошибка получения статистики: "+err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, result)
}

// GetCategoryByID @Summary Получение категории по ID
// @Description Получает информацию о категории по её ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} response.SuccessResponse{data=Category}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/categories/{id} [get]
func (h *Handlers) GetCategoryByID(c *gin.Context) {
	_ = c.Param("id") // Пока не используем ID
	
	// Получаем категорию по ID (упрощенная версия)
	// В реальной системе здесь должен быть поиск по ID
	response.Error(c, http.StatusNotImplemented, "Поиск по ID пока не реализован")
}

// RegisterRoutes регистрирует маршруты для категорий
func RegisterRoutes(router *gin.RouterGroup, handlers *Handlers) {
	categoriesGroup := router.Group("/categories")
	{
		// Публичные маршруты
        categoriesGroup.GET("", handlers.ListCategories)
        categoriesGroup.GET("/list", handlers.ListCategoryIDsAndNames)
		categoriesGroup.GET("/by-name/:category_name/characteristics", handlers.GetCategoryCharacteristics)
		categoriesGroup.GET("/by-id/:category_id/characteristics", handlers.GetCategoryCharacteristicsByID)
		categoriesGroup.GET("/characteristics_by_id/:category_id", handlers.GetCategoryCharacteristicsByID) // Новый простой маршрут
		categoriesGroup.GET("/by-name/:category_name/stats", handlers.GetCategoryStats)
		categoriesGroup.GET("/info/:id", handlers.GetCategoryByID) // Изменили маршрут чтобы избежать конфликта
		
		// Тестовый маршрут
		categoriesGroup.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Test route works!"})
		})
	}
}
