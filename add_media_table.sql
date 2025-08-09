-- Добавление таблицы media для хранения медиаданных продуктов

USE portaldata;

-- Создание таблицы media
CREATE TABLE IF NOT EXISTS media (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    image_urls JSON,
    video_urls JSON,
    model_3d_urls JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_product_media (product_id)
);

-- Добавление недостающих полей в таблицу products (если их нет)
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS brand_id INT NULL,
ADD COLUMN IF NOT EXISTS category_id INT NULL,
ADD COLUMN IF NOT EXISTS barcode VARCHAR(50) NULL,
ADD COLUMN IF NOT EXISTS status ENUM('pending', 'classified', 'not_classified') DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS user_id INT NULL;

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_media_product_id ON media(product_id);
CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_brand_category ON products(brand_id, category_id);
