package services

import (
	"database/sql"
	"errors"
	"portaldata-api/internal/database"
	"portaldata-api/internal/models"
	"portaldata-api/internal/utils"
)

type WarehouseService struct {
	db *database.DB
}

func NewWarehouseService(db *database.DB) *WarehouseService {
	return &WarehouseService{db: db}
}

func (s *WarehouseService) CreateWarehouse(req models.CreateWarehouseRequest, userID int64) (*models.Warehouse, error) {
	if req.Name == "" {
		return nil, errors.New("Требуется name")
	}
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO warehouses (user_id, address, latitude, longitude, working_hours, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`,
		userID, req.Address, req.Latitude, req.Longitude, req.WorkingHours,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.Warehouse{
		ID:           id,
		UserID:       userID,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      utils.PtrString(req.Address),
		WorkingHours: utils.PtrString(req.WorkingHours),
	}, nil
}

func (s *WarehouseService) UpdateWarehouse(id int64, req models.UpdateWarehouseRequest, userID int64) (*models.Warehouse, error) {
	if id == 0 {
		return nil, errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM warehouses WHERE id = $1", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return nil, errors.New("Доступ запрещён")
	} else if err != nil {
		return nil, err
	}
	_, err = s.db.Exec(
		`UPDATE warehouses SET address = $1, latitude = $2, longitude = $3, working_hours = $4 WHERE id = $5`,
		getString(req.Address), getFloat(req.Latitude), getFloat(req.Longitude), getString(req.WorkingHours), id,
	)
	if err != nil {
		return nil, err
	}
	return s.GetWarehouseByID(id)
}

func (s *WarehouseService) DeleteWarehouse(id int64, userID int64) error {
	if id == 0 {
		return errors.New("Требуется id")
	}
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM warehouses WHERE id = $1", id).Scan(&dbUserID)
	if err == sql.ErrNoRows || dbUserID != userID {
		return errors.New("Доступ запрещён")
	} else if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM warehouses WHERE id = $1", id)
	return err
}

func (s *WarehouseService) ListWarehouses(userID int64) ([]models.Warehouse, error) {
	rows, err := s.db.Query(`SELECT id, user_id, updated_at, created_at, longitude, latitude, wb_id, working_hours, address FROM warehouses WHERE user_id = $1 ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	warehouses := []models.Warehouse{}
	for rows.Next() {
		wh := models.Warehouse{}
		err := rows.Scan(&wh.ID, &wh.UserID, &wh.UpdatedAt, &wh.CreatedAt, &wh.Longitude, &wh.Latitude, &wh.WBID, &wh.WorkingHours, &wh.Address)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, wh)
	}
	return warehouses, nil
}

func (s *WarehouseService) GetWarehouseByID(id int64) (*models.Warehouse, error) {
	row := s.db.QueryRow(`SELECT id, user_id, updated_at, created_at, longitude, latitude, wb_id, working_hours, address FROM warehouses WHERE id = $1`, id)
	wh := models.Warehouse{}
	err := row.Scan(&wh.ID, &wh.UserID, &wh.UpdatedAt, &wh.CreatedAt, &wh.Longitude, &wh.Latitude, &wh.WBID, &wh.WorkingHours, &wh.Address)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func getString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func getFloat(ptr *float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return 0
} 