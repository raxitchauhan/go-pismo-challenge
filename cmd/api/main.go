package main

import (
	"context"
	"fmt"
	"go-pismo-challenge/pkg/config"
	"go-pismo-challenge/pkg/database"
	"go-pismo-challenge/pkg/handler"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/server"

	"github.com/rs/zerolog/log"
)

type Service struct {
	accountHandler *handler.Account
	trxHandler     *handler.Transaction
}

// @title Pismo API
// @version 1.0
// @description This is an API documentation for Pismo challenge
// @host localhost:3000
// @BasePath /v1
func main() {
	ctx := context.Background()
	done := NewService().Run(ctx)

	<-done
}

func NewService() *Service {
	cfg := config.LoadConfig()

	db, err := database.NewConnection(cfg.DSN(), cfg.DatabaseMaxOpenConns)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to establish db connection: %w", err))
	}

	accountRepo := repository.NewAccountRepo(db)
	accountHandler := handler.NewAccountHandler(accountRepo)

	trxRepo := repository.NewTransactionRepo(db)
	operationTypeRepo := repository.NewOperationTypeRepo(db)
	trxHandler := handler.NewTransactionHandler(trxRepo, accountRepo, operationTypeRepo)

	return &Service{
		accountHandler: accountHandler,
		trxHandler:     trxHandler,
	}
}

func (s *Service) Run(ctx context.Context) <-chan struct{} {
	if err := server.Start(s.accountHandler, s.trxHandler); err != nil {
		log.Fatal().Err(err)
	}

	return ctx.Done()
}
