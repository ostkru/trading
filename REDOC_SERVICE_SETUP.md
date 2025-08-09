# Настройка постоянной работы Redoc сервера

## ✅ Статус: НАСТРОЕНО И РАБОТАЕТ

### 📋 Что было сделано:

1. **Создан systemd сервис** `redoc-monitor.service` для автоматического мониторинга и перезапуска
2. **Создан скрипт мониторинга** `/var/www/trading/keep_redoc_alive.sh` который:
   - Проверяет доступность сервера каждые 60 секунд
   - Автоматически перезапускает сервер при сбоях
   - Ведет логи в `/var/www/trading/redoc.log`

### 🔧 Конфигурация:

**Сервис**: `/etc/systemd/system/redoc-monitor.service`
```ini
[Unit]
Description=Redoc Server Monitor
After=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/var/www/trading
ExecStart=/var/www/trading/keep_redoc_alive.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Скрипт мониторинга**: `/var/www/trading/keep_redoc_alive.sh`
- Проверяет `http://localhost:8090/` каждые 60 секунд
- Автоматически перезапускает сервер при недоступности
- Использует правильный PATH для Node.js v22.16.0

### 🚀 Управление сервисом:

```bash
# Статус сервиса
systemctl status redoc-monitor.service

# Запуск сервиса
systemctl start redoc-monitor.service

# Остановка сервиса
systemctl stop redoc-monitor.service

# Перезапуск сервиса
systemctl restart redoc-monitor.service

# Включение автозапуска
systemctl enable redoc-monitor.service
```

### 📊 Мониторинг:

- **Логи сервиса**: `journalctl -u redoc-monitor.service`
- **Логи redoc**: `/var/www/trading/redoc.log`
- **PID файл**: `/var/www/trading/redoc.pid`

### 🌐 Доступные адреса:

- **Локальный**: `http://localhost:8090/`
- **Внешний**: `http://92.53.64.38:8090/`
- **Спецификация**: `http://92.53.64.38:8090/spec.json`

### ✅ Преимущества настройки:

1. **Автоматический перезапуск** при сбоях
2. **Автозапуск** при перезагрузке системы
3. **Мониторинг** доступности сервера
4. **Логирование** всех событий
5. **Стабильная работа** 24/7

### 🔄 Как это работает:

1. Systemd запускает скрипт мониторинга
2. Скрипт каждые 60 секунд проверяет доступность сервера
3. Если сервер не отвечает, он автоматически перезапускается
4. Все действия логируются для диагностики

### 📝 Проверка работы:

```bash
# Проверка статуса
systemctl status redoc-monitor.service

# Проверка доступности
curl -s http://localhost:8090/ | head -5

# Проверка внешнего доступа
curl -s http://92.53.64.38:8090/ | head -5

# Просмотр логов
tail -f /var/www/trading/redoc.log
```

**Сервер настроен и работает стабильно!** 🎉
