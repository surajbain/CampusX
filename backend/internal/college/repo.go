package college

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
	ErrNotFound  = errors.New("college: not found")
	ErrSlugTaken = errors.New("college: slug already taken")
)

// ---- College CRUD ----

func (r *Repo) CreateCollege(ctx context.Context, c *models.College) error {
	err := r.db.WithContext(ctx).Create(c).Error
	if err != nil {
		if isUniqueViolation(err) {
			return ErrSlugTaken
		}
		return err
	}
	return nil
}

func (r *Repo) FindCollegeByID(ctx context.Context, id string) (*models.College, error) {
	var c models.College
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) FindCollegeBySlug(ctx context.Context, slug string) (*models.College, error) {
	var c models.College
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListColleges returns colleges matching filters.
// If city is provided, filters by city (case-insensitive).
func (r *Repo) ListColleges(ctx context.Context, city string, onlyActive bool, p pagination.Params) ([]models.College, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.College{})

	if onlyActive {
		q = q.Where("is_active = TRUE")
	}
	if city != "" {
		q = q.Where("LOWER(city) = LOWER(?)", city)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("name ASC")
	}

	var colleges []models.College
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&colleges).Error; err != nil {
		return nil, 0, err
	}
	return colleges, total, nil
}

func (r *Repo) UpdateCollege(ctx context.Context, id string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.College{}).
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

func (r *Repo) SoftDeleteCollege(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&models.College{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) ListCities(ctx context.Context) ([]string, error) {
	var cities []string
	err := r.db.WithContext(ctx).
		Model(&models.College{}).
		Where("is_active = TRUE").
		Distinct().
		Order("city ASC").
		Pluck("city", &cities).Error
	return cities, err
}

func (r *Repo) Stats(ctx context.Context) (map[string]int64, error) {
	out := map[string]int64{}

	var total int64
	if err := r.db.WithContext(ctx).Model(&models.College{}).Where("is_active = TRUE").Count(&total).Error; err != nil {
		return nil, err
	}
	out["colleges"] = total

	var events int64
	if err := r.db.WithContext(ctx).Model(&models.Event{}).
		Where("status = 'PUBLISHED'").Count(&events).Error; err != nil {
		return nil, err
	}
	out["published_events"] = events

	return out, nil
}

// ---- Admin user management (tenant-scoped) ----

func (r *Repo) CreateUser(ctx context.Context, u *models.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("email already registered")
		}
		return err
	}
	return nil
}

func (r *Repo) FindUserInCollege(ctx context.Context, collegeID, userID string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Scopes(scopes.College(collegeID)).
		Where("id = ?", userID).
		First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) ListUsersInCollege(ctx context.Context, collegeID string, p pagination.Params) ([]models.User, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.User{}).
		Scopes(scopes.College(collegeID))

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if order := p.OrderClause(); order != "" {
		q = q.Order(order)
	} else {
		q = q.Order("created_at DESC")
	}

	var users []models.User
	if err := q.Limit(p.Limit).Offset(p.Offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *Repo) UpdateUserInCollege(ctx context.Context, collegeID, userID string, updates map[string]any) error {
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Scopes(scopes.College(collegeID)).
		Where("id = ?", userID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) SoftDeleteUserInCollege(ctx context.Context, collegeID, userID string) error {
	res := r.db.WithContext(ctx).
		Scopes(scopes.College(collegeID)).
		Delete(&models.User{}, "id = ?", userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- helpers ----

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
