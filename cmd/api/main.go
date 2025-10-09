package main

import (
	"context"
	"fmt"
	"go-pismo-challenge/pkg/config"
	"go-pismo-challenge/pkg/database"
	"go-pismo-challenge/pkg/handler"
	"go-pismo-challenge/pkg/repository"
	"go-pismo-challenge/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
// @BasePath /api/v1
func main() {
	// create a root context that can be cancelled
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	NewService(ctx).Run(ctx)
}

func NewService(ctx context.Context) *Service {
	cfg := config.LoadConfig()

	db, err := database.NewConnection(cfg.DSN(), cfg.DatabaseMaxOpenConns)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to establish database connection: %w", err))
	}

	if err := database.MigrationVersionCheck(ctx, db, cfg.DatabaseMigrationTable, cfg.DatabaseMinVersion); err != nil {
		log.Fatal().Err(fmt.Errorf("failed while checking database migration version: %w", err))
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

func (s *Service) Run(ctx context.Context) {
	webServer := server.NewServer(s.accountHandler, s.trxHandler)
	go func() {
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err)
		}
	}()

	defer func() {
		// create a new context for shutdown timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := webServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("server shutdown failed")
		}
	}()

	<-ctx.Done()
}
