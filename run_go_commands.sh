#!/bin/bash

echo "=== ЗАПУСК GO API НА СЕРВЕРЕ ==="

cd /var/www/go

echo "1. Останавливаю старые процессы..."
pkill -f app
sleep 2

echo "2. Проверяю файлы..."
ls -la

echo "3. Обновляю зависимости..."
go mod tidy

echo "4. Собираю приложение..."
go build -o app cmd/api/main.go

echo "5. Запускаю приложение..."
nohup ./app > app.log 2>&1 &

echo "6. Жду запуска..."
sleep 5

echo "7. Проверяю процессы..."
ps aux | grep app

echo "8. Проверяю порт..."
netstat -tlnp | grep 8095

echo "9. Проверяю логи..."
tail -10 app.log

echo "=== ГОТОВО ===" 