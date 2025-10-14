package search

import (
	"time"
)

// SearchRequest запрос на поиск продуктов
type SearchRequest struct {
	Query         string                 `json:"query"`                    // Текстовый поиск
	CategoryID    *int64                 `json:"category_id,omitempty"`   // Фильтр по категории
	BrandID       *int64                 `json:"brand_id,omitempty"`     // Фильтр по бренду
	Brand         *string                `json:"brand,omitempty"`         // Фильтр по названию бренда
	Category      *string                `json:"category,omitempty"`      // Фильтр по названию категории
	Characteristics map[string]string     `json:"characteristics,omitempty"` // Фильтр по характеристикам
	PriceMin      *float64               `json:"price_min,omitempty"`     // Минимальная цена
	PriceMax      *float64               `json:"price_max,omitempty"`      // Максимальная цена
	Page          int                    `json:"page"`                     // Страница (начиная с 1)
	Limit         int                    `json:"limit"`                   // Количество результатов
	Sort          string                 `json:"sort,omitempty"`          // Сортировка: relevance, price_asc, price_desc, name_asc, name_desc
	Facets        bool                   `json:"facets,omitempty"`        // Включить агрегации для фильтров
}

// SearchResponse ответ на поисковый запрос
type SearchResponse struct {
	Products []SearchProduct `json:"products"`     // Найденные продукты
	Total    int64           `json:"total"`        // Общее количество результатов
	Page     int             `json:"page"`         // Текущая страница
	Limit    int             `json:"limit"`        // Количество на странице
	Facets   *SearchFacets   `json:"facets,omitempty"` // Агрегации для фильтров
}

// SearchProduct продукт в результатах поиска
type SearchProduct struct {
	ID             string                 `json:"id"`              // ID продукта
	Name           string                 `json:"name"`             // Название
	VendorCode     string                 `json:"vendor_code"`      // Артикул
	Brand          string                 `json:"brand"`            // Бренд
	BrandID        int64                  `json:"brand_id"`        // ID бренда
	Category       string                 `json:"category"`         // Категория
	CategoryID     string                 `json:"category_id"`      // ID категории
	Characteristics []Characteristic      `json:"characteristics"` // Характеристики
	ImageURLs      []string               `json:"image_urls"`      // URL изображений
	Score          float64                `json:"score,omitempty"` // Релевантность
}

// Characteristic характеристика продукта
type Characteristic struct {
	Name  string `json:"name"`  // Название характеристики
	Value string `json:"value"` // Значение характеристики
}

// SearchFacets агрегации для фильтров
type SearchFacets struct {
	Brands      []FacetItem `json:"brands"`       // Список брендов
	Categories  []FacetItem `json:"categories"`   // Список категорий
	PriceRanges []PriceRange `json:"price_ranges"` // Ценовые диапазоны
}

// FacetItem элемент агрегации
type FacetItem struct {
	Value string `json:"value"` // Значение
	Count int64  `json:"count"` // Количество
}

// PriceRange ценовой диапазон
type PriceRange struct {
	Min   float64 `json:"min"`   // Минимальная цена
	Max   float64 `json:"max"`   // Максимальная цена
	Count int64   `json:"count"` // Количество в диапазоне
}

// IndexProductRequest запрос на индексацию продукта
type IndexProductRequest struct {
	ID             int64                  `json:"id"`
	Name           string                 `json:"name"`
	VendorArticle  string                 `json:"vendor_article"`
	RecommendPrice float64                `json:"recommend_price"`
	Brand          *string                `json:"brand,omitempty"`
	Category       *string                `json:"category,omitempty"`
	BrandID        *int64                 `json:"brand_id,omitempty"`
	CategoryID     *int64                 `json:"category_id,omitempty"`
	Description    string                 `json:"description"`
	Barcode        *string                `json:"barcode,omitempty"`
	ImageURLs      []string               `json:"image_urls,omitempty"`
	VideoURLs      []string               `json:"video_urls,omitempty"`
	Model3DURLs    []string               `json:"model_3d_urls,omitempty"`
	Characteristics []Characteristic      `json:"characteristics,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	UserID         int64                  `json:"user_id"`
}
