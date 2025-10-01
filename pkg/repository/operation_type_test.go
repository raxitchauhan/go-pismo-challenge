package repository

import (
	"context"
	"errors"
	"go-pismo-challenge/pkg/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-jose/go-jose/v4/testutils/require"
	"github.com/stretchr/testify/suite"
)

type operationTypeSuite struct {
	suite.Suite
	repo OperationTypeConnector
	db   sqlmock.Sqlmock
}

func TestOperationType(t *testing.T) {
	suite.Run(t, new(operationTypeSuite))
}

func (s *operationTypeSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.repo = NewOperationTypeRepo(db)
	s.db = mock
}

func (s *operationTypeSuite) TearDownTest() {
	s.NoError(s.db.ExpectationsWereMet())
}

func (s *operationTypeSuite) TestGetOperationTypeSuccess() {
	ctx := context.Background()

	expected := model.OperationType{
		OperationTypeID: 1,
		IsCredit:        true,
	}
	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT operation_type_id, is_credit FROM transactions.operation_types WHERE operation_type_id = $1;`)).
		WithArgs(expected.OperationTypeID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"operation_type_id",
					"is_credit",
				}).
				AddRow(
					expected.OperationTypeID,
					expected.IsCredit,
				))

	got, err := s.repo.Get(ctx, expected.OperationTypeID)
	s.NoError(err)
	s.Equal(got, expected)
}

func (s *operationTypeSuite) TestGetOperationTypeError() {
	ctx := context.Background()
	mockError := errors.New("db error")

	s.db.ExpectQuery(regexp.QuoteMeta(`SELECT operation_type_id, is_credit FROM transactions.operation_types WHERE operation_type_id = $1;`)).
		WithArgs(1).
		WillReturnError(mockError)

	got, err := s.repo.Get(ctx, 1)
	s.Error(err)
	s.True(errors.Is(err, mockError))
	s.Equal(got, model.OperationType{})
}
