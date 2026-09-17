package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/logger"
	"github.com/campusx/api/pkg/response"
)

// Recovery catches panics, logs stack, returns a clean 500.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.From().Error().
					Interface("panic", r).
					Str("request_id", GetRequestID(c)).
					Str("path", c.Request.URL.Path).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorEnvelope{
					Success: false,
					Error: response.ErrorBody{
						Code:    string(apperror.CodeInternal),
						Message: "internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}
