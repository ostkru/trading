# Отчет об обновлении PUT методов для возврата только измененных данных

## Выполненные изменения

### 1. OpenAPI документация

#### Созданы новые схемы ответов для PUT методов:
- **`ProductUpdateResponse`** - для обновления продуктов
- **`OfferUpdateResponse`** - для обновления офферов  
- **`OrderStatusUpdateResponse`** - для изменения статуса заказов
- **`WarehouseUpdateResponse`** - для обновления складов

#### Обновлены все PUT методы в OpenAPI спецификации:
- **PUT `/products/{id}`** - теперь возвращает `ProductUpdateResponse`
- **PUT `/offers/{id}`** - теперь возвращает `OfferUpdateResponse`
- **PUT `/orders/{id}/status`** - теперь возвращает `OrderStatusUpdateResponse`
- **PUT `/warehouses/{id}`** - теперь возвращает `WarehouseUpdateResponse`

### 2. Go код

#### Модуль Offer (`/var/www/go/internal/modules/offer/`)

**Добавлены новые структуры в `model.go`:**
```go
type OfferUpdateResponse struct {
    OfferID       int64                  `json:"offer_id"`
    UpdatedFields map[string]interface{} `json:"updated_fields"`
    UpdatedAt     *string                `json:"updated_at"`
}
```

**Обновлен метод `UpdateOffer` в `service.go`:**
- Изменен возвращаемый тип с `*Offer` на `*OfferUpdateResponse`
- Добавлено отслеживание измененных полей в `updatedFields`
- Возвращается только ID оффера, измененные поля и время обновления

**Обновлен обработчик `UpdateOffer` в `handler.go`:**
- Изменена переменная с `offer` на `response`
- Обновлен возврат ответа

#### Модуль Order (`/var/www/go/internal/modules/order/`)

**Добавлена новая структура в `model.go`:**
```go
type OrderStatusUpdateResponse struct {
    OrderID         int64      `json:"order_id"`
    OrderStatus     string     `json:"order_status"`
    StatusReason    *string    `json:"status_reason,omitempty"`
    StatusChangedAt *time.Time `json:"status_changed_at,omitempty"`
    StatusChangedBy *int64     `json:"status_changed_by,omitempty"`
}
```

**Обновлен метод `UpdateOrderStatus` в `service.go`:**
- Изменен возвращаемый тип с `*Order` на `*OrderStatusUpdateResponse`
- Убрано получение полного объекта заказа
- Добавлено получение только данных статуса из БД
- Возвращается ID заказа, новый статус, причина изменения, время изменения и кто изменил

**Обновлен обработчик `UpdateOrderStatus` в `handler.go`:**
- Изменена переменная с `order` на `response`
- Обновлен возврат ответа

#### Модуль Warehouse (`/var/www/go/internal/modules/warehouse/`)

**Добавлена новая структура в `model.go`:**
```go
type WarehouseUpdateResponse struct {
    WarehouseID   int64                  `json:"warehouse_id"`
    UpdatedFields map[string]interface{} `json:"updated_fields"`
    UpdatedAt     *string                `json:"updated_at"`
}
```

**Обновлен метод `UpdateWarehouse` в `service.go`:**
- Изменен возвращаемый тип с `*Warehouse` на `*WarehouseUpdateResponse`
- Добавлено отслеживание измененных полей в `updatedFields`
- Изменена логика обновления для обработки только переданных полей
- Добавлен импорт `strings` для работы с SQL запросами
- Возвращается только ID склада, измененные поля и время обновления

**Обновлен обработчик `UpdateWarehouse` в `handler.go`:**
- Изменена переменная с `warehouse` на `response`
- Обновлен возврат ответа

## Структура новых ответов

### Для офферов и складов:
```json
{
  "offer_id": 123,
  "updated_fields": {
    "price_per_unit": 1500.0,
    "available_lots": 10
  },
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Для статуса заказов:
```json
{
  "order_id": 456,
  "order_status": "confirmed",
  "status_reason": "Заказ подтвержден продавцом",
  "status_changed_at": "2024-01-15T10:30:00Z",
  "status_changed_by": 789
}
```

## Преимущества изменений

1. **Эффективность**: Возвращаются только измененные данные, что уменьшает размер ответа
2. **Информативность**: Клиент получает четкую информацию о том, что именно было изменено
3. **Производительность**: Меньше данных передается по сети
4. **Прозрачность**: Легко отследить, какие поля были обновлены

## Файлы, которые были изменены

### OpenAPI документация:
- `/var/www/go/openapi.json` - добавлены новые схемы ответов и обновлены PUT методы

### Go код:
- `/var/www/go/internal/modules/offer/model.go` - добавлена `OfferUpdateResponse`
- `/var/www/go/internal/modules/offer/service.go` - обновлен метод `UpdateOffer`
- `/var/www/go/internal/modules/offer/handler.go` - обновлен обработчик `UpdateOffer`
- `/var/www/go/internal/modules/order/model.go` - добавлена `OrderStatusUpdateResponse`
- `/var/www/go/internal/modules/order/service.go` - обновлен метод `UpdateOrderStatus`
- `/var/www/go/internal/modules/order/handler.go` - обновлен обработчик `UpdateOrderStatus`
- `/var/www/go/internal/modules/warehouse/model.go` - добавлена `WarehouseUpdateResponse`
- `/var/www/go/internal/modules/warehouse/service.go` - обновлен метод `UpdateWarehouse`
- `/var/www/go/internal/modules/warehouse/handler.go` - обновлен обработчик `UpdateWarehouse`

## Примечание

Все PUT методы теперь возвращают только измененные данные вместо полных объектов, что делает API более эффективным и информативным для клиентов. 