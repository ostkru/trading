#!/bin/bash

# Скрипт для остановки Redoc сервера

PORT=8182
PID_FILE="redoc.pid"

echo "🛑 Остановка Redoc сервера..."

# Проверяем, есть ли PID файл
if [ ! -f "$PID_FILE" ]; then
    echo "ℹ️  PID файл не найден. Проверяем процессы на порту $PORT..."
    
    # Ищем процесс по порту
    PID=$(lsof -ti:$PORT 2>/dev/null)
    if [ -z "$PID" ]; then
        echo "✅ Redoc сервер не запущен"
        exit 0
    else
        echo "🔍 Найден процесс на порту $PORT (PID: $PID)"
    fi
else
    PID=$(cat "$PID_FILE")
    echo "📋 Используем PID из файла: $PID"
fi

# Останавливаем процесс
if [ ! -z "$PID" ]; then
    echo "🔄 Остановка процесса $PID..."
    kill $PID
    
    # Ждем завершения
    sleep 2
    
    # Проверяем, что процесс остановлен
    if ps -p $PID > /dev/null 2>&1; then
        echo "⚠️  Процесс не остановился, принудительно завершаем..."
        kill -9 $PID
        sleep 1
    fi
    
    # Финальная проверка
    if ps -p $PID > /dev/null 2>&1; then
        echo "❌ Не удалось остановить процесс $PID"
        exit 1
    else
        echo "✅ Redoc сервер успешно остановлен"
        rm -f "$PID_FILE"
    fi
else
    echo "❌ Не удалось найти процесс для остановки"
    exit 1
fi

# Проверяем, что порт освободился
if lsof -ti:$PORT > /dev/null 2>&1; then
    echo "⚠️  Порт $PORT все еще занят"
else
    echo "✅ Порт $PORT освобожден"
fi 