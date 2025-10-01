package repository

import (
	"context"
	"errors"
	"go-pismo-challenge/pkg/model"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jose/go-jose/v4/testutils/require"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
)

type transactionSuite struct {
	suite.Suite
	repo TransactionConnector
	db   sqlmock.Sqlmock
}

func TestTransaction(t *testing.T) {
	suite.Run(t, new(transactionSuite))
}

func (s *transactionSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.repo = NewTransactionRepo(db)
	s.db = mock
}

func (s *transactionSuite) TearDownTest() {
	s.NoError(s.db.ExpectationsWereMet())
}

func (s *transactionSuite) TestCreateSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	request := model.Transaction{
		UUID:            mockUUID,
		AccountUUID:     uuid.NewV5(mockUUID, "account"),
		OperationTypeID: 2,
		Amount:          -11.9,
		EventDate:       now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions.transaction (uuid, account_uuid, operation_type_id, amount, event_date) 
					values ($1, $2, $3, $4, $5) ON CONFLICT (uuid) DO NOTHING;`)).
		WithArgs(
			request.UUID.String(),
			request.AccountUUID.String(),
			request.OperationTypeID,
			request.Amount,
			request.EventDate,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Create(ctx, request)
	s.NoError(err)
}

func (s *transactionSuite) TestCreateError() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	mockError := errors.New("db error")

	request := model.Transaction{
		UUID:            mockUUID,
		AccountUUID:     uuid.NewV5(mockUUID, "account"),
		OperationTypeID: 2,
		Amount:          -11.9,
		EventDate:       now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO transactions.transaction (uuid, account_uuid, operation_type_id, amount, event_date) 
					values ($1, $2, $3, $4, $5) ON CONFLICT (uuid) DO NOTHING;`)).
		WithArgs(
			request.UUID.String(),
			request.AccountUUID.String(),
			request.OperationTypeID,
			request.Amount,
			request.EventDate,
		).WillReturnError(mockError)

	err := s.repo.Create(ctx, request)
	s.Error(err)
	s.True(errors.Is(err, mockError))
}

func (s *transactionSuite) TestCheckIdempotencySuccess() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM transactions.transaction where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.NoError(err)
}

func (s *transactionSuite) TestCheckIdempotencyDuplicate() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM transactions.transaction where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, ErrDuplicate))
}

func (s *transactionSuite) TestCheckIdempotencyError() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM transactions.transaction where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(mockError)

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, mockError))
}
