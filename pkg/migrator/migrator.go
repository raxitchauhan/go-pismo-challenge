package migrator

import (
	"context"
	"fmt"
	"go-pismo-challenge/pkg/database"
	"io/fs"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Run(ctx context.Context, dsn string, migrationFS fs.FS) error {
	db, err := database.NewConnection(dsn, 1)
	if err != nil {
		return fmt.Errorf("failed to establish db connection: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf(
			"unable to apply sql dialect to migrations package %q: %w",
			"postgres",
			err,
		)
	}

	goose.SetTableName("db_version_pismo")
	goose.SetLogger(&logger{})
	goose.SetVerbose(true)

	goose.SetBaseFS(migrationFS)

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("unable to apply database migrations: %w", err)
	}

	return nil
}
