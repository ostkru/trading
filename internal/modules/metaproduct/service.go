package metaproduct

import (
	"fmt"
	"log"

	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateProduct(req CreateProductRequest, userID int64) (*Product, error) {
	query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, description, user_id, category_id, brand_id) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)`
	
	result, err := s.db.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, req.Description, userID)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return s.GetProduct(id)
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	var product Product
	err := s.db.QueryRow("SELECT id, name, vendor_article, recommend_price, brand, category, description, created_at, updated_at, user_id, category_id, brand_id FROM products WHERE id = ?", id).Scan(
		&product.ID,
		&product.Name,
		&product.VendorArticle,
		&product.RecommendPrice,
		&product.Brand,
		&product.Category,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.UserID,
		&product.CategoryID,
		&product.BrandID,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) ListProducts(page, limit int, owner string, userID int64) (*ProductListResponse, error) {
	offset := (page - 1) * limit
	var where string
	var args []interface{}
	if owner == "my" {
		where = " WHERE user_id = ?"
		args = append(args, userID)
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM products"+where, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, vendor_article, recommend_price, brand, category, description, created_at, updated_at, user_id, category_id, brand_id
		FROM products` + where + `
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.VendorArticle,
			&product.RecommendPrice,
			&product.Brand,
			&product.Category,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.UserID,
			&product.CategoryID,
			&product.BrandID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return &ProductListResponse{
		Data:  products,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *Service) UpdateProduct(id int64, req UpdateProductRequest, userID int64) (*Product, error) {
	// Проверяем, что продукт принадлежит пользователю
	var existingProduct Product
	query := `SELECT id, user_id FROM products WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&existingProduct.ID, &existingProduct.UserID)
	if err != nil {
		return nil, err
	}
	
	if existingProduct.UserID != userID {
		return nil, fmt.Errorf("недостаточно прав для обновления продукта")
	}
	
	query = `UPDATE products SET 
		name = COALESCE(?, name),
		vendor_article = COALESCE(?, vendor_article),
		recommend_price = COALESCE(?, recommend_price),
		brand = COALESCE(?, brand),
		category = COALESCE(?, category),
		description = COALESCE(?, description),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?`
	
	_, err = s.db.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, req.Description, id, userID)
	if err != nil {
		return nil, err
	}
	
	return s.GetProduct(id)
}

func (s *Service) DeleteProduct(id int64, userID int64) error {
	// Проверяем, что продукт принадлежит пользователю
	var existingProduct Product
	query := `SELECT id, user_id FROM products WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&existingProduct.ID, &existingProduct.UserID)
	if err != nil {
		return err
	}
	
	if existingProduct.UserID != userID {
		return fmt.Errorf("недостаточно прав для удаления продукта")
	}
	
	_, err = s.db.Exec("DELETE FROM products WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func (s *Service) CreateProducts(req CreateProductsRequest, userID int64) ([]Product, error) {
	if len(req.Products) == 0 {
		return []Product{}, nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var products []Product
	for _, productReq := range req.Products {
		query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, description, user_id, category_id, brand_id) 
		          VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)`
		
		result, err := tx.Exec(query, productReq.Name, productReq.VendorArticle, productReq.RecommendPrice, productReq.Brand, productReq.Category, productReq.Description, userID)
		if err != nil {
			return nil, err
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		
		product := Product{
			ID:             id,
			Name:           productReq.Name,
			VendorArticle:  productReq.VendorArticle,
			RecommendPrice: productReq.RecommendPrice,
			Brand:          productReq.Brand,
			Category:       productReq.Category,
			Description:    productReq.Description,
		}
		products = append(products, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return products, nil
} 