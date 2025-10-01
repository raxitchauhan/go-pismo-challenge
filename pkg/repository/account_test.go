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

type accountSuite struct {
	suite.Suite
	repo AccountConnector
	db   sqlmock.Sqlmock
}

func TestAccount(t *testing.T) {
	suite.Run(t, new(accountSuite))
}

func (s *accountSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.repo = NewAccountRepo(db)
	s.db = mock
}

func (s *accountSuite) TearDownTest() {
	s.NoError(s.db.ExpectationsWereMet())
}

func (s *accountSuite) TestCreateSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	request := model.Account{
		UUID:           mockUUID,
		DocumentNumber: "doc",
		CreatedAt:      now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts.account (uuid, document_number, created_at) 
				values ($1, $2, $3) ON CONFLICT (uuid) DO NOTHING`)).
		WithArgs(
			request.UUID.String(),
			request.DocumentNumber,
			request.CreatedAt,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Create(ctx, request)
	s.NoError(err)
}

func (s *accountSuite) TestCreateError() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	mockError := errors.New("db error")

	request := model.Account{
		UUID:           mockUUID,
		DocumentNumber: "doc",
		CreatedAt:      now,
	}

	s.db.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts.account (uuid, document_number, created_at) 
				values ($1, $2, $3) ON CONFLICT (uuid) DO NOTHING`)).
		WithArgs(
			request.UUID.String(),
			request.DocumentNumber,
			request.CreatedAt,
		).WillReturnError(mockError)

	err := s.repo.Create(ctx, request)
	s.Error(err)
	s.True(errors.Is(err, mockError))
}

func (s *accountSuite) TestGetAccountSuccess() {
	ctx := context.Background()
	now := time.Now()
	mockUUID := uuid.NewV5(uuid.Nil, "")

	expected := model.Account{
		UUID:           mockUUID,
		DocumentNumber: "abc",
		CreatedAt:      now,
	}
	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT uuid, document_number, created_at FROM accounts.account where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"uuid",
					"document_number",
					"created_at",
				}).
				AddRow(
					expected.UUID.String(),
					expected.DocumentNumber,
					expected.CreatedAt,
				))

	got, err := s.repo.Get(ctx, mockUUID.String())
	s.NoError(err)
	s.Equal(got, expected)
}

func (s *accountSuite) TestGetAccountError() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT uuid, document_number, created_at FROM accounts.account where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(mockError)

	got, err := s.repo.Get(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, mockError))
	s.Equal(got, model.Account{})
}

func (s *accountSuite) TestCheckIdempotencySuccess() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM accounts.account where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.NoError(err)
}

func (s *accountSuite) TestCheckIdempotencyDuplicate() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM accounts.account where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, ErrDuplicate))
}

func (s *accountSuite) TestCheckIdempotencyError() {
	ctx := context.Background()
	mockUUID := uuid.NewV5(uuid.Nil, "")
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(COUNT(*), 0) as count FROM accounts.account where uuid = $1;`)).
		WithArgs(mockUUID.String()).
		WillReturnError(mockError)

	err := s.repo.CheckIdempotency(ctx, mockUUID.String())
	s.Error(err)
	s.True(errors.Is(err, mockError))
}
