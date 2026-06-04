-- +goose Up
CREATE TABLE IF NOT EXISTS outbox_events (
    id SERIAL PRIMARY KEY,
    payment_id INTEGER NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    payload TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    topic VARCHAR(255)
);

-- Индекс для быстрого поиска pending событий с next_retry_at <= NOW()
CREATE INDEX IF NOT EXISTS idx_outbox_events_status_next_retry
    ON outbox_events(status, next_retry_at, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_outbox_events_status_next_retry;
DROP TABLE IF EXISTS outbox_events;