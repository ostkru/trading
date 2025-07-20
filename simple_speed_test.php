<?php
echo "=== ПРОСТОЙ ТЕСТ СКОРОСТИ ===\n\n";

$api_key = "f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53";

// Тест PHP API через HTTPS
echo "PHP API (HTTPS):\n";
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
echo "  Время: {$php_time}ms (HTTP: $http_code)\n\n";

// Тест Go API
echo "Go API (HTTP):\n";
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
echo "  Время: {$go_time}ms (HTTP: $http_code)\n\n";

// Сравнение
echo "=== РЕЗУЛЬТАТЫ ===\n";
echo "PHP API: {$php_time}ms\n";
echo "Go API: {$go_time}ms\n";

if ($php_time < $go_time) {
    $diff = $go_time - $php_time;
    $percent = ($diff / $php_time) * 100;
    echo "PHP быстрее на {$diff}ms ({$percent}%)\n";
} else {
    $diff = $php_time - $go_time;
    $percent = ($diff / $go_time) * 100;
    echo "Go быстрее на {$diff}ms ({$percent}%)\n";
}
?> 