<?php
/**
 * Тест системы статусов заказов
 */

class OrderStatusTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $users = [
        'buyer' => [
            'name' => 'clear13808',
            'api_token' => '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27'
        ],
        'seller' => [
            'name' => 'veriy47043', 
            'api_token' => 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5'
        ]
    ];
    
    private $testResults = [];
    private $createdOrder = null;
    private $createdOffer = null;

    public function runAllTests() {
        echo "🚀 ТЕСТИРОВАНИЕ СИСТЕМЫ СТАТУСОВ ЗАКАЗОВ\n";
        echo "==========================================\n\n";

        // 1. Создание предложения продавцом
        $this->testCreateOffer();
        
        // 2. Создание заказа покупателем
        $this->testCreateOrder();
        
        // 3. Тестирование изменения статусов
        $this->testStatusChanges();
        
        // 4. Тестирование прав доступа
        $this->testAccessRights();
        
        // Вывод результатов
        $this->printResults();
    }

    private function testCreateOffer() {
        echo "📦 1. СОЗДАНИЕ ПРЕДЛОЖЕНИЯ\n";
        echo "---------------------------\n";
        
        // Сначала создаем продукт
        $productData = [
            'name' => 'Тестовый продукт для статусов',
            'vendor_article' => 'STATUS-TEST-001',
            'recommend_price' => 100.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт для тестирования статусов'
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['seller']['api_token']);
        $this->assertTest('Создание продукта', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $productId = $response['data']['id'];
            
            // Сначала создаем склад
            $warehouseData = [
                'name' => 'Тестовый склад',
                'address' => 'Москва, ул. Тестовая, 1',
                'latitude' => 55.7558,
                'longitude' => 37.6176,
                'working_hours' => '9:00-18:00'
            ];
            
            $response = $this->makeRequest('POST', '/warehouses', $warehouseData, $this->users['seller']['api_token']);
            $this->assertTest('Создание склада', $response['status'] === 201, $response);
            
            if ($response['status'] === 201) {
                $warehouseId = $response['data']['id'];
                
                // Создаем предложение
                $offerData = [
                    'product_id' => $productId,
                    'warehouse_id' => $warehouseId,
                    'offer_type' => 'sell',
                    'price_per_unit' => 100.00,
                    'units_per_lot' => 10,
                    'available_lots' => 5,
                    'tax_nds' => 20.00,
                    'is_public' => true
                ];
            }
            
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['seller']['api_token']);
            $this->assertTest('Создание предложения', $response['status'] === 201, $response);
            
            if ($response['status'] === 201) {
                $this->createdOffer = $response['data']['offer_id'];
            }
        }
        
        echo "\n";
    }

    private function testCreateOrder() {
        echo "📋 2. СОЗДАНИЕ ЗАКАЗА\n";
        echo "----------------------\n";
        
        if (!$this->createdOffer) {
            echo "❌ Пропуск: предложение не создано\n\n";
            return;
        }
        
        $orderData = [
            'offer_id' => $this->createdOffer,
            'quantity' => 1
        ];
        
        $response = $this->makeRequest('POST', '/orders', $orderData, $this->users['buyer']['api_token']);
        $this->assertTest('Создание заказа', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdOrder = $response['data']['order_id'];
            
            // Проверяем, что статус по умолчанию - pending
            $this->assertTest('Статус по умолчанию pending', 
                isset($response['data']['order_status']) && $response['data']['order_status'] === 'pending', $response);
        }
        
        echo "\n";
    }

    private function testStatusChanges() {
        echo "🔄 3. ТЕСТИРОВАНИЕ ИЗМЕНЕНИЯ СТАТУСОВ\n";
        echo "----------------------------------------\n";
        
        if (!$this->createdOrder) {
            echo "❌ Пропуск: заказ не создан\n\n";
            return;
        }
        
        // Продавец подтверждает заказ
        $statusData = ['status' => 'confirmed', 'reason' => 'Заказ подтвержден'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['seller']['api_token']);
        $this->assertTest('Продавец подтверждает заказ', $response['status'] === 200, $response);
        
        // Продавец начинает обработку
        $statusData = ['status' => 'processing', 'reason' => 'Заказ в обработке'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['seller']['api_token']);
        $this->assertTest('Продавец начинает обработку', $response['status'] === 200, $response);
        
        // Продавец отправляет товар
        $statusData = ['status' => 'shipped', 'reason' => 'Товар отправлен'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['seller']['api_token']);
        $this->assertTest('Продавец отправляет товар', $response['status'] === 200, $response);
        
        // Покупатель подтверждает получение
        $statusData = ['status' => 'delivered', 'reason' => 'Товар получен'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['buyer']['api_token']);
        $this->assertTest('Покупатель подтверждает получение', $response['status'] === 200, $response);
        
        echo "\n";
    }

    private function testAccessRights() {
        echo "🔒 4. ТЕСТИРОВАНИЕ ПРАВ ДОСТУПА\n";
        echo "----------------------------------\n";
        
        if (!$this->createdOrder) {
            echo "❌ Пропуск: заказ не создан\n\n";
            return;
        }
        
        // Покупатель пытается подтвердить заказ (должно быть запрещено)
        $statusData = ['status' => 'confirmed', 'reason' => 'Попытка подтверждения покупателем'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['buyer']['api_token']);
        $this->assertTest('Покупатель не может подтвердить заказ', $response['status'] === 403, $response);
        
        // Продавец пытается подтвердить доставку (должно быть запрещено)
        $statusData = ['status' => 'delivered', 'reason' => 'Попытка подтверждения продавцом'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['seller']['api_token']);
        $this->assertTest('Продавец не может подтвердить доставку', $response['status'] === 403, $response);
        
        // Попытка установить недопустимый статус
        $statusData = ['status' => 'invalid_status', 'reason' => 'Недопустимый статус'];
        $response = $this->makeRequest('PUT', '/orders/' . $this->createdOrder . '/status', $statusData, $this->users['seller']['api_token']);
        $this->assertTest('Недопустимый статус отклоняется', $response['status'] === 400, $response);
        
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
        echo sprintf("%-50s %s\n", $testName, $status);
        
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

// Запуск тестов
$test = new OrderStatusTest();
$test->runAllTests();
?> 