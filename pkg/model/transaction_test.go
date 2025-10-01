package model

import (
	"go-pismo-challenge/pkg/util"
	"testing"
)

// TestTransactionRequest_Validate tests the Validate method of TransactionRequest
func TestTransactionRequest_Validate(t *testing.T) {
	tests := []struct {
		name           string
		transactionReq TransactionRequest
		wantErr        []util.FieldError // expected validation errors
	}{
		{
			name: "Valid transaction request",
			transactionReq: TransactionRequest{
				AccountUUID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationTypeID: 1,
				Amount:          100.0,
				IdempotencyKey:  "some-string",
			},
			wantErr: nil, // no errors expected
		},
		{
			name: "Missing idempotency key",
			transactionReq: TransactionRequest{
				AccountUUID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationTypeID: 1,
				Amount:          100.0,
				IdempotencyKey:  "",
			},
			wantErr: []util.FieldError{
				{Field: "idempotency_key", Message: "field is required"},
			},
		},
		{
			name: "Invalid UUID format",
			transactionReq: TransactionRequest{
				AccountUUID:     "invalid-uuid-format",
				OperationTypeID: 1,
				Amount:          100.0,
				IdempotencyKey:  "some-string",
			},
			wantErr: []util.FieldError{
				{Field: "account_uuid", Message: "invalid UUID format"},
			},
		},
		{
			name: "Negative operation type ID",
			transactionReq: TransactionRequest{
				AccountUUID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationTypeID: -1,
				Amount:          100.0,
				IdempotencyKey:  "some-string",
			},
			wantErr: []util.FieldError{
				{Field: "operation_type_id", Message: "must be a positive integer: -1"},
			},
		},
		{
			name: "Negative amount",
			transactionReq: TransactionRequest{
				AccountUUID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationTypeID: 1,
				Amount:          -50.0,
				IdempotencyKey:  "some-string",
			},
			wantErr: []util.FieldError{
				{Field: "amount", Message: "must be a non-negative value: -50.00"},
			},
		},
	}

	// loop over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.transactionReq.Validate()
			if len(gotErr) != len(tt.wantErr) {
				t.Errorf("Expected %d errors, got %d", len(tt.wantErr), len(gotErr))

				return
			}
			for i, err := range gotErr {
				if err.Field != tt.wantErr[i].Field {
					t.Errorf("Expected error %v, got %v", tt.wantErr[i], err)
				}
			}
		})
	}
}
