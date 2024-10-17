-- +goose Up
-- The table stores data about all registered chat users
CREATE TABLE chat_user (
    id SERIAL PRIMARY KEY,      
    name VARCHAR(25) NOT NULL,         
    password VARCHAR(25) NOT NULL,     
    email VARCHAR(50),                 
    role VARCHAR(5) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE chat_user;
