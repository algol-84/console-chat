-- +goose Up
-- The table stores created chat rooms
CREATE TABLE chat_room (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- The table stores all chat messages
CREATE TABLE chat_message (
    id SERIAL PRIMARY KEY,
    chat_id SERIAL,
    from_user_id SERIAL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

-- +goose Down
drop table chat_room;
drop table chat_message;