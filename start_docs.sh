#!/bin/bash

# Скрипт для запуска документации API на порту 8090

echo "🚀 Запуск документации API на порту 8090..."

# Проверяем, что порт свободен
if netstat -tlnp | grep -q ":8090"; then
    echo "❌ Порт 8090 уже занят!"
    netstat -tlnp | grep ":8090"
    exit 1
fi

# Запускаем Python HTTP сервер для отображения HTML документации
echo "📖 Запуск HTTP сервера для документации..."
python3 -m http.server 8090 --bind 0.0.0.0 --directory /var/www/gogo/go-mod &

# Ждем запуска
sleep 3

# Проверяем, что сервер запустился
if netstat -tlnp | grep -q ":8090"; then
    echo "✅ Документация успешно запущена!"
    echo "🌐 Доступна по адресу: http://localhost:8090"
    echo "📋 Файлы документации:"
    echo "   - http://localhost:8090/redoc-documentation.html (полная документация)"
    echo "   - http://localhost:8090/docs.html (краткая документация)"
    echo "   - http://localhost:8090/index.html (главная страница)"
else
    echo "❌ Ошибка запуска документации"
    exit 1
fi 