-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS transactions;

CREATE TABLE IF NOT EXISTS transactions.operation_types (
    operation_type_id SERIAL PRIMARY KEY,
    description TEXT UNIQUE NOT NULL,
    is_credit BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions.transaction (
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    account_uuid TEXT NOT NULL,
    operation_type_id INTEGER NOT NULL,
    amount DECIMAL(10,2) DEFAULT 0,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP SCHEMA IF EXISTS transactions CASCADE;

-- +goose StatementEnd