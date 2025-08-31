#!/bin/bash

# 🧪 Скрипт запуска всех тестов Trading API

echo "🧪 Запуск всех тестов Trading API..."
echo "=================================="

# Проверяем, что мы в папке tests
if [ ! -f "README.md" ]; then
    echo "❌ Запускайте скрипт из папки tests/"
    exit 1
fi

# Создаем папку для результатов
mkdir -p results
timestamp=$(date +"%Y%m%d_%H%M%S")

echo "📅 Время запуска: $(date)"
echo ""

# 1. API тесты
echo "🚀 Запуск API тестов..."
cd api
if [ -f "test_redis_rate_limiting.php" ]; then
    echo "  - Redis Rate Limiting тесты..."
    php test_redis_rate_limiting.php > "../results/api_redis_${timestamp}.log" 2>&1
    if [ $? -eq 0 ]; then
        echo "  ✅ Redis тесты завершены"
    else
        echo "  ❌ Redis тесты завершились с ошибками"
    fi
else
    echo "  ⚠️ Redis тесты не найдены"
fi
cd ..

echo ""

# 2. Интеграционные тесты
echo "🔗 Запуск интеграционных тестов..."
cd integration
if [ -f "comprehensive_improved.php" ]; then
    echo "  - Комплексные тесты..."
    php comprehensive_improved.php > "../results/integration_${timestamp}.log" 2>&1
    if [ $? -eq 0 ]; then
        echo "  ✅ Интеграционные тесты завершены"
    else
        echo "  ❌ Интеграционные тесты завершились с ошибками"
    fi
else
    echo "  ⚠️ Комплексные тесты не найдены"
fi
cd ..

echo ""

# 3. Модульные тесты (если есть)
echo "🧩 Проверка модульных тестов..."
cd unit
if [ "$(ls -A)" ]; then
    echo "  - Модульные тесты найдены"
    for test_file in *.php; do
        if [ -f "$test_file" ]; then
            echo "    - Запуск $test_file..."
            php "$test_file" > "../results/unit_${test_file%.php}_${timestamp}.log" 2>&1
        fi
    done
else
    echo "  ℹ️ Модульные тесты не найдены"
fi
cd ..

echo ""

# 4. Тесты производительности (если есть)
echo "⚡ Проверка тестов производительности..."
cd performance
if [ "$(ls -A)" ]; then
    echo "  - Тесты производительности найдены"
    for test_file in *.php; do
        if [ -f "$test_file" ]; then
            echo "    - Запуск $test_file..."
            php "$test_file" > "../results/performance_${test_file%.php}_${timestamp}.log" 2>&1
        fi
    done
else
    echo "  ℹ️ Тесты производительности не найдены"
fi
cd ..

echo ""
echo "=================================="
echo "📊 Результаты тестов сохранены в папку results/"
echo "📅 Время завершения: $(date)"
echo "🎉 Все тесты завершены!"
