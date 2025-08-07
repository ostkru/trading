#!/bin/bash

# Скрипт для запуска Redoc документации
# Автор: AI Assistant
# Дата: $(date)

echo "🚀 Запуск Redoc документации на порту 8090..."

# Проверяем, что npx доступен
if ! command -v npx &> /dev/null; then
    echo "❌ Ошибка: npx не найден. Установите Node.js и npm"
    exit 1
fi

# Проверяем наличие openapi.json
if [ ! -f "openapi.json" ]; then
    echo "❌ Ошибка: Файл openapi.json не найден"
    exit 1
fi

# Останавливаем предыдущий процесс, если он запущен
pkill -f "redoc-cli serve" 2>/dev/null
sleep 2

# Запускаем redoc сервер с внешним доступом
echo "📖 Запуск сервера документации с внешним доступом..."
npx redoc-cli serve openapi.json --port 8090 --host 0.0.0.0 > redoc.log 2>&1 &

# Ждем немного и проверяем статус
sleep 3
if curl -s http://localhost:8090 > /dev/null; then
    echo "✅ Redoc документация успешно запущена!"
    echo "🌐 Локальный доступ: http://localhost:8090"
    echo "🌐 Внешний доступ: http://92.53.64.38:8090"
    echo "📝 Логи сохраняются в файл: redoc.log"
else
    echo "❌ Ошибка запуска сервера. Проверьте логи в redoc.log"
    exit 1
fi
