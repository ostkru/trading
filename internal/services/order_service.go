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

func (s *OrderService) CreateOrder(initiatorUserID int64, req models.CreateOrderRequest) (*models.Order, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var offer models.Offer
	var counterpartyUserID int64
	offerQuery := "SELECT user_id, price_per_unit, units_per_lot, max_shipping_days, available_lots, offer_type FROM offers WHERE offer_id = $1 FOR UPDATE"
	err = tx.QueryRow(offerQuery, req.OfferID).Scan(&counterpartyUserID, &offer.PricePerUnit, &offer.UnitsPerLot, &offer.MaxShippingDays, &offer.AvailableLots, &offer.OfferType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("offer not found")
		}
		return nil, err
	}

	if initiatorUserID == counterpartyUserID {
		return nil, errors.New("cannot create order on your own offer")
	}

	if offer.AvailableLots < req.LotCount {
		return nil, fmt.Errorf("not enough lots available. available: %d, requested: %d", offer.AvailableLots, req.LotCount)
	}

	newAvailableLots := offer.AvailableLots - req.LotCount
	_, err = tx.Exec("UPDATE offers SET available_lots = $1 WHERE offer_id = $2", newAvailableLots, req.OfferID)
	if err != nil {
		return nil, err
	}

	var orderType string
	if offer.OfferType == "sell" {
		orderType = "buy"
	} else {
		orderType = "sell"
	}

	query := `INSERT INTO orders (initiator_user_id, counterparty_user_id, offer_id, lot_count, order_type, price_per_unit, units_per_lot, max_shipping_days)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              RETURNING order_id, total_amount, is_multi, order_time, order_status, created_at, updated_at`

	var order models.Order
	order.InitiatorUserID = initiatorUserID
	order.CounterpartyUserID = &counterpartyUserID
	order.OfferID = &req.OfferID
	order.LotCount = req.LotCount
	order.OrderType = orderType
	order.PricePerUnit = offer.PricePerUnit
	order.UnitsPerLot = offer.UnitsPerLot
	order.MaxShippingDays = offer.MaxShippingDays

	err = tx.QueryRow(query,
		order.InitiatorUserID,
		order.CounterpartyUserID,
		order.OfferID,
		order.LotCount,
		order.OrderType,
		order.PricePerUnit,
		order.UnitsPerLot,
		order.MaxShippingDays,
	).Scan(&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OrderTime, &order.OrderStatus, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrderByID(orderID, userID int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days, created_at, updated_at
              FROM orders
              WHERE order_id = $1 AND (initiator_user_id = $2 OR counterparty_user_id = $2)`

	err := s.db.QueryRow(query, orderID, userID).Scan(
		&order.OrderID,
		&order.TotalAmount,
		&order.IsMulti,
		&order.OfferID,
		&order.InitiatorUserID,
		&order.CounterpartyUserID,
		&order.OrderTime,
		&order.PricePerUnit,
		&order.UnitsPerLot,
		&order.LotCount,
		&order.Notes,
		&order.OrderType,
		&order.PaymentMethod,
		&order.OrderStatus,
		&order.ShippingAddress,
		&order.TrackingNumber,
		&order.MaxShippingDays,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found or access denied")
		}
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrder(orderID int64, userID int64) (*models.GetOrderResponse, error) {
	if orderID == 0 {
		return nil, errors.New("order_id required")
	}
	order := models.Order{}
	err := s.db.QueryRow(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days FROM orders WHERE order_id = $1`, orderID).Scan(
		&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &order.OrderStatus, &order.ShippingAddress, &order.TrackingNumber, &order.MaxShippingDays,
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

func (s *OrderService) ListOrders(userID int64, status *string, role string, page, perPage int) ([]models.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var whereParts []string
	params := []interface{}{}
	paramIdx := 1

	switch role {
	case "initiator":
		whereParts = append(whereParts, fmt.Sprintf("initiator_user_id = $%d", paramIdx))
		params = append(params, userID)
		paramIdx++
	case "counterparty":
		whereParts = append(whereParts, fmt.Sprintf("counterparty_user_id = $%d", paramIdx))
		params = append(params, userID)
		paramIdx++
	default: // "all"
		whereParts = append(whereParts, fmt.Sprintf("(initiator_user_id = $%d OR counterparty_user_id = $%d)", paramIdx, paramIdx+1))
		params = append(params, userID, userID)
		paramIdx += 2
	}

	if status != nil && *status != "" {
		whereParts = append(whereParts, fmt.Sprintf("order_status = $%d", paramIdx))
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

	queryParamsForOrder := append(params, perPage, offset)
	orderQuery := fmt.Sprintf(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days FROM orders%s ORDER BY order_id DESC LIMIT $%d OFFSET $%d`, where, paramIdx, paramIdx+1)
	rows, err := s.db.Query(orderQuery, queryParamsForOrder...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	orders := []models.Order{}
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &order.OrderStatus, &order.ShippingAddress, &order.TrackingNumber, &order.MaxShippingDays)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	return orders, total, nil
} 