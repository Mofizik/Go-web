-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id        TEXT PRIMARY KEY,
    item      TEXT NOT NULL,
    quantity  INT  NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;