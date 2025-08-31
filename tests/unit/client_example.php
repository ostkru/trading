<?php
/**
 * ПРИМЕР ПРАВИЛЬНОГО ИСПОЛЬЗОВАНИЯ API СОЗДАНИЯ ПРОДУКТОВ
 * Для клиента, который жалуется на ошибку "цена должна быть положительной"
 */

$baseUrl = 'http://localhost:8095/api/v1';
$apiToken = 'YOUR_API_TOKEN_HERE'; // Замените на ваш API токен

echo "📚 ПРИМЕРЫ ПРАВИЛЬНОГО ИСПОЛЬЗОВАНИЯ API СОЗДАНИЯ ПРОДУКТОВ\n";
echo "================================================================\n\n";

echo "🔑 ВАЖНО: Замените YOUR_API_TOKEN_HERE на ваш реальный API токен!\n\n";

// Пример 1: Правильное создание продукта
echo "✅ ПРИМЕР 1: Правильное создание продукта\n";
echo "----------------------------------------\n";
$data = [
    'name' => 'Телевизор Samsung 32" Smart TV',
    'vendor_article' => 'SAMSUNG-32D3',
    'recommend_price' => 15999.99,  // ✅ Положительная цена
    'brand' => 'Samsung',
    'category' => 'Электроника',
    'description' => '32-дюймовый смарт-телевизор с поддержкой Smart TV'
];

echo "Отправляемые данные:\n";
echo json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE) . "\n\n";

// Пример 2: Продукт с медиаданными
echo "✅ ПРИМЕР 2: Продукт с медиаданными\n";
echo "----------------------------------\n";
$data2 = [
    'name' => 'Смартфон iPhone 15 Pro',
    'vendor_article' => 'IPHONE-15-PRO-128',
    'recommend_price' => 89999.00,  // ✅ Положительная цена
    'brand' => 'Apple',
    'category' => 'Смартфоны',
    'description' => 'Смартфон с чипом A17 Pro и камерой 48 Мп',
    'image_urls' => [
        'https://example.com/iphone15pro-1.jpg',
        'https://example.com/iphone15pro-2.jpg'
    ],
    'video_urls' => [
        'https://example.com/iphone15pro-review.mp4'
    ]
];

echo "Отправляемые данные:\n";
echo json_encode($data2, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE) . "\n\n";

// Пример 3: Минимальный продукт
echo "✅ ПРИМЕР 3: Минимальный продукт (только обязательные поля)\n";
echo "---------------------------------------------------------\n";
$data3 = [
    'name' => 'Книга "Война и мир"',
    'vendor_article' => 'BOOK-WAR-PEACE',
    'recommend_price' => 599.99,  // ✅ Положительная цена
    'brand' => 'Издательство',
    'category' => 'Книги'
];

echo "Отправляемые данные:\n";
echo json_encode($data3, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE) . "\n\n";

// Примеры НЕПРАВИЛЬНОГО использования
echo "❌ ПРИМЕРЫ НЕПРАВИЛЬНОГО ИСПОЛЬЗОВАНИЯ (ВЫЗЫВАЮТ ОШИБКУ)\n";
echo "==========================================================\n\n";

echo "❌ ОШИБКА: Нулевая цена\n";
echo "recommend_price: 0  ← Это вызовет ошибку 'Цена должна быть положительной'\n\n";

echo "❌ ОШИБКА: Отрицательная цена\n";
echo "recommend_price: -100  ← Это вызовет ошибку 'Цена должна быть положительной'\n\n";

echo "❌ ОШИБКА: Отсутствует цена\n";
echo "Поле recommend_price не указано  ← Это вызовет ошибку 'Цена должна быть положительной'\n\n";

echo "❌ ОШИБКА: Цена как строка\n";
echo "recommend_price: '100'  ← Это может вызвать ошибку валидации\n\n";

// Функция для отправки запроса
function makeRequest($method, $endpoint, $data = null, $apiToken = null) {
    global $baseUrl;
    
    $url = $baseUrl . $endpoint;
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Content-Type: application/json',
        'Authorization: Bearer ' . $apiToken
    ]);
    
    if ($method === 'POST') {
        curl_setopt($ch, CURLOPT_POST, true);
        if ($data) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
        }
    }
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return [
        'status' => $httpCode,
        'body' => $response
    ];
}

echo "🔧 КОД ДЛЯ ОТПРАВКИ ЗАПРОСА (PHP):\n";
echo "====================================\n";
echo "<?php\n";
echo "\$data = [\n";
echo "    'name' => 'Название продукта',\n";
echo "    'vendor_article' => 'АРТИКУЛ',\n";
echo "    'recommend_price' => 100.50,  // ← ОБЯЗАТЕЛЬНО положительное число\n";
echo "    'brand' => 'Бренд',\n";
echo "    'category' => 'Категория'\n";
echo "];\n\n";
echo "\$response = makeRequest('POST', '/products', \$data, \$apiToken);\n";
echo "?>\n\n";

echo "🔧 КОД ДЛЯ ОТПРАВКИ ЗАПРОСА (JavaScript):\n";
echo "==========================================\n";
echo "const data = {\n";
echo "    name: 'Название продукта',\n";
echo "    vendor_article: 'АРТИКУЛ',\n";
echo "    recommend_price: 100.50,  // ← ОБЯЗАТЕЛЬНО положительное число\n";
echo "    brand: 'Бренд',\n";
echo "    category: 'Категория'\n";
echo "};\n\n";
echo "fetch('/api/v1/products', {\n";
echo "    method: 'POST',\n";
echo "    headers: {\n";
echo "        'Content-Type': 'application/json',\n";
echo "        'Authorization': 'Bearer ' + apiToken\n";
echo "    },\n";
echo "    body: JSON.stringify(data)\n";
echo "});\n\n";

echo "🔧 КОД ДЛЯ ОТПРАВКИ ЗАПРОСА (cURL):\n";
echo "====================================\n";
echo "curl -X POST \"http://localhost:8095/api/v1/products\" \\\n";
echo "  -H \"Authorization: Bearer YOUR_API_TOKEN\" \\\n";
echo "  -H \"Content-Type: application/json\" \\\n";
echo "  -d '{\n";
echo "    \"name\": \"Название продукта\",\n";
echo "    \"vendor_article\": \"АРТИКУЛ\",\n";
echo "    \"recommend_price\": 100.50,\n";
echo "    \"brand\": \"Бренд\",\n";
echo "    \"category\": \"Категория\"\n";
echo "  }'\n\n";

echo "📋 ЧЕКЛИСТ ПРОВЕРКИ ПЕРЕД ОТПРАВКОЙ:\n";
echo "====================================\n";
echo "✅ name - не пустая строка\n";
echo "✅ vendor_article - не пустая строка\n";
echo "✅ recommend_price - число больше 0 (например: 100.50, 15999.99)\n";
echo "✅ brand - не пустая строка\n";
echo "✅ category - не пустая строка\n";
echo "✅ API токен указан в заголовке Authorization\n";
echo "✅ Content-Type: application/json\n\n";

echo "🚨 ЧАСТЫЕ ОШИБКИ КЛИЕНТОВ:\n";
echo "==========================\n";
echo "1. recommend_price = 0 (должно быть > 0)\n";
echo "2. recommend_price = -100 (должно быть > 0)\n";
echo "3. recommend_price = '100' (должно быть число, не строка)\n";
echo "4. Отсутствует поле recommend_price\n";
echo "5. Неправильный API токен\n";
echo "6. Неправильный Content-Type\n\n";

echo "📞 ЕСЛИ ПРОБЛЕМА ОСТАЕТСЯ:\n";
echo "==========================\n";
echo "1. Проверьте логи вашего приложения\n";
echo "2. Убедитесь, что API доступен по адресу: $baseUrl\n";
echo "3. Проверьте правильность API токена\n";
echo "4. Убедитесь, что все обязательные поля заполнены\n";
echo "5. Проверьте, что цена является положительным числом\n";
?>
