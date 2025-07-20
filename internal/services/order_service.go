package services

import (
	"database/sql"
	"errors"
	"fmt"
	"portaldata-api/internal/database"
	"portaldata-api/internal/models"
	"strings"
)

type OrderService struct {
	db *database.DB
}

func NewOrderService(db *database.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) GetOrder(orderID int64, userID int64) (*models.GetOrderResponse, error) {
	if orderID == 0 {
		return nil, errors.New("order_id required")
	}
	order := models.Order{}
	err := s.db.QueryRow(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number FROM orders WHERE order_id = $1`, orderID).Scan(
		&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &order.OrderStatus, &order.ShippingAddress, &order.TrackingNumber,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("Order not found")
	} else if err != nil {
		return nil, err
	}
	if order.InitiatorUserID != userID && (order.CounterpartyUserID == nil || *order.CounterpartyUserID != userID) {
		return nil, errors.New("Forbidden")
	}
	rows, err := s.db.Query(`SELECT id, order_id, offer_id, qty, price_per_unit, created_at, status FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.OrderItem{}
	for rows.Next() {
		item := models.OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.OfferID, &item.Qty, &item.PricePerUnit, &item.CreatedAt, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &models.GetOrderResponse{Order: order, OrderItems: items}, nil
}

func (s *OrderService) ListOrders(userID int64, status *string, page, perPage int) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	whereParts := []string{"(initiator_user_id = $1 OR counterparty_user_id = $2)"}
	params := []interface{}{userID, userID}
	paramIdx := 3
	if status != nil {
		whereParts = append(whereParts, "order_status = $"+fmt.Sprint(paramIdx))
		params = append(params, *status)
		paramIdx++
	}
	where := " WHERE " + strings.Join(whereParts, " AND ")
	var total int
	totalQuery := "SELECT COUNT(*) FROM orders" + where
	err := s.db.QueryRow(totalQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	params = append(params, perPage, offset)
	orderQuery := fmt.Sprintf(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number FROM orders%s ORDER BY order_id DESC LIMIT $%d OFFSET $%d`, where, paramIdx, paramIdx+1)
	rows, err := s.db.Query(orderQuery, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	orders := []models.Order{}
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &order.OrderStatus, &order.ShippingAddress, &order.TrackingNumber)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	return orders, total, nil
} 