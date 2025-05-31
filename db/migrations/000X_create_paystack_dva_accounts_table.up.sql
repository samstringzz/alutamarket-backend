CREATE TABLE IF NOT EXISTS paystack_dva_accounts (
    id VARCHAR(100) PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id),
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(100) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    bank_code VARCHAR(20),
    email VARCHAR(100) NOT NULL,
    customer_code VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_store_dva UNIQUE (store_id)
);

-- Create indexes for faster queries
CREATE INDEX idx_paystack_dva_store_id ON paystack_dva_accounts(store_id);
CREATE INDEX idx_paystack_dva_email ON paystack_dva_accounts(email);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_paystack_dva_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_paystack_dva_accounts_updated_at
    BEFORE UPDATE ON paystack_dva_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_paystack_dva_updated_at(); 