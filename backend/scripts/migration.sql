-- DROP TABLE IF EXISTS users;

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT)
);