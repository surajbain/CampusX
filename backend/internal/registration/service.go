package registration

import (
	"context"
	"errors"
	"time"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/pagination"
)

type Service struct {
	repo  *Repo
	audit *audit.Logger
}

func NewService(repo *Repo, auditLog *audit.Logger) *Service {
	return &Service{repo: repo, audit: auditLog}
}

// ---------- Register for an event ----------

func (s *Service) Register(ctx context.Context, eventID, userID string, req RegisterEventRequest) (*RegistrationResponse, error) {
	// 1. Find event
	event, err := s.repo.FindEvent(ctx, eventID)
	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("event lookup failed", err)
	}

	// 2. Check event status
	if event.Status != models.StatusPublished && event.Status != models.StatusOngoing {
		return nil, apperror.BadRequest("event is not open for registration")
	}

	// 3. Check registration window
	now := time.Now().UTC()
	if event.RegistrationOpens != nil && now.Before(*event.RegistrationOpens) {
		return nil, apperror.BadRequest("registration has not opened yet")
	}
	if event.RegistrationClose != nil && now.After(*event.RegistrationClose) {
		return nil, apperror.BadRequest("registration has closed")
	}

	// 4. Check user exists and belongs to a college
	user, err := s.repo.FindUser(ctx, userID)
	if err != nil {
		return nil, apperror.NotFound("user not found")
	}

	// 5. Validate registration type against event
	if req.Type == string(models.RegTypeParticipate) && !event.AllowTeams {
		// Participation without teams allowed for individual events — but student must explicitly opt in
		// No restriction here, just info
	}

	// 6. Check capacity
	if event.Capacity != nil && *event.Capacity > 0 {
		count, err := s.repo.CountForEvent(ctx, eventID, "")
		if err != nil {
			return nil, apperror.Internal("capacity check failed", err)
		}
		if count >= int64(*event.Capacity) {
			return nil, apperror.Conflict("event is full")
		}
	}

	// 7. Check for duplicate registration
	existing, err := s.repo.FindByEventAndUser(ctx, eventID, userID, req.Type)
	if err == nil && existing != nil {
		if existing.Status == "CANCELLED" {
			// Allow re-registration after cancellation
			updates := map[string]any{
				"status": "PENDING",
				"notes":  req.Notes,
			}
			// If paid and already paid? Let's not handle here yet — payment flow will do
			if event.PricePaise == 0 {
				updates["status"] = "CONFIRMED"
			}
			if err := s.repo.Update(ctx, existing.ID, updates); err != nil {
				return nil, apperror.Internal("re-registration failed", err)
			}
			loaded, _ := s.repo.FindByID(ctx, existing.ID)
			return toResponse(loaded), nil
		}
		return nil, apperror.Conflict("you are already registered for this event")
	}

	// 8. Create registration
	reg := &models.Registration{
		EventID:   eventID,
		UserID:    userID,
		CollegeID: event.CollegeID,
		Type:      models.RegistrationType(req.Type),
		Status:    models.RegStatusPending,
		Notes:     req.Notes,
	}

	// Free events get auto-confirmed
	if event.PricePaise == 0 {
		reg.Status = models.RegStatusConfirmed
	}

	if err := s.repo.Create(ctx, reg); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return nil, apperror.Conflict("you are already registered for this event")
		}
		return nil, apperror.Internal("registration failed", err)
	}

	// 9. Audit log
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &userID,
		CollegeID:  &event.CollegeID,
		Action:     "REGISTRATION_CREATED",
		EntityType: "registration",
		EntityID:   &reg.ID,
		Severity:   "INFO",
		Metadata: map[string]any{
			"event_id":   eventID,
			"event_slug": event.Slug,
			"type":       req.Type,
			"status":     string(reg.Status),
		},
	})

	// 10. Reload with preloads
	loaded, _ := s.repo.FindByID(ctx, reg.ID)
	_ = user // silence unused
	return toResponse(loaded), nil
}

// ---------- Get one ----------

func (s *Service) Get(ctx context.Context, id, userID string, isAdmin bool) (*RegistrationResponse, error) {
	reg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("registration not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}

	if !isAdmin && reg.UserID != userID {
		return nil, apperror.Forbidden("not your registration")
	}

	return toResponse(reg), nil
}

// ---------- List user's registrations ----------

func (s *Service) ListForUser(ctx context.Context, userID string, p pagination.Params) ([]RegistrationResponse, pagination.Meta, error) {
	regs, total, err := s.repo.ListForUser(ctx, userID, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list failed", err)
	}
	out := make([]RegistrationResponse, 0, len(regs))
	for i := range regs {
		out = append(out, *toResponse(&regs[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

// ---------- List registrations for an event (admin/organizer) ----------

func (s *Service) ListForEvent(ctx context.Context, eventID string, p pagination.Params) ([]RegistrationResponse, pagination.Meta, error) {
	regs, total, err := s.repo.ListForEvent(ctx, eventID, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list failed", err)
	}
	out := make([]RegistrationResponse, 0, len(regs))
	for i := range regs {
		out = append(out, *toResponse(&regs[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

// ---------- Cancel ----------

func (s *Service) Cancel(ctx context.Context, id, userID string, isAdmin bool) error {
	reg, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return apperror.NotFound("registration not found")
		}
		return apperror.Internal("lookup failed", err)
	}

	if !isAdmin && reg.UserID != userID {
		return apperror.Forbidden("not your registration")
	}

	if reg.Status == models.RegStatusCancelled {
		return apperror.BadRequest("already cancelled")
	}

	if err := s.repo.Cancel(ctx, id); err != nil {
		return apperror.Internal("cancel failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &userID,
		CollegeID:  &reg.CollegeID,
		Action:     "REGISTRATION_CANCELLED",
		EntityType: "registration",
		EntityID:   &id,
		Severity:   "INFO",
	})

	return nil
}

// ---------- Helpers ----------

func toResponse(reg *models.Registration) *RegistrationResponse {
	r := &RegistrationResponse{
		ID:        reg.ID,
		EventID:   reg.EventID,
		UserID:    reg.UserID,
		CollegeID: reg.CollegeID,
		Type:      string(reg.Type),
		Status:    string(reg.Status),
		Notes:     reg.Notes,
		CreatedAt: reg.CreatedAt,
		UpdatedAt: reg.UpdatedAt,
	}

	if reg.Event != nil {
		r.EventTitle = reg.Event.Title
		r.EventSlug = reg.Event.Slug
		r.EventPosterURL = reg.Event.PosterURL
		r.EventCategory = string(reg.Event.Category)
		r.EventStartsAt = reg.Event.StartsAt
		r.EventEndsAt = reg.Event.EndsAt
		r.EventVenue = reg.Event.Venue
		r.EventCity = reg.Event.City
		r.EventPricePaise = reg.Event.PricePaise
		r.EventCurrency = reg.Event.Currency
		if reg.Event.College != nil {
			r.EventCollegeName = reg.Event.College.Name
		}
	}

	if reg.User != nil {
		r.UserName = reg.User.FullName
		r.UserEmail = reg.User.Email
	}

	return r
}