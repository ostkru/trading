#!/bin/bash

# Скрипт для тестирования продакшена API
# Использование: ./test_production.sh

echo "🧪 ТЕСТИРОВАНИЕ ПРОДАКШЕНА API"
echo "================================"
echo ""

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Базовый URL продакшена
PROD_URL="https://api.portaldata.ru/v1/trading"
API_KEY="f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53"

# Функция для выполнения запроса
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -n "$data" ]; then
        curl -s -X "$method" \
             -H "Content-Type: application/json" \
             -H "Authorization: Bearer $API_KEY" \
             -d "$data" \
             "$PROD_URL$endpoint"
    else
        curl -s -X "$method" \
             -H "Authorization: Bearer $API_KEY" \
             "$PROD_URL$endpoint"
    fi
}

# Функция для проверки статуса
check_status() {
    local test_name=$1
    local response=$2
    local expected_status=$3
    
    status=$(echo "$response" | jq -r '.status // .error // "unknown"')
    
    if [ "$status" = "$expected_status" ]; then
        echo -e "${GREEN}✅ $test_name: OK${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name: FAILED (expected: $expected_status, got: $status)${NC}"
        echo "Response: $response"
        return 1
    fi
}

echo "1. Проверка доступности API..."
response=$(make_request "GET" "/offers/public")
check_status "API доступен" "$response" "200"

echo ""
echo "2. Тестирование продуктов..."
product_data='{
    "name": "Тестовый продукт продакшена",
    "vendor_article": "PROD-TEST-001",
    "recommend_price": 150.50,
    "brand": "TestBrand",
    "category": "TestCategory",
    "description": "Описание тестового продукта для продакшена"
}'

response=$(make_request "POST" "/products" "$product_data")
check_status "Создание продукта" "$response" "201"

if [ $? -eq 0 ]; then
    product_id=$(echo "$response" | jq -r '.id')
    echo "   Создан продукт с ID: $product_id"
    
    # Получение продукта
    response=$(make_request "GET" "/products/$product_id")
    check_status "Получение продукта" "$response" "200"
    
    # Обновление продукта
    update_data='{"name": "Обновленный продукт", "description": "Обновленное описание"}'
    response=$(make_request "PUT" "/products/$product_id" "$update_data")
    check_status "Обновление продукта" "$response" "200"
fi

echo ""
echo "3. Тестирование складов..."
warehouse_data='{
    "name": "Тестовый склад продакшена",
    "address": "г. Москва, ул. Тестовая, д. 1",
    "latitude": 55.7558,
    "longitude": 37.6176,
    "working_hours": "09:00-18:00"
}'

response=$(make_request "POST" "/warehouses" "$warehouse_data")
check_status "Создание склада" "$response" "201"

if [ $? -eq 0 ]; then
    warehouse_id=$(echo "$response" | jq -r '.id')
    echo "   Создан склад с ID: $warehouse_id"
    
    # Получение склада
    response=$(make_request "GET" "/warehouses/$warehouse_id")
    check_status "Получение склада" "$response" "200"
fi

echo ""
echo "4. Тестирование офферов..."
offer_data='{
    "product_id": 1,
    "offer_type": "sell",
    "price_per_unit": 100.00,
    "units_per_lot": 10,
    "available_lots": 5,
    "is_public": true
}'

response=$(make_request "POST" "/offers" "$offer_data")
check_status "Создание оффера" "$response" "201"

if [ $? -eq 0 ]; then
    offer_id=$(echo "$response" | jq -r '.id')
    echo "   Создан оффер с ID: $offer_id"
    
    # Получение оффера
    response=$(make_request "GET" "/offers/$offer_id")
    check_status "Получение оффера" "$response" "200"
fi

echo ""
echo "5. Тестирование публичных офферов..."
response=$(make_request "GET" "/offers/public")
check_status "Публичные офферы" "$response" "200"

echo ""
echo "6. Тестирование статистики..."
response=$(make_request "GET" "/statistics")
check_status "Статистика пользователя" "$response" "200"

echo ""
echo "7. Проверка документации..."
doc_response=$(curl -s -o /dev/null -w "%{http_code}" "$PROD_URL/docs")
if [ "$doc_response" = "200" ]; then
    echo -e "${GREEN}✅ Документация доступна${NC}"
else
    echo -e "${YELLOW}⚠️  Документация недоступна (статус: $doc_response)${NC}"
fi

echo ""
echo "8. Проверка CORS..."
cors_response=$(curl -s -X OPTIONS \
    -H "Origin: https://example.com" \
    -H "Access-Control-Request-Method: GET" \
    -H "Access-Control-Request-Headers: Authorization" \
    "$PROD_URL/products" -o /dev/null -w "%{http_code}")

if [ "$cors_response" = "204" ] || [ "$cors_response" = "200" ]; then
    echo -e "${GREEN}✅ CORS настроен корректно${NC}"
else
    echo -e "${YELLOW}⚠️  CORS может быть не настроен (статус: $cors_response)${NC}"
fi

echo ""
echo "9. Проверка SSL..."
ssl_check=$(curl -s -o /dev/null -w "%{http_code}" "$PROD_URL/products")
if [ "$ssl_check" != "000" ]; then
    echo -e "${GREEN}✅ SSL работает корректно${NC}"
else
    echo -e "${RED}❌ SSL не работает${NC}"
fi

echo ""
echo "🎉 ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo "=========================="
echo ""
echo "Для полного тестирования запустите:"
echo "php comprehensive_api_test.php"
echo ""
echo "Документация доступна по адресу:"
echo "$PROD_URL/docs" 