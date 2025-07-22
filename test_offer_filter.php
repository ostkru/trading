<?php

$baseUrl = 'http://localhost:8095';
$token = '80479fe392866b79e55c1640c107ee96c6aa25b7f8acf627a5ef226a5d8d1a27';

function makeRequest($method, $url, $data = null, $headers = []) {
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
    
    if ($data) {
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        $headers[] = 'Content-Type: application/json';
    }
    
    if (!empty($headers)) {
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    }
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return [
        'code' => $httpCode,
        'body' => json_decode($response, true)
    ];
}

echo "=== Тест фильтрации офферов ===\n\n";

// Тест 1: Получение только моих офферов (по умолчанию)
echo "1. Получение только моих офферов (filter=my):\n";
$response = makeRequest('GET', $baseUrl . '/api/v1/offers?filter=my', null, ["Authorization: Bearer $token"]);
echo "HTTP Code: " . $response['code'] . "\n";
if ($response['code'] === 200) {
    echo "Количество офферов: " . count($response['body']['offers']) . "\n";
    echo "Всего: " . $response['body']['total'] . "\n";
} else {
    echo "Ошибка: " . json_encode($response['body']) . "\n";
}
echo "\n";

// Тест 2: Получение чужих офферов
echo "2. Получение чужих офферов (filter=others):\n";
$response = makeRequest('GET', $baseUrl . '/api/v1/offers?filter=others', null, ["Authorization: Bearer $token"]);
echo "HTTP Code: " . $response['code'] . "\n";
if ($response['code'] === 200) {
    echo "Количество офферов: " . count($response['body']['offers']) . "\n";
    echo "Всего: " . $response['body']['total'] . "\n";
} else {
    echo "Ошибка: " . json_encode($response['body']) . "\n";
}
echo "\n";

// Тест 3: Получение всех офферов
echo "3. Получение всех офферов (filter=all):\n";
$response = makeRequest('GET', $baseUrl . '/api/v1/offers?filter=all', null, ["Authorization: Bearer $token"]);
echo "HTTP Code: " . $response['code'] . "\n";
if ($response['code'] === 200) {
    echo "Количество офферов: " . count($response['body']['offers']) . "\n";
    echo "Всего: " . $response['body']['total'] . "\n";
} else {
    echo "Ошибка: " . json_encode($response['body']) . "\n";
}
echo "\n";

// Тест 4: Получение без параметра filter (должен вернуть мои офферы по умолчанию)
echo "4. Получение без параметра filter (по умолчанию):\n";
$response = makeRequest('GET', $baseUrl . '/api/v1/offers', null, ["Authorization: Bearer $token"]);
echo "HTTP Code: " . $response['code'] . "\n";
if ($response['code'] === 200) {
    echo "Количество офферов: " . count($response['body']['offers']) . "\n";
    echo "Всего: " . $response['body']['total'] . "\n";
} else {
    echo "Ошибка: " . json_encode($response['body']) . "\n";
}
echo "\n";

// Тест 5: Неверный фильтр
echo "5. Неверный фильтр (filter=invalid):\n";
$response = makeRequest('GET', $baseUrl . '/api/v1/offers?filter=invalid', null, ["Authorization: Bearer $token"]);
echo "HTTP Code: " . $response['code'] . "\n";
if ($response['code'] === 200) {
    echo "Количество офферов: " . count($response['body']['offers']) . "\n";
    echo "Всего: " . $response['body']['total'] . "\n";
    echo "Примечание: Неверный фильтр должен вернуть мои офферы по умолчанию\n";
} else {
    echo "Ошибка: " . json_encode($response['body']) . "\n";
}
echo "\n";

echo "=== Тест завершен ===\n";
?> 