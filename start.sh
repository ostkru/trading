#!/bin/bash

echo "🚀 Запуск Trading API..."
echo "📁 Директория: $(pwd)"

# Остановить существующий процесс
pkill -f "./app" 2>/dev/null
sleep 2

# Собрать приложение
echo "🔨 Сборка приложения..."
go build -o app ./cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "✅ Сборка успешна"
    
    # Запустить приложение
    echo "🚀 Запуск сервера..."
    ./app &
    
    # Проверить, что сервер запустился
    sleep 3
    if pgrep -f "./app" > /dev/null; then
        echo "✅ Сервер запущен успешно"
        echo "🌐 API доступен по адресу: http://localhost:8095"
        echo "📋 Документация: http://localhost:8095/docs.html"
    else
        echo "❌ Ошибка запуска сервера"
        exit 1
    fi
else
    echo "❌ Ошибка сборки"
    exit 1
fi
