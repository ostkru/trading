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

CREATE TABLE api_rate_limits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    request_count INTEGER DEFAULT 1,
    minute_count INTEGER DEFAULT 1,
    day_count INTEGER DEFAULT 1,
    last_request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    minute_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    day_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_urls JSON, -- ссылки на изображения товара (массив URL)
    video_urls JSON, -- ссылки на видео обзоры (массив URL)
    model_3d_urls JSON, -- ссылки на 3д модели (массив URL)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    media_type VARCHAR(50) NOT NULL CHECK (media_type IN ('image', 'video', '3d_model')),
    url VARCHAR(500) NOT NULL,
    title VARCHAR(255),
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    file_size BIGINT,
    mime_type VARCHAR(100),
    duration INTEGER, -- для видео в секундах
    width INTEGER, -- ширина изображения/видео
    height INTEGER, -- высота изображения/видео
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    UNIQUE(product_id, url)
); 