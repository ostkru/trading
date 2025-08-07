#!/bin/bash

echo "🧪 Запуск тестов Trading API..."

# Проверить, что сервер запущен
if ! pgrep -f "./app" > /dev/null; then
    echo "❌ Сервер не запущен. Запустите ./start.sh сначала"
    exit 1
fi

echo "✅ Сервер запущен, запуск тестов..."

# Запустить PHP тесты
php comprehensive_api_test.php

echo "🎉 Тестирование завершено"
