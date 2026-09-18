package payment

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/logger"
)

type Service struct {
	repo     *Repo
	provider Provider
	audit    *audit.Logger
}

func NewService(repo *Repo, provider Provider, auditLog *audit.Logger) *Service {
	return &Service{repo: repo, provider: provider, audit: auditLog}
}

func (s *Service) CreatePayment(ctx context.Context, userID string, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	reg, err := s.repo.FindRegistration(ctx, req.RegistrationID)
	if err != nil {
		if errors.Is(err, ErrRegistrationNotFound) {
			return nil, apperror.NotFound("registration not found")
		}
		return nil, apperror.Internal("registration lookup failed", err)
	}

	if reg.UserID != userID {
		return nil, apperror.Forbidden("not your registration")
	}
	if reg.Status == models.RegStatusConfirmed {
		return nil, apperror.Conflict("already confirmed")
	}
	if reg.Status == models.RegStatusCancelled {
		return nil, apperror.BadRequest("registration cancelled")
	}
	if reg.Event == nil || reg.Event.PricePaise <= 0 {
		return nil, apperror.BadRequest("event is free — no payment needed")
	}

	existing, err := s.repo.FindByRegistrationID(ctx, req.RegistrationID)
	if err == nil && existing != nil {
		if existing.Status == models.PayCaptured {
			return nil, apperror.Conflict("already paid")
		}
		if existing.Status == models.PayCreated && existing.GatewayOrderID != nil {
			return &CreatePaymentResponse{
				PaymentID:      existing.ID,
				GatewayOrderID: *existing.GatewayOrderID,
				AmountPaise:    existing.AmountPaise,
				Currency:       existing.Currency,
				Provider:       string(existing.Provider),
				IdempotencyKey: existing.IdempotencyKey,
			}, nil
		}
	}

	order, err := s.provider.CreateOrder(ctx, CreateOrderRequest{
		AmountPaise: reg.Event.PricePaise,
		Currency:    reg.Event.Currency,
		Receipt:     reg.ID,
		Notes: map[string]string{
			"registration_id": reg.ID,
			"event_id":        reg.EventID,
			"user_id":         userID,
		},
	})
	if err != nil {
		return nil, apperror.Internal("order creation failed", err)
	}

	idempotencyKey := uuid.NewString()
	p := &models.Payment{
		RegistrationID: reg.ID,
		UserID:         userID,
		EventID:        reg.EventID,
		Provider:       models.PaymentProvider(s.provider.Name()),
		GatewayOrderID: &order.GatewayOrderID,
		IdempotencyKey: idempotencyKey,
		AmountPaise:    order.AmountPaise,
		Currency:       order.Currency,
		Status:         models.PayCreated,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return nil, apperror.Conflict("payment already exists")
		}
		return nil, apperror.Internal("payment persistence failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &userID,
		CollegeID:  &reg.CollegeID,
		Action:     "PAYMENT_CREATED",
		EntityType: "payment",
		EntityID:   &p.ID,
		Severity:   "INFO",
		Metadata: map[string]any{
			"registration_id": reg.ID,
			"amount_paise":    p.AmountPaise,
			"provider":        s.provider.Name(),
		},
	})

	return &CreatePaymentResponse{
		PaymentID:      p.ID,
		GatewayOrderID: order.GatewayOrderID,
		AmountPaise:    p.AmountPaise,
		Currency:       p.Currency,
		Provider:       string(p.Provider),
		PublicKey:      order.PublicKey,
		IdempotencyKey: idempotencyKey,
	}, nil
}

func (s *Service) VerifyPayment(ctx context.Context, userID string, req VerifyPaymentAPIRequest) (*PaymentResponse, error) {
	p, err := s.repo.FindByID(ctx, req.PaymentID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("payment not found")
		}
		return nil, apperror.Internal("payment lookup failed", err)
	}

	if p.UserID != userID {
		return nil, apperror.Forbidden("not your payment")
	}

	if p.Status == models.PayCaptured {
		return toResponse(p), nil
	}

	if err := s.provider.VerifyPayment(ctx, VerifyPaymentRequest{
		GatewayOrderID:   req.GatewayOrderID,
		GatewayPaymentID: req.GatewayPaymentID,
		Signature:        req.Signature,
	}); err != nil {
		if errors.Is(err, ErrInvalidSignature) {
			s.audit.Log(ctx, audit.Entry{
				ActorID:    &userID,
				Action:     "PAYMENT_SIGNATURE_INVALID",
				EntityType: "payment",
				EntityID:   &p.ID,
				Severity:   "CRITICAL",
			})
			return nil, apperror.Unauthorized("invalid payment signature")
		}
		return nil, apperror.Internal("signature verify failed", err)
	}

	fetched, err := s.provider.FetchPayment(ctx, req.GatewayPaymentID)
	if err != nil {
		return nil, apperror.Internal("gateway fetch failed", err)
	}

	if fetched.Status != "captured" && fetched.Status != "authorized" {
		return nil, apperror.BadRequest("payment not successful: " + fetched.Status)
	}

	if fetched.AmountPaise > 0 && fetched.AmountPaise != p.AmountPaise {
		s.audit.Log(ctx, audit.Entry{
			ActorID:    &userID,
			Action:     "PAYMENT_AMOUNT_MISMATCH",
			EntityType: "payment",
			EntityID:   &p.ID,
			Severity:   "CRITICAL",
		})
		return nil, apperror.BadRequest("payment amount mismatch")
	}

		// 7. Mark paid (atomic)
	if err := s.repo.MarkPaid(ctx, p.ID, p.RegistrationID, req.GatewayPaymentID, fetched.RawData); err != nil {
		// Log full error for debugging
		logger.From().Error().
			Err(err).
			Str("payment_id", p.ID).
			Str("registration_id", p.RegistrationID).
			Str("gateway_payment_id", req.GatewayPaymentID).
			Msg("MarkPaid failed — transaction rolled back")

		if errors.Is(err, ErrDuplicate) {
			s.audit.Log(ctx, audit.Entry{
				ActorID:    &userID,
				Action:     "PAYMENT_REPLAY_DETECTED",
				EntityType: "payment",
				EntityID:   &p.ID,
				Severity:   "CRITICAL",
				Metadata: map[string]any{
					"gateway_payment_id": req.GatewayPaymentID,
				},
			})
			return nil, apperror.Conflict("payment ID already used")
		}
		return nil, apperror.Internal("mark paid failed", err)
	}
	
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &userID,
		Action:     "PAYMENT_VERIFIED",
		EntityType: "payment",
		EntityID:   &p.ID,
		Severity:   "INFO",
		Metadata: map[string]any{
			"gateway_payment_id": req.GatewayPaymentID,
		},
	})

	loaded, _ := s.repo.FindByID(ctx, p.ID)
	return toResponse(loaded), nil
}

func (s *Service) Get(ctx context.Context, id, userID string, isAdmin bool) (*PaymentResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("payment not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	if !isAdmin && p.UserID != userID {
		return nil, apperror.Forbidden("not your payment")
	}
	return toResponse(p), nil
}

func toResponse(p *models.Payment) *PaymentResponse {
	r := &PaymentResponse{
		ID:             p.ID,
		RegistrationID: p.RegistrationID,
		UserID:         p.UserID,
		EventID:        p.EventID,
		Provider:       string(p.Provider),
		AmountPaise:    p.AmountPaise,
		Currency:       p.Currency,
		Status:         string(p.Status),
		VerifiedAt:     p.VerifiedAt,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	if p.GatewayPaymentID != nil {
		r.GatewayPaymentID = *p.GatewayPaymentID
	}
	if p.GatewayOrderID != nil {
		r.GatewayOrderID = *p.GatewayOrderID
	}
	return r
}