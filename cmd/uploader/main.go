package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"portaldata-api/internal/config"
	"portaldata-api/internal/services"
)

type Uploader struct {
	apiURL    string
	apiKey    string
	parserURL string
}

func NewUploader(apiURL, apiKey, parserURL string) *Uploader {
	return &Uploader{
		apiURL:    apiURL,
		apiKey:    apiKey,
		parserURL: parserURL,
	}
}

func (u *Uploader) UploadProducts() error {
	offset := 0
	limit := 10
	totalProcessed := 0

	log.Printf("Starting product upload process...")

	for {
		// Получаем данные из parser
		parserURL := fmt.Sprintf("%s?path=tn/020225.txt&limit=%d&offset=%d&exclude=Описание,Наличие,Валюта,URL,Изображения",
			u.parserURL, limit, offset)

		log.Printf("Fetching products from: %s", parserURL)

		resp, err := http.Get(parserURL)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("parser returned status: %d", resp.StatusCode)
		}

		var parserResponse services.ParserResponse
		if err := json.NewDecoder(resp.Body).Decode(&parserResponse); err != nil {
			return fmt.Errorf("failed to decode parser response: %w", err)
		}

		if len(parserResponse.Products) == 0 {
			log.Printf("No more products to process. Total processed: %d", totalProcessed)
			break
		}

		// Преобразуем данные для API
		var products []map[string]interface{}
		for _, p := range parserResponse.Products {
			product := map[string]interface{}{
				"name":           p.Name,
				"vendor_article": p.VendorArticle,
				"price":          p.Price,
				"brand":          p.Brand,
				"category":       p.Category,
				"description":    p.Description,
			}
			products = append(products, product)
		}

		// Отправляем данные в API
		if err := u.sendToAPI(products); err != nil {
			return fmt.Errorf("failed to send to API: %w", err)
		}

		totalProcessed += len(products)
		log.Printf("Processed batch: %d products (offset: %d, total: %d)", 
			len(products), offset, totalProcessed)

		offset += limit

		// Небольшая пауза между запросами
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Upload completed. Total products processed: %d", totalProcessed)
	return nil
}

func (u *Uploader) sendToAPI(products []map[string]interface{}) error {
	payload := map[string]interface{}{
		"products": products,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	apiURL := fmt.Sprintf("%s?api_key=%s", u.apiURL, u.apiKey)
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = http.NoBody // Используем GET запрос для совместимости

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	log.Printf("Successfully sent %d products to API", len(products))
	return nil
}

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем uploader
	uploader := NewUploader(
		"http://localhost:8080/v1/products/products.php",
		cfg.APIKey,
		"http://localhost:8080/v1/products/parser.php",
	)

	// Запускаем загрузку
	if err := uploader.UploadProducts(); err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	log.Println("Upload completed successfully")
} 