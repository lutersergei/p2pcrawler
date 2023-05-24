-- +goose Up
CREATE TABLE IF NOT EXISTS price_history
(
    id             INTEGER PRIMARY KEY,
    best_price     REAL,
    surplus_amount REAL,
    username       TEXT,
    exchange       TEXT,
    raw_json       TEXT,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS price_history;