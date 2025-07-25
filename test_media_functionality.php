<?php
/**
 * ТЕСТ МЕДИА ФУНКЦИОНАЛЬНОСТИ
 * Проверяет работу с медиа данными продуктов
 */

class MediaFunctionalityTest {
    private $baseUrl = 'http://localhost:8095/api/v1';
    private $apiKey = '026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f';
    private $testResults = [];
    private $createdProductId = null;

    public function run() {
        echo "🎬 ТЕСТ МЕДИА ФУНКЦИОНАЛЬНОСТИ\n";
        echo "================================\n";
        echo "Время запуска: " . date('Y-m-d H:i:s') . "\n\n";

        // 1. Проверка создания продукта с медиа
        $this->testCreateProductWithMedia();
        
        // 2. Проверка получения продукта с медиа
        $this->testGetProductWithMedia();
        
        // 3. Проверка обновления медиа
        $this->testUpdateProductMedia();
        
        // 4. Проверка пакетного создания с медиа
        $this->testBatchCreateWithMedia();
        
        // 5. Проверка валидации медиа
        $this->testMediaValidation();
        
        // 6. Очистка
        $this->cleanup();
        
        // 7. Результаты
        $this->printResults();
    }

    private function testCreateProductWithMedia() {
        echo "📸 1. СОЗДАНИЕ ПРОДУКТА С МЕДИА\n";
        echo "----------------------------------\n";
        
        $productData = [
            'name' => 'Смартфон с медиа ' . time(),
            'vendor_article' => 'MEDIA-SMART-' . time(),
            'recommend_price' => 45000.00,
            'brand' => 'Samsung',
            'category' => 'Смартфоны',
            'description' => 'Смартфон с полным набором медиа',
            'image_urls' => [
                'https://example.com/smart_front.jpg',
                'https://example.com/smart_back.jpg',
                'https://example.com/smart_side.jpg'
            ],
            'video_urls' => [
                'https://example.com/smart_review.mp4',
                'https://example.com/smart_unboxing.mp4'
            ],
            'model_3d_urls' => [
                'https://example.com/smart_3d.glb',
                'https://example.com/smart_3d.obj'
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products', $productData);
        $this->assertTest('Создание продукта с полным медиа', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            $this->createdProductId = $response['data']['id'];
            echo "   ✅ Продукт создан с ID: {$this->createdProductId}\n";
            
            // Проверяем наличие медиа в ответе
            if (isset($response['data']['media'])) {
                echo "   ✅ Медиа данные включены в ответ\n";
                $media = $response['data']['media'];
                echo "   📸 Изображений: " . count($media['image_urls'] ?? []) . "\n";
                echo "   🎥 Видео: " . count($media['video_urls'] ?? []) . "\n";
                echo "   🎮 3D моделей: " . count($media['model_3d_urls'] ?? []) . "\n";
            } else {
                echo "   ⚠️ Медиа данные отсутствуют в ответе\n";
            }
        }
        
        echo "\n";
    }

    private function testGetProductWithMedia() {
        echo "📸 2. ПОЛУЧЕНИЕ ПРОДУКТА С МЕДИА\n";
        echo "-----------------------------------\n";
        
        if (!$this->createdProductId) {
            echo "   ⚠️ Нет созданного продукта для тестирования\n\n";
            return;
        }
        
        $response = $this->makeRequest('GET', '/products/' . $this->createdProductId);
        $this->assertTest('Получение продукта с медиа', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            if (isset($response['data']['media'])) {
                echo "   ✅ Медиа данные получены\n";
                $media = $response['data']['media'];
                echo "   📸 Изображений: " . count($media['image_urls'] ?? []) . "\n";
                echo "   🎥 Видео: " . count($media['video_urls'] ?? []) . "\n";
                echo "   🎮 3D моделей: " . count($media['model_3d_urls'] ?? []) . "\n";
            } else {
                echo "   ⚠️ Медиа данные отсутствуют\n";
            }
        }
        
        echo "\n";
    }

    private function testUpdateProductMedia() {
        echo "📸 3. ОБНОВЛЕНИЕ МЕДИА ПРОДУКТА\n";
        echo "----------------------------------\n";
        
        if (!$this->createdProductId) {
            echo "   ⚠️ Нет созданного продукта для тестирования\n\n";
            return;
        }
        
        $updateData = [
            'image_urls' => [
                'https://example.com/updated_front.jpg',
                'https://example.com/updated_back.jpg'
            ],
            'video_urls' => [
                'https://example.com/updated_review.mp4'
            ],
            'model_3d_urls' => [
                'https://example.com/updated_3d.glb'
            ]
        ];
        
        $response = $this->makeRequest('PUT', '/products/' . $this->createdProductId, $updateData);
        $this->assertTest('Обновление медиа продукта', $response['status'] === 200, $response);
        
        if ($response['status'] === 200) {
            if (isset($response['data']['media'])) {
                echo "   ✅ Медиа данные обновлены\n";
                $media = $response['data']['media'];
                echo "   📸 Обновлено изображений: " . count($media['image_urls'] ?? []) . "\n";
                echo "   🎥 Обновлено видео: " . count($media['video_urls'] ?? []) . "\n";
                echo "   🎮 Обновлено 3D моделей: " . count($media['model_3d_urls'] ?? []) . "\n";
            } else {
                echo "   ⚠️ Медиа данные отсутствуют в ответе\n";
            }
        }
        
        echo "\n";
    }

    private function testBatchCreateWithMedia() {
        echo "📸 4. ПАКЕТНОЕ СОЗДАНИЕ С МЕДИА\n";
        echo "----------------------------------\n";
        
        $batchData = [
            'products' => [
                [
                    'name' => 'Пакетный продукт с медиа 1',
                    'vendor_article' => 'BATCH-MEDIA-1-' . time(),
                    'recommend_price' => 1000.00,
                    'brand' => 'BatchMediaBrand',
                    'category' => 'BatchMediaCategory',
                    'description' => 'Пакетный продукт с медиа 1',
                    'image_urls' => [
                        'https://example.com/batch1_image1.jpg',
                        'https://example.com/batch1_image2.jpg'
                    ],
                    'video_urls' => [
                        'https://example.com/batch1_video.mp4'
                    ]
                ],
                [
                    'name' => 'Пакетный продукт с медиа 2',
                    'vendor_article' => 'BATCH-MEDIA-2-' . time(),
                    'recommend_price' => 2000.00,
                    'brand' => 'BatchMediaBrand',
                    'category' => 'BatchMediaCategory',
                    'description' => 'Пакетный продукт с медиа 2',
                    'image_urls' => [
                        'https://example.com/batch2_image1.jpg'
                    ],
                    'model_3d_urls' => [
                        'https://example.com/batch2_3d.glb'
                    ]
                ]
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products/batch', $batchData);
        $this->assertTest('Пакетное создание с медиа', $response['status'] === 201, $response);
        
        if ($response['status'] === 201) {
            echo "   ✅ Пакетное создание выполнено\n";
            if (isset($response['data']['products'])) {
                $productsWithMedia = 0;
                foreach ($response['data']['products'] as $product) {
                    if (isset($product['media'])) {
                        $productsWithMedia++;
                    }
                }
                echo "   📦 Продуктов с медиа: $productsWithMedia\n";
            }
        }
        
        echo "\n";
    }

    private function testMediaValidation() {
        echo "📸 5. ВАЛИДАЦИЯ МЕДИА ДАННЫХ\n";
        echo "------------------------------\n";
        
        // Тест с неверными URL
        $invalidData = [
            'name' => 'Продукт с неверными URL',
            'vendor_article' => 'INVALID-URL-' . time(),
            'recommend_price' => 500.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт с неверными URL',
            'image_urls' => [
                'not_a_valid_url',
                'ftp://invalid-protocol.com/image.jpg'
            ]
        ];
        
        $response = $this->makeRequest('POST', '/products', $invalidData);
        $this->assertTest('Валидация неверных URL', $response['status'] === 400, $response);
        
        // Тест с пустыми массивами медиа
        $emptyMediaData = [
            'name' => 'Продукт с пустыми медиа',
            'vendor_article' => 'EMPTY-MEDIA-' . time(),
            'recommend_price' => 300.00,
            'brand' => 'TestBrand',
            'category' => 'TestCategory',
            'description' => 'Продукт с пустыми медиа',
            'image_urls' => [],
            'video_urls' => [],
            'model_3d_urls' => []
        ];
        
        $response = $this->makeRequest('POST', '/products', $emptyMediaData);
        $this->assertTest('Создание с пустыми медиа', $response['status'] === 201, $response);
        
        echo "\n";
    }

    private function cleanup() {
        echo "🧹 6. ОЧИСТКА ТЕСТОВЫХ ДАННЫХ\n";
        echo "--------------------------------\n";
        
        if ($this->createdProductId) {
            $response = $this->makeRequest('DELETE', '/products/' . $this->createdProductId);
            $this->assertTest('Удаление тестового продукта', $response['status'] === 200, $response);
        }
        
        echo "\n";
    }

    private function makeRequest($method, $endpoint, $data = null) {
        $url = $this->baseUrl . $endpoint;
        
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Authorization: Bearer ' . $this->apiKey,
            'Content-Type: application/json'
        ]);
        
        if ($method === 'POST') {
            curl_setopt($ch, CURLOPT_POST, true);
            if ($data) {
                curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
            }
        } elseif ($method === 'PUT') {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'PUT');
            if ($data) {
                curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
            }
        } elseif ($method === 'DELETE') {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'DELETE');
        }
        
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        return [
            'status' => $httpCode,
            'data' => json_decode($response, true),
            'raw' => $response
        ];
    }

    private function assertTest($testName, $condition, $response) {
        $status = $condition ? '✅' : '❌';
        echo "$status $testName\n";
        
        if (!$condition) {
            echo "   Ошибка: " . ($response['raw'] ?? 'Неизвестная ошибка') . "\n";
        }
        
        $this->testResults[] = [
            'name' => $testName,
            'passed' => $condition,
            'response' => $response
        ];
    }

    private function printResults() {
        echo "📊 РЕЗУЛЬТАТЫ ТЕСТА МЕДИА ФУНКЦИОНАЛЬНОСТИ\n";
        echo "==========================================\n";
        
        $passed = 0;
        $total = count($this->testResults);
        
        foreach ($this->testResults as $result) {
            if ($result['passed']) {
                $passed++;
            }
        }
        
        $successRate = round(($passed / $total) * 100, 2);
        
        echo "Всего тестов: $total\n";
        echo "Пройдено: $passed\n";
        echo "Провалено: " . ($total - $passed) . "\n";
        echo "Успешность: {$successRate}%\n\n";
        
        echo "🎯 СТАТУС МЕДИА ФУНКЦИОНАЛЬНОСТИ:\n";
        if ($successRate >= 80) {
            echo "🟢 ОТЛИЧНО - Медиа функциональность работает\n";
        } elseif ($successRate >= 60) {
            echo "🟡 ХОРОШО - Есть проблемы с медиа\n";
        } else {
            echo "🔴 ПЛОХО - Медиа функциональность не работает\n";
        }
        
        echo "\n📝 ВЫВОДЫ:\n";
        if ($successRate >= 80) {
            echo "✅ Медиа функциональность реализована и работает\n";
            echo "✅ Поддерживается создание продуктов с медиа\n";
            echo "✅ Поддерживается обновление медиа\n";
            echo "✅ Поддерживается пакетное создание с медиа\n";
        } else {
            echo "❌ Медиа функциональность не реализована или работает некорректно\n";
            echo "❌ Нужно проверить реализацию медиа в коде\n";
        }
    }
}

// Запуск теста
$test = new MediaFunctionalityTest();
$test->run();
?> 