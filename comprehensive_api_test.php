<?php
/**
 * Комплексный мультитест API PortalData
 * Проверяет все методы с разных пользователей и сценариев
 */

class ComprehensiveAPITest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $users = [
        'user1' => [
            'name' => 'clear13808',
            'api_token' => '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27'
        ],
        'user2' => [
            'name' => 'veriy47043', 
            'api_token' => 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5'
        ],
        'user3' => [
            'name' => '5ij48m4gcm',
            'api_token' => '00601582c3163466e0fece95d8e2315cb1c66814066ad8e0566d2813614d9001'
        ]
    ];
    
    private $testResults = [];
    private $createdProducts = [];
    private $createdOffers = [];
    private $createdOrders = [];
    private $createdWarehouses = [];

    public function runAllTests() {
        echo "🚀 ЗАПУСК КОМПЛЕКСНОГО ТЕСТИРОВАНИЯ API\n";
        echo "==========================================\n\n";

        // 1. Базовые проверки
        $this->testBasicEndpoints();
        
        // 2. Тестирование продуктов
        $this->testProducts();
        
        // 3. Тестирование складов
        $this->testWarehouses();
        
        // 4. Тестирование предложений
        $this->testOffers();
        
        // 5. Тестирование заказов
        $this->testOrders();
        
        // 6. Тестирование публичных маршрутов
        $this->testPublicRoutes();
        
        // 7. Тестирование ошибок и валидации
        $this->testErrorScenarios();
        
        // 8. Тестирование безопасности
        $this->testSecurityScenarios();
        
        // Вывод результатов
        $this->printResults();
    }

    private function testBasicEndpoints() {
        echo "📋 1. БАЗОВЫЕ ПРОВЕРКИ\n";
        echo "------------------------\n";
        
        // Проверка основного endpoint
        $response = $this->makeRequest('GET', '/', null, null);
        $this->assertTest('Основной endpoint', $response['status'] === 200, $response);
        
        // Проверка доступности API
        $response = $this->makeRequest('GET', '/products', null, $this->users['user1']['api_token']);
        $this->assertTest('API доступен', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testProducts() {
        echo "📦 2. ТЕСТИРОВАНИЕ ПРОДУКТОВ\n";
        echo "------------------------------\n";
        
        // Создание продукта пользователем 1
        $productData = [
            'name' => 'Тестовый продукт User1',
            'vendor_article' => 'TEST-USER1-001',
            'recommend_price' => 150.50,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Описание тестового продукта от User1'
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user1']['api_token']);
        $this->assertTest('Создание продукта User1', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['user1'] = $response['data']['id'];
        }
        
        // Создание продукта пользователем 2
        $productData = [
            'name' => 'Тестовый продукт User2',
            'vendor_article' => 'TEST-USER2-001',
            'recommend_price' => 200.00,
            'brand' => 'User2Brand',
            'category' => 'User2Category',
            'description' => 'Описание тестового продукта от User2'
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user2']['api_token']);
        $this->assertTest('Создание продукта User2', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['user2'] = $response['data']['id'];
        }
        
        // Получение списка продуктов
        $response = $this->makeRequest('GET', '/products', null, $this->users['user1']['api_token']);
        $this->assertTest('Получение списка продуктов', $response['status'] === 200, $response);
        
        // Получение конкретного продукта
        if (isset($this->createdProducts['user1'])) {
            $response = $this->makeRequest('GET', '/products/' . $this->createdProducts['user1'], null, $this->users['user1']['api_token']);
            $this->assertTest('Получение конкретного продукта', $response['status'] === 200, $response);
        }
        
        // Обновление продукта
        if (isset($this->createdProducts['user1'])) {
            $updateData = [
                'name' => 'Обновленный продукт User1',
                'recommend_price' => 175.00
            ];
            $response = $this->makeRequest('PUT', '/products/' . $this->createdProducts['user1'], $updateData, $this->users['user1']['api_token']);
            $this->assertTest('Обновление продукта', $response['status'] === 200, $response);
        }
        
        // Batch создание продуктов
        $batchData = [
            'products' => [
                [
                    'name' => 'Batch Product 1',
                    'vendor_article' => 'BATCH-001',
                    'recommend_price' => 100.00,
                    'brand' => 'BatchBrand',
                    'category' => 'BatchCategory',
                    'description' => 'Первый продукт из batch'
                ],
                [
                    'name' => 'Batch Product 2',
                    'vendor_article' => 'BATCH-002',
                    'recommend_price' => 150.00,
                    'brand' => 'BatchBrand',
                    'category' => 'BatchCategory',
                    'description' => 'Второй продукт из batch'
                ]
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products/batch', $batchData, $this->users['user1']['api_token']);
        $this->assertTest('Batch создание продуктов', $response['status'] === 201, $response);
        
        echo "\n";
    }

    private function testWarehouses() {
        echo "🏭 3. ТЕСТИРОВАНИЕ СКЛАДОВ\n";
        echo "----------------------------\n";
        
        // Создание склада пользователем 1
        $warehouseData = [
            'name' => 'Склад User1',
            'address' => 'Москва, ул. Тестовая, 1',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'working_hours' => '09:00-18:00'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['user1']['api_token']);
        $this->assertTest('Создание склада User1', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['user1'] = $response['data']['id'];
        }
        
        // Создание склада пользователем 2
        $warehouseData = [
            'name' => 'Склад User2',
            'address' => 'СПб, ул. Тестовая, 2',
            'latitude' => 59.9311,
            'longitude' => 30.3609,
            'working_hours' => '10:00-19:00'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['user2']['api_token']);
        $this->assertTest('Создание склада User2', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['user2'] = $response['data']['id'];
        }
        
        // Получение списка складов
        $response = $this->makeRequest('GET', '/warehouses', null, $this->users['user1']['api_token']);
        $this->assertTest('Получение списка складов', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testOffers() {
        echo "💰 4. ТЕСТИРОВАНИЕ ПРЕДЛОЖЕНИЙ\n";
        echo "--------------------------------\n";
        
        // Создание предложения пользователем 1
        if (isset($this->createdProducts['user1']) && isset($this->createdWarehouses['user1'])) {
            $offerData = [
                'product_id' => $this->createdProducts['user1'],
                'offer_type' => 'sale',
                'price_per_unit' => 180.00,
                'available_lots' => 10,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['user1'],
                'is_public' => true,
                'max_shipping_days' => 5
            ];
            
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user1']['api_token']);
            $this->assertTest('Создание предложения User1', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['user1'] = $response['data']['offer_id'];
            }
        }
        
        // Создание предложения пользователем 2
        if (isset($this->createdProducts['user2']) && isset($this->createdWarehouses['user2'])) {
            $offerData = [
                'product_id' => $this->createdProducts['user2'],
                'offer_type' => 'sale',
                'price_per_unit' => 250.00,
                'available_lots' => 5,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['user2'],
                'is_public' => false,
                'max_shipping_days' => 3
            ];
            
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user2']['api_token']);
            $this->assertTest('Создание предложения User2', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['user2'] = $response['data']['offer_id'];
            }
        }
        
        // Получение списка предложений пользователя
        $response = $this->makeRequest('GET', '/offers', null, $this->users['user1']['api_token']);
        $this->assertTest('Получение списка предложений', $response['status'] === 200, $response);
        
        // Обновление предложения
        if (isset($this->createdOffers['user1'])) {
            $updateData = [
                'price_per_unit' => 190.00,
                'available_lots' => 8
            ];
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['user1'], $updateData, $this->users['user1']['api_token']);
            $this->assertTest('Обновление предложения', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function testOrders() {
        echo "📋 5. ТЕСТИРОВАНИЕ ЗАКАЗОВ\n";
        echo "-----------------------------\n";
        
        // Создание заказа пользователем 2 на предложение пользователя 1
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 2
            ];
            
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
            $this->assertTest('Создание заказа User2 на предложение User1', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOrders['user2_on_user1'] = $response['data']['order_id'];
            }
        }
        
        // Создание заказа пользователем 3 на предложение пользователя 1
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 1
            ];
            
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user3']['api_token']);
            $this->assertTest('Создание заказа User3 на предложение User1', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOrders['user3_on_user1'] = $response['data']['order_id'];
            }
        }
        
        // Получение заказа
        if (isset($this->createdOrders['user2_on_user1'])) {
            $response = $this->makeRequest('GET', '/orders/' . $this->createdOrders['user2_on_user1'], null, $this->users['user2']['api_token']);
            $this->assertTest('Получение заказа', $response['status'] === 200, $response);
        }
        
        // Получение списка заказов
        $response = $this->makeRequest('GET', '/orders', null, $this->users['user2']['api_token']);
        $this->assertTest('Получение списка заказов', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testPublicRoutes() {
        echo "🌐 6. ТЕСТИРОВАНИЕ ПУБЛИЧНЫХ МАРШРУТОВ\n";
        echo "----------------------------------------\n";
        
        // Публичные предложения без авторизации
        $response = $this->makeRequest('GET', '/offers/public', null, null);
        $this->assertTest('Публичные предложения без авторизации', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testErrorScenarios() {
        echo "❌ 7. ТЕСТИРОВАНИЕ ОШИБОК И ВАЛИДАЦИИ\n";
        echo "----------------------------------------\n";
        
        // Попытка создания продукта без авторизации
        $productData = ['name' => 'Test Product'];
        $response = $this->makeRequest('POST', '/products', $productData, null);
        $this->assertTest('Создание продукта без авторизации (должно быть 401)', $response['status'] === 401, $response);
        
        // Попытка создания продукта с некорректными данными
        $productData = ['name' => '']; // Пустое имя
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user1']['api_token']);
        $this->assertTest('Создание продукта с пустым именем', $response['status'] === 400, $response);
        
        // Попытка получения несуществующего продукта
        $response = $this->makeRequest('GET', '/products/99999', null, $this->users['user1']['api_token']);
        $this->assertTest('Получение несуществующего продукта', $response['status'] === 404, $response);
        
        // Попытка создания заказа на собственное предложение
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 1
            ];
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user1']['api_token']);
            $this->assertTest('Заказ на собственное предложение (должно быть запрещено)', $response['status'] === 400, $response);
        }
        
        // Попытка создания заказа с некорректными данными
        $orderData = ['offer_id' => 0, 'quantity' => 0];
        $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user1']['api_token']);
        $this->assertTest('Создание заказа с некорректными данными', $response['status'] === 400, $response);
        
        // Попытка заказа больше лотов, чем доступно
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 999 // Очень много
            ];
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
            $this->assertTest('Заказ больше лотов, чем доступно', $response['status'] === 400, $response);
        }
        
        echo "\n";
    }

    private function testSecurityScenarios() {
        echo "🔒 8. ТЕСТИРОВАНИЕ БЕЗОПАСНОСТИ\n";
        echo "----------------------------------\n";
        
        // Попытка доступа к чужим данным
        if (isset($this->createdProducts['user1'])) {
            $response = $this->makeRequest('GET', '/products/' . $this->createdProducts['user1'], null, $this->users['user2']['api_token']);
            $this->assertTest('Доступ к чужому продукту (должен быть разрешен для чтения)', $response['status'] === 200, $response);
        }
        
        // Попытка обновления чужого продукта
        if (isset($this->createdProducts['user1'])) {
            $updateData = ['name' => 'Взломанный продукт'];
            $response = $this->makeRequest('PUT', '/products/' . $this->createdProducts['user1'], $updateData, $this->users['user2']['api_token']);
            $this->assertTest('Обновление чужого продукта (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка удаления чужого продукта
        if (isset($this->createdProducts['user1'])) {
            $response = $this->makeRequest('DELETE', '/products/' . $this->createdProducts['user1'], null, $this->users['user2']['api_token']);
            $this->assertTest('Удаление чужого продукта (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка обновления чужого предложения
        if (isset($this->createdOffers['user1'])) {
            $updateData = ['price_per_unit' => 999.99];
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['user1'], $updateData, $this->users['user2']['api_token']);
            $this->assertTest('Обновление чужого предложения (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка удаления чужого предложения
        if (isset($this->createdOffers['user1'])) {
            $response = $this->makeRequest('DELETE', '/offers/' . $this->createdOffers['user1'], null, $this->users['user2']['api_token']);
            $this->assertTest('Удаление чужого предложения (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка обновления чужого склада
        if (isset($this->createdWarehouses['user1'])) {
            $updateData = ['name' => 'Взломанный склад'];
            $response = $this->makeRequest('PUT', '/warehouses/' . $this->createdWarehouses['user1'], $updateData, $this->users['user2']['api_token']);
            $this->assertTest('Обновление чужого склада (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка удаления чужого склада
        if (isset($this->createdWarehouses['user1'])) {
            $response = $this->makeRequest('DELETE', '/warehouses/' . $this->createdWarehouses['user1'], null, $this->users['user2']['api_token']);
            $this->assertTest('Удаление чужого склада (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Попытка доступа с неверным API ключом
        $response = $this->makeRequest('GET', '/products', null, 'invalid_api_key');
        $this->assertTest('Доступ с неверным API ключом', $response['status'] === 401, $response);
        
        echo "\n";
    }

    private function makeRequest($method, $endpoint, $data = null, $apiToken = null) {
        $url = $this->baseUrl . $endpoint;
        
        $headers = ['Content-Type: application/json'];
        if ($apiToken) {
            $headers[] = 'Authorization: Bearer ' . $apiToken;
        }
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        
        if ($data && in_array($method, ['POST', 'PUT'])) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
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
        $status = $condition ? '✅ ПРОЙДЕН' : '❌ ПРОВАЛЕН';
        echo sprintf("%-60s %s\n", $testName, $status);
        
        if (!$condition) {
            echo "   Ошибка: " . ($response['data']['error'] ?? 'Неизвестная ошибка') . "\n";
        }
        
        $this->testResults[] = [
            'name' => $testName,
            'passed' => $condition,
            'response' => $response
        ];
    }

    private function printResults() {
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ\n";
        echo "============================\n";
        
        $passed = 0;
        $total = count($this->testResults);
        
        foreach ($this->testResults as $result) {
            if ($result['passed']) {
                $passed++;
            }
        }
        
        echo "Всего тестов: $total\n";
        echo "Пройдено: $passed\n";
        echo "Провалено: " . ($total - $passed) . "\n";
        echo "Процент успеха: " . round(($passed / $total) * 100, 2) . "%\n\n";
        
        if ($passed === $total) {
            echo "🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!\n";
        } else {
            echo "⚠️  НЕКОТОРЫЕ ТЕСТЫ ПРОВАЛЕНЫ\n";
        }
    }
}

// Запуск тестирования
$test = new ComprehensiveAPITest();
$test->runAllTests();
?> 