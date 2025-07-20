<?php
echo "=== ТЕСТ СКОРОСТИ PHP vs GO API ===\n\n";

$api_key = "f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53";

// Тест PHP API
echo "Тестирование PHP API (products.php?action=list):\n";
$php_url = "http://92.53.64.38/products.php?action=list&api_key={$api_key}";
$php_total = 0;
$php_requests = 5;

for ($i = 1; $i <= $php_requests; $i++) {
    $start = microtime(true);
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $php_url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 10);
    curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true); // Следуем редиректам
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "X-API-Key: {$api_key}",
        "Content-Type: application/json"
    ]);
    $response = curl_exec($ch);
    $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    $end = microtime(true);
    $time = ($end - $start) * 1000;
    $php_total += $time;
    echo "  Запрос $i: {$time}ms (HTTP: $http_code)\n";
}

$php_avg = $php_total / $php_requests;
echo "  Среднее время PHP: {$php_avg}ms\n\n";

// Тест Go API
echo "Тестирование Go API (port 8095):\n";
$go_url = "http://92.53.64.38:8095/products?api_key={$api_key}";
$go_total = 0;
$go_requests = 5;

for ($i = 1; $i <= $go_requests; $i++) {
    $start = microtime(true);
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $go_url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_TIMEOUT, 10);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "X-API-Key: {$api_key}",
        "Content-Type: application/json"
    ]);
    $response = curl_exec($ch);
    $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    $end = microtime(true);
    $time = ($end - $start) * 1000;
    $go_total += $time;
    echo "  Запрос $i: {$time}ms (HTTP: $http_code)\n";
}

$go_avg = $go_total / $go_requests;
echo "  Среднее время Go: {$go_avg}ms\n\n";

// Сравнение
echo "=== РЕЗУЛЬТАТЫ ===\n";
echo "PHP API среднее время: {$php_avg}ms\n";
echo "Go API среднее время: {$go_avg}ms\n";

if ($php_avg < $go_avg) {
    $diff = $go_avg - $php_avg;
    $percent = ($diff / $php_avg) * 100;
    echo "PHP быстрее на {$diff}ms ({$percent}%)\n";
} else {
    $diff = $php_avg - $go_avg;
    $percent = ($diff / $go_avg) * 100;
    echo "Go быстрее на {$diff}ms ({$percent}%)\n";
}
?> 