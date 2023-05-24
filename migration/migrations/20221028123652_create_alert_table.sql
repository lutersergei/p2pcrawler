-- +goose Up
CREATE TABLE IF NOT EXISTS alert_request
(
    id         INTEGER PRIMARY KEY,
    price      REAL,
    username   TEXT,
    exchange   TEXT,
    move_type  TEXT,
    deal_type  TEXT,
    status     TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS alert_request;