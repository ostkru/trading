# Systemd сервисы
.PHONY: services-install services-start services-stop services-restart services-status services-uninstall

# Установка systemd сервисов
services-install:
	@echo "🔧 Установка systemd сервисов..."
	@sudo ./scripts/manage-services.sh install

# Запуск сервисов
services-start:
	@echo "🚀 Запуск сервисов..."
	@sudo ./scripts/manage-services.sh start

# Остановка сервисов
services-stop:
	@echo "⏹️ Остановка сервисов..."
	@sudo ./scripts/manage-services.sh stop

# Перезапуск сервисов
services-restart:
	@echo "🔄 Перезапуск сервисов..."
	@sudo ./scripts/manage-services.sh restart

# Статус сервисов
services-status:
	@echo "📊 Статус сервисов..."
	@sudo ./scripts/manage-services.sh status

# Удаление сервисов
services-uninstall:
	@echo "🗑️ Удаление сервисов..."
	@sudo ./scripts/manage-services.sh uninstall
