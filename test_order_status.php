<?php
/**
 * Тест системы статусов заказов
 * Проверяет создание заказа и изменение его статусов
 * Включает замеры скорости выполнения
 */

class OrderStatusTest {
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
    private $createdProduct = null;
    private $createdWarehouse = null;
    private $createdOffer = null;
    private $createdOrder = null;

    public function runAllTests() {
        $totalStartTime = microtime(true);
        
        echo "🚀 ТЕСТ СИСТЕМЫ СТАТУСОВ ЗАКАЗОВ\n";
        echo "==================================\n\n";

        // 1. Создание необходимых ресурсов
        $this->createTestResources();
        
        // 2. Создание заказа
        $this->testOrderCreation();
        
        // 3. Тестирование изменения статусов
        $this->testStatusChanges();
        
        // 4. Тестирование ошибок
        $this->testErrorScenarios();
        
        $totalEndTime = microtime(true);
        $this->performanceMetrics['total_time'] = round(($totalEndTime - $totalStartTime) * 1000, 2);
        
        // Вывод результатов
        $this->printResults();
    }

    private function createTestResources() {
        echo "📦 СОЗДАНИЕ ТЕСТОВЫХ РЕСУРСОВ\n";
        echo "------------------------------\n";
        
        // Создание продукта
        $productData = [
            'name' => 'Тестовый продукт для статусов',
            'vendor_article' => 'STATUS-TEST-001',
            'recommend_price' => 100.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт для тестирования статусов заказов'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание продукта'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание продукта', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProduct = $response['data']['id'];
        }
        
        // Создание склада
        $warehouseData = [
            'name' => 'Тестовый склад для статусов',
            'address' => 'ул. Тестовая, 123',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'working_hours' => '09:00-18:00'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Создание склада'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Создание склада', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouse = $response['data']['id'];
        }
        
        // Создание предложения
        if ($this->createdProduct && $this->createdWarehouse) {
            $offerData = [
                'product_id' => $this->createdProduct,
                'offer_type' => 'sale',
                'price_per_unit' => 100.00,
                'available_lots' => 10,
                'tax_nds' => 20,
                'units_per_lot' => 1,
                'warehouse_id' => $this->createdWarehouse,
                'is_public' => true,
                'max_shipping_days' => 3
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание предложения'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание предложения', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffer = $response['data']['offer_id'];
            }
        }
        
        echo "\n";
    }

    private function testOrderCreation() {
        echo "📋 СОЗДАНИЕ ЗАКАЗА\n";
        echo "------------------\n";
        
        if ($this->createdOffer) {
            $orderData = [
                'offer_id' => $this->createdOffer,
                'quantity' => 2
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['user2']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Создание заказа'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Создание заказа', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOrder = $response['data']['order_id'];
                
                // Проверка статуса по умолчанию
                if (isset($response['data']['order_status'])) {
                    $this->assertTest('Статус по умолчанию pending', $response['data']['order_status'] === 'pending', $response);
                }
            }
        }
        
        echo "\n";
    }

    private function testStatusChanges() {
        echo "🔄 ИЗМЕНЕНИЕ СТАТУСОВ ЗАКАЗА\n";
        echo "-----------------------------\n";
        
        if (!$this->createdOrder) {
            echo "❌ Заказ не создан, пропускаем тесты статусов\n\n";
            return;
        }
        
        // Продавец подтверждает заказ
        $statusData = [
            'status' => 'confirmed',
            'reason' => 'Заказ подтвержден продавцом'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса на confirmed'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса на confirmed', $response['status'] === 200, $response);
        
        // Продавец переводит в обработку
        $statusData = [
            'status' => 'processing',
            'reason' => 'Заказ взят в обработку'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса на processing'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса на processing', $response['status'] === 200, $response);
        
        // Продавец отправляет заказ
        $statusData = [
            'status' => 'shipped',
            'reason' => 'Заказ отправлен'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса на shipped'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса на shipped', $response['status'] === 200, $response);
        
        // Покупатель подтверждает доставку
        $statusData = [
            'status' => 'delivered',
            'reason' => 'Заказ получен'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['user2']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса на delivered'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса на delivered', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testErrorScenarios() {
        echo "❌ ТЕСТИРОВАНИЕ ОШИБОК\n";
        echo "------------------------\n";
        
        if (!$this->createdOrder) {
            echo "❌ Заказ не создан, пропускаем тесты ошибок\n\n";
            return;
        }
        
        // Попытка изменить статус неавторизованным пользователем
        $statusData = [
            'status' => 'confirmed',
            'reason' => 'Попытка неавторизованного изменения'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, null);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса без авторизации'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса без авторизации', $response['status'] === 401, $response);
        
        // Попытка изменить статус несуществующего заказа
        $statusData = [
            'status' => 'confirmed',
            'reason' => 'Попытка изменить несуществующий заказ'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/999999/status', $statusData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Изменение статуса несуществующего заказа'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Изменение статуса несуществующего заказа', $response['status'] === 404, $response);
        
        // Попытка установить недопустимый статус
        $statusData = [
            'status' => 'invalid_status',
            'reason' => 'Попытка установить недопустимый статус'
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Установка недопустимого статуса'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Установка недопустимого статуса', $response['status'] === 400, $response);
        
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
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ СТАТУСОВ ЗАКАЗОВ\n";
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
$test = new OrderStatusTest();
$test->runAllTests();
?> 