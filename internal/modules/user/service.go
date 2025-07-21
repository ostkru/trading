package user

import (
	"database/sql"
)
type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetUserIDByAPIKey(apiKey string) (int64, error) {
	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE api_token = $1", apiKey).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
} 