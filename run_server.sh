#!/bin/bash

# Простой скрипт для запуска сервера в фоне
cd /var/www/go

# Остановка существующих процессов
pkill -f app_8095 2>/dev/null
sleep 2

# Запуск сервера в фоне
echo "Запуск PortalData API сервера в фоне..."
nohup ./app_8095 > server.log 2>&1 &
SERVER_PID=$!

echo "Сервер запущен с PID: $SERVER_PID"
echo "Логи: /var/www/go/server.log"
echo "Для остановки: pkill -f app_8095"
echo "Для просмотра логов: tail -f /var/www/go/server.log" 