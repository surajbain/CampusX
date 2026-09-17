package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/campusx/api/pkg/db"
	redisPkg "github.com/campusx/api/pkg/redis"
	"github.com/campusx/api/pkg/response"
)

type Handler struct {
	DB    *gorm.DB
	Redis *redis.Client
	Env   string
}

func NewHandler(gdb *gorm.DB, rdb *redis.Client, env string) *Handler {
	return &Handler{DB: gdb, Redis: rdb, Env: env}
}

// Live — process is up. Used by orchestrators to decide restart.
func (h *Handler) Live(c *gin.Context) {
	response.OK(c, gin.H{
		"status": "ok",
		"env":    h.Env,
	})
}

// Ready — dependencies are reachable. Used for traffic routing.
func (h *Handler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{}
	allOK := true

	if err := db.Health(ctx, h.DB); err != nil {
		checks["database"] = gin.H{"status": "down", "error": err.Error()}
		allOK = false
	} else {
		checks["database"] = gin.H{"status": "up"}
	}

	if err := redisPkg.Health(ctx, h.Redis); err != nil {
		checks["redis"] = gin.H{"status": "down", "error": err.Error()}
		allOK = false
	} else {
		checks["redis"] = gin.H{"status": "up"}
	}

	if !allOK {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"status":  "not_ready",
			"checks":  checks,
		})
		return
	}

	response.OK(c, gin.H{
		"status": "ready",
		"checks": checks,
	})
}

// RegisterRoutes wires health endpoints.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Live)
	r.GET("/ready", h.Ready)
}
