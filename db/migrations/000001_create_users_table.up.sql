CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    campus VARCHAR(100) NOT NULL,
    usertype VARCHAR(50) NOT NULL,
    active BOOLEAN DEFAULT true,
    twofa BOOLEAN DEFAULT false,
    store_name VARCHAR(255),
    store_email VARCHAR(255),
    has_physical_address BOOLEAN DEFAULT false,
    access_token TEXT,
    refresh_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);