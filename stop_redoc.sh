#!/bin/bash

# Скрипт для остановки Redoc документации
# Автор: AI Assistant
# Дата: $(date)

echo "🛑 Остановка Redoc документации..."

# Ищем процессы redoc-cli
REDOC_PIDS=$(pgrep -f "redoc-cli serve")

if [ -z "$REDOC_PIDS" ]; then
    echo "ℹ️  Redoc сервер не запущен"
    exit 0
fi

# Останавливаем процессы
echo "📋 Найдены процессы Redoc: $REDOC_PIDS"
for pid in $REDOC_PIDS; do
    echo "🔄 Остановка процесса $pid..."
    kill $pid 2>/dev/null
done

# Ждем завершения процессов
sleep 2

# Проверяем, что процессы остановлены
REMAINING_PIDS=$(pgrep -f "redoc-cli serve")
if [ -z "$REMAINING_PIDS" ]; then
    echo "✅ Redoc сервер успешно остановлен"
else
    echo "⚠️  Некоторые процессы не остановлены: $REMAINING_PIDS"
    echo "🔄 Принудительная остановка..."
    kill -9 $REMAINING_PIDS 2>/dev/null
    echo "✅ Все процессы остановлены"
fi

# Проверяем, что порт 8090 свободен
if ! netstat -tlnp 2>/dev/null | grep -q ":8090"; then
    echo "✅ Порт 8090 свободен"
else
    echo "⚠️  Порт 8090 все еще занят"
fi

