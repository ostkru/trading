<?php
echo "=== ПОЛНОЕ ТЕСТИРОВАНИЕ API ===\n\n";

$base_url = "http://localhost:8095/api/v1";
$api_key = "f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53";
$headers = ["Authorization: Bearer $api_key", "Content-Type: application/json"];

sleep(10);

function test_endpoint($method, $url, $data = null) {
    global $headers;
    $ch = curl_init();
    $options = [
        CURLOPT_URL => $url,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_CUSTOMREQUEST => $method,
        CURLOPT_HTTPHEADER => $headers,
    ];

    if ($data) {
        $options[CURLOPT_POSTFIELDS] = json_encode($data);
    }

    curl_setopt_array($ch, $options);
    $response = curl_exec($ch);
    $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);

    echo "{$method} {$url} | HTTP: {$http_code}\n";
    // echo "Response: " . mb_strimwidth($response, 0, 150, "...") . "\n\n";
    
    $decoded_response = json_decode($response, true);
    return [$http_code, $decoded_response];
}

// --- Metaproducts ---
echo "--- 1. Metaproducts ---\n";
$metaproduct_data = [
    "name" => "Тестовый метапродукт " . time(),
    "vendor_article" => "TEST-MP-" . time(),
    "recommend_price" => 1999.99,
    "brand" => "TestBrand",
    "category" => "TestCategory",
    "description" => "Описание тестового метапродукта"
];
list($code, $response) = test_endpoint('POST', "{$base_url}/metaproducts", $metaproduct_data);
$metaproduct_id = $response['id'] ?? null;

if ($metaproduct_id) {
    test_endpoint('GET', "{$base_url}/metaproducts");
    test_endpoint('GET', "{$base_url}/metaproducts/{$metaproduct_id}");
    test_endpoint('PUT', "{$base_url}/metaproducts/{$metaproduct_id}", ["description" => "Обновленное описание"]);
    test_endpoint('DELETE', "{$base_url}/metaproducts/{$metaproduct_id}");
}
test_endpoint('POST', "{$base_url}/metaproducts/batch", [
    "products" => [$metaproduct_data]
]);


// --- Warehouses ---
echo "\n--- 2. Warehouses ---\n";
$warehouse_data = [
    "name" => "Тестовый склад " . time(),
    "address" => "Тестовый адрес, 123",
    "latitude" => 55.7558,
    "longitude" => 37.6173,
    "working_hours" => "Пн-Пт 09:00-18:00"
];
list($code, $response) = test_endpoint('POST', "{$base_url}/warehouses", $warehouse_data);
$warehouse_id = $response['id'] ?? null;

if ($warehouse_id) {
    test_endpoint('GET', "{$base_url}/warehouses");
    test_endpoint('PUT', "{$base_url}/warehouses/{$warehouse_id}", ["address" => "Новый тестовый адрес"]);
    // test_endpoint('DELETE', "{$base_url}/warehouses/{$warehouse_id}"); // Don't delete so we can use it for offers
}


// --- Offers ---
echo "\n--- 3. Offers ---\n";
if ($metaproduct_id && $warehouse_id) {
    $offer_data = [
        "product_id" => $metaproduct_id,
        "warehouse_id" => $warehouse_id,
        "supplier_id" => 1,
        "offer_type" => "sell",
        "price_per_unit" => 1500.00,
        "available_lots" => 100,
        "tax_nds" => 20,
        "units_per_lot" => 1,
        "max_shipping_days" => 7
    ];
    list($code, $response) = test_endpoint('POST', "{$base_url}/offers", $offer_data);
    $offer_id = $response['offer_id'] ?? null;

    if ($offer_id) {
        test_endpoint('GET', "{$base_url}/offers");
        test_endpoint('GET', "{$base_url}/offers/{$offer_id}");
        test_endpoint('PUT', "{$base_url}/offers/{$offer_id}", ["available_lots" => 90]);
        // test_endpoint('DELETE', "{$base_url}/offers/{$offer_id}"); // Don't delete so we can use it for orders
    }
}
test_endpoint('GET', "{$base_url}/offers/public");
test_endpoint('GET', "{$base_url}/offers/wb_stock?product_id=1&warehouse_id=1&supplier_id=1");


// --- Orders ---
echo "\n--- 4. Orders ---\n";
if (isset($offer_id)) {
    // This should fail, as user can't order their own offer
    test_endpoint('POST', "{$base_url}/orders", ["offer_id" => $offer_id, "quantity" => 1]);
}
test_endpoint('GET', "{$base_url}/orders");
test_endpoint('GET', "{$base_url}/orders/1"); // Check for a specific order (might 404)


// --- Cleanup ---
echo "\n--- 5. Cleanup ---\n";
if (isset($offer_id)) {
    test_endpoint('DELETE', "{$base_url}/offers/{$offer_id}");
}
if ($warehouse_id) {
     test_endpoint('DELETE', "{$base_url}/warehouses/{$warehouse_id}");
}

echo "\n=== ТЕСТИРОВАНИЕ ЗАВЕРШЕНО ===\n";

?> 