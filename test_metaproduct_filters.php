<?php

function sendRequest($url, $apiKey = null) {
    $ch = curl_init($url);
    $headers = [
        'Content-Type: application/json',
    ];
    if ($apiKey) {
        $headers[] = 'Authorization: Bearer ' . $apiKey;
    }
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_VERBOSE, true);

    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

    if (curl_errno($ch)) {
        echo 'Curl error: ' . curl_error($ch);
    }

    curl_close($ch);

    return ['code' => $httpCode, 'response' => json_decode($response, true)];
}

$baseURL = 'http://127.0.0.1:8095/api/v1/metaproduct';
$apiKeyUserA = '6f336c5a-3a18-4941-b85a-8320e82c1629';
$apiKeyUserB = '8b4b7c65-6d6c-4f7d-8d4c-7a2e2d8d8e5a';

// Тест 1: Получение всех продуктов без авторизации
echo "Тест 1: Получение всех продуктов без авторизации...\n";
$result = sendRequest($baseURL);
if ($result['code'] === 200) {
    echo "УСПЕХ\n";
} else {
    echo "ПРОВАЛ (Код: {$result['code']})\n";
    print_r($result['response']);
}

// Тест 2: Получение 'своих' продуктов (User A)
echo "Тест 2: Получение 'своих' продуктов (User A)...\n";
$result = sendRequest($baseURL . '?owner=my', $apiKeyUserA);
if ($result['code'] === 200) {
    echo "УСПЕХ\n";
} else {
    echo "ПРОВАЛ (Код: {$result['code']})\n";
    print_r($result['response']);
}

// Тест 3: Получение 'чужих' продуктов (User A)
echo "Тест 3: Получение 'чужих' продуктов (User A)...\n";
$result = sendRequest($baseURL . '?owner=others', $apiKeyUserA);
if ($result['code'] === 200) {
    echo "УСПЕХ\n";
} else {
    echo "ПРОВАЛ (Код: {$result['code']})\n";
    print_r($result['response']);
}

// Тест 4: Поиск по артикулу 'test-article'
echo "Тест 4: Поиск по артикулу 'test-article'...\n";
$result = sendRequest($baseURL . '?search=test-article');
if ($result['code'] === 200) {
    echo "УСПЕХ\n";
} else {
    echo "ПРОВАЛ (Код: {$result['code']})\n";
    print_r($result['response']);
}

// Тест 5: Поиск по штрихкоду '123456789'
echo "Тест 5: Поиск по штрихкоду '123456789'...\n";
$result = sendRequest($baseURL . '?barcode=123456789');
if ($result['code'] === 200) {
    echo "УСПЕХ\n";
} else {
    echo "ПРОВАЛ (Код: {$result['code']})\n";
    print_r($result['response']);
} 