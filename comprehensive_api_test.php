<?php
/**
 * Комплексный мультитест API PortalData
 * Проверяет все методы с разных пользователей и сценариев
 * Включает замеры скорости выполнения
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
    private $performanceMetrics = [];

    public function runAllTests() {
        $totalStartTime = microtime(true);
        
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
        
        $totalEndTime = microtime(true);
        $this->performanceMetrics['total_time'] = round(($totalEndTime - $totalStartTime) * 1000, 2);
        
        // Вывод результатов
        $this->printResults();
    }

    private function testBasicEndpoints() {
        echo "📋 1. БАЗОВЫЕ ПРОВЕРКИ\n";
        echo "------------------------\n";
        
        // Проверка основного endpoint
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/', null, null);
        $endTime = microtime(true);
        $this->performanceMetrics['Основной endpoint'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Основной endpoint', $response['status'] === 200, $response);
        
        // Проверка доступности API
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/products', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['API доступен'] = round(($endTime - $startTime) * 1000, 2);
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
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта User1'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта User1', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['user1'] = $response['data']['id'];
        }
        
        // Создание продукта пользователем 2
        $productData = [
            'name' => 'Тестовый продукт User2',
            'vendor_article' => 'TEST-USER2-001',
            'recommend_price' => 200.75,
            'brand' => 'TestBrand2',
            'category' => 'TestCategory2',
            'description' => 'Описание тестового продукта от User2'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта User2'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта User2', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['user2'] = $response['data']['id'];
        }
        
        // Получение списка продуктов
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/products', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение списка продуктов'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение списка продуктов', $response['status'] === 200, $response);
        
        // Получение продукта по ID
        if (isset($this->createdProducts['user1'])) {
            $startTime = microtime(true);
            $response = $this->makeRequest('GET', '/products/' . $this->createdProducts['user1'], null, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Получение продукта по ID'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Получение продукта по ID', $response['status'] === 200, $response);
        }
        
        // Обновление продукта
        if (isset($this->createdProducts['user1'])) {
            $updateData = [
                'name' => 'Обновленный продукт User1',
                'recommend_price' => 175.25
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/products/' . $this->createdProducts['user1'], $updateData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Обновление продукта'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Обновление продукта', $response['status'] === 200, $response);
        }
        
        // Создание продукта с пустым именем (должно быть запрещено)
        $invalidProductData = [
            'name' => '',
            'vendor_article' => 'TEST-EMPTY-001',
            'recommend_price' => 100.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $invalidProductData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта с пустым именем'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта с пустым именем', $response['status'] === 400, $response);
        
        // Обновление чужого продукта (должно быть запрещено)
        if (isset($this->createdProducts['user1']) && isset($this->createdProducts['user2'])) {
            $updateData = [
                'name' => 'Попытка обновить чужой продукт',
                'recommend_price' => 999.99
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/products/' . $this->createdProducts['user1'], $updateData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Обновление чужого продукта (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Обновление чужого продукта (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Удаление чужого продукта (должно быть запрещено)
        if (isset($this->createdProducts['user1']) && isset($this->createdProducts['user2'])) {
            $startTime = microtime(true);
            $response = $this->makeRequest('DELETE', '/products/' . $this->createdProducts['user1'], null, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Удаление чужого продукта (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Удаление чужого продукта (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        echo "\n";
    }

    private function testWarehouses() {
        echo "🏭 3. ТЕСТИРОВАНИЕ СКЛАДОВ\n";
        echo "----------------------------\n";
        
        // Создание склада User1
        $warehouseData = [
            'name' => 'Склад User1',
            'address' => 'ул. Тестовая, 1',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'working_hours' => '09:00-18:00'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание склада User1'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание склада User1', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['user1'] = $response['data']['id'];
        }
        
        // Создание склада User2
        $warehouseData = [
            'name' => 'Склад User2',
            'address' => 'ул. Тестовая, 2',
            'latitude' => 55.7600,
            'longitude' => 37.6200,
            'working_hours' => '10:00-19:00'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание склада User2'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание склада User2', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['user2'] = $response['data']['id'];
        }
        
        // Получение списка складов
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/warehouses', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение списка складов'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение списка складов', $response['status'] === 200, $response);
        
        // Обновление чужого склада (должно быть запрещено)
        if (isset($this->createdWarehouses['user1']) && isset($this->createdWarehouses['user2'])) {
            $updateData = [
                'name' => 'Попытка обновить чужой склад',
                'address' => 'ул. Взломанная, 999'
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/warehouses/' . $this->createdWarehouses['user1'], $updateData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Обновление чужого склада (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Обновление чужого склада (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Удаление чужого склада (должно быть запрещено)
        if (isset($this->createdWarehouses['user1']) && isset($this->createdWarehouses['user2'])) {
            $startTime = microtime(true);
            $response = $this->makeRequest('DELETE', '/warehouses/' . $this->createdWarehouses['user1'], null, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Удаление чужого склада (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Удаление чужого склада (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        echo "\n";
    }

    private function testOffers() {
        echo "📋 4. ТЕСТИРОВАНИЕ ПРЕДЛОЖЕНИЙ\n";
        echo "--------------------------------\n";
        
        // Создание предложения
        if (isset($this->createdProducts['user1']) && isset($this->createdWarehouses['user1'])) {
            $offerData = [
                'product_id' => $this->createdProducts['user1'],
                'offer_type' => 'sale',
                'price_per_unit' => 100.00,
                'available_lots' => 10,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['user1'],
                'is_public' => true,
                'max_shipping_days' => 3
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание предложения'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание предложения', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['user1'] = $response['data']['offer_id'];
            }
        }
        
        // Создание предложения на покупку
        if (isset($this->createdProducts['user2']) && isset($this->createdWarehouses['user2'])) {
            $offerData = [
                'product_id' => $this->createdProducts['user2'],
                'offer_type' => 'buy',
                'price_per_unit' => 150.00,
                'available_lots' => 5,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['user2'],
                'is_public' => true,
                'max_shipping_days' => 5
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание предложения на покупку'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание предложения на покупку', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['user2'] = $response['data']['offer_id'];
            }
        }
        
        // Получение списка предложений пользователя
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение списка предложений'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение списка предложений', $response['status'] === 200, $response);
        
        // Тестирование фильтрации офферов
        echo "   🔍 Тестирование фильтрации офферов:\n";
        
        // Фильтр "my" - только мои офферы
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers?filter=my', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Фильтр офферов: my'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Фильтр офферов: my (только мои)', $response['status'] === 200, $response);
        
        // Фильтр "others" - чужие офферы
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers?filter=others', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Фильтр офферов: others'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Фильтр офферов: others (чужие)', $response['status'] === 200, $response);
        
        // Фильтр "all" - все офферы
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers?filter=all', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Фильтр офферов: all'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Фильтр офферов: all (все)', $response['status'] === 200, $response);
        
        // Без параметра filter (должен вернуть мои офферы по умолчанию)
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Фильтр офферов: по умолчанию'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Фильтр офферов: по умолчанию (my)', $response['status'] === 200, $response);
        
        // Неверный фильтр (должен вернуть мои офферы по умолчанию)
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers?filter=invalid', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Фильтр офферов: неверный'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Фильтр офферов: неверный (должен вернуть my)', $response['status'] === 200, $response);
        
        // Получение публичных предложений
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers/public', null, null);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение публичных предложений'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение публичных предложений', $response['status'] === 200, $response);
        
        // Обновление чужого предложения (должно быть запрещено)
        if (isset($this->createdOffers['user1']) && isset($this->createdOffers['user2'])) {
            $updateData = [
                'price_per_unit' => 999.99,
                'available_lots' => 999
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['user1'], $updateData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Обновление чужого предложения (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Обновление чужого предложения (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        // Удаление чужого предложения (должно быть запрещено)
        if (isset($this->createdOffers['user1']) && isset($this->createdOffers['user2'])) {
            $startTime = microtime(true);
            $response = $this->makeRequest('DELETE', '/offers/' . $this->createdOffers['user1'], null, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Удаление чужого предложения (должно быть запрещено)'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Удаление чужого предложения (должно быть запрещено)', $response['status'] === 403, $response);
        }
        
        echo "\n";
    }

    private function testOrders() {
        echo "📦 5. ТЕСТИРОВАНИЕ ЗАКАЗОВ\n";
        echo "----------------------------\n";
        
        // Создание заказа
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 2
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание заказа'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание заказа', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOrders['user2'] = $response['data']['order_id'];
            }
        }
        
        // Получение списка заказов
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/orders', null, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение списка заказов'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение списка заказов', $response['status'] === 200, $response);
        
        // Получение заказа по ID
        if (isset($this->createdOrders['user2'])) {
            $startTime = microtime(true);
            $response = $this->makeRequest('GET', '/orders/' . $this->createdOrders['user2'], null, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Получение заказа по ID'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Получение заказа по ID', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function testPublicRoutes() {
        echo "🌐 6. ТЕСТИРОВАНИЕ ПУБЛИЧНЫХ МАРШРУТОВ\n";
        echo "----------------------------------------\n";
        
        // Проверка публичных предложений без авторизации
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/offers/public', null, null);
        $endTime = microtime(true);
        $this->performanceMetrics['Публичные предложения без авторизации'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Публичные предложения без авторизации', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testErrorScenarios() {
        echo "❌ 7. ТЕСТИРОВАНИЕ ОШИБОК И ВАЛИДАЦИИ\n";
        echo "----------------------------------------\n";
        
        // Попытка доступа без API ключа
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/products', null, null);
        $endTime = microtime(true);
        $this->performanceMetrics['Доступ без API ключа'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Доступ без API ключа', $response['status'] === 401, $response);
        
        // Попытка доступа с неверным API ключом
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/products', null, 'invalid_token');
        $endTime = microtime(true);
        $this->performanceMetrics['Доступ с неверным API ключом'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Доступ с неверным API ключом', $response['status'] === 401, $response);
        
        // Попытка получить несуществующий ресурс
        $startTime = microtime(true);
        $response = $this->makeRequest('GET', '/products/999999', null, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Получение несуществующего ресурса'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Получение несуществующего ресурса', $response['status'] === 404, $response);
        
        // Попытка создать заказ на несуществующее предложение
        $orderData = [
            'offer_id' => 999999,
            'quantity' => 1
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание заказа на несуществующее предложение'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание заказа на несуществующее предложение', $response['status'] === 404, $response);
        
        echo "\n";
    }

    private function testSecurityScenarios() {
        echo "🔒 8. ТЕСТИРОВАНИЕ БЕЗОПАСНОСТИ\n";
        echo "--------------------------------\n";
        
        // Попытка создать заказ на свое предложение
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 1
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание заказа на свое предложение'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание заказа на свое предложение', $response['status'] === 400, $response);
        }
        
        // Попытка создать заказ с превышением доступного количества
        if (isset($this->createdOffers['user1'])) {
            $orderData = [
                'offer_id' => $this->createdOffers['user1'],
                'quantity' => 999999
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание заказа с превышением количества'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание заказа с превышением количества', $response['status'] === 400, $response);
        }
        
        echo "\n";
    }

    private function makeRequest($method, $endpoint, $data = null, $apiToken = null) {
        $url = $this->baseUrl . $endpoint;
        
        $headers = [
            'Content-Type: application/json',
            'Accept: application/json'
        ];
        
        if ($apiToken) {
            $headers[] = 'Authorization: Bearer ' . $apiToken;
        }
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        curl_setopt($ch, CURLOPT_TIMEOUT, 30);
        
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
        $result = $condition ? '✅ ПРОЙДЕН' : '❌ ПРОВАЛЕН';
        $status = $response['status'];
        $message = isset($response['data']['error']) ? $response['data']['error'] : '';
        
        echo sprintf("%-50s %s (HTTP %d)", $testName, $result, $status);
        if ($message) {
            echo " - $message";
        }
        echo "\n";
        
        $this->testResults[] = [
            'name' => $testName,
            'passed' => $condition,
            'status' => $status,
            'message' => $message
        ];
    }

    private function printResults() {
        echo "\n" . str_repeat("=", 80) . "\n";
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ\n";
        echo str_repeat("=", 80) . "\n\n";
        
        $totalTests = count($this->testResults);
        $passedTests = count(array_filter($this->testResults, function($test) {
            return $test['passed'];
        }));
        $failedTests = $totalTests - $passedTests;
        $successRate = round(($passedTests / $totalTests) * 100, 2);
        
        echo "📈 ОБЩАЯ СТАТИСТИКА:\n";
        echo "   Всего тестов: $totalTests\n";
        echo "   Пройдено: $passedTests\n";
        echo "   Провалено: $failedTests\n";
        echo "   Успешность: $successRate%\n";
        echo "   Общее время выполнения: {$this->performanceMetrics['total_time']} мс\n\n";
        
        echo "⚡ МЕТРИКИ ПРОИЗВОДИТЕЛЬНОСТИ:\n";
        echo str_repeat("-", 80) . "\n";
        foreach ($this->performanceMetrics as $testName => $time) {
            if ($testName !== 'total_time') {
                echo sprintf("%-50s %6.2f мс\n", $testName, $time);
            }
        }
        echo str_repeat("-", 80) . "\n";
        
        if ($failedTests > 0) {
            echo "\n❌ ПРОВАЛЕННЫЕ ТЕСТЫ:\n";
            echo str_repeat("-", 80) . "\n";
            foreach ($this->testResults as $test) {
                if (!$test['passed']) {
                    echo sprintf("• %s (HTTP %d): %s\n", $test['name'], $test['status'], $test['message']);
                }
            }
        }
        
        echo "\n" . str_repeat("=", 80) . "\n";
    }
}

// Запуск тестов
$test = new ComprehensiveAPITest();
$test->runAllTests();
?> 