package ticket

import (
	"context"
	"errors"
	"time"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/logger"
	"github.com/campusx/api/pkg/pagination"
)

type Service struct {
	repo  *Repo
	qr    *QRManager
	audit *audit.Logger
}

func NewService(repo *Repo, qr *QRManager, auditLog *audit.Logger) *Service {
	return &Service{repo: repo, qr: qr, audit: auditLog}
}

// ---------- Auto-issue ticket when registration is CONFIRMED ----------

// IssueForRegistration is idempotent — if a ticket already exists, returns it.
func (s *Service) IssueForRegistration(ctx context.Context, regID string) (*TicketResponse, error) {
	// Idempotency: return existing if found
	if existing, err := s.repo.FindByRegistrationID(ctx, regID); err == nil {
		loaded, _ := s.repo.FindByID(ctx, existing.ID)
		return toResponse(loaded), nil
	}

	// Fetch registration
	reg, err := s.repo.FindRegistration(ctx, regID)
	if err != nil {
		if errors.Is(err, ErrRegistrationNotFound) {
			return nil, apperror.NotFound("registration not found")
		}
		return nil, apperror.Internal("registration lookup failed", err)
	}

	if reg.Status != models.RegStatusConfirmed {
		return nil, apperror.BadRequest("registration is not confirmed")
	}

	// Compute validity window: from event registration_closes (or now) to event ends_at
	var validFrom, validUntil *time.Time
	if reg.Event != nil {
		if reg.Event.RegistrationClose != nil {
			validFrom = reg.Event.RegistrationClose
		}
		if reg.Event.EndsAt != nil {
			validUntil = reg.Event.EndsAt
		}
	}

	// Generate QR
	payload := QRPayload{
		TicketID: regID, // will be replaced with actual ticket ID after creation
		UserID:   reg.UserID,
		EventID:  reg.EventID,
		IssuedAt: time.Now().Unix(),
	}
	payloadJSON, sig, err := s.qr.Sign(payload)
	if err != nil {
		return nil, apperror.Internal("QR sign failed", err)
	}

	// Create ticket
	ticketCode := GenerateTicketCode()
	t := &models.Ticket{
		RegistrationID: regID,
		UserID:         reg.UserID,
		EventID:        reg.EventID,
		TicketCode:     ticketCode,
		QRPayload:      payloadJSON,
		QRSignature:    sig,
		ValidFrom:      validFrom,
		ValidUntil:     validUntil,
		Status:         models.TicketActive,
		IssuedAt:       time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, t); err != nil {
		if errors.Is(err, ErrDuplicate) {
			// Race: another request created it. Fetch and return.
			existing, ferr := s.repo.FindByRegistrationID(ctx, regID)
			if ferr != nil {
				return nil, apperror.Internal("ticket race fetch failed", ferr)
			}
			loaded, _ := s.repo.FindByID(ctx, existing.ID)
			return toResponse(loaded), nil
		}
		return nil, apperror.Internal("ticket creation failed", err)
	}

	// Re-sign QR with the real ticket ID
	payload.TicketID = t.ID
	finalPayloadJSON, finalSig, err := s.qr.Sign(payload)
	if err != nil {
		// Non-fatal — log and continue with placeholder
		logger.From().Warn().Err(err).Msg("ticket: re-sign failed, using initial payload")
	} else {
		_ = s.repo.Update(ctx, t.ID, map[string]any{
			"qr_payload":   finalPayloadJSON,
			"qr_signature": finalSig,
		})
		t.QRPayload = finalPayloadJSON
		t.QRSignature = finalSig
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &reg.UserID,
		CollegeID:  &reg.CollegeID,
		Action:     "TICKET_ISSUED",
		EntityType: "ticket",
		EntityID:   &t.ID,
		Severity:   "INFO",
		Metadata: map[string]any{
			"registration_id": regID,
			"event_id":        reg.EventID,
			"ticket_code":     ticketCode,
		},
	})

	loaded, _ := s.repo.FindByID(ctx, t.ID)
	return toResponse(loaded), nil
}

// ---------- Read ----------

func (s *Service) Get(ctx context.Context, id, userID string, isAdmin bool) (*TicketResponse, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("ticket not found")
		}
		return nil, apperror.Internal("ticket lookup failed", err)
	}
	if !isAdmin && t.UserID != userID {
		return nil, apperror.Forbidden("not your ticket")
	}
	return toResponse(t), nil
}

func (s *Service) ListForUser(ctx context.Context, userID string, p pagination.Params) ([]TicketResponse, pagination.Meta, error) {
	tickets, total, err := s.repo.ListForUser(ctx, userID, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list failed", err)
	}
	out := make([]TicketResponse, 0, len(tickets))
	for i := range tickets {
		out = append(out, *toResponse(&tickets[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

// ---------- Verify (for scanner — Phase 11) ----------

// VerifyQR checks if a QR payload + signature is valid and returns the ticket.
func (s *Service) VerifyQR(ctx context.Context, payloadB64, signature string) (*TicketResponse, error) {
	// Initial signature check (without time window)
	payload, err := s.qr.Verify(payloadB64, signature, nil, nil)
	if err != nil {
		return nil, apperror.Unauthorized("invalid QR")
	}

	// Fetch ticket
	t, err := s.repo.FindByID(ctx, payload.TicketID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("ticket not found")
		}
		return nil, apperror.Internal("ticket lookup failed", err)
	}

	// Re-verify with time window
	if _, err := s.qr.Verify(payloadB64, signature, t.ValidFrom, t.ValidUntil); err != nil {
		if errors.Is(err, ErrQRNotYetValid) {
			return nil, apperror.BadRequest("ticket not yet valid")
		}
		if errors.Is(err, ErrQRExpired) {
			return nil, apperror.BadRequest("ticket expired")
		}
		return nil, apperror.Unauthorized("invalid QR")
	}

	return toResponse(t), nil
}

// ---------- Helpers ----------

func toResponse(t *models.Ticket) *TicketResponse {
	r := &TicketResponse{
		ID:             t.ID,
		RegistrationID: t.RegistrationID,
		UserID:         t.UserID,
		EventID:        t.EventID,
		TicketCode:     t.TicketCode,
		QRPayload:      t.QRPayload,
		QRSignature:    t.QRSignature,
		ValidFrom:      t.ValidFrom,
		ValidUntil:     t.ValidUntil,
		Status:         string(t.Status),
		IssuedAt:       t.IssuedAt,
		UsedAt:         t.UsedAt,
	}

	if t.Registration != nil {
				if t.Registration.Event != nil {
			e := t.Registration.Event
			r.EventTitle = e.Title
			r.EventSlug = e.Slug
			r.EventPosterURL = e.PosterURL
			r.EventCategory = string(e.Category)
			r.EventStartsAt = e.StartsAt
			r.EventEndsAt = e.EndsAt
			r.EventVenue = e.Venue
			r.EventCity = e.City
			if e.College != nil {
				r.EventCollegeName = e.College.Name
			}
			r.EventWhatsappLink = e.WhatsappLink   // ← YE NAYA
		}
		if t.Registration.User != nil {
			r.UserName = t.Registration.User.FullName
			r.UserEmail = t.Registration.User.Email
		}
	}

		if t.Registration != nil {
		if t.Registration.Event != nil {
			e := t.Registration.Event
			// ... existing fields ...
			r.EventWhatsappLink = e.WhatsappLink   // ← ye add karo
		}
	}

	return r
}