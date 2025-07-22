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
	
	query := `INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := s.db.Exec(query, userID, req.ProductID, req.OfferType, req.PricePerUnit, req.AvailableLots, req.TaxNDS, req.UnitsPerLot, req.WarehouseID, req.IsPublic, req.MaxShippingDays)
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

func (s *Service) ListOffers(userID int64, page, limit int) (*OfferListResponse, error) {
	offset := (page - 1) * limit

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM offers WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `SELECT offer_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, latitude, longitude, warehouse_id, offer_type, max_shipping_days 
	          FROM offers WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := s.db.Query(query, userID, limit, offset)
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