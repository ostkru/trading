#!/bin/bash

# Скрипт для ручного обновления документации
echo "🔄 Обновление документации..."

# Останавливаем текущий сервис
echo "⏹️  Остановка Redoc сервиса..."
systemctl stop redoc-docs.service

# Ждем завершения
sleep 2

# Запускаем сервис заново
echo "▶️  Запуск Redoc сервиса..."
systemctl start redoc-docs.service

# Проверяем статус
sleep 3
if systemctl is-active --quiet redoc-docs.service; then
    echo "✅ Документация успешно обновлена!"
    echo "🌐 Доступна по адресу: http://localhost:8090"
else
    echo "❌ Ошибка запуска сервиса"
    systemctl status redoc-docs.service
fi
