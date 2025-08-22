package warehouse

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	// Дополнительная валидация на уровне service
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("Требуется name")
	}
	if strings.TrimSpace(req.Address) == "" {
		return nil, errors.New("Требуется address")
	}
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("Latitude должен быть в диапазоне от -90 до 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, errors.New("Longitude должен быть в диапазоне от -180 до 180")
	}

	// MySQL не поддерживает RETURNING, используем отдельные запросы
	result, err := s.db.Exec(`INSERT INTO warehouses (user_id, name, address, latitude, longitude, working_hours) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, req.Name, req.Address, req.Latitude, req.Longitude, req.WorkingHours)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Получаем созданный склад
	var wh Warehouse
	err = s.db.QueryRow(`SELECT id, user_id, name, address, latitude, longitude, working_hours, wb_id, created_at, updated_at FROM warehouses WHERE id = ?`, id).
		Scan(&wh.ID, &wh.UserID, &wh.Name, &wh.Address, &wh.Latitude, &wh.Longitude, &wh.WorkingHours, &wh.WBID, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
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

	// MySQL не поддерживает RETURNING, используем UPDATE + SELECT
	_, err = s.db.Exec(`UPDATE warehouses SET name = ?, address = ?, latitude = ?, longitude = ?, working_hours = ? WHERE id = ?`,
		utils.Coalesce(req.Name, ""), utils.Coalesce(req.Address, ""), utils.Coalesce(req.Latitude, 0.0), utils.Coalesce(req.Longitude, 0.0), utils.Coalesce(req.WorkingHours, ""), id)
	if err != nil {
		return nil, err
	}

	// Получаем обновленный склад
	var wh Warehouse
	err = s.db.QueryRow(`SELECT id, user_id, name, address, latitude, longitude, working_hours, wb_id, created_at, updated_at FROM warehouses WHERE id = ?`, id).
		Scan(&wh.ID, &wh.UserID, &wh.Name, &wh.Address, &wh.Latitude, &wh.Longitude, &wh.WorkingHours, &wh.WBID, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func (s *Service) DeleteWarehouse(id int64, userID int64) error {
	if id == 0 {
		return errors.New("Требуется id")
	}

	// Проверяем существование склада и права доступа
	var dbUserID int64
	err := s.db.QueryRow("SELECT user_id FROM warehouses WHERE id = ?", id).Scan(&dbUserID)
	if err == sql.ErrNoRows {
		return errors.New("Склад с указанным ID не найден")
	} else if err != nil {
		return fmt.Errorf("Ошибка при проверке склада: %v", err)
	}

	if dbUserID != userID {
		return errors.New("Склад принадлежит другому пользователю")
	}

	// Проверяем наличие связанных офферов
	var offerCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM offers WHERE warehouse_id = ?", id).Scan(&offerCount)
	if err != nil {
		return fmt.Errorf("Ошибка при проверке связанных офферов: %v", err)
	}

	if offerCount > 0 {
		return errors.New("Нельзя удалить склад: есть связанные офферы")
	}

	// Удаляем склад
	_, err = s.db.Exec("DELETE FROM warehouses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении склада: %v", err)
	}

	return nil
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

// CreateBatchWarehouses создает несколько складов одновременно
func (s *Service) CreateBatchWarehouses(warehouses []CreateWarehouseRequest, userID int64) ([]Warehouse, error) {
	if len(warehouses) == 0 {
		return nil, errors.New("Список складов пуст")
	}

	if len(warehouses) > 100 {
		return nil, errors.New("Максимальное количество складов в пакете: 100")
	}

	var createdWarehouses []Warehouse

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	for _, req := range warehouses {
		// Валидация данных
		if req.Name == "" {
			return nil, errors.New("Требуется name для всех складов")
		}
		if req.Address == "" {
			return nil, errors.New("Требуется address для всех складов")
		}

		// Создаем склад
		result, err := tx.Exec(`INSERT INTO warehouses (user_id, name, address, latitude, longitude, working_hours) VALUES (?, ?, ?, ?, ?, ?)`,
			userID, req.Name, req.Address, req.Latitude, req.Longitude, req.WorkingHours)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания склада %s: %v", req.Name, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения ID склада %s: %v", req.Name, err)
		}

		// Получаем созданный склад
		var wh Warehouse
		err = tx.QueryRow(`SELECT id, user_id, name, address, latitude, longitude, working_hours, wb_id, created_at, updated_at FROM warehouses WHERE id = ?`, id).
			Scan(&wh.ID, &wh.UserID, &wh.Name, &wh.Address, &wh.Latitude, &wh.Longitude, &wh.WorkingHours, &wh.WBID, &wh.CreatedAt, &wh.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения созданного склада %s: %v", req.Name, err)
		}

		createdWarehouses = append(createdWarehouses, wh)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %v", err)
	}

	return createdWarehouses, nil
}
