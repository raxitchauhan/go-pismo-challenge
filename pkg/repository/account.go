package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-pismo-challenge/pkg/model"
)

type accountRepo struct {
	db *sql.DB
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -package mocks -destination=./mocks/account_mock.go -source=account.go
type AccountConnector interface {
	Create(ctx context.Context, a model.Account) error
	CheckIdempotency(ctx context.Context, accountUUID string) error
	Get(ctx context.Context, uuid string) (model.Account, error)
}

func NewAccountRepo(db *sql.DB) AccountConnector {
	return &accountRepo{
		db,
	}
}

func (a *accountRepo) Create(ctx context.Context, account model.Account) error {
	insertSQL := `INSERT INTO accounts.account (uuid, document_number, created_at) values ($1, $2, $3) ON CONFLICT (uuid) DO NOTHING;`

	_, err := a.db.ExecContext(ctx, insertSQL, account.UUID.String(), account.DocumentNumber, account.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	return nil
}

func (a *accountRepo) CheckIdempotency(ctx context.Context, accountUUID string) error {
	checkIdempotencySQL := `SELECT COALESCE(COUNT(*), 0) as count FROM accounts.account where uuid = $1;`

	rows := a.db.QueryRowContext(ctx, checkIdempotencySQL, accountUUID)
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

func (a *accountRepo) Get(ctx context.Context, uuid string) (model.Account, error) {
	getAccount := `SELECT uuid, document_number, created_at FROM accounts.account where uuid = $1;`

	rows := a.db.QueryRowContext(ctx, getAccount, uuid)
	if rows.Err() != nil {
		return model.Account{}, fmt.Errorf("failed to query account: %w", rows.Err())
	}
	var account model.Account
	if err := rows.Scan(
		&account.UUID,
		&account.DocumentNumber,
		&account.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Account{}, ErrNoRows
		}

		return model.Account{}, fmt.Errorf("failed to scan account: %w", err)
	}

	return account, nil
}
