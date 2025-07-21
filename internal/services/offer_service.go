package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"portaldata-api/internal/database"
	"portaldata-api/internal/models"
	"strings"
)

type OfferService struct {
	db *database.DB
}

func NewOfferService(db *database.DB) *OfferService {
	return &OfferService{db: db}
}

func (s *OfferService) CreateOffer(userID int64, req models.CreateOfferRequest) (*models.Offer, error) {
	if req.ProductID == 0 || req.OfferType == "" || req.PricePerUnit == 0 || req.AvailableLots == 0 || req.TaxNDS == 0 || req.UnitsPerLot == 0 || req.WarehouseID == 0 {
		return nil, errors.New("Требуются product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id")
	}
	var latitude, longitude float64
	err := s.db.QueryRow("SELECT latitude, longitude FROM warehouses WHERE id = $1", req.WarehouseID).Scan(&latitude, &longitude)
	if err == sql.ErrNoRows {
		return nil, errors.New("Склад не найден")
	} else if err != nil {
		return nil, err
	}
	query := `
        INSERT INTO offers (user_id, product_id, offer_type, price_per_unit, available_lots, tax_nds, units_per_lot, warehouse_id, is_public, max_shipping_days)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING offer_id, created_at, updated_at
    `
	var offer models.Offer
	offer.UserID = userID
	offer.ProductID = &req.ProductID
	offer.OfferType = req.OfferType
	offer.PricePerUnit = req.PricePerUnit
	offer.AvailableLots = req.AvailableLots
	offer.TaxNDS = req.TaxNDS
	offer.UnitsPerLot = req.UnitsPerLot
	offer.WarehouseID = &req.WarehouseID
	
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	offer.IsPublic = &isPublic
	
	maxShippingDays := 0
    if req.MaxShippingDays != nil {
        maxShippingDays = *req.MaxShippingDays
    }
    offer.MaxShippingDays = maxShippingDays


	err = s.db.QueryRow(query,
		offer.UserID,
		offer.ProductID,
		offer.OfferType,
		offer.PricePerUnit,
		offer.AvailableLots,
		offer.TaxNDS,
		offer.UnitsPerLot,
		offer.WarehouseID,
		offer.IsPublic,
		offer.MaxShippingDays,
	).Scan(&offer.OfferID, &offer.CreatedAt, &offer.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (s *OfferService) UpdateOffer(id int64, req models.UpdateOfferRequest, userID int64) (*models.Offer, error) {
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
	setClauses := []string{}
	params := []interface{}{}
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
		return nil, nil // Нет полей для обновления
	}
	params = append(params, id)
	query := "UPDATE offers SET " + strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE offer_id = $%d", idx)
	_, err = s.db.Exec(query, params...)
	if err != nil {
		return nil, err
	}
	return s.GetOfferByID(id)
}

func (s *OfferService) DeleteOffer(id int64, userID int64) error {
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
	// Проверка на связанные активные заказы
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
	Offers []models.Offer `json:"offers"`
	Total  int           `json:"total"`
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
}

func (s *OfferService) ListOffers(userID int64, page, limit int) (*OfferListResponse, error) {
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
	var offers []models.Offer
	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.OfferID, &offer.WBID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.CategoryID, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.OfferName, &offer.Status, &offer.MaxShippingDays,
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

func (s *OfferService) PublicListOffers() ([]models.Offer, error) {
	rows, err := s.db.Query(`SELECT offer_id, wb_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, category_id, latitude, longitude, warehouse_id, offer_type, offer_name, status, max_shipping_days FROM offers WHERE is_public = true ORDER BY offer_id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offers := []models.Offer{}
	for rows.Next() {
		offer := models.Offer{}
		err := rows.Scan(
			&offer.OfferID, &offer.WBID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.CategoryID, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.OfferName, &offer.Status, &offer.MaxShippingDays,
		)
		if err != nil {
			return nil, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (s *OfferService) WBStock(productID, warehouseID, supplierID int64) (int, error) {
	// Аналогично PHP: читаем кэш-файлы WB
	cacheDir := fmt.Sprintf("/var/www/api.portaldata.ru/v1/products/wb/sellers_stock/%d/", supplierID)
	productFile := fmt.Sprintf("%sproduct_%d.json", cacheDir, productID)
	mapFile := "/var/www/api.portaldata.ru/v1/products/wb/sellers_stock/wb2internal_warehouses.json"
	if _, err := os.Stat(productFile); os.IsNotExist(err) {
		return 0, errors.New("Нет кэша WB для данного товара или склада")
	}
	if _, err := os.Stat(mapFile); os.IsNotExist(err) {
		return 0, errors.New("Нет кэша WB для данного товара или склада")
	}
	prodData, err := os.ReadFile(productFile)
	if err != nil {
		return 0, err
	}
	mapData, err := os.ReadFile(mapFile)
	if err != nil {
		return 0, err
	}
	var prod struct {
		Sizes []struct {
			Stocks []struct {
				Wh  int `json:"wh"`
				Qty int `json:"qty"`
			} `json:"stocks"`
		} `json:"sizes"`
	}
	if err := json.Unmarshal(prodData, &prod); err != nil {
		return 0, err
	}
	var mapDataStruct struct {
		Warehouses []struct {
			ID int64 `json:"id"`
		} `json:"warehouses"`
	}
	if err := json.Unmarshal(mapData, &mapDataStruct); err != nil {
		return 0, err
	}
	warehouseMap := make(map[int64]int64)
	for _, w := range mapDataStruct.Warehouses {
		warehouseMap[w.ID] = w.ID
	}
	var totalQty int
	for _, size := range prod.Sizes {
		for _, stock := range size.Stocks {
			if int64(stock.Wh) == warehouseID {
				totalQty += stock.Qty
				break
			}
		}
	}
	return totalQty, nil
}

func (s *OfferService) GetOfferByID(id int64) (*models.Offer, error) {
	row := s.db.QueryRow(`SELECT offer_id, wb_id, user_id, updated_at, created_at, is_public, product_id, price_per_unit, tax_nds, units_per_lot, available_lots, category_id, latitude, longitude, warehouse_id, offer_type, offer_name, status, max_shipping_days FROM offers WHERE offer_id = $1`, id)
	offer := models.Offer{}
	err := row.Scan(
		&offer.OfferID, &offer.WBID, &offer.UserID, &offer.UpdatedAt, &offer.CreatedAt, &offer.IsPublic, &offer.ProductID, &offer.PricePerUnit, &offer.TaxNDS, &offer.UnitsPerLot, &offer.AvailableLots, &offer.CategoryID, &offer.Latitude, &offer.Longitude, &offer.WarehouseID, &offer.OfferType, &offer.OfferName, &offer.Status, &offer.MaxShippingDays,
	)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}