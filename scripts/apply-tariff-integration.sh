#!/bin/bash

# Скрипт для интеграции существующей таблицы tariffs с rate limiting
# Файл: scripts/apply-tariff-integration.sh

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
    
    if [ ! -f "add_tariff_integration.sql" ]; then
        log_error "Файл add_tariff_integration.sql не найден!"
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

# Проверка существования таблицы tariffs
check_tariffs_table() {
    log_info "Проверка существования таблицы tariffs..."
    
    if mysql -u root -p -e "USE portaldata; DESCRIBE tariffs;" > /dev/null 2>&1; then
        log_success "Таблица tariffs найдена"
        
        # Показываем текущие тарифы
        log_info "Текущие тарифы:"
        mysql -u root -p -e "USE portaldata; SELECT id, name, price, JSON_UNQUOTE(JSON_EXTRACT(features, '$.daily_requests_limit')) as daily_limit FROM tariffs ORDER BY price;" 2>/dev/null
    else
        log_error "Таблица tariffs не найдена!"
        exit 1
    fi
}

# Создание резервной копии базы данных
create_backup() {
    log_info "Создание резервной копии базы данных..."
    
    BACKUP_FILE="backup_portaldata_tariff_integration_$(date +%Y%m%d_%H%M%S).sql"
    
    if mysqldump -u root -p portaldata > "$BACKUP_FILE" 2>/dev/null; then
        log_success "Резервная копия создана: $BACKUP_FILE"
    else
        log_warning "Не удалось создать резервную копию. Продолжаем без неё."
    fi
}

# Применение интеграции
apply_integration() {
    log_info "Применение интеграции тарифов..."
    
    if mysql -u root -p portaldata < add_tariff_integration.sql; then
        log_success "Интеграция успешно применена"
    else
        log_error "Ошибка при применении интеграции"
        exit 1
    fi
}

# Проверка результатов интеграции
verify_integration() {
    log_info "Проверка результатов интеграции..."
    
    # Проверяем добавление поля tariff_id в таблицу users
    if mysql -u root -p -e "USE portaldata; DESCRIBE users;" | grep -q "tariff_id"; then
        log_success "Поле tariff_id добавлено в таблицу users"
    else
        log_error "Поле tariff_id не найдено в таблице users"
        exit 1
    fi
    
    # Проверяем создание функций
    if mysql -u root -p -e "USE portaldata; SHOW FUNCTION STATUS WHERE Name = 'GetUserDailyLimit';" | grep -q "GetUserDailyLimit"; then
        log_success "Функция GetUserDailyLimit создана"
    else
        log_error "Функция GetUserDailyLimit не найдена"
        exit 1
    fi
    
    if mysql -u root -p -e "USE portaldata; SHOW FUNCTION STATUS WHERE Name = 'GetUserMinuteLimit';" | grep -q "GetUserMinuteLimit"; then
        log_success "Функция GetUserMinuteLimit создана"
    else
        log_error "Функция GetUserMinuteLimit не найдена"
        exit 1
    fi
    
    # Проверяем создание процедур
    if mysql -u root -p -e "USE portaldata; SHOW PROCEDURE STATUS WHERE Name = 'ChangeUserTariff';" | grep -q "ChangeUserTariff"; then
        log_success "Процедура ChangeUserTariff создана"
    else
        log_error "Процедура ChangeUserTariff не найдена"
        exit 1
    fi
    
    # Проверяем создание представлений
    if mysql -u root -p -e "USE portaldata; SHOW TABLES LIKE 'user_tariffs_with_limits';" | grep -q "user_tariffs_with_limits"; then
        log_success "Представление user_tariffs_with_limits создано"
    else
        log_error "Представление user_tariffs_with_limits не найдено"
        exit 1
    fi
    
    # Проверяем количество пользователей с тарифами
    USER_COUNT=$(mysql -u root -p -e "USE portaldata; SELECT COUNT(*) FROM users WHERE tariff_id IS NOT NULL;" -s -N 2>/dev/null)
    log_success "Найдено $USER_COUNT пользователей с назначенными тарифами"
}

# Тестирование функций
test_functions() {
    log_info "Тестирование созданных функций..."
    
    # Тестируем получение лимитов для пользователя
    log_info "Тестирование функции GetUserDailyLimit:"
    mysql -u root -p -e "USE portaldata; SELECT id, username, GetUserDailyLimit(id) as daily_limit FROM users LIMIT 3;" 2>/dev/null
    
    log_info "Тестирование функции GetUserMinuteLimit:"
    mysql -u root -p -e "USE portaldata; SELECT id, username, GetUserMinuteLimit(id) as minute_limit FROM users LIMIT 3;" 2>/dev/null
    
    log_success "Функции работают корректно"
}

# Тестирование представлений
test_views() {
    log_info "Тестирование созданных представлений..."
    
    log_info "Представление user_tariffs_with_limits:"
    mysql -u root -p -e "USE portaldata; SELECT user_id, username, tariff_name, daily_limit, minute_limit FROM user_tariffs_with_limits LIMIT 3;" 2>/dev/null
    
    log_info "Представление tariff_usage_statistics:"
    mysql -u root -p -e "USE portaldata; SELECT tariff_id, tariff_name, daily_limit, user_count FROM tariff_usage_statistics;" 2>/dev/null
    
    log_success "Представления работают корректно"
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
    log_info "Начинаем интеграцию существующей таблицы tariffs с rate limiting..."
    
    check_files
    check_mysql_connection
    check_tariffs_table
    create_backup
    apply_integration
    verify_integration
    test_functions
    test_views
    restart_service
    check_api
    
    log_success "Интеграция тарифов успешно завершена!"
    log_info "Теперь система rate limiting использует лимиты из тарифов пользователей:"
    log_info "1. Пользователи связаны с тарифами через поле tariff_id"
    log_info "2. Лимиты берутся из JSON поля features в таблице tariffs"
    log_info "3. Функции GetUserDailyLimit() и GetUserMinuteLimit() автоматически рассчитывают лимиты"
    log_info "4. Rate limiting теперь динамический и зависит от тарифа пользователя"
    
    log_info "Для тестирования используйте:"
    log_info "  mysql -u root -p -e \"USE portaldata; SELECT * FROM user_tariffs_with_limits;\""
    log_info "  mysql -u root -p -e \"USE portaldata; CALL ChangeUserTariff(1, 2);\""
}

# Обработка сигналов
trap 'log_error "Интеграция прервана пользователем"; exit 1' INT TERM

# Запуск основной функции
main "$@"
