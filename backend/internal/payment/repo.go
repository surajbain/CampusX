package payment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/campusx/api/internal/models"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

var (
	ErrNotFound             = errors.New("payment: not found")
	ErrRegistrationNotFound = errors.New("payment: registration not found")
	ErrDuplicate            = errors.New("payment: duplicate gateway_payment_id or idempotency_key")
)

func (r *Repo) Create(ctx context.Context, p *models.Payment) error {
	err := r.db.WithContext(ctx).Create(p).Error
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *Repo) FindByID(ctx context.Context, id string) (*models.Payment, error) {
	var p models.Payment
	err := r.db.WithContext(ctx).
		Preload("Registration").
		Preload("Registration.Event").
		Preload("Registration.Event.College").
		Where("id = ?", id).
		First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) FindByRegistrationID(ctx context.Context, regID string) (*models.Payment, error) {
	var p models.Payment
	err := r.db.WithContext(ctx).
		Where("registration_id = ?", regID).
		Order("created_at DESC").
		First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) FindRegistration(ctx context.Context, id string) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.WithContext(ctx).
		Preload("Event").
		Where("id = ?", id).
		First(&reg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRegistrationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("id = ?", id).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) MarkPaid(ctx context.Context, paymentID, registrationID, gatewayPaymentID string, raw map[string]any) error {
	if registrationID == "" {
		return errors.New("payment: registrationID is empty")
	}
	if paymentID == "" {
		return errors.New("payment: paymentID is empty")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		// 1. Update payment
		paymentRes := tx.Model(&models.Payment{}).
			Where("id = ?", paymentID).
			Updates(map[string]any{
				"status":             "CAPTURED",
				"gateway_payment_id": gatewayPaymentID,
				"verified_at":        now,
				"raw_payload":        raw,
			})
		if paymentRes.Error != nil {
			if isUniqueViolation(paymentRes.Error) {
				return ErrDuplicate
			}
			return fmt.Errorf("payment update: %w", paymentRes.Error)
		}
		if paymentRes.RowsAffected == 0 {
			return errors.New("payment: not found for update")
		}

		// 2. Fetch registration, then Save() — GORM Update() silently fails for 0 rows
		var reg models.Registration
		if err := tx.Where("id = ?", registrationID).First(&reg).Error; err != nil {
			return fmt.Errorf("registration fetch: %w", err)
		}

		if reg.Status != models.RegStatusConfirmed {
			reg.Status = models.RegStatusConfirmed
			if err := tx.Save(&reg).Error; err != nil {
				return fmt.Errorf("registration save: %w", err)
			}
		}

		return nil
	})
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}