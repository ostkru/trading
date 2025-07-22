<?php
/**
 * Тест ролей офферов и автоматического заполнения координат
 * Проверяет логику офферов покупки/продажи и автоматическое заполнение координат
 * Включает замеры скорости выполнения
 */

class OfferRolesAndCoordinatesTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $users = [
        'user1' => [
            'name' => 'clear13808',
            'api_token' => '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27'
        ],
        'user2' => [
            'name' => 'veriy47043', 
            'api_token' => 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5'
        ]
    ];
    
    private $testResults = [];
    private $performanceMetrics = [];
    private $createdProducts = [];
    private $createdWarehouses = [];
    private $createdOffers = [];
    private $createdOrders = [];

    public function runAllTests() {
        $totalStartTime = microtime(true);
        
        echo "🚀 ТЕСТ РОЛЕЙ ОФФЕРОВ И КООРДИНАТ\n";
        echo "==================================\n\n";

        // 1. Создание тестовых ресурсов
        $this->createTestResources();
        
        // 2. Тестирование офферов продажи
        $this->testSaleOffers();
        
        // 3. Тестирование офферов покупки
        $this->testBuyOffers();
        
        // 4. Тестирование автоматического заполнения координат
        $this->testCoordinatePopulation();
        
        $totalEndTime = microtime(true);
        $this->performanceMetrics['total_time'] = round(($totalEndTime - $totalStartTime) * 1000, 2);
        
        // Вывод результатов
        $this->printResults();
    }

    private function createTestResources() {
        echo "📦 СОЗДАНИЕ ТЕСТОВЫХ РЕСУРСОВ\n";
        echo "------------------------------\n";
        
        // Создание продуктов
        $productData1 = [
            'name' => 'Продукт для продажи',
            'vendor_article' => 'SALE-TEST-001',
            'recommend_price' => 100.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт для тестирования офферов продажи'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $productData1, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта для продажи'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта для продажи', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['sale'] = $response['data']['id'];
        }
        
        $productData2 = [
            'name' => 'Продукт для покупки',
            'vendor_article' => 'BUY-TEST-001',
            'recommend_price' => 150.00,
            'brand' => 'TestBrand2',
            'category' => 'TestCategory2',
            'description' => 'Продукт для тестирования офферов покупки'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $productData2, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта для покупки'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта для покупки', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['buy'] = $response['data']['id'];
        }
        
        // Создание складов
        $warehouseData1 = [
            'name' => 'Склад для продажи',
            'address' => 'ул. Продажная, 1',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'working_hours' => '09:00-18:00'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData1, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание склада для продажи'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание склада для продажи', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['sale'] = $response['data']['id'];
        }
        
        $warehouseData2 = [
            'name' => 'Склад для покупки',
            'address' => 'ул. Покупная, 2',
            'latitude' => 55.7600,
            'longitude' => 37.6200,
            'working_hours' => '10:00-19:00'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData2, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание склада для покупки'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание склада для покупки', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['buy'] = $response['data']['id'];
        }
        
        echo "\n";
    }

    private function testSaleOffers() {
        echo "💰 ТЕСТИРОВАНИЕ ОФФЕРОВ ПРОДАЖИ\n";
        echo "--------------------------------\n";
        
        if (isset($this->createdProducts['sale']) && isset($this->createdWarehouses['sale'])) {
            // Создание оффера продажи
            $offerData = [
                'product_id' => $this->createdProducts['sale'],
                'offer_type' => 'sale',
                'price_per_unit' => 100.00,
                'available_lots' => 10,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['sale'],
                'is_public' => true,
                'max_shipping_days' => 3
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание оффера продажи'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание оффера продажи', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['sale'] = $response['data']['offer_id'];
                
                // Проверка автоматического заполнения координат
                if (isset($response['data']['latitude']) && isset($response['data']['longitude'])) {
                    $this->assertTest('Координаты заполнены автоматически', 
                        $response['data']['latitude'] == 55.7558 && $response['data']['longitude'] == 37.6176, $response);
                }
            }
            
            // Создание заказа на оффер продажи
            if (isset($this->createdOffers['sale'])) {
                $orderData = [
                    'offer_id' => $this->createdOffers['sale'],
                    'quantity' => 2
                ];
                
                $startTime = microtime(true);
                $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
                $endTime = microtime(true);
                $this->performanceMetrics['Создание заказа на оффер продажи'] = round(($endTime - $startTime) * 1000, 2);
                $this->assertTest('Создание заказа на оффер продажи', $response['status'] === 201, $response);
                if ($response['status'] === 201) {
                    $this->createdOrders['sale'] = $response['data']['order_id'];
                    
                    // Проверка ролей в заказе продажи
                    if (isset($response['data']['initiator_user_id']) && isset($response['data']['counterparty_user_id'])) {
                        $this->assertTest('User2 является покупателем в заказе продажи', 
                            $response['data']['initiator_user_id'] == 2, $response);
                        $this->assertTest('User1 является продавцом в заказе продажи', 
                            $response['data']['counterparty_user_id'] == 1, $response);
                    }
                    
                    if (isset($response['data']['order_type'])) {
                        $this->assertTest('Тип заказа продажи', $response['data']['order_type'] === 'buy', $response);
                    }
                }
            }
        }
        
        echo "\n";
    }

    private function testBuyOffers() {
        echo "🛒 ТЕСТИРОВАНИЕ ОФФЕРОВ ПОКУПКИ\n";
        echo "--------------------------------\n";
        
        if (isset($this->createdProducts['buy']) && isset($this->createdWarehouses['buy'])) {
            // Создание оффера покупки
            $offerData = [
                'product_id' => $this->createdProducts['buy'],
                'offer_type' => 'buy',
                'price_per_unit' => 150.00,
                'available_lots' => 5,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouses['buy'],
                'is_public' => true,
                'max_shipping_days' => 5
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание оффера покупки'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание оффера покупки', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['buy'] = $response['data']['offer_id'];
                
                // Проверка автоматического заполнения координат
                if (isset($response['data']['latitude']) && isset($response['data']['longitude'])) {
                    $this->assertTest('Координаты заполнены автоматически для покупки', 
                        $response['data']['latitude'] == 55.7600 && $response['data']['longitude'] == 37.6200, $response);
                }
            }
            
            // Создание заказа на оффер покупки
            if (isset($this->createdOffers['buy'])) {
                $orderData = [
                    'offer_id' => $this->createdOffers['buy'],
                    'quantity' => 1
                ];
                
                $startTime = microtime(true);
                $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user1']['api_token']);
                $endTime = microtime(true);
                $this->performanceMetrics['Создание заказа на оффер покупки'] = round(($endTime - $startTime) * 1000, 2);
                $this->assertTest('Создание заказа на оффер покупки', $response['status'] === 201, $response);
                if ($response['status'] === 201) {
                    $this->createdOrders['buy'] = $response['data']['order_id'];
                    
                    // Проверка ролей в заказе покупки
                    if (isset($response['data']['initiator_user_id']) && isset($response['data']['counterparty_user_id'])) {
                        $this->assertTest('User2 является покупателем в заказе покупки', 
                            $response['data']['initiator_user_id'] == 2, $response);
                        $this->assertTest('User1 является продавцом в заказе покупки', 
                            $response['data']['counterparty_user_id'] == 1, $response);
                    }
                    
                    if (isset($response['data']['order_type'])) {
                        $this->assertTest('Тип заказа покупки', $response['data']['order_type'] === 'sell', $response);
                    }
                }
            }
        }
        
        echo "\n";
    }

    private function testCoordinatePopulation() {
        echo "📍 ТЕСТИРОВАНИЕ АВТОМАТИЧЕСКОГО ЗАПОЛНЕНИЯ КООРДИНАТ\n";
        echo "----------------------------------------------------\n";
        
        // Обновление оффера с изменением склада
        if (isset($this->createdOffers['sale']) && isset($this->createdWarehouses['buy'])) {
            $updateData = [
                'warehouse_id' => $this->createdWarehouses['buy']
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['sale'], $updateData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Обновление координат при смене склада'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Обновление координат при смене склада', $response['status'] === 200, $response);
            
            if ($response['status'] === 200 && isset($response['data']['latitude']) && isset($response['data']['longitude'])) {
                $this->assertTest('Координаты обновлены после смены склада', 
                    $response['data']['latitude'] == 55.7600 && $response['data']['longitude'] == 37.6200, $response);
            }
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
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ РОЛЕЙ ОФФЕРОВ И КООРДИНАТ\n";
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
$test = new OfferRolesAndCoordinatesTest();
$test->runAllTests(); 