package payment

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	requireAuth gin.HandlerFunc,
) {
	g := r.Group("/payments")
	g.Use(requireAuth)
	g.POST("/create", h.Create)
	g.POST("/verify", h.Verify)
	g.GET("/:id", h.Get)
}