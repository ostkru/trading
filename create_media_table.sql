-- Создание таблицы media для хранения медиа-контента продуктов

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_urls JSON, -- ссылки на изображения товара (массив URL)
    video_urls JSON, -- ссылки на видео обзоры (массив URL)
    model_3d_urls JSON, -- ссылки на 3д модели (массив URL)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации запросов
CREATE INDEX idx_media_product_id ON media(product_id); 