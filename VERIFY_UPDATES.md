# Инструкция по проверке обновлений документации

## Адрес документации
```
http://localhost:8090
```

## Что проверить в документации

### 1. Общее описание API
В начале документации должен быть раздел:
```
## Медиаданные продуктов
- **image_urls** - Массив URL изображений (.jpg, .jpeg, .png, .gif, .webp)
- **video_urls** - Массив URL видео файлов (.mp4, .avi, .mov, .wmv, .flv)
- **model_3d_urls** - Массив URL 3D моделей (.obj, .fbx, .3ds, .dae, .stl)
```

### 2. Схема Product
В разделе "Schemas" → "Product" должны быть поля:
- `image_urls` - Массив URL изображений
- `video_urls` - Массив URL видео файлов
- `model_3d_urls` - Массив URL 3D моделей

### 3. CreateProductRequest
В разделе "Schemas" → "CreateProductRequest" должны быть поля:
- `image_urls` - Массив URL изображений
- `video_urls` - Массив URL видео файлов
- `model_3d_urls` - Массив URL 3D моделей

### 4. UpdateProductRequest
В разделе "Schemas" → "UpdateProductRequest" должны быть поля:
- `image_urls` - Массив URL изображений
- `video_urls` - Массив URL видео файлов
- `model_3d_urls` - Массив URL 3D моделей

### 5. BatchCreateProductRequest
Новая схема для пакетного создания продуктов

### 6. Endpoints
В разделе "Paths" должны быть обновлены:

#### POST /products
Описание должно включать: "с поддержкой медиаданных"

#### PUT /products/{id}
Описание должно включать: "с поддержкой медиаданных"

#### GET /products/{id}
Описание должно включать: "с медиаданными"

#### GET /products
Описание должно включать: "с медиаданными"

#### POST /products/batch
Новый endpoint для пакетного создания продуктов

## Проверка через API

### Проверка схемы Product
```bash
curl -s http://localhost:8090/spec.json | jq '.components.schemas.Product.properties | {image_urls, video_urls, model_3d_urls}'
```

### Проверка нового endpoint
```bash
curl -s http://localhost:8090/spec.json | jq '.paths."/products/batch"'
```

### Проверка описания API
```bash
curl -s http://localhost:8090/spec.json | jq '.info.description' | grep -i "медиа"
```

## Ожидаемый результат

✅ Все поля медиаданных присутствуют в схемах
✅ Описания методов обновлены
✅ Новый endpoint добавлен
✅ Общее описание API включает информацию о медиаданных
✅ Документация доступна по адресу http://localhost:8090

## Если обновления не видны

1. Остановить сервер: `./stop_redoc.sh`
2. Запустить заново: `./start_redoc.sh`
3. Очистить кэш браузера
4. Перезагрузить страницу документации
