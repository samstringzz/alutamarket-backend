CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    cart_id INTEGER NOT NULL,
    coupon TEXT,
    fee NUMERIC(10,2) DEFAULT 0,
    status VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    customer JSONB NOT NULL,
    seller_id VARCHAR(50),
    stores_id TEXT[] DEFAULT '{}',
    delivery_details JSONB,
    amount NUMERIC(10,2) NOT NULL,
    uuid VARCHAR(100) UNIQUE NOT NULL,
    payment_gateway VARCHAR(50),
    payment_method VARCHAR(50),
    trans_ref VARCHAR(100),
    trans_status VARCHAR(50),
    products JSONB DEFAULT '[]',

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(uuid) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_uuid ON orders(uuid);
CREATE INDEX idx_orders_cart_id ON orders(cart_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(status);