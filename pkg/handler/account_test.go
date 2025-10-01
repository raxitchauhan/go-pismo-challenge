package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-pismo-challenge/pkg/model"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/repository/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type accountTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	connector    *Account
	mockAccounts *mocks.MockAccountConnector
	router       *chi.Mux
	recoder      *httptest.ResponseRecorder
}

func TestAccountHnadler(t *testing.T) {
	suite.Run(t, new(accountTestSuite))
}

// Setup test suite
func (s *accountTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockAccounts = mocks.NewMockAccountConnector(s.ctrl)

	s.connector = NewAccountHandler(s.mockAccounts)
	s.recoder = httptest.NewRecorder()
	s.router = chi.NewRouter()

	s.router.Post("/accounts", s.connector.Create)
	s.router.Get("/accounts/{uuid}", s.connector.Get)
}

// Assert expectations
func (s *accountTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Success: An account was created
//
// Return: 201
func (s *accountTestSuite) TestCreateAccountSuccess() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/accounts",
		strings.NewReader(
			`{
				"document_number" : "abc",
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()
	accountUUID := getMockUUID()

	s.mockAccounts.EXPECT().CheckIdempotency(gomock.Any(), accountUUID.String()).Return(nil)
	s.mockAccounts.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, a model.Account) error {
			// validate fields
			if a.UUID != accountUUID || a.DocumentNumber != "abc" {
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

// BadRequest: `document_number` field was not passed in the request body
//
// Returns: 400
func (s *accountTestSuite) TestAccountBadRequestDocumentNumber() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/accounts",
		strings.NewReader(`{ "idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f" }`))
	s.Require().NoError(err)
	defer req.Body.Close()

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("document_number", string(resBody))
}

// BadRequest: `idempotency_key` field was not passed in the request body
//
// Returns: 400
func (s *accountTestSuite) TestAccountBadRequestIdempotency() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/accounts",
		strings.NewReader(`{ "document_number" : "abc" }`))
	s.Require().NoError(err)
	defer req.Body.Close()

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("idempotency_key", string(resBody))
}

// BadRequest: Duplicate request received, idempotency check returns duplicate error
//
// Returns: 400
func (s *accountTestSuite) TestAccountIdempotencyCheckFailed() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/accounts",
		strings.NewReader(
			`{
				"document_number" : "abc",
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()
	accountUUID := getMockUUID()

	s.mockAccounts.EXPECT().CheckIdempotency(gomock.Any(), accountUUID.String()).Return(repository.ErrDuplicate)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusBadRequest, s.recoder.Code)
	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("duplicate", string(resBody))
}

// InternalServerError: Account creation failed at database
//
// Return: 500
func (s *accountTestSuite) TestAccountCreateFailed() {
	mockDBError := errors.New("some db error")
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodPost, "/accounts",
		strings.NewReader(
			`{
				"document_number" : "abc",
				"idempotency_key": "bc1f3956-e92e-4666-a5cd-4cbbd937b17f"
			}`))
	s.Require().NoError(err)
	defer req.Body.Close()
	accountUUID := getMockUUID()

	s.mockAccounts.EXPECT().CheckIdempotency(gomock.Any(), accountUUID.String()).Return(nil)
	s.mockAccounts.EXPECT().Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, a model.Account) error {
			// validate fields
			if a.UUID != accountUUID || a.DocumentNumber != "abc" {
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

// Success: Get account by UUID
//
// Return: 200
func (s *accountTestSuite) TestGetAccountSuccess() {
	accountUUID := getMockUUID()

	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/accounts/"+accountUUID.String(), nil)
	s.Require().NoError(err)

	expected := model.Account{
		UUID:           accountUUID,
		DocumentNumber: "sample-number",
		CreatedAt:      time.Now(),
	}

	s.mockAccounts.EXPECT().Get(gomock.Any(), accountUUID.String()).Return(expected, nil)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusOK, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	expectedJson, err := json.Marshal(expected)
	s.NoError(err)
	s.JSONEq(string(expectedJson), string(resBody))
}

// InternalServerError: Failed to get account by UUID
//
// Return: 500
func (s *accountTestSuite) TestGetAccountFailure() {
	accountUUID := getMockUUID()
	mockDBError := errors.New("some-db-error")

	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/accounts/"+accountUUID.String(), nil)
	s.Require().NoError(err)

	s.mockAccounts.EXPECT().Get(gomock.Any(), accountUUID.String()).Return(model.Account{}, mockDBError)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusInternalServerError, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("internal_error", string(resBody))
}

// NotFound: Given account UUID was not found
//
// Return: 404
func (s *accountTestSuite) TestGetAccountFailureNotFound() {
	accountUUID := getMockUUID()

	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/accounts/"+accountUUID.String(), nil)
	s.Require().NoError(err)

	s.mockAccounts.EXPECT().Get(gomock.Any(), accountUUID.String()).Return(model.Account{}, repository.ErrNoRows)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)

	resBody, err := io.ReadAll(s.recoder.Body)
	s.NoError(err)
	s.Regexp("not_found", string(resBody))
}

// NotFound: Get account, incorrect URL param
//
// Return: 404
func (s *accountTestSuite) TestGetAccountFailureIncorrectURLNotFound() {
	req, err := http.NewRequestWithContext(s.T().Context(), http.MethodGet, "/accounts/", nil)
	s.Require().NoError(err)

	s.router.ServeHTTP(s.recoder, req)

	s.Equal(http.StatusNotFound, s.recoder.Code)
}

// Returns mock uuid
func getMockUUID() uuid.UUID {
	return uuid.NewV5(uuid.Nil, "bc1f3956-e92e-4666-a5cd-4cbbd937b17f")
}
