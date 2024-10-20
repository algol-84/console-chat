# chat-server

Адрес локального сервиса localhost:50052
Адрес удаленного сервиса 45.130.9.109:50062

Адрес локальной базы localhost:54322
Адрес удаленной базы 45.130.9.109:54332


Нужно спроектировать таблицу базы данных которая хранит список чатов. Каждый чат должен содержать имена юзеров этого чата
Answer
Проектирование таблицы базы данных для хранения списка чатов, где каждый чат содержит имена пользователей, требует учета структуры данных и взаимосвязей между пользователями и чатами. Ниже приведен пример схемы базы данных, которая эффективно организует эту информацию.
Структура базы данных
1. Таблица пользователей (Users)
Эта таблица будет хранить информацию о каждом пользователе.
sql
CREATE TABLE Users (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

2. Таблица чатов (Chats)
Эта таблица будет хранить информацию о каждом чате.
sql
CREATE TABLE Chats (
    chat_id INT PRIMARY KEY AUTO_INCREMENT,
    chat_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

3. Таблица участников чата (ChatParticipants)
Эта таблица будет связывать пользователей с чатами, позволяя одному пользователю участвовать в нескольких чатах и одному чату содержать нескольких пользователей.
sql
CREATE TABLE ChatParticipants (
    chat_id INT,
    user_id INT,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES Chats(chat_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

Пример использования
Вставка пользователей:
sql
INSERT INTO Users (username) VALUES ('Alice'), ('Bob'), ('Charlie');

Вставка чатов:
sql
INSERT INTO Chats (chat_name) VALUES ('Chat Group 1'), ('Chat Group 2');

Вставка участников чата:
sql
INSERT INTO ChatParticipants (chat_id, user_id) VALUES 
(1, 1),  -- Alice в Chat Group 1
(1, 2),  -- Bob в Chat Group 1
(2, 2),  -- Bob в Chat Group 2
(2, 3);  -- Charlie в Chat Group 2

Заключение
Эта структура базы данных позволяет эффективно хранить информацию о чатах и участниках. Использование отдельной таблицы для участников чата обеспечивает гибкость и масштабируемость. Вы можете легко добавлять новых пользователей и чаты, а также управлять участниками без необходимости изменения основной структуры таблиц.