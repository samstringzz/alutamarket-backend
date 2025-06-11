-- Add store_id column if it doesn't exist
ALTER TABLE paystack_dva_accounts ADD COLUMN IF NOT EXISTS store_id INTEGER;

-- Update existing Paystack DVA accounts with store IDs
UPDATE paystack_dva_accounts pda
SET store_id = s.id
FROM stores s
JOIN users u ON s.user_id = u.id
WHERE pda.email = u.email;

-- Add unique constraint on store_id
ALTER TABLE paystack_dva_accounts ADD CONSTRAINT unique_store_paystack_dva UNIQUE (store_id);

-- Add foreign key constraint
ALTER TABLE paystack_dva_accounts ADD CONSTRAINT fk_store_paystack_dva FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE;

-- Create index on store_id for better performance
CREATE INDEX IF NOT EXISTS idx_paystack_dva_accounts_store_id ON paystack_dva_accounts(store_id); 