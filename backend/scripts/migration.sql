-- Создание таблицы для хранения ссылок
CREATE TABLE IF NOT EXISTS assistant_links (
    id SERIAL PRIMARY KEY,               -- Уникальный идентификатор ссылки
    assistant_id UUID NOT NULL,           -- Идентификатор ассистента (связан с администратором)
    url TEXT NOT NULL,                    -- URL или ссылка
    created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT), -- Время создания ссылки
    updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT), -- Время последнего обновления
    FOREIGN KEY (assistant_id) REFERENCES admins(id) ON DELETE CASCADE -- Связь с таблицей админов
);

-- Таблица для статистики запросов и ответов
CREATE TABLE IF NOT EXISTS assistant_statistics (
    id SERIAL PRIMARY KEY,               -- Уникальный идентификатор записи
    link_id INT NOT NULL,                 -- Ссылка, к которой привязана статистика
    request_time BIGINT DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT), -- Время запроса
    request_text TEXT NOT NULL,           -- Запрос пользователя
    response_text TEXT NOT NULL,          -- Ответ модели
    FOREIGN KEY (link_id) REFERENCES assistant_links(id) ON DELETE CASCADE -- Связь с таблицей ссылок
);

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_assistant_links_assistant_id ON assistant_links (assistant_id);
CREATE INDEX IF NOT EXISTS idx_assistant_statistics_link_id ON assistant_statistics (link_id);
