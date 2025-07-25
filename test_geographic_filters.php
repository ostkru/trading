<?php
/**
 * Тест географических фильтров для offers
 */

class GeographicFilterTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $users = [
        'user1' => [
            'api_token' => '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f'
        ],
        'user2' => [
            'api_token' => '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f'
        ]
    ];

    public function runTests() {
        echo "🧭 ТЕСТИРОВАНИЕ ГЕОГРАФИЧЕСКИХ ФИЛЬТРОВ OFFERS\n";
        echo "================================================\n\n";

        $this->testBasicGeographicFilter();
        $this->testPriceFilter();
        $this->testCombinedFilters();
        $this->testPublicOffersWithFilters();
        $this->testInvalidFilters();
    }

    private function testBasicGeographicFilter() {
        echo "📍 Тест базового географического фильтра\n";
        echo "----------------------------------------\n";

        $filters = [
            'filter' => 'all',
            'geographic' => [
                'min_latitude' => 55.0,
                'max_latitude' => 56.0,
                'min_longitude' => 37.0,
                'max_longitude' => 38.0
            ]
        ];

        $response = $this->makeRequest('POST', '/offers/filter', $filters, $this->users['user1']['api_token']);
        $this->assertTest('Географический фильтр (Москва)', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            echo "   ✅ Найдено офферов в области: " . count($response['data']['offers']) . "\n";
        }
        echo "\n";
    }

    private function testPriceFilter() {
        echo "💰 Тест фильтра по цене\n";
        echo "------------------------\n";

        $filters = [
            'filter' => 'all',
            'price_min' => 100,
            'price_max' => 5000
        ];

        $response = $this->makeRequest('POST', '/offers/filter', $filters, $this->users['user1']['api_token']);
        $this->assertTest('Фильтр по цене (100-5000)', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            echo "   ✅ Найдено офферов в ценовом диапазоне: " . count($response['data']['offers']) . "\n";
        }
        echo "\n";
    }

    private function testCombinedFilters() {
        echo "🔍 Тест комбинированных фильтров\n";
        echo "--------------------------------\n";

        $filters = [
            'filter' => 'all',
            'offer_type' => 'sale',
            'geographic' => [
                'min_latitude' => 0,
                'max_latitude' => 90,
                'min_longitude' => 0,
                'max_longitude' => 180
            ],
            'price_min' => 500,
            'available_lots' => 1
        ];

        $response = $this->makeRequest('POST', '/offers/filter', $filters, $this->users['user1']['api_token']);
        $this->assertTest('Комбинированный фильтр', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            echo "   ✅ Найдено офферов с комбинированными фильтрами: " . count($response['data']['offers']) . "\n";
        }
        echo "\n";
    }

    private function testPublicOffersWithFilters() {
        echo "🌐 Тест публичных офферов с фильтрами\n";
        echo "-------------------------------------\n";

        $filters = [
            'offer_type' => 'buy',
            'geographic' => [
                'min_latitude' => 55.0,
                'max_latitude' => 56.0,
                'min_longitude' => 37.0,
                'max_longitude' => 38.0
            ],
            'price_max' => 3000
        ];

        $response = $this->makeRequest('POST', '/offers/public/filter', $filters, null);
        $this->assertTest('Публичные офферы с фильтрами', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            echo "   ✅ Найдено публичных офферов: " . count($response['data']['offers']) . "\n";
        }
        echo "\n";
    }

    private function testInvalidFilters() {
        echo "❌ Тест некорректных фильтров\n";
        echo "------------------------------\n";

        // Тест с некорректным offer_type
        $filters = [
            'filter' => 'all',
            'offer_type' => 'invalid_type'
        ];

        $response = $this->makeRequest('POST', '/offers/filter', $filters, $this->users['user1']['api_token']);
        $this->assertTest('Некорректный offer_type', $response['status'] === 400, $response);

        // Тест с некорректным JSON
        $response = $this->makeRequest('POST', '/offers/filter', 'invalid json', $this->users['user1']['api_token']);
        $this->assertTest('Некорректный JSON', $response['status'] === 400, $response);

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
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        curl_setopt($ch, CURLOPT_MAXREDIRS, 3);

        if ($method === 'POST') {
            curl_setopt($ch, CURLOPT_POST, true);
            if ($data !== null) {
                if (is_string($data)) {
                    curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
                } else {
                    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
                }
            }
        }

        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        if ($response === false) {
            return ['status' => 0, 'error' => 'CURL error'];
        }

        $data = json_decode($response, true);
        if ($data === null) {
            $data = ['raw' => $response];
        }

        return ['status' => $httpCode, 'data' => $data];
    }

    private function assertTest($testName, $condition, $response) {
        if ($condition) {
            echo "   ✅ $testName\n";
        } else {
            echo "   ❌ $testName\n";
            echo "      HTTP Status: " . $response['status'] . "\n";
            if (isset($response['data']['error'])) {
                echo "      Error: " . $response['data']['error'] . "\n";
            }
        }
    }
}

// Запуск тестов
$test = new GeographicFilterTest();
$test->runTests();

echo "🎉 Тестирование географических фильтров завершено!\n";
?> 