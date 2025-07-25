<?php
/**
 * COMPREHENSIVE API TEST 2025
 * Полный тест всех возможностей PortalData API
 * Включает проверку медиа функциональности
 */

class ComprehensiveAPITest2025 {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
    private $testResults = [];
    private $createdIds = [];
    private $performanceData = [];
    private $startTime;

    public function run() {
        $this->startTime = microtime(true);
        
        echo "🚀 COMPREHENSIVE API TEST 2025\n";
        echo "===============================\n";
        echo "Время запуска: " . date('Y-m-d H:i:s') . "\n";
        echo "API URL: {$this->baseUrl}\n";
        echo "API Key: " . substr($this->apiKey, 0, 20) . "...\n\n";

        // 1. Проверка доступности сервера
        $this->testServerAvailability();
        
        // 2. Тест продуктов (с медиа)
        $this->testProducts();
        
        // 3. Тест складов
        $this->testWarehouses();
        
        // 4. Тест офферов
        $this->testOffers();
        
        // 5. Тест заказов
        $this->testOrders();
        
        // 6. Тест публичных endpoints
        $this->testPublicEndpoints();
        
        // 7. Тест ошибок и валидации
        $this->testErrorHandling();
        
        // 8. Тест производительности
        $this->testPerformance();
        
        // 9. Очистка тестовых данных
        $this->cleanupTestData();
        
        // 10. Вывод результатов
        $this->printResults();
    }

    private function testServerAvailability() {
        echo "🔍 1. ПРОВЕРКА ДОСТУПНОСТИ СЕРВЕРА\n";
        echo "------------------------------------\n";
        
        // Проверка основного endpoint
        $response = $this->makeRequest('GET', '');
        $this->assertTest('Основной endpoint', $response['status'] === 200, $response);
        
        // Проверка Swagger
        $response = $this->makeRequest('GET', '/swagger/index.html');
        $this->assertTest('Swagger UI', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testProducts() {
        echo "📦 2. ТЕСТИРОВАНИЕ ПРОДУКТОВ\n";
        echo "------------------------------\n";
        
        // Создание продукта без медиа
        $productData = [
            'name' => 'Тестовый продукт ' . time(),
            'vendor_article' => 'TEST-' . time(),
            'recommend_price' => 1500.50,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Описание тестового продукта'
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData);
        $this->assertTest('Создание продукта', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdIds['product'] = $response['data']['id'];
        }
        
        // Создание продукта с медиа (если поддерживается)
        $productWithMedia = [
            'name' => 'Продукт с медиа ' . time(),
            'vendor_article' => 'MEDIA-' . time(),
            'recommend_price' => 2500.00,
            'brand' => 'MediaBrand',
            'category' => 'MediaCategory',
            'description' => 'Продукт с медиа контентом',
            'image_urls' => [
                'https://example.com/image1.jpg',
                'https://example.com/image2.jpg'
            ],
            'video_urls' => [
                'https://example.com/video1.mp4'
            ],
            'model_3d_urls' => [
                'https://example.com/model1.glb'
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products', $productWithMedia);
        $this->assertTest('Создание продукта с медиа', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdIds['product_with_media'] = $response['data']['id'];
        }
        
        // Получение списка продуктов
        $response = $this->makeRequest('GET', '/products');
        $this->assertTest('Получение списка продуктов', $response['status'] === 200, $response);
        
        // Получение конкретного продукта
        if (isset($this->createdIds['product'])) {
            $response = $this->makeRequest('GET', '/products/' . $this->createdIds['product']);
            $this->assertTest('Получение продукта по ID', $response['status'] === 200, $response);
        }
        
        // Обновление продукта
        if (isset($this->createdIds['product'])) {
            $updateData = [
                'name' => 'Обновленный продукт ' . time(),
                'description' => 'Обновленное описание'
            ];
            
            $response = $this->makeRequest('PUT', '/products/' . $this->createdIds['product'], $updateData);
            $this->assertTest('Обновление продукта', $response['status'] === 200, $response);
        }
        
        // Пакетное создание продуктов
        $batchData = [
            'products' => [
                [
                    'name' => 'Пакетный продукт 1',
                    'vendor_article' => 'BATCH-001-' . time(),
                    'recommend_price' => 100.00,
                    'brand' => 'BatchBrand',
                    'category' => 'BatchCategory',
                    'description' => 'Пакетный продукт 1'
                ],
                [
                    'name' => 'Пакетный продукт 2',
                    'vendor_article' => 'BATCH-002-' . time(),
                    'recommend_price' => 200.00,
                    'brand' => 'BatchBrand',
                    'category' => 'BatchCategory',
                    'description' => 'Пакетный продукт 2'
                ]
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products/batch', $batchData);
        $this->assertTest('Пакетное создание продуктов', $response['status'] === 201, $response);
        
        echo "\n";
    }

    private function testWarehouses() {
        echo "🏭 3. ТЕСТИРОВАНИЕ СКЛАДОВ\n";
        echo "----------------------------\n";
        
        // Создание склада
        $warehouseData = [
            'name' => 'Тестовый склад ' . time(),
            'address' => 'ул. Тестовая, 123',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'contact_phone' => '+7-999-123-45-67',
            'contact_email' => 'warehouse@test.com'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData);
        $this->assertTest('Создание склада', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdIds['warehouse'] = $response['data']['id'];
        }
        
        // Получение списка складов
        $response = $this->makeRequest('GET', '/warehouses');
        $this->assertTest('Получение списка складов', $response['status'] === 200, $response);
        
        // Обновление склада
        if (isset($this->createdIds['warehouse'])) {
            $updateData = [
                'name' => 'Обновленный склад ' . time(),
                'contact_phone' => '+7-999-987-65-43'
            ];
            
            $response = $this->makeRequest('PUT', '/warehouses/' . $this->createdIds['warehouse'], $updateData);
            $this->assertTest('Обновление склада', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function testOffers() {
        echo "📋 4. ТЕСТИРОВАНИЕ ОФФЕРОВ\n";
        echo "-----------------------------\n";
        
        // Создание оффера
        $offerData = [
            'product_id' => isset($this->createdIds['product']) ? $this->createdIds['product'] : 1,
            'warehouse_id' => isset($this->createdIds['warehouse']) ? $this->createdIds['warehouse'] : 1,
            'price_per_unit' => 1500.00,
            'tax_nds' => 20,
            'units_per_lot' => 1,
            'available_lots' => 10,
            'offer_type' => 'sale',
            'max_shipping_days' => 3
        ];
        
        $response = $this->makeRequest('POST', '/offers', $offerData);
        $this->assertTest('Создание оффера', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdIds['offer'] = $response['data']['offer_id'];
        }
        
        // Получение списка офферов
        $response = $this->makeRequest('GET', '/offers');
        $this->assertTest('Получение списка офферов', $response['status'] === 200, $response);
        
        // Получение конкретного оффера
        if (isset($this->createdIds['offer'])) {
            $response = $this->makeRequest('GET', '/offers/' . $this->createdIds['offer']);
            $this->assertTest('Получение оффера по ID', $response['status'] === 200, $response);
        }
        
        // Обновление оффера
        if (isset($this->createdIds['offer'])) {
            $updateData = [
                'price_per_unit' => 1600.00,
                'available_lots' => 15
            ];
            
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdIds['offer'], $updateData);
            $this->assertTest('Обновление оффера', $response['status'] === 200, $response);
        }
        
        // Пакетное создание офферов
        $batchOffersData = [
            'offers' => [
                [
                    'product_id' => isset($this->createdIds['product']) ? $this->createdIds['product'] : 1,
                    'warehouse_id' => isset($this->createdIds['warehouse']) ? $this->createdIds['warehouse'] : 1,
                    'price_per_unit' => 1000.00,
                    'tax_nds' => 20,
                    'units_per_lot' => 1,
                    'available_lots' => 5,
                    'offer_type' => 'sale',
                    'max_shipping_days' => 2
                ],
                [
                    'product_id' => isset($this->createdIds['product']) ? $this->createdIds['product'] : 1,
                    'warehouse_id' => isset($this->createdIds['warehouse']) ? $this->createdIds['warehouse'] : 1,
                    'price_per_unit' => 2000.00,
                    'tax_nds' => 20,
                    'units_per_lot' => 1,
                    'available_lots' => 8,
                    'offer_type' => 'buy',
                    'max_shipping_days' => 5
                ]
            ]
        ];
        
        $response = $this->makeRequest('POST', '/offers/batch', $batchOffersData);
        $this->assertTest('Пакетное создание офферов', $response['status'] === 201, $response);
        
        echo "\n";
    }

    private function testOrders() {
        echo "📦 5. ТЕСТИРОВАНИЕ ЗАКАЗОВ\n";
        echo "-----------------------------\n";
        
        // Создание заказа
        $orderData = [
            'offer_id' => isset($this->createdIds['offer']) ? $this->createdIds['offer'] : 1,
            'quantity' => 2,
            'delivery_address' => 'ул. Заказная, 456',
            'contact_phone' => '+7-999-111-22-33',
            'contact_email' => 'order@test.com'
        ];
        
        $response = $this->makeRequest('POST', '/orders', $orderData);
        $this->assertTest('Создание заказа', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdIds['order'] = $response['data']['id'];
        }
        
        // Получение списка заказов
        $response = $this->makeRequest('GET', '/orders');
        $this->assertTest('Получение списка заказов', $response['status'] === 200, $response);
        
        // Получение конкретного заказа
        if (isset($this->createdIds['order'])) {
            $response = $this->makeRequest('GET', '/orders/' . $this->createdIds['order']);
            $this->assertTest('Получение заказа по ID', $response['status'] === 200, $response);
        }
        
        // Обновление статуса заказа
        if (isset($this->createdIds['order'])) {
            $statusData = [
                'status' => 'processing'
            ];
            
            $response = $this->makeRequest('PUT', '/orders/' . $this->createdIds['order'] . '/status', $statusData);
            $this->assertTest('Обновление статуса заказа', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function testPublicEndpoints() {
        echo "🌐 6. ТЕСТИРОВАНИЕ ПУБЛИЧНЫХ ENDPOINTS\n";
        echo "----------------------------------------\n";
        
        // Публичные офферы (без авторизации)
        $response = $this->makeRequest('GET', '/offers/public', null, null);
        $this->assertTest('Публичные офферы', $response['status'] === 200, $response);
        
        // WB Stock
        $response = $this->makeRequest('GET', '/offers/wb_stock?product_id=1&warehouse_id=1&supplier_id=42009');
        $this->assertTest('WB Stock', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testErrorHandling() {
        echo "⚠️ 7. ТЕСТИРОВАНИЕ ОБРАБОТКИ ОШИБОК\n";
        echo "--------------------------------------\n";
        
        // Неверный API ключ
        $response = $this->makeRequest('GET', '/products', null, 'invalid_key');
        $this->assertTest('Неверный API ключ', $response['status'] === 401, $response);
        
        // Несуществующий ресурс
        $response = $this->makeRequest('GET', '/products/999999');
        $this->assertTest('Несуществующий продукт', $response['status'] === 404, $response);
        
        // Неверные данные
        $invalidData = [
            'name' => '', // пустое имя
            'vendor_article' => 'TEST'
        ];
        
        $response = $this->makeRequest('POST', '/products', $invalidData);
        $this->assertTest('Неверные данные продукта', $response['status'] === 400, $response);
        
        // Неверный метод
        $response = $this->makeRequest('PATCH', '/products/1');
        $this->assertTest('Неверный HTTP метод', $response['status'] === 404, $response);
        
        echo "\n";
    }

    private function testPerformance() {
        echo "⚡ 8. ТЕСТИРОВАНИЕ ПРОИЗВОДИТЕЛЬНОСТИ\n";
        echo "--------------------------------------\n";
        
        // Тест скорости получения списка продуктов
        $startTime = microtime(true);
        for ($i = 0; $i < 5; $i++) {
            $this->makeRequest('GET', '/products');
        }
        $endTime = microtime(true);
        $avgTime = (($endTime - $startTime) / 5) * 1000;
        
        $this->performanceData['Среднее время запроса продуктов'] = round($avgTime, 2);
        $this->assertTest('Производительность продуктов', $avgTime < 1000, ['time' => $avgTime]);
        
        // Тест скорости получения публичных офферов
        $startTime = microtime(true);
        for ($i = 0; $i < 5; $i++) {
            $this->makeRequest('GET', '/offers/public');
        }
        $endTime = microtime(true);
        $avgTime = (($endTime - $startTime) / 5) * 1000;
        
        $this->performanceData['Среднее время запроса публичных офферов'] = round($avgTime, 2);
        $this->assertTest('Производительность публичных офферов', $avgTime < 1000, ['time' => $avgTime]);
        
        echo "\n";
    }

    private function cleanupTestData() {
        echo "🧹 9. ОЧИСТКА ТЕСТОВЫХ ДАННЫХ\n";
        echo "--------------------------------\n";
        
        // Удаление созданных ресурсов
        if (isset($this->createdIds['order'])) {
            $response = $this->makeRequest('DELETE', '/orders/' . $this->createdIds['order']);
            $this->assertTest('Удаление заказа', $response['status'] === 200, $response);
        }
        
        if (isset($this->createdIds['offer'])) {
            $response = $this->makeRequest('DELETE', '/offers/' . $this->createdIds['offer']);
            $this->assertTest('Удаление оффера', $response['status'] === 200, $response);
        }
        
        if (isset($this->createdIds['product'])) {
            $response = $this->makeRequest('DELETE', '/products/' . $this->createdIds['product']);
            $this->assertTest('Удаление продукта', $response['status'] === 200, $response);
        }
        
        if (isset($this->createdIds['warehouse'])) {
            $response = $this->makeRequest('DELETE', '/warehouses/' . $this->createdIds['warehouse']);
            $this->assertTest('Удаление склада', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function makeRequest($method, $endpoint, $data = null, $apiKey = null) {
        $url = $this->baseUrl . $endpoint;
        $apiKey = $apiKey ?: $this->apiKey;
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Authorization: Bearer ' . $apiKey,
            'Content-Type: application/json'
        ]);
        
        if ($method === 'POST') {
            curl_setopt($ch, CURLOPT_POST, true);
            if ($data) {
                curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
            }
        } elseif ($method === 'PUT') {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'PUT');
            if ($data) {
                curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
            }
        } elseif ($method === 'DELETE') {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'DELETE');
        }
        
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        return [
            'status' => $httpCode,
            'data' => json_decode($response, true),
            'raw' => $response
        ];
    }

    private function assertTest($testName, $condition, $response) {
        $status = $condition ? '✅' : '❌';
        echo "$status $testName\n";
        
        if (!$condition) {
            echo "   Ошибка: " . ($response['raw'] ?? 'Неизвестная ошибка') . "\n";
        }
        
        $this->testResults[] = [
            'name' => $testName,
            'passed' => $condition,
            'response' => $response
        ];
    }

    private function printResults() {
        $endTime = microtime(true);
        $totalTime = round(($endTime - $this->startTime) * 1000, 2);
        
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ\n";
        echo "==========================\n";
        
        $passed = 0;
        $total = count($this->testResults);
        
        foreach ($this->testResults as $result) {
            if ($result['passed']) {
                $passed++;
            }
        }
        
        $successRate = round(($passed / $total) * 100, 2);
        
        echo "Всего тестов: $total\n";
        echo "Пройдено: $passed\n";
        echo "Провалено: " . ($total - $passed) . "\n";
        echo "Успешность: {$successRate}%\n";
        echo "Общее время: {$totalTime}ms\n\n";
        
        echo "📈 МЕТРИКИ ПРОИЗВОДИТЕЛЬНОСТИ:\n";
        foreach ($this->performanceData as $metric => $value) {
            echo "$metric: {$value}ms\n";
        }
        
        echo "\n🎯 СТАТУС API:\n";
        if ($successRate >= 90) {
            echo "🟢 ОТЛИЧНО - API работает стабильно\n";
        } elseif ($successRate >= 70) {
            echo "🟡 ХОРОШО - Есть незначительные проблемы\n";
        } else {
            echo "🔴 ПЛОХО - Критические проблемы в API\n";
        }
    }
}

// Запуск теста
$test = new ComprehensiveAPITest2025();
$test->run();
?> 