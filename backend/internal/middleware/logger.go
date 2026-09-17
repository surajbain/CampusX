package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/campusx/api/pkg/logger"
)

// Logger logs each request with structured fields.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		log := logger.From()
		evt := log.Info()
		if c.Writer.Status() >= 500 {
			evt = log.Error()
		} else if c.Writer.Status() >= 400 {
			evt = log.Warn()
		}

		evt.
			Str("request_id", GetRequestID(c)).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", c.Writer.Status()).
			Int("bytes", c.Writer.Size()).
			Dur("latency_ms", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("http_request")

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error().
					Str("request_id", GetRequestID(c)).
					Err(e.Err).
					Msg("gin_error")
			}
		}
	}
}

// ensure zerolog import is used even if logger package changes
var _ = zerolog.InfoLevel
