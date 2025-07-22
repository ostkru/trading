#!/bin/bash

# Применяем описания к параметрам Product
jq '.components.schemas.Product.properties.name.description = "Название продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Product.properties.vendor_article.description = "Артикул поставщика"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Product.properties.recommend_price.description = "Рекомендуемая цена"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Product.properties.brand.description = "Бренд продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Product.properties.category.description = "Категория продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Product.properties.description.description = "Описание продукта и его характеристики"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам CreateProductRequest
jq '.components.schemas.CreateProductRequest.properties.name.description = "Название продукта (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateProductRequest.properties.vendor_article.description = "Артикул поставщика (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateProductRequest.properties.recommend_price.description = "Рекомендуемая цена"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateProductRequest.properties.brand.description = "Бренд продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateProductRequest.properties.category.description = "Категория продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateProductRequest.properties.description.description = "Описание продукта и его характеристики"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам UpdateProductRequest
jq '.components.schemas.UpdateProductRequest.properties.name.description = "Новое название продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateProductRequest.properties.vendor_article.description = "Новый артикул поставщика"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateProductRequest.properties.recommend_price.description = "Новая рекомендуемая цена"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateProductRequest.properties.brand.description = "Новый бренд продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateProductRequest.properties.category.description = "Новая категория продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateProductRequest.properties.description.description = "Новое описание продукта и его характеристики"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам Offer
jq '.components.schemas.Offer.properties.offer_id.description = "Уникальный идентификатор оффера"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.user_id.description = "ID пользователя-создателя"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.is_public.description = "Публичный ли оффер"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.product_id.description = "ID связанного продукта"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.price_per_unit.description = "Цена за единицу"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.tax_nds.description = "НДС в процентах"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.units_per_lot.description = "Количество единиц в лоте"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.available_lots.description = "Доступное количество лотов"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.latitude.description = "Широта склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.longitude.description = "Долгота склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.warehouse_id.description = "ID склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.offer_type.description = "Тип оффера (sale/buy)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Offer.properties.max_shipping_days.description = "Максимальное количество дней доставки"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам CreateOfferRequest
jq '.components.schemas.CreateOfferRequest.properties.product_id.description = "ID продукта (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.offer_type.description = "Тип оффера: sale или buy (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.price_per_unit.description = "Цена за единицу (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.available_lots.description = "Доступное количество лотов (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.tax_nds.description = "НДС в процентах (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.units_per_lot.description = "Количество единиц в лоте (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.warehouse_id.description = "ID склада (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.is_public.description = "Публичный ли оффер"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOfferRequest.properties.max_shipping_days.description = "Максимальное количество дней доставки"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам UpdateOfferRequest
jq '.components.schemas.UpdateOfferRequest.properties.price_per_unit.description = "Новая цена за единицу"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.available_lots.description = "Новое количество доступных лотов"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.tax_nds.description = "Новый НДС в процентах"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.units_per_lot.description = "Новое количество единиц в лоте"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.is_public.description = "Новый статус публичности"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.max_shipping_days.description = "Новое максимальное количество дней доставки"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOfferRequest.properties.warehouse_id.description = "Новый ID склада"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам Order
jq '.components.schemas.Order.properties.order_id.description = "Уникальный идентификатор заказа"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.total_amount.description = "Общая сумма заказа"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.offer_id.description = "ID оффера"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.initiator_user_id.description = "ID покупателя (инициатора)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.counterparty_user_id.description = "ID продавца"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.order_time.description = "Время создания заказа"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.price_per_unit.description = "Цена за единицу"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.units_per_lot.description = "Количество единиц в лоте"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.lot_count.description = "Количество лотов"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.notes.description = "Примечания к заказу"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.order_type.description = "Тип заказа"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.payment_method.description = "Способ оплаты"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.order_status.description = "Статус заказа"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.status_reason.description = "Причина изменения статуса"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.status_changed_at.description = "Время изменения статуса"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.status_changed_by.description = "Кто изменил статус"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.shipping_address.description = "Адрес доставки"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.tracking_number.description = "Номер отслеживания"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Order.properties.max_shipping_days.description = "Максимальное количество дней доставки"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам CreateOrderRequest
jq '.components.schemas.CreateOrderRequest.properties.offer_id.description = "ID оффера (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateOrderRequest.properties.quantity.description = "Количество лотов (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам UpdateOrderStatusRequest
jq '.components.schemas.UpdateOrderStatusRequest.properties.status.description = "Новый статус заказа (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateOrderStatusRequest.properties.reason.description = "Причина изменения статуса"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам Warehouse
jq '.components.schemas.Warehouse.properties.id.description = "Уникальный идентификатор склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.user_id.description = "ID владельца склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.name.description = "Название склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.longitude.description = "Долгота склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.latitude.description = "Широта склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.wb_id.description = "ID склада в Wildberries"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.working_hours.description = "Рабочие часы"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.Warehouse.properties.address.description = "Адрес склада"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам CreateWarehouseRequest
jq '.components.schemas.CreateWarehouseRequest.properties.name.description = "Название склада (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateWarehouseRequest.properties.address.description = "Адрес склада (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateWarehouseRequest.properties.latitude.description = "Широта склада (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateWarehouseRequest.properties.longitude.description = "Долгота склада (обязательно)"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.CreateWarehouseRequest.properties.working_hours.description = "Рабочие часы"' openapi.json > temp.json && mv temp.json openapi.json

# Применяем описания к параметрам UpdateWarehouseRequest
jq '.components.schemas.UpdateWarehouseRequest.properties.name.description = "Новое название склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateWarehouseRequest.properties.address.description = "Новый адрес склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateWarehouseRequest.properties.latitude.description = "Новая широта склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateWarehouseRequest.properties.longitude.description = "Новая долгота склада"' openapi.json > temp.json && mv temp.json openapi.json
jq '.components.schemas.UpdateWarehouseRequest.properties.working_hours.description = "Новые рабочие часы"' openapi.json > temp.json && mv temp.json openapi.json

echo "Все описания применены!" 