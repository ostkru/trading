package metaproduct

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)
type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateMetaproduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Валидация обязательных полей
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Название продукта обязательно"})
		return
	}
	product, err := h.service.CreateProduct(req, userID.(int64))
	if err != nil {
		log.Printf("CreateMetaproduct error: %v", err)
		// Проверяем, является ли ошибка ошибкой валидации
		if strings.Contains(err.Error(), "уже существует") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handlers) GetMetaproduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handlers) ListMetaproducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	owner := c.DefaultQuery("owner", "all")
	response, err := h.service.ListProducts(page, limit, owner, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) UpdateMetaproduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
	var req UpdateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	updatedProduct, err := h.service.UpdateProduct(id, req, userID.(int64))
    if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "продукт не найден"})
			return
		}
		if strings.Contains(err.Error(), "недостаточно прав") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "уже существует") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("UpdateMetaproduct error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusOK, updatedProduct)
}

func (h *Handlers) DeleteMetaproduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
	if err := h.service.DeleteProduct(id, userID.(int64)); err != nil {
		if strings.Contains(err.Error(), "недостаточно прав") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *Handlers) CreateMetaproducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	var req CreateProductsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	products, err := h.service.CreateProducts(req, userID.(int64))
    if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
        return
    }
    c.JSON(http.StatusCreated, products)
}
