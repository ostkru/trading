#!/bin/bash

# Скрипт для мониторинга и перезапуска redoc сервера
LOG_FILE="/var/www/trading/redoc.log"
PID_FILE="/var/www/trading/redoc.pid"

# Функция для запуска сервера
start_server() {
    echo "$(date): Запуск redoc сервера..." >> $LOG_FILE
    cd /var/www/trading
    nohup npx redoc-cli serve openapi.json --port 8090 --host 0.0.0.0 >> $LOG_FILE 2>&1 &
    echo $! > $PID_FILE
    echo "$(date): Сервер запущен с PID $(cat $PID_FILE)" >> $LOG_FILE
}

# Функция для остановки сервера
stop_server() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        if kill -0 $PID 2>/dev/null; then
            kill $PID
            echo "$(date): Сервер остановлен (PID: $PID)" >> $LOG_FILE
        fi
        rm -f $PID_FILE
    fi
}

# Функция для проверки работы сервера
check_server() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        if kill -0 $PID 2>/dev/null; then
            # Проверяем, отвечает ли сервер на HTTP запросы
            if curl -s http://localhost:8090/ > /dev/null 2>&1; then
                return 0
            else
                echo "$(date): Сервер не отвечает на HTTP запросы" >> $LOG_FILE
                return 1
            fi
        else
            echo "$(date): Процесс не найден (PID: $PID)" >> $LOG_FILE
            return 1
        fi
    else
        echo "$(date): PID файл не найден" >> $LOG_FILE
        return 1
    fi
}

# Основной цикл мониторинга
while true; do
    if ! check_server; then
        echo "$(date): Сервер не работает, перезапуск..." >> $LOG_FILE
        stop_server
        sleep 2
        start_server
    fi
    sleep 30
done
