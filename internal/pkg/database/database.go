package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"portaldata-api/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewConnection(config config.DatabaseConfig) (*DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройки connection pooling для оптимизации производительности
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(10)                 // Максимальное количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения
	db.SetConnMaxIdleTime(3 * time.Minute) // Время неактивности соединения

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создаем таблицы если их нет
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database connected successfully")
	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	// Таблицы уже существуют в базе данных
	return nil
}
