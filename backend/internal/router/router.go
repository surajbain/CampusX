package router

import (
	"os"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/auth"
	"github.com/campusx/api/internal/college"
	"github.com/campusx/api/internal/event"
	"github.com/campusx/api/internal/health"
	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/payment"
	"github.com/campusx/api/internal/registration"
	"github.com/campusx/api/internal/ticket"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/jwt"
	"github.com/campusx/api/pkg/response"
	"github.com/campusx/api/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Options struct {
	Config    *config.Config
	DB        *gorm.DB
	Redis     *redis.Client
	JWT       *jwt.Issuer
	UseTestRL bool
}

func Build(opts Options) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(opts.Config.HTTP.CORSOrigins))

	health.NewHandler(opts.DB, opts.Redis, opts.Config.App.Env).RegisterRoutes(r)

	auditLog := audit.NewLogger(opts.DB)

	// Auth
	authRepo := auth.NewRepo(opts.DB)
	authSvc := auth.NewService(authRepo, opts.JWT, session.NewStore(opts.Redis), auditLog)
	authHandler := auth.NewHandler(authSvc)
	limiters := buildLimiters(opts.Redis, opts.UseTestRL)

	// College
	colRepo := college.NewRepo(opts.DB)
	colSvc := college.NewService(colRepo, auditLog)
	colHandler := college.NewHandler(colSvc)

	// Event
	eventRepo := event.NewRepo(opts.DB)
	eventSvc := event.NewService(eventRepo, auditLog)
	eventHandler := event.NewHandler(eventSvc)

	// Registration
	regRepo := registration.NewRepo(opts.DB)
	regSvc := registration.NewService(regRepo, auditLog)
	regHandler := registration.NewHandler(regSvc)

	// Payment
	payProvider := buildPaymentProvider()
	payRepo := payment.NewRepo(opts.DB)
	paySvc := payment.NewService(payRepo, payProvider, auditLog)
	payHandler := payment.NewHandler(paySvc)

	// Routes
	v1 := r.Group("/api/v1")
	requireAuth := middleware.RequireAuth(opts.JWT)

	auth.RegisterRoutes(v1, authHandler, limiters, requireAuth)
	college.RegisterRoutes(v1, colHandler, requireAuth)
	event.RegisterRoutes(v1, eventHandler, requireAuth)
	registration.RegisterRoutes(v1, regHandler, requireAuth)
	payment.RegisterRoutes(v1, payHandler, requireAuth)
		// ---- Ticket management ----
	qrManager := ticket.NewQRManager(opts.Config.Security.QRHMACSecret)
	ticketRepo := ticket.NewRepo(opts.DB)
	ticketSvc := ticket.NewService(ticketRepo, qrManager, auditLog)
	ticketHandler := ticket.NewHandler(ticketSvc)
	ticket.RegisterRoutes(v1, ticketHandler, requireAuth)

	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, apperror.NotFound("route not found"))
	})

	return r
}

func buildPaymentProvider() payment.Provider {
	provider := os.Getenv("PAYMENT_PROVIDER")
	if provider == "razorpay" {
		return payment.NewRazorpayProvider(
			os.Getenv("RAZORPAY_KEY_ID"),
			os.Getenv("RAZORPAY_KEY_SECRET"),
			os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		)
	}
	return payment.NewMockProvider("mock_secret_dev_only")
}

func buildLimiters(rdb *redis.Client, useTest bool) *middleware.Limiters {
	if useTest {
		return middleware.NewTestLimiters(rdb)
	}
	return middleware.NewLimiters(rdb)
}