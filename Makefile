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

# OpenAPI Documentation Generator
.PHONY: openapi-generator
openapi-generator:
	@echo "🚀 Сборка генератора OpenAPI документации..."
	@go build -o openapi-generator ./cmd/openapi-generator
	@echo "✅ Генератор собран: ./openapi-generator"

.PHONY: generate-docs
generate-docs: openapi-generator
	@echo "📚 Генерация OpenAPI документации..."
	@./openapi-generator
	@echo "✅ Документация сгенерирована:"
	@echo "   - openapi_generated.json"
	@echo "   - api_documentation.html"

.PHONY: clean-docs
clean-docs:
	@echo "🧹 Очистка сгенерированной документации..."
	@rm -f openapi-generator openapi_generated.json api_documentation.html
	@echo "✅ Документация очищена"

# Тестирование
.PHONY: test test-api test-integration test-unit test-performance test-all
test:
	@echo "🧪 Запуск тестов..."
	@cd tests && make test-all

test-api:
	@echo "🧪 Запуск API тестов..."
	@cd tests && make test-api

test-integration:
	@echo "🧪 Запуск интеграционных тестов..."
	@cd tests && make test-integration

test-unit:
	@echo "🧪 Запуск модульных тестов..."
	@cd tests && make test-unit

test-performance:
	@echo "🧪 Запуск тестов производительности..."
	@cd tests && make test-performance

test-all:
	@echo "🧪 Запуск всех тестов..."
	@cd tests && make test-all
