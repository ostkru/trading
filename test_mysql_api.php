<?php
/**
 * Тестирование Go API с MySQL и тестовыми пользователями
 */

// Тестовые API ключи пользователей
$test_users = [
    'user1' => 'f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53',
    'user2' => '6f336c5a-3a18-4941-b85a-8320e82c1629',
    'user3' => '8b4b7c65-6d6c-4f7d-8d4c-7a2e2d8d8e5a',
    'user4' => 'f9c912b6989eb166ee48ec6d8f07a2b0d29d5efc8ae1c2e44fac9fe8c4d4a0b5',
    'user5' => '00601582c3163466e0fece95d8e2315cb1c66814066ad8e0566d2813614d9001',
    'user6' => 'f456d094d3581bc14bd4f5d9bd474db9cfe8966583412b9dea6a7abc00bfa8df'
];

$base_url = 'http://localhost:8095';
$results = [];

echo "=== Тестирование Go API с MySQL ===\n";
echo "Время: " . date('Y-m-d H:i:s') . "\n\n";

// Тест 1: Проверка доступности API
echo "1. Проверка доступности API...\n";
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $base_url . '/');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 10);
$response = curl_exec($ch);
$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($http_code == 200) {
    echo "✅ API доступен (HTTP: $http_code)\n";
    $results['api_available'] = true;
} else {
    echo "❌ API недоступен (HTTP: $http_code)\n";
    $results['api_available'] = false;
    exit(1);
}

// Тест 2: Тестирование с каждым пользователем
foreach ($test_users as $username => $api_key) {
    echo "\n2. Тестирование пользователя: $username\n";
    echo "API Key: " . substr($api_key, 0, 20) . "...\n";
    
    // Тест аутентификации
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $base_url . '/api/v1/products');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 10);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $api_key,
        'Content-Type: application/json'
    ]);
    $response = curl_exec($ch);
    $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    if ($http_code == 200) {
        echo "✅ Аутентификация успешна (HTTP: $http_code)\n";
        $results[$username]['auth'] = true;
        
        // Тест создания продукта
        $product_data = [
            'name' => 'Тестовый продукт ' . $username,
            'vendor_article' => 'TEST-' . strtoupper($username) . '-001',
            'recommend_price' => 100.50,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Тестовое описание продукта для ' . $username
        ];
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $base_url . '/api/v1/products');
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($product_data));
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Authorization: Bearer ' . $api_key,
            'Content-Type: application/json'
        ]);
        $response = curl_exec($ch);
        $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if ($http_code == 201 || $http_code == 200) {
            echo "✅ Создание продукта успешно (HTTP: $http_code)\n";
            $results[$username]['create_product'] = true;
            
            $response_data = json_decode($response, true);
            if (isset($response_data['id'])) {
                $product_id = $response_data['id'];
                echo "   Создан продукт с ID: $product_id\n";
                
                // Тест получения продукта
                $ch = curl_init();
                curl_setopt($ch, CURLOPT_URL, $base_url . '/api/v1/products/' . $product_id);
                curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
                curl_setopt($ch, CURLOPT_TIMEOUT, 10);
                curl_setopt($ch, CURLOPT_HTTPHEADER, [
                    'Authorization: Bearer ' . $api_key,
                    'Content-Type: application/json'
                ]);
                $response = curl_exec($ch);
                $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
                curl_close($ch);
                
                if ($http_code == 200) {
                    echo "✅ Получение продукта успешно (HTTP: $http_code)\n";
                    $results[$username]['get_product'] = true;
                } else {
                    echo "❌ Ошибка получения продукта (HTTP: $http_code)\n";
                    $results[$username]['get_product'] = false;
                }
            }
        } else {
            echo "❌ Ошибка создания продукта (HTTP: $http_code)\n";
            echo "   Ответ: $response\n";
            $results[$username]['create_product'] = false;
        }
        
    } else {
        echo "❌ Ошибка аутентификации (HTTP: $http_code)\n";
        echo "   Ответ: $response\n";
        $results[$username]['auth'] = false;
    }
}

// Тест 3: Тестирование публичных маршрутов
echo "\n3. Тестирование публичных маршрутов...\n";
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $base_url . '/api/v1/offers/public');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 10);
$response = curl_exec($ch);
$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($http_code == 200) {
    echo "✅ Публичные маршруты работают (HTTP: $http_code)\n";
    $results['public_routes'] = true;
} else {
    echo "❌ Ошибка публичных маршрутов (HTTP: $http_code)\n";
    $results['public_routes'] = false;
}

// Итоговый отчет
echo "\n=== ИТОГОВЫЙ ОТЧЕТ ===\n";
echo "API доступен: " . ($results['api_available'] ? '✅' : '❌') . "\n";
echo "Публичные маршруты: " . ($results['public_routes'] ? '✅' : '❌') . "\n\n";

foreach ($test_users as $username => $api_key) {
    if (isset($results[$username])) {
        echo "Пользователь $username:\n";
        echo "  Аутентификация: " . ($results[$username]['auth'] ? '✅' : '❌') . "\n";
        if (isset($results[$username]['create_product'])) {
            echo "  Создание продукта: " . ($results[$username]['create_product'] ? '✅' : '❌') . "\n";
        }
        if (isset($results[$username]['get_product'])) {
            echo "  Получение продукта: " . ($results[$username]['get_product'] ? '✅' : '❌') . "\n";
        }
        echo "\n";
    }
}

echo "Тестирование завершено!\n";
?> 