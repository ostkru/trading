#!/bin/bash

# Скрипт для генерации статической документации из openapi.json
cd /var/www/trading

echo "🔄 Генерация документации из openapi.json..."

# Проверяем, что openapi.json существует
if [ ! -f "openapi.json" ]; then
    echo "❌ Ошибка: openapi.json не найден!"
    exit 1
fi

# Генерируем HTML документацию используя старую версию Redoc
echo "📝 Генерация HTML документации..."
npx redoc-cli build openapi.json -o redoc-documentation.html

if [ $? -eq 0 ]; then
    echo "✅ Документация успешно сгенерирована!"
    echo "🌐 Доступна по адресу: http://docs.portaldata.ru:8090"
else
    echo "❌ Ошибка генерации документации"
    exit 1
fi
