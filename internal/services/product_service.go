package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"portaldata-api/internal/database"
	"portaldata-api/internal/models"
)

type ProductService struct {
	db *database.DB
}

func NewProductService(db *database.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (name, vendor_article, price, brand, category, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name, vendor_article, price, brand, category, description, created_at, updated_at
	`

	now := time.Now()
	var product models.Product

	err := s.db.QueryRow(
		query,
		req.Name,
		req.VendorArticle,
		req.Price,
		req.Brand,
		req.Category,
		req.Description,
		now,
		now,
	).Scan(
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
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	log.Printf("Created product: %s (ID: %d)", product.Name, product.ID)
	return &product, nil
}

func (s *ProductService) CreateProducts(req models.CreateProductsRequest) ([]models.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var products []models.Product
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
		RETURNING id, name, vendor_article, price, brand, category, description, created_at, updated_at
	`

	now := time.Now()
	for _, reqProduct := range req.Products {
		var product models.Product
		err := tx.QueryRow(
			query,
			reqProduct.Name,
			reqProduct.VendorArticle,
			reqProduct.Price,
			reqProduct.Brand,
			reqProduct.Category,
			reqProduct.Description,
			now,
			now,
		).Scan(
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
			return nil, fmt.Errorf("failed to create product %s: %w", reqProduct.Name, err)
		}

		products = append(products, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Created/Updated %d products", len(products))
	return products, nil
}

func (s *ProductService) GetProduct(id int) (*models.Product, error) {
	query := `
		SELECT id, name, vendor_article, price, brand, category, description, created_at, updated_at
		FROM products WHERE id = $1
	`

	var product models.Product
	err := s.db.QueryRow(query, id).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) ListProducts(page, limit int) (*models.ProductListResponse, error) {
	offset := (page - 1) * limit

	// Получаем общее количество
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Получаем товары
	query := `
		SELECT id, name, vendor_article, price, brand, category, description, created_at, updated_at
		FROM products 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
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
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return &models.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func (s *ProductService) UpdateProduct(id int, req models.UpdateProductRequest) (*models.Product, error) {
	query := `
		UPDATE products 
		SET name = COALESCE($1, name),
			vendor_article = COALESCE($2, vendor_article),
			price = COALESCE($3, price),
			brand = COALESCE($4, brand),
			category = COALESCE($5, category),
			description = COALESCE($6, description),
			updated_at = $7
		WHERE id = $8
		RETURNING id, name, vendor_article, price, brand, category, description, created_at, updated_at
	`

	var product models.Product
	err := s.db.QueryRow(
		query,
		req.Name,
		req.VendorArticle,
		req.Price,
		req.Brand,
		req.Category,
		req.Description,
		time.Now(),
		id,
	).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
} 