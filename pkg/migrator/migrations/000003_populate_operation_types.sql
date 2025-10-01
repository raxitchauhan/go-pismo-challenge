-- +goose Up
-- +goose StatementBegin
INSERT INTO transactions.operation_types (description, is_credit) values ('CASH_PURCHASE', false);
INSERT INTO transactions.operation_types (description, is_credit) values ('INSTALLMENT_PURCHASE', false);
INSERT INTO transactions.operation_types (description, is_credit) values ('WITHDRAWAL', false);
INSERT INTO transactions.operation_types (description, is_credit) values ('PAYMENT', true);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE transactions.operation_types;

-- +goose StatementEnd