package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetUserIDByAPIKey(apiKey string) (int64, error) {
	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_token = ? AND is_active = TRUE", apiKey).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUserByID получает пользователя по ID
func (s *Service) GetUserByID(userID int64) (*User, error) {
	var user User
	err := s.db.QueryRow(`
		SELECT id, username, email, api_token, tariff_id, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.APIToken,
		&user.TariffID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %v", err)
	}
	return &user, nil
}

// CreateUser создает нового пользователя
func (s *Service) CreateUser(req *CreateUserRequest) (*User, error) {
	// Проверяем уникальность username
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности username: %v", err)
	}
	if count > 0 {
		return nil, errors.New("Пользователь с таким username уже существует")
	}

	// Проверяем уникальность email
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности email: %v", err)
	}
	if count > 0 {
		return nil, errors.New("Пользователь с таким email уже существует")
	}

	// Проверяем уникальность API токена
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE api_token = ?", req.APIToken).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки уникальности API токена: %v", err)
	}
	if count > 0 {
		return nil, errors.New("Пользователь с таким API токеном уже существует")
	}

	// Устанавливаем дефолтный тариф если не указан
	tariffID := req.TariffID
	if tariffID == 0 {
		tariffID = 1 // Базовый тариф
	}

	// Создаем пользователя
	query := `INSERT INTO users (username, email, api_token, tariff_id, is_active) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, req.Username, req.Email, req.APIToken, tariffID, true)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ID пользователя: %v", err)
	}

	return s.GetUserByID(userID)
}

// UpdateUser обновляет пользователя
func (s *Service) UpdateUser(userID int64, req *UpdateUserRequest) (*User, error) {
	// Проверяем существование пользователя
	_, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Формируем SET части запроса
	var setParts []string
	var args []interface{}

	if req.Username != nil {
		// Проверяем уникальность username
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", *req.Username, userID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки уникальности username: %v", err)
		}
		if count > 0 {
			return nil, errors.New("Пользователь с таким username уже существует")
		}

		setParts = append(setParts, "username = ?")
		args = append(args, *req.Username)
	}
	if req.Email != nil {
		// Проверяем уникальность email
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", *req.Email, userID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки уникальности email: %v", err)
		}
		if count > 0 {
			return nil, errors.New("Пользователь с таким email уже существует")
		}

		setParts = append(setParts, "email = ?")
		args = append(args, *req.Email)
	}
	if req.APIToken != nil {
		// Проверяем уникальность API токена
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE api_token = ? AND id != ?", *req.APIToken, userID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки уникальности API токена: %v", err)
		}
		if count > 0 {
			return nil, errors.New("Пользователь с таким API токеном уже существует")
		}

		setParts = append(setParts, "api_token = ?")
		args = append(args, *req.APIToken)
	}
	if req.TariffID != nil {
		setParts = append(setParts, "tariff_id = ?")
		args = append(args, *req.TariffID)
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = ?")
		args = append(args, *req.IsActive)
	}

	if len(setParts) == 0 {
		return s.GetUserByID(userID)
	}

	// Выполняем обновление
	args = append(args, userID)
	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id = ?"

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления пользователя: %v", err)
	}

	return s.GetUserByID(userID)
}

// DeleteUser удаляет пользователя
func (s *Service) DeleteUser(userID int64) error {
	// Проверяем существование пользователя
	_, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Удаляем пользователя (каскадное удаление через внешние ключи)
	_, err = s.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %v", err)
	}

	return nil
}
