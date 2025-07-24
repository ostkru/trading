# PortalData API Server - Управление

## 🚀 Быстрый старт

### Запуск сервера
```bash
cd /var/www/go
./manage_server.sh start
```

### Проверка статуса
```bash
./manage_server.sh status
```

### Просмотр логов
```bash
./manage_server.sh logs
```

### Тестирование API
```bash
./manage_server.sh test
```

## 📋 Доступные команды

| Команда | Описание |
|---------|----------|
| `start` | Запустить сервер |
| `stop` | Остановить сервер |
| `restart` | Перезапустить сервер |
| `status` | Показать статус |
| `logs` | Показать логи в реальном времени |
| `test` | Протестировать API |
| `build` | Собрать сервер |

## 🔧 Альтернативные способы запуска

### 1. Простой запуск в фоне
```bash
./run_server.sh
```

### 2. Systemd сервис (автозапуск)
```bash
# Включить автозапуск
sudo systemctl enable portaldata-api.service

# Запустить сервис
sudo systemctl start portaldata-api.service

# Проверить статус
sudo systemctl status portaldata-api.service

# Остановить сервис
sudo systemctl stop portaldata-api.service
```

### 3. Ручной запуск
```bash
# Запуск в фоне
nohup ./app_8095 > server.log 2>&1 &

# Запуск с логированием
./app_8095 > server.log 2>&1 &

# Запуск в терминале (для отладки)
./app_8095
```

## 📊 Мониторинг

### Проверка процессов
```bash
ps aux | grep app_8095
```

### Проверка портов
```bash
netstat -tlnp | grep 8095
lsof -i :8095
```

### Просмотр логов
```bash
# Последние строки
tail -20 server.log

# Логи в реальном времени
tail -f server.log

# Поиск ошибок
grep -i error server.log
```

## 🛠️ Устранение проблем

### Сервер не запускается
1. Проверьте, что порт 8095 свободен:
   ```bash
   lsof -i :8095
   ```

2. Убейте процессы, занимающие порт:
   ```bash
   pkill -f app_8095
   ```

3. Пересоберите сервер:
   ```bash
   ./manage_server.sh build
   ```

### Сервер зависает
1. Остановите все процессы:
   ```bash
   pkill -f app_8095
   ```

2. Запустите с помощью скрипта управления:
   ```bash
   ./manage_server.sh start
   ```

### Проблемы с базой данных
1. Проверьте подключение к MySQL:
   ```bash
   mysql -u root -p123456 -e "SELECT 1;"
   ```

2. Проверьте переменные окружения в коде

## 📁 Структура файлов

```
/var/www/go/
├── app_8095                    # Исполняемый файл сервера
├── manage_server.sh            # Основной скрипт управления
├── run_server.sh               # Простой скрипт запуска
├── app_8095_improved          # Улучшенная версия запуска
├── server.log                  # Логи сервера
├── server.pid                  # PID файл (создается автоматически)
├── portaldata-api.service      # Systemd сервис
└── README_SERVER.md           # Эта документация
```

## 🔐 API Ключи

Для работы с API требуется авторизация:
- **API Key**: `026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f`
- **Header**: `Authorization: Bearer <token>` или `X-API-KEY: <token>`

## 🌐 Доступные эндпоинты

- **API Base URL**: `http://localhost:8095/api/v1`
- **Swagger UI**: `http://localhost:8095/swagger/index.html`
- **Products**: `/api/v1/products`
- **Offers**: `/api/v1/offers`
- **Orders**: `/api/v1/orders`
- **Warehouses**: `/api/v1/warehouses`

## 📝 Примеры использования

### Тестирование API с curl
```bash
# Получить продукты (требует авторизацию)
curl -H "Authorization: Bearer 026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f" \
     http://localhost:8095/api/v1/products

# Создать продукт
curl -X POST \
     -H "Authorization: Bearer 026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Product","article":"TEST001","brand":"TestBrand"}' \
     http://localhost:8095/api/v1/products
```

## ⚡ Автоматический перезапуск

Сервер настроен на автоматический перезапуск при сбоях:
- **Systemd**: Автоматический перезапуск через 5 секунд
- **Скрипт управления**: Ручной перезапуск через `restart`
- **Логирование**: Все ошибки записываются в `server.log` 