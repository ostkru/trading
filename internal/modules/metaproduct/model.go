package metaproduct

import "time"

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
}

type CreateProductRequest struct {
	Name           string  `json:"name"`
	VendorArticle  string  `json:"vendor_article"`
	RecommendPrice float64 `json:"recommend_price"`
	Brand          string  `json:"brand"`
	Category       string  `json:"category"`
	Description    string  `json:"description"`
}

type UpdateProductRequest struct {
	Name           *string  `json:"name,omitempty"`
	VendorArticle  *string  `json:"vendor_article,omitempty"`
	RecommendPrice *float64 `json:"recommend_price,omitempty"`
	Brand          *string  `json:"brand,omitempty"`
	Category       *string  `json:"category,omitempty"`
	Description    *string  `json:"description,omitempty"`
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