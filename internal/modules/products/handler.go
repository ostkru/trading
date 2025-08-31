package products

import (
	"fmt"
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

// @POST /products
// Создает новый продукт в системе
// Требует аутентификации и валидации данных
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

	// Валидация данных
	if req.Name == "" {
		response.BadRequest(c, "Требуется name")
		return
	}
	if req.VendorArticle == "" {
		response.BadRequest(c, "Требуется vendor_article")
		return
	}
	if req.RecommendPrice <= 0 {
		response.BadRequest(c, "Цена должна быть положительной")
		return
	}
	if req.Brand == "" {
		response.BadRequest(c, "Требуется brand")
		return
	}
	if req.Category == "" {
		response.BadRequest(c, "Требуется category")
		return
	}

	product, err := h.service.CreateProduct(&req, userID.(int64))
	if err != nil {
		// Проверяем ошибки валидации медиаданных
		if strings.Contains(err.Error(), "некорректный URL") {
			response.BadRequest(c, err.Error())
			return
		}
		log.Printf("CreateMetaproduct error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusCreated, product)
}

// @GET /products/{id}
// Получает продукт по его уникальному идентификатору
// Возвращает 404 если продукт не найден
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

// @GET /products
// Получает список всех продуктов с возможностью фильтрации
// Поддерживает пагинацию и сортировку
func (h *Handlers) ListMetaproducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	owner := c.DefaultQuery("owner", "my")
	uploadStatus := c.DefaultQuery("upload_status", "")

	// Валидация параметра owner
	validOwners := []string{"my", "all", "others", "pending"}
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

	// Валидация параметра upload_status
	validUploadStatuses := []string{"", "processing", "classified", "not_classified"}
	isValidUploadStatus := false
	for _, validStatus := range validUploadStatuses {
		if uploadStatus == validStatus {
			isValidUploadStatus = true
			break
		}
	}
	if !isValidUploadStatus {
		uploadStatus = "" // По умолчанию без фильтра по статусу
	}

	productsResponse, err := h.service.ListProducts(page, limit, owner, uploadStatus, userID.(int64))
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

	// Валидация данных
	if req.Name != nil && *req.Name == "" {
		response.BadRequest(c, "Имя не может быть пустым")
		return
	}
	if req.VendorArticle != nil && *req.VendorArticle == "" {
		response.BadRequest(c, "Артикул не может быть пустым")
		return
	}
	if req.RecommendPrice != nil && *req.RecommendPrice <= 0 {
		response.BadRequest(c, "Цена должна быть положительной")
		return
	}
	if req.Brand != nil && *req.Brand == "" {
		response.BadRequest(c, "Бренд не может быть пустым")
		return
	}
	if req.Category != nil && *req.Category == "" {
		response.BadRequest(c, "Категория не может быть пустым")
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

	// Валидация данных для каждого продукта
	for i, product := range req.Products {
		if product.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Продукт %d: требуется name", i+1)})
			return
		}
		if product.VendorArticle == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Продукт %d: требуется vendor_article", i+1)})
			return
		}
		if product.RecommendPrice <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Продукт %d: цена должна быть положительной", i+1)})
			return
		}
		if product.Brand == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Продукт %d: требуется brand", i+1)})
			return
		}
		if product.Category == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Продукт %d: требуется category", i+1)})
			return
		}
	}

	products, err := h.service.CreateProducts(req, userID.(int64))
	if err != nil {
		// Проверяем ошибки валидации медиаданных
		if strings.Contains(err.Error(), "некорректный URL") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, products)
}
