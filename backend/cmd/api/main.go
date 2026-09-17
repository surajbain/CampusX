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

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/auth"
	"github.com/campusx/api/internal/college"
	"github.com/campusx/api/internal/event"
	"github.com/campusx/api/internal/health"
	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/migrations"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/db"
	"github.com/campusx/api/pkg/jwt"
	"github.com/campusx/api/pkg/logger"
	redisPkg "github.com/campusx/api/pkg/redis"
	"github.com/campusx/api/pkg/response"
	"github.com/campusx/api/pkg/session"
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

	// ------- MIGRATIONS -------
	if os.Getenv("RUN_MIGRATIONS") == "true" {
		log.Info().Msg("running migrations")
		if err := migrations.Run(ctx, gdb, &log); err != nil {
			log.Fatal().Err(err).Msg("migrations failed")
		}
	}
	// --------------------------

	rdb, err := redisPkg.Connect(ctx, cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("redis connect failed")
	}
	defer func() { _ = redisPkg.Close(rdb) }()

	// ------- AUTH WIRING -------
	jwtIssuer := jwt.NewIssuer(jwt.Config{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessTTL:     cfg.JWT.AccessTTL,
		RefreshTTL:    cfg.JWT.RefreshTTL,
		Issuer:        cfg.JWT.Issuer,
	})
	sessionStore := session.NewStore(rdb)
	auditLog := audit.NewLogger(gdb)

	authRepo := auth.NewRepo(gdb)
	authSvc := auth.NewService(authRepo, jwtIssuer, sessionStore, auditLog)
	authHandler := auth.NewHandler(authSvc)

	limiters := middleware.NewLimiters(rdb)
	// ---------------------------

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.HTTP.CORSOrigins))

	health.NewHandler(gdb, rdb, cfg.App.Env).RegisterRoutes(router)

	// ------- API v1 ROUTES -------
	v1 := router.Group("/api/v1")
	auth.RegisterRoutes(v1, authHandler, limiters, middleware.RequireAuth(jwtIssuer))
	// -----------------------------

	// ---- College management ----
	collegeRepo := college.NewRepo(gdb)
	collegeSvc := college.NewService(collegeRepo, auditLog)
	collegeHandler := college.NewHandler(collegeSvc)
	college.RegisterRoutes(v1, collegeHandler, middleware.RequireAuth(jwtIssuer))
	router.NoRoute(func(c *gin.Context) {
		response.Fail(c, apperror.NotFound("route not found"))
	})

	// ---- Event management ----
	eventRepo := event.NewRepo(gdb)
	eventSvc := event.NewService(eventRepo, auditLog)
	eventHandler := event.NewHandler(eventSvc)
	event.RegisterRoutes(v1, eventHandler, middleware.RequireAuth(jwtIssuer))

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
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
