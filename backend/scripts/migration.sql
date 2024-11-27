DROP TABLE IF EXISTS uploaded_files, messages, chats, users;
SELECT * FROM users WHERE id = 1;

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Создаем таблицу чатов
CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Внешний ключ, указывающий на таблицу users
    title TEXT DEFAULT 'New Chat',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Создаем таблицу сообщений
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE, -- Внешний ключ, указывающий на таблицу chats
    sender TEXT NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW()
);

-- Создаем таблицу загруженных файлов
CREATE TABLE IF NOT EXISTS uploaded_files (
    id SERIAL PRIMARY KEY,
    unique_id UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id INT REFERENCES users(id) ON DELETE SET NULL, -- Внешний ключ, указывающий на таблицу users
    file_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
