package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"portaldata-api/internal/pkg/database"
	"portaldata-api/internal/modules/offer"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateOrder(initiatorUserID int64, req CreateOrderRequest) (*Order, error) {
	if req.OfferID == 0 || req.LotCount == 0 {
		return nil, errors.New("Требуются offer_id и lot_count")
	}

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
			return nil, errors.New("Предложение не найдено")
		}
		return nil, err
	}

	if initiatorUserID == counterpartyUserID {
		return nil, errors.New("Нельзя создать заказ на собственное предложение")
	}

	if offer.AvailableLots < req.LotCount {
		return nil, fmt.Errorf("недостаточно лотов. доступно: %d, запрошено: %d", offer.AvailableLots, req.LotCount)
	}

	newAvailableLots := offer.AvailableLots - req.LotCount
	_, err = tx.Exec("UPDATE offers SET available_lots = ? WHERE offer_id = ?", newAvailableLots, req.OfferID)
	if err != nil {
		return nil, err
	}

	// Определяем роли в зависимости от типа оффера
	var sellerUserID, buyerUserID int64
	var orderType string
	
	if offer.OfferType == "sale" || offer.OfferType == "sell" {
		// Оффер на продажу: владелец оффера = продавец, создатель заказа = покупатель
		sellerUserID = counterpartyUserID
		buyerUserID = initiatorUserID
		orderType = "buy"
	} else if offer.OfferType == "buy" || offer.OfferType == "purchase" {
		// Оффер на покупку: владелец оффера = покупатель, создатель заказа = продавец
		sellerUserID = initiatorUserID
		buyerUserID = counterpartyUserID
		orderType = "sell"
	} else {
		return nil, fmt.Errorf("неподдерживаемый тип оффера: %s", offer.OfferType)
	}

	// Рассчитываем total_amount
	totalAmount := offer.PricePerUnit * float64(offer.UnitsPerLot) * float64(req.LotCount)
	
	query := `INSERT INTO orders (initiator_user_id, counterparty_user_id, offer_id, lot_count, order_type, price_per_unit, units_per_lot, max_shipping_days, total_amount, order_status, status_changed_by)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := tx.Exec(query,
		buyerUserID,      // initiator_user_id = покупатель
		sellerUserID,     // counterparty_user_id = продавец
		req.OfferID,
		req.LotCount,
		orderType,
		offer.PricePerUnit,
		offer.UnitsPerLot,
		offer.MaxShippingDays,
		totalAmount,
		OrderStatusPending,
		buyerUserID,      // статус изменен покупателем
	)
	if err != nil {
		return nil, err
	}
	
	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var order Order
	order.OrderID = orderID
	order.InitiatorUserID = buyerUserID      // покупатель
	order.CounterpartyUserID = &sellerUserID // продавец
	order.OfferID = &req.OfferID
	order.LotCount = req.LotCount
	order.OrderType = &orderType
	order.PricePerUnit = offer.PricePerUnit
	order.UnitsPerLot = offer.UnitsPerLot
	order.MaxShippingDays = &offer.MaxShippingDays
	order.TotalAmount = totalAmount
	orderStatus := OrderStatusPending
	order.OrderStatus = &orderStatus

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Service) GetOrderByID(orderID, userID int64) (*Order, error) {
	var order Order
	query := `SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, status_reason, status_changed_at, status_changed_by, shipping_address, tracking_number, max_shipping_days, created_at, updated_at
              FROM orders
              WHERE order_id = ? AND (initiator_user_id = ? OR counterparty_user_id = ?)`

	err := s.db.QueryRow(query, orderID, userID, userID).Scan(
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
		&order.StatusReason,
		&order.StatusChangedAt,
		&order.StatusChangedBy,
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
	return &order, nil
}

func (s *Service) GetOrder(orderID int64, userID int64) (*GetOrderResponse, error) {
	if orderID == 0 {
		return nil, errors.New("Требуется order_id")
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

func (s *Service) UpdateOrderStatus(orderID, userID int64, req UpdateOrderStatusRequest) (*Order, error) {
	// Проверяем, что заказ существует и пользователь имеет к нему доступ
	order, err := s.GetOrderByID(orderID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("заказ не найден")
		}
		return nil, err
	}

	// Проверяем, что статус валидный
	validStatuses := map[string]bool{
		OrderStatusPending:    true,
		OrderStatusConfirmed:  true,
		OrderStatusProcessing: true,
		OrderStatusShipped:    true,
		OrderStatusDelivered:  true,
		OrderStatusCancelled:  true,
		OrderStatusRejected:   true,
	}
	if !validStatuses[req.Status] {
		return nil, errors.New("недопустимый статус заказа")
	}

	// Проверяем права на изменение статуса
	if !s.canChangeStatus(order, userID, req.Status) {
		return nil, errors.New("недостаточно прав для изменения статуса")
	}

	// Обновляем статус
	query := `UPDATE orders SET order_status = ?, status_reason = ?, status_changed_at = CURRENT_TIMESTAMP, status_changed_by = ? WHERE order_id = ?`
	result, err := s.db.Exec(query, req.Status, req.Reason, userID, orderID)
	if err != nil {
		return nil, err
	}
	
	// Проверяем, была ли обновлена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	
	if rowsAffected == 0 {
		return nil, errors.New("заказ не найден")
	}

	// Возвращаем обновленный заказ
	updatedOrder, err := s.GetOrderByID(orderID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("заказ не найден после обновления")
		}
		return nil, err
	}
	
	return updatedOrder, nil
}

func (s *Service) canChangeStatus(order *Order, userID int64, newStatus string) bool {
	// Продавец (counterparty_user_id) может изменять статус
	if order.CounterpartyUserID != nil && *order.CounterpartyUserID == userID {
		switch newStatus {
		case OrderStatusConfirmed, OrderStatusProcessing, OrderStatusShipped, OrderStatusCancelled, OrderStatusRejected:
			return true
		}
	}

	// Покупатель (initiator_user_id) может изменять статус
	if order.InitiatorUserID == userID {
		switch newStatus {
		case OrderStatusDelivered, OrderStatusCancelled:
			return true
		}
	}

	return false
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
	orderQuery := fmt.Sprintf(`SELECT order_id, total_amount, is_multi, offer_id, initiator_user_id, counterparty_user_id, order_time, price_per_unit, units_per_lot, lot_count, notes, order_type, payment_method, order_status, shipping_address, tracking_number, max_shipping_days FROM orders%s ORDER BY order_id DESC LIMIT ? OFFSET ?`, where)
	rows, err := s.db.Query(orderQuery, queryParamsForOrder...)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	orders := []Order{}
	for rows.Next() {
		order := Order{}
		err := rows.Scan(&order.OrderID, &order.TotalAmount, &order.IsMulti, &order.OfferID, &order.InitiatorUserID, &order.CounterpartyUserID, &order.OrderTime, &order.PricePerUnit, &order.UnitsPerLot, &order.LotCount, &order.Notes, &order.OrderType, &order.PaymentMethod, &order.OrderStatus, &order.ShippingAddress, &order.TrackingNumber, &order.MaxShippingDays)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	return orders, total, nil
} 