package main

import (
	"context"
	"go-pismo-challenge/pkg/config"
	"go-pismo-challenge/pkg/migrator"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()

	// run migration
	if err := migrator.Run(context.Background(), cfg.DSN(), migrator.FS); err != nil {
		log.Fatal().Err(err).Msg("failed to apply migration")
	}

	log.Info().Msg("migration applied successfully")
}
