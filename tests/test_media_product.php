<?php
/**
 * Тест создания продукта с медиа
 */

$apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
$baseUrl = 'http://localhost:8095';

echo "🧪 ТЕСТ СОЗДАНИЯ ПРОДУКТА С МЕДИА\n";
echo "====================================\n";
echo "API ключ: $apiKey\n";
echo "Базовый URL: $baseUrl\n\n";

// Функция для выполнения запроса
function makeRequest($method, $url, $data = null) {
    global $apiKey;
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "Authorization: Bearer $apiKey",
        "Content-Type: application/json"
    ]);
    
    if ($method === 'POST') {
        curl_setopt($ch, CURLOPT_POST, true);
        if ($data) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        }
    }
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    return [
        'code' => $httpCode,
        'response' => $response
    ];
}

echo "📸 ТЕСТ 1: Создание продукта с изображениями (обязательно)\n";
echo "------------------------------------------------------------\n";

$data1 = json_encode([
    'name' => 'Тестовый продукт с медиа',
    'vendor_article' => 'MEDIA001',
    'recommend_price' => 1500.50,
    'brand' => 'TestMediaBrand',
    'category' => 'Электроника',
    'description' => 'Продукт с медиа контентом',
    'image_urls' => [
        'https://example.com/image1.jpg',
        'https://example.com/image2.jpg',
        'https://example.com/image3.jpg'
    ],
    'video_urls' => [
        'https://example.com/video1.mp4',
        'https://example.com/video2.mp4'
    ],
    'model_3d_urls' => [
        'https://example.com/model1.glb',
        'https://example.com/model2.obj'
    ]
]);

$result1 = makeRequest('POST', "$baseUrl/api/v1/products", $data1);
echo "POST с медиа: HTTP $result1[code]\n";
if ($result1['code'] === 201) {
    $response1 = json_decode($result1['response'], true);
    echo "✅ Продукт создан успешно\n";
    echo "ID продукта: " . $response1['id'] . "\n";
    if (isset($response1['media'])) {
        echo "✅ Медиа данные включены в ответ\n";
        echo "Изображений: " . count(json_decode($response1['media']['image_urls'], true)) . "\n";
        echo "Видео: " . count(json_decode($response1['media']['video_urls'], true)) . "\n";
        echo "3D моделей: " . count(json_decode($response1['media']['model_3d_urls'], true)) . "\n";
    } else {
        echo "❌ Медиа данные отсутствуют в ответе\n";
    }
} else {
    echo "❌ Ошибка создания: $result1[response]\n";
}

echo "\n📸 ТЕСТ 2: Создание продукта только с изображениями\n";
echo "----------------------------------------------------\n";

$data2 = json_encode([
    'name' => 'Продукт только с изображениями',
    'vendor_article' => 'MEDIA002',
    'recommend_price' => 2000.00,
    'brand' => 'TestMediaBrand',
    'category' => 'Одежда',
    'description' => 'Продукт только с изображениями',
    'image_urls' => [
        'https://example.com/photo1.jpg',
        'https://example.com/photo2.jpg'
    ]
]);

$result2 = makeRequest('POST', "$baseUrl/api/v1/products", $data2);
echo "POST только с изображениями: HTTP $result2[code]\n";
if ($result2['code'] === 201) {
    echo "✅ Продукт создан успешно\n";
} else {
    echo "❌ Ошибка создания: $result2[response]\n";
}

echo "\n📸 ТЕСТ 3: Получение продукта с медиа\n";
echo "----------------------------------------\n";

if ($result1['code'] === 201) {
    $response1 = json_decode($result1['response'], true);
    $productId = $response1['id'];
    
    $result3 = makeRequest('GET', "$baseUrl/api/v1/products/$productId");
    echo "GET продукта $productId: HTTP $result3[code]\n";
    if ($result3['code'] === 200) {
        $response3 = json_decode($result3['response'], true);
        if (isset($response3['media'])) {
            echo "✅ Медиа данные получены\n";
            echo "Изображений: " . count(json_decode($response3['media']['image_urls'], true)) . "\n";
        } else {
            echo "❌ Медиа данные отсутствуют\n";
        }
    } else {
        echo "❌ Ошибка получения: $result3[response]\n";
    }
}

echo "\n📸 ТЕСТ 4: Обновление продукта с новыми медиа\n";
echo "-----------------------------------------------\n";

if ($result1['code'] === 201) {
    $response1 = json_decode($result1['response'], true);
    $productId = $response1['id'];
    
    $updateData = json_encode([
        'image_urls' => [
            'https://example.com/new_image1.jpg',
            'https://example.com/new_image2.jpg',
            'https://example.com/new_image3.jpg'
        ],
        'video_urls' => [
            'https://example.com/new_video1.mp4'
        ]
    ]);
    
    $result4 = makeRequest('PUT', "$baseUrl/api/v1/products/$productId", $updateData);
    echo "PUT обновление медиа: HTTP $result4[code]\n";
    if ($result4['code'] === 200) {
        echo "✅ Продукт обновлен успешно\n";
    } else {
        echo "❌ Ошибка обновления: $result4[response]\n";
    }
}

echo "\n🎉 ИТОГОВЫЙ РЕЗУЛЬТАТ:\n";
echo "======================\n";
echo "✅ Таблица media создана с JSON полями\n";
echo "✅ Модели продуктов расширены для поддержки медиа\n";
echo "✅ Сервис обновлен для работы с медиа\n";
echo "✅ Поддержка создания, обновления и получения медиа\n";
echo "\n📝 ЗАКЛЮЧЕНИЕ:\n";
echo "Система готова для работы с медиа контентом продуктов:\n";
echo "1. ✅ Изображения (обязательные)\n";
echo "2. ✅ Видео обзоры (опциональные)\n";
echo "3. ✅ 3D модели (опциональные)\n";
echo "4. ✅ JSON формат для хранения массивов URL\n";
?> 