package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func NewConnection(dsn string, maxConnections int) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database connection")
	}

	db.SetMaxOpenConns(maxConnections)

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Check database migration version against the minimum required version
func MigrationVersionCheck(ctx context.Context, db *sql.DB, tableName string, minVersion int) error {
	schemaVersionQuery := fmt.Sprintf(
		`SELECT version_id FROM %s WHERE is_applied = true AND version_id >= %d;`,
		tableName,
		minVersion)

	var res any
	if err := db.QueryRowContext(ctx, schemaVersionQuery).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: expected %s version to be >= %d", errors.New("version incompatible"), tableName, minVersion)
		}

		return fmt.Errorf("error while checking version: %w", err)
	}
	log.Info().Msg("database migration version check successful")

	return nil
}
