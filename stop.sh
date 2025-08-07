#!/bin/bash

echo "🛑 Остановка Trading API..."

# Найти и остановить процесс
if pgrep -f "./app" > /dev/null; then
    echo "�� Найден процесс приложения"
    pkill -f "./app"
    sleep 2
    
    # Проверить, что процесс остановлен
    if pgrep -f "./app" > /dev/null; then
        echo "⚠️  Процесс не остановлен, принудительное завершение..."
        pkill -9 -f "./app"
    fi
    
    echo "✅ Приложение остановлено"
else
    echo "ℹ️  Процесс приложения не найден"
fi
