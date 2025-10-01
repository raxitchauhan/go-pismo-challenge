package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-pismo-challenge/pkg/model"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/util"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type Transaction struct {
	trxRepo           repository.TransactionConnector
	accountRepo       repository.AccountConnector
	operationTypeRepo repository.OperationTypeConnector
}

func NewTransactionHandler(
	t repository.TransactionConnector,
	a repository.AccountConnector,
	o repository.OperationTypeConnector,
) *Transaction {
	return &Transaction{
		trxRepo:           t,
		accountRepo:       a,
		operationTypeRepo: o,
	}
}

// @Summary Returns Generated Transaction UUID
// @Description Create a transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body model.TransactionRequest true "Add transaction request"
// @Param operation_type_id query string true "User status" Enum(active, inactive, suspended)
// @Success 201 {object} model.TransactionResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /transactions [post]
func (t *Transaction) Create(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to decode request body",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	if vErr := req.Validate(); len(vErr) > 0 {
		err = util.WriteJSONError(w,
			http.StatusBadRequest,
			util.ErrorDescription{
				Code:    validationError,
				Status:  http.StatusBadRequest,
				Title:   failedToCreateTrx,
				Details: "failed to validate request body",
			},
			vErr...)
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	trxUUID := uuid.NewV5(uuid.Nil, req.IdempotencyKey)
	if err := t.trxRepo.CheckIdempotency(r.Context(), trxUUID.String()); err != nil {
		status := http.StatusInternalServerError
		code := internalError
		if errors.Is(err, repository.ErrDuplicate) {
			status = http.StatusBadRequest
			code = badRequest
		}

		err = util.WriteJSONError(w, status, util.ErrorDescription{
			Status:  status,
			Code:    code,
			Title:   failedToCreateTrx,
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	// validate account
	if _, err := t.accountRepo.Get(r.Context(), req.AccountUUID); err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			err = util.WriteJSONError(w, http.StatusBadRequest, util.ErrorDescription{
				Status:  http.StatusBadRequest,
				Code:    badRequest,
				Title:   "failed to validate account",
				Details: fmt.Sprintf("account not found for account_uuid: '%s'", req.AccountUUID),
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to write error response")
			}

			return
		}

		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to create transaction",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	// validate operation type
	operationType, err := t.operationTypeRepo.Get(r.Context(), req.OperationTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			err = util.WriteJSONError(w, http.StatusBadRequest, util.ErrorDescription{
				Status:  http.StatusBadRequest,
				Code:    badRequest,
				Title:   "failed to validate operation_type_id",
				Details: fmt.Sprintf("invalid operation type: %d", req.OperationTypeID),
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to write error response")
			}

			return
		}

		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to create transaction",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	// ignoring the error as we have already validated the field
	accoutUUID, _ := uuid.FromString(req.AccountUUID)
	trx := model.Transaction{
		UUID:            trxUUID,
		AccountUUID:     accoutUUID,
		OperationTypeID: req.OperationTypeID,
		Amount:          resolveAmount(req.Amount, operationType.IsCredit),
		EventDate:       time.Now(),
	}

	if err := t.trxRepo.Create(r.Context(), trx); err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   failedToCreateTrx,
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, model.TransactionResponse{UUID: trxUUID.String()}); err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to write response",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error response")
		}

		return
	}
}

// returns negative amount for debit operation type, else positive
func resolveAmount(amount float64, isCredit bool) float64 {
	if isCredit {
		return amount
	}

	return amount * (-1)
}
