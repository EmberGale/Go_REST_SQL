-- +goose Up
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    person VARCHAR(255) NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(255) NOT NULL,
    user_id INTEGER,
    time TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS payments;
