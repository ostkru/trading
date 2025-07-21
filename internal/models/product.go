package models

import (
	"time"
)

type Product struct {
	ID             int        `json:"id" db:"id"`
	Name           string     `json:"name" db:"name" validate:"required"`
	VendorArticle  string     `json:"vendor_article" db:"vendor_article" validate:"required"`
	RecommendPrice *float64   `json:"recommend_price" db:"recommend_price" validate:"required,min=0"`
	Brand          string     `json:"brand" db:"brand" validate:"required"`
	Category       string     `json:"category" db:"category" validate:"required"`
	Description    string     `json:"description" db:"description"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateProductRequest struct {
	Name           string  `json:"name" validate:"required"`
	VendorArticle  string  `json:"vendor_article" validate:"required"`
	RecommendPrice float64 `json:"recommend_price" validate:"required,min=0"`
	Brand          string  `json:"brand" validate:"required"`
	Category       string  `json:"category" validate:"required"`
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

type APIResponse struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Code  int         `json:"code,omitempty"`
} 