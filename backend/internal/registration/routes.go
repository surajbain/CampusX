package registration

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/models"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	requireAuth gin.HandlerFunc,
) {
	// ---- User registrations ----
	regs := r.Group("/registrations")
	regs.Use(requireAuth)
	regs.GET("/my", h.ListMine)
	regs.GET("/:id", h.Get)
	regs.DELETE("/:id", h.Cancel)

	// ---- Register for event (student) ----
	events := r.Group("/events")
	events.Use(requireAuth)
	events.POST("/:id/register", h.Register)

	// ---- List event registrations (admin/organizer) ----
	admin := r.Group("/events/:id/registrations")
	admin.Use(requireAuth)
	admin.Use(middleware.RequireRole(
		models.RoleOrganizer,
		models.RoleCollegeAdmin,
		models.RoleSuperAdmin,
	))
	admin.GET("", h.ListForEvent)
}