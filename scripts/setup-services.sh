#!/bin/bash

# Скрипт для автоматической настройки systemd сервисов
# Автоматически определяет пути и настраивает сервисы

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Проверка прав root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "Этот скрипт должен быть запущен с правами root"
        exit 1
    fi
}

# Определение путей
detect_paths() {
    log "Определение путей..."
    
    # Текущая директория
    CURRENT_DIR=$(pwd)
    log "Текущая директория: $CURRENT_DIR"
    
    # Поиск исполняемого файла api
    if [[ -f "$CURRENT_DIR/api" ]]; then
        API_PATH="$CURRENT_DIR/api"
        WORKING_DIR="$CURRENT_DIR"
        log "Найден исполняемый файл: $API_PATH"
    elif [[ -f "$CURRENT_DIR/cmd/api/api" ]]; then
        API_PATH="$CURRENT_DIR/cmd/api/api"
        WORKING_DIR="$CURRENT_DIR"
        log "Найден исполняемый файл: $API_PATH"
    else
        # Поиск в системе
        API_PATH=$(which api 2>/dev/null || echo "")
        if [[ -n "$API_PATH" ]]; then
            WORKING_DIR=$(dirname "$API_PATH")
            log "Найден исполняемый файл в системе: $API_PATH"
        else
            error "Исполняемый файл 'api' не найден!"
            error "Убедитесь, что приложение скомпилировано"
            exit 1
        fi
    fi
    
    # Проверка Redis
    if ! command -v redis-server &> /dev/null; then
        error "Redis не установлен! Установите Redis:"
        error "Ubuntu/Debian: sudo apt install redis-server"
        error "CentOS/RHEL: sudo yum install redis"
        exit 1
    fi
    
    if ! command -v redis-cli &> /dev/null; then
        error "Redis CLI не установлен!"
        exit 1
    fi
    
    log "Redis найден: $(which redis-server)"
    
    # Проверка systemd
    if ! command -v systemctl &> /dev/null; then
        error "Systemd не найден! Этот скрипт работает только с systemd"
        exit 1
    fi
    
    log "Systemd найден: $(which systemctl)"
}

# Создание сервисов
create_services() {
    log "Создание systemd сервисов..."
    
    # Создание Redis сервиса
    cat > /etc/systemd/system/redis-trading.service << EOF
[Unit]
Description=Redis Server for Trading Application
Documentation=https://redis.io/documentation
After=network.target
Wants=network.target

[Service]
Type=forking
User=root
Group=root
PIDFile=/var/run/redis-trading.pid
ExecStart=/usr/bin/redis-server --daemonize yes --pidfile /var/run/redis-trading.pid --port 6379 --bind 127.0.0.1
ExecStop=/bin/kill -s QUIT \$MAINPID
Restart=always
RestartSec=3
TimeoutStopSec=10
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    
    # Создание Trading API сервиса
    cat > /etc/systemd/system/trading-api.service << EOF
[Unit]
Description=Trading API Service
Documentation=https://github.com/your-repo/trading-api
After=network.target redis-trading.service
Requires=redis-trading.service
Wants=network.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=$WORKING_DIR
ExecStart=$API_PATH
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
TimeoutStartSec=30
TimeoutStopSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=trading-api
Environment=GIN_MODE=release

# Защита от перезапуска
StartLimitInterval=60
StartLimitBurst=3

[Install]
WantedBy=multi-user.target
EOF
    
    log "Сервисы созданы:"
    log "  - /etc/systemd/system/redis-trading.service"
    log "  - /etc/systemd/system/trading-api.service"
}

# Установка сервисов
install_services() {
    log "Установка systemd сервисов..."
    
    # Перезагрузка systemd
    systemctl daemon-reload
    
    # Включение автозапуска
    systemctl enable redis-trading.service
    systemctl enable trading-api.service
    
    log "Сервисы установлены и включены для автозапуска"
}

# Проверка установки
verify_installation() {
    log "Проверка установки..."
    
    # Проверка файлов сервисов
    if [[ ! -f "/etc/systemd/system/redis-trading.service" ]]; then
        error "Redis сервис не создан!"
        exit 1
    fi
    
    if [[ ! -f "/etc/systemd/system/trading-api.service" ]]; then
        error "Trading API сервис не создан!"
        exit 1
    fi
    
    # Проверка статуса
    if systemctl is-enabled redis-trading.service >/dev/null 2>&1; then
        log "✅ Redis сервис включен для автозапуска"
    else
        error "❌ Redis сервис не включен для автозапуска"
        exit 1
    fi
    
    if systemctl is-enabled trading-api.service >/dev/null 2>&1; then
        log "✅ Trading API сервис включен для автозапуска"
    else
        error "❌ Trading API сервис не включен для автозапуска"
        exit 1
    fi
    
    log "✅ Установка завершена успешно!"
}

# Основная функция
main() {
    log "🚀 Настройка systemd сервисов для Trading API"
    log "=============================================="
    
    check_root
    detect_paths
    create_services
    install_services
    verify_installation
    
    log ""
    log "🎉 Сервисы успешно настроены!"
    log ""
    log "📋 Следующие шаги:"
    log "1. Запустить сервисы: sudo ./scripts/manage-services.sh start"
    log "2. Проверить статус: sudo ./scripts/manage-services.sh status"
    log "3. Остановить сервисы: sudo ./scripts/manage-services.sh stop"
    log ""
    log "📚 Документация: SYSTEMD_SERVICES_README.md"
}

# Запуск скрипта
main "$@"
