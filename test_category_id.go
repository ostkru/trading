package main

import (
	"fmt"
	"portaldata-api/internal/utils"
)

func main() {
	fmt.Println("=== Тестирование алгоритма генерации category_id ===")
	
	// Тестовые категории
	testCategories := []string{
		"2384 - сверла",
		"2197 - шуруповерты", 
		"1168 - шлифовальные машины",
		"2294 - отвертки",
		"4473 - насосы погружные",
		"2382 - биты для шуруповерта",
		"519 - кабели",
		"1318 - сварочные аппараты",
		"4009 - диски пильные",
		"2066 - молотки",
		"3830 - буры",
		"4010 - диски алмазные",
		"1166 - перфораторы",
		"770 - домкраты",
		"655 - наборы инструментов",
		"4099 - круги шлифовальные",
		"1535 - лески для триммеров",
		"2223 - триммеры садовые",
		"2060 - ключи гаечные",
		"1708 - смесители",
	}
	
	fmt.Println("\nГенерация category_id для тестовых категорий:")
	fmt.Println("Категория -> Category ID")
	fmt.Println("=" + string(make([]byte, 50)))
	
	for _, category := range testCategories {
		categoryID := utils.GenerateCategoryID(category)
		fmt.Printf("%-30s -> %d\n", category, categoryID)
	}
	
	// Тестирование стабильности
	fmt.Println("\n=== Тест стабильности ===")
	testCategory := "1318 - сварочные аппараты"
	
	for i := 0; i < 5; i++ {
		categoryID := utils.GenerateCategoryID(testCategory)
		fmt.Printf("Попытка %d: %s -> %d\n", i+1, testCategory, categoryID)
	}
	
	// Тестирование с разными регистрами и пробелами
	fmt.Println("\n=== Тест нормализации ===")
	variations := []string{
		"1318 - сварочные аппараты",
		"1318 - СВАРОЧНЫЕ АППАРАТЫ",
		"  1318 - сварочные аппараты  ",
		"1318   -   сварочные   аппараты",
		"1318-сварочные-аппараты",
	}
	
	for _, variation := range variations {
		categoryID := utils.GenerateCategoryID(variation)
		fmt.Printf("%-35s -> %d\n", variation, categoryID)
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
	
	fmt.Println("\n=== Тест завершен ===")
}

