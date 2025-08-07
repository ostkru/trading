#!/bin/bash

echo "🧪 ТЕСТИРОВАНИЕ СИСТЕМЫ СТАТУСОВ ПРОДУКТОВ"
echo "============================================="

API_KEY="f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53"
BASE_URL="http://localhost:8095"

echo ""
echo "1. Проверка отображения статуса в продуктах"
echo "-------------------------------------------"
curl -s "$BASE_URL/api/v1/products/1?api_key=$API_KEY" | jq '.status'

echo ""
echo "2. Тест фильтра по статусу 'not_classified'"
echo "-------------------------------------------"
curl -s "$BASE_URL/api/v1/products?owner=not_classified&api_key=$API_KEY" | jq '.total'

echo ""
echo "3. Тест фильтра по статусу 'classified'"
echo "----------------------------------------"
curl -s "$BASE_URL/api/v1/products?owner=classified&api_key=$API_KEY" | jq '.total'

echo ""
echo "4. Тест фильтра по статусу 'pending'"
echo "------------------------------------"
curl -s "$BASE_URL/api/v1/products?owner=pending&api_key=$API_KEY" | jq '.total'

echo ""
echo "5. Тест создания оффера с неклассифицированным продуктом"
echo "--------------------------------------------------------"
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $API_KEY" \
  -d '{"product_id": 1, "offer_type": "sale", "price_per_unit": 100, "available_lots": 5, "tax_nds": 20, "units_per_lot": 1, "warehouse_id": 1}' \
  "$BASE_URL/api/v1/offers" | jq '.error'

echo ""
echo "6. Тест создания оффера с классифицированным продуктом"
echo "------------------------------------------------------"
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $API_KEY" \
  -d '{"product_id": 423, "offer_type": "sale", "price_per_unit": 100, "available_lots": 5, "tax_nds": 20, "units_per_lot": 1, "warehouse_id": 1}' \
  "$BASE_URL/api/v1/offers" | jq '.offer_id'

echo ""
echo "7. Тест создания нового продукта (должен получить статус 'pending')"
echo "-------------------------------------------------------------------"
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $API_KEY" \
  -d '{"name": "Тестовый продукт для статусов", "vendor_article": "STATUS-TEST-001", "recommend_price": 150, "brand": "TestBrand", "category": "TestCategory", "description": "Тест статусов"}' \
  "$BASE_URL/api/v1/products" | jq '.status'

echo ""
echo "✅ ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo "==========================" 