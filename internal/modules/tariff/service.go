package tariff

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateTariff создает новый тариф
func (s *Service) CreateTariff(req *CreateTariffRequest) (*Tariff, error) {
	// Валидация входных данных
	if req.Name == "" {
		return nil, errors.New("Требуется название тарифа")
	}
	if req.MinuteLimit <= 0 {
		return nil, errors.New("Минутный лимит должен быть положительным")
	}
	if req.DayLimit <= 0 {
		return nil, errors.New("Дневной лимит должен быть положительным")
	}

	// Проверяем уникальность названия
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tariffs WHERE name = ?", req.Name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности названия: %v", err)
	}
	if count > 0 {
		return nil, errors.New("Тариф с таким названием уже существует")
	}

	// Устанавливаем значение по умолчанию для IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Создаем тариф (адаптировано под существующую структуру)
	features := fmt.Sprintf(`{"daily_requests_limit": %d, "description": "%s"}`, req.DayLimit, req.Description)
	query := `INSERT INTO tariffs (name, price, duration_days, features, is_active) 
	          VALUES (?, 0.00, 365, ?, ?)`

	result, err := s.db.Exec(query, req.Name, features, isActive)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания тарифа: %v", err)
	}

	tariffID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ID тарифа: %v", err)
	}

	// Получаем созданный тариф
	return s.GetTariff(tariffID)
}

// GetTariff получает тариф по ID
func (s *Service) GetTariff(id int64) (*Tariff, error) {
	var tariff Tariff
	err := s.db.QueryRow(`
		SELECT id, name, description, minute_limit, day_limit, is_active, created_at, updated_at
		FROM tariffs WHERE id = ?
	`, id).Scan(
		&tariff.ID,
		&tariff.Name,
		&tariff.Description,
		&tariff.MinuteLimit,
		&tariff.DayLimit,
		&tariff.IsActive,
		&tariff.CreatedAt,
		&tariff.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Тариф не найден")
		}
		return nil, fmt.Errorf("ошибка получения тарифа: %v", err)
	}

	// Отладочный лог
	log.Printf("DEBUG: GetTariff ID=%d, MinuteLimit=%d, DayLimit=%d", tariff.ID, tariff.MinuteLimit, tariff.DayLimit)

	if tariff.MinuteLimit < 60 {
		log.Printf("DEBUG: MinuteLimit %d < 60, setting to 60", tariff.MinuteLimit)
		tariff.MinuteLimit = 60
	}

	return &tariff, nil
}

// ListTariffs получает список тарифов
func (s *Service) ListTariffs(page, limit int, activeOnly bool) (*TariffListResponse, error) {
	offset := (page - 1) * limit
	var whereClause string
	var args []interface{}

	if activeOnly {
		whereClause = "WHERE is_active = ?"
		args = append(args, true)
	}

	// Получаем общее количество
	var total int
	countQuery := "SELECT COUNT(*) FROM tariffs " + whereClause
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("ошибка подсчета тарифов: %v", err)
	}

	// Получаем список тарифов
	query := `
		SELECT id, name, price, duration_days, features, is_active, created_at, updated_at
		FROM tariffs ` + whereClause + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка тарифов: %v", err)
	}
	defer rows.Close()

	var tariffs []Tariff
	for rows.Next() {
		var tariff Tariff
		var features string
		err := rows.Scan(
			&tariff.ID,
			&tariff.Name,
			&tariff.Description, // Используем для price
			&tariff.MinuteLimit, // Используем для duration_days
			&features,
			&tariff.IsActive,
			&tariff.CreatedAt,
			&tariff.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования тарифа: %v", err)
		}

		// Парсим JSON для получения лимитов
		var featuresData map[string]interface{}
		if err := json.Unmarshal([]byte(features), &featuresData); err == nil {
			if dailyLimit, ok := featuresData["daily_requests_limit"].(float64); ok {
				tariff.DayLimit = int(dailyLimit)
			}
			if desc, ok := featuresData["description"].(string); ok {
				tariff.Description = desc
			}
		}

		// Минутный лимит рассчитываем из дневного
		tariff.MinuteLimit = tariff.DayLimit / 1440
		if tariff.MinuteLimit < 60 {
			tariff.MinuteLimit = 60
		}

		tariffs = append(tariffs, tariff)
	}

	return &TariffListResponse{
		Tariffs: tariffs,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// UpdateTariff обновляет тариф
func (s *Service) UpdateTariff(id int64, req *UpdateTariffRequest) (*Tariff, error) {
	// Проверяем существование тарифа
	_, err := s.GetTariff(id)
	if err != nil {
		return nil, err
	}

	// Формируем SET части запроса
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		// Проверяем уникальность названия
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM tariffs WHERE name = ? AND id != ?", *req.Name, id).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки уникальности названия: %v", err)
		}
		if count > 0 {
			return nil, errors.New("Тариф с таким названием уже существует")
		}

		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.MinuteLimit != nil {
		if *req.MinuteLimit <= 0 {
			return nil, errors.New("Минутный лимит должен быть положительным")
		}
		setParts = append(setParts, "minute_limit = ?")
		args = append(args, *req.MinuteLimit)
	}
	if req.DayLimit != nil {
		if *req.DayLimit <= 0 {
			return nil, errors.New("Дневной лимит должен быть положительным")
		}
		setParts = append(setParts, "day_limit = ?")
		args = append(args, *req.DayLimit)
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *req.IsActive)
	}

	if len(setParts) == 0 {
		return s.GetTariff(id)
	}

	// Выполняем обновление
	args = append(args, id)
	query := "UPDATE tariffs SET " + strings.Join(setParts, ", ") + " WHERE id = ?"

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления тарифа: %v", err)
	}

	return s.GetTariff(id)
}

// DeleteTariff удаляет тариф
func (s *Service) DeleteTariff(id int64) error {
	// Проверяем существование тарифа
	_, err := s.GetTariff(id)
	if err != nil {
		return err
	}

	// Проверяем, используется ли тариф пользователями
	var userCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE tariff_id = ?", id).Scan(&userCount)
	if err != nil {
		return fmt.Errorf("ошибка проверки использования тарифа: %v", err)
	}
	if userCount > 0 {
		return errors.New("Нельзя удалить тариф: он используется пользователями")
	}

	// Удаляем тариф
	_, err = s.db.Exec("DELETE FROM tariffs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления тарифа: %v", err)
	}

	return nil
}

// GetUserTariffInfo получает информацию о тарифе пользователя
func (s *Service) GetUserTariffInfo(userID int64) (*UserTariffInfo, error) {
	var info UserTariffInfo
	var features string
	err := s.db.QueryRow(`
		SELECT 
			u.id, u.username, u.email, u.tariff_id,
			t.name as tariff_name, t.features, t.is_active
		FROM users u
		LEFT JOIN tariffs t ON u.tariff_id = t.id
		WHERE u.id = ?
	`, userID).Scan(
		&info.UserID,
		&info.Username,
		&info.Email,
		&info.TariffID,
		&info.TariffName,
		&features,
		&info.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения информации о тарифе пользователя: %v", err)
	}

	// Парсим JSON для получения лимитов
	var featuresData map[string]interface{}
	if err := json.Unmarshal([]byte(features), &featuresData); err == nil {
		if dailyLimit, ok := featuresData["daily_requests_limit"].(float64); ok {
			info.DayLimit = int(dailyLimit)
		}
	}

	// Минутный лимит рассчитываем из дневного
	info.MinuteLimit = info.DayLimit / 1440
	if info.MinuteLimit < 60 {
		info.MinuteLimit = 60
	}

	return &info, nil
}

// ChangeUserTariff изменяет тариф пользователя
func (s *Service) ChangeUserTariff(req *ChangeUserTariffRequest) error {
	// Проверяем существование пользователя
	var userExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", req.UserID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %v", err)
	}
	if !userExists {
		return errors.New("Пользователь не найден")
	}

	// Проверяем существование тарифа
	var tariffExists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tariffs WHERE id = ? AND is_active = TRUE)", req.TariffID).Scan(&tariffExists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования тарифа: %v", err)
	}
	if !tariffExists {
		return errors.New("Тариф не найден или неактивен")
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	// Обновляем тариф пользователя
	_, err = tx.Exec("UPDATE users SET tariff_id = ?, updated_at = ? WHERE id = ?",
		req.TariffID, time.Now(), req.UserID)
	if err != nil {
		return fmt.Errorf("ошибка обновления тарифа пользователя: %v", err)
	}

	// Сбрасываем счетчики rate limiting для пользователя
	_, err = tx.Exec("UPDATE api_rate_limits SET minute_count = 0, day_count = 0, updated_at = ? WHERE user_id = ?",
		time.Now(), req.UserID)
	if err != nil {
		log.Printf("Предупреждение: ошибка сброса счетчиков rate limiting: %v", err)
		// Не прерываем транзакцию, так как это не критично
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %v", err)
	}

	return nil
}

// GetTariffLimits получает лимиты тарифа пользователя
func (s *Service) GetTariffLimits(userID int64) (*TariffLimits, error) {
	var limits TariffLimits
	err := s.db.QueryRow(`
		SELECT t.minute_limit, t.day_limit
		FROM users u
		JOIN tariffs t ON u.tariff_id = t.id
		WHERE u.id = ? AND u.is_active = TRUE AND t.is_active = TRUE
	`, userID).Scan(&limits.MinuteLimit, &limits.DayLimit)

	if err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем дефолтные лимиты если пользователь не найден
			return &TariffLimits{
				MinuteLimit: 1000,
				DayLimit:    10000,
			}, nil
		}
		return nil, fmt.Errorf("ошибка получения лимитов тарифа: %v", err)
	}

	return &limits, nil
}

// GetTariffUsageStats получает статистику использования тарифов
func (s *Service) GetTariffUsageStats() ([]TariffUsageStats, error) {
	query := `
		SELECT 
			t.id as tariff_id,
			t.name as tariff_name,
			COUNT(u.id) as user_count,
			t.minute_limit,
			t.day_limit,
			COALESCE(AVG(arl.minute_count), 0) as avg_minute_usage,
			COALESCE(AVG(arl.day_count), 0) as avg_day_usage,
			MAX(arl.last_request_time) as last_activity
		FROM tariffs t
		LEFT JOIN users u ON t.id = u.tariff_id AND u.is_active = TRUE
		LEFT JOIN api_rate_limits arl ON u.id = arl.user_id
		WHERE t.is_active = TRUE
		GROUP BY t.id, t.name, t.minute_limit, t.day_limit
		ORDER BY t.created_at ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статистики тарифов: %v", err)
	}
	defer rows.Close()

	var stats []TariffUsageStats
	for rows.Next() {
		var stat TariffUsageStats
		var lastActivity sql.NullString

		err := rows.Scan(
			&stat.TariffID,
			&stat.TariffName,
			&stat.UserCount,
			&stat.MinuteLimit,
			&stat.DayLimit,
			&stat.AvgMinuteUsage,
			&stat.AvgDayUsage,
			&lastActivity,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования статистики: %v", err)
		}

		if lastActivity.Valid {
			stat.LastActivity = &lastActivity.String
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
