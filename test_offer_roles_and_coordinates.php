<?php

class OfferRolesAndCoordinatesTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $user1Token = '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27';
    private $user2Token = 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5';
    
    private $testResults = [];
    private $totalTests = 0;
    private $passedTests = 0;
    
    public function run() {
        echo "🚀 ТЕСТИРОВАНИЕ ЛОГИКИ РОЛЕЙ И КООРДИНАТ\n";
        echo "==========================================\n\n";
        
        $this->testSaleOfferOrder();
        $this->testBuyOfferOrder();
        $this->testCoordinatesAutoFill();
        
        $this->printResults();
    }
    
    private function testSaleOfferOrder() {
        echo "📦 1. ТЕСТ ЗАКАЗА НА ОФФЕР ПРОДАЖИ\n";
        echo "------------------------------------\n";
        
        // Создаем оффер на продажу от User2
        $this->createSaleOffer();
        
        // User1 создает заказ на оффер User2
        $this->createOrderOnSaleOffer();
        
        // Проверяем роли в заказе
        $this->checkSaleOrderRoles();
    }
    
    private function testBuyOfferOrder() {
        echo "\n📦 2. ТЕСТ ЗАКАЗА НА ОФФЕР ПОКУПКИ\n";
        echo "-------------------------------------\n";
        
        // Создаем оффер на покупку от User2
        $this->createBuyOffer();
        
        // User1 создает заказ на оффер покупки User2
        $this->createOrderOnBuyOffer();
        
        // Проверяем роли в заказе
        $this->checkBuyOrderRoles();
    }
    
    private function testCoordinatesAutoFill() {
        echo "\n📍 3. ТЕСТ АВТОМАТИЧЕСКОГО ЗАПОЛНЕНИЯ КООРДИНАТ\n";
        echo "------------------------------------------------\n";
        
        $this->testCreateOfferWithCoordinates();
        $this->testUpdateOfferWithNewWarehouse();
    }
    
    private function createSaleOffer() {
        $data = [
            'product_id' => 19,
            'offer_type' => 'sale',
            'price_per_unit' => 300.0,
            'available_lots' => 10,
            'tax_nds' => 20,
            'units_per_lot' => 1,
            'warehouse_id' => 3,
            'is_public' => true,
            'max_shipping_days' => 5
        ];
        
        $response = $this->makeRequest('POST', '/offers', $data, $this->user2Token);
        $this->assertTest('Создание оффера на продажу', 
            isset($response['offer_id']), $response);
    }
    
    private function createBuyOffer() {
        $data = [
            'product_id' => 19,
            'offer_type' => 'buy',
            'price_per_unit' => 250.0,
            'available_lots' => 5,
            'tax_nds' => 20,
            'units_per_lot' => 1,
            'warehouse_id' => 3,
            'is_public' => true,
            'max_shipping_days' => 3
        ];
        
        $response = $this->makeRequest('POST', '/offers', $data, $this->user2Token);
        $this->assertTest('Создание оффера на покупку', 
            isset($response['offer_id']), $response);
    }
    
    private function createOrderOnSaleOffer() {
        // Получаем последний оффер на продажу
        $offers = $this->makeRequest('GET', '/offers', [], $this->user2Token);
        $saleOffer = null;
        foreach ($offers['offers'] as $offer) {
            if ($offer['offer_type'] === 'sale') {
                $saleOffer = $offer;
                break;
            }
        }
        
        if (!$saleOffer) {
            $this->assertTest('Найден оффер на продажу', false, 'Оффер на продажу не найден');
            return;
        }
        
        $data = [
            'offer_id' => $saleOffer['offer_id'],
            'quantity' => 2
        ];
        
        $response = $this->makeRequest('POST', '/orders', $data, $this->user1Token);
        $this->assertTest('Создание заказа на оффер продажи', 
            isset($response['order_id']), $response);
        
        // Сохраняем заказ для проверки ролей
        $this->saleOrder = $response;
    }
    
    private function createOrderOnBuyOffer() {
        // Получаем последний оффер на покупку
        $offers = $this->makeRequest('GET', '/offers', [], $this->user2Token);
        $buyOffer = null;
        foreach ($offers['offers'] as $offer) {
            if ($offer['offer_type'] === 'buy') {
                $buyOffer = $offer;
                break;
            }
        }
        
        if (!$buyOffer) {
            $this->assertTest('Найден оффер на покупку', false, 'Оффер на покупку не найден');
            return;
        }
        
        $data = [
            'offer_id' => $buyOffer['offer_id'],
            'quantity' => 1
        ];
        
        $response = $this->makeRequest('POST', '/orders', $data, $this->user1Token);
        $this->assertTest('Создание заказа на оффер покупки', 
            isset($response['order_id']), $response);
        
        // Сохраняем заказ для проверки ролей
        $this->buyOrder = $response;
    }
    
    private function checkSaleOrderRoles() {
        if (!isset($this->saleOrder)) {
            $this->assertTest('Проверка ролей в заказе продажи', false, 'Заказ не создан');
            return;
        }
        
        // В заказе на продажу: User1 (создатель) = покупатель, User2 (владелец оффера) = продавец
        $this->assertTest('User1 является покупателем в заказе продажи', 
            $this->saleOrder['initiator_user_id'] == 1, $this->saleOrder);
        
        $this->assertTest('User2 является продавцом в заказе продажи', 
            $this->saleOrder['counterparty_user_id'] == 2, $this->saleOrder);
        
        $this->assertTest('Тип заказа продажи - buy', 
            $this->saleOrder['order_type'] === 'buy', $this->saleOrder);
    }
    
    private function checkBuyOrderRoles() {
        if (!isset($this->buyOrder)) {
            $this->assertTest('Проверка ролей в заказе покупки', false, 'Заказ не создан');
            return;
        }
        
        // В заказе на покупку: User1 (создатель) = продавец, User2 (владелец оффера) = покупатель
        $this->assertTest('User1 является продавцом в заказе покупки', 
            $this->buyOrder['initiator_user_id'] == 1, $this->buyOrder);
        
        $this->assertTest('User2 является покупателем в заказе покупки', 
            $this->buyOrder['counterparty_user_id'] == 2, $this->buyOrder);
        
        $this->assertTest('Тип заказа покупки - sell', 
            $this->buyOrder['order_type'] === 'sell', $this->buyOrder);
    }
    
    private function testCreateOfferWithCoordinates() {
        $data = [
            'product_id' => 23,
            'offer_type' => 'sale',
            'price_per_unit' => 400.0,
            'available_lots' => 3,
            'tax_nds' => 20,
            'units_per_lot' => 1,
            'warehouse_id' => 5, // Склад с координатами
            'is_public' => true,
            'max_shipping_days' => 7
        ];
        
        $response = $this->makeRequest('POST', '/offers', $data, $this->user1Token);
        $this->assertTest('Создание оффера с автоматическими координатами', 
            isset($response['offer_id']), $response);
        
        if (isset($response['offer_id'])) {
            $this->assertTest('Координаты автоматически заполнены', 
                isset($response['latitude']) && isset($response['longitude']), $response);
            
            $this->assertTest('Координаты соответствуют складу', 
                $response['latitude'] == 59.93110000 && $response['longitude'] == 30.36090000, $response);
        }
    }
    
    private function testUpdateOfferWithNewWarehouse() {
        // Получаем последний оффер
        $offers = $this->makeRequest('GET', '/offers', [], $this->user1Token);
        if (empty($offers['offers'])) {
            $this->assertTest('Обновление координат оффера', false, 'Нет офферов для обновления');
            return;
        }
        
        $offer = $offers['offers'][0];
        
        $data = [
            'warehouse_id' => 1 // Склад с другими координатами
        ];
        
        $response = $this->makeRequest('PUT', "/offers/{$offer['offer_id']}", $data, $this->user1Token);
        $this->assertTest('Обновление оффера с новым складом', 
            isset($response['offer_id']), $response);
        
        if (isset($response['offer_id'])) {
            $this->assertTest('Координаты обновились', 
                $response['latitude'] == 55.75580000 && $response['longitude'] == 37.61760000, $response);
        }
    }
    
    private function makeRequest($method, $endpoint, $data = [], $token = null) {
        $url = $this->baseUrl . $endpoint;
        
        $headers = ['Content-Type: application/json'];
        if ($token) {
            $headers[] = "Authorization: Bearer $token";
        }
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        
        if (!empty($data)) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        }
        
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if ($httpCode >= 400) {
            return ['error' => $response];
        }
        
        return json_decode($response, true) ?: [];
    }
    
    private function assertTest($name, $condition, $response = null) {
        $this->totalTests++;
        
        if ($condition) {
            echo "✅ $name\n";
            $this->passedTests++;
        } else {
            echo "❌ $name\n";
            if ($response && isset($response['error'])) {
                echo "   Ошибка: {$response['error']}\n";
            }
        }
    }
    
    private function printResults() {
        echo "\n📊 РЕЗУЛЬТАТЫ ТЕСТИРОВАНИЯ\n";
        echo "============================\n";
        echo "Всего тестов: {$this->totalTests}\n";
        echo "Пройдено: {$this->passedTests}\n";
        echo "Провалено: " . ($this->totalTests - $this->passedTests) . "\n";
        echo "Процент успеха: " . round(($this->passedTests / $this->totalTests) * 100, 2) . "%\n\n";
        
        if ($this->passedTests == $this->totalTests) {
            echo "🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!\n";
        } else {
            echo "⚠️  НЕКОТОРЫЕ ТЕСТЫ ПРОВАЛЕНЫ\n";
        }
    }
}

// Запуск теста
$test = new OfferRolesAndCoordinatesTest();
$test->run(); 