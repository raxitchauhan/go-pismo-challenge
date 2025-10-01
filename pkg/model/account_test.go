package model

import (
	"go-pismo-challenge/pkg/util"
	"testing"
)

// TestAccountRequest_Validate tests the Validate method of AccountRequest
func TestAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		accountReq AccountRequest
		wantErr    []util.FieldError // expected validation errors
	}{
		{
			name: "Valid account request",
			accountReq: AccountRequest{
				DocumentNumber: "12345",
				IdempotencyKey: "abcde-12345",
			},
			wantErr: nil, // no errors expected
		},
		{
			name: "Missing document number",
			accountReq: AccountRequest{
				DocumentNumber: "",
				IdempotencyKey: "abcde-12345",
			},
			wantErr: []util.FieldError{
				{Field: "document_number", Message: "field is required"},
			},
		},
		{
			name: "Missing idempotency key",
			accountReq: AccountRequest{
				DocumentNumber: "12345",
				IdempotencyKey: "",
			},
			wantErr: []util.FieldError{
				{Field: "idempotency_key", Message: "field is required"},
			},
		},
		{
			name: "Missing both fields",
			accountReq: AccountRequest{
				DocumentNumber: "",
				IdempotencyKey: "",
			},
			wantErr: []util.FieldError{
				{Field: "idempotency_key", Message: "field is required"},
				{Field: "document_number", Message: "field is required"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.accountReq.Validate()
			if len(gotErr) != len(tt.wantErr) {
				t.Errorf("Expected %d errors, got %d", len(tt.wantErr), len(gotErr))

				return
			}
			for i, err := range gotErr {
				if err.Field != tt.wantErr[i].Field || err.Message != tt.wantErr[i].Message {
					t.Errorf("Expected error %v, got %v", tt.wantErr[i], err)
				}
			}
		})
	}
}
