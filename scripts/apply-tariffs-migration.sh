#!/bin/bash

# Скрипт для применения миграции тарифов
# Файл: scripts/apply-tariffs-migration.sh

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода сообщений
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка наличия необходимых файлов
check_files() {
    log_info "Проверка наличия необходимых файлов..."
    
    if [ ! -f "add_tariffs_table.sql" ]; then
        log_error "Файл add_tariffs_table.sql не найден!"
        exit 1
    fi
    
    log_success "Все необходимые файлы найдены"
}

# Проверка подключения к MySQL
check_mysql_connection() {
    log_info "Проверка подключения к MySQL..."
    
    if ! mysql -u root -p -e "SELECT 1;" > /dev/null 2>&1; then
        log_error "Не удается подключиться к MySQL. Проверьте настройки подключения."
        exit 1
    fi
    
    log_success "Подключение к MySQL успешно"
}

# Создание резервной копии базы данных
create_backup() {
    log_info "Создание резервной копии базы данных..."
    
    BACKUP_FILE="backup_portaldata_$(date +%Y%m%d_%H%M%S).sql"
    
    if mysqldump -u root -p portaldata > "$BACKUP_FILE" 2>/dev/null; then
        log_success "Резервная копия создана: $BACKUP_FILE"
    else
        log_warning "Не удалось создать резервную копию. Продолжаем без неё."
    fi
}

# Применение миграции
apply_migration() {
    log_info "Применение миграции тарифов..."
    
    if mysql -u root -p portaldata < add_tariffs_table.sql; then
        log_success "Миграция успешно применена"
    else
        log_error "Ошибка при применении миграции"
        exit 1
    fi
}

# Проверка результатов миграции
verify_migration() {
    log_info "Проверка результатов миграции..."
    
    # Проверяем создание таблицы tariffs
    if mysql -u root -p -e "USE portaldata; DESCRIBE tariffs;" > /dev/null 2>&1; then
        log_success "Таблица tariffs создана"
    else
        log_error "Таблица tariffs не найдена"
        exit 1
    fi
    
    # Проверяем добавление поля tariff_id в таблицу users
    if mysql -u root -p -e "USE portaldata; DESCRIBE users;" | grep -q "tariff_id"; then
        log_success "Поле tariff_id добавлено в таблицу users"
    else
        log_error "Поле tariff_id не найдено в таблице users"
        exit 1
    fi
    
    # Проверяем создание таблицы api_rate_limits
    if mysql -u root -p -e "USE portaldata; DESCRIBE api_rate_limits;" > /dev/null 2>&1; then
        log_success "Таблица api_rate_limits создана"
    else
        log_error "Таблица api_rate_limits не найдена"
        exit 1
    fi
    
    # Проверяем количество тарифов
    TARIFF_COUNT=$(mysql -u root -p -e "USE portaldata; SELECT COUNT(*) FROM tariffs;" -s -N 2>/dev/null)
    if [ "$TARIFF_COUNT" -ge 5 ]; then
        log_success "Создано $TARIFF_COUNT тарифов"
    else
        log_warning "Создано только $TARIFF_COUNT тарифов (ожидалось минимум 5)"
    fi
    
    # Проверяем количество пользователей
    USER_COUNT=$(mysql -u root -p -e "USE portaldata; SELECT COUNT(*) FROM users;" -s -N 2>/dev/null)
    log_success "Найдено $USER_COUNT пользователей"
}

# Перезапуск сервиса
restart_service() {
    log_info "Перезапуск сервиса portaldata-api..."
    
    if systemctl is-active --quiet portaldata-api; then
        if systemctl restart portaldata-api; then
            log_success "Сервис portaldata-api перезапущен"
        else
            log_error "Ошибка при перезапуске сервиса"
            exit 1
        fi
    else
        log_warning "Сервис portaldata-api не запущен"
    fi
}

# Проверка работоспособности API
check_api() {
    log_info "Проверка работоспособности API..."
    
    # Ждем запуска сервиса
    sleep 5
    
    if curl -s http://localhost:8080/ > /dev/null; then
        log_success "API доступен"
    else
        log_warning "API недоступен. Проверьте логи сервиса."
    fi
}

# Основная функция
main() {
    log_info "Начинаем применение миграции тарифов..."
    
    check_files
    check_mysql_connection
    create_backup
    apply_migration
    verify_migration
    restart_service
    check_api
    
    log_success "Миграция тарифов успешно завершена!"
    log_info "Теперь вы можете:"
    log_info "1. Управлять тарифами через API endpoints"
    log_info "2. Назначать тарифы пользователям"
    log_info "3. Мониторить использование через статистику"
    log_info "4. Проверять лимиты через rate limiting"
}

# Обработка сигналов
trap 'log_error "Миграция прервана пользователем"; exit 1' INT TERM

# Запуск основной функции
main "$@"
