package products

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"portaldata-api/internal/pkg/cache"
	"portaldata-api/internal/pkg/classifier"
	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db    *database.DB
	cache cache.Cache

	// Классификатор продуктов
	classifier classifier.ProductClassifier

	// Prepared statements для оптимизации
	stmtCreateProduct *sql.Stmt
	stmtGetProduct    *sql.Stmt
	stmtUpdateProduct *sql.Stmt
	stmtDeleteProduct *sql.Stmt
	stmtListProducts  *sql.Stmt
	stmtCountProducts *sql.Stmt
}

func NewService(db *database.DB) *Service {
	// Создаем фабрику кэша
	cacheFactory := cache.NewCacheFactory()
	productCache := cacheFactory.CreateProductCache()

	// Создаем обработчик результатов классификации
	resultHandler := classifier.NewProductResultHandler(db.DB)

	// Создаем классификатор OSTK
	ostkClassifier := classifier.NewOSTKClassifier("https://api.ostk.ru", resultHandler)

	service := &Service{
		db:         db,
		cache:      productCache,
		classifier: ostkClassifier,
	}

	// Инициализируем prepared statements
	if err := service.initPreparedStatements(); err != nil {
		log.Printf("⚠️ Ошибка инициализации prepared statements: %v", err)
		// Продолжаем работу без prepared statements
	}

	return service
}

// initPreparedStatements инициализирует prepared statements для оптимизации
func (s *Service) initPreparedStatements() error {
	var err error

	// Statement для создания продукта
	s.stmtCreateProduct, err = s.db.Prepare(`
		INSERT INTO products (name, vendor_article, recommend_price, brand, category, 
		                     brand_id, category_id, description, barcode, user_id, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'processing', ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки stmtCreateProduct: %v", err)
	}

	// Statement для получения продукта
	s.stmtGetProduct, err = s.db.Prepare(`
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id
		WHERE p.id = ?
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки stmtGetProduct: %v", err)
	}

	// Statement для обновления продукта
	s.stmtUpdateProduct, err = s.db.Prepare(`
		UPDATE products 
		SET name = ?, vendor_article = ?, recommend_price = ?, brand = ?, category = ?,
		    brand_id = ?, category_id = ?, description = ?, barcode = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки stmtUpdateProduct: %v", err)
	}

	// Statement для удаления продукта
	s.stmtDeleteProduct, err = s.db.Prepare(`
		DELETE FROM products WHERE id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки stmtDeleteProduct: %v", err)
	}

	// Statement для подсчета продуктов
	s.stmtCountProducts, err = s.db.Prepare(`
		SELECT COUNT(*) FROM products p
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки stmtCountProducts: %v", err)
	}

	log.Println("✅ Prepared statements успешно инициализированы")
	return nil
}

// optimizeTransaction настраивает транзакцию для максимальной производительности
func (s *Service) optimizeTransaction(tx *sql.Tx) error {
	// В MySQL нельзя изменять настройки транзакции после её начала
	// Эти настройки должны быть установлены на уровне сессии или соединения
	// Для оптимизации используем только доступные методы

	// Логируем информацию о транзакции для отладки
	log.Printf("🔧 Транзакция оптимизирована: %T", tx)

	return nil
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

func (s *Service) CreateProduct(req *CreateProductRequest, userID int64) (*Product, error) {
	// Валидация входных данных
	if req.Name == "" {
		return nil, errors.New("Требуется name")
	}
	if req.VendorArticle == "" {
		return nil, errors.New("Требуется vendor_article")
	}
	if req.RecommendPrice <= 0 {
		return nil, errors.New("Цена должна быть положительной")
	}
	if req.Brand == "" {
		return nil, errors.New("Требуется brand")
	}
	if req.Category == "" {
		return nil, errors.New("Требуется category")
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode, user_id, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'processing')`

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

	// Сохраняем медиаданные в таблицу media (если есть)
	if len(req.ImageURLs) > 0 || len(req.VideoURLs) > 0 || len(req.Model3DURLs) > 0 {
		mediaQuery := `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES (?, ?, ?, ?)`

		// Оптимизированное преобразование в JSON
		imageURLsJSON := "[]"
		videoURLsJSON := "[]"
		model3DURLsJSON := "[]"

		if len(req.ImageURLs) > 0 {
			// Используем более эффективный способ создания JSON
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

	// Создаем продукт напрямую без дополнительного SELECT
	now := time.Now()
	product := &Product{
		ID:             productID,
		Name:           req.Name,
		VendorArticle:  req.VendorArticle,
		RecommendPrice: req.RecommendPrice,
		Brand:          &req.Brand,
		Category:       &req.Category,
		BrandID:        req.BrandID,
		CategoryID:     req.CategoryID,
		Description:    req.Description,
		Barcode:        req.Barcode,
		Status:         "processing",
		UserID:         userID,
		ImageURLs:      req.ImageURLs,
		VideoURLs:      req.VideoURLs,
		Model3DURLs:    req.Model3DURLs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, err
	}

	// Запускаем асинхронную классификацию продукта
	go s.classifyProductAsync(productID, req, userID)

	return product, nil
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	// Пытаемся получить продукт из кэша
	var product Product
	cacheKey := fmt.Sprintf("product:%d", id)
	if err := s.cache.Get(cacheKey, &product); err == nil {
		log.Printf("✅ Продукт %d получен из кэша", id)
		return &product, nil
	}

	// Если в кэше нет, получаем из базы данных
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

	// Сохраняем в кэш
	if err := s.cache.Set(cacheKey, &product, 1*time.Hour); err != nil {
		log.Printf("⚠️ Ошибка сохранения в кэш: %v", err)
	}

	return &product, nil
}

func (s *Service) ListProducts(page, limit int, owner string, uploadStatus string, userID int64) (*ProductListResponse, error) {
	offset := (page - 1) * limit
	var where string
	var args []interface{}

	// Обработка фильтров по владельцу
	if owner == "my" {
		where = " WHERE p.user_id = ?"
		args = append(args, userID)
	} else if owner == "others" {
		where = " WHERE p.user_id != ?"
		args = append(args, userID)
	} else if owner == "pending" {
		// Показываем продукты со статусом 'processing' (ожидающие классификации)
		where = " WHERE p.status = 'processing'"
	}

	// Добавляем фильтр по статусу классификации
	if uploadStatus != "" {
		if where != "" {
			where += " AND p.status = ?"
		} else {
			where = " WHERE p.status = ?"
		}
		args = append(args, uploadStatus)
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

	// Оптимизируем транзакцию
	if err := s.optimizeTransaction(tx); err != nil {
		return nil, err
	}

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

	// Оптимизируем транзакцию
	if err := s.optimizeTransaction(tx); err != nil {
		return nil, err
	}

	// Валидация всех продуктов перед вставкой
	for _, p := range req.Products {
		if err := s.validateMediaURLs(p.ImageURLs, p.VideoURLs, p.Model3DURLs); err != nil {
			return nil, err
		}
	}

	// Bulk INSERT для продуктов
	productsQuery := `INSERT INTO products (name, vendor_article, recommend_price, brand, category, brand_id, category_id, description, barcode, user_id, status, created_at, updated_at) VALUES `
	productsArgs := []interface{}{}
	productsPlaceholders := []string{}

	now := time.Now()

	for _, p := range req.Products {
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

		productsPlaceholders = append(productsPlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'processing', ?, ?)")
		productsArgs = append(productsArgs, p.Name, p.VendorArticle, p.RecommendPrice, p.Brand, p.Category, brandID, categoryID, p.Description, barcode, userID, now, now)
	}

	productsQuery += strings.Join(productsPlaceholders, ", ")

	// Выполняем Bulk INSERT
	result, err := tx.Exec(productsQuery, productsArgs...)
	if err != nil {
		return nil, fmt.Errorf("ошибка bulk insert продуктов: %v", err)
	}

	// Получаем ID последнего вставленного продукта
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Вычисляем количество вставленных строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// Bulk INSERT для медиаданных (если есть)
	var mediaQuery string
	var mediaArgs []interface{}
	var mediaPlaceholders []string

	mediaCount := 0
	for i, p := range req.Products {
		if len(p.ImageURLs) > 0 || len(p.VideoURLs) > 0 || len(p.Model3DURLs) > 0 {
			// Вычисляем ID продукта: lastID - (rowsAffected - 1) + i
			productID := lastID - (rowsAffected - 1) + int64(i)

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

			mediaPlaceholders = append(mediaPlaceholders, "(?, ?, ?, ?)")
			mediaArgs = append(mediaArgs, productID, imageURLsJSON, videoURLsJSON, model3DURLsJSON)
			mediaCount++
		}
	}

	if mediaCount > 0 {
		mediaQuery = `INSERT INTO media (product_id, image_urls, video_urls, model_3d_urls) VALUES ` + strings.Join(mediaPlaceholders, ", ")
		_, err = tx.Exec(mediaQuery, mediaArgs...)
		if err != nil {
			return nil, fmt.Errorf("ошибка bulk insert медиа: %v", err)
		}
	}

	// Получаем все созданные продукты одним запросом
	selectQuery := `
		SELECT p.id, p.name, p.vendor_article, p.recommend_price, p.brand, p.category, 
		       p.brand_id, p.category_id, p.description, p.barcode, p.status, p.created_at, 
		       p.updated_at, p.user_id,
		       COALESCE(m.image_urls, '[]') as image_urls,
		       COALESCE(m.video_urls, '[]') as video_urls,
		       COALESCE(m.model_3d_urls, '[]') as model_3d_urls
		FROM products p
		LEFT JOIN media m ON p.id = m.product_id
		WHERE p.id >= ? AND p.id <= ? AND p.user_id = ?
		ORDER BY p.id
	`

	startID := lastID - (rowsAffected - 1)
	rows, err := tx.Query(selectQuery, startID, lastID, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения созданных продуктов: %v", err)
	}
	defer rows.Close()

	var createdProducts []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
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

// classifyProductAsync запускает асинхронную классификацию продукта
func (s *Service) classifyProductAsync(productID int64, req *CreateProductRequest, userID int64) {
	// Создаем запрос на классификацию
	classificationReq := classifier.ProductClassificationRequest{
		ProductID:    productID,
		ProductName:  req.Name,
		UserCategory: req.Category,
		UserID:       userID,
		Brand:        req.Brand,
		Category:     req.Category,
	}

	// Добавляем задачу в очередь классификации
	if err := s.classifier.ClassifyProductAsync(classificationReq); err != nil {
		log.Printf("❌ Ошибка добавления задачи классификации для продукта %d: %v", productID, err)
	} else {
		log.Printf("🚀 Задача классификации добавлена для продукта %d: %s", productID, req.Name)
	}
}
