-- +goose Up
-- The table stores data about all registered chat users
CREATE TABLE chat_user (
    id SERIAL PRIMARY KEY,      
    name TEXT NOT NULL,         
    password TEXT NOT NULL,     
    email TEXT,                 
    role TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE chat_user;
