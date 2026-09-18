package ticket

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
	g := r.Group("/tickets")
	g.Use(requireAuth)

	// User tickets
	g.GET("/my", h.ListMine)
	g.GET("/:id", h.Get)

	// Manual issue (for testing/admin)
	admin := g.Group("")
	admin.Use(middleware.RequireRole(
		models.RoleCollegeAdmin,
		models.RoleSuperAdmin,
		models.RoleOrganizer,
	))
	admin.POST("/issue/:registration_id", h.IssueForRegistration)

	// Verify (scanner — Phase 11)
	scanner := g.Group("/verify")
	scanner.Use(middleware.RequireRole(
		models.RoleVolunteer,
		models.RoleOrganizer,
		models.RoleCollegeAdmin,
		models.RoleSuperAdmin,
	))
	scanner.POST("", h.VerifyQR)
}