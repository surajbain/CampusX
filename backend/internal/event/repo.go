package event

import (
	"context"
	"errors"
	"strings"
	"time"

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
	ErrNotFound  = errors.New("event: not found")
	ErrSlugTaken = errors.New("event: slug already taken for this college")
)

// ---- Filters ----

type ListFilters struct {
	CollegeID  string
	Category   string
	City       string
	Status     string
	Query      string     // free-text on title + description
	From       *time.Time // starts_at >= from
	To         *time.Time // starts_at <= to
	Featured   *bool
	OnlyPublic bool   // force status=PUBLISHED (for public list)
	OwnerID    string // if set, include owner's drafts too
}

// ---- Create / Read / Update / Delete ----

func (r *Repo) Create(ctx context.Context, e *models.Event) error {
	err := r.db.WithContext(ctx).Create(e).Error
	if err != nil {
		if isUniqueViolation(err) {
			return ErrSlugTaken
		}
		return err
	}
	return nil
}

func (r *Repo) FindByID(ctx context.Context, id string) (*models.Event, error) {
	var e models.Event
	err := r.db.WithContext(ctx).
		Preload("College").
		Where("id = ?", id).
		First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repo) FindBySlug(ctx context.Context, collegeID, slug string) (*models.Event, error) {
	var e models.Event
	err := r.db.WithContext(ctx).
		Preload("College").
		Where("college_id = ? AND slug = ?", collegeID, slug).
		First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// SlugExists checks if a slug is taken in the given college (optionally excluding an ID).
func (r *Repo) SlugExists(ctx context.Context, collegeID, slug, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Event{}).
		Where("college_id = ? AND slug = ?", collegeID, slug)
	if excludeID != "" {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.Event{}).
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

func (r *Repo) SoftDelete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&models.Event{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- List ----

func (r *Repo) List(ctx context.Context, f ListFilters, p pagination.Params) ([]models.Event, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Event{}).
		Preload("College")

	if f.CollegeID != "" {
		q = q.Scopes(scopes.College(f.CollegeID))
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.City != "" {
		q = q.Where("LOWER(city) = LOWER(?)", f.City)
	}
	if f.Featured != nil {
		if *f.Featured {
			q = q.Where("is_featured = TRUE")
		} else {
			q = q.Where("is_featured = FALSE")
		}
	}

	// Visibility
	if f.OnlyPublic {
		q = q.Where("status = 'PUBLISHED'")
	} else if f.OwnerID != "" {
		// Include owner's drafts + all published
		q = q.Where("status = 'PUBLISHED' OR created_by = ?", f.OwnerID)
	}

	if f.Status != "" && !f.OnlyPublic {
		q = q.Where("status = ?", f.Status)
	}

	if f.Query != "" {
		like := "%" + strings.ToLower(f.Query) + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", like, like)
	}
	if f.From != nil {
		q = q.Where("starts_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("starts_at <= ?", *f.To)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Default sort: upcoming events first.
	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("starts_at ASC NULLS LAST, created_at DESC")
	}

	var events []models.Event
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

// ---- Helpers ----

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
