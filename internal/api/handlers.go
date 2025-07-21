package api

import (
	"net/http"
	"strconv"

	"portaldata-api/internal/models"
	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	productService *services.ProductService
	parserService  *services.ParserService
	apiKey         string
}

func NewHandlers(productService *services.ProductService, parserService *services.ParserService, apiKey string) *Handlers {
	return &Handlers{
		productService: productService,
		parserService:  parserService,
		apiKey:        apiKey,
	}
}

// Middleware для проверки API ключа
func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("api_key")
		if apiKey == "" {
			apiKey = c.GetHeader("Authorization")
			if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
				apiKey = apiKey[7:]
			}
		}

		if apiKey != h.apiKey {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				OK:    false,
				Error: "Invalid API key",
				Code:  401,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Создание одного товара
func (h *Handlers) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid request data: " + err.Error(),
			Code:  400,
		})
		return
	}

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  500,
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		OK:   true,
		Data: product,
	})
}

// Создание нескольких товаров
func (h *Handlers) CreateProducts(c *gin.Context) {
	var req models.CreateProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid request data: " + err.Error(),
			Code:  400,
		})
		return
	}

	products, err := h.productService.CreateProducts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  500,
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		OK:   true,
		Data: products,
	})
}

// Получение товара по ID
func (h *Handlers) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid product ID",
			Code:  400,
		})
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  404,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		OK:   true,
		Data: product,
	})
}

// Список товаров с пагинацией
func (h *Handlers) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	owner := c.DefaultQuery("owner", "all")
	userID, _ := c.Get("user_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := h.productService.ListProducts(page, limit, owner, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  500,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		OK:   true,
		Data: response,
	})
}

// Обновление товара
func (h *Handlers) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid product ID",
			Code:  400,
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid request data: " + err.Error(),
			Code:  400,
		})
		return
	}

	product, err := h.productService.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  404,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		OK:   true,
		Data: product,
	})
}

// Удаление товара
func (h *Handlers) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Invalid product ID",
			Code:  400,
		})
		return
	}

	err = h.productService.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  404,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		OK: true,
	})
}

// Парсинг CSV файла
func (h *Handlers) ParseCSV(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			OK:    false,
			Error: "Path parameter is required",
			Code:  400,
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем список исключенных полей
	excludeStr := c.Query("exclude")
	var excludeFields []string
	if excludeStr != "" {
		excludeFields = []string{"Описание", "Наличие", "Валюта", "URL", "Изображения"}
	}

	response, err := h.parserService.ParseCSV(path, limit, offset, excludeFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			OK:    false,
			Error: err.Error(),
			Code:  500,
		})
		return
	}

	c.JSON(http.StatusOK, response)
} 