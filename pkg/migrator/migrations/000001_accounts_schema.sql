-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS accounts;

CREATE TABLE IF NOT EXISTS accounts.account (
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    document_number TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP SCHEMA IF EXISTS accounts CASCADE;

-- +goose StatementEnd