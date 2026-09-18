package registration

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/internal/scopes"
	"github.com/campusx/api/pkg/pagination"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

var (
	ErrNotFound       = errors.New("registration: not found")
	ErrAlreadyExists  = errors.New("registration: already exists")
	ErrEventNotOpen   = errors.New("registration: event not open")
	ErrEventFull      = errors.New("registration: event full")
	ErrEventNotFound  = errors.New("registration: event not found")
)

// ---- Create ----

func (r *Repo) Create(ctx context.Context, reg *models.Registration) error {
	err := r.db.WithContext(ctx).Create(reg).Error
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

// ---- Find ----

func (r *Repo) FindByID(ctx context.Context, id string) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("Event.College").
		Preload("User").
		Where("id = ?", id).
		First(&reg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *Repo) FindByEventAndUser(ctx context.Context, eventID, userID, regType string) (*models.Registration, error) {
	var reg models.Registration
	err := r.db.WithContext(ctx).
		Where("event_id = ? AND user_id = ? AND type = ?", eventID, userID, regType).
		First(&reg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

// ---- List for user ----

func (r *Repo) ListForUser(ctx context.Context, userID string, p pagination.Params) ([]models.Registration, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Registration{}).
		Preload("Event").
		Preload("Event.College").
		Where("user_id = ?", userID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("created_at DESC")
	}

	var regs []models.Registration
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&regs).Error; err != nil {
		return nil, 0, err
	}
	return regs, total, nil
}

// ---- List for event (admin/organizer) ----

func (r *Repo) ListForEvent(ctx context.Context, eventID string, p pagination.Params) ([]models.Registration, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Registration{}).
		Preload("User").
		Where("event_id = ?", eventID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("created_at DESC")
	}

	var regs []models.Registration
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&regs).Error; err != nil {
		return nil, 0, err
	}
	return regs, total, nil
}

// ---- Count ----

func (r *Repo) CountForEvent(ctx context.Context, eventID string, regType string) (int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Registration{}).
		Where("event_id = ?", eventID)
	if regType != "" {
		q = q.Where("type = ?", regType)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ---- Update ----

func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.Registration{}).
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

// ---- Cancel ----

func (r *Repo) Cancel(ctx context.Context, id string) error {
	return r.Update(ctx, id, map[string]any{"status": "CANCELLED"})
}

// ---- Lookup ----

func (r *Repo) FindEvent(ctx context.Context, id string) (*models.Event, error) {
	var e models.Event
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repo) FindUser(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ---- Ensure user belongs to college ----

func (r *Repo) UserInCollege(ctx context.Context, userID, collegeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Scopes(scopes.College(collegeID)).
		Where("id = ?", userID).
		Count(&count).Error
	return count > 0, err
}

// ---- Helpers ----

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}