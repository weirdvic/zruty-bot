-- Создание таблицы для хранения конфигурации бота
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    username TEXT,
    is_bot BOOLEAN DEFAULT FALSE,
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    check_passed_at TIMESTAMP DEFAULT NULL
);

-- Создание таблицы администраторов
CREATE TABLE IF NOT EXISTS admins (
    user_id INTEGER PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создание таблицы групп
CREATE TABLE IF NOT EXISTS chats (
    chat_id INTEGER PRIMARY KEY,
    title TEXT,
    username TEXT,
    chat_type TEXT
);

CREATE TABLE IF NOT EXISTS user_chats (
    user_id INTEGER,
    chat_id INTEGER,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_id) REFERENCES chats(chat_id)
);