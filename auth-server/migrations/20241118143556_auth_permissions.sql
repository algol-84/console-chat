-- +goose Up
CREATE TABLE user_permissions (
    endpoint TEXT,
    role TEXT NOT NULL
);

INSERT INTO user_permissions (endpoint, role) VALUES 
('/chat_v1.ChatV1/CreateChat', 'ADMIN'),
('/chat_v1.ChatV1/ConnectChat', 'USER'),
('/chat_v1.ChatV1/SendMessage', 'USER');  

-- +goose Down
DROP TABLE user_permissions;
