package categories

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"portaldata-api/internal/pkg/database"
)

type Service struct {
	db    *database.DB
	redis *redis.Client
}

func NewService(db *database.DB, redis *redis.Client) *Service {
	return &Service{db: db, redis: redis}
}

// ListAllCategoryIDsAndNames возвращает все уникальные пары category_id + category_name из MySQL с кэшированием
func (s *Service) ListAllCategoryIDsAndNames(ctx context.Context) (*ListCategoryIDNameResponse, error) {
	// Проверить кэш в Redis
	cacheKey := "categories:list"
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var result ListCategoryIDNameResponse
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				fmt.Printf("DEBUG: Categories list served from Redis cache\n")
				return &result, nil
			}
		}
	}

	// Если нет в кэше - получить из MySQL
	fmt.Printf("DEBUG: Categories list not in cache, fetching from MySQL\n")
    rows, err := s.db.QueryContext(ctx, `
        SELECT DISTINCT category_id, category_name
        FROM products
        WHERE category_id IS NOT NULL
          AND category_name IS NOT NULL
          AND category_name <> ''
        ORDER BY category_name
    `)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
    }
    defer rows.Close()

    result := &ListCategoryIDNameResponse{Categories: make([]CategoryIDName, 0)}
    for rows.Next() {
        var id sql.NullInt64
        var name sql.NullString
        if err := rows.Scan(&id, &name); err != nil {
            return nil, fmt.Errorf("ошибка чтения строки: %v", err)
        }
        if id.Valid && name.Valid {
            result.Categories = append(result.Categories, CategoryIDName{ID: id.Int64, Name: name.String})
        }
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка курсора: %v", err)
    }
    result.Total = int64(len(result.Categories))

	// Сохранить в кэш на 1 час
	if s.redis != nil {
		data, err := json.Marshal(result)
		if err == nil {
			s.redis.Set(ctx, cacheKey, data, time.Hour)
			fmt.Printf("DEBUG: Categories list cached in Redis for 1 hour\n")
		}
	}

    return result, nil
}

// ListCategories получает список категорий с фильтрацией и пагинацией
func (s *Service) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	// Устанавливаем значения по умолчанию
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Строим WHERE условия
	whereConditions := []string{}
	args := []interface{}{}

	// Поиск по названию
	if req.Search != "" {
		whereConditions = append(whereConditions, "category LIKE ?")
		args = append(args, "%"+req.Search+"%")
	}

	// Фильтр по родительской категории (пока не реализован)
	// if req.ParentID != nil {
	//     whereConditions = append(whereConditions, "parent_id = ?")
	//     args = append(args, *req.ParentID)
	// }

	// Фильтр по активности (пока не реализован)
	// if req.Active != nil {
	//     whereConditions = append(whereConditions, "is_active = ?")
	//     args = append(args, *req.Active)
	// }

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Строим ORDER BY
	orderBy := "category ASC"
	switch req.Sort {
	case "name_asc":
		orderBy = "category ASC"
	case "name_desc":
		orderBy = "category DESC"
	case "products_asc":
		orderBy = "COUNT(*) ASC"
	case "products_desc":
		orderBy = "COUNT(*) DESC"
	}

	// Получаем общее количество
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT category) 
		FROM products 
		%s`, whereClause)
	
	var total int64
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("ошибка подсчета категорий: %v", err)
	}

	// Получаем категории с подсчетом продуктов
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT 
			category,
			COUNT(*) as product_count
		FROM products 
		%s
		GROUP BY category
		ORDER BY %s
		LIMIT ? OFFSET ?`, whereClause, orderBy)

	args = append(args, req.Limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения категорий: %v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		
		err := rows.Scan(
			&category.Name,
			&category.ProductCount,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования категории: %v", err)
		}

		// Генерируем ID для категории (если это WB категория)
		if strings.HasPrefix(strings.ToLower(category.Name), "wb:") {
			category.ID = generateCategoryID(category.Name)
		}

		category.CreatedAt = time.Now() // Упрощенная версия
		category.UpdatedAt = time.Now() // Упрощенная версия
		category.IsActive = true // Все категории из продуктов считаются активными

		categories = append(categories, category)
	}

	return &ListCategoriesResponse{
		Categories: categories,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

// GetCategoryCharacteristics получает характеристики для категории
func (s *Service) GetCategoryCharacteristics(ctx context.Context, categoryName string) (*GetCategoryCharacteristicsResponse, error) {
	// Получаем информацию о категории
	categoryQuery := `
		SELECT 
			category,
			COUNT(*) as product_count,
			MIN(created_at) as first_created,
			MAX(updated_at) as last_updated
		FROM products 
		WHERE category = ?
		GROUP BY category`

	var category Category
	var firstCreated, lastUpdated time.Time
	var productCount int64

	err := s.db.QueryRow(categoryQuery, categoryName).Scan(
		&category.Name,
		&productCount,
		&firstCreated,
		&lastUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("категория не найдена")
		}
		return nil, fmt.Errorf("ошибка получения категории: %v", err)
	}

	// Генерируем ID для категории
	if strings.HasPrefix(strings.ToLower(category.Name), "wb:") {
		category.ID = generateCategoryID(category.Name)
	}

	category.ProductCount = productCount
	category.CreatedAt = firstCreated
	category.UpdatedAt = lastUpdated
	category.IsActive = true

	// Получаем характеристики из OpenSearch (если доступен)
	characteristics, err := s.getCharacteristicsFromOpenSearch(ctx, categoryName)
	if err != nil {
		// Если OpenSearch недоступен, возвращаем пустой список характеристик
		characteristics = []Characteristic{}
	}

	return &GetCategoryCharacteristicsResponse{
		Category:        category,
		Characteristics: characteristics,
	}, nil
}

// GetCategoryStats получает статистику по категории
func (s *Service) GetCategoryStats(ctx context.Context, categoryName string) (*GetCategoryStatsResponse, error) {
	// Получаем статистику по продуктам
	productStatsQuery := `
		SELECT 
			category,
			COUNT(*) as product_count,
			AVG(recommend_price) as avg_price,
			MIN(recommend_price) as min_price,
			MAX(recommend_price) as max_price
		FROM products 
		WHERE category = ? AND recommend_price IS NOT NULL
		GROUP BY category`

	var stats CategoryStats
	var avgPrice, minPrice, maxPrice sql.NullFloat64

	err := s.db.QueryRow(productStatsQuery, categoryName).Scan(
		&stats.CategoryName,
		&stats.ProductCount,
		&avgPrice,
		&minPrice,
		&maxPrice,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("категория не найдена")
		}
		return nil, fmt.Errorf("ошибка получения статистики: %v", err)
	}

	// Генерируем ID для категории
	if strings.HasPrefix(strings.ToLower(stats.CategoryName), "wb:") {
		stats.CategoryID = generateCategoryID(stats.CategoryName)
	}

	if avgPrice.Valid {
		stats.AvgPrice = avgPrice.Float64
	}
	if minPrice.Valid {
		stats.MinPrice = minPrice.Float64
	}
	if maxPrice.Valid {
		stats.MaxPrice = maxPrice.Float64
	}

	// Получаем статистику по офферам
	offerStatsQuery := `
		SELECT COUNT(*) as offer_count
		FROM offers o
		JOIN products p ON o.product_id = p.id
		WHERE p.category = ?`

	var offerCount int64
	err = s.db.QueryRow(offerStatsQuery, categoryName).Scan(&offerCount)
	if err != nil {
		// Если ошибка, не критично - просто не добавляем статистику по офферам
		offerCount = 0
	}

	stats.OfferCount = offerCount

	return &GetCategoryStatsResponse{
		Stats: stats,
	}, nil
}

// getCharacteristicsFromOpenSearch получает характеристики из OpenSearch
func (s *Service) getCharacteristicsFromOpenSearch(ctx context.Context, categoryName string) ([]Characteristic, error) {
	// Получаем характеристики из OpenSearch через агрегацию
	opensearchURL := "http://localhost:9200"
	
	// Строим запрос для получения характеристик категории
	query := map[string]interface{}{
		"size": 1, // Получаем хотя бы один результат для извлечения названия категории
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"category": categoryName,
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"characteristics": map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "characteristics",
				},
				"aggs": map[string]interface{}{
					"characteristic_names": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "characteristics.name",
							"size":  1000,
						},
						"aggs": map[string]interface{}{
							"values": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "characteristics.value.keyword",
									"size":  100,
								},
							},
						},
					},
				},
			},
		},
	}
	
	// Выполняем запрос к OpenSearch
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %v", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", opensearchURL+"/products/_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка OpenSearch: статус %d", resp.StatusCode)
	}
	
	// Парсим ответ
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %v", err)
	}
	
	// Извлекаем характеристики
	characteristics := []Characteristic{}
	
	if aggs, ok := searchResponse["aggregations"].(map[string]interface{}); ok {
		if charAggs, ok := aggs["characteristics"].(map[string]interface{}); ok {
			if charNames, ok := charAggs["characteristic_names"].(map[string]interface{}); ok {
				if buckets, ok := charNames["buckets"].([]interface{}); ok {
					for _, bucket := range buckets {
						if bucketMap, ok := bucket.(map[string]interface{}); ok {
							charName := bucketMap["key"].(string)
							docCount := int64(bucketMap["doc_count"].(float64))
							
							// Определяем тип характеристики на основе значений
							charType := s.determineCharacteristicType(bucketMap)
							fmt.Printf("DEBUG: Characteristic '%s' determined type: %s\n", charName, charType)
							
							// Извлекаем опции для select/multiselect
							options := s.extractCharacteristicOptions(bucketMap)
							
							// Определяем единицы измерения
							unit := s.extractCharacteristicUnit(charName)
							
							characteristic := Characteristic{
								ID:         int64(len(characteristics) + 1),
								CategoryID: 0, // Будет установлен позже
								Name:       charName,
								Type:       charType,
								Required:   false, // По умолчанию не обязательная
								Options:    options,
								Unit:       unit,
								Description: fmt.Sprintf("Характеристика найдена в %d продуктах", docCount),
								CreatedAt:  time.Now(),
								UpdatedAt:  time.Now(),
							}
							
							characteristics = append(characteristics, characteristic)
						}
					}
				}
			}
		}
	}
	
	return characteristics, nil
}

// GetCategoryCharacteristicsByID получает характеристики для категории по ID из OpenSearch
func (s *Service) GetCategoryCharacteristicsByID(ctx context.Context, categoryID int64) (*GetCategoryCharacteristicsResponse, error) {
	fmt.Printf("DEBUG: GetCategoryCharacteristicsByID called with categoryID: %d\n", categoryID)
	
	// Получаем характеристики из OpenSearch через агрегацию по category_id
	opensearchURL := "http://localhost:9200"
	
	// Строим запрос для получения товаров категории
	query := map[string]interface{}{
		"size": 100, // Получаем больше товаров для извлечения характеристик
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"category_id": categoryID, // Используем как число
						},
					},
				},
			},
		},
	}
	
	// Выполняем запрос к OpenSearch
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %v", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", opensearchURL+"/products/_search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка OpenSearch: статус %d", resp.StatusCode)
	}
	
	// Парсим ответ
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %v", err)
	}
	
	// Проверяем, есть ли результаты и получаем название категории
	var categoryName string
	if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
		fmt.Printf("DEBUG: Hits found: %+v\n", hits)
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if count, ok := total["value"].(float64); ok && count == 0 {
				return nil, fmt.Errorf("категория не найдена")
			}
		}
		
		// Получаем информацию о категории из первого результата
		if hitList, ok := hits["hits"].([]interface{}); ok && len(hitList) > 0 {
			fmt.Printf("DEBUG: Found %d hits\n", len(hitList))
			if firstHit, ok := hitList[0].(map[string]interface{}); ok {
				if source, ok := firstHit["_source"].(map[string]interface{}); ok {
					fmt.Printf("DEBUG: Source fields: %+v\n", source)
					// Ищем название категории в поле category_name (или entity если есть)
					if name, ok := source["category_name"].(string); ok && name != "" {
						categoryName = name
						fmt.Printf("DEBUG: Using category_name: %s\n", name)
					} else if name, ok := source["entity"].(string); ok && name != "" {
						categoryName = name
						fmt.Printf("DEBUG: Using entity: %s\n", name)
					} else {
						fmt.Printf("DEBUG: No valid category name found\n")
					}
				}
			}
		} else {
			fmt.Printf("DEBUG: No hits found\n")
		}
	} else {
		fmt.Printf("DEBUG: No hits in response\n")
	}
	
    if categoryName == "" {
        categoryName = fmt.Sprintf("Категория %d", categoryID)
        fmt.Printf("DEBUG: Using fallback category name: %s\n", categoryName)
    } else {
        fmt.Printf("DEBUG: Using found category name: %s\n", categoryName)
    }

    // Переопределяем название категории значением из MySQL (products.category_name), если доступно
    // Требование: брать человекочитаемое название из колонки category_name, а не из category
    var dbCategoryName sql.NullString
    err = s.db.QueryRow(
        `SELECT category_name FROM products WHERE category_id = ? AND category_name IS NOT NULL AND category_name != '' LIMIT 1`,
        categoryID,
    ).Scan(&dbCategoryName)
    if err == nil && dbCategoryName.Valid {
        if cn := strings.TrimSpace(dbCategoryName.String); cn != "" {
            fmt.Printf("DEBUG: Overriding category name from DB category_name: %s\n", cn)
            categoryName = cn
        }
    } else if err != nil && err != sql.ErrNoRows {
        fmt.Printf("DEBUG: Error reading category_name from DB: %v\n", err)
    }
	
	// Извлекаем характеристики из всех найденных товаров
	characteristics := []Characteristic{}
	charMap := make(map[string]map[string]int) // charName -> value -> count
	
	if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
		if hitList, ok := hits["hits"].([]interface{}); ok {
			fmt.Printf("DEBUG: Found %d products\n", len(hitList))
			for _, hit := range hitList {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if source, ok := hitMap["_source"].(map[string]interface{}); ok {
						// Сначала пробуем получить характеристики из поля characteristics
						if charList, ok := source["characteristics"].([]interface{}); ok && len(charList) > 0 {
							for _, char := range charList {
								if charMapItem, ok := char.(map[string]interface{}); ok {
									if name, ok := charMapItem["name"].(string); ok {
										if value, ok := charMapItem["value"].(string); ok {
											if charMap[name] == nil {
												charMap[name] = make(map[string]int)
											}
											charMap[name][value]++
										}
									}
								}
							}
						} else {
							// Если characteristics пустое, парсим из description
							if description, ok := source["description"].(string); ok && description != "" {
								s.parseCharacteristicsFromDescription(description, charMap)
							}
						}
					}
				}
			}
		}
	}
	
	// Преобразуем собранные характеристики в структуры
	i := 1
	for charName, values := range charMap {
		// Определяем тип характеристики на основе значений
		charType := s.determineCharacteristicTypeFromValues(values)
		fmt.Printf("DEBUG: Characteristic '%s' determined type: %s\n", charName, charType)
		
		// Извлекаем опции и определяем единицы измерения
		options := make([]string, 0, len(values))
		var unit string
		
		// Для числовых характеристик извлекаем единицы измерения
		if charType == "number" {
			for value := range values {
				options = append(options, value)
				if unit == "" {
					_, extractedUnit := s.extractNumericValue(value)
					if extractedUnit != "" {
						unit = extractedUnit
					}
				}
			}
		} else {
			// Для нечисловых характеристик просто добавляем все значения
			for value := range values {
				options = append(options, value)
			}
			// Определяем единицы измерения по названию характеристики
			unit = s.extractCharacteristicUnit(charName)
		}
		
		characteristic := Characteristic{
			ID:         int64(i),
			CategoryID: categoryID,
			Name:       charName,
			Type:       charType,
			Required:   false,
			Options:    options,
			Unit:       unit,
			Description: fmt.Sprintf("Характеристика найдена в %d уникальных значениях", len(values)),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		characteristics = append(characteristics, characteristic)
		i++
	}
	
	category := Category{
		ID:          categoryID,
		Name:        categoryName,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProductCount: 0, // Можно добавить подсчет из hits.total
	}
	
	return &GetCategoryCharacteristicsResponse{
		Category:        category,
		Characteristics: characteristics,
	}, nil
}

// determineCharacteristicType определяет тип характеристики на основе значений
func (s *Service) determineCharacteristicType(bucket map[string]interface{}) string {
	if values, ok := bucket["values"].(map[string]interface{}); ok {
		if valueBuckets, ok := values["buckets"].([]interface{}); ok {
			if len(valueBuckets) == 0 {
				return "text"
			}
			
			// Проверяем первые несколько значений
			for i, valueBucket := range valueBuckets {
				if i >= 5 { // Проверяем только первые 5 значений
					break
				}
				
				if valueMap, ok := valueBucket.(map[string]interface{}); ok {
					value := valueMap["key"].(string)
					
					// Проверяем, является ли значение числом (включая с единицами измерения)
					if s.isNumericValue(value) {
						return "number"
					}
					
					// Проверяем булевые значения
					if value == "да" || value == "нет" || value == "true" || value == "false" {
						return "boolean"
					}
					
					// Проверяем multiselect (значения через запятую, но не десятичные числа)
					if strings.Contains(value, ",") {
						// Проверяем, не является ли это десятичным числом
						if !strings.Contains(value, ".") && !strings.Contains(value, " ") {
							// Если содержит запятую, но нет точки и пробелов - это десятичное число
							continue
						}
						// Если содержит запятую и пробелы - это multiselect
						if strings.Contains(value, " ") {
							return "multiselect"
						}
					}
				}
			}
			
			// Если много уникальных значений, это select
			if len(valueBuckets) > 10 {
				return "select"
			}
			
			return "text"
		}
	}
	
	return "text"
}

// isNumericValue проверяет, является ли значение числовым (включая с единицами измерения)
func (s *Service) isNumericValue(value string) bool {
	// Убираем пробелы
	value = strings.TrimSpace(value)
	fmt.Printf("DEBUG: Checking if '%s' is numeric\n", value)
	
	// Проверяем, является ли значение чистым числом
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		fmt.Printf("DEBUG: '%s' is pure number\n", value)
		return true
	}
	
	// Проверяем числовые значения с единицами измерения
	// Паттерны: "10 мм", "25.5 кг", "100 Вт", "220 В", "50 Гц", "1.5 А", "12 В", "24 В"
	numericPatterns := []string{
		" мм", " кг", " г", " мг", " т", " л", " мл", " м", " см", " км",
		" В", " А", " Вт", " Гц", " МГц", " ГГц", " МГц", " кГц",
		" °C", " °F", " %", " шт", " пк", " уп", " р", " руб", " $", " €",
		" дБ", " ГБ", " МБ", " КБ", " ТБ", " МП", " пикс", " пикселей",
		" об/мин", " об/с", " м/с", " км/ч", " м/ч", " м²", " м³", " см²", " см³",
	}
	
	for _, pattern := range numericPatterns {
		if strings.HasSuffix(value, pattern) {
			// Извлекаем числовую часть
			numericPart := strings.TrimSuffix(value, pattern)
			if _, err := strconv.ParseFloat(numericPart, 64); err == nil {
				return true
			}
		}
	}
	
	return false
}

// extractNumericValue извлекает числовое значение из строки с единицами измерения
func (s *Service) extractNumericValue(value string) (float64, string) {
	// Убираем пробелы
	value = strings.TrimSpace(value)
	
	// Проверяем, является ли значение чистым числом
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, ""
	}
	
	// Проверяем числовые значения с единицами измерения
	numericPatterns := []string{
		" мм", " кг", " г", " мг", " т", " л", " мл", " м", " см", " км",
		" В", " А", " Вт", " Гц", " МГц", " ГГц", " МГц", " кГц",
		" °C", " °F", " %", " шт", " пк", " уп", " р", " руб", " $", " €",
		" дБ", " ГБ", " МБ", " КБ", " ТБ", " МП", " пикс", " пикселей",
		" об/мин", " об/с", " м/с", " км/ч", " м/ч", " м²", " м³", " см²", " см³",
	}
	
	for _, pattern := range numericPatterns {
		if strings.HasSuffix(value, pattern) {
			// Извлекаем числовую часть
			numericPart := strings.TrimSuffix(value, pattern)
			if num, err := strconv.ParseFloat(numericPart, 64); err == nil {
				return num, strings.TrimSpace(pattern)
			}
		}
	}
	
	return 0, ""
}

// determineCharacteristicTypeFromValues определяет тип характеристики на основе значений
func (s *Service) determineCharacteristicTypeFromValues(values map[string]int) string {
	if len(values) == 0 {
		return "text"
	}
	
	// Проверяем первые несколько значений
	count := 0
	numericCount := 0
	booleanCount := 0
	multiselectCount := 0
	
	for value := range values {
		if count >= 10 { // Проверяем больше значений для лучшей точности
			break
		}
		
		// Проверяем, является ли значение числом (включая с единицами измерения)
		if s.isNumericValue(value) {
			numericCount++
		}
		
		// Проверяем булевые значения
		if value == "да" || value == "нет" || value == "true" || value == "false" {
			booleanCount++
		}
		
		// Проверяем multiselect (значения через запятую, но не десятичные числа)
		if strings.Contains(value, ",") {
			// Проверяем, не является ли это десятичным числом
			if !strings.Contains(value, ".") && !strings.Contains(value, " ") {
				// Если содержит запятую, но нет точки и пробелов - это десятичное число
				continue
			}
			// Если содержит запятую и пробелы - это multiselect
			if strings.Contains(value, " ") {
				multiselectCount++
			}
		}
		
		count++
	}
	
	// Определяем тип на основе преобладающих значений
	if booleanCount > count/2 {
		return "boolean"
	}
	
	if multiselectCount > count/3 {
		return "multiselect"
	}
	
	if numericCount > count/2 {
		return "number"
	}
	
	// Если много уникальных значений, это select
	if len(values) > 10 {
		return "select"
	}
	
	return "text"
}

// extractCharacteristicOptions извлекает опции для select/multiselect
func (s *Service) extractCharacteristicOptions(bucket map[string]interface{}) []string {
	options := []string{}
	optionSet := make(map[string]bool) // Для избежания дубликатов
	
	if values, ok := bucket["values"].(map[string]interface{}); ok {
		if valueBuckets, ok := values["buckets"].([]interface{}); ok {
			for _, valueBucket := range valueBuckets {
				if valueMap, ok := valueBucket.(map[string]interface{}); ok {
					value := valueMap["key"].(string)
					
					// Если значение содержит запятую (multiselect), разбиваем его
					if strings.Contains(value, ",") {
						parts := strings.Split(value, ",")
						for _, part := range parts {
							part = strings.TrimSpace(part)
							if part != "" && !optionSet[part] {
								options = append(options, part)
								optionSet[part] = true
							}
						}
					} else {
						// Обычное значение
						if !optionSet[value] {
							options = append(options, value)
							optionSet[value] = true
						}
					}
				}
			}
		}
	}
	
	return options
}

// extractCharacteristicUnit извлекает единицы измерения из названия характеристики
func (s *Service) extractCharacteristicUnit(charName string) string {
	// Простое извлечение единиц измерения из названия
	if strings.Contains(strings.ToLower(charName), "мощность") && strings.Contains(strings.ToLower(charName), "вт") {
		return "Вт"
	}
	if strings.Contains(strings.ToLower(charName), "напряжение") && strings.Contains(strings.ToLower(charName), "в") {
		return "В"
	}
	if strings.Contains(strings.ToLower(charName), "емкость") && strings.Contains(strings.ToLower(charName), "ач") {
		return "Ач"
	}
	if strings.Contains(strings.ToLower(charName), "диаметр") && strings.Contains(strings.ToLower(charName), "мм") {
		return "мм"
	}
	if strings.Contains(strings.ToLower(charName), "длина") && strings.Contains(strings.ToLower(charName), "см") {
		return "см"
	}
	if strings.Contains(strings.ToLower(charName), "вес") && strings.Contains(strings.ToLower(charName), "кг") {
		return "кг"
	}
	if strings.Contains(strings.ToLower(charName), "обороты") && strings.Contains(strings.ToLower(charName), "об/мин") {
		return "об/мин"
	}
	
	return ""
}

// generateCategoryID генерирует ID для WB категории
func generateCategoryID(categoryName string) int64 {
	if !strings.HasPrefix(strings.ToLower(categoryName), "wb:") {
		return 0
	}
	
	// Простая генерация ID на основе хэша названия
	// В реальной системе здесь должен быть тот же алгоритм, что и в utils
	return int64(len(categoryName) * 1000) // Упрощенная версия
}

// parseCharacteristicsFromDescription парсит характеристики из текстового описания
func (s *Service) parseCharacteristicsFromDescription(description string, charMap map[string]map[string]int) {
	// Разбиваем описание по точкам с запятой
	parts := strings.Split(description, ";")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Ищем двоеточие для разделения названия и значения
		if strings.Contains(part, ":") {
			colonIndex := strings.Index(part, ":")
			name := strings.TrimSpace(part[:colonIndex])
			value := strings.TrimSpace(part[colonIndex+1:])
			
			if name != "" && value != "" {
				if charMap[name] == nil {
					charMap[name] = make(map[string]int)
				}
				charMap[name][value]++
			}
		}
	}
}
