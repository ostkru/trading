package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ParserService struct {
	dataPath string
}

type CSVProduct struct {
	Name         string  `json:"name"`
	VendorArticle string  `json:"vendor_article"`
	Price        float64 `json:"price"`
	Brand        string  `json:"brand"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
}

type ParserResponse struct {
	Products []CSVProduct `json:"products"`
	Total    int          `json:"total"`
	Offset   int          `json:"offset"`
	Limit    int          `json:"limit"`
}

func NewParserService(dataPath string) *ParserService {
	return &ParserService{dataPath: dataPath}
}

func (s *ParserService) ParseCSV(filePath string, limit, offset int, excludeFields []string) (*ParserResponse, error) {
	// Полный путь к файлу
	fullPath := fmt.Sprintf("%s/%s", s.dataPath, filePath)
	
	log.Printf("Parsing CSV file: %s", fullPath)
	
	// Открываем файл
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Создаем CSV reader
	reader := csv.NewReader(file)
	reader.Comma = ';' // Разделитель точка с запятой
	reader.LazyQuotes = true

	// Читаем заголовки
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers: %w", err)
	}

	// Создаем карту исключенных полей для быстрого поиска
	excludeMap := make(map[string]bool)
	for _, field := range excludeFields {
		excludeMap[field] = true
	}

	// Находим индексы нужных полей
	fieldIndexes := make(map[string]int)
	for i, header := range headers {
		if !excludeMap[header] {
			fieldIndexes[header] = i
		}
	}

	// Пропускаем строки до offset
	for i := 0; i < offset; i++ {
		_, err := reader.Read()
		if err != nil {
			break // Достигнут конец файла
		}
	}

	// Читаем данные с лимитом
	var products []CSVProduct
	count := 0
	for count < limit {
		record, err := reader.Read()
		if err != nil {
			break // Достигнут конец файла
		}

		product := s.parseRecord(record, fieldIndexes)
		if product != nil {
			products = append(products, *product)
			count++
		}
	}

	// Подсчитываем общее количество строк в файле
	total, err := s.countTotalLines(fullPath)
	if err != nil {
		log.Printf("Warning: failed to count total lines: %v", err)
		total = offset + len(products)
	}

	log.Printf("Parsed %d products (offset: %d, limit: %d, total: %d)", 
		len(products), offset, limit, total)

	return &ParserResponse{
		Products: products,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}, nil
}

func (s *ParserService) parseRecord(record []string, fieldIndexes map[string]int) *CSVProduct {
	product := &CSVProduct{}

	// Извлекаем поля по индексам
	if idx, ok := fieldIndexes["Название"]; ok && idx < len(record) {
		product.Name = strings.TrimSpace(record[idx])
	}
	if idx, ok := fieldIndexes["Артикул"]; ok && idx < len(record) {
		product.VendorArticle = strings.TrimSpace(record[idx])
	}
	if idx, ok := fieldIndexes["Цена"]; ok && idx < len(record) {
		if price, err := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); err == nil {
			product.Price = price
		}
	}
	if idx, ok := fieldIndexes["Бренд"]; ok && idx < len(record) {
		product.Brand = strings.TrimSpace(record[idx])
	}
	if idx, ok := fieldIndexes["Категория"]; ok && idx < len(record) {
		product.Category = strings.TrimSpace(record[idx])
	}
	if idx, ok := fieldIndexes["Описание"]; ok && idx < len(record) {
		product.Description = strings.TrimSpace(record[idx])
	}

	// Проверяем обязательные поля
	if product.Name == "" || product.VendorArticle == "" || product.Brand == "" {
		return nil // Пропускаем товары без обязательных полей
	}

	return product
}

func (s *ParserService) countTotalLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	
	count := 0
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}
		count++
	}

	return count - 1, nil // Вычитаем заголовок
} 