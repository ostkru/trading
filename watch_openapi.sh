#!/bin/bash

# Скрипт для автоматического обновления документации при изменении openapi.json
OPENAPI_FILE="/var/www/trading/openapi.json"
LOG_FILE="/var/log/openapi-watcher.log"

echo "$(date): Запуск мониторинга openapi.json..." >> "$LOG_FILE"

# Функция для генерации документации
generate_docs() {
    echo "$(date): Обнаружены изменения в openapi.json, генерация документации..." >> "$LOG_FILE"
    
    # Генерируем новую документацию
    cd /var/www/trading
    if npx redoc-cli build openapi.json -o redoc-documentation.html; then
        echo "$(date): Документация успешно сгенерирована" >> "$LOG_FILE"
    else
        echo "$(date): Ошибка генерации документации" >> "$LOG_FILE"
    fi
}

# Проверяем, что файл существует
if [ ! -f "$OPENAPI_FILE" ]; then
    echo "$(date): Ошибка: файл $OPENAPI_FILE не найден" >> "$LOG_FILE"
    exit 1
fi

# Получаем начальный хеш файла
INITIAL_HASH=$(md5sum "$OPENAPI_FILE" | awk '{print $1}')
echo "$(date): Начальный хеш: $INITIAL_HASH" >> "$LOG_FILE"

# Бесконечный цикл мониторинга
while true; do
    sleep 5  # Проверяем каждые 5 секунд
    
    # Проверяем, что файл все еще существует
    if [ ! -f "$OPENAPI_FILE" ]; then
        echo "$(date): Файл $OPENAPI_FILE удален, ожидание..." >> "$LOG_FILE"
        continue
    fi
    
    # Получаем текущий хеш
    CURRENT_HASH=$(md5sum "$OPENAPI_FILE" | awk '{print $1}')
    
    # Если хеш изменился, генерируем новую документацию
    if [ "$CURRENT_HASH" != "$INITIAL_HASH" ]; then
        echo "$(date): Обнаружены изменения в openapi.json" >> "$LOG_FILE"
        echo "$(date): Старый хеш: $INITIAL_HASH" >> "$LOG_FILE"
        echo "$(date): Новый хеш: $CURRENT_HASH" >> "$LOG_FILE"
        
        generate_docs
        
        # Обновляем хеш
        INITIAL_HASH="$CURRENT_HASH"
        echo "$(date): Хеш обновлен: $INITIAL_HASH" >> "$LOG_FILE"
    fi
done
