-- Drop the existing table if it exists
DROP TABLE IF EXISTS paystack_dva_accounts;

-- Create the table with all required fields
CREATE TABLE paystack_dva_accounts (
    id VARCHAR(100) PRIMARY KEY,
    store_id INT NOT NULL,
    account_number VARCHAR(20),
    account_name VARCHAR(100),
    bank_name VARCHAR(100),
    bank_code VARCHAR(20),
    email VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id)
); 