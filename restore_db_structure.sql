-- Восстановление структуры БД PortalData

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    vendor_article VARCHAR NOT NULL,
    recommend_price NUMERIC,
    brand VARCHAR,
    category VARCHAR,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    user_id INTEGER,
    category_id INTEGER,
    brand_id INTEGER
);

CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    updated_at TIMESTAMP,
    created_at TIMESTAMP,
    longitude NUMERIC,
    latitude NUMERIC,
    wb_id VARCHAR,
    working_hours VARCHAR,
    address VARCHAR
);

CREATE TABLE offers (
    offer_id SERIAL PRIMARY KEY,
    user_id INTEGER,
    is_public BOOLEAN,
    product_id INTEGER REFERENCES products(id),
    price_per_unit NUMERIC,
    tax_nds INTEGER,
    units_per_lot INTEGER,
    available_lots INTEGER,
    latitude NUMERIC,
    longitude NUMERIC,
    warehouse_id INTEGER REFERENCES warehouses(id),
    offer_type VARCHAR,
    max_shipping_days INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    total_amount NUMERIC,
    is_multi BOOLEAN,
    offer_id INTEGER REFERENCES offers(offer_id),
    initiator_user_id INTEGER,
    counterparty_user_id INTEGER,
    order_time TIMESTAMP,
    price_per_unit NUMERIC,
    units_per_lot INTEGER,
    lot_count INTEGER,
    notes TEXT,
    order_type VARCHAR,
    payment_method VARCHAR,
    order_status VARCHAR,
    shipping_address VARCHAR,
    tracking_number VARCHAR,
    max_shipping_days INTEGER
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(order_id),
    offer_id INTEGER REFERENCES offers(offer_id),
    qty INTEGER,
    price_per_unit NUMERIC,
    created_at TIMESTAMP,
    status VARCHAR
); 