#!/bin/bash

# Скрипт для запуска Redoc сервера с документацией PortalData API

PORT=8182
OPENAPI_FILE="openapi.json"
PID_FILE="redoc.pid"

echo "🚀 Запуск Redoc сервера для PortalData API..."

# Проверяем, что файл OpenAPI существует
if [ ! -f "$OPENAPI_FILE" ]; then
    echo "❌ Ошибка: Файл $OPENAPI_FILE не найден!"
    exit 1
fi

# Проверяем, не запущен ли уже сервер
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "⚠️  Redoc сервер уже запущен на порту $PORT (PID: $PID)"
        echo "🌐 Документация доступна по адресу: http://localhost:$PORT"
        exit 0
    else
        echo "🧹 Удаляем старый PID файл..."
        rm -f "$PID_FILE"
    fi
fi

# Запускаем Redoc сервер на всех интерфейсах
echo "📖 Запуск документации на порту $PORT (доступно извне)..."
npx redoc-cli serve "$OPENAPI_FILE" --port "$PORT" --host "0.0.0.0" &
REDOC_PID=$!

# Сохраняем PID
echo $REDOC_PID > "$PID_FILE"

# Ждем немного для запуска
sleep 3

# Проверяем, что сервер запустился
if ps -p $REDOC_PID > /dev/null 2>&1; then
    echo "✅ Redoc сервер успешно запущен!"
    echo "🌐 Документация доступна по адресу: http://localhost:$PORT"
    echo "📋 PID: $REDOC_PID"
    echo "💡 Для остановки сервера выполните: ./stop_redoc.sh"
else
    echo "❌ Ошибка запуска Redoc сервера"
    rm -f "$PID_FILE"
    exit 1
fi 