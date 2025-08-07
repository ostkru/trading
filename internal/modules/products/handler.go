package products

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"portaldata-api/internal/pkg/response"

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
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	product, err := h.service.CreateProduct(req, userID.(int64))
	if err != nil {
		if err.Error() == "Требуется name" {
			response.BadRequest(c, err.Error())
			return
		}
		log.Printf("CreateMetaproduct error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusCreated, product)
}

func (h *Handlers) GetMetaproduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный формат ID")
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		response.NotFound(c, "Продукт не найден")
		return
	}
	response.SuccessWithData(c, http.StatusOK, product)
}

func (h *Handlers) ListMetaproducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	owner := c.DefaultQuery("owner", "my")

	// Валидация параметра owner
	validOwners := []string{"my", "all", "others", "pending", "not_classified", "classified"}
	isValidOwner := false
	for _, validOwner := range validOwners {
		if owner == validOwner {
			isValidOwner = true
			break
		}
	}
	if !isValidOwner {
		owner = "my" // По умолчанию показываем только свои продукты
	}

	productsResponse, err := h.service.ListProducts(page, limit, owner, userID.(int64))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusOK, productsResponse)
}

func (h *Handlers) UpdateMetaproduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный формат ID")
		return
	}
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	updatedProduct, err := h.service.UpdateProduct(id, req, userID.(int64))
	if err != nil {
		if strings.Contains(err.Error(), "принадлежит другому пользователю") {
			response.Forbidden(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "не найден") {
			response.NotFound(c, err.Error())
			return
		}
		log.Printf("UpdateMetaproduct error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusOK, updatedProduct)
}

func (h *Handlers) DeleteMetaproduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный формат ID")
		return
	}
	if err := h.service.DeleteProduct(id, userID.(int64)); err != nil {
		if strings.Contains(err.Error(), "принадлежит другому пользователю") {
			response.Forbidden(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "не найден") {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, http.StatusOK, "Product deleted")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, products)
}
