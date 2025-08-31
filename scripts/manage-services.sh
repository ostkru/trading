#!/bin/bash

# Скрипт для управления сервисами Trading API
# Использование: ./manage-services.sh [start|stop|restart|status|install|uninstall]

SERVICE_NAME="trading-api"
REDIS_SERVICE_NAME="redis-trading"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
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

# Удаление сервисов
uninstall_services() {
    log "Удаление systemd сервисов..."
    
    # Остановка сервисов
    systemctl stop trading-api.service 2>/dev/null
    systemctl stop redis-trading.service 2>/dev/null
    
    # Отключение автозапуска
    systemctl disable trading-api.service 2>/dev/null
    systemctl disable redis-trading.service 2>/dev/null
    
    # Удаление файлов сервисов
    rm -f /etc/systemd/system/trading-api.service
    rm -f /etc/systemd/system/redis-trading.service
    
    # Перезагрузка systemd
    systemctl daemon-reload
    
    log "Сервисы удалены"
}

# Запуск сервисов
start_services() {
    log "Запуск сервисов..."
    
    # Запуск Redis
    if systemctl start redis-trading.service; then
        log "Redis запущен успешно"
    else
        error "Ошибка запуска Redis"
        exit 1
    fi
    
    # Ожидание готовности Redis
    log "Ожидание готовности Redis..."
    sleep 3
    
    # Проверка Redis
    if redis-cli ping >/dev/null 2>&1; then
        log "Redis готов к работе"
    else
        error "Redis не отвечает"
        exit 1
    fi
    
    # Запуск API
    if systemctl start trading-api.service; then
        log "Trading API запущен успешно"
    else
        error "Ошибка запуска Trading API"
        exit 1
    fi
    
    # Ожидание готовности API
    log "Ожидание готовности API..."
    sleep 5
    
    # Проверка API
    if curl -s http://localhost:8095/ >/dev/null 2>&1; then
        log "Trading API готов к работе"
    else
        warning "Trading API может быть еще не готов"
    fi
}

# Остановка сервисов
stop_services() {
    log "Остановка сервисов..."
    
    # Остановка API
    if systemctl stop trading-api.service; then
        log "Trading API остановлен"
    fi
    
    # Остановка Redis
    if systemctl stop redis-trading.service; then
        log "Redis остановлен"
    fi
}

# Перезапуск сервисов
restart_services() {
    log "Перезапуск сервисов..."
    stop_services
    sleep 2
    start_services
}

# Статус сервисов
status_services() {
    log "Статус сервисов:"
    echo "----------------------------------------"
    
    # Статус Redis
    echo "Redis Service:"
    systemctl status redis-trading.service --no-pager -l
    
    echo "----------------------------------------"
    
    # Статус API
    echo "Trading API Service:"
    systemctl status trading-api.service --no-pager -l
    
    echo "----------------------------------------"
    
    # Проверка портов
    echo "Проверка портов:"
    if netstat -tlnp | grep :6379 >/dev/null; then
        echo "✅ Redis (6379): активен"
    else
        echo "❌ Redis (6379): неактивен"
    fi
    
    if netstat -tlnp | grep :8095 >/dev/null; then
        echo "✅ Trading API (8095): активен"
    else
        echo "❌ Trading API (8095): неактивен"
    fi
    
    echo "----------------------------------------"
    
    # Проверка Redis
    echo "Проверка Redis:"
    if redis-cli ping >/dev/null 2>&1; then
        echo "✅ Redis отвечает: $(redis-cli ping)"
    else
        echo "❌ Redis не отвечает"
    fi
    
    # Проверка API
    echo "Проверка API:"
    if curl -s http://localhost:8095/ >/dev/null 2>&1; then
        echo "✅ API отвечает"
    else
        echo "❌ API не отвечает"
    fi
}

# Основная логика
main() {
    check_root
    
    case "$1" in
        start)
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        status)
            status_services
            ;;
        install)
            install_services
            ;;
        uninstall)
            uninstall_services
            ;;
        *)
            echo "Использование: $0 {start|stop|restart|status|install|uninstall}"
            echo ""
            echo "Команды:"
            echo "  start     - Запустить все сервисы"
            echo "  stop      - Остановить все сервисы"
            echo "  restart   - Перезапустить все сервисы"
            echo "  status    - Показать статус сервисов"
            echo "  install   - Установить systemd сервисы"
            echo "  uninstall - Удалить systemd сервисы"
            exit 1
            ;;
    esac
}

# Запуск скрипта
main "$@"
