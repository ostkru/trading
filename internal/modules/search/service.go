package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Service сервис для работы с поиском
type Service struct {
	opensearchURL string
	httpClient    *http.Client
}

// NewService создает новый сервис поиска
func NewService(opensearchURL string) *Service {
	return &Service{
		opensearchURL: opensearchURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchProducts выполняет поиск продуктов
func (s *Service) SearchProducts(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// Построить поисковый запрос
	query := s.buildSearchQuery(req)
	
	// Выполнить поиск
	searchData := map[string]interface{}{
		"query": query,
		"size":  req.Limit,
		"from":  (req.Page - 1) * req.Limit,
		"sort":  s.buildSort(req.Sort),
	}
	
	// Добавить агрегации если нужны
	if req.Facets {
		searchData["aggs"] = s.buildAggregations()
	}
	
	// Выполнить запрос к OpenSearch
	resp, err := s.executeSearch(ctx, searchData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	
	// Преобразовать ответ
	return s.parseSearchResponse(resp, req)
}

// IndexProduct индексирует продукт в OpenSearch
func (s *Service) IndexProduct(ctx context.Context, product *IndexProductRequest) error {
	// Преобразовать продукт в формат OpenSearch
	doc := s.convertProductToDocument(product)
	
	// Выполнить индексацию
	url := fmt.Sprintf("%s/products/_doc/%d", s.opensearchURL, product.ID)
	return s.indexDocument(ctx, url, doc)
}

// UpdateProduct обновляет продукт в индексе
func (s *Service) UpdateProduct(ctx context.Context, product *IndexProductRequest) error {
	return s.IndexProduct(ctx, product)
}

// DeleteProduct удаляет продукт из индекса
func (s *Service) DeleteProduct(ctx context.Context, productID int64) error {
	url := fmt.Sprintf("%s/products/_doc/%d", s.opensearchURL, productID)
	return s.deleteDocument(ctx, url)
}

// GetIndexStats возвращает статистику индекса
func (s *Service) GetIndexStats(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/products/_stats", s.opensearchURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opensearch error: %s", string(body))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result, nil
}

// buildSearchQuery строит поисковый запрос
func (s *Service) buildSearchQuery(req *SearchRequest) map[string]interface{} {
	must := []map[string]interface{}{}
	should := []map[string]interface{}{}
	filter := []map[string]interface{}{}
	
	// Текстовый поиск
	if req.Query != "" {
		should = append(should, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^3", "brand^2", "category^2", "characteristics.value^1.5"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}
	
	// Фильтры
	if req.CategoryID != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": *req.CategoryID,
			},
		})
	}
	
	if req.BrandID != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"brand_id": *req.BrandID,
			},
		})
	}
	
	if req.Brand != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"brand": *req.Brand,
			},
		})
	}
	
	if req.Category != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"category_name": *req.Category,
			},
		})
	}
	
	// Фильтры по характеристикам
	if len(req.Characteristics) > 0 {
		for name, value := range req.Characteristics {
			filter = append(filter, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "characteristics",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"characteristics.name": name}},
								{"match": map[string]interface{}{"characteristics.value": value}},
							},
						},
					},
				},
			})
		}
	}
	
	// Ценовые фильтры (если есть поле price в индексе)
	if req.PriceMin != nil || req.PriceMax != nil {
		priceRange := map[string]interface{}{}
		if req.PriceMin != nil {
			priceRange["gte"] = *req.PriceMin
		}
		if req.PriceMax != nil {
			priceRange["lte"] = *req.PriceMax
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}
	
	// Построить финальный запрос
	query := map[string]interface{}{
		"bool": map[string]interface{}{},
	}
	
	if len(must) > 0 {
		query["bool"].(map[string]interface{})["must"] = must
	}
	if len(should) > 0 {
		query["bool"].(map[string]interface{})["should"] = should
		query["bool"].(map[string]interface{})["minimum_should_match"] = 1
	}
	if len(filter) > 0 {
		query["bool"].(map[string]interface{})["filter"] = filter
	}
	
	// Если нет условий, возвращаем match_all
	if len(must) == 0 && len(should) == 0 && len(filter) == 0 {
		return map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}
	
	return query
}

// buildSort строит сортировку
func (s *Service) buildSort(sort string) []map[string]interface{} {
	switch sort {
	case "price_asc":
		return []map[string]interface{}{
			{"price": map[string]interface{}{"order": "asc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "price_desc":
		return []map[string]interface{}{
			{"price": map[string]interface{}{"order": "desc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "name_asc":
		return []map[string]interface{}{
			{"name.keyword": map[string]interface{}{"order": "asc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "name_desc":
		return []map[string]interface{}{
			{"name.keyword": map[string]interface{}{"order": "desc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	default: // relevance
		return []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	}
}

// buildAggregations строит агрегации для фильтров
func (s *Service) buildAggregations() map[string]interface{} {
	return map[string]interface{}{
		"brands": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "brand",
				"size":  50,
			},
		},
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category_name",
				"size":  50,
			},
		},
		"price_ranges": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "price",
				"ranges": []map[string]interface{}{
					{"to": 1000},
					{"from": 1000, "to": 5000},
					{"from": 5000, "to": 10000},
					{"from": 10000},
				},
			},
		},
	}
}

// executeSearch выполняет поисковый запрос
func (s *Service) executeSearch(ctx context.Context, searchData map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(searchData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search data: %w", err)
	}
	
	url := fmt.Sprintf("%s/products/_search", s.opensearchURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opensearch error: %s", string(body))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result, nil
}

// parseSearchResponse парсит ответ поиска
func (s *Service) parseSearchResponse(resp map[string]interface{}, req *SearchRequest) (*SearchResponse, error) {
	hits, ok := resp["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	
	total, _ := hits["total"].(map[string]interface{})
	totalValue, _ := total["value"].(float64)
	
	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}
	
	products := make([]SearchProduct, 0, len(hitsList))
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		
		product := s.convertDocumentToProduct(source)
		if score, ok := hitMap["_score"].(float64); ok {
			product.Score = score
		}
		
		products = append(products, product)
	}
	
	response := &SearchResponse{
		Products: products,
		Total:    int64(totalValue),
		Page:     req.Page,
		Limit:    req.Limit,
	}
	
	// Парсить агрегации если есть
	if req.Facets {
		response.Facets = s.parseFacets(resp)
	}
	
	return response, nil
}

// convertProductToDocument преобразует продукт в документ OpenSearch
func (s *Service) convertProductToDocument(product *IndexProductRequest) map[string]interface{} {
	doc := map[string]interface{}{
		"id":           strconv.FormatInt(product.ID, 10),
		"name":         product.Name,
		"vendor_code":  product.VendorArticle,
		"price":       product.RecommendPrice,
		"description":  product.Description,
		"timestamp":    time.Now().Unix(),
	}
	
	if product.Brand != nil {
		doc["brand"] = *product.Brand
	}
	if product.Category != nil {
		doc["category_name"] = *product.Category
	}
	if product.BrandID != nil {
		doc["brand_id"] = *product.BrandID
	}
	if product.CategoryID != nil {
		doc["category_id"] = strconv.FormatInt(*product.CategoryID, 10)
	}
	if product.Barcode != nil {
		doc["barcode"] = *product.Barcode
	}
	
	if len(product.ImageURLs) > 0 {
		doc["image_urls"] = product.ImageURLs
	}
	
	if len(product.Characteristics) > 0 {
		characteristics := make([]map[string]interface{}, 0, len(product.Characteristics))
		for _, char := range product.Characteristics {
			characteristics = append(characteristics, map[string]interface{}{
				"name":  char.Name,
				"value": char.Value,
			})
		}
		doc["characteristics"] = characteristics
	}
	
	return doc
}

// convertDocumentToProduct преобразует документ OpenSearch в продукт
func (s *Service) convertDocumentToProduct(source map[string]interface{}) SearchProduct {
	product := SearchProduct{}
	
	if id, ok := source["id"].(string); ok {
		product.ID = id
	}
	if name, ok := source["name"].(string); ok {
		product.Name = name
	}
	if vendorCode, ok := source["vendor_code"].(string); ok {
		product.VendorCode = vendorCode
	}
	if brand, ok := source["brand"].(string); ok {
		product.Brand = brand
	}
	if category, ok := source["category_name"].(string); ok {
		product.Category = category
	}
	if categoryID, ok := source["category_id"].(string); ok {
		product.CategoryID = categoryID
	}
	if brandID, ok := source["brand_id"].(float64); ok {
		product.BrandID = int64(brandID)
	}
	
	// Обработка изображений
	if imageURLs, ok := source["image_urls"].([]interface{}); ok {
		product.ImageURLs = make([]string, 0, len(imageURLs))
		for _, url := range imageURLs {
			if urlStr, ok := url.(string); ok {
				product.ImageURLs = append(product.ImageURLs, urlStr)
			}
		}
	}
	
	// Обработка характеристик
	if characteristics, ok := source["characteristics"].([]interface{}); ok {
		product.Characteristics = make([]Characteristic, 0, len(characteristics))
		for _, char := range characteristics {
			if charMap, ok := char.(map[string]interface{}); ok {
				characteristic := Characteristic{}
				if name, ok := charMap["name"].(string); ok {
					characteristic.Name = name
				}
				if value, ok := charMap["value"].(string); ok {
					characteristic.Value = value
				}
				product.Characteristics = append(product.Characteristics, characteristic)
			}
		}
	}
	
	return product
}

// parseFacets парсит агрегации
func (s *Service) parseFacets(resp map[string]interface{}) *SearchFacets {
	aggs, ok := resp["aggregations"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	facets := &SearchFacets{}
	
	// Бренды
	if brands, ok := aggs["brands"].(map[string]interface{}); ok {
		if buckets, ok := brands["buckets"].([]interface{}); ok {
			facets.Brands = make([]FacetItem, 0, len(buckets))
			for _, bucket := range buckets {
				if bucketMap, ok := bucket.(map[string]interface{}); ok {
					item := FacetItem{}
					if key, ok := bucketMap["key"].(string); ok {
						item.Value = key
					}
					if count, ok := bucketMap["doc_count"].(float64); ok {
						item.Count = int64(count)
					}
					facets.Brands = append(facets.Brands, item)
				}
			}
		}
	}
	
	// Категории
	if categories, ok := aggs["categories"].(map[string]interface{}); ok {
		if buckets, ok := categories["buckets"].([]interface{}); ok {
			facets.Categories = make([]FacetItem, 0, len(buckets))
			for _, bucket := range buckets {
				if bucketMap, ok := bucket.(map[string]interface{}); ok {
					item := FacetItem{}
					if key, ok := bucketMap["key"].(string); ok {
						item.Value = key
					}
					if count, ok := bucketMap["doc_count"].(float64); ok {
						item.Count = int64(count)
					}
					facets.Categories = append(facets.Categories, item)
				}
			}
		}
	}
	
	return facets
}

// indexDocument индексирует документ
func (s *Service) indexDocument(ctx context.Context, url string, doc map[string]interface{}) error {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch error: %s", string(body))
	}
	
	return nil
}

// deleteDocument удаляет документ
func (s *Service) deleteDocument(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch error: %s", string(body))
	}
	
	return nil
}
