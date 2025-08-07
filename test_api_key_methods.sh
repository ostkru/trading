#!/bin/bash

echo "🧪 ТЕСТИРОВАНИЕ РАЗЛИЧНЫХ СПОСОБОВ ПЕРЕДАЧИ API КЛЮЧА"
echo "=================================================="

API_KEY="f428fbc16a97b9e2a55717bd34e97537ec34cb8c04a5f32eeb4e88c9ee998a53"
BASE_URL="http://localhost:8095"

echo ""
echo "1. Тест с Authorization: Bearer header"
echo "----------------------------------------"
curl -s -w "HTTP Status: %{http_code}\n" \
  -H "Authorization: Bearer $API_KEY" \
  "$BASE_URL/api/v1/offers" | head -2

echo ""
echo "2. Тест с X-API-KEY header"
echo "----------------------------"
curl -s -w "HTTP Status: %{http_code}\n" \
  -H "X-API-KEY: $API_KEY" \
  "$BASE_URL/api/v1/offers" | head -2

echo ""
echo "3. Тест с api_key в GET параметрах"
echo "-----------------------------------"
curl -s -w "HTTP Status: %{http_code}\n" \
  "$BASE_URL/api/v1/offers?api_key=$API_KEY" | head -2

echo ""
echo "4. Тест без авторизации (должен вернуть ошибку)"
echo "------------------------------------------------"
curl -s -w "HTTP Status: %{http_code}\n" \
  "$BASE_URL/api/v1/offers"

echo ""
echo "5. Тест с неверным API ключом"
echo "------------------------------"
curl -s -w "HTTP Status: %{http_code}\n" \
  -H "Authorization: Bearer invalid_key" \
  "$BASE_URL/api/v1/offers"

echo ""
echo "✅ ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo "==========================" 