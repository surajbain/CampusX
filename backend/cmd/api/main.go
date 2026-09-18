package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/migrations"
	"github.com/campusx/api/internal/router"
	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/db"
	"github.com/campusx/api/pkg/jwt"
	"github.com/campusx/api/pkg/logger"
	redisPkg "github.com/campusx/api/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	log := logger.Init(cfg.App.Env, cfg.App.Name)
	log.Info().Msg("starting CampusX API")

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gdb, err := db.Connect(ctx, cfg.DB, cfg.App.Env)
	if err != nil {
		log.Fatal().Err(err).Msg("db connect failed")
	}
	defer func() { _ = db.Close(gdb) }()

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		log.Info().Msg("running migrations")
		if err := migrations.Run(ctx, gdb, &log); err != nil {
			log.Fatal().Err(err).Msg("migrations failed")
		}
	}

	rdb, err := redisPkg.Connect(ctx, cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("redis connect failed")
	}
	defer func() { _ = redisPkg.Close(rdb) }()

	jwtIssuer := jwt.NewIssuer(jwt.Config{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessTTL:     cfg.JWT.AccessTTL,
		RefreshTTL:    cfg.JWT.RefreshTTL,
		Issuer:        cfg.JWT.Issuer,
	})

	// Build router using shared function
	routerEngine := router.Build(router.Options{
		Config: cfg,
		DB:     gdb,
		Redis:  rdb,
		JWT:    jwtIssuer,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      routerEngine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("http server listening")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("http server error")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	} else {
		log.Info().Msg("server stopped cleanly")
	}
}