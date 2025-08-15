#!/bin/bash

cd /var/www/trading
export PATH="/root/.nvm/versions/node/v22.16.0/bin:$PATH"

# Проверяем, что openapi.json существует
if [ ! -f "openapi.json" ]; then
    echo "Ошибка: openapi.json не найден!"
    exit 1
fi

# Останавливаем nginx если он запущен на порту 8090
if netstat -tlnp | grep -q ":8090.*nginx"; then
    echo "Остановка nginx на порту 8090..."
    systemctl stop nginx
    sleep 2
fi

# Очищаем кэш Node.js для свежего запуска
export NODE_OPTIONS="--max-old-space-size=512"

# Запускаем Redoc с оптимизированными настройками
echo "Запуск Redoc на порту 8090..."
exec npx redoc-cli serve openapi.json \
    --port 8090 \
    --host 0.0.0.0 \
    --watch \
    --no-cache
