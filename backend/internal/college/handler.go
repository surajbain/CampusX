package college

import (
	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/middleware"
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

// ---- College CRUD ----

func (h *Handler) Create(c *gin.Context) {
	var req CreateCollegeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	actorID := middleware.MustUserID(c)
	res, err := h.svc.CreateCollege(c.Request.Context(), req, actorID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) List(c *gin.Context) {
	allowedSorts := map[string]string{
		"name":       "name",
		"city":       "city",
		"created_at": "created_at",
	}
	p := pagination.Parse(c, allowedSorts)
	city := c.Query("city")

	res, meta, err := h.svc.ListColleges(c.Request.Context(), city, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

func (h *Handler) GetByID(c *gin.Context) {
	res, err := h.svc.GetCollege(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	res, err := h.svc.GetCollegeBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateCollegeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	actorID := middleware.MustUserID(c)
	actorRole := string(middleware.MustUserRole(c))
	res, err := h.svc.UpdateCollege(c.Request.Context(), c.Param("id"), req, actorID, actorRole)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Delete(c *gin.Context) {
	actorID := middleware.MustUserID(c)
	if err := h.svc.DeleteCollege(c.Request.Context(), c.Param("id"), actorID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "college deleted"})
}

func (h *Handler) ListCities(c *gin.Context) {
	cities, err := h.svc.ListCities(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"cities": cities})
}

func (h *Handler) Stats(c *gin.Context) {
	st, err := h.svc.Stats(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, st)
}

// ---- Admin user management ----

func (h *Handler) CreateUser(c *gin.Context) {
	collegeID := c.Param("id")
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	actorID := middleware.MustUserID(c)
	res, err := h.svc.CreateUserInCollege(c.Request.Context(), collegeID, req, actorID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) ListUsers(c *gin.Context) {
	collegeID := c.Param("id")
	allowedSorts := map[string]string{
		"created_at": "created_at",
		"email":      "email",
		"role":       "role",
	}
	p := pagination.Parse(c, allowedSorts)

	res, meta, err := h.svc.ListUsersInCollege(c.Request.Context(), collegeID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

func (h *Handler) GetUser(c *gin.Context) {
	collegeID := c.Param("id")
	userID := c.Param("user_id")
	res, err := h.svc.GetUserInCollege(c.Request.Context(), collegeID, userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	collegeID := c.Param("id")
	userID := c.Param("user_id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	actorID := middleware.MustUserID(c)
	res, err := h.svc.UpdateUserInCollege(c.Request.Context(), collegeID, userID, req, actorID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	collegeID := c.Param("id")
	userID := c.Param("user_id")
	actorID := middleware.MustUserID(c)
	if err := h.svc.DeleteUserInCollege(c.Request.Context(), collegeID, userID, actorID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "user deleted"})
}
