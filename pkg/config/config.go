package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME,required"`

	DatabaseHost           string `env:"DATABASE_HOST,required"`
	DatabaseName           string `env:"DATABASE_NAME,required"`
	DatabaseUserName       string `env:"DATABASE_USERNAME,required"`
	DatabasePassword       string `env:"DATABASE_PASSWORD,required"`
	DatabaseSSLMode        string `env:"DATABASE_SSL_MODE,required"`
	DatabaseSSLRootCert    string `env:"DATABASE_SSL_ROOT_CERT,required"`
	DatabaseMaxOpenConns   int    `env:"DATABASE_MAX_OPEN_CONNS,required"`
	DatabaseMigrationTable string `env:"DATABASE_MIGRATION_TABLE,required"`
	DatabaseMinVersion     int    `env:"DATABASE_MIN_VERSION,required"`
}

func LoadConfig() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(fmt.Errorf("failed to load config: %w", err))
	}

	return cfg
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s&sslrootcert=%s&application_name=%s",
		c.DatabaseUserName,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabaseName,
		c.DatabaseSSLMode,
		c.DatabaseSSLRootCert,
		c.ServiceName,
	)
}
