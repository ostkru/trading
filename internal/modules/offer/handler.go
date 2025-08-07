package offer

import (
	"database/sql"
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

func (h *Handlers) CreateOffer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Некорректные данные запроса")
		return
	}

	offer, err := h.service.CreateOffer(req, userID.(int64))
	if err != nil {
		log.Printf("CreateOffer error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusCreated, offer)
}

func (h *Handlers) UpdateOffer(c *gin.Context) {
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
	var req UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	offer, err := h.service.UpdateOffer(id, req, userID.(int64))
	if err != nil {
		if err.Error() == "Доступ запрещён" {
			response.Forbidden(c, "Доступ запрещен")
			return
		}
		log.Printf("UpdateOffer error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusOK, offer)
}

func (h *Handlers) DeleteOffer(c *gin.Context) {
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

	if err := h.service.DeleteOffer(id, userID.(int64)); err != nil {
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

	response.SuccessWithMessage(c, http.StatusOK, "Offer deleted")
}

func (h *Handlers) ListOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Параметр фильтрации
	filter := c.DefaultQuery("filter", "my")
	if filter != "my" && filter != "others" && filter != "all" {
		// Используем "my" по умолчанию для неверных значений
		filter = "my"
	}

	// Параметр типа оффера
	offerType := c.Query("offer_type")
	if offerType != "" && offerType != "sale" && offerType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.ListOffers(userID.(int64), page, limit, filter, offerType)
	if err != nil {
		log.Printf("ListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// ListOffersWithFilters обрабатывает запросы с расширенными фильтрами
func (h *Handlers) ListOffersWithFilters(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Парсим JSON с фильтрами
	var filters OfferFilterRequest
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат фильтров"})
		return
	}

	// Валидация фильтров
	if filters.Filter != "" && filters.Filter != "my" && filters.Filter != "others" && filters.Filter != "all" {
		filters.Filter = "my" // По умолчанию
	}
	if filters.OfferType != "" && filters.OfferType != "sale" && filters.OfferType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.ListOffersWithFilters(userID.(int64), page, limit, &filters)
	if err != nil {
		log.Printf("ListOffersWithFilters error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) GetOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный формат ID")
		return
	}
	offer, err := h.service.GetOfferByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFound(c, "Оффер не найден")
			return
		}
		log.Printf("GetOffer error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}
	response.SuccessWithData(c, http.StatusOK, offer)
}

func (h *Handlers) WBStock(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID продукта")
		return
	}
	warehouseID, err := strconv.ParseInt(c.Query("warehouse_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID склада")
		return
	}
	supplierID, err := strconv.ParseInt(c.Query("supplier_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID поставщика")
		return
	}

	stock, err := h.service.WBStock(productID, warehouseID, supplierID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusOK, stock)
}

func (h *Handlers) PublicListOffers(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Создаем фильтры из query параметров
	filters := &OfferFilterRequest{}

	// Парсим offer_type
	if offerType := c.Query("offer_type"); offerType != "" {
		filters.OfferType = offerType
	}

	// Парсим ценовые фильтры
	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			filters.PriceMin = &priceMin
		}
	}
	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			filters.PriceMax = &priceMax
		}
	}

	// Парсим available_lots
	if availableLotsStr := c.Query("available_lots"); availableLotsStr != "" {
		if availableLots, err := strconv.Atoi(availableLotsStr); err == nil {
			filters.AvailableLots = &availableLots
		}
	}

	// Парсим brand_id
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		if brandID, err := strconv.ParseInt(brandIDStr, 10, 64); err == nil {
			filters.BrandID = &brandID
		}
	}

	// Парсим category_id
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	// Парсим product_name
	if productName := c.Query("product_name"); productName != "" {
		filters.ProductName = &productName
	}

	// Парсим vendor_article
	if vendorArticle := c.Query("vendor_article"); vendorArticle != "" {
		filters.VendorArticle = &vendorArticle
	}

	// Парсим warehouse_id
	if warehouseIDStr := c.Query("warehouse_id"); warehouseIDStr != "" {
		if warehouseID, err := strconv.ParseInt(warehouseIDStr, 10, 64); err == nil {
			filters.WarehouseID = &warehouseID
		}
	}

	// Парсим tax_nds
	if taxNDSStr := c.Query("tax_nds"); taxNDSStr != "" {
		if taxNDS, err := strconv.Atoi(taxNDSStr); err == nil {
			filters.TaxNDS = &taxNDS
		}
	}

	// Парсим units_per_lot
	if unitsPerLotStr := c.Query("units_per_lot"); unitsPerLotStr != "" {
		if unitsPerLot, err := strconv.Atoi(unitsPerLotStr); err == nil {
			filters.UnitsPerLot = &unitsPerLot
		}
	}

	// Парсим max_shipping_days
	if maxShippingDaysStr := c.Query("max_shipping_days"); maxShippingDaysStr != "" {
		if maxShippingDays, err := strconv.Atoi(maxShippingDaysStr); err == nil {
			filters.MaxShippingDays = &maxShippingDays
		}
	}

	// Парсим географические фильтры
	if minLatStr := c.Query("min_latitude"); minLatStr != "" {
		if minLat, err := strconv.ParseFloat(minLatStr, 64); err == nil {
			if maxLatStr := c.Query("max_latitude"); maxLatStr != "" {
				if maxLat, err := strconv.ParseFloat(maxLatStr, 64); err == nil {
					if minLngStr := c.Query("min_longitude"); minLngStr != "" {
						if minLng, err := strconv.ParseFloat(minLngStr, 64); err == nil {
							if maxLngStr := c.Query("max_longitude"); maxLngStr != "" {
								if maxLng, err := strconv.ParseFloat(maxLngStr, 64); err == nil {
									filters.Geographic = &GeographicFilter{
										MinLatitude:  minLat,
										MaxLatitude:  maxLat,
										MinLongitude: minLng,
										MaxLongitude: maxLng,
									}
								}
							}
						}
					}
				}
			}
		}
	}

	response, err := h.service.PublicListOffersWithFilters(page, limit, filters)
	if err != nil {
		log.Printf("PublicListOffers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// PublicListOffersWithFilters обрабатывает запросы публичных офферов с фильтрами
func (h *Handlers) PublicListOffersWithFilters(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Парсим JSON с фильтрами
	var filters OfferFilterRequest
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат фильтров"})
		return
	}

	// Валидация фильтров
	if filters.OfferType != "" && filters.OfferType != "sale" && filters.OfferType != "buy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое значение offer_type. Допустимые значения: sale, buy"})
		return
	}

	response, err := h.service.PublicListOffersWithFilters(page, limit, &filters)
	if err != nil {
		log.Printf("PublicListOffersWithFilters error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) CreateOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}

	var req CreateOffersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Некорректные данные запроса")
		return
	}

	offers, err := h.service.CreateOffers(req, userID.(int64))
	if err != nil {
		log.Printf("CreateOffers error: %v", err)
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithData(c, http.StatusCreated, gin.H{"offers": offers})
}
