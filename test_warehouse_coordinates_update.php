<?php
/**
 * ТЕСТ ОБНОВЛЕНИЯ КООРДИНАТ ПРИ СМЕНЕ СКЛАДА
 * Проверяет, что координаты оффера автоматически обновляются при изменении warehouse_id
 */

class WarehouseCoordinatesUpdateTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $users = [
        'user1' => [
            'name' => 'clear13808',
            'api_token' => '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27'
        ]
    ];
    
    private $testResults = [];
    private $createdProducts = [];
    private $createdOffers = [];
    private $createdWarehouses = [];
    private $performanceMetrics = [];

    public function runAllTests() {
        echo "🏭 ТЕСТ ОБНОВЛЕНИЯ КООРДИНАТ ПРИ СМЕНЕ СКЛАДА\n";
        echo "================================================\n\n";

        // 1. Создание тестовых ресурсов
        $this->createTestResources();
        
        // 2. Тест изменения координат при смене склада
        $this->testWarehouseChangeCoordinates();
        
        // 3. Тест с разными координатами складов
        $this->testDifferentWarehouseCoordinates();
        
        // 4. Вывод результатов
        $this->printResults();
    }

    private function createTestResources() {
        echo "📦 СОЗДАНИЕ ТЕСТОВЫХ РЕСУРСОВ\n";
        echo "--------------------------------\n";
        
        // Создание продукта
        $productData = [
            'name' => 'Тестовый продукт для координат',
            'vendor_article' => 'COORD-TEST-' . time(),
            'recommend_price' => 100.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт для тестирования координат'
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData, $this->users['user1']['api_token']);
        $this->assertTest('Создание продукта', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdProducts['main'] = $response['data']['id'];
        }
        
        // Создание первого склада (Москва)
        $warehouse1Data = [
            'name' => 'Склад Москва',
            'address' => 'Москва, Красная площадь, 1',
            'latitude' => 55.7558,
            'longitude' => 37.6176,
            'working_hours' => '09:00-18:00'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouse1Data, $this->users['user1']['api_token']);
        $this->assertTest('Создание склада 1 (Москва)', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['moscow'] = $response['data']['id'];
        }
        
        // Создание второго склада (Санкт-Петербург)
        $warehouse2Data = [
            'name' => 'Склад Санкт-Петербург',
            'address' => 'Санкт-Петербург, Невский проспект, 1',
            'latitude' => 59.9311,
            'longitude' => 30.3609,
            'working_hours' => '10:00-19:00'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouse2Data, $this->users['user1']['api_token']);
        $this->assertTest('Создание склада 2 (СПб)', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['spb'] = $response['data']['id'];
        }
        
        // Создание третьего склада (Екатеринбург)
        $warehouse3Data = [
            'name' => 'Склад Екатеринбург',
            'address' => 'Екатеринбург, ул. Ленина, 1',
            'latitude' => 56.8519,
            'longitude' => 60.6122,
            'working_hours' => '08:00-17:00'
        ];
        
        $response = $this->makeRequest('POST', '/warehouses', $warehouse3Data, $this->users['user1']['api_token']);
        $this->assertTest('Создание склада 3 (Екатеринбург)', $response['status'] === 201, $response);
        if ($response['status'] === 201) {
            $this->createdWarehouses['ekb'] = $response['data']['id'];
        }
        
        // Создание оффера с первым складом
        if (isset($this->createdProducts['main']) && isset($this->createdWarehouses['moscow'])) {
            $offerData = [
                'product_id' => $this->createdProducts['main'],
                'offer_type' => 'sale',
                'price_per_unit' => 150.00,
                'units_per_lot' => 1,
                'available_lots' => 10,
                'tax_nds' => 20,
                'warehouse_id' => $this->createdWarehouses['moscow'],
                'is_public' => true,
                'max_shipping_days' => 7
            ];
            
            $response = $this->makeRequest('POST', '/offers', $offerData, $this->users['user1']['api_token']);
            $this->assertTest('Создание оффера с московским складом', $response['status'] === 201, $response);
            if ($response['status'] === 201) {
                $this->createdOffers['main'] = $response['data']['offer_id'];
            }
        }
        
        echo "\n";
    }

    private function testWarehouseChangeCoordinates() {
        echo "📍 ТЕСТ ИЗМЕНЕНИЯ КООРДИНАТ ПРИ СМЕНЕ СКЛАДА\n";
        echo "------------------------------------------------\n";
        
        if (!isset($this->createdOffers['main']) || !isset($this->createdWarehouses['spb'])) {
            echo "❌ Не удалось создать необходимые ресурсы для теста\n";
            return;
        }
        
        // Получаем исходные координаты оффера
        $response = $this->makeRequest('GET', '/offers/' . $this->createdOffers['main'], null, $this->users['user1']['api_token']);
        $this->assertTest('Получение исходного оффера', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            $originalLatitude = $response['data']['latitude'];
            $originalLongitude = $response['data']['longitude'];
            
            echo "   📍 Исходные координаты: $originalLatitude, $originalLongitude (Москва)\n";
            
            // Меняем склад на СПб
            $updateData = [
                'warehouse_id' => $this->createdWarehouses['spb']
            ];
            
            $startTime = microtime(true);
            $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['main'], $updateData, $this->users['user1']['api_token']);
            $endTime = microtime(true);
            $this->performanceMetrics['Смена склада на СПб'] = round(($endTime - $startTime) * 1000, 2);
            $this->assertTest('Смена склада на СПб', $response['status'] === 200, $response);
            
            if ($response['status'] === 200) {
                $newLatitude = $response['data']['latitude'];
                $newLongitude = $response['data']['longitude'];
                
                echo "   📍 Новые координаты: $newLatitude, $newLongitude (СПб)\n";
                
                // Проверяем, что координаты изменились
                $coordinatesChanged = ($newLatitude != $originalLatitude) || ($newLongitude != $originalLongitude);
                $this->assertTest('Координаты изменились при смене склада', $coordinatesChanged, $response);
                
                // Проверяем, что новые координаты соответствуют СПб
                $spbCoordinates = ($newLatitude == 59.9311) && ($newLongitude == 30.3609);
                $this->assertTest('Координаты соответствуют СПб', $spbCoordinates, $response);
            }
        }
        
        echo "\n";
    }

    private function testDifferentWarehouseCoordinates() {
        echo "🌍 ТЕСТ С РАЗНЫМИ КООРДИНАТАМИ СКЛАДОВ\n";
        echo "----------------------------------------\n";
        
        if (!isset($this->createdOffers['main']) || !isset($this->createdWarehouses['ekb'])) {
            echo "❌ Не удалось создать необходимые ресурсы для теста\n";
            return;
        }
        
        // Меняем склад на Екатеринбург
        $updateData = [
            'warehouse_id' => $this->createdWarehouses['ekb']
        ];
        
        $startTime = microtime(true);
        $response = $this->makeRequest('PUT', '/offers/' . $this->createdOffers['main'], $updateData, $this->users['user1']['api_token']);
        $endTime = microtime(true);
        $this->performanceMetrics['Смена склада на Екатеринбург'] = round(($endTime - $startTime) * 1000, 2);
        $this->assertTest('Смена склада на Екатеринбург', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            $newLatitude = $response['data']['latitude'];
            $newLongitude = $response['data']['longitude'];
            
            echo "   📍 Координаты после смены на Екатеринбург: $newLatitude, $newLongitude\n";
            
            // Проверяем, что координаты соответствуют Екатеринбургу
            $ekbCoordinates = ($newLatitude == 56.8519) && ($newLongitude == 60.6122);
            $this->assertTest('Координаты соответствуют Екатеринбургу', $ekbCoordinates, $response);
            
            // Проверяем, что координаты не соответствуют СПб
            $notSpbCoordinates = ($newLatitude != 59.9311) || ($newLongitude != 30.3609);
            $this->assertTest('Координаты не соответствуют СПб', $notSpbCoordinates, $response);
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
        } elseif ($method === 'GET') {
            // GET request
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
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТА ОБНОВЛЕНИЯ КООРДИНАТ\n";
        echo str_repeat("=", 80) . "\n\n";
        
        $totalTests = count($this->testResults);
        $passedTests = count(array_filter($this->testResults, function($test) {
            return $test['passed'];
        }));
        $failedTests = $totalTests - $passedTests;
        $successRate = round(($passedTests / $totalTests) * 100, 2);
        
        echo "📈 СТАТИСТИКА:\n";
        echo "   Всего тестов: $totalTests\n";
        echo "   Пройдено: $passedTests\n";
        echo "   Провалено: $failedTests\n";
        echo "   Успешность: $successRate%\n\n";
        
        echo "⚡ МЕТРИКИ ПРОИЗВОДИТЕЛЬНОСТИ:\n";
        echo str_repeat("-", 80) . "\n";
        foreach ($this->performanceMetrics as $testName => $time) {
            echo sprintf("%-50s %6.2f мс\n", $testName, $time);
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
        echo "🎯 ВЫВОДЫ:\n";
        if ($successRate >= 90) {
            echo "✅ Автоматическое обновление координат работает корректно!\n";
        } elseif ($successRate >= 80) {
            echo "⚠️  Есть небольшие проблемы с обновлением координат.\n";
        } else {
            echo "❌ Требуется доработка логики обновления координат.\n";
        }
        echo str_repeat("=", 80) . "\n";
    }
}

// Запуск тестов
$test = new WarehouseCoordinatesUpdateTest();
$test->runAllTests();
?> 