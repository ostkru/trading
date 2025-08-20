package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ClassificationResult результат классификации от API catformat
type ClassificationResult struct {
	Status          string      `json:"status,omitempty"`
	Request         string      `json:"request,omitempty"`
	FoundCategory   string      `json:"found_category,omitempty"`
	FoundCategoryID interface{} `json:"found_category_id,omitempty"`
	FoundBrand      string      `json:"found_brand,omitempty"`
	FoundBrandID    interface{} `json:"brand_id,omitempty"`
	Accuracy        float64     `json:"accuracy,omitempty"`
	BrandAccuracy   float64     `json:"brand_accuracy,omitempty"`
	UserCategory    string      `json:"user_category,omitempty"`
	User            string      `json:"user,omitempty"`
	Error           string      `json:"error,omitempty"`
}

func main() {
	fmt.Println("🔍 Тестируем парсинг JSON ответа API...")

	// Тестируем с разными продуктами
	testProducts := []struct {
		name     string
		category string
		brand    string
	}{
		{"Телевизор Samsung 55QN90B Neo QLED", "Телевизоры", "Samsung"},
		{"iPhone 15 Pro Max", "Смартфоны", "Apple"},
		{"Dell Latitude 5520", "Ноутбуки", "Dell"},
	}

	for i, product := range testProducts {
		fmt.Printf("\n=== Тест %d: %s ===\n", i+1, product.name)

		// Вызываем API
		apiURL := "https://api.ostk.ru/products/catformat.php"
		params := url.Values{}
		params.Set("product_name", product.name)
		params.Set("user", "1")
		params.Set("user_category", product.category)
		params.Set("brand", product.brand)
		params.Set("findbrand", "1")
		params.Set("findcat", "1")

		resp, err := http.Get(apiURL + "?" + params.Encode())
		if err != nil {
			fmt.Printf("❌ Ошибка HTTP запроса: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ Ошибка чтения ответа: %v\n", err)
			continue
		}

		fmt.Printf("📡 Ответ API: %s\n", string(body))

		// Парсим JSON
		var result ClassificationResult
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("❌ Ошибка парсинга JSON: %v\n", err)
			continue
		}

		fmt.Printf("📊 Результат парсинга:\n")
		fmt.Printf("  Status: %s\n", result.Status)
		fmt.Printf("  FoundCategory: %s\n", result.FoundCategory)
		fmt.Printf("  FoundCategoryID: %v (тип: %T)\n", result.FoundCategoryID, result.FoundCategoryID)
		fmt.Printf("  FoundBrand: %s\n", result.FoundBrand)
		fmt.Printf("  FoundBrandID: %v (тип: %T)\n", result.FoundBrandID, result.FoundBrandID)
		fmt.Printf("  Accuracy: %f\n", result.Accuracy)
		fmt.Printf("  BrandAccuracy: %f\n", result.BrandAccuracy)

		// Проверяем типы данных
		if result.Status == "found" {
			fmt.Printf("✅ Классификация успешна\n")

			// Проверяем category_id
			if result.FoundCategoryID != nil {
				switch v := result.FoundCategoryID.(type) {
				case string:
					fmt.Printf("⚠️  FoundCategoryID - строка: %s\n", v)
				case float64:
					fmt.Printf("✅ FoundCategoryID - число: %f\n", v)
				default:
					fmt.Printf("❓ FoundCategoryID - неожиданный тип: %T\n", v)
				}
			}

			// Проверяем brand_id
			if result.FoundBrandID != nil {
				switch v := result.FoundBrandID.(type) {
				case string:
					fmt.Printf("⚠️  FoundBrandID - строка: %s\n", v)
				case float64:
					fmt.Printf("✅ FoundBrandID - число: %f\n", v)
				default:
					fmt.Printf("❓ FoundBrandID - неожиданный тип: %T\n", v)
				}
			}
		} else {
			fmt.Printf("❌ Классификация не удалась: %s\n", result.Status)
		}
	}
}
