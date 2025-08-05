package metaproduct

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

// validateMediaURLs проверяет корректность URL медиа файлов
func (s *Service) validateMediaURLs(imageURLs, videoURLs, model3DURLs []string) error {
	// Проверяем изображения
	for _, url := range imageURLs {
		if err := s.validateURL(url, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}); err != nil {
			return fmt.Errorf("некорректный URL изображения %s: %v", url, err)
		}
	}

	// Проверяем видео
	for _, url := range videoURLs {
		if err := s.validateURL(url, []string{".mp4", ".avi", ".mov", ".wmv", ".flv"}); err != nil {
			return fmt.Errorf("некорректный URL видео %s: %v", url, err)
		}
	}

	// Проверяем 3D модели
	for _, url := range model3DURLs {
		if err := s.validateURL(url, []string{".obj", ".fbx", ".3ds", ".dae", ".stl"}); err != nil {
			return fmt.Errorf("некорректный URL 3D модели %s: %v", url, err)
		}
	}

	return nil
}

// validateURL проверяет корректность URL и расширения файла
func (s *Service) validateURL(urlStr string, allowedExtensions []string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("некорректный формат URL: %v", err)
	}

	// Проверяем протокол
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("неподдерживаемый протокол: %s", parsedURL.Scheme)
	}

	// Проверяем наличие хоста
	if parsedURL.Host == "" {
		return fmt.Errorf("отсутствует хост в URL")
	}

	// Проверяем расширение файла
	path := strings.ToLower(parsedURL.Path)
	hasValidExtension := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(path, ext) {
			hasValidExtension = true
			break
		}
	}

	if !hasValidExtension {
		return fmt.Errorf("неподдерживаемое расширение файла. Разрешены: %v", allowedExtensions)
	}

	return nil
}

func (s *Service) CreateProduct(req CreateProductRequest, userID int64) (*Product, error) {
	if req.Name == "" {
		return nil, errors.New("Требуется name")
	}

	// Валидация медиа URL
	if err := s.validateMediaURLs(req.ImageURLs, req.VideoURLs, req.Model3DURLs); err != nil {
		return nil, err
	}

	query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, user_id) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Обработка NULL значений для brand_id и category_id
	var brandID, categoryID interface{}
	if req.BrandID != nil {
		brandID = *req.BrandID
	} else {
		brandID = nil
	}
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	} else {
		categoryID = nil
	}

	result, err := s.db.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, brandID, categoryID, req.Description, userID)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert id: %v", err)
		return nil, err
	}

	// Получаем созданный продукт
	var product Product
	err = s.db.QueryRow("SELECT id, name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, created_at, updated_at, user_id FROM products WHERE id = ?", productID).Scan(
		&product.ID,
		&product.Name,
		&product.VendorArticle,
		&product.RecommendPrice,
		&product.Brand,
		&product.Category,
		&product.BrandID,
		&product.CategoryID,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.UserID,
	)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}
	return &product, nil
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	var product Product
	err := s.db.QueryRow("SELECT id, name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, created_at, updated_at, user_id FROM products WHERE id = ?", id).Scan(
		&product.ID,
		&product.Name,
		&product.VendorArticle,
		&product.RecommendPrice,
		&product.Brand,
		&product.Category,
		&product.BrandID,
		&product.CategoryID,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.UserID,
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
		SELECT id, name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, created_at, updated_at, user_id
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

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice, &product.Brand, &product.Category, &product.BrandID, &product.CategoryID, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.UserID); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return &ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func (s *Service) UpdateProduct(id int64, req UpdateProductRequest, userID int64) (*Product, error) {
	if id == 0 {
		return nil, errors.New("Требуется id")
	}

	// Проверяем существование продукта и права доступа
	var productUserID int64
	err := s.db.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&productUserID)
	if err == sql.ErrNoRows {
		return nil, errors.New("Продукт с указанным ID не найден")
	} else if err != nil {
		return nil, fmt.Errorf("Ошибка при проверке продукта: %v", err)
	}

	if productUserID != userID {
		return nil, errors.New("Продукт принадлежит другому пользователю")
	}

	// Формируем SET части запроса
	var setParts []string
	var args []interface{}
	argId := 1

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
		argId++
	}
	if req.VendorArticle != nil {
		setParts = append(setParts, "vendor_article = ?")
		args = append(args, *req.VendorArticle)
		argId++
	}
	if req.RecommendPrice != nil {
		setParts = append(setParts, "recommend_price = ?")
		args = append(args, *req.RecommendPrice)
		argId++
	}
	if req.Brand != nil {
		setParts = append(setParts, "brand = ?")
		args = append(args, *req.Brand)
		argId++
	}
	if req.Category != nil {
		setParts = append(setParts, "category = ?")
		args = append(args, *req.Category)
		argId++
	}
	if req.BrandID != nil {
		setParts = append(setParts, "brand_id = ?")
		args = append(args, *req.BrandID)
		argId++
	}
	if req.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *req.CategoryID)
		argId++
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
		argId++
	}

	if len(setParts) == 0 {
		return nil, nil
	}

	args = append(args, id)
	query := "UPDATE products SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Получаем обновленный продукт
	var product Product
	err = s.db.QueryRow("SELECT id, name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, created_at, updated_at, user_id FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice, &product.Brand, &product.Category, &product.BrandID, &product.CategoryID, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.UserID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) DeleteProduct(id int64, userID int64) error {
	if id == 0 {
		return errors.New("Требуется id")
	}

	// Проверяем существование продукта и права доступа
	var productUserID int64
	err := s.db.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&productUserID)
	if err == sql.ErrNoRows {
		return errors.New("Продукт с указанным ID не найден")
	} else if err != nil {
		return fmt.Errorf("Ошибка при проверке продукта: %v", err)
	}

	if productUserID != userID {
		return errors.New("Продукт принадлежит другому пользователю")
	}

	// Проверяем наличие связанных офферов
	var offerCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM offers WHERE product_id = ?", id).Scan(&offerCount)
	if err != nil {
		return fmt.Errorf("Ошибка при проверке связанных офферов: %v", err)
	}

	if offerCount > 0 {
		return errors.New("Нельзя удалить продукт: есть связанные офферы")
	}

	// Удаляем продукт
	_, err = s.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении продукта: %v", err)
	}

	return nil
}

func (s *Service) CreateProducts(req CreateProductsRequest, userID int64) ([]Product, error) {
	if len(req.Products) == 0 {
		return nil, errors.New("Требуется хотя бы один продукт")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var createdProducts []Product

	for _, p := range req.Products {
		// Валидация медиа URL
		if err := s.validateMediaURLs(p.ImageURLs, p.VideoURLs, p.Model3DURLs); err != nil {
			return nil, err
		}

		query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, user_id) 
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := tx.Exec(query, p.Name, p.VendorArticle, p.RecommendPrice, p.Brand, p.Category, p.BrandID, p.CategoryID, p.Description, userID)
		if err != nil {
			return nil, err
		}

		productID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		// Получаем созданный продукт
		var product Product
		err = tx.QueryRow("SELECT id, name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, created_at, updated_at, user_id FROM products WHERE id = ?", productID).Scan(&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice, &product.Brand, &product.Category, &product.BrandID, &product.CategoryID, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.UserID)
		if err != nil {
			return nil, err
		}
		createdProducts = append(createdProducts, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdProducts, nil
}
