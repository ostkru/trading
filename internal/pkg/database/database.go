package database

import (
	"database/sql"
	"fmt"
	"log"

	"portaldata-api/internal/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewConnection(config config.DatabaseConfig) (*DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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
	// Создаем таблицу products
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		vendor_article VARCHAR(255) NOT NULL,
		recommend_price DECIMAL(10,2) NOT NULL,
		brand VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY unique_brand_vendor (brand, vendor_article)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Создаем индексы если их нет
	indexes := []string{
		"CREATE INDEX idx_products_brand_vendor ON products(brand, vendor_article)",
		"CREATE INDEX idx_products_name ON products(name)",
		"CREATE INDEX idx_products_category ON products(category)",
	}

	for _, indexQuery := range indexes {
		_, err := db.Exec(indexQuery)
		if err != nil {
			// Игнорируем ошибку если индекс уже существует
			log.Printf("Index creation warning: %v", err)
		}
	}

	// Создаем таблицу users если её нет
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		api_token VARCHAR(255) UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(usersQuery)
	return err
} 