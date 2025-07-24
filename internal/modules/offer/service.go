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

	// Валидация offer_type
	if req.OfferType != "sale" && req.OfferType != "buy" {
		return nil, errors.New("offer_type должен быть 'sale' или 'buy'")
	}

	// Устанавливаем значение по умолчанию для is_public
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	query := `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING offer_id, user_id, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, warehouse_id, offer_type, max_shipping_days`
	var offer Offer
	err := s.db.QueryRow(query, userID, req.ProductID, req.OfferType, req.PricePerUnit, req.AvailableLots, req.TaxNDS, req.UnitsPerLot, req.WarehouseID, isPublic, req.MaxShippingDays).Scan(
		&offer.OfferID, &offer.UserID, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.WarehouseID, &offer.OfferType, &offer.MaxShippingDays,
	)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (s *Service) UpdateOffer(id int64, req UpdateOfferRequest, userID int64) (*Offer, error) {
	if id == 0 {
		return nil, errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM offers WHERE offer_id = $1", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return nil, errors.New("Доступ запрещён")
	} else if err != nil {
		return nil, err
	}
	var setClauses []string
	var params []interface{}
	idx := 1
	if req.PricePerUnit != nil {
		setClauses = append(setClauses, fmt.Sprintf("price_per_unit = $%d", idx))
		params = append(params, *req.PricePerUnit)
		idx++
	}
	if req.AvailableLots != nil {
		setClauses = append(setClauses, fmt.Sprintf("available_lots = $%d", idx))
		params = append(params, *req.AvailableLots)
		idx++
	}
	if req.TaxNDS != nil {
		setClauses = append(setClauses, fmt.Sprintf("tax_nds = $%d", idx))
		params = append(params, *req.TaxNDS)
		idx++
	}
	if req.UnitsPerLot != nil {
		setClauses = append(setClauses, fmt.Sprintf("units_per_lot = $%d", idx))
		params = append(params, *req.UnitsPerLot)
		idx++
	}
	if req.IsPublic != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_public = $%d", idx))
		params = append(params, *req.IsPublic)
		idx++
	}
	if req.MaxShippingDays != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_shipping_days = $%d", idx))
		params = append(params, *req.MaxShippingDays)
		idx++
	}
	if len(setClauses) == 0 {
		return nil, nil
	}
	params = append(params, id)
	query := "UPDATE offers SET " + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE offer_id = $%d", idx)
	_, err = s.db.Exec(query, params...)
	if err != nil {
		return nil, err
	}
	return s.GetOfferByID(id)
}

func (s *Service) GetOfferByID(id int64) (*Offer, error) {
	query := `SELECT offer_id, wb_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, category_id, latitude, longitude, warehouse_id, offer_type, offer_name, status, max_shipping_days FROM offers WHERE offer_id = $1`
	row := s.db.QueryRow(query, id)
	var offer Offer
	err := row.Scan(
		&offer.OfferID,
		&offer.WBID,
		&offer.UserID,
		&offer.UpdatedAt,
		&offer.CreatedAt,
		&offer.IsPublic,
		&offer.ProductID,
		&offer.PricePerUnit,
		&offer.TaxNDS,
		&offer.UnitsPerLot,
		&offer.AvailableLots,
		&offer.CategoryID,
		&offer.Latitude,
		&offer.Longitude,
		&offer.WarehouseID,
		&offer.OfferType,
		&offer.OfferName,
		&offer.Status,
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
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM offers WHERE offer_id = $1", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return errors.New("Доступ запрещён")
	} else if err != nil {
		return err
	}
	var cnt int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE offer_id = $1 AND order_status IN ('pending','active')`, id).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("Нельзя удалить оффер: есть связанные активные заказы")
	}
	_, err = s.db.Exec("DELETE FROM offers WHERE offer_id = $1", id)
	return err
}

type OfferListResponse struct {
	Offers []Offer `json:"offers"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

func (s *Service) ListOffers(userID int64, page, limit int) (*OfferListResponse, error) {
	offset := (page - 1) * limit
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM offers WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, err
	}
	query := `SELECT offer_id, wb_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, category_id, latitude, longitude, warehouse_id, offer_type, offer_name, status, max_shipping_days FROM offers WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(&offer.OfferID, &offer.WBID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.CategoryID, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.OfferName, &offer.Status, &offer.MaxShippingDays)
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
	rows, err := s.db.Query(`SELECT offer_id, wb_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, category_id, latitude, longitude, warehouse_id, offer_type, offer_name, status, max_shipping_days FROM offers WHERE is_public = true ORDER BY offer_id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []Offer
	for rows.Next() {
		var offer Offer
		err := rows.Scan(&offer.OfferID, &offer.WBID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.CategoryID, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.OfferName, &offer.Status, &offer.MaxShippingDays)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (s *Service) WBStock(productID, warehouseID, supplierID int64) (int, error) {
	var quantity int
	err := s.db.QueryRow("SELECT quantity FROM offers WHERE product_id = $1 AND warehouse_id = $2 AND supplier_id = $3", productID, warehouseID, supplierID).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return quantity, nil
}
