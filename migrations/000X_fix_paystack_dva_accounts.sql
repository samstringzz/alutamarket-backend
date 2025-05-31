-- Drop the existing table if it exists
DROP TABLE IF EXISTS paystack_dva_accounts;

-- Create the table with proper auto-increment
CREATE TABLE paystack_dva_accounts (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 