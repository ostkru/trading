package services

import (
	"database/sql"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserIDByAPIKey(apiKey string) (int64, error) {
	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_token = $1", apiKey).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
} 