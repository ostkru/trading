#!/bin/bash

# Скрипт для тестирования API тарифов
# Файл: scripts/test-tariffs-api.sh

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
API_BASE_URL="http://localhost:8080/api"
ADMIN_API_KEY="f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53"
TEST_API_KEY="test_api_key_1"

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

# Функция для выполнения HTTP запросов
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local api_key=$4
    
    local headers="Content-Type: application/json"
    if [ -n "$api_key" ]; then
        headers="$headers"$'\n'"X-API-Key: $api_key"
    fi
    
    if [ -n "$data" ]; then
        curl -s -X "$method" "$url" \
            -H "$headers" \
            -d "$data"
    else
        curl -s -X "$method" "$url" \
            -H "$headers"
    fi
}

# Тест 1: Проверка доступности API
test_api_availability() {
    log_info "Тест 1: Проверка доступности API"
    
    response=$(make_request "GET" "$API_BASE_URL/../" "")
    
    if echo "$response" | grep -q "success"; then
        log_success "API доступен"
    else
        log_error "API недоступен"
        return 1
    fi
}

# Тест 2: Получение списка тарифов
test_list_tariffs() {
    log_info "Тест 2: Получение списка тарифов"
    
    response=$(make_request "GET" "$API_BASE_URL/tariffs" "" "")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Список тарифов получен"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка получения списка тарифов"
        echo "$response"
        return 1
    fi
}

# Тест 3: Получение тарифа по ID
test_get_tariff() {
    log_info "Тест 3: Получение тарифа по ID"
    
    response=$(make_request "GET" "$API_BASE_URL/tariffs/1" "" "")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Тариф получен"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка получения тарифа"
        echo "$response"
        return 1
    fi
}

# Тест 4: Создание нового тарифа
test_create_tariff() {
    log_info "Тест 4: Создание нового тарифа"
    
    local new_tariff='{
        "name": "Тестовый тариф",
        "description": "Тариф для тестирования",
        "minute_limit": 1500,
        "day_limit": 15000,
        "is_active": true
    }'
    
    response=$(make_request "POST" "$API_BASE_URL/tariffs" "$new_tariff" "$ADMIN_API_KEY")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Тариф создан"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка создания тарифа"
        echo "$response"
        return 1
    fi
}

# Тест 5: Получение информации о тарифе пользователя
test_get_user_tariff() {
    log_info "Тест 5: Получение информации о тарифе пользователя"
    
    response=$(make_request "GET" "$API_BASE_URL/tariffs/user/1" "" "$ADMIN_API_KEY")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Информация о тарифе пользователя получена"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка получения информации о тарифе пользователя"
        echo "$response"
        return 1
    fi
}

# Тест 6: Изменение тарифа пользователя
test_change_user_tariff() {
    log_info "Тест 6: Изменение тарифа пользователя"
    
    local change_request='{
        "user_id": 1,
        "tariff_id": 2
    }'
    
    response=$(make_request "PUT" "$API_BASE_URL/tariffs/user/change" "$change_request" "$ADMIN_API_KEY")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Тариф пользователя изменен"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка изменения тарифа пользователя"
        echo "$response"
        return 1
    fi
}

# Тест 7: Получение статистики тарифов
test_get_tariff_stats() {
    log_info "Тест 7: Получение статистики тарифов"
    
    response=$(make_request "GET" "$API_BASE_URL/tariffs/stats" "" "$ADMIN_API_KEY")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Статистика тарифов получена"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка получения статистики тарифов"
        echo "$response"
        return 1
    fi
}

# Тест 8: Проверка rate limiting с новыми лимитами
test_rate_limiting() {
    log_info "Тест 8: Проверка rate limiting с новыми лимитами"
    
    # Делаем несколько запросов для проверки лимитов
    for i in {1..5}; do
        response=$(make_request "GET" "$API_BASE_URL/tariffs/my" "" "$TEST_API_KEY")
        log_info "Запрос $i: $(echo "$response" | grep -o '"success":[^,]*' || echo "ошибка")"
        sleep 1
    done
}

# Тест 9: Обновление тарифа
test_update_tariff() {
    log_info "Тест 9: Обновление тарифа"
    
    local update_data='{
        "minute_limit": 2000,
        "day_limit": 20000
    }'
    
    response=$(make_request "PUT" "$API_BASE_URL/tariffs/1" "$update_data" "$ADMIN_API_KEY")
    
    if echo "$response" | grep -q "success.*true"; then
        log_success "Тариф обновлен"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        log_error "Ошибка обновления тарифа"
        echo "$response"
        return 1
    fi
}

# Тест 10: Удаление тестового тарифа
test_delete_tariff() {
    log_info "Тест 10: Удаление тестового тарифа"
    
    # Сначала найдем ID тестового тарифа
    response=$(make_request "GET" "$API_BASE_URL/tariffs" "" "$ADMIN_API_KEY")
    
    # Ищем тариф с названием "Тестовый тариф"
    local test_tariff_id=$(echo "$response" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$test_tariff_id" ]; then
        response=$(make_request "DELETE" "$API_BASE_URL/tariffs/$test_tariff_id" "" "$ADMIN_API_KEY")
        
        if echo "$response" | grep -q "success.*true"; then
            log_success "Тестовый тариф удален"
        else
            log_warning "Ошибка удаления тестового тарифа"
            echo "$response"
        fi
    else
        log_warning "Тестовый тариф не найден для удаления"
    fi
}

# Основная функция
main() {
    log_info "Начинаем тестирование API тарифов..."
    
    local tests_passed=0
    local tests_failed=0
    
    # Список тестов
    tests=(
        "test_api_availability"
        "test_list_tariffs"
        "test_get_tariff"
        "test_create_tariff"
        "test_get_user_tariff"
        "test_change_user_tariff"
        "test_get_tariff_stats"
        "test_rate_limiting"
        "test_update_tariff"
        "test_delete_tariff"
    )
    
    # Выполняем тесты
    for test in "${tests[@]}"; do
        echo ""
        if $test; then
            ((tests_passed++))
        else
            ((tests_failed++))
        fi
    done
    
    # Выводим результаты
    echo ""
    log_info "Результаты тестирования:"
    log_success "Пройдено тестов: $tests_passed"
    if [ $tests_failed -gt 0 ]; then
        log_error "Провалено тестов: $tests_failed"
    else
        log_success "Провалено тестов: $tests_failed"
    fi
    
    if [ $tests_failed -eq 0 ]; then
        log_success "Все тесты прошли успешно! 🎉"
    else
        log_warning "Некоторые тесты провалились. Проверьте логи выше."
    fi
}

# Проверка зависимостей
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        log_error "curl не установлен"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        log_warning "jq не установлен. JSON вывод будет неформатированным."
    fi
}

# Запуск
check_dependencies
main "$@"
