-- +goose Up
CREATE TABLE user_permissions (
    endpoint TEXT,
    role TEXT NOT NULL
);

INSERT INTO user_permissions (endpoint, role) VALUES 
('/user_v1.UserV1/Create', 'ADMIN'),
('/user_v1.UserV1/Get', 'USER'),
('/user_v1.UserV1/Update', 'ADMIN'),
('/user_v1.UserV1/Delete', 'ADMIN');

-- +goose Down
DROP TABLE user_permissions;
