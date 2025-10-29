package categories

import (
	"time"
)

// Category представляет категорию продукта
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProductCount int64    `json:"product_count,omitempty"`
}

// Characteristic представляет характеристику категории
type Characteristic struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // text, number, boolean, select, multiselect
	Required    bool      `json:"required"`
	Options     []string  `json:"options,omitempty"` // для select/multiselect
	Unit        string    `json:"unit,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryWithCharacteristics представляет категорию с характеристиками
type CategoryWithCharacteristics struct {
	Category       Category        `json:"category"`
	Characteristics []Characteristic `json:"characteristics"`
}

// CategoryStats представляет статистику по категории
type CategoryStats struct {
	CategoryID    int64 `json:"category_id"`
	CategoryName  string `json:"category_name"`
	ProductCount  int64 `json:"product_count"`
	OfferCount    int64 `json:"offer_count"`
	AvgPrice      float64 `json:"avg_price,omitempty"`
	MinPrice      float64 `json:"min_price,omitempty"`
	MaxPrice      float64 `json:"max_price,omitempty"`
}

// ListCategoriesRequest запрос на получение списка категорий
type ListCategoriesRequest struct {
	Page     int    `json:"page" form:"page"`
	Limit    int    `json:"limit" form:"limit"`
	Search   string `json:"search" form:"search"`
	ParentID *int64 `json:"parent_id" form:"parent_id"`
	Active   *bool  `json:"active" form:"active"`
	Sort     string `json:"sort" form:"sort"`
}

// ListCategoriesResponse ответ со списком категорий
type ListCategoriesResponse struct {
	Categories []Category `json:"categories"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}

// CategoryIDName представляет пару id+name для категории
type CategoryIDName struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// ListCategoryIDNameResponse ответ со списком пар id+name
type ListCategoryIDNameResponse struct {
    Categories []CategoryIDName `json:"categories"`
    Total      int64            `json:"total"`
}

// GetCategoryCharacteristicsRequest запрос на получение характеристик категории
type GetCategoryCharacteristicsRequest struct {
	CategoryID int64 `json:"category_id" form:"category_id"`
}

// GetCategoryCharacteristicsResponse ответ с характеристиками категории
type GetCategoryCharacteristicsResponse struct {
	Category        Category        `json:"category"`
	Characteristics []Characteristic `json:"characteristics"`
}

// GetCategoryStatsRequest запрос на получение статистики категории
type GetCategoryStatsRequest struct {
	CategoryID int64 `json:"category_id" form:"category_id"`
}

// GetCategoryStatsResponse ответ со статистикой категории
type GetCategoryStatsResponse struct {
	Stats CategoryStats `json:"stats"`
}
