<?php
/**
 * Тест для проверки лимитов API
 * API ключ: 026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f
 */

$apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
$baseUrl = 'http://localhost:8095';

function makeRequest($method, $endpoint, $data = null, $headers = []) {
    global $baseUrl, $apiKey;
    
    $url = $baseUrl . $endpoint;
    
    $defaultHeaders = [
        'Authorization: Bearer ' . $apiKey,
        'Content-Type: application/json'
    ];
    
    $headers = array_merge($defaultHeaders, $headers);
    
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
    $error = curl_error($ch);
    curl_close($ch);
    
    return [
        'status' => $httpCode,
        'response' => $response,
        'error' => $error
    ];
}

function testRateLimits() {
    echo "🚀 ТЕСТИРОВАНИЕ ЛИМИТОВ API\n";
    echo "==========================================\n\n";
    
    // Тест 1: Проверка базового доступа
    echo "📋 1. БАЗОВАЯ ПРОВЕРКА ДОСТУПА\n";
    echo "--------------------------------\n";
    
    $result = makeRequest('GET', '/api/v1/products');
    echo "Получение списка продуктов: ";
    if ($result['status'] === 200) {
        echo "✅ ПРОЙДЕН (HTTP {$result['status']})\n";
    } else {
        echo "❌ ПРОВАЛЕН (HTTP {$result['status']})\n";
        echo "Ответ: " . $result['response'] . "\n";
    }
    
    // Тест 2: Проверка минутных лимитов (все методы)
    echo "\n⏱️  2. ПРОВЕРКА МИНУТНЫХ ЛИМИТОВ\n";
    echo "--------------------------------\n";
    
    $minuteTests = [
        ['GET', '/api/v1/products', null, 'Получение продуктов'],
        ['GET', '/api/v1/offers', null, 'Получение офферов'],
        ['GET', '/api/v1/orders', null, 'Получение заказов'],
        ['GET', '/api/v1/warehouses', null, 'Получение складов'],
        ['POST', '/api/v1/products', ['name' => 'Test Product', 'vendor_article' => 'TEST001', 'recommend_price' => 100, 'brand' => 'TestBrand', 'category' => 'TestCategory', 'description' => 'Test description'], 'Создание продукта'],
    ];
    
    foreach ($minuteTests as $test) {
        $result = makeRequest($test[0], $test[1], isset($test[2]) ? $test[2] : null);
        echo "{$test[3]}: ";
        if ($result['status'] === 200 || $result['status'] === 201) {
            echo "✅ ПРОЙДЕН (HTTP {$result['status']})\n";
        } elseif ($result['status'] === 429) {
            echo "⚠️  ЛИМИТ ПРЕВЫШЕН (HTTP {$result['status']})\n";
        } else {
            echo "❌ ПРОВАЛЕН (HTTP {$result['status']})\n";
        }
        
        // Небольшая пауза между запросами
        usleep(100000); // 0.1 секунды
    }
    
    // Тест 3: Проверка дневных лимитов (методы all и public)
    echo "\n📅 3. ПРОВЕРКА ДНЕВНЫХ ЛИМИТОВ\n";
    echo "--------------------------------\n";
    
    $dailyTests = [
        ['GET', '/api/v1/offers?filter=all', 'Получение всех офферов'],
        ['GET', '/api/v1/offers/public', 'Получение публичных офферов'],
    ];
    
    foreach ($dailyTests as $test) {
        $result = makeRequest($test[0], $test[1]);
        echo "{$test[2]}: ";
        if ($result['status'] === 200) {
            echo "✅ ПРОЙДЕН (HTTP {$result['status']})\n";
        } elseif ($result['status'] === 429) {
            echo "⚠️  ДНЕВНОЙ ЛИМИТ ПРЕВЫШЕН (HTTP {$result['status']})\n";
        } else {
            echo "❌ ПРОВАЛЕН (HTTP {$result['status']})\n";
        }
        
        usleep(100000);
    }
    
    // Тест 4: Стресс-тест для проверки лимитов
    echo "\n💥 4. СТРЕСС-ТЕСТ ЛИМИТОВ\n";
    echo "--------------------------------\n";
    
    echo "Отправка множественных запросов для проверки лимитов...\n";
    
    $successCount = 0;
    $limitCount = 0;
    $errorCount = 0;
    
    for ($i = 1; $i <= 50; $i++) {
        $result = makeRequest('GET', '/api/v1/products');
        
        if ($result['status'] === 200) {
            $successCount++;
        } elseif ($result['status'] === 429) {
            $limitCount++;
        } else {
            $errorCount++;
        }
        
        if ($i % 10 === 0) {
            echo "Запрос {$i}: Успешно: {$successCount}, Лимит: {$limitCount}, Ошибки: {$errorCount}\n";
        }
        
        usleep(50000); // 0.05 секунды
    }
    
    echo "\n📊 ИТОГИ СТРЕСС-ТЕСТА:\n";
    echo "Успешных запросов: {$successCount}\n";
    echo "Запросов с превышением лимита: {$limitCount}\n";
    echo "Ошибок: {$errorCount}\n";
    
    // Тест 5: Проверка разных типов запросов
    echo "\n🔄 5. ПРОВЕРКА РАЗНЫХ ТИПОВ ЗАПРОСОВ\n";
    echo "--------------------------------\n";
    
    $differentTests = [
        ['GET', '/api/v1/products', null, 'GET запрос'],
        ['POST', '/api/v1/products', ['name' => 'Stress Test', 'vendor_article' => 'STRESS001', 'recommend_price' => 50, 'brand' => 'StressBrand', 'category' => 'StressCategory', 'description' => 'Stress test product'], 'POST запрос'],
        ['GET', '/api/v1/offers?filter=my', null, 'GET с параметрами'],
        ['GET', '/api/v1/offers?filter=others', null, 'GET с параметрами (другие)'],
    ];
    
    foreach ($differentTests as $test) {
        $result = makeRequest($test[0], $test[1], isset($test[2]) ? $test[2] : null);
        echo "{$test[3]}: ";
        if ($result['status'] === 200 || $result['status'] === 201) {
            echo "✅ ПРОЙДЕН (HTTP {$result['status']})\n";
        } elseif ($result['status'] === 429) {
            echo "⚠️  ЛИМИТ ПРЕВЫШЕН (HTTP {$result['status']})\n";
        } else {
            echo "❌ ПРОВАЛЕН (HTTP {$result['status']})\n";
        }
        
        usleep(100000);
    }
    
    echo "\n✅ ТЕСТИРОВАНИЕ ЛИМИТОВ ЗАВЕРШЕНО\n";
    echo "==========================================\n";
}

// Запуск тестов
testRateLimits();
?> 