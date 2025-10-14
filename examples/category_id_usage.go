package main

import (
	"encoding/json"
	"fmt"
	"log"
	"portaldata-api/internal/modules/products"
	"portaldata-api/internal/utils"
)

func main() {
	fmt.Println("=== Примеры использования автоматической генерации category_id ===")
	
	// Пример 1: Создание продукта с WB категорией (должен генерировать ID)
	fmt.Println("\n1. Создание продукта с WB категорией:")
	createReq := products.CreateProductRequest{
		Name:           "Сварочный аппарат инверторный",
		VendorArticle:  "SW-200",
		RecommendPrice: 15000.00,
		Brand:          "Ресанта",
		Category:       "wb: 1318 - сварочные аппараты",
		Description:    "Профессиональный сварочный аппарат",
	}
	
	// Автоматическая генерация category_id
	createReq.GenerateCategoryID()
	
	fmt.Printf("   Категория: %s\n", createReq.Category)
	fmt.Printf("   Сгенерированный category_id: %d\n", *createReq.CategoryID)
	
	// Пример 2: Создание продукта с НЕ-WB категорией (НЕ должен генерировать ID)
	fmt.Println("\n2. Создание продукта с НЕ-WB категорией:")
	createReqNonWB := products.CreateProductRequest{
		Name:           "Сверло по металлу",
		VendorArticle:  "DR-10",
		RecommendPrice: 150.00,
		Brand:          "Bosch",
		Category:       "2384 - сверла", // НЕ-WB категория
		Description:    "Сверло по металлу 10мм",
	}
	
	// Автоматическая генерация НЕ произойдет для НЕ-WB категории
	createReqNonWB.GenerateCategoryID()
	
	fmt.Printf("   Категория: %s\n", createReqNonWB.Category)
	fmt.Printf("   Сгенерированный category_id: %v (должен быть nil)\n", createReqNonWB.CategoryID)
	
	// Пример 3: Создание продукта с явным category_id
	fmt.Println("\n3. Создание продукта с явным category_id:")
	createReqWithID := products.CreateProductRequest{
		Name:           "Сверло по металлу",
		VendorArticle:  "DR-10",
		RecommendPrice: 150.00,
		Brand:          "Bosch",
		Category:       "wb: 2384 - сверла",
		CategoryID:     int64Ptr(1234567890), // Явно указанный ID
		Description:    "Сверло по металлу 10мм",
	}
	
	// Генерация не произойдет, так как ID уже указан
	createReqWithID.GenerateCategoryID()
	
	fmt.Printf("   Категория: %s\n", createReqWithID.Category)
	fmt.Printf("   Используемый category_id: %d\n", *createReqWithID.CategoryID)
	
	// Пример 4: Обновление продукта с изменением категории на WB
	fmt.Println("\n4. Обновление продукта с изменением категории на WB:")
	updateReq := products.UpdateProductRequest{
		Category: stringPtr("wb: 2384 - сверла"),
	}
	
	// Автоматическая генерация нового category_id
	updateReq.GenerateCategoryID()
	
	fmt.Printf("   Новая категория: %s\n", *updateReq.Category)
	fmt.Printf("   Сгенерированный category_id: %d\n", *updateReq.CategoryID)
	
	// Пример 5: Массовое создание продуктов с WB категориями
	fmt.Println("\n5. Массовое создание продуктов с WB категориями:")
	bulkReq := products.CreateProductsRequest{
		Products: []products.CreateProductRequest{
			{
				Name:           "Молоток слесарный",
				VendorArticle:  "HM-500",
				RecommendPrice: 800.00,
				Brand:          "Stanley",
				Category:       "wb: 2066 - молотки",
				Description:    "Молоток слесарный 500г",
			},
			{
				Name:           "Отвертка крестовая",
				VendorArticle:  "SD-PH2",
				RecommendPrice: 120.00,
				Brand:          "Wera",
				Category:       "wb: 2294 - отвертки",
				Description:    "Отвертка крестовая PH2",
			},
		},
	}
	
	// Автоматическая генерация для всех продуктов
	for i := range bulkReq.Products {
		bulkReq.Products[i].GenerateCategoryID()
	}
	
	for i, product := range bulkReq.Products {
		fmt.Printf("   Продукт %d: %s\n", i+1, product.Name)
		fmt.Printf("   Категория: %s\n", product.Category)
		fmt.Printf("   category_id: %d\n", *product.CategoryID)
	}
	
	// Пример 6: Прямое использование утилит
	fmt.Println("\n6. Прямое использование утилит:")
	
	// Генерация ID из WB категории
	wbCategoryID := utils.GenerateCategoryID("wb: 1318 - сварочные аппараты")
	fmt.Printf("   GenerateCategoryID('wb: 1318 - сварочные аппараты') = %d\n", wbCategoryID)
	
	// Генерация ID из НЕ-WB категории (должен вернуть 0)
	nonWBCategoryID := utils.GenerateCategoryID("1318 - сварочные аппараты")
	fmt.Printf("   GenerateCategoryID('1318 - сварочные аппараты') = %d (должен быть 0)\n", nonWBCategoryID)
	
	// Генерация ID из частей
	partsID := utils.GenerateCategoryIDFromParts("1318", "сварочные аппараты")
	fmt.Printf("   GenerateCategoryIDFromParts('1318', 'сварочные аппараты') = %d\n", partsID)
	
	// Создание названия из частей
	categoryName := utils.GenerateCategoryNameFromParts("1318", "сварочные аппараты")
	fmt.Printf("   GenerateCategoryNameFromParts('1318', 'сварочные аппараты') = '%s'\n", categoryName)
	
	// Пример 7: JSON сериализация
	fmt.Println("\n7. JSON сериализация:")
	jsonData, err := json.MarshalIndent(createReq, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   JSON:\n%s\n", string(jsonData))
	
	fmt.Println("\n=== Примеры завершены ===")
}

// Вспомогательные функции для создания указателей
func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
