# 🧪 Makefile для тестов Trading API

.PHONY: test test-api test-integration test-unit test-performance test-all clean-results

# Запуск всех тестов
test: test-all

# Запуск API тестов
test-api:
	@echo "🚀 Запуск API тестов..."
	@cd api && php test_redis_rate_limiting.php

# Запуск интеграционных тестов
test-integration:
	@echo "🔗 Запуск интеграционных тестов..."
	@cd integration && php comprehensive_improved.php

# Запуск модульных тестов
test-unit:
	@echo "🧩 Запуск модульных тестов..."
	@if [ -d "unit" ] && [ "$(ls -A unit)" ]; then \
		cd unit && for file in *.php; do \
			if [ -f "$$file" ]; then \
				echo "  - $$file"; \
				php "$$file"; \
			fi; \
		done; \
	else \
		echo "ℹ️ Модульные тесты не найдены"; \
	fi

# Запуск тестов производительности
test-performance:
	@echo "⚡ Запуск тестов производительности..."
	@if [ -d "performance" ] && [ "$(ls -A performance)" ]; then \
		cd performance && for file in *.php; do \
			if [ -f "$$file" ]; then \
				echo "  - $$file"; \
				php "$$file"; \
			fi; \
		done; \
	else \
		echo "ℹ️ Тесты производительности не найдены"; \
	fi

# Запуск всех тестов через bash скрипт
test-all:
	@echo "🧪 Запуск всех тестов..."
	@./run_tests.sh

# Очистка результатов тестов
clean-results:
	@echo "🧹 Очистка результатов тестов..."
	@rm -rf results/*
	@echo "✅ Результаты очищены"

# Помощь
help:
	@echo "🧪 Доступные команды для тестов:"
	@echo ""
	@echo "  test-api        - Запуск API тестов"
	@echo "  test-integration - Запуск интеграционных тестов"
	@echo "  test-unit       - Запуск модульных тестов"
	@echo "  test-performance - Запуск тестов производительности"
	@echo "  test-all        - Запуск всех тестов"
	@echo "  clean-results   - Очистка результатов"
	@echo "  help            - Показать эту справку"



