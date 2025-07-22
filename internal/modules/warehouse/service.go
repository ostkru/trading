package warehouse

import (
	"database/sql"
	"errors"

	"portaldata-api/internal/pkg/database"
	"portaldata-api/internal/utils"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateWarehouse(req CreateWarehouseRequest, userID int64) (*Warehouse, error) {
	if req.Name == "" {
		return nil, errors.New("Требуется name")
	}
	query := `INSERT INTO warehouses (user_id, name, address, latitude, longitude, working_hours) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, userID, req.Name, req.Address, req.Latitude, req.Longitude, req.WorkingHours)
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetWarehouseByID(id)
}

func (s *Service) UpdateWarehouse(id int64, req UpdateWarehouseRequest, userID int64) (*Warehouse, error) {
	if id == 0 {
		return nil, errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM warehouses WHERE id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return nil, errors.New("Доступ запрещён")
	} else if err != nil {
		return nil, err
	}
	
	query := `UPDATE warehouses SET name = ?, address = ?, latitude = ?, longitude = ?, working_hours = ? WHERE id = ?`
	_, err = s.db.Exec(query, utils.Coalesce(req.Name, ""), utils.Coalesce(req.Address, ""), utils.Coalesce(req.Latitude, 0.0), utils.Coalesce(req.Longitude, 0.0), utils.Coalesce(req.WorkingHours, ""), id)
	if err != nil {
		return nil, err
	}
	
	return s.GetWarehouseByID(id)
}

func (s *Service) DeleteWarehouse(id int64, userID int64) error {
	if id == 0 {
		return errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM warehouses WHERE id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return errors.New("Доступ запрещён")
	} else if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM warehouses WHERE id = ?", id)
	return err
}

type WarehouseListResponse struct {
	Warehouses []Warehouse `json:"warehouses"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

func (s *Service) ListWarehouses(userID int64, page, limit int) (*WarehouseListResponse, error) {
	offset := (page - 1) * limit
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM warehouses WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, err
	}
	query := `SELECT id, user_id, name, address, latitude, longitude, working_hours, wb_id, created_at, updated_at FROM warehouses WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var warehouses []Warehouse
	for rows.Next() {
		var wh Warehouse
		err := rows.Scan(&wh.ID, &wh.UserID, &wh.Name, &wh.Address, &wh.Latitude, &wh.Longitude, &wh.WorkingHours, &wh.WBID, &wh.CreatedAt, &wh.UpdatedAt)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, wh)
	}
	return &WarehouseListResponse{
		Warehouses: warehouses,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *Service) GetWarehouseByID(id int64) (*Warehouse, error) {
	row := s.db.QueryRow(`SELECT id, user_id, name, address, latitude, longitude, working_hours, wb_id, created_at, updated_at FROM warehouses WHERE id = ?`, id)
	wh := Warehouse{}
	err := row.Scan(&wh.ID, &wh.UserID, &wh.Name, &wh.Address, &wh.Latitude, &wh.Longitude, &wh.WorkingHours, &wh.WBID, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
} 