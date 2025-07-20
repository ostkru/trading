package api

import (
	"net/http"
	"strconv"

	"portaldata-api/internal/models"
	"portaldata-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ProductHandlers struct {
	service *services.ProductService
}

func NewProductHandlers(service *services.ProductService) *ProductHandlers {
	return &ProductHandlers{service: service}
}

func (h *ProductHandlers) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса: " + err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	_ = userID // если нужно использовать user_id при создании
	product, err := h.service.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandlers) CreateProducts(c *gin.Context) {
	var req models.CreateProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	products, err := h.service.CreateProducts(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, products)
}

func (h *ProductHandlers) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID мета-объекта"})
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Мета-объект не найден"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandlers) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	response, err := h.service.ListProducts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandlers) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID мета-объекта"})
		return
	}
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandlers) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID мета-объекта"})
		return
	}
	err = h.service.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
} 