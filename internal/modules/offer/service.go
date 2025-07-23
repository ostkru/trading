package offer

import (
	"database/sql"
	"errors"
	"strings"

	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateOffer(req CreateOfferRequest, userID int64) (*Offer, error) {
	if req.ProductID == 0 || req.OfferType == "" || req.PricePerUnit == 0 || req.AvailableLots == 0 || req.TaxNDS == 0 || req.UnitsPerLot == 0 || req.WarehouseID == 0 {
		return nil, errors.New("Требуются product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id")
	}
	
	// Получаем координаты склада
	var latitude, longitude *float64
	err := s.db.QueryRow("SELECT latitude, longitude FROM warehouses WHERE id = ?", req.WarehouseID).Scan(&latitude, &longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Склад не найден")
		}
		return nil, err
	}
	
	query := `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days, latitude, longitude) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// Устанавливаем значение по умолчанию для is_public если не указано
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	

	
	query = `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days, latitude, longitude) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := s.db.Exec(query, userID, req.ProductID, req.OfferType, req.PricePerUnit, req.AvailableLots, req.TaxNDS, req.UnitsPerLot, req.WarehouseID, isPublic, req.MaxShippingDays, latitude, longitude)
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetOfferByID(id)
}

func (s *Service) UpdateOffer(id int64, req UpdateOfferRequest, userID int64) (*Offer, error) {
	if id == 0 {
		return nil, errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM offers WHERE offer_id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return nil, errors.New("Доступ запрещён")
	} else if err != nil {
		return nil, err
	}
	
	var setClauses []string
	var params []interface{}
	
	if req.PricePerUnit != nil {
		setClauses = append(setClauses, "price_per_unit = ?")
		params = append(params, *req.PricePerUnit)
	}
	if req.AvailableLots != nil {
		setClauses = append(setClauses, "available_lots = ?")
		params = append(params, *req.AvailableLots)
	}
	if req.TaxNDS != nil {
		setClauses = append(setClauses, "tax_nds = ?")
		params = append(params, *req.TaxNDS)
	}
	if req.UnitsPerLot != nil {
		setClauses = append(setClauses, "units_per_lot = ?")
		params = append(params, *req.UnitsPerLot)
	}
	if req.IsPublic != nil {
		setClauses = append(setClauses, "is_public = ?")
		params = append(params, *req.IsPublic)
	}
	if req.MaxShippingDays != nil {
		setClauses = append(setClauses, "max_shipping_days = ?")
		params = append(params, *req.MaxShippingDays)
	}
	
	// Если изменяется warehouse_id, обновляем координаты
	if req.WarehouseID != nil {
		// Получаем координаты нового склада
		var latitude, longitude *float64
		err := s.db.QueryRow("SELECT latitude, longitude FROM warehouses WHERE id = ?", *req.WarehouseID).Scan(&latitude, &longitude)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("Склад не найден")
			}
			return nil, err
		}
		
		setClauses = append(setClauses, "warehouse_id = ?", "latitude = ?", "longitude = ?")
		params = append(params, *req.WarehouseID, latitude, longitude)
	}
	
	if len(setClauses) == 0 {
		return nil, nil
	}
	
	params = append(params, id)
	query := "UPDATE offers SET " + strings.Join(setClauses, ", ") + " WHERE offer_id = ?"
	_, err = s.db.Exec(query, params...)
	if err != nil {
		return nil, err
	}
	return s.GetOfferByID(id)
}

func (s *Service) GetOfferByID(id int64) (*Offer, error) {
	query := `SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days FROM offers WHERE offer_id = ?`
	row := s.db.QueryRow(query, id)
	var offer Offer
	err := row.Scan(
		&offer.OfferID,
		&offer.UserID,
		&offer.UpdatedAt,
		&offer.CreatedAt,
		&offer.IsPublic,
		&offer.ProductID,
		&offer.PricePerUnit,
		&offer.TaxNDS,
		&offer.UnitsPerLot,
		&offer.AvailableLots,
		&offer.Latitude,
		&offer.Longitude,
		&offer.WarehouseID,
		&offer.OfferType,
		&offer.MaxShippingDays,
	)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (s *Service) DeleteOffer(id int64, userID int64) error {
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM offers WHERE offer_id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return errors.New("Доступ запрещён")
	} else if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM offers WHERE offer_id = ?", id)
	return err
}

type OfferListResponse struct {
	Offers []Offer `json:"offers"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

func (s *Service) ListOffers(userID int64, page, limit int, filter string) (*OfferListResponse, error) {
	offset := (page - 1) * limit

	var whereClause string
	var countParams []interface{}
	var queryParams []interface{}

	switch filter {
	case "my":
		whereClause = "WHERE user_id = ?"
		countParams = append(countParams, userID)
		queryParams = append(queryParams, userID)
	case "others":
		whereClause = "WHERE user_id != ?"
		countParams = append(countParams, userID)
		queryParams = append(queryParams, userID)
	case "all":
		whereClause = ""
	default:
		whereClause = "WHERE user_id = ?"
		countParams = append(countParams, userID)
		queryParams = append(queryParams, userID)
	}

	// Подсчет общего количества
	var total int
	countQuery := "SELECT COUNT(*) FROM offers " + whereClause
	err := s.db.QueryRow(countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Получение офферов
	query := `SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days 
	          FROM offers ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	queryParams = append(queryParams, limit, offset)
	rows, err := s.db.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(
			&offer.OfferID,
			&offer.UserID,
			&offer.UpdatedAt,
			&offer.CreatedAt,
			&offer.IsPublic,
			&offer.ProductID,
			&offer.PricePerUnit,
			&offer.TaxNDS,
			&offer.UnitsPerLot,
			&offer.AvailableLots,
			&offer.Latitude,
			&offer.Longitude,
			&offer.WarehouseID,
			&offer.OfferType,
			&offer.MaxShippingDays,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}

	return &OfferListResponse{
		Offers: offers,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

func (s *Service) PublicListOffers() ([]Offer, error) {
	query := `SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days 
	          FROM offers WHERE is_public = true ORDER BY created_at DESC`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(
			&offer.OfferID,
			&offer.UserID,
			&offer.UpdatedAt,
			&offer.CreatedAt,
			&offer.IsPublic,
			&offer.ProductID,
			&offer.PricePerUnit,
			&offer.TaxNDS,
			&offer.UnitsPerLot,
			&offer.AvailableLots,
			&offer.Latitude,
			&offer.Longitude,
			&offer.WarehouseID,
			&offer.OfferType,
			&offer.MaxShippingDays,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

func (s *Service) WBStock(productID, warehouseID, supplierID int64) (int, error) {
	var stock int
	err := s.db.QueryRow("SELECT available_lots FROM offers WHERE product_id = ? AND warehouse_id = ? AND user_id = ?", productID, warehouseID, supplierID).Scan(&stock)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func (s *Service) CreateOffers(req CreateOffersRequest, userID int64) ([]Offer, error) {
	// Проверяем количество офферов
	if len(req.Offers) == 0 {
		return nil, errors.New("Список офферов не может быть пустым")
	}
	if len(req.Offers) > 100 {
		return nil, errors.New("Максимальное количество офферов за один запрос: 100")
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var createdOffers []Offer

	// Подготавливаем запрос для вставки
	query := `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days, latitude, longitude) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, offerReq := range req.Offers {
		// Проверяем обязательные поля
		if offerReq.ProductID == 0 || offerReq.OfferType == "" || offerReq.PricePerUnit == 0 || offerReq.AvailableLots == 0 || offerReq.TaxNDS == 0 || offerReq.UnitsPerLot == 0 || offerReq.WarehouseID == 0 {
			return nil, errors.New("Требуются product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id")
		}

		// Получаем координаты склада
		var latitude, longitude *float64
		err := tx.QueryRow("SELECT latitude, longitude FROM warehouses WHERE id = ?", offerReq.WarehouseID).Scan(&latitude, &longitude)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("Склад не найден")
			}
			return nil, err
		}

		// Устанавливаем значение по умолчанию для is_public если не указано
		isPublic := true
		if offerReq.IsPublic != nil {
			isPublic = *offerReq.IsPublic
		}

		// Вставляем оффер
		result, err := tx.Exec(query, userID, offerReq.ProductID, offerReq.OfferType, offerReq.PricePerUnit, offerReq.AvailableLots, offerReq.TaxNDS, offerReq.UnitsPerLot, offerReq.WarehouseID, isPublic, offerReq.MaxShippingDays, latitude, longitude)
		if err != nil {
			return nil, err
		}

		// Получаем ID созданного оффера
		offerID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		// Получаем полные данные созданного оффера
		var offer Offer
		err = tx.QueryRow(`SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days 
		                   FROM offers WHERE offer_id = ?`, offerID).Scan(
			&offer.OfferID,
			&offer.UserID,
			&offer.UpdatedAt,
			&offer.CreatedAt,
			&offer.IsPublic,
			&offer.ProductID,
			&offer.PricePerUnit,
			&offer.TaxNDS,
			&offer.UnitsPerLot,
			&offer.AvailableLots,
			&offer.Latitude,
			&offer.Longitude,
			&offer.WarehouseID,
			&offer.OfferType,
			&offer.MaxShippingDays,
		)
		if err != nil {
			return nil, err
		}

		createdOffers = append(createdOffers, offer)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return createdOffers, nil
}