# Адрес документации API

## Сохраненный адрес документации на порту 8090

### Локальный доступ
```
http://localhost:8090
```

### Внешний доступ
```
http://92.53.64.38:8090
```

## Информация о сервере документации

- **Порт**: 8090
- **Хост**: 0.0.0.0 (внешний доступ)
- **Технология**: Redoc
- **Файл спецификации**: openapi.json
- **Логи**: redoc.log

## Запуск документации

### Автоматический запуск
```bash
./start_redoc.sh
```

### Остановка
```bash
./stop_redoc.sh
```

## Обновления документации

Документация была обновлена с учетом новых параметров медиаданных:

### Добавленные возможности
- **image_urls** - Массив URL изображений
- **video_urls** - Массив URL видео файлов  
- **model_3d_urls** - Массив URL 3D моделей
- **Пакетное создание продуктов** - Новый endpoint `/products/batch`

### Поддерживаемые форматы
- **Изображения**: .jpg, .jpeg, .png, .gif, .webp
- **Видео**: .mp4, .avi, .mov, .wmv, .flv
- **3D модели**: .obj, .fbx, .3ds, .dae, .stl

### Валидация
- Проверка URL и расширений файлов
- Поддержка только HTTP/HTTPS протоколов
- Транзакционная безопасность операций

## Структура медиаданных

Медиаданные сохраняются в отдельной таблице `media` и связаны с продуктами через `product_id`:

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

## Примеры использования

### Создание продукта с медиаданными
```json
{
    "name": "Смартфон Galaxy Pro",
    "vendor_article": "GALAXY001",
    "recommend_price": 45000.00,
    "brand": "Samsung",
    "category": "Смартфоны",
    "description": "Мощный смартфон с отличной камерой",
    "image_urls": [
        "https://example.com/galaxy_front.jpg",
        "https://example.com/galaxy_back.jpg"
    ],
    "video_urls": [
        "https://example.com/galaxy_review.mp4"
    ],
    "model_3d_urls": [
        "https://example.com/galaxy_3d_model.glb"
    ]
}
```

### Пакетное создание продуктов
```json
{
    "products": [
        {
            "name": "Продукт 1",
            "vendor_article": "PROD001",
            "recommend_price": 2500.00,
            "brand": "Brand1",
            "category": "Категория1",
            "image_urls": ["https://example.com/prod1.jpg"]
        },
        {
            "name": "Продукт 2", 
            "vendor_article": "PROD002",
            "recommend_price": 3500.00,
            "brand": "Brand2",
            "category": "Категория2",
            "video_urls": ["https://example.com/prod2.mp4"]
        }
    ]
}
```

## Тестирование

Для тестирования медиаданных используйте файл:
```bash
php test_media_integration.php
```

## Дополнительная документация

- [MEDIA_INTEGRATION.md](./MEDIA_INTEGRATION.md) - Подробная документация по интеграции медиаданных
- [MEDIA_INTEGRATION_SUMMARY.md](./MEDIA_INTEGRATION_SUMMARY.md) - Краткое резюме интеграции
