package search

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"portaldata-api/internal/pkg/response"
)

// Handlers обработчики для поиска
type Handlers struct {
	service *Service
}

// NewHandlers создает новые обработчики
func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// SearchProducts @Summary Поиск продуктов
// @Description Выполняет полнотекстовый поиск продуктов по различным критериям
// @Tags search
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param category_id query int false "ID категории"
// @Param brand_id query int false "ID бренда"
// @Param brand query string false "Название бренда"
// @Param category query string false "Название категории"
// @Param price_min query number false "Минимальная цена"
// @Param price_max query number false "Максимальная цена"
// @Param page query int false "Страница (начиная с 1)" default(1)
// @Param limit query int false "Количество результатов" default(20)
// @Param sort query string false "Сортировка (relevance, price_asc, price_desc, name_asc, name_desc)" default(relevance)
// @Param facets query bool false "Включить агрегации для фильтров" default(false)
// @Success 200 {object} response.Response{data=SearchResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/search/products [get]
func (h *Handlers) SearchProducts(c *gin.Context) {
	// Парсинг параметров запроса
	req := &SearchRequest{
		Page:  1,
		Limit: 20,
		Sort:  "relevance",
	}
	
	// Текстовый поиск
	if query := c.Query("query"); query != "" {
		req.Query = query
	}
	
	// Фильтры
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			req.CategoryID = &categoryID
		}
	}
	
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		if brandID, err := strconv.ParseInt(brandIDStr, 10, 64); err == nil {
			req.BrandID = &brandID
		}
	}
	
	if brand := c.Query("brand"); brand != "" {
		req.Brand = &brand
	}
	
	if category := c.Query("category"); category != "" {
		req.Category = &category
	}
	
	// Ценовые фильтры
	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			req.PriceMin = &priceMin
		}
	}
	
	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			req.PriceMax = &priceMax
		}
	}
	
	// Пагинация
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			req.Limit = limit
		}
	}
	
	// Сортировка
	if sort := c.Query("sort"); sort != "" {
		req.Sort = sort
	}
	
	// Агрегации
	if facetsStr := c.Query("facets"); facetsStr != "" {
		if facets, err := strconv.ParseBool(facetsStr); err == nil {
			req.Facets = facets
		}
	}
	
	// Выполнить поиск
	result, err := h.service.SearchProducts(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка поиска: "+err.Error())
		return
	}
	
	response.SuccessWithData(c, http.StatusOK, result)
}

// SearchProductsByCharacteristics @Summary Поиск продуктов по характеристикам
// @Description Выполняет поиск продуктов по техническим характеристикам
// @Tags search
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Параметры поиска"
// @Success 200 {object} response.Response{data=SearchResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/search/products/characteristics [post]
func (h *Handlers) SearchProductsByCharacteristics(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные данные: "+err.Error())
		return
	}
	
	// Установить значения по умолчанию
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Sort == "" {
		req.Sort = "relevance"
	}
	
	// Выполнить поиск
	result, err := h.service.SearchProducts(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка поиска: "+err.Error())
		return
	}
	
	response.SuccessWithData(c, http.StatusOK, result)
}

// IndexProduct @Summary Индексация продукта
// @Description Добавляет или обновляет продукт в поисковом индексе
// @Tags search
// @Accept json
// @Produce json
// @Param request body IndexProductRequest true "Данные продукта"
// @Success 200 {object} response.Response{data=string}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/search/index [post]
func (h *Handlers) IndexProduct(c *gin.Context) {
	var req IndexProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Неверные данные: "+err.Error())
		return
	}
	
	// Индексировать продукт
	if err := h.service.IndexProduct(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка индексации: "+err.Error())
		return
	}
	
	response.SuccessWithMessage(c, http.StatusOK, "Продукт успешно проиндексирован")
}

// DeleteProduct @Summary Удаление продукта из индекса
// @Description Удаляет продукт из поискового индекса
// @Tags search
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} response.Response{data=string}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/search/index/{id} [delete]
func (h *Handlers) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Неверный ID продукта: "+err.Error())
		return
	}
	
	// Удалить продукт из индекса
	if err := h.service.DeleteProduct(c.Request.Context(), productID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка удаления: "+err.Error())
		return
	}
	
	response.SuccessWithMessage(c, http.StatusOK, "Продукт успешно удален из индекса")
}

// GetSearchStats @Summary Статистика поиска
// @Description Возвращает статистику поискового индекса
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 500 {object} response.Response
// @Router /api/search/stats [get]
func (h *Handlers) GetSearchStats(c *gin.Context) {
	// Получить статистику индекса
	stats, err := h.service.GetIndexStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Ошибка получения статистики: "+err.Error())
		return
	}
	
	response.SuccessWithData(c, http.StatusOK, stats)
}

// RegisterRoutes регистрирует маршруты поиска
func RegisterRoutes(router *gin.RouterGroup, handlers *Handlers) {
	searchGroup := router.Group("/search")
	{
		// Публичные маршруты (без авторизации)
		searchGroup.GET("/products", handlers.SearchProducts)
		searchGroup.POST("/products/characteristics", handlers.SearchProductsByCharacteristics)
		
		// Защищенные маршруты (требуют авторизации)
		searchGroup.POST("/index", handlers.IndexProduct)
		searchGroup.DELETE("/index/:id", handlers.DeleteProduct)
		searchGroup.GET("/stats", handlers.GetSearchStats)
	}
}
