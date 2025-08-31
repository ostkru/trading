package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OpenAPI структуры
type OpenAPISpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       Info                   `json:"info"`
	Servers    []Server               `json:"servers"`
	Paths      map[string]PathItem    `json:"paths"`
	Components Components             `json:"components"`
	Tags       []Tag                  `json:"tags"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Summary string     `json:"summary,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	OperationID string                `json:"operationId"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Schema      *SchemaRef  `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

type RequestBody struct {
	Description string                  `json:"description"`
	Content     map[string]MediaType   `json:"content"`
	Required    bool                    `json:"required"`
}

type MediaType struct {
	Schema  *SchemaRef    `json:"schema"`
	Example interface{}   `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

type Response struct {
	Description string                  `json:"description"`
	Content     map[string]MediaType   `json:"content,omitempty"`
	Headers     map[string]Header      `json:"headers,omitempty"`
}

type Header struct {
	Description string     `json:"description"`
	Schema      *SchemaRef `json:"schema"`
}

type SchemaRef struct {
	Ref         string                 `json:"$ref,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]*SchemaRef  `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Items       *SchemaRef             `json:"items,omitempty"`
	Example     interface{}            `json:"example,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
}

type Example struct {
	Summary       string      `json:"summary"`
	Description   string      `json:"description"`
	Value         interface{} `json:"value"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

type Components struct {
	Schemas         map[string]*SchemaRef         `json:"schemas"`
	Parameters      map[string]*Parameter         `json:"parameters"`
	RequestBodies   map[string]*RequestBody      `json:"requestBodies"`
	Responses       map[string]*Response         `json:"responses"`
	SecuritySchemes map[string]*SecurityScheme   `json:"securitySchemes"`
}

type SecurityScheme struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Name        string `json:"name,omitempty"`
	In          string `json:"in,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Go структуры для анализа
type GoStruct struct {
	Name       string
	Fields     []GoField
	Comment    string
	Package    string
}

type GoField struct {
	Name       string
	Type       string
	Tag        string
	Comment    string
	Required   bool
	Validation []string
}

type GoHandler struct {
	Name       string
	Method     string
	Path       string
	Comment    string
	Request    *GoStruct
	Response   *GoStruct
	Parameters []GoField
}

// Генератор OpenAPI
type OpenAPIGenerator struct {
	Spec        *OpenAPISpec
	GoStructs   map[string]*GoStruct
	GoHandlers  []*GoHandler
	Examples    map[string]interface{}
}

func NewOpenAPIGenerator() *OpenAPIGenerator {
	return &OpenAPIGenerator{
		Spec: &OpenAPISpec{
			OpenAPI: "3.0.3",
			Info: Info{
				Title:       "Trading API",
				Description: "API для торговой платформы",
				Version:     "1.0.0",
			},
			Servers: []Server{
				{
					URL:         "https://api.portaldata.ru/v1/trading",
					Description: "Production server",
				},
				{
					URL:         "http://localhost:8095",
					Description: "Local development server",
				},
			},
			Paths:      make(map[string]PathItem),
			Components: Components{
				Schemas:         make(map[string]*SchemaRef),
				Parameters:      make(map[string]*Parameter),
				RequestBodies:   make(map[string]*RequestBody),
				Responses:       make(map[string]*Response),
				SecuritySchemes: make(map[string]*SecurityScheme),
			},
			Tags: []Tag{
				{Name: "Products", Description: "Управление продуктами"},
				{Name: "Offers", Description: "Управление предложениями"},
				{Name: "Orders", Description: "Управление заказами"},
				{Name: "Warehouses", Description: "Управление складами"},
				{Name: "Rate Limiting", Description: "Управление лимитами API"},
			},
		},
		GoStructs:  make(map[string]*GoStruct),
		GoHandlers: make([]*GoHandler, 0),
		Examples:   make(map[string]interface{}),
	}
}

func (g *OpenAPIGenerator) Generate() error {
	// 1. Сканируем Go файлы
	if err := g.scanGoFiles(); err != nil {
		return fmt.Errorf("ошибка сканирования Go файлов: %w", err)
	}

	// 2. Генерируем схемы из Go структур
	if err := g.generateSchemas(); err != nil {
		return fmt.Errorf("ошибка генерации схем: %w", err)
	}

	// 3. Генерируем пути из Go обработчиков
	if err := g.generatePaths(); err != nil {
		return fmt.Errorf("ошибка генерации путей: %w", err)
	}

	// 4. Генерируем примеры
	if err := g.generateExamples(); err != nil {
		return fmt.Errorf("ошибка генерации примеров: %w", err)
	}

	// 5. Добавляем стандартные компоненты
	g.addStandardComponents()

	return nil
}

func (g *OpenAPIGenerator) scanGoFiles() error {
	// Сканируем internal/modules для поиска структур и обработчиков
	modulesDir := "internal/modules"
	
	return filepath.Walk(modulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if err := g.parseGoFile(path); err != nil {
				log.Printf("Предупреждение: ошибка парсинга %s: %v", path, err)
			}
		}
		
		return nil
	})
}

func (g *OpenAPIGenerator) parseGoFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

			// Анализируем AST
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.TypeSpec:
				if structType, ok := x.Type.(*ast.StructType); ok {
					g.parseStruct(x.Name.Name, node.Name.Name, structType)
				}
			case *ast.FuncDecl:
				if x.Recv != nil {
					g.parseHandler(x)
				}
			}
			return true
		})

	return nil
}

func (g *OpenAPIGenerator) parseStruct(name, packageName string, structType *ast.StructType) {
	goStruct := &GoStruct{
		Name:    name,
		Package: packageName,
		Fields:  make([]GoField, 0),
	}

	for _, field := range structType.Fields.List {
		if field.Names != nil {
			goField := GoField{
				Name:     field.Names[0].Name,
				Type:     g.getTypeString(field.Type),
				Tag:      g.getTagString(field.Tag),
				Comment:  g.getCommentString(field.Doc),
				Required: g.isFieldRequired(field.Tag),
			}
			
			// Парсим валидацию из тегов
			goField.Validation = g.parseValidation(field.Tag)
			
			goStruct.Fields = append(goStruct.Fields, goField)
		}
	}

	g.GoStructs[name] = goStruct
}

func (g *OpenAPIGenerator) parseHandler(funcDecl *ast.FuncDecl) {
	// Анализируем комментарии для поиска HTTP методов и путей
	comment := g.getCommentString(funcDecl.Doc)
	
	// Ищем HTTP метод и путь в комментариях
	method, path := g.extractHTTPInfo(comment)
	if method == "" || path == "" {
		return
	}

	handler := &GoHandler{
		Name:    funcDecl.Name.Name,
		Method:  method,
		Path:    path,
		Comment: comment,
	}

	// Анализируем параметры функции
	g.analyzeHandlerParams(funcDecl, handler)
	
	g.GoHandlers = append(g.GoHandlers, handler)
}

func (g *OpenAPIGenerator) extractHTTPInfo(comment string) (method, path string) {
	// Ищем паттерны типа @GET /products или @POST /offers
	re := regexp.MustCompile(`@(GET|POST|PUT|DELETE|PATCH)\s+([^\s]+)`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

func (g *OpenAPIGenerator) analyzeHandlerParams(funcDecl *ast.FuncDecl, handler *GoHandler) {
	// Анализируем параметры функции для определения request/response структур
	for _, param := range funcDecl.Type.Params.List {
		if len(param.Names) > 0 {
			paramName := param.Names[0].Name
			paramType := g.getTypeString(param.Type)
			
			// Определяем тип параметра
			switch {
			case strings.Contains(paramName, "c") || strings.Contains(paramType, "Context"):
				// Пропускаем контекст
			case strings.Contains(paramType, "Request"):
				handler.Request = g.findStructByName(paramType)
			case strings.Contains(paramType, "Response"):
				handler.Response = g.findStructByName(paramType)
			default:
				// Другие параметры
				handler.Parameters = append(handler.Parameters, GoField{
					Name: paramName,
					Type: paramType,
				})
			}
		}
	}
}

func (g *OpenAPIGenerator) getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + g.getTypeString(t.X)
	case *ast.ArrayType:
		return "[]" + g.getTypeString(t.Elt)
	case *ast.SelectorExpr:
		return g.getTypeString(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func (g *OpenAPIGenerator) getTagString(tag *ast.BasicLit) string {
	if tag != nil {
		return strings.Trim(tag.Value, "`")
	}
	return ""
}

func (g *OpenAPIGenerator) getCommentString(doc *ast.CommentGroup) string {
	if doc != nil {
		comments := make([]string, 0)
		for _, comment := range doc.List {
			comments = append(comments, strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")))
		}
		return strings.Join(comments, " ")
	}
	return ""
}

func (g *OpenAPIGenerator) isFieldRequired(tag *ast.BasicLit) bool {
	if tag == nil {
		return false
	}
	
	tagStr := strings.Trim(tag.Value, "`")
	return strings.Contains(tagStr, `binding:"required"`) || 
		   strings.Contains(tagStr, `json:",required"`)
}

func (g *OpenAPIGenerator) parseValidation(tag *ast.BasicLit) []string {
	if tag == nil {
		return nil
	}
	
	tagStr := strings.Trim(tag.Value, "`")
	validations := make([]string, 0)
	
	// Ищем validation теги
	if strings.Contains(tagStr, "min:") {
		validations = append(validations, "min")
	}
	if strings.Contains(tagStr, "max:") {
		validations = append(validations, "max")
	}
	if strings.Contains(tagStr, "email") {
		validations = append(validations, "email")
	}
	if strings.Contains(tagStr, "url") {
		validations = append(validations, "url")
	}
	
	return validations
}

func (g *OpenAPIGenerator) findStructByName(name string) *GoStruct {
	// Убираем указатели и пакеты
	cleanName := strings.TrimPrefix(name, "*")
	if idx := strings.LastIndex(cleanName, "."); idx != -1 {
		cleanName = cleanName[idx+1:]
	}
	
	return g.GoStructs[cleanName]
}

func (g *OpenAPIGenerator) generateSchemas() error {
	for name, goStruct := range g.GoStructs {
		schema := g.convertGoStructToSchema(goStruct)
		g.Spec.Components.Schemas[name] = schema
	}
	return nil
}

func (g *OpenAPIGenerator) convertGoStructToSchema(goStruct *GoStruct) *SchemaRef {
	schema := &SchemaRef{
		Type:        "object",
		Description: goStruct.Comment,
		Properties:  make(map[string]*SchemaRef),
		Required:    make([]string, 0),
	}

	for _, field := range goStruct.Fields {
		fieldSchema := g.convertGoFieldToSchema(field)
		schema.Properties[field.Name] = fieldSchema
		
		if field.Required {
			schema.Required = append(schema.Required, field.Name)
		}
	}

	return schema
}

func (g *OpenAPIGenerator) convertGoFieldToSchema(field GoField) *SchemaRef {
	schema := &SchemaRef{
		Description: field.Comment,
	}

	// Определяем тип поля
	switch {
	case strings.HasPrefix(field.Type, "[]"):
		schema.Type = "array"
		schema.Items = &SchemaRef{Type: g.getBasicType(strings.TrimPrefix(field.Type, "[]"))}
	case field.Type == "string":
		schema.Type = "string"
		if len(field.Validation) > 0 {
			schema.Enum = g.generateEnumValues(field)
		}
	case field.Type == "int" || field.Type == "int64":
		schema.Type = "integer"
		schema.Format = "int64"
	case field.Type == "float64":
		schema.Type = "number"
		schema.Format = "double"
	case field.Type == "bool":
		schema.Type = "boolean"
	case field.Type == "time.Time":
		schema.Type = "string"
		schema.Format = "date-time"
	default:
		// Ссылка на другую структуру
		schema.Ref = "#/components/schemas/" + field.Type
	}

	// Добавляем примеры
	schema.Example = g.generateFieldExample(field)

	return schema
}

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
	default:
		return "string"
	}
}

func (g *OpenAPIGenerator) generateEnumValues(field GoField) []interface{} {
	// Генерируем примеры enum значений на основе валидации
	examples := make([]interface{}, 0)
	
	for _, validation := range field.Validation {
		switch validation {
		case "min":
			examples = append(examples, "min_value")
		case "max":
			examples = append(examples, "max_value")
		case "email":
			examples = append(examples, "user@example.com")
		case "url":
			examples = append(examples, "https://example.com")
		}
	}
	
	if len(examples) == 0 {
		examples = append(examples, "example_value")
	}
	
	return examples
}

func (g *OpenAPIGenerator) generateFieldExample(field GoField) interface{} {
	switch field.Type {
	case "string":
		if strings.Contains(strings.ToLower(field.Name), "name") {
			return "Example Name"
		}
		if strings.Contains(strings.ToLower(field.Name), "email") {
			return "user@example.com"
		}
		if strings.Contains(strings.ToLower(field.Name), "url") {
			return "https://example.com"
		}
		return "example_string"
	case "int", "int64":
		return 123
	case "float64":
		return 123.45
	case "bool":
		return true
	case "time.Time":
		return "2025-01-01T00:00:00Z"
	default:
		return "example_value"
	}
}

func (g *OpenAPIGenerator) generatePaths() error {
	for _, handler := range g.GoHandlers {
		path := handler.Path
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		operation := g.convertHandlerToOperation(handler)
		
		pathItem := g.Spec.Paths[path]
		switch handler.Method {
		case "GET":
			pathItem.Get = operation
		case "POST":
			pathItem.Post = operation
		case "PUT":
			pathItem.Put = operation
		case "DELETE":
			pathItem.Delete = operation
		case "PATCH":
			pathItem.Patch = operation
		}
		
		pathItem.Summary = handler.Comment
		g.Spec.Paths[path] = pathItem
	}
	
	return nil
}

func (g *OpenAPIGenerator) convertHandlerToOperation(handler *GoHandler) *Operation {
	operation := &Operation{
		Tags:        []string{g.getTagFromPath(handler.Path)},
		Summary:     handler.Comment,
		Description: handler.Comment,
		OperationID: handler.Name,
		Responses:   g.generateResponses(handler),
	}

	// Добавляем параметры
	if len(handler.Parameters) > 0 {
		operation.Parameters = g.convertParameters(handler.Parameters)
	}

	// Добавляем request body для POST/PUT/PATCH
	if handler.Request != nil && (handler.Method == "POST" || handler.Method == "PUT" || handler.Method == "PATCH") {
		operation.RequestBody = g.generateRequestBody(handler.Request)
	}

	return operation
}

func (g *OpenAPIGenerator) getTagFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		// Преобразуем в заголовок
		return strings.Title(parts[0])
	}
	return "General"
}

func (g *OpenAPIGenerator) convertParameters(fields []GoField) []Parameter {
	parameters := make([]Parameter, 0)
	
	for _, field := range fields {
		param := Parameter{
			Name:        field.Name,
			In:          "query",
			Description: field.Comment,
			Required:    field.Required,
			Schema:      g.convertGoFieldToSchema(field),
			Example:     g.generateFieldExample(field),
		}
		parameters = append(parameters, param)
	}
	
	return parameters
}

func (g *OpenAPIGenerator) generateRequestBody(goStruct *GoStruct) *RequestBody {
	return &RequestBody{
		Description: fmt.Sprintf("Данные для %s", goStruct.Name),
		Required:    true,
		Content: map[string]MediaType{
			"application/json": {
				Schema: &SchemaRef{
					Ref: "#/components/schemas/" + goStruct.Name,
				},
				Example: g.generateStructExample(goStruct),
			},
		},
	}
}

func (g *OpenAPIGenerator) generateResponses(handler *GoHandler) map[string]Response {
	responses := make(map[string]Response)
	
	// Успешный ответ
	successResponse := Response{
		Description: "Успешный ответ",
		Content: map[string]MediaType{
			"application/json": {
				Schema: &SchemaRef{
					Type: "object",
					Properties: map[string]*SchemaRef{
						"success": {
							Type:    "boolean",
							Example: true,
						},
						"data": {
							Ref: "#/components/schemas/" + handler.Response.Name,
						},
					},
				},
				Example: g.generateSuccessResponseExample(handler),
			},
		},
	}
	
	responses["200"] = successResponse
	
	// Ошибки
	errorResponses := map[string]string{
		"400": "Некорректный запрос",
		"401": "Не авторизован",
		"403": "Доступ запрещен",
		"404": "Не найдено",
		"500": "Внутренняя ошибка сервера",
	}
	
	for code, description := range errorResponses {
		responses[code] = Response{
			Description: description,
			Content: map[string]MediaType{
				"application/json": {
					Schema: &SchemaRef{
						Type: "object",
						Properties: map[string]*SchemaRef{
							"error": {
								Type:    "string",
								Example: description,
							},
						},
					},
					Example: map[string]interface{}{
						"error": description,
					},
				},
			},
		}
	}
	
	return responses
}

func (g *OpenAPIGenerator) generateStructExample(goStruct *GoStruct) interface{} {
	example := make(map[string]interface{})
	
	for _, field := range goStruct.Fields {
		example[field.Name] = g.generateFieldExample(field)
	}
	
	return example
}

func (g *OpenAPIGenerator) generateSuccessResponseExample(handler *GoHandler) interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    g.generateStructExample(handler.Response),
	}
}

func (g *OpenAPIGenerator) generateExamples() error {
	// Генерируем примеры для основных структур
	g.Examples["Product"] = map[string]interface{}{
		"id":          1,
		"name":        "Пример продукта",
		"article":     "PROD-001",
		"brand":       "Примерный бренд",
		"category":    "Электроника",
		"price":       999.99,
		"description": "Описание продукта",
		"user_id":     1,
	}
	
	g.Examples["Offer"] = map[string]interface{}{
		"id":             1,
		"product_id":     1,
		"type":           "sale",
		"price":          999.99,
		"lot_count":      10,
		"vat":            true,
		"delivery_days":  3,
		"user_id":        1,
		"warehouse_id":   1,
	}
	
	g.Examples["Warehouse"] = map[string]interface{}{
		"id":        1,
		"name":      "Главный склад",
		"address":   "ул. Примерная, 1",
		"latitude":  55.7558,
		"longitude": 37.6176,
		"user_id":  1,
	}
	
	return nil
}

func (g *OpenAPIGenerator) addStandardComponents() {
	// Добавляем стандартные схемы безопасности
	g.Spec.Components.SecuritySchemes["ApiKeyAuth"] = &SecurityScheme{
		Type:        "apiKey",
		Description: "API ключ для аутентификации",
		Name:        "X-API-KEY",
		In:          "header",
	}
	
	// Добавляем стандартные параметры
	g.Spec.Components.Parameters["APIKey"] = &Parameter{
		Name:        "X-API-KEY",
		In:          "header",
		Description: "API ключ для аутентификации",
		Required:    true,
		Schema: &SchemaRef{
			Type: "string",
		},
		Example: "your_api_key_here",
	}
}

func (g *OpenAPIGenerator) SaveToFile(filename string) error {
	// Сортируем пути для стабильного вывода
	sortedPaths := make([]string, 0, len(g.Spec.Paths))
	for path := range g.Spec.Paths {
		sortedPaths = append(sortedPaths, path)
	}
	
	// Создаем временную копию с отсортированными путями
	tempSpec := *g.Spec
	tempSpec.Paths = make(map[string]PathItem)
	for _, path := range sortedPaths {
		tempSpec.Paths[path] = g.Spec.Paths[path]
	}
	
	// Сохраняем в файл
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	return encoder.Encode(tempSpec)
}

func (g *OpenAPIGenerator) GenerateHTMLDocs(filename string) error {
	// Создаем простую HTML документацию
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - API Документация</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .info { background: #e8f4fd; padding: 15px; border-radius: 5px; margin: 10px 0; }
        .schemas { margin-top: 30px; }
        .schema { border: 1px solid #ddd; margin: 10px 0; border-radius: 5px; }
        .schema-header { padding: 10px; background: #f8f9fa; font-weight: bold; }
        .schema-body { padding: 10px; }
        .field { margin: 5px 0; padding: 5px; background: #f8f9fa; border-radius: 3px; }
        .type { color: #0066cc; font-weight: bold; }
        .required { color: #cc0000; font-weight: bold; }
        .example { color: #666; font-style: italic; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
            <p>%s</p>
            <p><strong>Версия:</strong> %s</p>
        </div>

        <div class="info">
            <h2>📡 Серверы</h2>
            <ul>`, g.Spec.Info.Title, g.Spec.Info.Title, g.Spec.Info.Description, g.Spec.Info.Version)

	// Добавляем серверы
	for _, server := range g.Spec.Servers {
		html += fmt.Sprintf(`<li><strong>%s</strong> - %s</li>`, server.URL, server.Description)
	}
	html += `</ul></div>`

	// Добавляем схемы
	html += `<div class="schemas">
            <h2>📋 Схемы данных</h2>`

	for name, schema := range g.Spec.Components.Schemas {
		html += fmt.Sprintf(`
            <div class="schema">
                <div class="schema-header">%s</div>
                <div class="schema-body">`, name)

		if schema.Description != "" {
			html += fmt.Sprintf(`<p><strong>Описание:</strong> %s</p>`, schema.Description)
		}

		if schema.Properties != nil {
			html += `<p><strong>Поля:</strong></p>`
			for fieldName, field := range schema.Properties {
				required := ""
				if len(schema.Required) > 0 {
					for _, req := range schema.Required {
						if req == fieldName {
							required = " <span class='required'>(обязательное)</span>"
							break
						}
					}
				}

				example := ""
				if field.Example != nil {
					example = fmt.Sprintf(" <span class='example'>пример: %v</span>", field.Example)
				}

				html += fmt.Sprintf(`<div class="field">
                    <strong>%s</strong> <span class="type">(%s)</span>%s%s
                </div>`, fieldName, field.Type, required, example)
			}
		}

		html += `</div></div>`
	}

	html += `
        </div>
    </div>
</body>
</html>`

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(html)
	return err
}

func main() {
	log.Println("🚀 Запуск генератора OpenAPI документации...")
	
	generator := NewOpenAPIGenerator()
	
	if err := generator.Generate(); err != nil {
		log.Fatalf("❌ Ошибка генерации: %v", err)
	}
	
	// Сохраняем OpenAPI спецификацию
	if err := generator.SaveToFile("openapi_generated.json"); err != nil {
		log.Fatalf("❌ Ошибка сохранения OpenAPI файла: %v", err)
	}
	log.Println("✅ OpenAPI спецификация сохранена в openapi_generated.json")
	
	// Генерируем HTML документацию
	if err := generator.GenerateHTMLDocs("api_documentation.html"); err != nil {
		log.Fatalf("❌ Ошибка генерации HTML документации: %v", err)
	}
	log.Println("✅ HTML документация сохранена в api_documentation.html")
	
	log.Println("🎉 Генерация завершена успешно!")
}
