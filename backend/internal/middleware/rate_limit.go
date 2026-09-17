package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/response"
)

// RateLimitConfig describes a single limiter.
type RateLimitConfig struct {
	Name   string        // bucket namespace
	Limit  int           // max requests
	Window time.Duration // in this window
}

// Limiter is a Redis sliding-window limiter.
type Limiter struct {
	rdb    *redis.Client
	config RateLimitConfig
}

// Limiters bundles all app rate limiters.
type Limiters struct {
	AuthRegister *Limiter
	AuthLogin    *Limiter
	AuthRefresh  *Limiter
}

// NewLimiters returns production rate limiters.
func NewLimiters(rdb *redis.Client) *Limiters {
	return &Limiters{
		AuthRegister: NewLimiter(rdb, RateLimitConfig{Name: "auth_register", Limit: 3, Window: time.Hour}),
		AuthLogin:    NewLimiter(rdb, RateLimitConfig{Name: "auth_login", Limit: 5, Window: 15 * time.Minute}),
		AuthRefresh:  NewLimiter(rdb, RateLimitConfig{Name: "auth_refresh", Limit: 20, Window: time.Hour}),
	}
}

// NewTestLimiters returns effectively-infinite limiters for use in tests only.
func NewTestLimiters(rdb *redis.Client) *Limiters {
	return &Limiters{
		AuthRegister: NewLimiter(rdb, RateLimitConfig{Name: "test_register", Limit: 100000, Window: time.Hour}),
		AuthLogin:    NewLimiter(rdb, RateLimitConfig{Name: "test_login", Limit: 100000, Window: time.Hour}),
		AuthRefresh:  NewLimiter(rdb, RateLimitConfig{Name: "test_refresh", Limit: 100000, Window: time.Hour}),
	}
}

func NewLimiter(rdb *redis.Client, cfg RateLimitConfig) *Limiter {
	return &Limiter{rdb: rdb, config: cfg}
}

// ApplyByIP rate-limits by client IP. Use as gin middleware.
func (l *Limiter) ApplyByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rl:%s:ip:%s", l.config.Name, c.ClientIP())
		ok, remaining, err := l.allow(c.Request.Context(), key)
		if err != nil {
			// Fail-open: if Redis is down, don't block traffic.
			c.Next()
			return
		}
		setRateHeaders(c, l.config.Limit, remaining, l.config.Window)
		if !ok {
			response.Fail(c, apperror.New(apperror.CodeTooManyReqs, "too many requests, please try again later"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// ApplyByKey rate-limits by an arbitrary key extractor.
func (l *Limiter) ApplyByKey(keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := keyFn(c)
		if k == "" {
			c.Next()
			return
		}
		key := fmt.Sprintf("rl:%s:%s", l.config.Name, k)
		ok, remaining, err := l.allow(c.Request.Context(), key)
		if err != nil {
			c.Next()
			return
		}
		setRateHeaders(c, l.config.Limit, remaining, l.config.Window)
		if !ok {
			response.Fail(c, apperror.New(apperror.CodeTooManyReqs, "too many requests, please try again later"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// allow increments the counter and returns (allowed, remaining, err).
func (l *Limiter) allow(ctx context.Context, key string) (bool, int, error) {
	pipe := l.rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, l.config.Window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}
	count := int(incr.Val())
	if count > l.config.Limit {
		return false, 0, nil
	}
	return true, l.config.Limit - count, nil
}

func setRateHeaders(c *gin.Context, limit, remaining int, window time.Duration) {
	c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(window).Unix())))
}
