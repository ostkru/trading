package order

import "time"

type Order struct {
	OrderID            int64     `json:"order_id"`
	TotalAmount        float64   `json:"total_amount"`
	IsMulti            bool      `json:"is_multi"`
	OfferID            *int64    `json:"offer_id,omitempty"`
	InitiatorUserID    int64     `json:"initiator_user_id"`
	CounterpartyUserID *int64    `json:"counterparty_user_id,omitempty"`
	OrderTime          time.Time `json:"order_time"`
	PricePerUnit       float64   `json:"price_per_unit"`
	UnitsPerLot        int       `json:"units_per_lot"`
	LotCount           int       `json:"lot_count"`
	Notes              *string   `json:"notes,omitempty"`
	OrderType          string    `json:"order_type"`
	PaymentMethod      *string   `json:"payment_method,omitempty"`
	OrderStatus        string    `json:"order_status"`
	ShippingAddress    *string   `json:"shipping_address,omitempty"`
	TrackingNumber     *string   `json:"tracking_number,omitempty"`
	MaxShippingDays    int       `json:"max_shipping_days"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

type OrderItem struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	OfferID      int64     `json:"offer_id"`
	Qty          int       `json:"qty"`
	PricePerUnit float64   `json:"price_per_unit"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

type CreateOrderRequest struct {
	OfferID  int64 `json:"offer_id"`
	LotCount int   `json:"quantity"`
}

type GetOrderResponse struct {
	Order      Order       `json:"order"`
	OrderItems []OrderItem `json:"order_items"`
}
