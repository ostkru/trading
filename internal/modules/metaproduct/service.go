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
	// Проверяем уникальность комбинации vendor_article + brand
	var existingID int64
	err := s.db.QueryRow("SELECT id FROM products WHERE vendor_article = ? AND brand = ?", req.VendorArticle, req.Brand).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("продукт с артикулом '%s' и брендом '%s' уже существует", req.VendorArticle, req.Brand)
	}
	
	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	// Создаем продукт
	query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, description, user_id, category_id, brand_id) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL)`
	
	result, err := tx.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, req.Description, userID)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}
	
	productID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	// Создаем медиа, если есть данные
	if len(req.ImageURLs) > 0 || len(req.VideoURLs) > 0 || len(req.Model3DURLs) > 0 {
		mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`
		_, err = tx.Exec(mediaQuery, productID, req.ImageURLs, req.VideoURLs, req.Model3DURLs)
		if err != nil {
			log.Printf("Error creating media: %v", err)
			return nil, err
		}
	}
	
	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	
	return s.GetProduct(productID)
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
	
	// Получаем медиа для продукта
	var media Media
	err = s.db.QueryRow("SELECT id, product_id, image_urls, video_urls, model_3d_urls, created_at, updated_at FROM media WHERE product_id = ?", id).Scan(
		&media.ID,
		&media.ProductID,
		&media.ImageURLs,
		&media.VideoURLs,
		&media.Model3DURLs,
		&media.CreatedAt,
		&media.UpdatedAt,
	)
	if err == nil {
		product.Media = &media
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
		
		// Получаем медиа для продукта
		var media Media
		err = s.db.QueryRow("SELECT id, product_id, image_urls, video_urls, model_3d_urls, created_at, updated_at FROM media WHERE product_id = ?", product.ID).Scan(
			&media.ID,
			&media.ProductID,
			&media.ImageURLs,
			&media.VideoURLs,
			&media.Model3DURLs,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err == nil {
			product.Media = &media
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
	query := `SELECT id, user_id, vendor_article, brand FROM products WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&existingProduct.ID, &existingProduct.UserID, &existingProduct.VendorArticle, &existingProduct.Brand)
	if err != nil {
		return nil, err
	}
	
	if existingProduct.UserID != userID {
		return nil, fmt.Errorf("недостаточно прав для обновления продукта")
	}
	
	// Проверяем уникальность, если изменяются vendor_article или brand
	newVendorArticle := existingProduct.VendorArticle
	newBrand := existingProduct.Brand
	if req.VendorArticle != nil {
		newVendorArticle = *req.VendorArticle
	}
	if req.Brand != nil {
		newBrand = *req.Brand
	}
	
	// Если изменились vendor_article или brand, проверяем уникальность
	if newVendorArticle != existingProduct.VendorArticle || newBrand != existingProduct.Brand {
		var duplicateID int64
		err := s.db.QueryRow("SELECT id FROM products WHERE vendor_article = ? AND brand = ? AND id != ?", newVendorArticle, newBrand, id).Scan(&duplicateID)
		if err == nil {
			return nil, fmt.Errorf("продукт с артикулом '%s' и брендом '%s' уже существует", newVendorArticle, newBrand)
		}
	}
	
	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	// Обновляем продукт
	query = `UPDATE products SET 
		name = COALESCE(?, name),
		vendor_article = COALESCE(?, vendor_article),
		recommend_price = COALESCE(?, recommend_price),
		brand = COALESCE(?, brand),
		category = COALESCE(?, category),
		description = COALESCE(?, description),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?`
	
	_, err = tx.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, req.Description, id, userID)
	if err != nil {
		return nil, err
	}
	
	// Обновляем или создаем медиа, если есть данные
	if req.ImageURLs != nil || req.VideoURLs != nil || req.Model3DURLs != nil {
		// Проверяем, существует ли медиа для этого продукта
		var mediaID int64
		err = tx.QueryRow("SELECT id FROM media WHERE product_id = ?", id).Scan(&mediaID)
		
		if err != nil {
			// Медиа не существует, создаем новое
			mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`
			_, err = tx.Exec(mediaQuery, id, req.ImageURLs, req.VideoURLs, req.Model3DURLs)
		} else {
			// Медиа существует, обновляем
			mediaQuery := `UPDATE media SET 
				image_urls = COALESCE(?, image_urls),
				video_urls = COALESCE(?, video_urls),
				model_3d_urls = COALESCE(?, model_3d_urls),
				updated_at = CURRENT_TIMESTAMP
				WHERE product_id = ?`
			_, err = tx.Exec(mediaQuery, req.ImageURLs, req.VideoURLs, req.Model3DURLs, id)
		}
		
		if err != nil {
			log.Printf("Error updating media: %v", err)
			return nil, err
		}
	}
	
	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
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

	// Проверяем уникальность всех продуктов перед вставкой
	for _, productReq := range req.Products {
		var existingID int64
		err := tx.QueryRow("SELECT id FROM products WHERE vendor_article = ? AND brand = ?", productReq.VendorArticle, productReq.Brand).Scan(&existingID)
		if err == nil {
			return nil, fmt.Errorf("продукт с артикулом '%s' и брендом '%s' уже существует", productReq.VendorArticle, productReq.Brand)
		}
	}

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
		
		// Создаем медиа, если есть данные
		if len(productReq.ImageURLs) > 0 || len(productReq.VideoURLs) > 0 || len(productReq.Model3DURLs) > 0 {
			mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`
			_, err = tx.Exec(mediaQuery, id, productReq.ImageURLs, productReq.VideoURLs, productReq.Model3DURLs)
			if err != nil {
				log.Printf("Error creating media for product %d: %v", id, err)
				return nil, err
			}
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