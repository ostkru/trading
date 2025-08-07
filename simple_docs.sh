#!/bin/bash

echo "�� Запуск простого HTTP сервера для документации..."

# Останавливаем все процессы на порту 8090
pkill -f "python3 -m http.server 8090" 2>/dev/null
pkill -f "nginx.*8090" 2>/dev/null

# Ждем освобождения порта
sleep 2

# Запускаем простой HTTP сервер с таймаутом
timeout 3600 python3 -m http.server 8090 --bind 0.0.0.0 --directory /var/www/gogo/go-mod &

# Ждем запуска
sleep 3

# Проверяем
if netstat -tlnp | grep -q ":8090"; then
    echo "✅ Документация запущена на порту 8090"
    echo "🌐 Доступна по адресу: http://92.53.64.38:8090"
    echo "📋 Файлы:"
    echo "   - http://92.53.64.38:8090/docs.html"
    echo "   - http://92.53.64.38:8090/redoc-documentation.html"
    echo "   - http://92.53.64.38:8090/index.html"
else
    echo "❌ Ошибка запуска документации"
fi
