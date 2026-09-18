package ticket

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

// GET /api/v1/tickets/my
func (h *Handler) ListMine(c *gin.Context) {
	userID := middleware.MustUserID(c)
	p := pagination.Parse(c, map[string]string{
		"issued_at": "issued_at",
	})
	res, meta, err := h.svc.ListForUser(c.Request.Context(), userID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

// GET /api/v1/tickets/:id
func (h *Handler) Get(c *gin.Context) {
	userID := middleware.MustUserID(c)
	role := middleware.MustUserRole(c)
	isAdmin := role == models.RoleCollegeAdmin ||
		role == models.RoleSuperAdmin ||
		role == models.RoleOrganizer ||
		role == models.RoleVolunteer

	res, err := h.svc.Get(c.Request.Context(), c.Param("id"), userID, isAdmin)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// POST /api/v1/tickets/:id/issue  (internal — call from registration flow)
func (h *Handler) IssueForRegistration(c *gin.Context) {
	regID := c.Param("registration_id")
	res, err := h.svc.IssueForRegistration(c.Request.Context(), regID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

// POST /api/v1/tickets/verify  (for scanner — Phase 11)
type verifyRequest struct {
	PayloadB64 string `json:"payload_b64" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}

func (h *Handler) VerifyQR(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.BadRequest("invalid request body"))
		return
	}
	res, err := h.svc.VerifyQR(c.Request.Context(), req.PayloadB64, req.Signature)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}