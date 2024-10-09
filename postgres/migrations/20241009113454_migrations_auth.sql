-- +goose Up
-- The table stores data about all registered chat users
CREATE TYPE user_roles AS ENUM ('ADMIN', 'USER');
CREATE TABLE chat_user (
    id SERIAL PRIMARY KEY,      
    name TEXT NOT NULL,         
    password TEXT NOT NULL,     
    email TEXT,                 
    role user_roles,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE chat_user;
