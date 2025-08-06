<?php
/**
 * Тест исправлений API
 */

$baseUrl = 'http://localhost:8095';
$apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';

echo "🔧 ТЕСТ ИСПРАВЛЕНИЙ API\n";
echo "========================\n\n";

// 1. Проверка основного endpoint
echo "1. Проверка основного endpoint:\n";
$response = curl_get_contents("$baseUrl/");
if ($response) {
    $data = json_decode($response, true);
    if (isset($data['message']) && $data['message'] === 'API ПорталДанных.РФ доступен') {
        echo "✅ Основной endpoint работает\n";
    } else {
        echo "❌ Основной endpoint не работает\n";
    }
} else {
    echo "❌ Сервер не отвечает\n";
}

// 2. Проверка публичных офферов
echo "\n2. Проверка публичных офферов:\n";
$response = curl_get_contents("$baseUrl/api/v1/offers/public");
if ($response) {
    $data = json_decode($response, true);
    if (isset($data['offers'])) {
        echo "✅ Публичные офферы работают\n";
    } else {
        echo "❌ Публичные офферы не работают\n";
    }
} else {
    echo "❌ Публичные офферы недоступны\n";
}

// 3. Проверка создания заказа
echo "\n3. Проверка создания заказа:\n";
$orderData = [
    'offer_id' => 1,
    'quantity' => 1,
    'delivery_address' => 'Тестовый адрес',
    'contact_phone' => '+7-999-123-45-67',
    'contact_email' => 'test@example.com'
];

$response = curl_post_contents("$baseUrl/api/v1/orders", $orderData, $apiKey);
if ($response) {
    $data = json_decode($response, true);
    if (isset($data['id'])) {
        echo "✅ Создание заказа работает\n";
    } else {
        echo "❌ Создание заказа не работает: " . ($data['error'] ?? 'неизвестная ошибка') . "\n";
    }
} else {
    echo "❌ Создание заказа недоступно\n";
}

// 4. Проверка валидации медиа
echo "\n4. Проверка валидации медиа:\n";
$invalidMediaData = [
    'name' => 'Тест валидации',
    'vendor_article' => 'TEST-VALID',
    'recommend_price' => 100,
    'brand' => 'TestBrand',
    'category' => 'TestCategory',
    'description' => 'Тест валидации медиа',
    'image_urls' => ['not_a_valid_url', 'ftp://invalid.com/image.jpg']
];

$response = curl_post_contents("$baseUrl/api/v1/products", $invalidMediaData, $apiKey);
if ($response) {
    $data = json_decode($response, true);
    if (isset($data['error'])) {
        echo "✅ Валидация медиа работает\n";
    } else {
        echo "❌ Валидация медиа не работает\n";
    }
} else {
    echo "❌ Валидация медиа недоступна\n";
}

echo "\n🎯 ИТОГОВЫЙ РЕЗУЛЬТАТ:\n";
echo "======================\n";

function curl_get_contents($url) {
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 5);
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return $httpCode === 200 ? $response : false;
}

function curl_post_contents($url, $data, $apiKey) {
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $apiKey,
        'Content-Type: application/json'
    ]);
    curl_setopt($ch, CURLOPT_TIMEOUT, 5);
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return $httpCode >= 200 && $httpCode < 300 ? $response : false;
}
?> 