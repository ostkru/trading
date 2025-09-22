package offer

import (
	"database/sql"
	"errors"
	"fmt"
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

	// Валидация цены (проверяем первым)
	if req.PricePerUnit < 0 {
		return nil, errors.New("Цена не может быть отрицательной")
	}

	// Валидация offer_type
	if req.OfferType != "sale" && req.OfferType != "buy" {
		return nil, errors.New("offer_type должен быть 'sale' или 'buy'")
	}

	// Проверяем, что у продукта есть category_id и brand_id - только такие продукты могут быть добавлены в офферы
	var categoryID, brandID sql.NullInt64
	err := s.db.QueryRow("SELECT category_id, brand_id FROM products WHERE id = ?", req.ProductID).Scan(&categoryID, &brandID)
	if err == sql.ErrNoRows {
		return nil, errors.New("Продукт не найден")
	} else if err != nil {
		return nil, fmt.Errorf("Ошибка при проверке продукта: %v", err)
	}

	if !categoryID.Valid || !brandID.Valid {
		return nil, errors.New("Продукт требует классификации. Поля category_id и brand_id должны быть заполнены. Подготовка этих данных может занять некоторое время.")
	}

	// Проверяем, что склад существует
	var warehouseExists bool
	err = s.db.QueryRow("SELECT 1 FROM warehouses WHERE id = ?", req.WarehouseID).Scan(&warehouseExists)
	if err == sql.ErrNoRows {
		return nil, errors.New("Склад не найден")
	} else if err != nil {
		return nil, fmt.Errorf("Ошибка при проверке склада: %v", err)
	}

	// Устанавливаем значение по умолчанию для is_public
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	query := `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, userID, req.ProductID, req.OfferType, req.PricePerUnit, req.AvailableLots, req.TaxNDS, req.UnitsPerLot, req.WarehouseID, isPublic, req.MaxShippingDays)
	if err != nil {
		return nil, err
	}

	offerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Получаем созданный оффер
	return s.GetOfferByID(offerID)
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
	if req.WarehouseID != nil {
		setClauses = append(setClauses, "warehouse_id = ?")
		params = append(params, *req.WarehouseID)

		// Получаем координаты нового склада
		var latitude, longitude float64
		err := s.db.QueryRow("SELECT latitude, longitude FROM warehouses WHERE id = ?", *req.WarehouseID).Scan(&latitude, &longitude)
		if err == nil {
			setClauses = append(setClauses, "latitude = ?")
			setClauses = append(setClauses, "longitude = ?")
			params = append(params, latitude, longitude)
		}
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
	if id == 0 {
		return errors.New("Требуется id")
	}

	// Проверяем существование оффера и права доступа
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM offers WHERE offer_id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		return errors.New("Оффер с указанным ID не найден")
	} else if err != nil {
		return fmt.Errorf("Ошибка при проверке оффера: %v", err)
	}

	if dbUserID != userID {
		return errors.New("Оффер принадлежит другому пользователю")
	}

	// Проверяем наличие активных заказов
	var cnt int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE offer_id = ? AND order_status IN ('pending','active')`, id).Scan(&cnt)
	if err != nil {
		return fmt.Errorf("Ошибка при проверке заказов: %v", err)
	}
	if cnt > 0 {
		return errors.New("Нельзя удалить оффер: есть связанные активные заказы")
	}

	// Удаляем оффер
	_, err = s.db.Exec("DELETE FROM offers WHERE offer_id = ?", id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении оффера: %v", err)
	}

	return nil
}

type OfferListResponse struct {
	Offers []Offer `json:"offers"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

func (s *Service) ListOffers(userID int64, page, limit int, filter, offerType string) (*OfferListResponse, error) {
	return s.ListOffersWithFilters(userID, page, limit, &OfferFilterRequest{
		Filter:    filter,
		OfferType: offerType,
	})
}

// ListOffersWithFilters возвращает список офферов с расширенными фильтрами
func (s *Service) ListOffersWithFilters(userID int64, page, limit int, filters *OfferFilterRequest) (*OfferListResponse, error) {
	offset := (page - 1) * limit

	// Формируем WHERE условия в зависимости от фильтра
	var whereConditions []string
	var params []interface{}

	// Базовый фильтр по пользователю
	switch filters.Filter {
	case "my":
		whereConditions = append(whereConditions, "user_id = ?")
		params = append(params, userID)
	case "others":
		whereConditions = append(whereConditions, "user_id != ?")
		params = append(params, userID)
		whereConditions = append(whereConditions, "is_public = true")
	case "all":
		// Для "all" показываем только публичные офферы других пользователей + все свои
		whereConditions = append(whereConditions, "(user_id = ? OR (user_id != ? AND is_public = true))")
		params = append(params, userID, userID)
	default:
		// По умолчанию используем "my"
		whereConditions = append(whereConditions, "user_id = ?")
		params = append(params, userID)
	}

	// Добавляем фильтр по типу оффера
	if filters.OfferType != "" {
		whereConditions = append(whereConditions, "offer_type = ?")
		params = append(params, filters.OfferType)
	}

	// Добавляем географический фильтр
	if filters.Geographic != nil {
		whereConditions = append(whereConditions, "latitude >= ? AND latitude <= ? AND longitude >= ? AND longitude <= ?")
		params = append(params, filters.Geographic.MinLatitude, filters.Geographic.MaxLatitude,
			filters.Geographic.MinLongitude, filters.Geographic.MaxLongitude)
	}

	// Добавляем фильтр по цене
	if filters.PriceMin != nil {
		whereConditions = append(whereConditions, "price_per_unit >= ?")
		params = append(params, *filters.PriceMin)
	}
	if filters.PriceMax != nil {
		whereConditions = append(whereConditions, "price_per_unit <= ?")
		params = append(params, *filters.PriceMax)
	}

	// Добавляем фильтр по количеству доступных лотов
	if filters.AvailableLots != nil {
		whereConditions = append(whereConditions, "available_lots >= ?")
		params = append(params, *filters.AvailableLots)
	}

	// Формируем WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Подсчет общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM offers %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, params...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Получение офферов с пагинацией
	query := fmt.Sprintf(`SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days FROM offers %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	params = append(params, limit, offset)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(&offer.OfferID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.MaxShippingDays)
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

func (s *Service) WBStock(productID, warehouseID, supplierID int64) (int, error) {
	query := `SELECT available_lots FROM offers WHERE product_id = ? AND warehouse_id = ? AND user_id = ? AND is_public = true LIMIT 1`
	var stock int
	err := s.db.QueryRow(query, productID, warehouseID, supplierID).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Возвращаем 0 вместо ошибки
		}
		return 0, err
	}
	return stock, nil
}

// PublicListOffers возвращает список публичных офферов
func (s *Service) PublicListOffers(page, limit int) (*OfferListResponse, error) {
	return s.PublicListOffersWithFilters(page, limit, &OfferFilterRequest{})
}

// PublicListOffersWithFilters возвращает список публичных офферов с фильтрами
func (s *Service) PublicListOffersWithFilters(page, limit int, filters *OfferFilterRequest) (*OfferListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Формируем WHERE условия
	var whereConditions []string
	var params []interface{}

	// Базовое условие - только публичные офферы
	whereConditions = append(whereConditions, "o.is_public = true")

	// Добавляем фильтр по типу оффера
	if filters.OfferType != "" {
		whereConditions = append(whereConditions, "o.offer_type = ?")
		params = append(params, filters.OfferType)
	}

	// Добавляем географический фильтр
	if filters.Geographic != nil {
		whereConditions = append(whereConditions, "w.latitude >= ? AND w.latitude <= ? AND w.longitude >= ? AND w.longitude <= ?")
		params = append(params, filters.Geographic.MinLatitude, filters.Geographic.MaxLatitude,
			filters.Geographic.MinLongitude, filters.Geographic.MaxLongitude)
	}

	// Добавляем фильтр по цене
	if filters.PriceMin != nil {
		whereConditions = append(whereConditions, "o.price_per_unit >= ?")
		params = append(params, *filters.PriceMin)
	}
	if filters.PriceMax != nil {
		whereConditions = append(whereConditions, "o.price_per_unit <= ?")
		params = append(params, *filters.PriceMax)
	}

	// Добавляем фильтр по количеству доступных лотов
	if filters.AvailableLots != nil {
		whereConditions = append(whereConditions, "o.available_lots >= ?")
		params = append(params, *filters.AvailableLots)
	}

	// Добавляем фильтр по ID бренда
	if filters.BrandID != nil {
		whereConditions = append(whereConditions, "p.brand_id = ?")
		params = append(params, *filters.BrandID)
	}

	// Добавляем фильтр по ID категории
	if filters.CategoryID != nil {
		whereConditions = append(whereConditions, "p.category_id = ?")
		params = append(params, *filters.CategoryID)
	}

	// Добавляем фильтр по названию продукта (поиск)
	if filters.ProductName != nil && *filters.ProductName != "" {
		whereConditions = append(whereConditions, "p.name LIKE ?")
		params = append(params, "%"+*filters.ProductName+"%")
	}

	// Добавляем фильтр по артикулу производителя
	if filters.VendorArticle != nil && *filters.VendorArticle != "" {
		whereConditions = append(whereConditions, "p.vendor_article LIKE ?")
		params = append(params, "%"+*filters.VendorArticle+"%")
	}

	// Добавляем фильтр по ID склада
	if filters.WarehouseID != nil {
		whereConditions = append(whereConditions, "o.warehouse_id = ?")
		params = append(params, *filters.WarehouseID)
	}

	// Добавляем фильтр по НДС
	if filters.TaxNDS != nil {
		whereConditions = append(whereConditions, "o.tax_nds = ?")
		params = append(params, *filters.TaxNDS)
	}

	// Добавляем фильтр по количеству единиц в лоте
	if filters.UnitsPerLot != nil {
		whereConditions = append(whereConditions, "o.units_per_lot = ?")
		params = append(params, *filters.UnitsPerLot)
	}

	// Добавляем фильтр по максимальным дням доставки
	if filters.MaxShippingDays != nil {
		whereConditions = append(whereConditions, "o.max_shipping_days <= ?")
		params = append(params, *filters.MaxShippingDays)
	}

	// Формируем WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Получаем общее количество публичных офферов с фильтрами
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM offers o
		LEFT JOIN products p ON o.product_id = p.id
		LEFT JOIN warehouses w ON o.warehouse_id = w.id
		%s`, whereClause)
	var total int
	err := s.db.QueryRow(countQuery, params...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Получаем публичные офферы с пагинацией и фильтрами
	query := fmt.Sprintf(`
		SELECT o.offer_id, o.user_id, o.created_at, o.updated_at, o.is_public, 
		       o.product_id, o.price_per_unit, o.tax_nds, o.units_per_lot, 
		       o.available_lots, o.warehouse_id, o.offer_type, o.max_shipping_days,
		       p.name as product_name, p.vendor_article, p.recommend_price,
		       w.name as warehouse_name, w.address, w.latitude, w.longitude
		FROM offers o
		LEFT JOIN products p ON o.product_id = p.id
		LEFT JOIN warehouses w ON o.warehouse_id = w.id
		%s
		ORDER BY o.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	params = append(params, limit, offset)
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []Offer
	for rows.Next() {
		var offer Offer
		var productName, vendorArticle sql.NullString
		var recommendPrice sql.NullFloat64
		var warehouseName, warehouseAddress sql.NullString
		var warehouseLat, warehouseLng sql.NullFloat64

		err := rows.Scan(
			&offer.OfferID, &offer.UserID, &offer.CreatedAt, &offer.UpdatedAt, &offer.IsPublic,
			&offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot,
			&offer.AvailableLots, &offer.WarehouseID, &offer.OfferType, &offer.MaxShippingDays,
			&productName, &vendorArticle, &recommendPrice,
			&warehouseName, &warehouseAddress, &warehouseLat, &warehouseLng,
		)
		if err != nil {
			return nil, err
		}

		// Добавляем информацию о продукте и складе
		if productName.Valid {
			offer.ProductName = &productName.String
		}
		if vendorArticle.Valid {
			offer.VendorArticle = &vendorArticle.String
		}
		if recommendPrice.Valid {
			offer.RecommendPrice = &recommendPrice.Float64
		}
		if warehouseName.Valid {
			offer.WarehouseName = &warehouseName.String
		}
		if warehouseAddress.Valid {
			offer.WarehouseAddress = &warehouseAddress.String
		}
		if warehouseLat.Valid {
			offer.Latitude = &warehouseLat.Float64
		}
		if warehouseLng.Valid {
			offer.Longitude = &warehouseLng.Float64
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

func (s *Service) CreateOffers(req CreateOffersRequest, userID int64) ([]Offer, error) {
	var offers []Offer

	for _, offerReq := range req.Offers {
		offer, err := s.CreateOffer(offerReq, userID)
		if err != nil {
			return nil, err
		}
		offers = append(offers, *offer)
	}

	return offers, nil
}
