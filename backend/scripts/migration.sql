-- Удаление таблиц, если они существуют
-- DROP TABLE IF EXISTS assistant_statistics CASCADE;
-- DROP TABLE IF EXISTS assistant_links CASCADE;
-- DROP TABLE IF EXISTS admins CASCADE;

-- Создание таблицы для администраторов с UUID в качестве id
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  -- Используем UUID в качестве идентификатора
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT)
);

-- Создание таблицы для ссылок ассистента с UUID
CREATE TABLE IF NOT EXISTS assistant_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),   -- UUID вместо SERIAL
    assistant_id UUID NOT NULL,                        -- UUID для связи с ассистентом
    url TEXT NOT NULL,                                 -- URL для ассистента
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT),
    updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT),
    FOREIGN KEY (assistant_id) REFERENCES admins(id) ON DELETE CASCADE
);

-- Создание таблицы для статистики запросов и ответов
CREATE TABLE IF NOT EXISTS assistant_statistics (
    id SERIAL PRIMARY KEY,                            -- SERIAL первичный ключ
    link_id UUID NOT NULL,                            -- UUID вместо INT
    request_time BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT),
    request_text TEXT NOT NULL,
    response_text TEXT NOT NULL,
    FOREIGN KEY (link_id) REFERENCES assistant_links(id) ON DELETE CASCADE  -- Ссылаемся на UUID
);

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_assistant_links_assistant_id ON assistant_links (assistant_id);
CREATE INDEX IF NOT EXISTS idx_assistant_statistics_link_id ON assistant_statistics (link_id);
