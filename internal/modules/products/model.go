package products

import (
	"time"
	"portaldata-api/internal/utils"
)

type Product struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	VendorArticle  string    `json:"vendor_article"`
	RecommendPrice float64   `json:"recommend_price"`
	Brand          *string   `json:"brand,omitempty"`
	Category       *string   `json:"category,omitempty"`
	BrandID        *int64    `json:"brand_id,omitempty"`
	CategoryID     *int64    `json:"category_id,omitempty"`
	Description    string    `json:"description"`
	Barcode        *string   `json:"barcode,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         int64     `json:"user_id"`
	// Медиа поля (могут быть JSON строками или слайсами)
	ImageURLs   interface{} `json:"image_urls,omitempty"`
	VideoURLs   interface{} `json:"video_urls,omitempty"`
	Model3DURLs interface{} `json:"model_3d_urls,omitempty"`
}

type CreateProductRequest struct {
	Name           string  `json:"name"`
	VendorArticle  string  `json:"vendor_article"`
	RecommendPrice float64 `json:"recommend_price"`
	Brand          string  `json:"brand"`
	Category       string  `json:"category"`
	BrandID        *int64  `json:"brand_id,omitempty"`
	CategoryID     *int64  `json:"category_id,omitempty"`
	Description    string  `json:"description"`
	Barcode        *string `json:"barcode,omitempty"`
	// Медиа поля
	ImageURLs   []string `json:"image_urls,omitempty"`
	VideoURLs   []string `json:"video_urls,omitempty"`
	Model3DURLs []string `json:"model_3d_urls,omitempty"`
}

// GenerateCategoryID автоматически генерирует category_id из названия категории
// Работает только для WB категорий в формате "wb: 1318 - сварочные аппараты"
// если category_id не был предоставлен
func (r *CreateProductRequest) GenerateCategoryID() {
	if r.CategoryID == nil && r.Category != "" {
		categoryID := utils.GenerateCategoryID(r.Category)
		// Генерируем ID только если это WB категория (categoryID > 0)
		if categoryID > 0 {
			r.CategoryID = &categoryID
		}
	}
}

type UpdateProductRequest struct {
	Name           *string  `json:"name,omitempty"`
	VendorArticle  *string  `json:"vendor_article,omitempty"`
	RecommendPrice *float64 `json:"recommend_price,omitempty"`
	Brand          *string  `json:"brand,omitempty"`
	Category       *string  `json:"category,omitempty"`
	BrandID        *int64   `json:"brand_id,omitempty"`
	CategoryID     *int64   `json:"category_id,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Barcode        *string  `json:"barcode,omitempty"`
	// Медиа поля
	ImageURLs   *[]string `json:"image_urls,omitempty"`
	VideoURLs   *[]string `json:"video_urls,omitempty"`
	Model3DURLs *[]string `json:"model_3d_urls,omitempty"`
}

// GenerateCategoryID автоматически генерирует category_id из названия категории
// Работает только для WB категорий в формате "wb: 1318 - сварочные аппараты"
// если category_id не был предоставлен, но category был изменен
func (r *UpdateProductRequest) GenerateCategoryID() {
	if r.CategoryID == nil && r.Category != nil && *r.Category != "" {
		categoryID := utils.GenerateCategoryID(*r.Category)
		// Генерируем ID только если это WB категория (categoryID > 0)
		if categoryID > 0 {
			r.CategoryID = &categoryID
		}
	}
}

type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

type CreateProductsRequest struct {
	Products []CreateProductRequest `json:"products" validate:"required,dive"`
}
