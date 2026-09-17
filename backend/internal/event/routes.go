package event

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/models"
)

// RegisterRoutes wires all event endpoints.
func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	requireAuth gin.HandlerFunc,
) {
	g := r.Group("/events")

	// ---- Public (no auth) ----
	g.GET("", h.List)
	g.GET("/featured", h.Featured)
	g.GET("/trending", h.Trending)
	g.GET("/categories", h.Categories)
	g.GET("/slug/:slug", h.GetBySlug)
	g.GET("/:id", h.GetByID)

	// ---- Authenticated ----
	auth := g.Group("")
	auth.Use(requireAuth)

	// My events (organizer's own, including drafts).
	auth.GET("/my", h.ListMine)

	// Create: ORGANIZER / COLLEGE_ADMIN / SUPER_ADMIN
	creators := auth.Group("")
	creators.Use(middleware.RequireRole(
		models.RoleOrganizer,
		models.RoleCollegeAdmin,
		models.RoleSuperAdmin,
	))
	creators.POST("", h.Create)

	// Update / lifecycle / delete: owner checked inside service.
	editors := auth.Group("")
	editors.Use(middleware.RequireRole(
		models.RoleOrganizer,
		models.RoleCollegeAdmin,
		models.RoleSuperAdmin,
	))
	editors.PATCH("/:id", h.Update)
	editors.DELETE("/:id", h.Delete)
	editors.POST("/:id/publish", h.Publish)
	editors.POST("/:id/cancel", h.Cancel)
	editors.POST("/:id/complete", h.Complete)
}
