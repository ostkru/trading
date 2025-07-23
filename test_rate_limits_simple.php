<?php
/**
 * Простой тест для демонстрации текущего состояния лимитов API
 * API ключ: 026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f
 */

$apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
$baseUrl = 'http://localhost:8095';

function makeRequest($method, $endpoint, $data = null) {
    global $baseUrl, $apiKey;
    
    $url = $baseUrl . $endpoint;
    
    $headers = [
        'Authorization: Bearer ' . $apiKey,
        'Content-Type: application/json'
    ];
    
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
    }
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return [
        'status' => $httpCode,
        'response' => $response
    ];
}

echo "🚀 ТЕСТ ТЕКУЩЕГО СОСТОЯНИЯ ЛИМИТОВ API\n";
echo "==========================================\n";
echo "API ключ: $apiKey\n";
echo "Базовый URL: $baseUrl\n\n";

echo "📋 ТЕКУЩЕЕ СОСТОЯНИЕ:\n";
echo "--------------------------------\n";
echo "⚠️  ВНИМАНИЕ: Система лимитов еще не реализована!\n";
echo "Все запросы будут проходить успешно (HTTP 200/201)\n";
echo "Код 429 (Too Many Requests) не будет возвращаться\n\n";

// Тест 1: Базовая проверка
echo "1️⃣ БАЗОВАЯ ПРОВЕРКА ДОСТУПА\n";
$result = makeRequest('GET', '/api/v1/products');
echo "GET /api/v1/products: HTTP {$result['status']}\n";
if ($result['status'] === 200) {
    echo "✅ Успешно - API доступен\n";
} else {
    echo "❌ Ошибка доступа\n";
}

// Тест 2: Стресс-тест для демонстрации отсутствия лимитов
echo "\n2️⃣ СТРЕСС-ТЕСТ (50 запросов подряд)\n";
echo "Демонстрация отсутствия лимитов...\n";

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
    
    usleep(10000); // 0.01 секунды
}

echo "\n📊 ИТОГИ СТРЕСС-ТЕСТА:\n";
echo "Успешных запросов: {$successCount}\n";
echo "Запросов с превышением лимита: {$limitCount}\n";
echo "Ошибок: {$errorCount}\n";

if ($limitCount === 0) {
    echo "\n⚠️  ВЫВОД: Лимиты НЕ РАБОТАЮТ\n";
    echo "Все {$successCount} запросов прошли успешно\n";
    echo "Это означает, что система лимитов еще не реализована\n";
} else {
    echo "\n✅ ВЫВОД: Лимиты РАБОТАЮТ\n";
    echo "Обнаружено {$limitCount} запросов с превышением лимита\n";
}

// Тест 3: Проверка разных типов запросов
echo "\n3️⃣ ПРОВЕРКА РАЗНЫХ ТИПОВ ЗАПРОСОВ\n";

$tests = [
    ['GET', '/api/v1/products', 'Получение продуктов'],
    ['GET', '/api/v1/offers', 'Получение офферов'],
    ['GET', '/api/v1/orders', 'Получение заказов'],
    ['GET', '/api/v1/warehouses', 'Получение складов'],
    ['GET', '/api/v1/offers?filter=all', 'Получение всех офферов'],
    ['GET', '/api/v1/offers/public', 'Получение публичных офферов'],
];

foreach ($tests as $test) {
    $result = makeRequest($test[0], $test[1]);
    echo "{$test[2]}: HTTP {$result['status']}";
    
    if ($result['status'] === 200) {
        echo " ✅\n";
    } elseif ($result['status'] === 429) {
        echo " ⚠️ (ЛИМИТ)\n";
    } else {
        echo " ❌\n";
    }
    
    usleep(50000); // 0.05 секунды
}

echo "\n✅ ТЕСТ ЗАВЕРШЕН\n";
echo "==========================================\n";
echo "Для реализации лимитов необходимо:\n";
echo "1. Создать middleware для проверки лимитов\n";
echo "2. Интегрировать с таблицей api_rate_limits\n";
echo "3. Настроить логику минутных и дневных лимитов\n";
?> 