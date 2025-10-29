package characteristics

import (
	"time"
)

// Characteristic представляет характеристику продукта
type Characteristic struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Type        string    `json:"type"` // text, number, boolean, select, multiselect
	Required    bool      `json:"required"`
	Options     []string  `json:"options,omitempty"` // для select/multiselect
	Unit        string    `json:"unit,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CharacteristicValue представляет значение характеристики
type CharacteristicValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// CharacteristicStats представляет статистику по характеристике
type CharacteristicStats struct {
	Name        string                `json:"name"`
	TotalCount  int64                 `json:"total_count"`
	UniqueCount int64                 `json:"unique_count"`
	Values      []CharacteristicValue `json:"values"`
	Type        string                `json:"type"`
}

// ListCharacteristicsRequest запрос на получение списка характеристик
type ListCharacteristicsRequest struct {
	Page       int    `json:"page" form:"page"`
	Limit      int    `json:"limit" form:"limit"`
	Search     string `json:"search" form:"search"`
	CategoryID *int64 `json:"category_id" form:"category_id"`
	Type       string `json:"type" form:"type"`
	Sort       string `json:"sort" form:"sort"`
}

// ListCharacteristicsResponse ответ со списком характеристик
type ListCharacteristicsResponse struct {
	Characteristics []CharacteristicStats `json:"characteristics"`
	Total           int64                 `json:"total"`
	Page            int                   `json:"page"`
	Limit           int                   `json:"limit"`
}

// GetCharacteristicValuesRequest запрос на получение значений характеристики
type GetCharacteristicValuesRequest struct {
	Name        string `json:"name" form:"name"`
	CategoryID  *int64 `json:"category_id" form:"category_id"`
	Page        int    `json:"page" form:"page"`
	Limit       int    `json:"limit" form:"limit"`
	Search      string `json:"search" form:"search"`
}

// GetCharacteristicValuesResponse ответ со значениями характеристики
type GetCharacteristicValuesResponse struct {
	CharacteristicName string                `json:"characteristic_name"`
	Values             []CharacteristicValue `json:"values"`
	Total              int64                 `json:"total"`
	Page               int                   `json:"page"`
	Limit              int                   `json:"limit"`
}

// GetCategoryCharacteristicsRequest запрос на получение характеристик категории
type GetCategoryCharacteristicsRequest struct {
	CategoryName string `json:"category_name" form:"category_name"`
	Page         int    `json:"page" form:"page"`
	Limit        int    `json:"limit" form:"limit"`
}

// GetCategoryCharacteristicsResponse ответ с характеристиками категории
type GetCategoryCharacteristicsResponse struct {
	CategoryName     string                `json:"category_name"`
	Characteristics  []CharacteristicStats `json:"characteristics"`
	Total            int64                 `json:"total"`
	Page             int                   `json:"page"`
	Limit            int                   `json:"limit"`
}
