package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"portaldata-api/internal/modules/offer"
	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateOrder(initiatorUserID int64, req CreateOrderRequest) (*Order, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var offer offer.Offer
	var counterpartyUserID int64
	offerQuery := "SELECT user_id, price_per_unit, units_per_lot, max_shipping_days, available_lots, offer_type FROM offers WHERE offer_id = ? FOR UPDATE"
	err = tx.QueryRow(offerQuery, req.OfferID).Scan(&counterpartyUserID, &offer.PricePerUnit, &offer.UnitsPerLot, &offer.MaxShippingDays, &offer.AvailableLots, &offer.OfferType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("offer not found")
		}
		return nil, err
	}

	if initiatorUserID == counterpartyUserID {
		return nil, errors.New("Нельзя создать заказ на собственное предложение")
	}

	if offer.AvailableLots < req.LotCount {
		return nil, fmt.Errorf("not enough lots available. available: %d, requested: %d", offer.AvailableLots, req.LotCount)
	}

	newAvailableLots := offer.AvailableLots - req.LotCount
	_, err = tx.Exec("UPDATE offers SET available_lots = ? WHERE offer_id = ?", newAvailableLots, req.OfferID)
	if err != nil {
		return nil, err
	}

    var orderType string
    switch offer.OfferType {
    case "sale":
        // Заказ к офферу на продажу — это покупка
        orderType = "buy"
    case "buy":
        // Заказ к офферу на покупку — это продажа
        orderType = "sell"
    default:
        return nil, errors.New("unsupported offer_type for order: must be 'sale' or 'buy'")
    }

	// Вычисляем total_amount
	calculatedTotalAmount := offer.PricePerUnit * float64(req.LotCount) * float64(offer.UnitsPerLot)

	query := `INSERT INTO orders (initiator_user_id, counterparty_user_id, offer_id, lot_count, order_type, price_per_unit, units_per_lot, max_shipping_days, total_amount)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var order Order
	order.InitiatorUserID = initiatorUserID
	order.CounterpartyUserID = &counterpartyUserID
	order.OfferID = &req.OfferID
	order.LotCount = req.LotCount
	order.OrderType = orderType
	order.PricePerUnit = offer.PricePerUnit
	order.UnitsPerLot = offer.UnitsPerLot
	order.MaxShippingDays = offer.MaxShippingDays

	result, err := tx.Exec(query,
		order.InitiatorUserID,
		order.CounterpartyUserID,
		order.OfferID,
		order.LotCount,
		order.OrderType,
		order.PricePerUnit,
		order.UnitsPerLot,
		order.MaxShippingDays,
		calculatedTotalAmount,
	)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Получаем созданный заказ
	var totalAmount sql.NullFloat64
	var isMulti sql.NullBool
	var orderStatus sql.NullString
	err = tx.QueryRow("SELECT order_id, total_amount, is_multi, order_time, order_status, created_at, updated_at FROM orders WHERE order_id = ?", orderID).Scan(
		&order.OrderID, &totalAmount, &isMulti, &order.OrderTime, &orderStatus, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Обработка NULL значения для total_amount
	if totalAmount.Valid {
		order.TotalAmount = &totalAmount.Float64
	}

	// Обработка NULL значения для is_multi
	if isMulti.Valid {
		order.IsMulti = &isMulti.Bool
	}

	// Обработка NULL значения для order_status
	if orderStatus.Valid {
		order.OrderStatus = &orderStatus.String
	}

	return &order, nil
}

func (s *Service) GetOrderByID(orderID, userID int64) (*Order, error) {
	var order Order
	query := `SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days, created_at, updated_at
              FROM orders
              WHERE order_id = ? AND (initiator_user_id = ? OR counterparty_user_id = ?)`

	var totalAmount sql.NullFloat64
	var isMulti sql.NullBool
	var orderStatus sql.NullString
	err := s.db.QueryRow(query, orderID, userID, userID).Scan(
		&order.OrderID,
		&totalAmount,
		&isMulti,
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
		&orderStatus,
		&order.ShippingAddress,
		&order.TrackingNumber,
		&order.MaxShippingDays,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Обработка NULL значения для total_amount
	if totalAmount.Valid {
		order.TotalAmount = &totalAmount.Float64
	}

	// Обработка NULL значения для is_multi
	if isMulti.Valid {
		order.IsMulti = &isMulti.Bool
	}

	// Обработка NULL значения для order_status
	if orderStatus.Valid {
		order.OrderStatus = &orderStatus.String
	}

	return &order, nil
}

func (s *Service) GetOrder(orderID int64, userID int64) (*GetOrderResponse, error) {
	if orderID == 0 {
		return nil, errors.New("order_id required")
	}
	order, err := s.GetOrderByID(orderID, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`SELECT id, order_id, offer_id, qty, price_per_unit, created_at, status FROM order_items WHERE order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OrderItem{}
	for rows.Next() {
		item := OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.OfferID, &item.Qty, &item.PricePerUnit, &item.CreatedAt, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &GetOrderResponse{Order: *order, OrderItems: items}, nil
}

func (s *Service) ListOrders(userID int64, status *string, role string, page, perPage int) ([]Order, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var whereParts []string
	params := []interface{}{}

	switch role {
	case "initiator":
		whereParts = append(whereParts, "initiator_user_id = ?")
		params = append(params, userID)
	case "counterparty":
		whereParts = append(whereParts, "counterparty_user_id = ?")
		params = append(params, userID)
	default:
		whereParts = append(whereParts, "(initiator_user_id = ? OR counterparty_user_id = ?)")
		params = append(params, userID, userID)
	}

	if status != nil && *status != "" {
		whereParts = append(whereParts, "order_status = ?")
		params = append(params, *status)
	}

	where := " WHERE " + strings.Join(whereParts, " AND ")
	var total int
	totalQuery := "SELECT COUNT(*) FROM orders" + where
	err := s.db.QueryRow(totalQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryParamsForOrder := append(params, perPage, offset)
	orderQuery := fmt.Sprintf(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days, created_at, updated_at FROM orders%s ORDER BY order_id DESC LIMIT ? OFFSET ?`, where)
	rows, err := s.db.Query(orderQuery, queryParamsForOrder...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	orders := []Order{}
	for rows.Next() {
		order := Order{}
		var totalAmount sql.NullFloat64
		var isMulti sql.NullBool
		var orderStatus sql.NullString
		err := rows.Scan(&order.OrderID, &totalAmount, &isMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &orderStatus, &order.ShippingAddress, &order.TrackingNumber, &order.MaxShippingDays, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		// Обработка NULL значения для total_amount
		if totalAmount.Valid {
			order.TotalAmount = &totalAmount.Float64
		}

		// Обработка NULL значения для is_multi
		if isMulti.Valid {
			order.IsMulti = &isMulti.Bool
		}

		// Обработка NULL значения для order_status
		if orderStatus.Valid {
			order.OrderStatus = &orderStatus.String
		}

		orders = append(orders, order)
	}
	return orders, total, nil
}

func (s *Service) UpdateOrderStatus(orderID, userID int64, status string) (*Order, error) {
	if orderID == 0 {
		return nil, errors.New("Требуется order_id")
	}

	// Проверяем, что заказ принадлежит пользователю
	var orderUserID int64
	err := s.db.QueryRow("SELECT initiator_user_id FROM orders WHERE order_id = ?", orderID).Scan(&orderUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Order not found")
		}
		return nil, err
	}
	if orderUserID != userID {
		return nil, errors.New("Access denied")
	}

	// Обновляем статус заказа
	_, err = s.db.Exec("UPDATE orders SET order_status = ? WHERE order_id = ?", status, orderID)
	if err != nil {
		return nil, err
	}

	// Получаем обновленный заказ
	return s.GetOrderByID(orderID, userID)
}
