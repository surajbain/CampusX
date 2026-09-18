package payment

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	userID := middleware.MustUserID(c)
	res, err := h.svc.CreatePayment(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) Verify(c *gin.Context) {
	var req VerifyPaymentAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	userID := middleware.MustUserID(c)
	res, err := h.svc.VerifyPayment(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	userID := middleware.MustUserID(c)
	role := middleware.MustUserRole(c)
	isAdmin := role == models.RoleCollegeAdmin ||
		role == models.RoleSuperAdmin ||
		role == models.RoleOrganizer

	res, err := h.svc.Get(c.Request.Context(), c.Param("id"), userID, isAdmin)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}