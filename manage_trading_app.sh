#!/bin/bash

# Скрипт управления торговым приложением
# Использование: ./manage_trading_app.sh [start|stop|restart|status|logs]

SERVICE_NAME="trading-app"

case "$1" in
    start)
        echo "🚀 Запуск торгового приложения..."
        systemctl start $SERVICE_NAME
        sleep 2
        if systemctl is-active --quiet $SERVICE_NAME; then
            echo "✅ Приложение успешно запущено на порту 8095"
            lsof -i :8095
        else
            echo "❌ Ошибка запуска приложения"
            systemctl status $SERVICE_NAME
        fi
        ;;
    stop)
        echo "🛑 Остановка торгового приложения..."
        systemctl stop $SERVICE_NAME
        echo "✅ Приложение остановлено"
        ;;
    restart)
        echo "🔄 Перезапуск торгового приложения..."
        systemctl restart $SERVICE_NAME
        sleep 2
        if systemctl is-active --quiet $SERVICE_NAME; then
            echo "✅ Приложение успешно перезапущено на порту 8095"
            lsof -i :8095
        else
            echo "❌ Ошибка перезапуска приложения"
            systemctl status $SERVICE_NAME
        fi
        ;;
    status)
        echo "📊 Статус торгового приложения:"
        systemctl status $SERVICE_NAME --no-pager -l
        echo ""
        echo "🌐 Порт 8095:"
        lsof -i :8095 2>/dev/null || echo "Порт свободен"
        ;;
    logs)
        echo "📝 Последние логи торгового приложения:"
        journalctl -u $SERVICE_NAME -n 30 --no-pager
        ;;
    *)
        echo "Использование: $0 {start|stop|restart|status|logs}"
        echo ""
        echo "Команды:"
        echo "  start   - запустить приложение"
        echo "  stop    - остановить приложение"
        echo "  restart - перезапустить приложение"
        echo "  status  - показать статус"
        echo "  logs    - показать логи"
        exit 1
        ;;
esac
