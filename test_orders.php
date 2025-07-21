<?php

function sendRequest($method, $url, $apiKey, $data = null) {
    $ch = curl_init();
    $headers = [
        'Content-Type: application/json',
        'Authorization: Bearer ' . $apiKey,
    ];

    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_VERBOSE, true); // Включаем подробный вывод

    if ($method === 'POST') {
        curl_setopt($ch, CURLOPT_POST, true);
        if ($data) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        }
    }

    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    
    if (curl_errno($ch)) {
        echo 'Curl error: ' . curl_error($ch) . "\n";
    }
    
    curl_close($ch);

    return ['code' => $httpCode, 'body' => json_decode($response, true)];
}

// --- НАСТРОЙКИ ---
$baseUrl = 'http://localhost:8095/api/v1';
// ЗАМЕНИТЕ НА ВАШИ РЕАЛЬНЫЕ КЛЮЧИ
$apiKeyUserA = 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5'; // Пользователь, который создает оффер
$apiKeyUserB = '00601582c3163466e0fece95d8e2315cb1c66814066ad8e0566d2813614d9001'; // Пользователь, который создает заказ
$apiKeyUserC = 'f456d094d3581bc14bd4f5d9bd474db9cfe8966583412b9dea6a7abc00bfa8df'; // Посторонний пользователь

// --- 1. Пользователь A создает оффер ---
echo "Шаг 1: Пользователь A создает оффер...\n";
$offerData = [
    'product_id' => 1,
    'offer_type' => 'sell',
    'price_per_unit' => 150.75,
    'available_lots' => 10,
    'tax_nds' => 20,
    'units_per_lot' => 1,
    'warehouse_id' => 1,
];
$createOfferResponse = sendRequest('POST', $baseUrl . '/offers', $apiKeyUserA, $offerData);
if ($createOfferResponse['code'] !== 201) {
    die("Ошибка: не удалось создать оффер. Код: " . $createOfferResponse['code'] . "\nТело: " . print_r($createOfferResponse['body'], true));
}
$offerId = $createOfferResponse['body']['offer_id'];
echo "Оффер создан. ID: $offerId\n\n";

// --- 2. Пользователь B создает заказ на этот оффер ---
echo "Шаг 2: Пользователь B создает заказ на оффер $offerId...\n";
$orderData = [
    'offer_id' => $offerId,
    'lot_count' => 2,
];
$createOrderResponse = sendRequest('POST', $baseUrl . '/orders', $apiKeyUserB, $orderData);
if ($createOrderResponse['code'] !== 201) {
    die("Ошибка: не удалось создать заказ. Код: " . $createOrderResponse['code'] . "\nТело: " . print_r($createOrderResponse['body'], true));
}
$orderId = $createOrderResponse['body']['order_id'];
echo "Заказ создан. ID: $orderId\n\n";

// --- 3. Проверки ---
echo "Шаг 3: Проверки...\n";

// 3.1 Пользователь A (контрагент)
$ordersUserA = sendRequest('GET', $baseUrl . '/orders?role=counterparty', $apiKeyUserA);
$foundA = false;
foreach ($ordersUserA['body']['orders'] as $order) {
    if ($order['order_id'] === $orderId) $foundA = true;
}
echo " - Пользователь A видит заказ как контрагент: " . ($foundA ? "УСПЕХ" : "ПРОВАЛ") . "\n";

// 3.2 Пользователь B (инициатор)
$ordersUserB = sendRequest('GET', $baseUrl . '/orders?role=initiator', $apiKeyUserB);
$foundB = false;
foreach ($ordersUserB['body']['orders'] as $order) {
    if ($order['order_id'] === $orderId) $foundB = true;
}
echo " - Пользователь B видит заказ как инициатор: " . ($foundB ? "УСПЕХ" : "ПРОВАЛ") . "\n";

// 3.3 Пользователь C (посторонний)
$ordersUserC = sendRequest('GET', $baseUrl . '/orders', $apiKeyUserC);
$foundC = false;
foreach ($ordersUserC['body']['orders'] as $order) {
    if ($order['order_id'] === $orderId) $foundC = true;
}
echo " - Посторонний пользователь C НЕ видит заказ: " . (!$foundC ? "УСПЕХ" : "ПРОВАЛ") . "\n";

// 3.4 Проверка доступа к заказу по ID
$getOrderA = sendRequest('GET', $baseUrl . '/orders/' . $orderId, $apiKeyUserA);
echo " - Пользователь A имеет доступ к заказу по ID: " . ($getOrderA['code'] === 200 ? "УСПЕХ" : "ПРОВАЛ") . "\n";

$getOrderB = sendRequest('GET', $baseUrl . '/orders/' . $orderId, $apiKeyUserB);
echo " - Пользователь B имеет доступ к заказу по ID: " . ($getOrderB['code'] === 200 ? "УСПЕХ" : "ПРОВАЛ") . "\n";

$getOrderC = sendRequest('GET', $baseUrl . '/orders/' . $orderId, $apiKeyUserC);
echo " - Пользователь C НЕ имеет доступа к заказу по ID: " . ($getOrderC['code'] !== 200 ? "УСПЕХ" : "ПРОВАЛ") . "\n";

?> 