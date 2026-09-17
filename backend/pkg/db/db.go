package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/logger"
)

// Connect opens a GORM connection with retry-on-startup.
// It waits for Postgres to be ready (important for docker-compose).
func Connect(ctx context.Context, cfg config.DBConfig, env string) (*gorm.DB, error) {
	log := logger.From()

	logLevel := gormlogger.Warn
	if env == "development" {
		logLevel = gormlogger.Info
	}

	gormCfg := &gorm.Config{
		Logger:                 gormlogger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	var gdb *gorm.DB
	var err error
	backoff := 500 * time.Millisecond
	maxAttempts := 10

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		gdb, err = gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
		if err == nil {
			sqlDB, sErr := gdb.DB()
			if sErr == nil {
				if pErr := sqlDB.PingContext(ctx); pErr == nil {
					// Configure pool.
					sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
					sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
					sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
					sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

					log.Info().
						Str("host", cfg.Host).
						Str("db", cfg.Name).
						Int("max_open", cfg.MaxOpenConns).
						Msg("database connected")
					return gdb, nil
				}
			}
		}

		log.Warn().
			Err(err).
			Int("attempt", attempt).
			Dur("backoff", backoff).
			Msg("database not ready, retrying")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
		backoff *= 2
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}
	}

	return nil, fmt.Errorf("database connection failed after %d attempts: %w", maxAttempts, err)
}

// Health pings the DB for /ready checks.
func Health(ctx context.Context, gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close closes the underlying connection pool.
func Close(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
