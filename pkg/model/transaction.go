package model

import (
	"fmt"
	"go-pismo-challenge/pkg/util"
	"time"

	"github.com/gofrs/uuid"
)

type TransactionRequest struct {
	AccountUUID     string  `json:"account_uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	OperationTypeID int     `json:"operation_type_id" example:"1" format:"int64"`
	Amount          float64 `json:"amount" example:"1.1" format:"float64"`
	IdempotencyKey  string  `json:"idempotency_key" example:"some-string" format:"string"`
}

func (t TransactionRequest) Validate() []util.FieldError {
	vErr := make([]util.FieldError, 0)
	if t.IdempotencyKey == "" {
		vErr = append(vErr, util.FieldError{
			Field:   "idempotency_key",
			Message: "field is required",
		})
	}
	if _, err := uuid.FromString(t.AccountUUID); err != nil {
		vErr = append(vErr, util.FieldError{
			Field:   "account_uuid",
			Message: fmt.Sprintf("invalid uuid: '%s'", t.AccountUUID),
		})
	}
	if t.OperationTypeID <= 0 {
		vErr = append(vErr, util.FieldError{
			Field:   "operation_type_id",
			Message: fmt.Sprintf("field is required and non-negative: %d", t.OperationTypeID),
		})
	}
	if t.Amount < 0 {
		vErr = append(vErr, util.FieldError{
			Field:   "amount",
			Message: fmt.Sprintf("field should be non-negative: %0.2f", t.Amount),
		})
	}

	return vErr
}

type TransactionResponse struct {
	UUID string `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

type Transaction struct {
	UUID            uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	AccountUUID     uuid.UUID `json:"account_uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	OperationTypeID int       `json:"operation_type_id" example:"1" format:"int64"`
	Amount          float64   `json:"amount" example:"1.1" format:"float64"`
	EventDate       time.Time `json:"event_date" example:"2025-10-01T06:22:46.931755Z" format:"time"`
}
