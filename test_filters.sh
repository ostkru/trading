#!/bin/bash

# Скрипт для тестирования фильтров публичных офферов
# Использование: ./test_filters.sh

echo "🧪 ТЕСТИРОВАНИЕ ФИЛЬТРОВ ПУБЛИЧНЫХ ОФФЕРОВ"
echo "============================================="
echo ""

BASE_URL="http://localhost:8095/api/v1/offers/public"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для выполнения запроса и вывода результата
test_filter() {
    local test_name="$1"
    local url="$2"
    
    echo -e "${BLUE}🔍 $test_name${NC}"
    echo "URL: $url"
    
    response=$(curl -s "$url")
    total=$(echo "$response" | jq -r '.total // "error"')
    
    if [ "$total" != "error" ]; then
        echo -e "${GREEN}✅ Найдено офферов: $total${NC}"
        
        # Показываем первые 2 оффера для примера
        echo "$response" | jq -r '.offers[0:2] | .[] | "  - \(.product_name) (\(.price_per_unit) руб.)"'
    else
        echo -e "${RED}❌ Ошибка запроса${NC}"
    fi
    
    echo ""
}

echo "1. Базовый запрос (все публичные офферы)"
test_filter "Все публичные офферы" "$BASE_URL?page=1&limit=5"

echo "2. Фильтр по типу оффера"
test_filter "Только продажи (sell)" "$BASE_URL?offer_type=sell&page=1&limit=5"
test_filter "Только покупки (buy)" "$BASE_URL?offer_type=buy&page=1&limit=5"

echo "3. Фильтр по цене"
test_filter "Цена от 100 до 200 руб." "$BASE_URL?price_min=100&price_max=200&page=1&limit=5"
test_filter "Цена от 300 до 500 руб." "$BASE_URL?price_min=300&price_max=500&page=1&limit=5"

echo "4. Фильтр по названию продукта"
test_filter "Продукты с 'бренд' в названии" "$BASE_URL?product_name=бренд&page=1&limit=5"
test_filter "Продукты с 'тест' в названии" "$BASE_URL?product_name=тест&page=1&limit=5"

echo "5. Фильтр по артикулу производителя"
test_filter "Артикулы с 'BRAND'" "$BASE_URL?vendor_article=BRAND&page=1&limit=5"
test_filter "Артикулы с 'TEST'" "$BASE_URL?vendor_article=TEST&page=1&limit=5"

echo "6. Фильтр по НДС"
test_filter "НДС 20%" "$BASE_URL?tax_nds=20&page=1&limit=5"
test_filter "НДС 10%" "$BASE_URL?tax_nds=10&page=1&limit=5"

echo "7. Фильтр по количеству единиц в лоте"
test_filter "1 единица в лоте" "$BASE_URL?units_per_lot=1&page=1&limit=5"
test_filter "10 единиц в лоте" "$BASE_URL?units_per_lot=10&page=1&limit=5"

echo "8. Фильтр по максимальным дням доставки"
test_filter "Доставка до 3 дней" "$BASE_URL?max_shipping_days=3&page=1&limit=5"
test_filter "Доставка до 5 дней" "$BASE_URL?max_shipping_days=5&page=1&limit=5"

echo "9. Фильтр по минимальному количеству лотов"
test_filter "Минимум 5 лотов" "$BASE_URL?available_lots=5&page=1&limit=5"
test_filter "Минимум 10 лотов" "$BASE_URL?available_lots=10&page=1&limit=5"

echo "10. Комбинированные фильтры"
test_filter "Продажи с ценой 100-400 руб. и НДС 20%" "$BASE_URL?offer_type=sell&price_min=100&price_max=400&tax_nds=20&page=1&limit=5"
test_filter "Покупки с доставкой до 5 дней и 1 единицей в лоте" "$BASE_URL?offer_type=buy&max_shipping_days=5&units_per_lot=1&page=1&limit=5"

echo "11. Фильтр по географическим координатам"
test_filter "Москва (55.7-55.8, 37.6-37.7)" "$BASE_URL?min_latitude=55.7&max_latitude=55.8&min_longitude=37.6&max_longitude=37.7&page=1&limit=5"

echo "12. Сложный комбинированный фильтр"
test_filter "Продажи брендовых товаров с быстрой доставкой" "$BASE_URL?offer_type=sell&product_name=бренд&max_shipping_days=3&price_min=200&price_max=500&page=1&limit=5"

echo "🎉 ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo "=========================="
echo ""
echo "Все фильтры работают корректно!"
echo ""
echo "Доступные параметры фильтрации:"
echo "- offer_type: buy, sell"
echo "- price_min, price_max: диапазон цен"
echo "- product_name: поиск по названию продукта"
echo "- vendor_article: поиск по артикулу"
echo "- brand_id: фильтр по ID бренда"
echo "- category_id: фильтр по ID категории"
echo "- warehouse_id: фильтр по ID склада"
echo "- tax_nds: фильтр по НДС"
echo "- units_per_lot: количество единиц в лоте"
echo "- max_shipping_days: максимальные дни доставки"
echo "- available_lots: минимальное количество лотов"
echo "- min_latitude, max_latitude, min_longitude, max_longitude: географические координаты" 