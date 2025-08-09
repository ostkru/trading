#!/bin/bash

# Скрипт для поддержания работы redoc сервера
LOG_FILE="/var/www/trading/redoc.log"
PID_FILE="/var/www/trading/redoc.pid"

# Устанавливаем правильный PATH
export PATH="/root/.nvm/versions/node/v22.16.0/bin:$PATH"

while true; do
    # Проверяем, работает ли сервер
    if ! curl -s http://localhost:8090/ > /dev/null 2>&1; then
        echo "$(date): Сервер не отвечает, перезапуск..." >> $LOG_FILE
        
        # Останавливаем старый процесс
        if [ -f $PID_FILE ]; then
            kill $(cat $PID_FILE) 2>/dev/null
            rm -f $PID_FILE
        fi
        
        # Запускаем новый процесс
        cd /var/www/trading
        nohup npx redoc-cli serve openapi.json --port 8090 --host 0.0.0.0 >> $LOG_FILE 2>&1 &
        echo $! > $PID_FILE
        echo "$(date): Сервер перезапущен с PID $(cat $PID_FILE)" >> $LOG_FILE
    fi
    
    sleep 60
done
