package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.BadRequest("invalid request body"))
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.BadRequest("invalid request body"))
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.BadRequest("invalid request body"))
		return
	}
	userID := middleware.MustUserID(c)
	if err := h.svc.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "logged out"})
}

func (h *Handler) Me(c *gin.Context) {
	userID := middleware.MustUserID(c)
	res, err := h.svc.Me(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.BadRequest("invalid request body"))
		return
	}
	userID := middleware.MustUserID(c)
	if err := h.svc.ChangePassword(c.Request.Context(), userID, req); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "password changed. Please log in again."})
}
