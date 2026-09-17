package auth

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
)

type Routes struct {
	handler  *Handler
	limiters *middleware.Limiters
}

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	limiters *middleware.Limiters,
	requireAuth gin.HandlerFunc,
) {
	g := r.Group("/auth")

	// Public endpoints (rate limited).
	g.POST("/register",
		limiters.AuthRegister.ApplyByIP(),
		h.Register,
	)
	g.POST("/login",
		limiters.AuthLogin.ApplyByIP(),
		h.Login,
	)
	g.POST("/refresh",
		limiters.AuthRefresh.ApplyByIP(),
		h.Refresh,
	)

	// Authenticated endpoints.
	auth := g.Group("")
	auth.Use(requireAuth)
	auth.POST("/logout", h.Logout)
	auth.GET("/me", h.Me)
	auth.POST("/change-password", h.ChangePassword)
}

var _ = time.Second
