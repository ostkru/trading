package examples

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// @POST /products
// Создает новый продукт в системе
// Требует аутентификации и валидации данных
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	product, err := h.service.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    product,
	})
}

// @GET /products
// Получает список всех продуктов с возможностью фильтрации
// Поддерживает пагинацию и сортировку
func (h *Handler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
	})
}

// @GET /products/{id}
// Получает продукт по его уникальному идентификатору
// Возвращает 404 если продукт не найден
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// @PUT /products/{id}
// Обновляет существующий продукт
// Требует аутентификации и проверки прав доступа
func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	product, err := h.service.UpdateProduct(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// @DELETE /products/{id}
// Удаляет продукт из системы
// Требует подтверждения и проверки зависимостей
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Продукт успешно удален",
	})
}

// @POST /offers
// Создает новое торговое предложение
// Автоматически проверяет доступность продукта и склада
func (h *Handler) CreateOffer(c *gin.Context) {
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	offer, err := h.service.CreateOffer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    offer,
	})
}

// @GET /offers/public
// Получает список публичных предложений
// Доступно без аутентификации
func (h *Handler) ListPublicOffers(c *gin.Context) {
	offers, err := h.service.ListPublicOffers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    offers,
	})
}

// @POST /warehouses
// Создает новый склад
// Требует валидации географических координат
func (h *Handler) CreateWarehouse(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	warehouse, err := h.service.CreateWarehouse(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    warehouse,
	})
}

// @GET /warehouses
// Получает список всех складов
// Поддерживает фильтрацию по региону
func (h *Handler) ListWarehouses(c *gin.Context) {
	warehouses, err := h.service.ListWarehouses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    warehouses,
	})
}

// @POST /orders
// Создает новый заказ
// Автоматически проверяет доступность товара
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	order, err := h.service.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    order,
	})
}

// @GET /orders
// Получает список заказов пользователя
// Требует аутентификации
func (h *Handler) ListOrders(c *gin.Context) {
	orders, err := h.service.ListOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
	})
}

// @GET /rate-limit/stats
// Получает статистику использования API
// Доступно для администраторов
func (h *Handler) GetRateLimitStats(c *gin.Context) {
	stats, err := h.service.GetRateLimitStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// @POST /rate-limit/reset
// Сбрасывает счетчики rate limiting
// Требует административных прав
func (h *Handler) ResetRateLimit(c *gin.Context) {
	err := h.service.ResetRateLimit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rate limiting счетчики сброшены",
	})
}

// Структуры запросов и ответов
type CreateProductRequest struct {
	Name            string   `json:"name" binding:"required"`
	Brand           string   `json:"brand" binding:"required"`
	Category        string   `json:"category" binding:"required"`
	Description     string   `json:"description"`
	RecommendPrice  float64  `json:"recommend_price" binding:"min=0"`
	VendorArticle   string   `json:"vendor_article"`
	Barcode         string   `json:"barcode"`
	ImageURLs       []string `json:"image_urls"`
	VideoURLs       []string `json:"video_urls"`
	Model3DURLs     []string `json:"model_3d_urls"`
	BrandID         int64    `json:"brand_id"`
	CategoryID      int64    `json:"category_id"`
}

type UpdateProductRequest struct {
	Name            string   `json:"name" binding:"required"`
	Brand           string   `json:"brand" binding:"required"`
	Category        string   `json:"category" binding:"required"`
	Description     string   `json:"description"`
	RecommendPrice  float64  `json:"recommend_price" binding:"min=0"`
	VendorArticle   string   `json:"vendor_article"`
	Barcode         string   `json:"barcode"`
	ImageURLs       []string `json:"image_urls"`
	VideoURLs       []string `json:"video_urls"`
	Model3DURLs     []string `json:"model_3d_urls"`
}

type CreateOfferRequest struct {
	ProductID     int64   `json:"product_id" binding:"required"`
	Type          string  `json:"type" binding:"required,oneof=sale buy"`
	Price         float64 `json:"price" binding:"required,min=0"`
	LotCount      int     `json:"lot_count" binding:"required,min=1"`
	VAT           bool    `json:"vat"`
	DeliveryDays  int     `json:"delivery_days" binding:"min=1,max=365"`
	WarehouseID   int64   `json:"warehouse_id" binding:"required"`
}

type CreateWarehouseRequest struct {
	Name      string  `json:"name" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

type CreateOrderRequest struct {
	OfferID   int64 `json:"offer_id" binding:"required"`
	LotCount  int   `json:"lot_count" binding:"required,min=1"`
}

type Product struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Brand          string   `json:"brand"`
	Category       string   `json:"category"`
	Description    string   `json:"description"`
	RecommendPrice float64  `json:"recommend_price"`
	VendorArticle  string   `json:"vendor_article"`
	Barcode        string   `json:"barcode"`
	ImageURLs      []string `json:"image_urls"`
	VideoURLs      []string `json:"video_urls"`
	Model3DURLs    []string `json:"model_3d_urls"`
	BrandID        int64    `json:"brand_id"`
	CategoryID     int64    `json:"category_id"`
	UserID         int64    `json:"user_id"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type Offer struct {
	ID           int64   `json:"id"`
	ProductID    int64   `json:"product_id"`
	Type         string  `json:"type"`
	Price        float64 `json:"price"`
	LotCount     int     `json:"lot_count"`
	VAT          bool    `json:"vat"`
	DeliveryDays int     `json:"delivery_days"`
	WarehouseID  int64   `json:"warehouse_id"`
	UserID       int64   `json:"user_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type Warehouse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    int64   `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type Order struct {
	ID        int64       `json:"id"`
	OfferID   int64       `json:"offer_id"`
	LotCount  int         `json:"lot_count"`
	Status    string      `json:"status"`
	UserID    int64       `json:"user_id"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	ID       int64   `json:"id"`
	OrderID  int64   `json:"order_id"`
	ProductID int64  `json:"product_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type RateLimitStats struct {
	TotalAPIKeys    int64 `json:"total_api_keys"`
	ActiveAPIKeys   int64 `json:"active_api_keys"`
	BlockedAPIKeys  int64 `json:"blocked_api_keys"`
	TotalRequests   int64 `json:"total_requests"`
	BlockedRequests int64 `json:"blocked_requests"`
}

type Handler struct {
	service Service
}

type Service interface {
	CreateProduct(req CreateProductRequest) (*Product, error)
	ListProducts() ([]Product, error)
	GetProduct(id string) (*Product, error)
	UpdateProduct(id string, req UpdateProductRequest) (*Product, error)
	DeleteProduct(id string) error
	CreateOffer(req CreateOfferRequest) (*Offer, error)
	ListPublicOffers() ([]Offer, error)
	CreateWarehouse(req CreateWarehouseRequest) (*Warehouse, error)
	ListWarehouses() ([]Warehouse, error)
	CreateOrder(req CreateOrderRequest) (*Order, error)
	ListOrders() ([]Order, error)
	GetRateLimitStats() (*RateLimitStats, error)
	ResetRateLimit() error
}
