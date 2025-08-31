# 🚀 Автоматический запуск Redis и Trading API

Этот документ описывает, как настроить автоматический запуск Redis и Trading API через systemd.

## 📋 Требования

- Linux система с systemd
- Права root для установки сервисов
- Redis установлен (`redis-server`, `redis-cli`)
- Trading API скомпилирован (`./api`)

## 🔧 Установка

### 1. Автоматическая установка (рекомендуется)

```bash
# Перейти в директорию проекта
cd /var/www/trading

# Установить сервисы
sudo ./scripts/manage-services.sh install
```

### 2. Ручная установка

```bash
# Перезагрузить systemd
sudo systemctl daemon-reload

# Включить автозапуск
sudo systemctl enable redis-trading.service
sudo systemctl enable trading-api.service
```

## 🎮 Управление сервисами

### Запуск всех сервисов
```bash
sudo ./scripts/manage-services.sh start
```

### Остановка всех сервисов
```bash
sudo ./scripts/manage-services.sh stop
```

### Перезапуск всех сервисов
```bash
sudo ./scripts/manage-services.sh restart
```

### Проверка статуса
```bash
sudo ./scripts/manage-services.sh status
```

### Удаление сервисов
```bash
sudo ./scripts/manage-services.sh uninstall
```

## 📁 Структура файлов

```
/etc/systemd/system/
├── redis-trading.service      # Сервис Redis
└── trading-api.service        # Сервис Trading API

/var/www/trading/
├── scripts/
│   └── manage-services.sh     # Скрипт управления
└── SYSTEMD_SERVICES_README.md # Эта документация
```

## ⚙️ Конфигурация сервисов

### Redis Service (`redis-trading.service`)
- **Порт:** 6379
- **Привязка:** 127.0.0.1 (только локально)
- **PID файл:** `/var/run/redis-trading.pid`
- **Автоперезапуск:** Да
- **Тип:** Forking (демон)

### Trading API Service (`trading-api.service`)
- **Порт:** 8095
- **Рабочая директория:** `/var/www/trading`
- **Зависимость:** Требует `redis-trading.service`
- **Автоперезапуск:** Да
- **Тип:** Simple
- **Режим:** Production (`GIN_MODE=release`)

## 🔄 Порядок запуска

1. **Сеть** (`network.target`)
2. **Redis** (`redis-trading.service`)
3. **Trading API** (`trading-api.service`)

## 📊 Мониторинг

### Логи Redis
```bash
sudo journalctl -u redis-trading.service -f
```

### Логи Trading API
```bash
sudo journalctl -u trading-api.service -f
```

### Статус сервисов
```bash
sudo systemctl status redis-trading.service
sudo systemctl status trading-api.service
```

## 🚨 Устранение неполадок

### Redis не запускается
```bash
# Проверить логи
sudo journalctl -u redis-trading.service -n 50

# Проверить порт
sudo netstat -tlnp | grep :6379

# Проверить права на PID файл
ls -la /var/run/redis-trading.pid
```

### Trading API не запускается
```bash
# Проверить логи
sudo journalctl -u trading-api.service -n 50

# Проверить зависимости
sudo systemctl list-dependencies trading-api.service

# Проверить рабочую директорию
ls -la /var/www/trading/api
```

### Проблемы с портами
```bash
# Проверить занятые порты
sudo netstat -tlnp | grep -E ':(6379|8095)'

# Остановить конфликтующие процессы
sudo fuser -k 6379/tcp
sudo fuser -k 8095/tcp
```

## 🔒 Безопасность

- Redis привязан только к localhost (127.0.0.1)
- Сервисы запускаются от root (для продакшена рекомендуется создать отдельного пользователя)
- Ограничение на количество перезапусков (StartLimitBurst=3)

## 📝 Примечания

- При изменении конфигурации сервисов выполните `sudo systemctl daemon-reload`
- Для применения изменений перезапустите сервисы: `sudo ./scripts/manage-services.sh restart`
- Все логи записываются в systemd journal и доступны через `journalctl`

## 🆘 Поддержка

При возникновении проблем:
1. Проверьте логи: `sudo ./scripts/manage-services.sh status`
2. Проверьте зависимости: `sudo systemctl list-dependencies trading-api.service`
3. Проверьте конфигурацию: `sudo systemctl cat redis-trading.service`
