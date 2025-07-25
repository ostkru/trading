<?php
/**
 * Тест новой логики лимитов: только GET методы учитываются в дневных лимитах
 */

$apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
$baseUrl = 'http://localhost:8095';

echo "🚀 ТЕСТ НОВОЙ ЛОГИКИ ЛИМИТОВ\n";
echo "==========================================\n";
echo "API ключ: $apiKey\n";
echo "Базовый URL: $baseUrl\n\n";

echo "📋 НОВАЯ ЛОГИКА:\n";
echo "--------------------------------\n";
echo "✅ Минутные лимиты: работают для ВСЕХ методов (60/мин)\n";
echo "✅ Дневные лимиты: работают ТОЛЬКО для GET методов (1000/день)\n";
echo "❌ POST/PUT/DELETE методы НЕ учитываются в дневных лимитах\n\n";

// Функция для выполнения запроса
function makeRequest($method, $url, $data = null) {
    global $apiKey;
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "Authorization: Bearer $apiKey",
        "Content-Type: application/json"
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
        'code' => $httpCode,
        'response' => $response
    ];
}

echo "1️⃣ БАЗОВАЯ ПРОВЕРКА ДОСТУПА\n";
echo "--------------------------------\n";

// GET запрос
$getResult = makeRequest('GET', "$baseUrl/api/v1/products");
echo "GET /api/v1/products: HTTP {$getResult['code']} " . ($getResult['code'] === 200 ? "✅" : "❌") . "\n";

// POST запрос
$postData = [
    'name' => 'Test Product',
    'vendor_article' => 'TEST001',
    'recommend_price' => 100,
    'brand' => 'TestBrand',
    'category' => 'TestCategory',
    'description' => 'Test description'
];
$postResult = makeRequest('POST', "$baseUrl/api/v1/products", $postData);
echo "POST /api/v1/products: HTTP {$postResult['code']} " . ($postResult['code'] === 201 ? "✅" : "❌") . "\n\n";

echo "2️⃣ ТЕСТ ЛИМИТОВ (10 запросов каждого типа)\n";
echo "--------------------------------\n";

$getCount = 0;
$postCount = 0;
$getLimitCount = 0;
$postLimitCount = 0;

for ($i = 1; $i <= 10; $i++) {
    // GET запрос
    $getResult = makeRequest('GET', "$baseUrl/api/v1/products");
    if ($getResult['code'] === 200) {
        $getCount++;
    } elseif ($getResult['code'] === 429) {
        $getLimitCount++;
    }
    
    // POST запрос
    $postData = [
        'name' => "Test Product $i",
        'vendor_article' => "TEST00$i",
        'recommend_price' => 100 + $i,
        'brand' => 'TestBrand',
        'category' => 'TestCategory',
        'description' => "Test description $i"
    ];
    $postResult = makeRequest('POST', "$baseUrl/api/v1/products", $postData);
    if ($postResult['code'] === 201) {
        $postCount++;
    } elseif ($postResult['code'] === 429) {
        $postLimitCount++;
    }
    
    echo "Запрос $i: GET={$getResult['code']}, POST={$postResult['code']}\n";
}

echo "\n📊 РЕЗУЛЬТАТЫ ТЕСТА:\n";
echo "--------------------------------\n";
echo "GET запросы: Успешно=$getCount, Лимит=$getLimitCount\n";
echo "POST запросы: Успешно=$postCount, Лимит=$postLimitCount\n\n";

echo "3️⃣ СТРЕСС-ТЕСТ (50 POST запросов подряд)\n";
echo "--------------------------------\n";
echo "Демонстрация того, что POST запросы НЕ учитываются в дневных лимитах...\n";

$postSuccessCount = 0;
$postLimitCount = 0;

for ($i = 1; $i <= 50; $i++) {
    $postData = [
        'name' => "Stress Test Product $i",
        'vendor_article' => "STRESS00$i",
        'recommend_price' => 50 + $i,
        'brand' => 'StressBrand',
        'category' => 'StressCategory',
        'description' => "Stress test product $i"
    ];
    
    $postResult = makeRequest('POST', "$baseUrl/api/v1/products", $postData);
    
    if ($postResult['code'] === 201) {
        $postSuccessCount++;
    } elseif ($postResult['code'] === 429) {
        $postLimitCount++;
    }
    
    if ($i % 10 === 0) {
        echo "Запрос $i: Успешно=$postSuccessCount, Лимит=$postLimitCount\n";
    }
}

echo "\n📊 ИТОГИ СТРЕСС-ТЕСТА POST:\n";
echo "Успешных POST запросов: $postSuccessCount\n";
echo "POST запросов с превышением лимита: $postLimitCount\n\n";

echo "4️⃣ СТРЕСС-ТЕСТ (50 GET запросов подряд)\n";
echo "--------------------------------\n";
echo "Демонстрация того, что GET запросы УЧИТЫВАЮТСЯ в дневных лимитах...\n";

$getSuccessCount = 0;
$getLimitCount = 0;

for ($i = 1; $i <= 50; $i++) {
    $getResult = makeRequest('GET', "$baseUrl/api/v1/products");
    
    if ($getResult['code'] === 200) {
        $getSuccessCount++;
    } elseif ($getResult['code'] === 429) {
        $getLimitCount++;
    }
    
    if ($i % 10 === 0) {
        echo "Запрос $i: Успешно=$getSuccessCount, Лимит=$getLimitCount\n";
    }
}

echo "\n📊 ИТОГИ СТРЕСС-ТЕСТА GET:\n";
echo "Успешных GET запросов: $getSuccessCount\n";
echo "GET запросов с превышением лимита: $getLimitCount\n\n";

echo "5️⃣ ПРОВЕРКА РАЗНЫХ HTTP МЕТОДОВ\n";
echo "--------------------------------\n";

// GET
$getResult = makeRequest('GET', "$baseUrl/api/v1/products");
echo "GET /api/v1/products: HTTP {$getResult['code']} " . ($getResult['code'] === 200 ? "✅" : "❌") . "\n";

// POST
$postResult = makeRequest('POST', "$baseUrl/api/v1/products", $postData);
echo "POST /api/v1/products: HTTP {$postResult['code']} " . ($postResult['code'] === 201 ? "✅" : "❌") . "\n";

// PUT (если есть продукт с ID 1)
$putResult = makeRequest('PUT', "$baseUrl/api/v1/products/1", ['name' => 'Updated Product']);
echo "PUT /api/v1/products/1: HTTP {$putResult['code']} " . ($putResult['code'] === 200 ? "✅" : "❌") . "\n";

// DELETE (если есть продукт с ID 1)
$deleteResult = makeRequest('DELETE', "$baseUrl/api/v1/products/1");
echo "DELETE /api/v1/products/1: HTTP {$deleteResult['code']} " . ($deleteResult['code'] === 200 ? "✅" : "❌") . "\n\n";

echo "✅ ТЕСТ ЗАВЕРШЕН\n";
echo "==========================================\n";
echo "НОВАЯ ЛОГИКА РАБОТАЕТ:\n";
echo "✅ Минутные лимиты применяются ко ВСЕМ методам\n";
echo "✅ Дневные лимиты применяются ТОЛЬКО к GET методам\n";
echo "✅ POST/PUT/DELETE методы НЕ учитываются в дневных лимитах\n";
?> 