package registration

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/pagination"
	"github.com/campusx/api/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /api/v1/events/:id/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}

	eventID := c.Param("id")
	userID := middleware.MustUserID(c)

	res, err := h.svc.Register(c.Request.Context(), eventID, userID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

// GET /api/v1/registrations/my
func (h *Handler) ListMine(c *gin.Context) {
	userID := middleware.MustUserID(c)
	p := pagination.Parse(c, map[string]string{
		"created_at": "created_at",
	})
	res, meta, err := h.svc.ListForUser(c.Request.Context(), userID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

// GET /api/v1/registrations/:id
func (h *Handler) Get(c *gin.Context) {
	userID := middleware.MustUserID(c)
	role := middleware.MustUserRole(c)
	isAdmin := role == models.RoleCollegeAdmin || role == models.RoleSuperAdmin || role == models.RoleOrganizer

	res, err := h.svc.Get(c.Request.Context(), c.Param("id"), userID, isAdmin)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// DELETE /api/v1/registrations/:id
func (h *Handler) Cancel(c *gin.Context) {
	userID := middleware.MustUserID(c)
	role := middleware.MustUserRole(c)
	isAdmin := role == models.RoleCollegeAdmin || role == models.RoleSuperAdmin

	if err := h.svc.Cancel(c.Request.Context(), c.Param("id"), userID, isAdmin); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "registration cancelled"})
}

// GET /api/v1/events/:id/registrations (admin/organizer)
func (h *Handler) ListForEvent(c *gin.Context) {
	p := pagination.Parse(c, map[string]string{
		"created_at": "created_at",
	})
	eventID := c.Param("id")
	res, meta, err := h.svc.ListForEvent(c.Request.Context(), eventID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}