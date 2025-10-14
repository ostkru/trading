package main

import (
	"fmt"
	"portaldata-api/internal/utils"
)

func main() {
	fmt.Println("=== Тестирование алгоритма генерации category_id для WB категорий ===")
	
	// Тестовые WB категории
	wbCategories := []string{
		"wb: 1318 - сварочные аппараты",
		"wb: 2384 - сверла",
		"wb: 2197 - шуруповерты",
		"wb: 1168 - шлифовальные машины",
		"wb: 2294 - отвертки",
		"wb: 4473 - насосы погружные",
		"wb: 2382 - биты для шуруповерта",
		"wb: 519 - кабели",
		"wb: 4009 - диски пильные",
		"wb: 2066 - молотки",
	}
	
	// Тестовые НЕ-WB категории (не должны генерировать ID)
	nonWBCategories := []string{
		"1318 - сварочные аппараты",
		"Сварочные аппараты",
		"Инструменты",
		"Электроинструмент",
		"Ручной инструмент",
	}
	
	fmt.Println("\n=== WB категории (должны генерировать ID) ===")
	fmt.Println("Категория -> Category ID")
	fmt.Println("=" + string(make([]byte, 50)))
	
	for _, category := range wbCategories {
		categoryID := utils.GenerateCategoryID(category)
		fmt.Printf("%-35s -> %d\n", category, categoryID)
	}
	
	fmt.Println("\n=== НЕ-WB категории (НЕ должны генерировать ID) ===")
	fmt.Println("Категория -> Category ID")
	fmt.Println("=" + string(make([]byte, 50)))
	
	for _, category := range nonWBCategories {
		categoryID := utils.GenerateCategoryID(category)
		fmt.Printf("%-35s -> %d\n", category, categoryID)
	}
	
	// Тестирование парсинга WB категорий
	fmt.Println("\n=== Тест парсинга WB категорий ===")
	testWBCategories := []string{
		"wb: 1318 - сварочные аппараты",
		"WB: 2384 - сверла",
		"wb:2197-шуруповерты",
		"  wb: 1168 - шлифовальные машины  ",
		"wb: 2294   -   отвертки",
	}
	
	for _, category := range testWBCategories {
		subjectId, entity := utils.ParseWBCategory(category)
		categoryID := utils.GenerateCategoryID(category)
		fmt.Printf("'%s' -> subjectId='%s', entity='%s', ID=%d\n", 
			category, subjectId, entity, categoryID)
	}
	
	// Тестирование стабильности
	fmt.Println("\n=== Тест стабильности ===")
	testCategory := "wb: 1318 - сварочные аппараты"
	
	for i := 0; i < 5; i++ {
		categoryID := utils.GenerateCategoryID(testCategory)
		fmt.Printf("Попытка %d: %s -> %d\n", i+1, testCategory, categoryID)
	}
	
	// Тестирование с разными регистрами и пробелами
	fmt.Println("\n=== Тест нормализации WB категорий ===")
	variations := []string{
		"wb: 1318 - сварочные аппараты",
		"WB: 1318 - СВАРОЧНЫЕ АППАРАТЫ",
		"  wb: 1318 - сварочные аппараты  ",
		"wb:1318-сварочные-аппараты",
		"wb:  1318   -   сварочные   аппараты",
	}
	
	for _, variation := range variations {
		categoryID := utils.GenerateCategoryID(variation)
		fmt.Printf("%-40s -> %d\n", variation, categoryID)
	}
	
	// Тестирование GenerateCategoryIDFromParts
	fmt.Println("\n=== Тест GenerateCategoryIDFromParts ===")
	testParts := []struct {
		subjectId string
		entity    string
	}{
		{"1318", "сварочные аппараты"},
		{"2384", "сверла"},
		{"2197", "шуруповерты"},
		{"1168", "шлифовальные машины"},
	}
	
	for _, parts := range testParts {
		categoryID := utils.GenerateCategoryIDFromParts(parts.subjectId, parts.entity)
		categoryName := utils.GenerateCategoryNameFromParts(parts.subjectId, parts.entity)
		fmt.Printf("%s + %s -> %s (ID: %d)\n", parts.subjectId, parts.entity, categoryName, categoryID)
	}
	
	// Проверка совместимости с PHP
	fmt.Println("\n=== Проверка совместимости с PHP ===")
	phpTestCases := []struct {
		subjectId string
		entity    string
		expected  int64
	}{
		{"1318", "сварочные аппараты", 1496952285},
		{"2384", "сверла", 2457353255},
		{"2197", "шуруповерты", 1319357246},
		{"1168", "шлифовальные машины", 3770456496},
	}
	
	for _, testCase := range phpTestCases {
		actual := utils.GenerateCategoryIDFromParts(testCase.subjectId, testCase.entity)
		status := "✅"
		if actual != testCase.expected {
			status = "❌"
		}
		fmt.Printf("%s %s + %s -> %d (ожидалось: %d)\n", 
			status, testCase.subjectId, testCase.entity, actual, testCase.expected)
	}
	
	fmt.Println("\n=== Тест завершен ===")
}

