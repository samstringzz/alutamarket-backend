CREATE TABLE IF NOT EXISTS store_earnings (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending/released/reversed
    transaction_type VARCHAR(50) NOT NULL, -- order/direct_transfer
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
); 