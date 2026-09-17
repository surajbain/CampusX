package event

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/pagination"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

type Service struct {
	repo  *Repo
	audit *audit.Logger
}

func NewService(repo *Repo, auditLog *audit.Logger) *Service {
	return &Service{repo: repo, audit: auditLog}
}

// ---------- Create ----------

func (s *Service) Create(ctx context.Context, collegeID, actorID string, req CreateEventRequest) (*EventResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}

	slug, err := s.resolveSlug(ctx, collegeID, req.Title, req.Slug, "")
	if err != nil {
		return nil, err
	}

	e := &models.Event{
		CollegeID:         collegeID,
		CreatedBy:         &actorID,
		Slug:              slug,
		Title:             strings.TrimSpace(req.Title),
		Description:       req.Description,
		Category:          models.EventCategory(req.Category),
		PosterURL:         req.PosterURL,
		Venue:             req.Venue,
		City:              titleCase(req.City),
		Status:            models.StatusDraft,
		StartsAt:          req.StartsAt,
		EndsAt:            req.EndsAt,
		RegistrationOpens: req.RegistrationOpens,
		RegistrationClose: req.RegistrationCloses,
		PricePaise:        req.PricePaise,
		Currency:          "INR",
		Capacity:          req.Capacity,
		AllowTeams:        req.AllowTeams,
		TeamSizeMin:       req.TeamSizeMin,
		TeamSizeMax:       req.TeamSizeMax,
		PrizePoolPaise:    req.PrizePoolPaise,
		ContactEmail:      strings.ToLower(strings.TrimSpace(req.ContactEmail)),
		WhatsappLink:      req.WhatsappLink,
		Rules:             req.Rules,
		ScheduleJSON:      req.ScheduleJSON,
		FAQJSON:           req.FAQJSON,
		IsFeatured:        req.IsFeatured,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		if errors.Is(err, ErrSlugTaken) {
			return nil, apperror.Conflict("slug already in use for this college")
		}
		return nil, apperror.Internal("create event failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &collegeID,
		Action:     "EVENT_CREATED",
		EntityType: "event",
		EntityID:   &e.ID,
		Severity:   "INFO",
		Metadata:   map[string]any{"slug": e.Slug, "title": e.Title},
	})

	// Reload with College preloaded.
	loaded, _ := s.repo.FindByID(ctx, e.ID)
	return toResponse(loaded), nil
}

// ---------- Read ----------

func (s *Service) Get(ctx context.Context, id string) (*EventResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	return toResponse(e), nil
}

// GetBySlug requires both college_id and slug because slugs are per-college unique.
func (s *Service) GetBySlug(ctx context.Context, collegeID, slug string) (*EventResponse, error) {
	e, err := s.repo.FindBySlug(ctx, collegeID, strings.ToLower(slug))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	return toResponse(e), nil
}

// ListPublic returns events visible to everyone (PUBLISHED only).
func (s *Service) ListPublic(ctx context.Context, f ListFilters, p pagination.Params) ([]EventResponse, pagination.Meta, error) {
	f.OnlyPublic = true
	return s.list(ctx, f, p)
}

// ListForUser returns published events + the user's own drafts (if they own some).
func (s *Service) ListForUser(ctx context.Context, userID string, f ListFilters, p pagination.Params) ([]EventResponse, pagination.Meta, error) {
	f.OwnerID = userID
	return s.list(ctx, f, p)
}

func (s *Service) list(ctx context.Context, f ListFilters, p pagination.Params) ([]EventResponse, pagination.Meta, error) {
	events, total, err := s.repo.List(ctx, f, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list failed", err)
	}
	out := make([]EventResponse, 0, len(events))
	for i := range events {
		out = append(out, *toResponse(&events[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

// ---------- Update ----------

func (s *Service) Update(ctx context.Context, id, actorID, actorRole string, req UpdateEventRequest) (*EventResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}

	// Ownership check (SUPER_ADMIN bypasses).
	if actorRole != string(models.RoleSuperAdmin) {
		if e.CreatedBy == nil || *e.CreatedBy != actorID {
			return nil, apperror.Forbidden("only the event owner can update this event")
		}
	}

	// Only DRAFT events can be edited freely. Once PUBLISHED, restrict changes.
	if e.Status != models.StatusDraft && actorRole != string(models.RoleSuperAdmin) {
		allowed := []string{"description", "rules", "poster_url", "contact_email", "whatsapp_link"}
		if !onlyAllowedFields(req, allowed) {
			return nil, apperror.Conflict("published events can only update: description, rules, poster, contact info")
		}
	}

	updates := buildUpdates(req)
	if len(updates) == 0 {
		return nil, apperror.BadRequest("no fields to update")
	}

	if err := s.validateDatesAfterUpdate(e, req); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, apperror.Internal("update failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &e.CollegeID,
		Action:     "EVENT_UPDATED",
		EntityType: "event",
		EntityID:   &id,
		Severity:   "INFO",
		Metadata:   map[string]any{"fields": keys(updates)},
	})

	loaded, _ := s.repo.FindByID(ctx, id)
	return toResponse(loaded), nil
}

// ---------- Lifecycle ----------

func (s *Service) Publish(ctx context.Context, id, actorID, actorRole string) (*EventResponse, error) {
	return s.transition(ctx, id, actorID, actorRole, models.StatusPublished)
}

func (s *Service) Cancel(ctx context.Context, id, actorID, actorRole string) (*EventResponse, error) {
	return s.transition(ctx, id, actorID, actorRole, models.StatusCancelled)
}

func (s *Service) Complete(ctx context.Context, id, actorID, actorRole string) (*EventResponse, error) {
	return s.transition(ctx, id, actorID, actorRole, models.StatusCompleted)
}

func (s *Service) transition(ctx context.Context, id, actorID, actorRole string, to models.EventStatus) (*EventResponse, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}

	// Ownership.
	if actorRole != string(models.RoleSuperAdmin) {
		if e.CreatedBy == nil || *e.CreatedBy != actorID {
			return nil, apperror.Forbidden("only the event owner can change status")
		}
	}

	// State machine check.
	if err := CanTransition(e.Status, to); err != nil {
		return nil, err
	}

	// Publishing requires required fields.
	if to == models.StatusPublished {
		if err := requirePublishFields(e); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, id, map[string]any{"status": string(to)}); err != nil {
		return nil, apperror.Internal("update status failed", err)
	}

	action := "EVENT_" + strings.ToUpper(string(to))
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &e.CollegeID,
		Action:     action,
		EntityType: "event",
		EntityID:   &id,
		Severity:   "INFO",
		Metadata:   map[string]any{"from": string(e.Status), "to": string(to)},
	})

	loaded, _ := s.repo.FindByID(ctx, id)
	return toResponse(loaded), nil
}

// ---------- Delete ----------

func (s *Service) Delete(ctx context.Context, id, actorID, actorRole string) error {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return apperror.NotFound("event not found")
		}
		return apperror.Internal("lookup failed", err)
	}
	if actorRole != string(models.RoleSuperAdmin) {
		if e.CreatedBy == nil || *e.CreatedBy != actorID {
			return apperror.Forbidden("only the event owner can delete this event")
		}
	}
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return apperror.Internal("delete failed", err)
	}
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &e.CollegeID,
		Action:     "EVENT_DELETED",
		EntityType: "event",
		EntityID:   &id,
		Severity:   "WARN",
	})
	return nil
}

// ---------- Discovery helpers ----------

func (s *Service) Categories() []string {
	return []string{
		string(models.CatHackathon),
		string(models.CatCultural),
		string(models.CatSports),
		string(models.CatWorkshop),
		string(models.CatTechFest),
		string(models.CatOther),
	}
}

// ---------- Validation ----------

func validateCreate(req CreateEventRequest) error {
	if !slugRe.MatchString(strings.ToLower(strings.TrimSpace(req.Slug))) && req.Slug != "" {
		return apperror.BadRequest("slug must be lowercase alphanumeric with hyphens")
	}

	if req.RegistrationOpens != nil && req.RegistrationCloses != nil {
		if !req.RegistrationOpens.Before(*req.RegistrationCloses) {
			return apperror.BadRequest("registration_opens must be before registration_closes")
		}
	}
	if req.RegistrationCloses != nil && req.StartsAt != nil {
		if !req.RegistrationCloses.Before(*req.StartsAt) {
			return apperror.BadRequest("registration_closes must be before starts_at")
		}
	}
	if req.StartsAt != nil && req.EndsAt != nil {
		if !req.StartsAt.Before(*req.EndsAt) {
			return apperror.BadRequest("starts_at must be before ends_at")
		}
	}

	if req.AllowTeams {
		if req.TeamSizeMin == nil || req.TeamSizeMax == nil {
			return apperror.BadRequest("team_size_min and team_size_max are required when allow_teams=true")
		}
		if *req.TeamSizeMin > *req.TeamSizeMax {
			return apperror.BadRequest("team_size_min must be <= team_size_max")
		}
	}
	return nil
}

func (s *Service) validateDatesAfterUpdate(e *models.Event, req UpdateEventRequest) error {
	startsAt := e.StartsAt
	if req.StartsAt != nil {
		startsAt = req.StartsAt
	}
	endsAt := e.EndsAt
	if req.EndsAt != nil {
		endsAt = req.EndsAt
	}
	regOpens := e.RegistrationOpens
	if req.RegistrationOpens != nil {
		regOpens = req.RegistrationOpens
	}
	regCloses := e.RegistrationClose
	if req.RegistrationCloses != nil {
		regCloses = req.RegistrationCloses
	}

	if regOpens != nil && regCloses != nil && !regOpens.Before(*regCloses) {
		return apperror.BadRequest("registration_opens must be before registration_closes")
	}
	if regCloses != nil && startsAt != nil && !regCloses.Before(*startsAt) {
		return apperror.BadRequest("registration_closes must be before starts_at")
	}
	if startsAt != nil && endsAt != nil && !startsAt.Before(*endsAt) {
		return apperror.BadRequest("starts_at must be before ends_at")
	}
	return nil
}

func requirePublishFields(e *models.Event) error {
	if e.Title == "" {
		return apperror.BadRequest("title is required to publish")
	}
	if e.StartsAt == nil {
		return apperror.BadRequest("starts_at is required to publish")
	}
	if e.Venue == "" && e.City == "" {
		return apperror.BadRequest("venue or city is required to publish")
	}
	return nil
}

// ---------- Helpers ----------

func (s *Service) resolveSlug(ctx context.Context, collegeID, title, requested, excludeID string) (string, error) {
	slug := strings.TrimSpace(requested)
	if slug == "" {
		slug = slugify(title)
	}
	slug = strings.ToLower(slug)

	if !slugRe.MatchString(slug) {
		return "", apperror.BadRequest("could not derive a valid slug")
	}

	exists, err := s.repo.SlugExists(ctx, collegeID, slug, excludeID)
	if err != nil {
		return "", apperror.Internal("slug check failed", err)
	}
	if !exists {
		return slug, nil
	}

	for i := 2; i <= 100; i++ {
		candidate := fmt.Sprintf("%s-%d", slug, i)
		exists, err := s.repo.SlugExists(ctx, collegeID, candidate, excludeID)
		if err != nil {
			return "", apperror.Internal("slug check failed", err)
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", apperror.Conflict("could not generate unique slug")
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevHyphen := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevHyphen = false
		case r == ' ' || r == '-' || r == '_':
			if !prevHyphen && b.Len() > 0 {
				b.WriteRune('-')
				prevHyphen = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 160 {
		out = out[:160]
	}
	return out
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func buildUpdates(req UpdateEventRequest) map[string]any {
	u := map[string]any{}
	if req.Title != nil {
		u["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		u["description"] = *req.Description
	}
	if req.Category != nil {
		u["category"] = *req.Category
	}
	if req.PosterURL != nil {
		u["poster_url"] = *req.PosterURL
	}
	if req.Venue != nil {
		u["venue"] = *req.Venue
	}
	if req.City != nil {
		u["city"] = titleCase(*req.City)
	}
	if req.StartsAt != nil {
		u["starts_at"] = *req.StartsAt
	}
	if req.EndsAt != nil {
		u["ends_at"] = *req.EndsAt
	}
	if req.RegistrationOpens != nil {
		u["registration_opens"] = *req.RegistrationOpens
	}
	if req.RegistrationCloses != nil {
		u["registration_closes"] = *req.RegistrationCloses
	}
	if req.PricePaise != nil {
		u["price_paise"] = *req.PricePaise
	}
	if req.Capacity != nil {
		u["capacity"] = *req.Capacity
	}
	if req.AllowTeams != nil {
		u["allow_teams"] = *req.AllowTeams
	}
	if req.TeamSizeMin != nil {
		u["team_size_min"] = *req.TeamSizeMin
	}
	if req.TeamSizeMax != nil {
		u["team_size_max"] = *req.TeamSizeMax
	}
	if req.PrizePoolPaise != nil {
		u["prize_pool_paise"] = *req.PrizePoolPaise
	}
	if req.ContactEmail != nil {
		u["contact_email"] = strings.ToLower(strings.TrimSpace(*req.ContactEmail))
	}
	if req.WhatsappLink != nil {
		u["whatsapp_link"] = *req.WhatsappLink
	}
	if req.Rules != nil {
		u["rules"] = *req.Rules
	}
	if req.IsFeatured != nil {
		u["is_featured"] = *req.IsFeatured
	}
	return u
}

func onlyAllowedFields(req UpdateEventRequest, allowed []string) bool {
	allowedSet := map[string]bool{}
	for _, a := range allowed {
		allowedSet[a] = true
	}
	if req.Title != nil && !allowedSet["title"] {
		return false
	}
	if req.Category != nil && !allowedSet["category"] {
		return false
	}
	if req.Venue != nil && !allowedSet["venue"] {
		return false
	}
	if req.City != nil && !allowedSet["city"] {
		return false
	}
	if req.StartsAt != nil && !allowedSet["starts_at"] {
		return false
	}
	if req.EndsAt != nil && !allowedSet["ends_at"] {
		return false
	}
	if req.RegistrationOpens != nil && !allowedSet["registration_opens"] {
		return false
	}
	if req.RegistrationCloses != nil && !allowedSet["registration_closes"] {
		return false
	}
	if req.PricePaise != nil && !allowedSet["price_paise"] {
		return false
	}
	if req.Capacity != nil && !allowedSet["capacity"] {
		return false
	}
	if req.AllowTeams != nil && !allowedSet["allow_teams"] {
		return false
	}
	if req.TeamSizeMin != nil && !allowedSet["team_size_min"] {
		return false
	}
	if req.TeamSizeMax != nil && !allowedSet["team_size_max"] {
		return false
	}
	if req.PrizePoolPaise != nil && !allowedSet["prize_pool_paise"] {
		return false
	}
	if req.IsFeatured != nil && !allowedSet["is_featured"] {
		return false
	}
	return true
}

func toResponse(e *models.Event) *EventResponse {
	r := &EventResponse{
		ID:                 e.ID,
		CollegeID:          e.CollegeID,
		Slug:               e.Slug,
		Title:              e.Title,
		Description:        e.Description,
		Category:           string(e.Category),
		PosterURL:          e.PosterURL,
		Venue:              e.Venue,
		City:               e.City,
		Status:             string(e.Status),
		StartsAt:           e.StartsAt,
		EndsAt:             e.EndsAt,
		RegistrationOpens:  e.RegistrationOpens,
		RegistrationCloses: e.RegistrationClose,
		PricePaise:         e.PricePaise,
		Currency:           e.Currency,
		Capacity:           e.Capacity,
		AllowTeams:         e.AllowTeams,
		TeamSizeMin:        e.TeamSizeMin,
		TeamSizeMax:        e.TeamSizeMax,
		PrizePoolPaise:     e.PrizePoolPaise,
		ContactEmail:       e.ContactEmail,
		WhatsappLink:       e.WhatsappLink,
		Rules:              e.Rules,
		Schedule:           e.ScheduleJSON,
		FAQ:                e.FAQJSON,
		IsFeatured:         e.IsFeatured,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
	if e.CreatedBy != nil {
		r.CreatedBy = *e.CreatedBy
	}
	if e.College != nil {
		r.CollegeName = e.College.Name
		r.CollegeSlug = e.College.Slug
		r.CollegeCity = e.College.City
	}
	return r
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

var _ = time.Now
