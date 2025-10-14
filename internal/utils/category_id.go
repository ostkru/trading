package utils

import (
	"hash/crc32"
	"strconv"
	"strings"
)

// GenerateCategoryID создает уникальный числовой ID категории из названия
// Работает только для категорий с префиксом "wb:" в формате "wb: 1318 - сварочные аппараты"
// Использует алгоритм CRC32 для генерации стабильного хэша
func GenerateCategoryID(categoryName string) int64 {
	// Проверяем, что категория имеет префикс "wb:"
	if !isWBCategory(categoryName) {
		// Если не WB категория, возвращаем 0 (не генерируем ID)
		return 0
	}
	
	// Извлекаем части из формата "wb: 1318 - сварочные аппараты"
	subjectId, entity := parseWBCategory(categoryName)
	if subjectId == "" || entity == "" {
		return 0
	}
	
	// Используем алгоритм из PHP: subjectId + "_" + entity
	combined := subjectId + "_" + entity
	
	// Создаем хэш CRC32
	hash := crc32.ChecksumIEEE([]byte(combined))
	
	// Возвращаем абсолютное значение как int64
	return int64(hash)
}

// isWBCategory проверяет, является ли категория WB категорией
func isWBCategory(categoryName string) bool {
	normalized := strings.TrimSpace(categoryName)
	return strings.HasPrefix(strings.ToLower(normalized), "wb:")
}

// parseWBCategory извлекает subjectId и entity из WB категории
// Формат: "wb: 1318 - сварочные аппараты" -> subjectId="1318", entity="сварочные аппараты"
func parseWBCategory(categoryName string) (subjectId, entity string) {
	// Убираем префикс "wb:" и нормализуем
	normalized := strings.TrimSpace(categoryName)
	if strings.HasPrefix(strings.ToLower(normalized), "wb:") {
		normalized = strings.TrimSpace(normalized[3:]) // Убираем "wb:"
	}
	
	// Ищем разделитель " - "
	parts := strings.Split(normalized, " - ")
	if len(parts) != 2 {
		return "", ""
	}
	
	subjectId = strings.TrimSpace(parts[0])
	entity = strings.TrimSpace(parts[1])
	
	return subjectId, entity
}

// normalizeCategoryName нормализует название категории для стабильного хэширования
func normalizeCategoryName(categoryName string) string {
	// Убираем лишние пробелы
	normalized := strings.TrimSpace(categoryName)
	
	// Приводим к нижнему регистру для консистентности
	normalized = strings.ToLower(normalized)
	
	// Убираем лишние пробелы между словами
	normalized = strings.Join(strings.Fields(normalized), " ")
	
	return normalized
}

// GenerateCategoryIDFromParts создает ID категории из частей (например, subjectId + entity)
// Это для совместимости с существующим алгоритмом из PHP
func GenerateCategoryIDFromParts(subjectId string, entity string) int64 {
	// Объединяем части с разделителем
	combined := subjectId + "_" + entity
	
	// Создаем хэш CRC32
	hash := crc32.ChecksumIEEE([]byte(combined))
	
	// Возвращаем абсолютное значение как int64
	return int64(hash)
}

// GenerateCategoryNameFromParts создает название категории из частей
func GenerateCategoryNameFromParts(subjectId string, entity string) string {
	return subjectId + " - " + entity
}

// ValidateCategoryID проверяет, что ID категории является валидным
func ValidateCategoryID(categoryID int64) bool {
	// ID должен быть положительным числом
	return categoryID > 0
}

// ParseWBCategory извлекает subjectId и entity из WB категории (публичная версия)
func ParseWBCategory(categoryName string) (subjectId, entity string) {
	return parseWBCategory(categoryName)
}

// IsWBCategory проверяет, является ли категория WB категорией (публичная версия)
func IsWBCategory(categoryName string) bool {
	return isWBCategory(categoryName)
}

// GetCategoryIDFromString пытается извлечь ID категории из строки
func GetCategoryIDFromString(categoryStr string) (int64, error) {
	// Пытаемся распарсить как число
	if id, err := strconv.ParseInt(categoryStr, 10, 64); err == nil {
		return id, nil
	}
	
	// Если не число, генерируем ID из строки
	return GenerateCategoryID(categoryStr), nil
}
