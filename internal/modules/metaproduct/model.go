package metaproduct

import (
	"encoding/json"
	"time"
)

type Media struct {
	ID          int64           `json:"id"`
	ProductID   int64           `json:"product_id"`
	ImageURLs   json.RawMessage `json:"image_urls,omitempty"`
	VideoURLs   json.RawMessage `json:"video_urls,omitempty"`
	Model3DURLs json.RawMessage `json:"model_3d_urls,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Product struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	VendorArticle  string    `json:"vendor_article"`
	RecommendPrice float64   `json:"recommend_price"`
	Brand          string    `json:"brand"`
	Category       string    `json:"category"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         int64     `json:"user_id"`
	CategoryID     *int64    `json:"category_id,omitempty"`
	BrandID        *int64    `json:"brand_id,omitempty"`
	Media          *Media    `json:"media,omitempty"`
}

type CreateProductRequest struct {
	Name           string          `json:"name"`
	VendorArticle  string          `json:"vendor_article"`
	RecommendPrice float64         `json:"recommend_price"`
	Brand          string          `json:"brand"`
	Category       string          `json:"category"`
	Description    string          `json:"description"`
	ImageURLs      json.RawMessage `json:"image_urls,omitempty"`      // обязательные ссылки на изображения
	VideoURLs      json.RawMessage `json:"video_urls,omitempty"`      // ссылки на видео обзоры
	Model3DURLs    json.RawMessage `json:"model_3d_urls,omitempty"`   // ссылки на 3д модели
}

type UpdateProductRequest struct {
	Name           *string          `json:"name,omitempty"`
	VendorArticle  *string          `json:"vendor_article,omitempty"`
	RecommendPrice *float64         `json:"recommend_price,omitempty"`
	Brand          *string          `json:"brand,omitempty"`
	Category       *string          `json:"category,omitempty"`
	Description    *string          `json:"description,omitempty"`
	ImageURLs      *json.RawMessage `json:"image_urls,omitempty"`      // ссылки на изображения
	VideoURLs      *json.RawMessage `json:"video_urls,omitempty"`      // ссылки на видео обзоры
	Model3DURLs    *json.RawMessage `json:"model_3d_urls,omitempty"`   // ссылки на 3д модели
}

type ProductListResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

type CreateProductsRequest struct {
	Products []CreateProductRequest `json:"products" validate:"required,dive"`
} 