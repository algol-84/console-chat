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

CREATE TABLE chat_permissions (
    endpoint TEXT,
    role TEXT NOT NULL
);

INSERT INTO chat_permissions (endpoint, role) VALUES 
('/user_v1.UserV1/Create', 'ADMIN'),
('/user_v1.UserV1/Get', 'USER'),
('/user_v1.UserV1/Update', 'ADMIN'),
('/user_v1.UserV1/Delete', 'ADMIN');

-- +goose Down
DROP TABLE chat_user;
DROP TABLE chat_permissions;
