# 🚀 Автоматический генератор OpenAPI документации

Автоматический инструмент для генерации OpenAPI спецификации на основе Go приложения с примерами запросов и ответов.

## ✨ Возможности

- **Автоматический анализ Go кода** - сканирует структуры и обработчики
- **Генерация OpenAPI 3.0.3** - полная совместимость со стандартом
- **Автоматические примеры** - создает примеры запросов и ответов
- **HTML документация** - генерирует красивую веб-страницу
- **Валидация тегов** - анализирует binding и json теги
- **Безопасность** - автоматически добавляет схемы аутентификации

## 🏗️ Архитектура

```
cmd/openapi-generator/
├── main.go                 # Основной генератор
├── internal/               # Внутренняя логика
│   ├── parser/            # Парсер Go AST
│   ├── generator/         # Генератор OpenAPI
│   └── examples/          # Генератор примеров
└── templates/             # HTML шаблоны
```

## 🚀 Быстрый старт

### 1. Сборка генератора

```bash
make openapi-generator
```

### 2. Генерация документации

```bash
make generate-docs
```

### 3. Очистка

```bash
make clean-docs
```

## 📝 Использование

### Автоматический режим

```bash
# Генерация всех файлов
./openapi-generator
```

### Ручной режим

```bash
# Только OpenAPI JSON
./openapi-generator --format json

# Только HTML
./openapi-generator --format html

# Оба формата
./openapi-generator --format all
```

## 🔍 Как это работает

### 1. Сканирование Go файлов

Генератор автоматически сканирует директорию `internal/modules/` и анализирует:

- **Структуры** - для генерации схем данных
- **Обработчики** - для генерации API endpoints
- **Комментарии** - для описаний и метаданных
- **Теги** - для валидации и требований

### 2. Анализ AST

```go
// Пример комментария для автоматического определения HTTP метода и пути
// @POST /products
func (h *Handler) CreateProduct(c *gin.Context) {
    // ...
}
```

### 3. Генерация схем

```go
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required"`
    Article     string  `json:"article" binding:"required"`
    Brand       string  `json:"brand" binding:"required"`
    Category    string  `json:"category" binding:"required"`
    Price       float64 `json:"price" binding:"required,min=0"`
    Description string  `json:"description"`
}
```

Автоматически генерируется:

```json
{
  "CreateProductRequest": {
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "description": "Название продукта"
      },
      "price": {
        "type": "number",
        "format": "double",
        "minimum": 0,
        "example": 123.45
      }
    },
    "required": ["name", "article", "brand", "category", "price"]
  }
}
```

### 4. Генерация примеров

Автоматически создаются примеры на основе:

- **Типов данных** - string, int, float64, bool
- **Названий полей** - name, email, url
- **Валидации** - min, max, email, url
- **Структуры** - вложенные объекты

## 📊 Выходные файлы

### 1. `openapi_generated.json`

Полная OpenAPI 3.0.3 спецификация:

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Trading API",
    "description": "API для торговой платформы",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.portaldata.ru/v1/trading",
      "description": "Production server"
    },
    {
      "url": "http://localhost:8095",
      "description": "Local development server"
    }
  ],
  "paths": {
    "/products": {
      "post": {
        "tags": ["Products"],
        "summary": "Создание продукта",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/CreateProductRequest"
              },
              "example": {
                "name": "Пример продукта",
                "price": 999.99
              }
            }
          }
        }
      }
    }
  }
}
```

### 2. `api_documentation.html`

Красивая HTML документация с:

- **Цветовая кодировка** HTTP методов
- **Примеры запросов** и ответов
- **Схемы данных** с описаниями
- **Адаптивный дизайн** для мобильных устройств

## ⚙️ Конфигурация

### Настройка серверов

```go
Servers: []Server{
    {
        URL:         "https://api.portaldata.ru/v1/trading",
        Description: "Production server",
    },
    {
        URL:         "http://localhost:8095",
        Description: "Local development server",
    },
}
```

### Настройка тегов

```go
Tags: []Tag{
    {Name: "Products", Description: "Управление продуктами"},
    {Name: "Offers", Description: "Управление предложениями"},
    {Name: "Orders", Description: "Управление заказами"},
    {Name: "Warehouses", Description: "Управление складами"},
    {Name: "Rate Limiting", Description: "Управление лимитами API"},
}
```

### Настройка безопасности

```go
SecuritySchemes: map[string]*SecurityScheme{
    "ApiKeyAuth": &SecurityScheme{
        Type:        "apiKey",
        Description: "API ключ для аутентификации",
        Name:        "X-API-KEY",
        In:          "header",
    },
}
```

## 🔧 Расширение функциональности

### Добавление новых типов данных

```go
func (g *OpenAPIGenerator) getBasicType(goType string) string {
    switch goType {
    case "string":
        return "string"
    case "int", "int64":
        return "integer"
    case "float64":
        return "number"
    case "bool":
        return "boolean"
    case "time.Time":
        return "string"
    case "uuid.UUID":
        return "string"
    case "decimal.Decimal":
        return "number"
    default:
        return "string"
    }
}
```

### Добавление новых валидаций

```go
func (g *OpenAPIGenerator) parseValidation(tag *ast.BasicLit) []string {
    if tag == nil {
        return nil
    }
    
    tagStr := strings.Trim(tag.Value, "`")
    validations := make([]string, 0)
    
    // Существующие валидации
    if strings.Contains(tagStr, "min:") {
        validations = append(validations, "min")
    }
    
    // Новые валидации
    if strings.Contains(tagStr, "regexp:") {
        validations = append(validations, "regexp")
    }
    if strings.Contains(tagStr, "unique") {
        validations = append(validations, "unique")
    }
    
    return validations
}
```

### Кастомные примеры

```go
func (g *OpenAPIGenerator) generateFieldExample(field GoField) interface{} {
    switch field.Type {
    case "string":
        if strings.Contains(strings.ToLower(field.Name), "phone") {
            return "+7 (999) 123-45-67"
        }
        if strings.Contains(strings.ToLower(field.Name), "inn") {
            return "1234567890"
        }
        // ... остальная логика
    }
}
```

## 📚 Интеграция с CI/CD

### GitHub Actions

```yaml
name: Generate OpenAPI Docs
on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Generate OpenAPI documentation
        run: |
          make generate-docs
      
      - name: Commit generated docs
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add openapi_generated.json api_documentation.html
          git commit -m "Auto-generate OpenAPI documentation" || exit 0
          git push
```

### GitLab CI

```yaml
generate-docs:
  stage: build
  image: golang:1.21
  script:
    - make generate-docs
  artifacts:
    paths:
      - openapi_generated.json
      - api_documentation.html
  only:
    - main
    - develop
```

## 🎯 Лучшие практики

### 1. Комментарии для обработчиков

```go
// @POST /products
// Создает новый продукт в системе
// Требует аутентификации и валидации данных
func (h *Handler) CreateProduct(c *gin.Context) {
    // ...
}
```

### 2. Описательные названия структур

```go
// Хорошо
type CreateProductRequest struct { ... }
type ProductResponse struct { ... }

// Плохо
type Request struct { ... }
type Response struct { ... }
```

### 3. Валидационные теги

```go
type Product struct {
    Name        string  `json:"name" binding:"required,min=1,max=100"`
    Price       float64 `json:"price" binding:"required,min=0"`
    Category    string  `json:"category" binding:"required,oneof=electronics clothing books"`
    Email       string  `json:"email" binding:"required,email"`
}
```

### 4. Группировка по тегам

```go
// @POST /products
// @tag Products
func (h *Handler) CreateProduct(c *gin.Context) { ... }

// @GET /products
// @tag Products
func (h *Handler) ListProducts(c *gin.Context) { ... }
```

## 🐛 Устранение неполадок

### Ошибка: "не удалось найти структуру"

```bash
# Проверьте, что структура определена в internal/modules/
find internal/modules -name "*.go" -exec grep -l "type.*struct" {} \;
```

### Ошибка: "не удалось определить HTTP метод"

```bash
# Убедитесь, что в комментариях есть @POST, @GET и т.д.
grep -r "@POST\|@GET\|@PUT\|@DELETE" internal/modules/
```

### Ошибка: "не удалось сгенерировать примеры"

```bash
# Проверьте типы полей в структурах
grep -r "type.*struct" internal/modules/ | head -5
```

## 🔮 Планы развития

- [ ] **Поддержка GraphQL** - генерация GraphQL схем
- [ ] **Интерактивная документация** - Swagger UI интеграция
- [ ] **Автоматические тесты** - генерация тестов на основе схем
- [ ] **Валидация** - проверка соответствия кода и документации
- [ ] **Многоязычность** - поддержка разных языков описаний
- [ ] **Плагины** - система расширений для кастомной логики

## 📞 Поддержка

Если у вас есть вопросы или предложения:

1. **Создайте Issue** в репозитории
2. **Напишите в чат** команды разработки
3. **Проверьте документацию** по ссылкам выше

---

**🎉 Автоматизируйте документацию и сосредоточьтесь на коде!**
