package handler

import (
	"context"
	"errors"
	"go-pismo-challenge/pkg/model"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/repository/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type transactionTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	connector          *Transaction
	mockAccounts       *mocks.MockAccountConnector
	mockTrx            *mocks.MockTransactionConnector
	mockOperationTypes *mocks.MockOperationTypeConnector
	router             *chi.Mux
	recoder            *httptest.ResponseRecorder
}

func TestTransactionHnadler(t *testing.T) {
	suite.Run(t, new(transactionTestSuite))
}

// Setup test suite
func (s *transactionTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockAccounts = mocks.NewMockAccountConnector(s.ctrl)
	s.mockTrx = mocks.NewMockTransactionConnector(s.ctrl)
	s.mockOperationTypes = mocks.NewMockOperationTypeConnector(s.ctrl)

	s.connector = NewTransactionHandler(s.mockTrx, s.mockAccounts, s.mockOperationTypes)
	s.recoder = httptest.NewRecorder()
	s.router = chi.NewRouter()

	s.router.Post("/transactions", s.connector.Create)
}

// Assert expectations
func (s *transactionTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Success: An account was created
//
// Return: 201
func (s *transactionTestSuite) TestCreateTrxSuccess() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	trxUUID := getMockTrxUUID()
	expectedOperationType := model.OperationType{
		OperationTypeID: 4,
		IsCredit:        true,
	}

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)

	s.mockOperationTypes.EXPECT().Get(gomock.Any(), 4).Return(expectedOperationType, nil)

	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, nil)

	s.mockTrx.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, t model.Transaction) error {
			// validate fields
			if t.UUID != trxUUID ||
				t.AccountUUID.String() != "e2a84838-88de-5fbc-8636-6ef49e26f00a" ||
				t.Amount != 1.1 ||
				t.OperationTypeID != 4 {
				return errors.New("incorrect params")
			}

			return nil
		})

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusCreated, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("uuid", string(resBody))
}

// BadRequest: `idempotency_key` field was not passed in the request body
//
// Returns: 400
func (s *transactionTestSuite) TestTransactionBadRequestIdempotency() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1
			}`))
	s.Require().NoError(err)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("idempotency_key", string(resBody))
}

// BadRequest: Duplicate request received, idempotency check returns duplicate error
//
// Returns: 400
func (s *transactionTestSuite) TestTransactionIdempotencyCheckFailed() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(repository.ErrDuplicate)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("duplicate", string(resBody))
}

// BadRequest: Validation failed on all fields
//
// Returns: 400
func (s *transactionTestSuite) TestTransactionValidationFailed() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "invalid-uuid",
				"amount": -9,
				"operation_type_id": -5
			}`))
	s.Require().NoError(err)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)

	s.Regexp("idempotency_key", string(resBody))
	s.Regexp("amount", string(resBody))
	s.Regexp("operation_type_id", string(resBody))
	s.Regexp("account_uuid", string(resBody))
	// sample error message looks something like this
	/*
		{
			"errors": [
				{
					"id": "41f3f002-66bb-4eb6-8890-86b03b370724",
					"code": "validation_error",
					"status": 400,
					"title": "failed to create transaction",
					"detail": "failed to validate request body",
					"source": {
						"field": "idempotency_key",
						"message": "field is required"
					}
				},
				{
					"id": "9959823c-0396-4938-a878-979c33a9a438",
					"code": "validation_error",
					"status": 400,
					"title": "failed to create transaction",
					"detail": "failed to validate request body",
					"source": {
						"field": "account_uuid",
						"message": "invalid uuid: 'invalid-uuid'"
					}
				},
				{
					"id": "1c9c0108-5424-4509-9f1e-c3c51b8c9301",
					"code": "validation_error",
					"status": 400,
					"title": "failed to create transaction",
					"detail": "failed to validate request body",
					"source": {
						"field": "operation_type_id",
						"message": "field is required and non-negative: -5"
					}
				},
				{
					"id": "63949b75-4fb0-4fd4-b76c-d0c9120b8cd2",
					"code": "validation_error",
					"status": 400,
					"title": "failed to create transaction",
					"detail": "failed to validate request body",
					"source": {
						"field": "amount",
						"message": "field should be non-negative: -9.00"
					}
				}
			]
		}
	*/
}

// InternalServerError: Account check failed, server error
//
// Returns: 500
func (s *transactionTestSuite) TestTransactionAccountCheckFailed() {
	mockDBError := errors.New("some-db-error")
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)
	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("internal_error", string(resBody))
}

// BadRequest: Account check failed, account was not found
//
// Returns: 400
func (s *transactionTestSuite) TestTransactionAccountCheckFailedNotFound() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)
	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("bad_request", string(resBody))
}

// InternalServerError: Operation Type check failed, server error
//
// Return: 500
func (s *transactionTestSuite) TestTransactionOperationTypeFailed() {
	mockDBError := errors.New("some-db-error")
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 4,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)

	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, nil)

	s.mockOperationTypes.EXPECT().Get(gomock.Any(), 4).Return(model.OperationType{}, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("internal_error", string(resBody))
}

// BadRequest: Operation Type ID was invalid, not found
//
// Return: 400
func (s *transactionTestSuite) TestTransactionOperationTypeInvalid() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 14,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)

	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, nil)

	s.mockOperationTypes.EXPECT().Get(gomock.Any(), 14).Return(model.OperationType{}, repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("bad_request", string(resBody))
}

// InternalServerError: DB error, failed to create transaction
//
// Return: 500
func (s *transactionTestSuite) TestTransactionFailed() {
	mockDBError := errors.New("some-db-error")
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/transactions",
		strings.NewReader(
			`{
				"account_uuid": "e2a84838-88de-5fbc-8636-6ef49e26f00a",
				"operation_type_id": 14,
				"amount": 1.1,
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()

	trxUUID := getMockTrxUUID()

	s.mockTrx.EXPECT().CheckIdempotency(gomock.Any(), trxUUID.String()).Return(nil)

	s.mockAccounts.EXPECT().Get(gomock.Any(), "e2a84838-88de-5fbc-8636-6ef49e26f00a").Return(model.Account{}, nil)

	s.mockOperationTypes.EXPECT().Get(gomock.Any(), 14).Return(model.OperationType{}, nil)

	s.mockTrx.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, t model.Transaction) error {
			// validate fields
			if t.UUID != trxUUID ||
				t.AccountUUID.String() != "e2a84838-88de-5fbc-8636-6ef49e26f00a" ||
				t.Amount != 1.1 ||
				t.OperationTypeID != 4 {
				return errors.New("incorrect params")
			}

			return mockDBError
		})

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("internal_error", string(resBody))
}

// Returns mock uuid
func getMockTrxUUID() uuid.UUID {
	return uuid.NewV5(uuid.Nil, "bc1f3956-e92e-4666-a5cd-4cbbd937b17f")
}
