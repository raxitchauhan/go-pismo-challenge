package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-pismo-challenge/pkg/model"
)

type transactionRepo struct {
	db *sql.DB
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -package mocks -destination=./mocks/transaction_mock.go -source=transaction.go
type TransactionConnector interface {
	Create(ctx context.Context, a model.Transaction) error
	CheckIdempotency(ctx context.Context, trxUUID string) error
}

func NewTransactionRepo(db *sql.DB) TransactionConnector {
	return &transactionRepo{
		db,
	}
}

func (a *transactionRepo) Create(ctx context.Context, transaction model.Transaction) error {
	insertSQL := `INSERT INTO transactions.transaction (uuid, account_uuid, operation_type_id, amount, event_date) 
					values ($1, $2, $3, $4, $5) ON CONFLICT (uuid) DO NOTHING;`

	_, err := a.db.ExecContext(ctx, insertSQL,
		transaction.UUID.String(),
		transaction.AccountUUID.String(),
		transaction.OperationTypeID,
		transaction.Amount,
		transaction.EventDate,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}

func (a *transactionRepo) CheckIdempotency(ctx context.Context, trxUUID string) error {
	checkIdempotencySQL := `SELECT COALESCE(COUNT(*), 0) as count FROM transactions.transaction where uuid = $1;`

	rows := a.db.QueryRowContext(ctx, checkIdempotencySQL, trxUUID)
	if rows.Err() != nil {
		return fmt.Errorf("failed to query: %w", rows.Err())
	}
	var count uint64
	if err := rows.Scan(&count); err != nil {
		return fmt.Errorf("failed to scan while idempotency check: %w", err)
	}

	if count > 0 {
		return ErrDuplicate
	}

	return nil
}
