CREATE DATABASE shop;

\c shop;

-- TODO: update schema
CREATE TABLE IF NOT EXISTS "users" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16),
    password VARCHAR(32),
);

INSERT INTO users (name, password) VALUES ('hello', '5d41402abc4b2a76b9719d911017c592');

CREATE USER shop_service WITH PASSWORD 'shop_service_pass';

GRANT CONNECT ON DATABASE shop TO shop_service;
GRANT USAGE ON SCHEMA public TO shop_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES
IN SCHEMA public TO shop_service; 
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO shop_service;

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_user INTEGER REFERENCES users(id) ON DELETE CASCADE,
    to_user INTEGER REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- TODO: add item type table