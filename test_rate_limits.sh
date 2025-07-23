#!/bin/bash

# Тест для проверки лимитов API
# API ключ: 026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f

API_KEY="026b26ac7a206c51a216b3280042cda5178710912da68ae696a713970034dd5f"
BASE_URL="http://localhost:8095"

echo "🚀 ТЕСТИРОВАНИЕ ЛИМИТОВ API"
echo "=========================================="
echo "API ключ: $API_KEY"
echo ""

# Функция для выполнения запроса
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    local url="$BASE_URL$endpoint"
    local headers=(
        "Authorization: Bearer $API_KEY"
        "Content-Type: application/json"
    )
    
    local curl_cmd="curl -s -w '\nHTTP_STATUS:%{http_code}\n' -X $method"
    
    for header in "${headers[@]}"; do
        curl_cmd="$curl_cmd -H '$header'"
    done
    
    if [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -d '$data'"
    fi
    
    curl_cmd="$curl_cmd '$url'"
    
    local response=$(eval $curl_cmd)
    local http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d: -f2)
    local body=$(echo "$response" | grep -v "HTTP_STATUS:")
    
    echo "$http_status"
}

# Тест 1: Базовая проверка доступа
echo "📋 1. БАЗОВАЯ ПРОВЕРКА ДОСТУПА"
echo "--------------------------------"

status=$(make_request "GET" "/api/v1/products")
echo "Получение списка продуктов: HTTP $status"

if [ "$status" = "200" ]; then
    echo "✅ ПРОЙДЕН"
else
    echo "❌ ПРОВАЛЕН"
fi

echo ""

# Тест 2: Проверка минутных лимитов
echo "⏱️  2. ПРОВЕРКА МИНУТНЫХ ЛИМИТОВ"
echo "--------------------------------"

endpoints=(
    "/api/v1/products"
    "/api/v1/offers"
    "/api/v1/orders"
    "/api/v1/warehouses"
)

for endpoint in "${endpoints[@]}"; do
    status=$(make_request "GET" "$endpoint")
    echo "GET $endpoint: HTTP $status"
    
    if [ "$status" = "200" ]; then
        echo "✅ ПРОЙДЕН"
    elif [ "$status" = "429" ]; then
        echo "⚠️  ЛИМИТ ПРЕВЫШЕН"
    else
        echo "❌ ПРОВАЛЕН"
    fi
    
    sleep 0.1
done

echo ""

# Тест 3: Проверка дневных лимитов
echo "📅 3. ПРОВЕРКА ДНЕВНЫХ ЛИМИТОВ"
echo "--------------------------------"

daily_endpoints=(
    "/api/v1/offers?filter=all"
    "/api/v1/offers/public"
)

for endpoint in "${daily_endpoints[@]}"; do
    status=$(make_request "GET" "$endpoint")
    echo "GET $endpoint: HTTP $status"
    
    if [ "$status" = "200" ]; then
        echo "✅ ПРОЙДЕН"
    elif [ "$status" = "429" ]; then
        echo "⚠️  ДНЕВНОЙ ЛИМИТ ПРЕВЫШЕН"
    else
        echo "❌ ПРОВАЛЕН"
    fi
    
    sleep 0.1
done

echo ""

# Тест 4: Стресс-тест
echo "💥 4. СТРЕСС-ТЕСТ ЛИМИТОВ"
echo "--------------------------------"

success_count=0
limit_count=0
error_count=0

echo "Отправка множественных запросов..."

for i in {1..20}; do
    status=$(make_request "GET" "/api/v1/products")
    
    if [ "$status" = "200" ]; then
        ((success_count++))
    elif [ "$status" = "429" ]; then
        ((limit_count++))
    else
        ((error_count++))
    fi
    
    if [ $((i % 5)) -eq 0 ]; then
        echo "Запрос $i: Успешно: $success_count, Лимит: $limit_count, Ошибки: $error_count"
    fi
    
    sleep 0.05
done

echo ""
echo "📊 ИТОГИ СТРЕСС-ТЕСТА:"
echo "Успешных запросов: $success_count"
echo "Запросов с превышением лимита: $limit_count"
echo "Ошибок: $error_count"

echo ""
echo "✅ ТЕСТИРОВАНИЕ ЛИМИТОВ ЗАВЕРШЕНО"
echo "==========================================" 