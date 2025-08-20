package classifier

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// Cache интерфейс для кэширования
type Cache interface {
	Delete(key string) error
}

// ProductResultHandler обработчик результатов классификации продуктов
type ProductResultHandler struct {
	db    *sql.DB
	cache Cache
}

// NewProductResultHandler создает новый обработчик результатов
func NewProductResultHandler(db *sql.DB) *ProductResultHandler {
	return &ProductResultHandler{
		db:    db,
		cache: nil, // Кэш может быть установлен позже
	}
}

// HandleClassificationResult обрабатывает результат классификации
// Обновляет статус на 'classified' только если вероятность >= 0.99 и ID числовые
func (h *ProductResultHandler) HandleClassificationResult(productID int64, result *ClassificationResult) error {
	log.Printf("🔍 Начинаю обработку результата классификации для продукта %d", productID)
	log.Printf("📊 Результат: Status=%s, Accuracy=%.2f, BrandAccuracy=%.2f", result.Status, result.Accuracy, result.BrandAccuracy)
	log.Printf("📊 FoundCategoryID: %v (тип: %T)", result.FoundCategoryID, result.FoundCategoryID)
	log.Printf("📊 FoundBrandID: %v (тип: %T)", result.FoundBrandID, result.FoundBrandID)

	if result.Status != "found" {
		log.Printf("⚠️ Классификация продукта %d не удалась: статус=%s", productID, result.Status)
		return nil
	}

	// Проверяем вероятность классификации
	var categoryID, brandID *int64

	// Проверяем категорию (вероятность должна быть >= 0.99 и ID числовой)
	if result.FoundCategoryID != nil && result.Accuracy >= 0.99 {
		// Проверяем, что ID категории - это число или строка с числом
		switch v := result.FoundCategoryID.(type) {
		case float64:
			// JSON числа парсятся как float64
			catID := int64(v)
			categoryID = &catID
			log.Printf("✅ Категория продукта %d найдена с вероятностью %.2f (ID: %d)",
				productID, result.Accuracy, catID)
		case string:
			// API может вернуть строку с числом - пытаемся преобразовать
			if catID, err := strconv.ParseInt(v, 10, 64); err == nil {
				categoryID = &catID
				log.Printf("✅ Категория продукта %d найдена с вероятностью %.2f (ID: %d, преобразован из строки)",
					productID, result.Accuracy, catID)
			} else {
				log.Printf("❌ Ошибка API: found_category_id не является числом, получено: %s", v)
				// Устанавливаем статус not_classified
				if err := h.setProductNotClassified(productID); err != nil {
					log.Printf("❌ Ошибка установки статуса not_classified для продукта %d: %v", productID, err)
				}
				return nil
			}
		default:
			log.Printf("❌ Неожиданный тип found_category_id: %T", result.FoundCategoryID)
			// Устанавливаем статус not_classified
			if err := h.setProductNotClassified(productID); err != nil {
				log.Printf("❌ Ошибка установки статуса not_classified для продукта %d: %v", productID, err)
			}
			return nil
		}
	} else if result.FoundCategoryID != nil {
		log.Printf("⚠️ Категория продукта %d найдена, но вероятность %.2f < 0.99",
			productID, result.Accuracy)
	} else {
		log.Printf("❌ Категория для продукта %d не найдена", productID)
	}

	// Проверяем бренд (если есть)
	if result.FoundBrandID != nil && result.BrandAccuracy >= 0.99 {
		// Проверяем, что ID бренда - это число или строка с числом
		switch v := result.FoundBrandID.(type) {
		case float64:
			// JSON числа парсятся как float64
			brID := int64(v)
			brandID = &brID
			log.Printf("✅ Бренд продукта %d найден с вероятностью %.2f (ID: %d)",
				productID, result.BrandAccuracy, brID)
		case string:
			// API может вернуть строку с числом - пытаемся преобразовать
			if brID, err := strconv.ParseInt(v, 10, 64); err == nil {
				brandID = &brID
				log.Printf("✅ Бренд продукта %d найден с вероятностью %.2f (ID: %d, преобразован из строки)",
					productID, result.BrandAccuracy, brID)
			} else {
				log.Printf("❌ Ошибка API: brand_id не является числом, получено: %s", v)
				// Устанавливаем статус not_classified
				if err := h.setProductNotClassified(productID); err != nil {
					log.Printf("❌ Ошибка установки статуса not_classified для продукта %d: %v", productID, err)
				}
				return nil
			}
		default:
			log.Printf("❌ Неожиданный тип brand_id: %T", result.FoundBrandID)
			// Устанавливаем статус not_classified
			if err := h.setProductNotClassified(productID); err != nil {
				log.Printf("❌ Ошибка установки статуса not_classified для продукта %d: %v", productID, err)
			}
			return nil
		}
	} else if result.FoundBrandID != nil {
		log.Printf("⚠️ Бренд продукта %d найден, но вероятность %.2f < 0.99",
			productID, result.Accuracy)
	} else {
		log.Printf("❌ Бренд для продукта %d не найден", productID)
	}

	// Проверяем, что ОБА ID найдены с вероятностью >= 0.99 (обязательно!)
	log.Printf("🔍 Проверяю результат: categoryID=%v, brandID=%v", categoryID, brandID)

	if categoryID != nil && brandID != nil {
		// Оба ID найдены - устанавливаем статус classified
		log.Printf("🎯 Оба ID найдены, устанавливаю статус classified для продукта %d", productID)
		if err := h.updateProductClassification(productID, categoryID, brandID); err != nil {
			log.Printf("❌ Ошибка обновления классификации продукта %d: %v", productID, err)
			return fmt.Errorf("ошибка обновления классификации продукта %d: %w", productID, err)
		}
		log.Printf("🎯 Продукт %d успешно классифицирован (category_id=%d, brand_id=%d)",
			productID, *categoryID, *brandID)
	} else {
		// Не все ID найдены - устанавливаем статус not_classified
		log.Printf("⏳ Продукт %d остается неклассифицированным (не все ID найдены)", productID)
		if categoryID == nil {
			log.Printf("   - categoryID не найден")
		}
		if brandID == nil {
			log.Printf("   - brandID не найден")
		}
		if err := h.setProductNotClassified(productID); err != nil {
			log.Printf("❌ Ошибка установки статуса not_classified для продукта %d: %v", productID, err)
		}
	}

	log.Printf("✅ Обработка результата классификации для продукта %d завершена", productID)
	return nil
}

// setProductNotClassified устанавливает статус продукта как not_classified
func (h *ProductResultHandler) setProductNotClassified(productID int64) error {
	// Начинаем транзакцию
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Обновляем статус на not_classified
	query := "UPDATE products SET status = 'not_classified', updated_at = NOW() WHERE id = ?"
	result, err := tx.Exec(query, productID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения UPDATE: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("продукт с ID %d не найден", productID)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	log.Printf("⚠️ Продукт %d установлен как not_classified (ошибка API)", productID)
	return nil
}

// updateProductClassification обновляет классификацию продукта в базе данных
func (h *ProductResultHandler) updateProductClassification(productID int64, categoryID, brandID *int64) error {
	// Начинаем транзакцию
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Формируем SQL запрос для обновления
	query := "UPDATE products SET status = 'classified'"
	var args []interface{}
	var placeholders []string

	// Добавляем обновление category_id если найден
	if categoryID != nil {
		placeholders = append(placeholders, "category_id = ?")
		args = append(args, *categoryID)
	}

	// Добавляем обновление brand_id если найден
	if brandID != nil {
		placeholders = append(placeholders, "brand_id = ?")
		args = append(args, *brandID)
	}

	// Добавляем updated_at
	placeholders = append(placeholders, "updated_at = NOW()")

	// Формируем полный запрос
	if len(placeholders) > 0 {
		query += ", " + placeholders[0]
		for i := 1; i < len(placeholders); i++ {
			query += ", " + placeholders[i]
		}
	}

	query += " WHERE id = ?"
	args = append(args, productID)

	// Выполняем обновление
	result, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка выполнения UPDATE: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("продукт с ID %d не найден", productID)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	log.Printf("✅ Продукт %d обновлен: category_id=%v, brand_id=%v, status=classified",
		productID, categoryID, brandID)

	// Инвалидируем кэш для этого продукта
	if h.cache != nil {
		cacheKey := fmt.Sprintf("product:%d", productID)
		if err := h.cache.Delete(cacheKey); err != nil {
			log.Printf("⚠️ Ошибка инвалидации кэша для продукта %d: %v", productID, err)
		} else {
			log.Printf("🗑️ Кэш для продукта %d инвалидирован", productID)
		}
	}

	return nil
}
