package classifier

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ProductClassificationRequest запрос на классификацию продукта
type ProductClassificationRequest struct {
	ProductID    int64  `json:"product_id"`
	ProductName  string `json:"product_name"`
	UserCategory string `json:"user_category"`
	UserID       int64  `json:"user_id"`
	Brand        string `json:"brand"`
	Category     string `json:"category"`
}

// ClassificationResult результат классификации от API catformat
type ClassificationResult struct {
	Status          string      `json:"status,omitempty"`
	Request         string      `json:"request,omitempty"`
	FoundCategory   string      `json:"found_category,omitempty"`
	FoundCategoryID interface{} `json:"found_category_id,omitempty"` // Может быть string или number
	FoundBrand      string      `json:"found_brand,omitempty"`
	FoundBrandID    interface{} `json:"brand_id,omitempty"` // API возвращает brand_id, а не found_brand_id
	Accuracy        float64     `json:"accuracy,omitempty"`
	BrandAccuracy   float64     `json:"brand_accuracy,omitempty"` // API возвращает brand_accuracy
	UserCategory    string      `json:"user_category,omitempty"`
	User            string      `json:"user,omitempty"`
	Error           string      `json:"error,omitempty"`
}

// ProductClassifier интерфейс для классификации продуктов
type ProductClassifier interface {
	ClassifyProduct(ctx context.Context, req ProductClassificationRequest) (*ClassificationResult, error)
	ClassifyProductAsync(req ProductClassificationRequest) error
}

// OSTKClassifier реализация классификатора через API OSTK
type OSTKClassifier struct {
	apiBaseURL string
	httpClient *http.Client
	// Канал для асинхронной обработки
	classificationQueue chan ProductClassificationRequest
	// Обработчик результатов
	resultHandler ResultHandler
}

// ResultHandler интерфейс для обработки результатов классификации
type ResultHandler interface {
	HandleClassificationResult(productID int64, result *ClassificationResult) error
}

// NewOSTKClassifier создает новый классификатор OSTK
func NewOSTKClassifier(apiBaseURL string, resultHandler ResultHandler) *OSTKClassifier {
	classifier := &OSTKClassifier{
		apiBaseURL: apiBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		classificationQueue: make(chan ProductClassificationRequest, 1000), // Буфер на 1000 задач
		resultHandler:       resultHandler,
	}

	// Запускаем воркер для обработки очереди
	go classifier.startWorker()

	return classifier
}

// ClassifyProductAsync добавляет задачу в очередь для асинхронной обработки
func (c *OSTKClassifier) ClassifyProductAsync(req ProductClassificationRequest) error {
	select {
	case c.classificationQueue <- req:
		log.Printf("✅ Задача классификации добавлена в очередь для продукта %d", req.ProductID)
		return nil
	default:
		return fmt.Errorf("очередь классификации переполнена")
	}
}

// startWorker запускает воркер для обработки очереди классификации
func (c *OSTKClassifier) startWorker() {
	log.Printf("🚀 Воркер классификации запущен")

	// Создаем контекст для воркера
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	for {
		select {
		case req, ok := <-c.classificationQueue:
			if !ok {
				log.Printf("🔒 Очередь классификации закрыта, завершаю воркер")
				return
			}

			log.Printf("🔍 Обрабатываю классификацию продукта %d: %s", req.ProductID, req.ProductName)

			// Создаем контекст с таймаутом для каждого запроса
			ctx, cancel := context.WithTimeout(workerCtx, 30*time.Second)

			// Выполняем классификацию
			result, err := c.ClassifyProduct(ctx, req)
			if err != nil {
				log.Printf("❌ Ошибка классификации продукта %d: %v", req.ProductID, err)
				cancel()
				continue
			}

			// Обрабатываем результат
			if err := c.resultHandler.HandleClassificationResult(req.ProductID, result); err != nil {
				log.Printf("❌ Ошибка обработки результата для продукта %d: %v", req.ProductID, err)
			}

			cancel() // Освобождаем ресурсы контекста

			// Небольшая пауза между запросами
			time.Sleep(100 * time.Millisecond)

		case <-workerCtx.Done():
			log.Printf("🔒 Воркер классификации остановлен")
			return
		}
	}
}

// ClassifyProduct выполняет синхронную классификацию продукта
func (c *OSTKClassifier) ClassifyProduct(ctx context.Context, req ProductClassificationRequest) (*ClassificationResult, error) {
	// Проверяем контекст перед началом
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Printf("🔍 Начинаю классификацию продукта %d: %s", req.ProductID, req.ProductName)

	// Формируем URL для API
	apiURL := fmt.Sprintf("%s/products/catformat.php", c.apiBaseURL)

	// Создаем параметры запроса
	params := url.Values{}
	params.Set("product_name", req.ProductName)
	params.Set("user", strconv.FormatInt(req.UserID, 10))
	params.Set("user_category", req.UserCategory)
	params.Set("brand", req.Brand) // Добавляем бренд как он называется у партнера для улучшения точности классификации
	params.Set("findbrand", "1")
	params.Set("findcat", "1")

	log.Printf("🔍 Отправляю запрос к API: %s?%s", apiURL, params.Encode())

	// Создаем HTTP запрос
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("🔍 Получен ответ от API: статус %d", resp.StatusCode)

	// Проверяем контекст после получения ответа
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	// Парсим ответ
	var result ClassificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON ответа: %w", err)
	}

	log.Printf("🔍 Результат парсинга JSON: %+v", result)

	// Логируем результат
	if result.Status == "found" {
		log.Printf("📊 Результат классификации продукта %d: категория=%s (ID: %v, точность=%.2f), бренд=%s (ID: %v, точность=%.2f)",
			req.ProductID, result.FoundCategory, result.FoundCategoryID, result.Accuracy, result.FoundBrand, result.FoundBrandID, result.BrandAccuracy)
	} else {
		log.Printf("📊 Результат классификации продукта %d: статус=%s", req.ProductID, result.Status)
	}

	return &result, nil
}

// Close закрывает классификатор
func (c *OSTKClassifier) Close() {
	close(c.classificationQueue)
	log.Printf("🔒 Классификатор закрыт")
}
