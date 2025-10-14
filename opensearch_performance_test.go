package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

// TestResult представляет результат одного теста
type TestResult struct {
	Query     string        `json:"query"`
	Duration  time.Duration `json:"duration_ms"`
	Hits      int           `json:"hits"`
	Error     string        `json:"error,omitempty"`
}

// PerformanceTest представляет набор тестов
type PerformanceTest struct {
	BaseURL string
	Tests   []TestResult
}

// NewPerformanceTest создает новый тест производительности
func NewPerformanceTest(baseURL string) *PerformanceTest {
	return &PerformanceTest{
		BaseURL: baseURL,
		Tests:   make([]TestResult, 0),
	}
}

// RunTest выполняет один тест
func (pt *PerformanceTest) RunTest(query, endpoint string) TestResult {
	start := time.Now()
	
	// Формируем URL
	url := fmt.Sprintf("%s/api/opensearch/?%s", pt.BaseURL, query)
	
	// Выполняем запрос
	resp, err := http.Get(url)
	if err != nil {
		return TestResult{
			Query:    query,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			Query:    query,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}
	
	// Парсим JSON ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return TestResult{
			Query:    query,
			Duration: time.Since(start),
			Error:    err.Error(),
		}
	}
	
	// Извлекаем количество результатов
	hits := 0
	if hitsData, ok := result["hits"].(map[string]interface{}); ok {
		if totalData, ok := hitsData["total"].(map[string]interface{}); ok {
			if value, ok := totalData["value"].(float64); ok {
				hits = int(value)
			}
		}
	}
	
	return TestResult{
		Query:    query,
		Duration: time.Since(start),
		Hits:     hits,
	}
}

// RunMultipleTests выполняет несколько тестов
func (pt *PerformanceTest) RunMultipleTests(query string, iterations int) []TestResult {
	results := make([]TestResult, 0, iterations)
	
	for i := 0; i < iterations; i++ {
		result := pt.RunTest(query, "/api/search/")
		results = append(results, result)
		
		// Небольшая пауза между запросами
		time.Sleep(100 * time.Millisecond)
	}
	
	return results
}

// RunComprehensiveTest выполняет комплексный тест производительности
func (pt *PerformanceTest) RunComprehensiveTest() {
	fmt.Println("🚀 Запуск комплексного теста производительности Go API")
	fmt.Println("=" * 60)
	
	// Тестовые запросы
	testQueries := []string{
		"query=дрель&size=10",
		"query=сварочный&size=20",
		"query=инструмент&size=50",
		"brand=DeWalt&size=30",
		"brand=Bosch&size=40",
		"category=электроинструмент&size=25",
		"price_min=1000&price_max=5000&size=15",
		"query=перфоратор&brand=Makita&size=20",
		"query=болгарка&price_min=2000&size=30",
		"query=отвертка&size=100",
	}
	
	// Выполняем тесты
	allResults := make([]TestResult, 0)
	
	for i, query := range testQueries {
		fmt.Printf("Тест %d/%d: %s\n", i+1, len(testQueries), query)
		
		// Выполняем 5 итераций для каждого запроса
		results := pt.RunMultipleTests(query, 5)
		allResults = append(allResults, results...)
		
		// Выводим статистику для этого запроса
		pt.printQueryStats(query, results)
		fmt.Println()
	}
	
	// Общая статистика
	pt.printOverallStats(allResults)
}

// printQueryStats выводит статистику для одного запроса
func (pt *PerformanceTest) printQueryStats(query string, results []TestResult) {
	if len(results) == 0 {
		fmt.Printf("  ❌ Нет результатов для запроса: %s\n", query)
		return
	}
	
	// Фильтруем успешные результаты
	successResults := make([]TestResult, 0)
	for _, result := range results {
		if result.Error == "" {
			successResults = append(successResults, result)
		}
	}
	
	if len(successResults) == 0 {
		fmt.Printf("  ❌ Все запросы завершились с ошибкой\n")
		return
	}
	
	// Вычисляем статистику
	durations := make([]time.Duration, len(successResults))
	for i, result := range successResults {
		durations[i] = result.Duration
	}
	
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	
	avg := time.Duration(0)
	for _, d := range durations {
		avg += d
	}
	avg = avg / time.Duration(len(durations))
	
	min := durations[0]
	max := durations[len(durations)-1]
	median := durations[len(durations)/2]
	
	fmt.Printf("  ✅ Успешных запросов: %d/%d\n", len(successResults), len(results))
	fmt.Printf("  ⏱️  Среднее время: %v\n", avg)
	fmt.Printf("  ⏱️  Минимальное время: %v\n", min)
	fmt.Printf("  ⏱️  Максимальное время: %v\n", max)
	fmt.Printf("  ⏱️  Медианное время: %v\n", median)
	fmt.Printf("  📊 Найдено результатов: %d\n", successResults[0].Hits)
}

// printOverallStats выводит общую статистику
func (pt *PerformanceTest) printOverallStats(results []TestResult) {
	if len(results) == 0 {
		fmt.Println("❌ Нет результатов для анализа")
		return
	}
	
	// Фильтруем успешные результаты
	successResults := make([]TestResult, 0)
	for _, result := range results {
		if result.Error == "" {
			successResults = append(successResults, result)
		}
	}
	
	if len(successResults) == 0 {
		fmt.Println("❌ Все запросы завершились с ошибкой")
		return
	}
	
	// Вычисляем общую статистику
	durations := make([]time.Duration, len(successResults))
	for i, result := range successResults {
		durations[i] = result.Duration
	}
	
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	
	avg := time.Duration(0)
	for _, d := range durations {
		avg += d
	}
	avg = avg / time.Duration(len(durations))
	
	min := durations[0]
	max := durations[len(durations)-1]
	median := durations[len(durations)/2]
	
	fmt.Println("📊 ОБЩАЯ СТАТИСТИКА GO API")
	fmt.Println("=" * 40)
	fmt.Printf("✅ Успешных запросов: %d/%d (%.1f%%)\n", 
		len(successResults), len(results), 
		float64(len(successResults))/float64(len(results))*100)
	fmt.Printf("⏱️  Среднее время: %v\n", avg)
	fmt.Printf("⏱️  Минимальное время: %v\n", min)
	fmt.Printf("⏱️  Максимальное время: %v\n", max)
	fmt.Printf("⏱️  Медианное время: %v\n", median)
	fmt.Printf("🚀 Запросов в секунду: %.1f\n", 
		float64(len(successResults))/avg.Seconds())
}

func main() {
	// Конфигурация
	baseURL := "http://localhost:8095"
	
	// Проверяем доступность API
	fmt.Println("🔍 Проверка доступности Go API...")
	resp, err := http.Get(baseURL + "/api/opensearch/health")
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к Go API: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		log.Fatalf("❌ Go API недоступен, статус: %d", resp.StatusCode)
	}
	
	fmt.Println("✅ Go API доступен")
	fmt.Println()
	
	// Создаем и запускаем тест
	test := NewPerformanceTest(baseURL)
	test.RunComprehensiveTest()
	
	fmt.Println()
	fmt.Println("🎯 Тест завершен!")
}
