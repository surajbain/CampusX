package college

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/models"
)

// RegisterRoutes wires all college endpoints.
//
// Public:
//
//	GET  /api/v1/colleges
//	GET  /api/v1/colleges/:id
//	GET  /api/v1/colleges/slug/:slug
//	GET  /api/v1/colleges/cities
//	GET  /api/v1/colleges/stats
//
// SUPER_ADMIN only:
//
//	POST   /api/v1/colleges
//	DELETE /api/v1/colleges/:id
//
// COLLEGE_ADMIN or SUPER_ADMIN (with tenant check):
//
//	PATCH  /api/v1/colleges/:id
//	POST   /api/v1/colleges/:college_id/users
//	GET    /api/v1/colleges/:college_id/users
//	GET    /api/v1/colleges/:college_id/users/:user_id
//	PATCH  /api/v1/colleges/:college_id/users/:user_id
//	DELETE /api/v1/colleges/:college_id/users/:user_id
func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	requireAuth gin.HandlerFunc,
) {
	g := r.Group("/colleges")

	// ---- Public (no auth) ----
	g.GET("", h.List)
	g.GET("/cities", h.ListCities)
	g.GET("/stats", h.Stats)
	g.GET("/slug/:slug", h.GetBySlug)
	g.GET("/:id", h.GetByID)

	// ---- SUPER_ADMIN only ----
	admin := g.Group("")
	admin.Use(requireAuth)
	admin.Use(middleware.RequireRole(models.RoleSuperAdmin))
	admin.POST("", h.Create)
	admin.DELETE("/:id", h.Delete)

	// ---- COLLEGE_ADMIN or SUPER_ADMIN ----
	tenantGroup := g.Group("")
	tenantGroup.Use(requireAuth)
	tenantGroup.Use(middleware.RequireRole(models.RoleCollegeAdmin, models.RoleSuperAdmin))
	tenantGroup.PATCH("/:id", h.Update)

	// ---- User management, tenant-scoped ----
	userGroup := g.Group("/:id/users")
	userGroup.Use(requireAuth)
	userGroup.Use(middleware.RequireRole(models.RoleCollegeAdmin, models.RoleSuperAdmin))
	userGroup.Use(middleware.RequireTenantParam())
	userGroup.POST("", h.CreateUser)
	userGroup.GET("", h.ListUsers)
	userGroup.GET("/:user_id", h.GetUser)
	userGroup.PATCH("/:user_id", h.UpdateUser)
	userGroup.DELETE("/:user_id", h.DeleteUser)
}
