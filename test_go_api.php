<?php
echo "=== ТЕСТ GO API С ПРАВИЛЬНОЙ ФУНКЦИОНАЛЬНОСТЬЮ ===\n\n";

$api_key = "f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53";

// Тест Go API - список продуктов
echo "1. Тест Go API - GET /products:\n";
$go_url = "http://92.53.64.38:8095/products?api_key={$api_key}";
$start = microtime(true);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $go_url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 10);
$response = curl_exec($ch);
$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);
$end = microtime(true);
$go_time = ($end - $start) * 1000;
echo "  Время: {$go_time}ms (HTTP: $http_code)\n";
echo "  Ответ: " . substr($response, 0, 200) . "...\n\n";

// Тест Go API - создание продукта
echo "2. Тест Go API - POST /products:\n";
$go_url = "http://92.53.64.38:8095/products?api_key={$api_key}";
$post_data = json_encode([
    "name" => "Test Product",
    "price" => 100.50,
    "description" => "Test description"
]);
$start = microtime(true);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $go_url);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $post_data);
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 10);
$response = curl_exec($ch);
$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);
$end = microtime(true);
$go_create_time = ($end - $start) * 1000;
echo "  Время: {$go_create_time}ms (HTTP: $http_code)\n";
echo "  Ответ: " . substr($response, 0, 200) . "...\n\n";

// Тест PHP API для сравнения
echo "3. Тест PHP API - GET /products.php:\n";
$php_url = "https://92.53.64.38/products.php?action=list&api_key={$api_key}";
$start = microtime(true);
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $php_url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_TIMEOUT, 10);
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
$response = curl_exec($ch);
$http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);
$end = microtime(true);
$php_time = ($end - $start) * 1000;
echo "  Время: {$php_time}ms (HTTP: $http_code)\n";
echo "  Ответ: " . substr($response, 0, 200) . "...\n\n";

// Сравнение результатов
echo "=== РЕЗУЛЬТАТЫ СРАВНЕНИЯ ===\n";
echo "Go API (GET): {$go_time}ms\n";
echo "Go API (POST): {$go_create_time}ms\n";
echo "PHP API (GET): {$php_time}ms\n\n";

if ($go_time > 0 && $php_time > 0) {
    if ($go_time < $php_time) {
        $diff = $php_time - $go_time;
        $percent = ($diff / $php_time) * 100;
        echo "Go API быстрее PHP API на {$diff}ms ({$percent}%)\n";
    } else {
        $diff = $go_time - $php_time;
        $percent = ($diff / $go_time) * 100;
        echo "PHP API быстрее Go API на {$diff}ms ({$percent}%)\n";
    }
} else {
    echo "Один из API не отвечает корректно\n";
}
?> 