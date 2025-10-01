package handler

import (
	"encoding/json"
	"errors"
	"go-pismo-challenge/pkg/model"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/util"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

type Account struct {
	accountRepo repository.AccountConnector
}

func NewAccountHandler(a repository.AccountConnector) *Account {
	return &Account{
		accountRepo: a,
	}
}

// @Summary Returns Generated Account UUID
// @Description Create an Account.
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body model.AccountRequest true "Add account request"
// @Success 201 {object} model.AccountResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /accounts [post]
func (a *Account) Create(w http.ResponseWriter, r *http.Request) {
	var req model.AccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to decode request body",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}

	vErr := req.Validate()
	if len(vErr) > 0 {
		err := util.WriteJSONError(w,
			http.StatusBadRequest,
			util.ErrorDescription{
				Code:    validationError,
				Status:  http.StatusBadRequest,
				Title:   failedToCreateAccount,
				Details: "failed to validate request body",
			},
			vErr...)
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}

	accountUUID := uuid.NewV5(uuid.Nil, req.IdempotencyKey)
	if err := a.accountRepo.CheckIdempotency(r.Context(), accountUUID.String()); err != nil {
		status := http.StatusInternalServerError
		code := internalError
		if errors.Is(err, repository.ErrDuplicate) {
			status = http.StatusBadRequest
			code = badRequest
		}

		err = util.WriteJSONError(w, status, util.ErrorDescription{
			Status:  status,
			Code:    code,
			Title:   failedToCreateAccount,
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}
	account := model.Account{
		UUID:           accountUUID,
		DocumentNumber: req.DocumentNumber,
		CreatedAt:      time.Now(),
	}

	if err := a.accountRepo.Create(r.Context(), account); err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   failedToCreateAccount,
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}

	if err := util.WriteJSON(w, http.StatusCreated, model.AccountResponse{UUID: accountUUID.String()}); err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to write response",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}
}

// @Summary Returns an Account
// @Description Get an Account by UUID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param   uuid path string true "Account UUID"
// @Success 200 {object} model.Account
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /accounts/{uuid} [get]
func (a *Account) Get(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		err := util.WriteJSONError(w, http.StatusNotFound, util.ErrorDescription{
			Status:  http.StatusNotFound,
			Code:    notFound,
			Title:   accountNotFound,
			Details: "path param 'uuid' cannot be empty",
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}
	}
	account, err := a.accountRepo.Get(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			err = util.WriteJSONError(w, http.StatusNotFound, util.ErrorDescription{
				Status:  http.StatusNotFound,
				Code:    notFound,
				Title:   accountNotFound,
				Details: err.Error(),
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to write error resposne")
			}

			return
		}

		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to get account",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}

	if err := util.WriteJSON(w, http.StatusOK, account); err != nil {
		err = util.WriteJSONError(w, http.StatusInternalServerError, util.ErrorDescription{
			Status:  http.StatusInternalServerError,
			Code:    internalError,
			Title:   "failed to write response",
			Details: err.Error(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to write error resposne")
		}

		return
	}
}
