#!/bin/bash

# Скрипт для запуска API ПорталДанных.РФ сервера
# Автоматический перезапуск при падении

LOG_FILE="/var/www/go/server.log"
PID_FILE="/var/www/go/server.pid"
MAX_RESTARTS=10
RESTART_DELAY=5

echo "$(date): Запуск API ПорталДанных.РФ сервера..." >> "$LOG_FILE"

# Функция для остановки сервера
stop_server() {
    echo "$(date): Получен сигнал остановки" >> "$LOG_FILE"
    if [ -f "$PID_FILE" ]; then
        kill $(cat "$PID_FILE") 2>/dev/null
        rm -f "$PID_FILE"
    fi
    exit 0
}

# Обработка сигналов
trap stop_server SIGTERM SIGINT

# Счетчик перезапусков
restart_count=0

while [ $restart_count -lt $MAX_RESTARTS ]; do
    echo "$(date): Запуск сервера (попытка $((restart_count + 1))/$MAX_RESTARTS)" >> "$LOG_FILE"
    
    # Запуск сервера в фоне
    cd /var/www/go
    ./app_8095 >> "$LOG_FILE" 2>&1 &
    server_pid=$!
    
    # Сохранение PID
    echo $server_pid > "$PID_FILE"
    
    # Ожидание завершения процесса
    wait $server_pid
    exit_code=$?
    
    echo "$(date): Сервер завершился с кодом $exit_code" >> "$LOG_FILE"
    
    # Удаление PID файла
    rm -f "$PID_FILE"
    
    # Проверка на нормальное завершение
    if [ $exit_code -eq 0 ]; then
        echo "$(date): Сервер завершен нормально" >> "$LOG_FILE"
        break
    fi
    
    restart_count=$((restart_count + 1))
    
    if [ $restart_count -lt $MAX_RESTARTS ]; then
        echo "$(date): Перезапуск через $RESTART_DELAY секунд..." >> "$LOG_FILE"
        sleep $RESTART_DELAY
    else
        echo "$(date): Достигнут лимит перезапусков ($MAX_RESTARTS)" >> "$LOG_FILE"
    fi
done

echo "$(date): Скрипт завершен" >> "$LOG_FILE" 