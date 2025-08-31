# 🚀 Быстрый запуск Trading API

## 📋 Минимальные требования

- Linux с systemd
- Redis (установлен и доступен)
- Go приложение скомпилировано в `./api`

## ⚡ Быстрый запуск (3 команды)

```bash
# 1. Клонировать/распаковать проект
cd /path/to/trading

# 2. Настроить сервисы (автоматически)
sudo ./scripts/setup-services.sh

# 3. Запустить сервисы
sudo ./scripts/manage-services.sh start
```

## 🔍 Что делает setup-services.sh

1. **Автоматически определяет** пути к исполняемым файлам
2. **Проверяет зависимости** (Redis, systemd)
3. **Создает systemd сервисы** с правильными путями
4. **Настраивает автозапуск** при загрузке системы

## 📁 Структура после установки

```
/etc/systemd/system/
├── redis-trading.service      # Redis сервис
└── trading-api.service        # Trading API сервис

/var/www/trading/ (или другая директория)
├── api                        # Исполняемый файл
├── scripts/
│   ├── setup-services.sh      # Настройка сервисов
│   └── manage-services.sh     # Управление сервисами
└── SYSTEMD_SERVICES_README.md # Подробная документация
```

## 🎯 Основные команды

```bash
# Управление сервисами
sudo ./scripts/manage-services.sh start      # Запуск
sudo ./scripts/manage-services.sh stop       # Остановка
sudo ./scripts/manage-services.sh restart    # Перезапуск
sudo ./scripts/manage-services.sh status     # Статус

# Через Makefile (если доступен)
make services-start
make services-status
```

## 🚨 Устранение проблем

### "Исполняемый файл 'api' не найден"
```bash
# Скомпилировать приложение
go build -o api ./cmd/api

# Или указать путь вручную
export API_PATH=/path/to/your/api
sudo ./scripts/setup-services.sh
```

### "Redis не установлен"
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install redis-server

# CentOS/RHEL
sudo yum install redis

# Проверить установку
redis-cli ping
```

### "Permission denied"
```bash
# Запустить с правами root
sudo ./scripts/setup-services.sh
```

## 📊 Проверка работы

```bash
# Статус сервисов
sudo ./scripts/manage-services.sh status

# Проверка портов
netstat -tlnp | grep -E ':(6379|8095)'

# Тест API
curl http://localhost:8095/
curl http://localhost:8095/api/v1/rate-limit/stats
```

## 🔄 Переустановка

```bash
# Удалить сервисы
sudo ./scripts/manage-services.sh uninstall

# Настроить заново
sudo ./scripts/setup-services.sh
```

## 📚 Дополнительная информация

- **Подробная документация:** `SYSTEMD_SERVICES_README.md`
- **Управление сервисами:** `scripts/manage-services.sh --help`
- **Логи:** `sudo journalctl -u trading-api.service -f`
