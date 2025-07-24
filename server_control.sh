#!/bin/bash

# Скрипт управления PortalData API сервером

PID_FILE="/var/www/go/server.pid"
LOG_FILE="/var/www/go/server.log"
SCRIPT_DIR="/var/www/go"

case "$1" in
    start)
        echo "Запуск PortalData API сервера..."
        if [ -f "$PID_FILE" ]; then
            echo "Сервер уже запущен (PID: $(cat $PID_FILE))"
            exit 1
        fi
        
        cd "$SCRIPT_DIR"
        nohup ./start_server.sh > /dev/null 2>&1 &
        echo "Сервер запущен в фоновом режиме"
        echo "Логи: $LOG_FILE"
        ;;
        
    stop)
        echo "Остановка PortalData API сервера..."
        if [ -f "$PID_FILE" ]; then
            kill $(cat "$PID_FILE") 2>/dev/null
            rm -f "$PID_FILE"
            echo "Сервер остановлен"
        else
            echo "Сервер не запущен"
        fi
        ;;
        
    restart)
        echo "Перезапуск PortalData API сервера..."
        $0 stop
        sleep 2
        $0 start
        ;;
        
    status)
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                echo "Сервер запущен (PID: $pid)"
                echo "Логи: $LOG_FILE"
                echo "Последние строки лога:"
                tail -5 "$LOG_FILE" 2>/dev/null || echo "Лог файл не найден"
            else
                echo "Сервер не запущен (PID файл устарел)"
                rm -f "$PID_FILE"
            fi
        else
            echo "Сервер не запущен"
        fi
        ;;
        
    logs)
        if [ -f "$LOG_FILE" ]; then
            tail -f "$LOG_FILE"
        else
            echo "Лог файл не найден: $LOG_FILE"
        fi
        ;;
        
    build)
        echo "Сборка сервера..."
        cd "$SCRIPT_DIR"
        go build -o app_8095 cmd/api/main.go
        if [ $? -eq 0 ]; then
            echo "Сборка завершена успешно"
        else
            echo "Ошибка сборки"
            exit 1
        fi
        ;;
        
    *)
        echo "Использование: $0 {start|stop|restart|status|logs|build}"
        echo ""
        echo "Команды:"
        echo "  start   - запустить сервер"
        echo "  stop    - остановить сервер"
        echo "  restart - перезапустить сервер"
        echo "  status  - показать статус сервера"
        echo "  logs    - показать логи в реальном времени"
        echo "  build   - собрать сервер"
        exit 1
        ;;
esac 