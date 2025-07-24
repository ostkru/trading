package database

import (
	"database/sql"
	"fmt"
	"log"

	"portaldata-api/internal/pkg/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection(config config.DatabaseConfig) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
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
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		vendor_article VARCHAR(255) NOT NULL,
		recommend_price DECIMAL(10,2) NOT NULL,
		brand VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(brand, vendor_article)
	);
	
	CREATE INDEX IF NOT EXISTS idx_products_brand_vendor ON products(brand, vendor_article);
	CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
	CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
	
	ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL;
	`

	_, err := db.Exec(query)
	return err
}
