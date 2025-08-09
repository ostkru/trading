# Интеграция медиаданных с продуктами

## Обзор

Система поддерживает хранение и управление медиаданными для продуктов через таблицу `media`. Медиаданные включают изображения, видео и 3D модели, хранящиеся в формате JSON массивов URL.

## Структура базы данных

### Таблица `media`
```sql
CREATE TABLE media (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    image_urls JSON,
    video_urls JSON,
    model_3d_urls JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_product_media (product_id)
);
```

## API Endpoints

### Создание продукта с медиаданными
```http
POST /api/v1/products
Content-Type: application/json
Authorization: Bearer {api_key}

{
    "name": "Смартфон Galaxy Pro",
    "vendor_article": "GALAXY001",
    "recommend_price": 45000.00,
    "brand": "Samsung",
    "category": "Смартфоны",
    "description": "Мощный смартфон с отличной камерой",
    "image_urls": [
        "https://example.com/galaxy_front.jpg",
        "https://example.com/galaxy_back.jpg",
        "https://example.com/galaxy_side.jpg"
    ],
    "video_urls": [
        "https://example.com/galaxy_review.mp4",
        "https://example.com/galaxy_unboxing.mp4"
    ],
    "model_3d_urls": [
        "https://example.com/galaxy_3d_model.glb",
        "https://example.com/galaxy_3d_model.obj"
    ]
}
```

### Получение продукта с медиаданными
```http
GET /api/v1/products/{id}
Authorization: Bearer {api_key}
```

### Обновление медиаданных продукта
```http
PUT /api/v1/products/{id}
Content-Type: application/json
Authorization: Bearer {api_key}

{
    "image_urls": [
        "https://example.com/new_galaxy_front.jpg",
        "https://example.com/new_galaxy_back.jpg"
    ],
    "video_urls": [
        "https://example.com/new_galaxy_review.mp4"
    ],
    "model_3d_urls": [
        "https://example.com/new_galaxy_3d_model.glb"
    ]
}
```

### Пакетное создание продуктов с медиаданными
```http
POST /api/v1/products/batch
Content-Type: application/json
Authorization: Bearer {api_key}

{
    "products": [
        {
            "name": "Продукт 1",
            "vendor_article": "PROD001",
            "recommend_price": 2500.00,
            "brand": "Brand1",
            "category": "Категория1",
            "description": "Описание продукта 1",
            "image_urls": [
                "https://example.com/prod1_1.jpg",
                "https://example.com/prod1_2.jpg"
            ],
            "video_urls": [
                "https://example.com/prod1_video.mp4"
            ]
        },
        {
            "name": "Продукт 2",
            "vendor_article": "PROD002",
            "recommend_price": 3500.00,
            "brand": "Brand2",
            "category": "Категория2",
            "description": "Описание продукта 2",
            "image_urls": [
                "https://example.com/prod2_1.jpg"
            ],
            "model_3d_urls": [
                "https://example.com/prod2_model.glb"
            ]
        }
    ]
}
```

## Поддерживаемые форматы файлов

### Изображения
- `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`

### Видео
- `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`

### 3D модели
- `.obj`, `.fbx`, `.3ds`, `.dae`, `.stl`

## Валидация

### URL валидация
- Поддерживаются только HTTP и HTTPS протоколы
- Обязательно наличие хоста в URL
- Проверяется расширение файла

### Примеры валидных URL
```
✅ https://example.com/image.jpg
✅ https://cdn.example.com/video.mp4
✅ https://models.example.com/model.glb
❌ ftp://example.com/file.jpg
❌ https://example.com/file.txt
❌ https://example.com/
```

## Транзакционная безопасность

Все операции с медиаданными выполняются в транзакциях:
- Создание продукта + медиаданные
- Обновление продукта + медиаданные
- Пакетное создание продуктов + медиаданные

При ошибке в любой части операции все изменения откатываются.

## Примеры использования

### Создание продукта только с изображениями
```json
{
    "name": "Простой продукт",
    "vendor_article": "SIMPLE001",
    "recommend_price": 1500.00,
    "brand": "SimpleBrand",
    "category": "Электроника",
    "description": "Продукт только с изображениями",
    "image_urls": [
        "https://example.com/simple1.jpg",
        "https://example.com/simple2.jpg"
    ]
}
```

### Обновление только медиаданных
```json
{
    "image_urls": [
        "https://example.com/new_image1.jpg",
        "https://example.com/new_image2.jpg"
    ]
}
```

### Очистка медиаданных (пустые массивы)
```json
{
    "image_urls": [],
    "video_urls": [],
    "model_3d_urls": []
}
```

## Обработка ошибок

### Ошибки валидации
```json
{
    "error": "некорректный URL изображения https://example.com/file.txt: неподдерживаемое расширение файла. Разрешены: [.jpg .jpeg .png .gif .webp]"
}
```

### Ошибки доступа
```json
{
    "error": "Продукт принадлежит другому пользователю"
}
```

### Ошибки базы данных
```json
{
    "error": "Ошибка при создании медиаданных"
}
```

## Тестирование

### Запуск тестов
```bash
php test_media_integration.php
```

### Проверка работы API
```bash
# Создание продукта с медиаданными
curl -X POST "http://localhost:8095/api/v1/products" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "name": "Тестовый продукт",
    "vendor_article": "TEST001",
    "recommend_price": 1000.00,
    "brand": "TestBrand",
    "category": "Тест",
    "description": "Тестовый продукт с медиа",
    "image_urls": ["https://example.com/test.jpg"]
  }'

# Получение продукта с медиаданными
curl "http://localhost:8095/api/v1/products/1" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

## Производительность

### Индексы
```sql
CREATE INDEX idx_media_product_id ON media(product_id);
CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_brand_category ON products(brand_id, category_id);
```

### Оптимизации
- LEFT JOIN для получения медиаданных
- COALESCE для обработки NULL значений
- JSON парсинг только при необходимости
- Транзакции для атомарности операций

## Безопасность

- Валидация URL и расширений файлов
- Проверка прав доступа к продуктам
- Санитизация JSON данных
- Защита от SQL инъекций через prepared statements
