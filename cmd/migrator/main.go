package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"portaldata-api/internal/config"
	"portaldata-api/internal/models"

	_ "github.com/lib/pq"
)

type Migrator struct {
	sourceDB *sql.DB
	targetDB *sql.DB
	apiURL   string
	apiKey   string
}

func NewMigrator(sourceDSN, targetDSN, apiURL, apiKey string) (*Migrator, error) {
	sourceDB, err := sql.Open("postgres", sourceDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to source DB: %w", err)
	}

	targetDB, err := sql.Open("postgres", targetDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target DB: %w", err)
	}

	return &Migrator{
		sourceDB: sourceDB,
		targetDB: targetDB,
		apiURL:   apiURL,
		apiKey:   apiKey,
	}, nil
}

func (m *Migrator) MigrateProducts() error {
	log.Println("Starting product migration...")

	// Получаем данные из старой БД
	rows, err := m.sourceDB.Query(`
		SELECT id, name, vendor_article, price, brand, category, description, created_at, updated_at
		FROM products
		ORDER BY id
	`)
	if err != nil {
		return fmt.Errorf("failed to query source DB: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.VendorArticle,
			&product.Price,
			&product.Brand,
			&product.Category,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	log.Printf("Found %d products to migrate", len(products))

	// Мигрируем данные в новую БД
	for i, product := range products {
		query := `
			INSERT INTO products (name, vendor_article, price, brand, category, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (brand, vendor_article) 
			DO UPDATE SET 
				name = EXCLUDED.name,
				price = EXCLUDED.price,
				category = EXCLUDED.category,
				description = EXCLUDED.description,
				updated_at = EXCLUDED.updated_at
		`

		_, err := m.targetDB.Exec(
			query,
			product.Name,
			product.VendorArticle,
			product.Price,
			product.Brand,
			product.Category,
			product.Description,
			product.CreatedAt,
			product.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert product %d: %w", product.ID, err)
		}

		if (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d products", i+1, len(products))
		}
	}

	log.Printf("Successfully migrated %d products", len(products))
	return nil
}

func (m *Migrator) TestAPI() error {
	log.Println("Testing API connection...")

	// Тестируем health endpoint
	resp, err := http.Get(m.apiURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to test health endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint returned status: %d", resp.StatusCode)
	}

	log.Println("API health check passed")

	// Тестируем создание товара
	testProduct := models.CreateProductRequest{
		Name:          "Test Product",
		VendorArticle: "TEST001",
		Price:         100.50,
		Brand:         "TestBrand",
		Category:      "TestCategory",
		Description:   "Test description",
	}

	payload := models.CreateProductsRequest{
		Products: []models.CreateProductRequest{testProduct},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal test data: %w", err)
	}

	apiURL := fmt.Sprintf("%s/v1/products/products.php?api_key=%s", m.apiURL, m.apiKey)
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = http.NoBody

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send test request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API test failed with status: %d", resp.StatusCode)
	}

	log.Println("API test passed")
	return nil
}

func (m *Migrator) Close() {
	if m.sourceDB != nil {
		m.sourceDB.Close()
	}
	if m.targetDB != nil {
		m.targetDB.Close()
	}
}

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Настройки подключения
	sourceDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, 
		cfg.Database.DBName, cfg.Database.SSLMode)
	
	targetDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"localhost", "5432", "postgres", "portaldata_password_2024", 
		"portaldata", "disable")

	// Создаем мигратор
	migrator, err := NewMigrator(sourceDSN, targetDSN, "http://localhost:8080", cfg.APIKey)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	// Тестируем API
	if err := migrator.TestAPI(); err != nil {
		log.Printf("API test failed: %v", err)
	}

	// Мигрируем данные
	if err := migrator.MigrateProducts(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
} 