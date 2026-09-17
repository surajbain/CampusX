package college

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/crypto"
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

// ---------- College CRUD ----------

func (s *Service) CreateCollege(ctx context.Context, req CreateCollegeRequest, actorID string) (*CollegeResponse, error) {
	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	if !slugRe.MatchString(slug) {
		return nil, apperror.BadRequest("slug must be lowercase alphanumeric with hyphens")
	}

	c := &models.College{
		Name:         strings.TrimSpace(req.Name),
		Slug:         slug,
		City:         titleCase(strings.TrimSpace(req.City)),
		State:        titleCase(strings.TrimSpace(req.State)),
		LogoURL:      strings.TrimSpace(req.LogoURL),
		Website:      strings.TrimSpace(req.Website),
		ContactEmail: strings.ToLower(strings.TrimSpace(req.ContactEmail)),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		IsActive:     true,
	}

	if err := s.repo.CreateCollege(ctx, c); err != nil {
		if errors.Is(err, ErrSlugTaken) {
			return nil, apperror.Conflict("slug already in use")
		}
		return nil, apperror.Internal("create college failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		Action:     "COLLEGE_CREATED",
		EntityType: "college",
		EntityID:   &c.ID,
		Severity:   "INFO",
		Metadata:   map[string]any{"slug": c.Slug},
	})

	return toCollegeResponse(c), nil
}

func (s *Service) GetCollege(ctx context.Context, id string) (*CollegeResponse, error) {
	c, err := s.repo.FindCollegeByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("college not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	return toCollegeResponse(c), nil
}

func (s *Service) GetCollegeBySlug(ctx context.Context, slug string) (*CollegeResponse, error) {
	c, err := s.repo.FindCollegeBySlug(ctx, strings.ToLower(slug))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("college not found")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	return toCollegeResponse(c), nil
}

func (s *Service) ListColleges(ctx context.Context, city string, p pagination.Params) ([]CollegeResponse, pagination.Meta, error) {
	colleges, total, err := s.repo.ListColleges(ctx, city, true, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list failed", err)
	}
	out := make([]CollegeResponse, 0, len(colleges))
	for i := range colleges {
		out = append(out, *toCollegeResponse(&colleges[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

func (s *Service) UpdateCollege(ctx context.Context, id string, req UpdateCollegeRequest, actorID, actorRole string) (*CollegeResponse, error) {
	// Tenant check for COLLEGE_ADMIN is done by middleware; SUPER_ADMIN bypasses.
	// Additional guard: only SUPER_ADMIN can change is_active (soft control).
	updates := map[string]any{}

	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.City != nil {
		updates["city"] = titleCase(strings.TrimSpace(*req.City))
	}
	if req.State != nil {
		updates["state"] = titleCase(strings.TrimSpace(*req.State))
	}
	if req.LogoURL != nil {
		updates["logo_url"] = strings.TrimSpace(*req.LogoURL)
	}
	if req.Website != nil {
		updates["website"] = strings.TrimSpace(*req.Website)
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = strings.ToLower(strings.TrimSpace(*req.ContactEmail))
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = strings.TrimSpace(*req.ContactPhone)
	}
	if req.IsActive != nil {
		if actorRole != string(models.RoleSuperAdmin) {
			return nil, apperror.Forbidden("only SUPER_ADMIN can change is_active")
		}
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		return nil, apperror.BadRequest("no fields to update")
	}

	if err := s.repo.UpdateCollege(ctx, id, updates); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("college not found")
		}
		return nil, apperror.Internal("update failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		Action:     "COLLEGE_UPDATED",
		EntityType: "college",
		EntityID:   &id,
		Severity:   "INFO",
		Metadata:   map[string]any{"fields": keys(updates)},
	})

	return s.GetCollege(ctx, id)
}

func (s *Service) DeleteCollege(ctx context.Context, id, actorID string) error {
	if err := s.repo.SoftDeleteCollege(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return apperror.NotFound("college not found")
		}
		return apperror.Internal("delete failed", err)
	}
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		Action:     "COLLEGE_DELETED",
		EntityType: "college",
		EntityID:   &id,
		Severity:   "WARN",
	})
	return nil
}

func (s *Service) ListCities(ctx context.Context) ([]string, error) {
	cities, err := s.repo.ListCities(ctx)
	if err != nil {
		return nil, apperror.Internal("cities failed", err)
	}
	return cities, nil
}

func (s *Service) Stats(ctx context.Context) (map[string]int64, error) {
	st, err := s.repo.Stats(ctx)
	if err != nil {
		return nil, apperror.Internal("stats failed", err)
	}
	return st, nil
}

// ---------- Admin user management ----------

func (s *Service) CreateUserInCollege(ctx context.Context, collegeID string, req CreateUserRequest, actorID string) (*AdminUserResponse, error) {
	hash, err := crypto.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal("hash failed", err)
	}

	u := &models.User{
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash:  hash,
		FullName:      strings.TrimSpace(req.FullName),
		Phone:         strings.TrimSpace(req.Phone),
		Role:          models.Role(req.Role),
		CollegeID:     &collegeID,
		IsActive:      true,
		EmailVerified: false, // admin-created users need to verify (or set true for now)
	}

	if err := s.repo.CreateUser(ctx, u); err != nil {
		if strings.Contains(err.Error(), "already registered") {
			return nil, apperror.Conflict("email already registered")
		}
		return nil, apperror.Internal("create user failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &collegeID,
		Action:     "USER_CREATED",
		EntityType: "user",
		EntityID:   &u.ID,
		Severity:   "INFO",
		Metadata:   map[string]any{"role": req.Role},
	})

	return toAdminUserResponse(u), nil
}

func (s *Service) ListUsersInCollege(ctx context.Context, collegeID string, p pagination.Params) ([]AdminUserResponse, pagination.Meta, error) {
	users, total, err := s.repo.ListUsersInCollege(ctx, collegeID, p)
	if err != nil {
		return nil, pagination.Meta{}, apperror.Internal("list users failed", err)
	}
	out := make([]AdminUserResponse, 0, len(users))
	for i := range users {
		out = append(out, *toAdminUserResponse(&users[i]))
	}
	return out, pagination.NewMeta(p, total), nil
}

func (s *Service) GetUserInCollege(ctx context.Context, collegeID, userID string) (*AdminUserResponse, error) {
	u, err := s.repo.FindUserInCollege(ctx, collegeID, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("user not found in this college")
		}
		return nil, apperror.Internal("lookup failed", err)
	}
	return toAdminUserResponse(u), nil
}

func (s *Service) UpdateUserInCollege(ctx context.Context, collegeID, userID string, req UpdateUserRequest, actorID string) (*AdminUserResponse, error) {
	// Prevent demoting self.
	if actorID == userID && req.Role != nil && *req.Role != string(models.RoleCollegeAdmin) {
		return nil, apperror.BadRequest("cannot change your own role")
	}

	updates := map[string]any{}
	if req.FullName != nil {
		updates["full_name"] = strings.TrimSpace(*req.FullName)
	}
	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		return nil, apperror.BadRequest("no fields to update")
	}

	if err := s.repo.UpdateUserInCollege(ctx, collegeID, userID, updates); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, apperror.NotFound("user not found in this college")
		}
		return nil, apperror.Internal("update failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &collegeID,
		Action:     "USER_UPDATED",
		EntityType: "user",
		EntityID:   &userID,
		Severity:   "INFO",
		Metadata:   map[string]any{"fields": keys(updates)},
	})

	return s.GetUserInCollege(ctx, collegeID, userID)
}

func (s *Service) DeleteUserInCollege(ctx context.Context, collegeID, userID, actorID string) error {
	if actorID == userID {
		return apperror.BadRequest("cannot delete yourself")
	}
	if err := s.repo.SoftDeleteUserInCollege(ctx, collegeID, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			return apperror.NotFound("user not found in this college")
		}
		return apperror.Internal("delete failed", err)
	}
	s.audit.Log(ctx, audit.Entry{
		ActorID:    &actorID,
		CollegeID:  &collegeID,
		Action:     "USER_DELETED",
		EntityType: "user",
		EntityID:   &userID,
		Severity:   "WARN",
	})
	return nil
}

// ---------- helpers ----------

func toCollegeResponse(c *models.College) *CollegeResponse {
	return &CollegeResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		City:         c.City,
		State:        c.State,
		LogoURL:      c.LogoURL,
		Website:      c.Website,
		ContactEmail: c.ContactEmail,
		ContactPhone: c.ContactPhone,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toAdminUserResponse(u *models.User) *AdminUserResponse {
	r := &AdminUserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FullName:      u.FullName,
		Phone:         u.Phone,
		Role:          string(u.Role),
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	}
	if u.CollegeID != nil {
		r.CollegeID = *u.CollegeID
	}
	return r
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

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

var _ = time.Now
