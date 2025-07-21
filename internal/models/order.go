package models

type Order struct {
	OrderID            int64    `json:"order_id"`
	TotalAmount        *float64 `json:"total_amount,omitempty"`
	IsMulti            bool     `json:"is_multi"`
	OfferID            *int64   `json:"offer_id,omitempty"`
	InitiatorUserID    int64    `json:"initiator_user_id"`
	CounterpartyUserID *int64   `json:"counterparty_user_id,omitempty"`
	OrderTime          string   `json:"order_time"`
	PricePerUnit       float64  `json:"price_per_unit"`
	UnitsPerLot        int      `json:"units_per_lot"`
	LotCount           int      `json:"lot_count"`
	Notes              *string  `json:"notes,omitempty"`
	OrderType          string   `json:"order_type"`
	PaymentMethod      *string  `json:"payment_method,omitempty"`
	OrderStatus        string   `json:"order_status"`
	ShippingAddress    *string  `json:"shipping_address,omitempty"`
	TrackingNumber     *string  `json:"tracking_number,omitempty"`
	MaxShippingDays    int      `json:"max_shipping_days"`
	CreatedAt          *string `json:"created_at,omitempty"`
	UpdatedAt          *string `json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	OfferID  int64 `json:"offer_id"`
	LotCount int   `json:"lot_count"`
}

type OrderItem struct {
	ID           int64    `json:"id"`
	OrderID      int64    `json:"order_id"`
	OfferID      int64    `json:"offer_id"`
	Qty          int      `json:"qty"`
	PricePerUnit float64  `json:"price_per_unit"`
	CreatedAt    *string  `json:"created_at,omitempty"`
	Status       *string  `json:"status,omitempty"`
}

type GetOrderResponse struct {
	Order      Order       `json:"order"`
	OrderItems []OrderItem `json:"order_items"`
} 