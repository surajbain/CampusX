package event

import (
	"time"

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

// ---------- Create ----------

func (h *Handler) Create(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	// College ID comes from JWT (tenant-scoped).
	collegeID := middleware.MustCollegeID(c)
	actorID := middleware.MustUserID(c)

	res, err := h.svc.Create(c.Request.Context(), collegeID, actorID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, res)
}

// ---------- Read ----------

func (h *Handler) GetByID(c *gin.Context) {
	res, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	collegeID := c.Query("college_id")
	if collegeID == "" {
		response.Fail(c, apperror.BadRequest("college_id query param is required"))
		return
	}
	res, err := h.svc.GetBySlug(c.Request.Context(), collegeID, c.Param("slug"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// ---------- List (public) ----------

func (h *Handler) List(c *gin.Context) {
	filters := parseFilters(c)
	p := parsePagination(c)

	res, meta, err := h.svc.ListPublic(c.Request.Context(), filters, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

// ---------- List (my events) ----------

func (h *Handler) ListMine(c *gin.Context) {
	userID := middleware.MustUserID(c)
	filters := parseFilters(c)
	p := parsePagination(c)

	res, meta, err := h.svc.ListForUser(c.Request.Context(), userID, filters, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

// ---------- Featured / Trending ----------

func (h *Handler) Featured(c *gin.Context) {
	f := parseFilters(c)
	tr := true
	f.Featured = &tr
	f.OnlyPublic = true
	p := parsePagination(c)

	res, meta, err := h.svc.ListPublic(c.Request.Context(), f, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

func (h *Handler) Trending(c *gin.Context) {
	// Phase 5: same as list, sorted by starts_at ASC.
	// Phase 17: replace with registration-count ranking.
	f := parseFilters(c)
	f.OnlyPublic = true
	p := parsePagination(c)

	res, meta, err := h.svc.ListPublic(c.Request.Context(), f, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.WithMeta(c, 200, res, meta)
}

// ---------- Categories ----------

func (h *Handler) Categories(c *gin.Context) {
	response.OK(c, gin.H{"categories": h.svc.Categories()})
}

// ---------- Update ----------

func (h *Handler) Update(c *gin.Context) {
	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, apperror.WithDetails(
			apperror.CodeUnprocessable, "invalid request body",
			map[string]any{"error": err.Error()},
		))
		return
	}
	actorID := middleware.MustUserID(c)
	actorRole := string(middleware.MustUserRole(c))

	res, err := h.svc.Update(c.Request.Context(), c.Param("id"), actorID, actorRole, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// ---------- Lifecycle ----------

func (h *Handler) Publish(c *gin.Context) {
	h.transition(c, models.StatusPublished)
}

func (h *Handler) Cancel(c *gin.Context) {
	h.transition(c, models.StatusCancelled)
}

func (h *Handler) Complete(c *gin.Context) {
	h.transition(c, models.StatusCompleted)
}

func (h *Handler) transition(c *gin.Context, to models.EventStatus) {
	actorID := middleware.MustUserID(c)
	actorRole := string(middleware.MustUserRole(c))

	var (
		res *EventResponse
		err error
	)
	switch to {
	case models.StatusPublished:
		res, err = h.svc.Publish(c.Request.Context(), c.Param("id"), actorID, actorRole)
	case models.StatusCancelled:
		res, err = h.svc.Cancel(c.Request.Context(), c.Param("id"), actorID, actorRole)
	case models.StatusCompleted:
		res, err = h.svc.Complete(c.Request.Context(), c.Param("id"), actorID, actorRole)
	}
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// ---------- Delete ----------

func (h *Handler) Delete(c *gin.Context) {
	actorID := middleware.MustUserID(c)
	actorRole := string(middleware.MustUserRole(c))
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), actorID, actorRole); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "event deleted"})
}

// ---------- Helpers ----------

func parseFilters(c *gin.Context) ListFilters {
	f := ListFilters{
		CollegeID: c.Query("college_id"),
		Category:  c.Query("category"),
		City:      c.Query("city"),
		Status:    c.Query("status"),
		Query:     c.Query("q"),
	}
	if v := c.Query("featured"); v == "true" {
		b := true
		f.Featured = &b
	} else if v == "false" {
		b := false
		f.Featured = &b
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		} else if t, err := time.Parse("2006-01-02", v); err == nil {
			f.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		} else if t, err := time.Parse("2006-01-02", v); err == nil {
			// inclusive end of day
			t = t.Add(24*time.Hour - time.Second)
			f.To = &t
		}
	}
	return f
}

func parsePagination(c *gin.Context) pagination.Params {
	allowed := map[string]string{
		"starts_at":  "starts_at",
		"created_at": "created_at",
		"title":      "title",
		"city":       "city",
	}
	return pagination.Parse(c, allowed)
}
