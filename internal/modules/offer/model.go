package offer

type Offer struct {
	OfferID        int64    `json:"offer_id"`
	UserID         int64    `json:"user_id"`
	UpdatedAt      *string  `json:"updated_at,omitempty"`
	CreatedAt      *string  `json:"created_at,omitempty"`
	IsPublic       *bool    `json:"is_public,omitempty"`
	ProductID      *int64   `json:"product_id,omitempty"`
	PricePerUnit   float64  `json:"price_per_unit"`
	TaxNDS         int      `json:"tax_nds"`
	UnitsPerLot    int      `json:"units_per_lot"`
	AvailableLots  int      `json:"available_lots"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	WarehouseID    *int64   `json:"warehouse_id,omitempty"`
	OfferType      string   `json:"offer_type"`
	MaxShippingDays int     `json:"max_shipping_days"`
}

type CreateOfferRequest struct {
	ProductID     int64  `json:"product_id"`
	OfferType     string `json:"offer_type"`
	PricePerUnit  float64 `json:"price_per_unit"`
	AvailableLots int    `json:"available_lots"`
	TaxNDS        int    `json:"tax_nds"`
	UnitsPerLot   int    `json:"units_per_lot"`
	WarehouseID   int64  `json:"warehouse_id"`
	IsPublic      *bool  `json:"is_public,omitempty"`
	MaxShippingDays *int `json:"max_shipping_days,omitempty"`
}

type UpdateOfferRequest struct {
	PricePerUnit  *float64 `json:"price_per_unit,omitempty"`
	AvailableLots *int     `json:"available_lots,omitempty"`
	TaxNDS        *int     `json:"tax_nds,omitempty"`
	UnitsPerLot   *int     `json:"units_per_lot,omitempty"`
	IsPublic      *bool    `json:"is_public,omitempty"`
	MaxShippingDays *int   `json:"max_shipping_days,omitempty"`
} 