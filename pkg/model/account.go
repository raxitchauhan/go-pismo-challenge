package model

import (
	"go-pismo-challenge/pkg/util"
	"time"

	"github.com/gofrs/uuid"
)

type AccountRequest struct {
	DocumentNumber string `json:"document_number" example:"some-string" format:"string"`
	IdempotencyKey string `json:"idempotency_key" example:"some-string" format:"string"`
}

type AccountResponse struct {
	UUID string `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

type Account struct {
	UUID           uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	DocumentNumber string    `json:"document_number" example:"some-string" format:"string"`
	CreatedAt      time.Time `json:"created_at" example:"2025-10-01T06:22:46.931755Z" format:"time"`
}

func (a AccountRequest) Validate() []util.FieldError {
	err := make([]util.FieldError, 0)
	if a.IdempotencyKey == "" {
		err = append(err, util.FieldError{
			Field:   "idempotency_key",
			Message: "field is required",
		})
	}
	if a.DocumentNumber == "" {
		err = append(err, util.FieldError{
			Field:   "document_number",
			Message: "field is required",
		})
	}

	return err
}
