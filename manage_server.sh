#!/bin/bash

# Скрипт управления PortalData API сервером
# Использование: ./manage_server.sh [start|stop|restart|status|logs|test]

PID_FILE="/var/www/go/server.pid"
LOG_FILE="/var/www/go/server.log"
SCRIPT_DIR="/var/www/go"

case "$1" in
    start)
        echo "🚀 Запуск PortalData API сервера..."
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                echo "❌ Сервер уже запущен (PID: $pid)"
                exit 1
            else
                rm -f "$PID_FILE"
            fi
        fi
        
        cd "$SCRIPT_DIR"
        nohup ./app_8095 > "$LOG_FILE" 2>&1 &
        server_pid=$!
        echo $server_pid > "$PID_FILE"
        
        sleep 2
        if kill -0 "$server_pid" 2>/dev/null; then
            echo "✅ Сервер успешно запущен (PID: $server_pid)"
            echo "📋 Логи: $LOG_FILE"
            echo "🌐 URL: http://localhost:8095"
        else
            echo "❌ Ошибка запуска сервера"
            rm -f "$PID_FILE"
            exit 1
        fi
        ;;
        
    stop)
        echo "🛑 Остановка PortalData API сервера..."
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                kill "$pid"
                echo "✅ Сервер остановлен (PID: $pid)"
            else
                echo "⚠️  Процесс не найден, но PID файл существует"
            fi
            rm -f "$PID_FILE"
        else
            # Попробуем найти и остановить процесс
            pkill -f app_8095 2>/dev/null
            echo "✅ Все процессы сервера остановлены"
        fi
        ;;
        
    restart)
        echo "🔄 Перезапуск PortalData API сервера..."
        $0 stop
        sleep 3
        $0 start
        ;;
        
    status)
        echo "📊 Статус PortalData API сервера:"
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                echo "✅ Сервер запущен (PID: $pid)"
                echo "📋 Логи: $LOG_FILE"
                echo "🌐 URL: http://localhost:8095"
                echo ""
                echo "📈 Последние строки лога:"
                tail -5 "$LOG_FILE" 2>/dev/null || echo "Лог файл не найден"
            else
                echo "❌ Сервер не запущен (PID файл устарел)"
                rm -f "$PID_FILE"
            fi
        else
            # Проверим, есть ли процессы
            if pgrep -f app_8095 > /dev/null; then
                echo "⚠️  Сервер запущен, но PID файл отсутствует"
                pgrep -f app_8095
            else
                echo "❌ Сервер не запущен"
            fi
        fi
        ;;
        
    logs)
        if [ -f "$LOG_FILE" ]; then
            echo "📋 Логи сервера (Ctrl+C для выхода):"
            tail -f "$LOG_FILE"
        else
            echo "❌ Лог файл не найден: $LOG_FILE"
        fi
        ;;
        
    test)
        echo "🧪 Тестирование API сервера..."
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            if kill -0 "$pid" 2>/dev/null; then
                echo "✅ Сервер запущен, тестируем API..."
                response=$(curl -s -w "%{http_code}" http://localhost:8095/api/v1/products)
                if [[ $response == *"401"* ]]; then
                    echo "✅ API отвечает (требуется авторизация)"
                elif [[ $response == *"200"* ]]; then
                    echo "✅ API отвечает успешно"
                else
                    echo "❌ API не отвечает"
                fi
            else
                echo "❌ Сервер не запущен"
            fi
        else
            echo "❌ Сервер не запущен"
        fi
        ;;
        
    build)
        echo "🔨 Сборка сервера..."
        cd "$SCRIPT_DIR"
        go build -o app_8095 cmd/api/main.go
        if [ $? -eq 0 ]; then
            echo "✅ Сборка завершена успешно"
        else
            echo "❌ Ошибка сборки"
            exit 1
        fi
        ;;
        
    *)
        echo "📖 Использование: $0 {start|stop|restart|status|logs|test|build}"
        echo ""
        echo "Команды:"
        echo "  start   - запустить сервер"
        echo "  stop    - остановить сервер"
        echo "  restart - перезапустить сервер"
        echo "  status  - показать статус сервера"
        echo "  logs    - показать логи в реальном времени"
        echo "  test    - протестировать API"
        echo "  build   - собрать сервер"
        echo ""
        echo "Примеры:"
        echo "  $0 start    # Запустить сервер"
        echo "  $0 status   # Проверить статус"
        echo "  $0 logs     # Смотреть логи"
        echo "  $0 test     # Тестировать API"
        exit 1
        ;;
esac 