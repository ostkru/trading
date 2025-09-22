package user

import "time"

// User представляет пользователя
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	APIToken  string    `json:"api_token"`
	TariffID  int64     `json:"tariff_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	APIToken string `json:"api_token" binding:"required,min=1"`
	TariffID int64  `json:"tariff_id,omitempty"`
}

// UpdateUserRequest представляет запрос на обновление пользователя
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=1,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	APIToken *string `json:"api_token,omitempty" binding:"omitempty,min=1"`
	TariffID *int64  `json:"tariff_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UserListResponse представляет ответ со списком пользователей
type UserListResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
