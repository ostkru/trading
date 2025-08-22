<?php
/**
 * Скрипт для быстрого запуска тестов Redis Rate Limiting
 */

require_once 'test_redis_rate_limiting.php';

function main() {
    $baseUrl = $argv[1] ?? 'http://localhost:8080';
    
    echo "🚀 ТЕСТИРОВАНИЕ REDIS RATE LIMITING\n";
    echo "=====================================\n\n";
    echo "🌐 Base URL: $baseUrl\n";
    echo "⏰ Время: " . date('Y-m-d H:i:s') . "\n\n";
    
    // Проверяем доступность Redis API
    echo "🔍 Проверка доступности Redis Rate Limiting API...\n";
    
    $ch = curl_init();
    curl_setopt_array($ch, [
        CURLOPT_URL => $baseUrl . '/api/v1/rate-limit/stats',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_TIMEOUT => 5,
        CURLOPT_HEADER => false,
        CURLOPT_NOBODY => false
    ]);
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    if ($httpCode === 200) {
        echo "✅ Redis Rate Limiting API доступен\n\n";
    } else {
        echo "⚠️  Redis Rate Limiting API недоступен (HTTP $httpCode)\n";
        echo "💡 Убедитесь, что сервер запущен и Redis подключен\n\n";
    }
    
    // Запускаем тесты
    try {
        $tester = new RedisRateLimitingTest($baseUrl);
        $success = $tester->runTests();
        
        echo "\n" . str_repeat("=", 60) . "\n";
        if ($success) {
            echo "🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ УСПЕШНО!";
        } else {
            echo "😞 НЕКОТОРЫЕ ТЕСТЫ НЕ ПРОШЛИ!";
        }
        echo "\n" . str_repeat("=", 60) . "\n";
        
        return $success ? 0 : 1;
        
    } catch (Exception $e) {
        echo "❌ ОШИБКА ВЫПОЛНЕНИЯ ТЕСТОВ: " . $e->getMessage() . "\n";
        return 1;
    }
}

// Запуск только если файл вызван напрямую
if (basename(__FILE__) === basename($_SERVER['SCRIPT_NAME'])) {
    exit(main());
}

?>
