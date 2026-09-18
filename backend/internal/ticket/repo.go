package ticket

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/pagination"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

var (
	ErrNotFound          = errors.New("ticket: not found")
	ErrRegistrationNotFound = errors.New("ticket: registration not found")
	ErrDuplicate         = errors.New("ticket: duplicate")
)

func (r *Repo) Create(ctx context.Context, t *models.Ticket) error {
	err := r.db.WithContext(ctx).Create(t).Error
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *Repo) FindByID(ctx context.Context, id string) (*models.Ticket, error) {
	var t models.Ticket
	err := r.db.WithContext(ctx).
		Preload("Registration").
		Preload("Registration.Event").
		Preload("Registration.Event.College").
		Preload("Registration.User").
		Where("id = ?", id).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) FindByRegistrationID(ctx context.Context, regID string) (*models.Ticket, error) {
	var t models.Ticket
	err := r.db.WithContext(ctx).
		Where("registration_id = ?", regID).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) FindByTicketCode(ctx context.Context, code string) (*models.Ticket, error) {
	var t models.Ticket
	err := r.db.WithContext(ctx).
		Where("ticket_code = ?", code).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) FindRegistration(ctx context.Context, id string) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("Event.College").
		Preload("User").
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

// ListForUser returns all tickets for a given user.
func (r *Repo) ListForUser(ctx context.Context, userID string, p pagination.Params) ([]models.Ticket, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Preload("Registration").
		Preload("Registration.Event").
		Preload("Registration.Event.College").
		Where("user_id = ?", userID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("issued_at DESC")
	}

	var tickets []models.Ticket
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.Ticket{}).
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

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}