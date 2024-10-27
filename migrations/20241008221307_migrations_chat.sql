-- +goose Up
-- Таблица содержит список всех созданных чатов
CREATE TABLE chat_room (
    chat_id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица связывает пользователей с чатами
-- Каждый юзер может состоять в нескольких чатах, но нельзя добавить одного юзера в один чат дважды
CREATE TABLE chat_user (
    chat_id INT NOT NULL REFERENCES chat_room(chat_id) ON DELETE CASCADE,
    user_name TEXT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (chat_id, user_name)
);

CREATE TABLE chat_log (
    chat_id INT NOT NULL,
    info TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
drop table chat_user;
drop table chat_room;
drop table chat_log;