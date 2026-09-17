package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init sets up global zerolog logger. Call once at startup.
func Init(env, appName string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var logger zerolog.Logger
	if env == "development" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	logger = logger.With().
		Str("service", appName).
		Str("env", env).
		Logger()

	log.Logger = logger
	return logger
}

// From returns the global logger (safe to call anywhere).
func From() *zerolog.Logger {
	return &log.Logger
}
