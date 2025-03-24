CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    images TEXT[] NOT NULL,
    thumbnail VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) DEFAULT 0,
    status BOOLEAN DEFAULT true,
    always_available BOOLEAN DEFAULT false,
    quantity INTEGER NOT NULL DEFAULT 0,
    file VARCHAR(255),
    store VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    subcategory VARCHAR(255) NOT NULL,
    type VARCHAR(50),
    unit_sold INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_store ON products(store);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);