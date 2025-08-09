package products

import (
	"database/sql"
	"encoding/json"
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
		if err := s.validateURL(url, []string{".obj", ".fbx", ".3ds", ".dae", ".stl", ".glb"}); err != nil {
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

// parseMediaJSON парсит JSON строки медиаданных в слайсы строк
func (s *Service) parseMediaJSON(product *Product) error {
	// Парсим ImageURLs
	if imageURLsStr, ok := product.ImageURLs.(string); ok {
		var imageURLs []string
		if err := json.Unmarshal([]byte(imageURLsStr), &imageURLs); err != nil {
			return fmt.Errorf("error parsing image_urls JSON: %v", err)
		}
		product.ImageURLs = imageURLs
	} else if imageURLsBytes, ok := product.ImageURLs.([]byte); ok {
		var imageURLs []string
		if err := json.Unmarshal(imageURLsBytes, &imageURLs); err != nil {
			return fmt.Errorf("error parsing image_urls JSON: %v", err)
		}
		product.ImageURLs = imageURLs
	}

	// Парсим VideoURLs
	if videoURLsStr, ok := product.VideoURLs.(string); ok {
		var videoURLs []string
		if err := json.Unmarshal([]byte(videoURLsStr), &videoURLs); err != nil {
			return fmt.Errorf("error parsing video_urls JSON: %v", err)
		}
		product.VideoURLs = videoURLs
	} else if videoURLsBytes, ok := product.VideoURLs.([]byte); ok {
		var videoURLs []string
		if err := json.Unmarshal(videoURLsBytes, &videoURLs); err != nil {
			return fmt.Errorf("error parsing video_urls JSON: %v", err)
		}
		product.VideoURLs = videoURLs
	}

	// Парсим Model3DURLs
	if model3DURLsStr, ok := product.Model3DURLs.(string); ok {
		var model3DURLs []string
		if err := json.Unmarshal([]byte(model3DURLsStr), &model3DURLs); err != nil {
			return fmt.Errorf("error parsing model_3d_urls JSON: %v", err)
		}
		product.Model3DURLs = model3DURLs
	} else if model3DURLsBytes, ok := product.Model3DURLs.([]byte); ok {
		var model3DURLs []string
		if err := json.Unmarshal(model3DURLsBytes, &model3DURLs); err != nil {
			return fmt.Errorf("error parsing model_3d_urls JSON: %v", err)
		}
		product.Model3DURLs = model3DURLs
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

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode, user_id, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`

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

	// Обработка NULL значений для barcode
	var barcode interface{}
	if req.Barcode != nil {
		barcode = *req.Barcode
	} else {
		barcode = nil
	}

	result, err := tx.Exec(query, req.Name, req.VendorArticle, req.RecommendPrice, req.Brand, req.Category, brandID, categoryID, req.Description, barcode, userID)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert id: %v", err)
		return nil, err
	}

	// Сохраняем медиаданные в таблицу media
	if len(req.ImageURLs) > 0 || len(req.VideoURLs) > 0 || len(req.Model3DURLs) > 0 {
		mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`

		// Преобразуем слайсы в JSON строки
		imageURLsJSON := "[]"
		videoURLsJSON := "[]"
		model3DURLsJSON := "[]"

		if len(req.ImageURLs) > 0 {
			imageURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(req.ImageURLs, `","`))
		}
		if len(req.VideoURLs) > 0 {
			videoURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(req.VideoURLs, `","`))
		}
		if len(req.Model3DURLs) > 0 {
			model3DURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(req.Model3DURLs, `","`))
		}

		_, err = tx.Exec(mediaQuery, productID, imageURLsJSON, videoURLsJSON, model3DURLsJSON)
		if err != nil {
			log.Printf("Error creating media: %v", err)
			return nil, err
		}
	}

	// Получаем созданный продукт с медиаданными
	var product Product
	err = tx.QueryRow(`
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id
		WHERE p.id = ?
	`, productID).Scan(
		&product.ID,
		&product.Name,
		&product.VendorArticle,
		&product.RecommendPrice,
		&product.Brand,
		&product.Category,
		&product.BrandID,
		&product.CategoryID,
		&product.Description,
		&product.Barcode,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.UserID,
		&product.ImageURLs,
		&product.VideoURLs,
		&product.Model3DURLs,
	)
	if err != nil {
		log.Printf("Error getting created product: %v", err)
		return nil, err
	}

	// Парсим JSON строки в слайсы
	if err := s.parseMediaJSON(&product); err != nil {
		log.Printf("Error parsing media JSON: %v", err)
		return nil, err
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	return &product, nil
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	var product Product
	err := s.db.QueryRow(`
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id
		WHERE p.id = ?
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.VendorArticle,
		&product.RecommendPrice,
		&product.Brand,
		&product.Category,
		&product.BrandID,
		&product.CategoryID,
		&product.Description,
		&product.Barcode,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.UserID,
		&product.ImageURLs,
		&product.VideoURLs,
		&product.Model3DURLs,
	)
	if err != nil {
		return nil, err
	}

	// Парсим JSON строки в слайсы
	if err := s.parseMediaJSON(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *Service) ListProducts(page, limit int, owner string, userID int64) (*ProductListResponse, error) {
	offset := (page - 1) * limit
	var where string
	var args []interface{}

	// Обработка фильтров по владельцу и статусу
	if owner == "my" {
		where = " WHERE p.user_id = ?"
		args = append(args, userID)
	} else if owner == "others" {
		where = " WHERE p.user_id != ?"
		args = append(args, userID)
	} else if owner == "pending" {
		// Показываем продукты со статусом 'pending' (ожидающие классификации)
		where = " WHERE p.status = 'pending'"
	} else if owner == "not_classified" {
		// Показываем продукты со статусом 'pending', у которых нет category_id или brand_id
		where = " WHERE p.status = 'pending' AND (p.category_id IS NULL OR p.brand_id IS NULL)"
	} else if owner == "classified" {
		// Показываем продукты со статусом 'pending', у которых есть и category_id, и brand_id
		where = " WHERE p.status = 'pending' AND p.category_id IS NOT NULL AND p.brand_id IS NOT NULL"
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM products p"+where, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id` + where + `
		ORDER BY p.created_at DESC 
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
		if err := rows.Scan(
			&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice,
			&product.Brand, &product.Category, &product.BrandID, &product.CategoryID,
			&product.Description, &product.Barcode, &product.Status, &product.CreatedAt,
			&product.UpdatedAt, &product.UserID,
			&product.ImageURLs, &product.VideoURLs, &product.Model3DURLs,
		); err != nil {
			return nil, err
		}

		// Парсим JSON строки в слайсы
		if err := s.parseMediaJSON(&product); err != nil {
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

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Формируем SET части запроса для products
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.VendorArticle != nil {
		setParts = append(setParts, "vendor_article = ?")
		args = append(args, *req.VendorArticle)
	}
	if req.RecommendPrice != nil {
		setParts = append(setParts, "recommend_price = ?")
		args = append(args, *req.RecommendPrice)
	}
	if req.Brand != nil {
		setParts = append(setParts, "brand = ?")
		args = append(args, *req.Brand)
	}
	if req.Category != nil {
		setParts = append(setParts, "category = ?")
		args = append(args, *req.Category)
	}
	if req.BrandID != nil {
		setParts = append(setParts, "brand_id = ?")
		args = append(args, *req.BrandID)
	}
	if req.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *req.CategoryID)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Barcode != nil {
		setParts = append(setParts, "barcode = ?")
		args = append(args, *req.Barcode)
	}

	// Обновляем продукт если есть изменения
	if len(setParts) > 0 {
		args = append(args, id)
		query := "UPDATE products SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	// Обновляем медиаданные если они предоставлены
	if req.ImageURLs != nil || req.VideoURLs != nil || req.Model3DURLs != nil {
		// Проверяем существование записи в media
		var mediaID int
		err = tx.QueryRow("SELECT id FROM media WHERE product_id = ?", id).Scan(&mediaID)

		if err == sql.ErrNoRows {
			// Создаем новую запись в media
			imageURLsJSON := "[]"
			videoURLsJSON := "[]"
			model3DURLsJSON := "[]"

			if req.ImageURLs != nil && len(*req.ImageURLs) > 0 {
				imageURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(*req.ImageURLs, `","`))
			}
			if req.VideoURLs != nil && len(*req.VideoURLs) > 0 {
				videoURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(*req.VideoURLs, `","`))
			}
			if req.Model3DURLs != nil && len(*req.Model3DURLs) > 0 {
				model3DURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(*req.Model3DURLs, `","`))
			}

			_, err = tx.Exec("INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)",
				id, imageURLsJSON, videoURLsJSON, model3DURLsJSON)
			if err != nil {
				return nil, err
			}
		} else if err == nil {
			// Обновляем существующую запись в media
			var setMediaParts []string
			var mediaArgs []interface{}

			if req.ImageURLs != nil {
				setMediaParts = append(setMediaParts, "image_urls = ?")
				if len(*req.ImageURLs) > 0 {
					mediaArgs = append(mediaArgs, fmt.Sprintf(`["%s"]`, strings.Join(*req.ImageURLs, `","`)))
				} else {
					mediaArgs = append(mediaArgs, "[]")
				}
			}
			if req.VideoURLs != nil {
				setMediaParts = append(setMediaParts, "video_urls = ?")
				if len(*req.VideoURLs) > 0 {
					mediaArgs = append(mediaArgs, fmt.Sprintf(`["%s"]`, strings.Join(*req.VideoURLs, `","`)))
				} else {
					mediaArgs = append(mediaArgs, "[]")
				}
			}
			if req.Model3DURLs != nil {
				setMediaParts = append(setMediaParts, "model_3d_urls = ?")
				if len(*req.Model3DURLs) > 0 {
					mediaArgs = append(mediaArgs, fmt.Sprintf(`["%s"]`, strings.Join(*req.Model3DURLs, `","`)))
				} else {
					mediaArgs = append(mediaArgs, "[]")
				}
			}

			if len(setMediaParts) > 0 {
				mediaArgs = append(mediaArgs, id)
				mediaQuery := "UPDATE media SET " + strings.Join(setMediaParts, ", ") + " WHERE product_id = ?"
				_, err = tx.Exec(mediaQuery, mediaArgs...)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	// Получаем обновленный продукт с медиаданными
	var product Product
	err = tx.QueryRow(`
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id
		WHERE p.id = ?
	`, id).Scan(
		&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice,
		&product.Brand, &product.Category, &product.BrandID, &product.CategoryID,
		&product.Description, &product.Barcode, &product.Status, &product.CreatedAt,
		&product.UpdatedAt, &product.UserID,
		&product.ImageURLs, &product.VideoURLs, &product.Model3DURLs,
	)
	if err != nil {
		return nil, err
	}

	// Парсим JSON строки в слайсы
	if err := s.parseMediaJSON(&product); err != nil {
		return nil, err
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
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

		// Обработка NULL значений для brand_id, category_id и barcode
		var brandID, categoryID, barcode interface{}
		if p.BrandID != nil {
			brandID = *p.BrandID
		} else {
			brandID = nil
		}
		if p.CategoryID != nil {
			categoryID = *p.CategoryID
		} else {
			categoryID = nil
		}
		if p.Barcode != nil {
			barcode = *p.Barcode
		} else {
			barcode = nil
		}

		query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode, user_id, status) 
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`
		result, err := tx.Exec(query, p.Name, p.VendorArticle, p.RecommendPrice, p.Brand, p.Category, brandID, categoryID, p.Description, barcode, userID)
		if err != nil {
			return nil, err
		}

		productID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		// Сохраняем медиаданные в таблицу media
		if len(p.ImageURLs) > 0 || len(p.VideoURLs) > 0 || len(p.Model3DURLs) > 0 {
			mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`

			// Преобразуем слайсы в JSON строки
			imageURLsJSON := "[]"
			videoURLsJSON := "[]"
			model3DURLsJSON := "[]"

			if len(p.ImageURLs) > 0 {
				imageURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(p.ImageURLs, `","`))
			}
			if len(p.VideoURLs) > 0 {
				videoURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(p.VideoURLs, `","`))
			}
			if len(p.Model3DURLs) > 0 {
				model3DURLsJSON = fmt.Sprintf(`["%s"]`, strings.Join(p.Model3DURLs, `","`))
			}

			_, err = tx.Exec(mediaQuery, productID, imageURLsJSON, videoURLsJSON, model3DURLsJSON)
			if err != nil {
				return nil, err
			}
		}

		// Получаем созданный продукт с медиаданными
		var product Product
		err = tx.QueryRow(`
			SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
			       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
			       p.updated_at, p.user_id,
			       COALESCE(m.image_urls, '[]') as image_urls,
			       COALESCE(m.video_urls, '[]') as video_urls,
			       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
			FROM products p
			LEFT JOIN media m ON p.id = m.product_id
			WHERE p.id = ?
		`, productID).Scan(
			&product.ID, &product.Name, &product.VendorArticle, &product.RecommendPrice,
			&product.Brand, &product.Category, &product.BrandID, &product.CategoryID,
			&product.Description, &product.Barcode, &product.Status, &product.CreatedAt,
			&product.UpdatedAt, &product.UserID,
			&product.ImageURLs, &product.VideoURLs, &product.Model3DURLs,
		)
		if err != nil {
			return nil, err
		}

		// Парсим JSON строки в слайсы
		if err := s.parseMediaJSON(&product); err != nil {
			return nil, err
		}

		createdProducts = append(createdProducts, product)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdProducts, nil
}
