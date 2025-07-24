-- Восстановление структуры БД PortalData в MySQL

USE portaldata;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    vendor_article VARCHAR(255) NOT NULL,
    recommend_price DECIMAL(10,2),
    brand VARCHAR(255),
    category VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE warehouses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    longitude DECIMAL(10,8),
    latitude DECIMAL(10,8),
    wb_id VARCHAR(255),
    working_hours VARCHAR(255),
    address VARCHAR(255)
);

CREATE TABLE offers (
    offer_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    is_public BOOLEAN DEFAULT 1,
    product_id INT,
    price_per_unit DECIMAL(10,2),
    tax_nds INT,
    units_per_lot INT,
    available_lots INT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(10,8),
    warehouse_id INT,
    offer_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

CREATE TABLE orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    total_amount DECIMAL(10,2),
    is_multi BOOLEAN,
    offer_id INT,
    initiator_user_id INT,
    counterparty_user_id INT,
    order_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    price_per_unit DECIMAL(10,2),
    units_per_lot INT,
    lot_count INT,
    notes TEXT,
    order_type VARCHAR(50),
    payment_method VARCHAR(50),
    order_status VARCHAR(50),
    shipping_address VARCHAR(255),
    tracking_number VARCHAR(255),
    FOREIGN KEY (offer_id) REFERENCES offers(offer_id)
);

CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    offer_id INT,
    qty INT,
    price_per_unit DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50),
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (offer_id) REFERENCES offers(offer_id)
); 