package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-pismo-challenge/pkg/model"
)

type operationRepo struct {
	db *sql.DB
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -package mocks -destination=./mocks/operation_type_mock.go -source=operation_type.go
type OperationTypeConnector interface {
	Get(ctx context.Context, id int) (model.OperationType, error)
}

func NewOperationTypeRepo(db *sql.DB) OperationTypeConnector {
	return &operationRepo{
		db,
	}
}

func (o *operationRepo) Get(ctx context.Context, id int) (model.OperationType, error) {
	checkIdempotencySQL := `SELECT operation_type_id, is_credit FROM transactions.operation_types WHERE operation_type_id = $1;`

	rows := o.db.QueryRowContext(ctx, checkIdempotencySQL, id)
	if rows.Err() != nil {
		return model.OperationType{}, fmt.Errorf("failed to query: %w", rows.Err())
	}
	var ot model.OperationType
	if err := rows.Scan(
		&ot.OperationTypeID,
		&ot.IsCredit,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OperationType{}, ErrNoRows
		}

		return model.OperationType{}, fmt.Errorf("failed to scan while idempotency check: %w", err)
	}

	return ot, nil
}
