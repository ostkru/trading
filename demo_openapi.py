#!/usr/bin/env python3
"""
Демонстрация автоматического генератора OpenAPI документации
"""

import json
import os

def generate_openapi_demo():
    """Генерирует демонстрационную OpenAPI спецификацию"""
    
    print("🚀 Генерация демонстрационной OpenAPI документации...")
    
    # Создаем OpenAPI спецификацию
    openapi_spec = {
        "openapi": "3.0.3",
        "info": {
            "title": "Trading API - Демо версия",
            "description": "API для торговой платформы с автоматически сгенерированными примерами",
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
                    "description": "Создает новый продукт в системе. Требует аутентификации и валидации данных.",
                    "operationId": "CreateProduct",
                    "requestBody": {
                        "description": "Данные для создания продукта",
                        "required": True,
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/CreateProductRequest"
                                },
                                "example": {
                                    "name": "iPhone 15 Pro",
                                    "brand": "Apple",
                                    "category": "Электроника",
                                    "description": "Смартфон премиум класса",
                                    "recommend_price": 99999.99,
                                    "vendor_article": "IP15PRO-256",
                                    "barcode": "1234567890123",
                                    "image_urls": ["https://example.com/iphone1.jpg"],
                                    "video_urls": ["https://example.com/iphone1.mp4"],
                                    "model_3d_urls": ["https://example.com/iphone1.obj"]
                                }
                            }
                        }
                    },
                    "responses": {
                        "201": {
                            "description": "Продукт успешно создан",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "success": {
                                                "type": "boolean",
                                                "example": True
                                            },
                                            "data": {
                                                "$ref": "#/components/schemas/Product"
                                            }
                                        }
                                    },
                                    "example": {
                                        "success": True,
                                        "data": {
                                            "id": 1,
                                            "name": "iPhone 15 Pro",
                                            "brand": "Apple",
                                            "category": "Электроника",
                                            "description": "Смартфон премиум класса",
                                            "recommend_price": 99999.99,
                                            "vendor_article": "IP15PRO-256",
                                            "barcode": "1234567890123",
                                            "image_urls": ["https://example.com/iphone1.jpg"],
                                            "video_urls": ["https://example.com/iphone1.mp4"],
                                            "model_3d_urls": ["https://example.com/iphone1.obj"],
                                            "user_id": 1,
                                            "created_at": "2025-01-01T00:00:00Z",
                                            "updated_at": "2025-01-01T00:00:00Z"
                                        }
                                    }
                                }
                            }
                        },
                        "400": {
                            "description": "Некорректный запрос",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "error": {
                                                "type": "string",
                                                "example": "Некорректный запрос"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "get": {
                    "tags": ["Products"],
                    "summary": "Получение списка продуктов",
                    "description": "Получает список всех продуктов с возможностью фильтрации. Поддерживает пагинацию и сортировку.",
                    "operationId": "ListProducts",
                    "responses": {
                        "200": {
                            "description": "Список продуктов",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "success": {
                                                "type": "boolean",
                                                "example": True
                                            },
                                            "data": {
                                                "type": "array",
                                                "items": {
                                                    "$ref": "#/components/schemas/Product"
                                                }
                                            }
                                        }
                                    },
                                    "example": {
                                        "success": True,
                                        "data": [
                                            {
                                                "id": 1,
                                                "name": "iPhone 15 Pro",
                                                "brand": "Apple",
                                                "category": "Электроника",
                                                "recommend_price": 99999.99
                                            },
                                            {
                                                "id": 2,
                                                "name": "MacBook Pro 16",
                                                "brand": "Apple",
                                                "category": "Электроника",
                                                "recommend_price": 299999.99
                                            }
                                        ]
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "/offers/public": {
                "get": {
                    "tags": ["Offers"],
                    "summary": "Публичные предложения",
                    "description": "Получает список публичных предложений. Доступно без аутентификации.",
                    "operationId": "ListPublicOffers",
                    "responses": {
                        "200": {
                            "description": "Список публичных предложений",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "success": {
                                                "type": "boolean",
                                                "example": True
                                            },
                                            "data": {
                                                "type": "array",
                                                "items": {
                                                    "$ref": "#/components/schemas/Offer"
                                                }
                                            }
                                        }
                                    },
                                    "example": {
                                        "success": True,
                                        "data": [
                                            {
                                                "id": 1,
                                                "product_id": 1,
                                                "type": "sale",
                                                "price": 99999.99,
                                                "lot_count": 5,
                                                "vat": True,
                                                "delivery_days": 3,
                                                "warehouse_id": 1
                                            }
                                        ]
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "/warehouses": {
                "post": {
                    "tags": ["Warehouses"],
                    "summary": "Создание склада",
                    "description": "Создает новый склад. Требует валидации географических координат.",
                    "operationId": "CreateWarehouse",
                    "requestBody": {
                        "description": "Данные для создания склада",
                        "required": True,
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/CreateWarehouseRequest"
                                },
                                "example": {
                                    "name": "Главный склад Москва",
                                    "address": "ул. Тверская, 1, Москва",
                                    "latitude": 55.7558,
                                    "longitude": 37.6176
                                }
                            }
                        }
                    },
                    "responses": {
                        "201": {
                            "description": "Склад успешно создан",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "type": "object",
                                        "properties": {
                                            "success": {
                                                "type": "boolean",
                                                "example": True
                                            },
                                            "data": {
                                                "$ref": "#/components/schemas/Warehouse"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "components": {
            "schemas": {
                "CreateProductRequest": {
                    "type": "object",
                    "properties": {
                        "name": {
                            "type": "string",
                            "description": "Название продукта",
                            "example": "iPhone 15 Pro"
                        },
                        "brand": {
                            "type": "string",
                            "description": "Бренд продукта",
                            "example": "Apple"
                        },
                        "category": {
                            "type": "string",
                            "description": "Категория продукта",
                            "example": "Электроника"
                        },
                        "description": {
                            "type": "string",
                            "description": "Описание продукта",
                            "example": "Смартфон премиум класса"
                        },
                        "recommend_price": {
                            "type": "number",
                            "format": "double",
                            "description": "Рекомендуемая цена",
                            "minimum": 0,
                            "example": 99999.99
                        },
                        "vendor_article": {
                            "type": "string",
                            "description": "Артикул поставщика",
                            "example": "IP15PRO-256"
                        },
                        "barcode": {
                            "type": "string",
                            "description": "Штрих-код продукта",
                            "example": "1234567890123"
                        },
                        "image_urls": {
                            "type": "array",
                            "description": "URL изображений продукта",
                            "items": {
                                "type": "string"
                            },
                            "example": ["https://example.com/iphone1.jpg"]
                        },
                        "video_urls": {
                            "type": "array",
                            "description": "URL видео продукта",
                            "items": {
                                "type": "string"
                            },
                            "example": ["https://example.com/iphone1.mp4"]
                        },
                        "model_3d_urls": {
                            "type": "array",
                            "description": "URL 3D моделей продукта",
                            "items": {
                                "type": "string"
                            },
                            "example": ["https://example.com/iphone1.obj"]
                        }
                    },
                    "required": ["name", "brand", "category", "recommend_price"]
                },
                "Product": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "integer",
                            "format": "int64",
                            "example": 1
                        },
                        "name": {
                            "type": "string",
                            "description": "Название продукта",
                            "example": "iPhone 15 Pro"
                        },
                        "brand": {
                            "type": "string",
                            "description": "Бренд продукта",
                            "example": "Apple"
                        },
                        "category": {
                            "type": "string",
                            "description": "Категория продукта",
                            "example": "Электроника"
                        },
                        "description": {
                            "type": "string",
                            "description": "Описание продукта",
                            "example": "Смартфон премиум класса"
                        },
                        "recommend_price": {
                            "type": "number",
                            "format": "double",
                            "description": "Рекомендуемая цена",
                            "example": 99999.99
                        },
                        "user_id": {
                            "type": "integer",
                            "format": "int64",
                            "example": 1
                        },
                        "created_at": {
                            "type": "string",
                            "format": "date-time",
                            "description": "Дата создания",
                            "example": "2025-01-01T00:00:00Z"
                        },
                        "updated_at": {
                            "type": "string",
                            "format": "date-time",
                            "description": "Дата обновления",
                            "example": "2025-01-01T00:00:00Z"
                        }
                    }
                },
                "CreateWarehouseRequest": {
                    "type": "object",
                    "properties": {
                        "name": {
                            "type": "string",
                            "description": "Название склада",
                            "example": "Главный склад Москва"
                        },
                        "address": {
                            "type": "string",
                            "description": "Адрес склада",
                            "example": "ул. Тверская, 1, Москва"
                        },
                        "latitude": {
                            "type": "number",
                            "format": "double",
                            "description": "Широта",
                            "minimum": -90,
                            "maximum": 90,
                            "example": 55.7558
                        },
                        "longitude": {
                            "type": "number",
                            "format": "double",
                            "description": "Долгота",
                            "minimum": -180,
                            "maximum": 180,
                            "example": 37.6176
                        }
                    },
                    "required": ["name", "address", "latitude", "longitude"]
                },
                "Offer": {
                    "type": "object",
                    "properties": {
                        "id": {
                            "type": "integer",
                            "format": "int64",
                            "example": 1
                        },
                        "product_id": {
                            "type": "integer",
                            "format": "int64",
                            "description": "ID продукта",
                            "example": 1
                        },
                        "type": {
                            "type": "string",
                            "description": "Тип предложения",
                            "enum": ["sale", "buy"],
                            "example": "sale"
                        },
                        "price": {
                            "type": "number",
                            "format": "double",
                            "description": "Цена предложения",
                            "example": 99999.99
                        },
                        "lot_count": {
                            "type": "integer",
                            "description": "Количество лотов",
                            "minimum": 1,
                            "example": 5
                        },
                        "vat": {
                            "type": "boolean",
                            "description": "Включен ли НДС",
                            "example": True
                        },
                        "delivery_days": {
                            "type": "integer",
                            "description": "Дни доставки",
                            "minimum": 1,
                            "maximum": 365,
                            "example": 3
                        },
                        "warehouse_id": {
                            "type": "integer",
                            "format": "int64",
                            "description": "ID склада",
                            "example": 1
                        }
                    }
                }
            },
            "securitySchemes": {
                "ApiKeyAuth": {
                    "type": "apiKey",
                    "description": "API ключ для аутентификации",
                    "name": "X-API-KEY",
                    "in": "header"
                }
            }
        },
        "tags": [
            {"name": "Products", "description": "Управление продуктами"},
            {"name": "Offers", "description": "Управление предложениями"},
            {"name": "Warehouses", "description": "Управление складами"},
            {"name": "Orders", "description": "Управление заказами"}
        ]
    }
    
    # Сохраняем в файл
    with open("openapi_demo.json", "w", encoding="utf-8") as f:
        json.dump(openapi_spec, f, indent=2, ensure_ascii=False)
    
    print("✅ Демонстрационная OpenAPI спецификация сохранена в openapi_demo.json")
    
    # Создаем HTML документацию
    html_content = f"""
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{openapi_spec['info']['title']} - API Документация</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 0; padding: 20px; }}
        .container {{ max-width: 1200px; margin: 0 auto; }}
        .header {{ background: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }}
        .endpoint {{ border: 1px solid #ddd; margin: 10px 0; border-radius: 5px; }}
        .method {{ padding: 10px; font-weight: bold; color: white; }}
        .get {{ background: #61affe; }}
        .post {{ background: #49cc90; }}
        .put {{ background: #fca130; }}
        .delete {{ background: #f93e3e; }}
        .path {{ padding: 10px; background: #f8f9fa; font-family: monospace; }}
        .description {{ padding: 10px; }}
        .examples {{ padding: 10px; background: #f8f9fa; }}
        .example {{ margin: 10px 0; }}
        .code {{ background: #2d3748; color: #e2e8f0; padding: 10px; border-radius: 3px; font-family: monospace; }}
        .schemas {{ margin-top: 30px; }}
        .schema {{ border: 1px solid #ddd; margin: 10px 0; border-radius: 5px; }}
        .schema-header {{ padding: 10px; background: #f8f9fa; font-weight: bold; }}
        .schema-body {{ padding: 10px; }}
        .field {{ margin: 5px 0; padding: 5px; background: #f8f9fa; border-radius: 3px; }}
        .type {{ color: #0066cc; font-weight: bold; }}
        .required {{ color: #cc0000; font-weight: bold; }}
        .example {{ color: #666; font-style: italic; }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{openapi_spec['info']['title']}</h1>
            <p>{openapi_spec['info']['description']}</p>
            <p><strong>Версия:</strong> {openapi_spec['info']['version']}</p>
        </div>

        <div class="info">
            <h2>📡 Серверы</h2>
            <ul>
"""
    
    for server in openapi_spec['servers']:
        html_content += f'                <li><strong>{server["url"]}</strong> - {server["description"]}</li>\n'
    
    html_content += """            </ul>
        </div>

        <h2>🚀 API Endpoints</h2>
"""
    
    for path, path_item in openapi_spec['paths'].items():
        html_content += f'        <div class="endpoint">\n'
        
        if 'post' in path_item:
            html_content += f'            <div class="method post">POST</div>\n'
            html_content += f'            <div class="path">{path}</div>\n'
            html_content += f'            <div class="description">\n'
            html_content += f'                <strong>{path_item["post"]["summary"]}</strong><br>\n'
            html_content += f'                {path_item["post"]["description"]}\n'
            html_content += f'            </div>\n'
            
            if 'requestBody' in path_item['post']:
                html_content += f'            <div class="examples">\n'
                html_content += f'                <strong>Пример запроса:</strong>\n'
                html_content += f'                <div class="example">\n'
                html_content += f'                    <div class="code">{json.dumps(path_item["post"]["requestBody"]["content"]["application/json"]["example"], indent=2, ensure_ascii=False)}</div>\n'
                html_content += f'                </div>\n'
                html_content += f'            </div>\n'
        
        if 'get' in path_item:
            html_content += f'            <div class="method get">GET</div>\n'
            html_content += f'            <div class="path">{path}</div>\n'
            html_content += f'            <div class="description">\n'
            html_content += f'                <strong>{path_item["get"]["summary"]}</strong><br>\n'
            html_content += f'                {path_item["get"]["description"]}\n'
            html_content += f'            </div>\n'
        
        html_content += '        </div>\n'
    
    html_content += """
        <div class="schemas">
            <h2>📋 Схемы данных</h2>
"""
    
    for name, schema in openapi_spec['components']['schemas'].items():
        html_content += f'            <div class="schema">\n'
        html_content += f'                <div class="schema-header">{name}</div>\n'
        html_content += f'                <div class="schema-body">\n'
        html_content += f'                    <strong>Тип:</strong> {schema["type"]}<br>\n'
        
        if 'properties' in schema:
            html_content += '                    <strong>Поля:</strong>\n'
            for field_name, field in schema['properties'].items():
                required = ""
                if 'required' in schema and field_name in schema['required']:
                    required = ' <span class="required">(обязательное)</span>'
                
                example = ""
                if 'example' in field:
                    example = f' <span class="example">пример: {field["example"]}</span>'
                
                html_content += f'                    <div class="field">\n'
                html_content += f'                        <strong>{field_name}</strong> <span class="type">({field["type"]})</span>{required}{example}\n'
                html_content += f'                    </div>\n'
        
        html_content += '                </div>\n'
        html_content += '            </div>\n'
    
    html_content += """
        </div>
    </div>
</body>
</html>
"""
    
    with open("api_documentation_demo.html", "w", encoding="utf-8") as f:
        f.write(html_content)
    
    print("✅ HTML документация сохранена в api_documentation_demo.html")
    print("🎉 Генерация завершена успешно!")

if __name__ == "__main__":
    generate_openapi_demo()
