package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/campusx/api/internal/audit"
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/crypto"
	"github.com/campusx/api/pkg/jwt"
	"github.com/campusx/api/pkg/session"
)

type Service struct {
	repo     *Repo
	jwt      *jwt.Issuer
	sessions *session.Store
	audit    *audit.Logger
}

func NewService(repo *Repo, jwtIssuer *jwt.Issuer, sessions *session.Store, auditLog *audit.Logger) *Service {
	return &Service{
		repo:     repo,
		jwt:      jwtIssuer,
		sessions: sessions,
		audit:    auditLog,
	}
}

// Register creates a new STUDENT account.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Basic college existence check.
	var collegeCount int64
	if err := s.repo.db.WithContext(ctx).
		Model(&models.College{}).
		Where("id = ? AND is_active = TRUE", req.CollegeID).
		Count(&collegeCount).Error; err != nil {
		return nil, apperror.Internal("college lookup failed", err)
	}
	if collegeCount == 0 {
		return nil, apperror.BadRequest("invalid or inactive college")
	}

	hash, err := crypto.Hash(req.Password)
	if err != nil {
		return nil, apperror.Internal("password hashing failed", err)
	}

	user := &models.User{
		Email:         email,
		PasswordHash:  hash,
		FullName:      strings.TrimSpace(req.FullName),
		Phone:         req.Phone,
		Role:          models.RoleStudent,
		CollegeID:     &req.CollegeID,
		IsActive:      true,
		EmailVerified: false,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return nil, apperror.Conflict("email already registered")
		}
		return nil, apperror.Internal("user creation failed", err)
	}

	s.audit.Log(ctx, audit.Entry{
		ActorID:    &user.ID,
		CollegeID:  &req.CollegeID,
		Action:     "USER_REGISTERED",
		EntityType: "user",
		EntityID:   &user.ID,
		Severity:   "INFO",
	})

	return s.issueTokensFor(ctx, user)
}

// Login authenticates and returns tokens.
func (s *Service) Login(ctx context.Context, req LoginRequest, ip, ua string) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			s.audit.Log(ctx, audit.Entry{
				Action:    "LOGIN_FAILED",
				Severity:  "WARN",
				IPAddress: ip,
				UserAgent: ua,
				Metadata:  map[string]any{"email": email, "reason": "user_not_found"},
			})
			return nil, apperror.Unauthorized("invalid credentials")
		}
		return nil, apperror.Internal("login lookup failed", err)
	}

	if !user.IsActive {
		s.audit.Log(ctx, audit.Entry{
			ActorID:   &user.ID,
			Action:    "LOGIN_FAILED",
			Severity:  "WARN",
			IPAddress: ip,
			UserAgent: ua,
			Metadata:  map[string]any{"reason": "inactive"},
		})
		return nil, apperror.Forbidden("account disabled")
	}

	ok, err := crypto.Verify(req.Password, user.PasswordHash)
	if err != nil || !ok {
		s.audit.Log(ctx, audit.Entry{
			ActorID:   &user.ID,
			Action:    "LOGIN_FAILED",
			Severity:  "WARN",
			IPAddress: ip,
			UserAgent: ua,
			Metadata:  map[string]any{"reason": "bad_password"},
		})
		return nil, apperror.Unauthorized("invalid credentials")
	}

	_ = s.repo.TouchLastLogin(ctx, user.ID)

	s.audit.Log(ctx, audit.Entry{
		ActorID:   &user.ID,
		Action:    "LOGIN_SUCCESS",
		Severity:  "INFO",
		IPAddress: ip,
		UserAgent: ua,
	})

	return s.issueTokensFor(ctx, user)
}

// Refresh rotates a refresh token: old one is revoked, new pair issued.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, apperror.Unauthorized("invalid refresh token")
	}

	if _, err := s.sessions.Verify(ctx, claims.TokenID, claims.UserID); err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			// Possible reuse attack — revoke all sessions for this user.
			_ = s.sessions.RevokeAll(ctx, claims.UserID)
			s.audit.Log(ctx, audit.Entry{
				ActorID:  &claims.UserID,
				Action:   "REFRESH_REUSE_DETECTED",
				Severity: "CRITICAL",
				Metadata: map[string]any{"jti": claims.TokenID},
			})
			return nil, apperror.Unauthorized("refresh token revoked")
		}
		return nil, apperror.Internal("session verify failed", err)
	}

	// Revoke old refresh token (rotation).
	_ = s.sessions.Revoke(ctx, claims.TokenID, claims.UserID)

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperror.Unauthorized("user not found")
	}
	if !user.IsActive {
		return nil, apperror.Forbidden("account disabled")
	}

	return s.issueTokensFor(ctx, user)
}

// Logout revokes a single refresh token.
func (s *Service) Logout(ctx context.Context, userID, refreshToken string) error {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return apperror.Unauthorized("invalid refresh token")
	}
	if claims.UserID != userID {
		return apperror.Forbidden("token does not belong to you")
	}
	if err := s.sessions.Revoke(ctx, claims.TokenID, userID); err != nil {
		return apperror.Internal("revoke failed", err)
	}
	s.audit.Log(ctx, audit.Entry{
		ActorID:  &userID,
		Action:   "LOGOUT",
		Severity: "INFO",
	})
	return nil
}

// ChangePassword updates the user's password and revokes all sessions.
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("user not found")
	}

	ok, err := crypto.Verify(req.CurrentPassword, user.PasswordHash)
	if err != nil || !ok {
		return apperror.Unauthorized("current password is incorrect")
	}

	newHash, err := crypto.Hash(req.NewPassword)
	if err != nil {
		return apperror.Internal("hash failed", err)
	}
	if err := s.repo.UpdatePassword(ctx, userID, newHash); err != nil {
		return apperror.Internal("update failed", err)
	}

	// Force re-login on all devices.
	_ = s.sessions.RevokeAll(ctx, userID)

	s.audit.Log(ctx, audit.Entry{
		ActorID:  &userID,
		Action:   "PASSWORD_CHANGED",
		Severity: "WARN",
	})
	return nil
}

// Me returns the current user.
func (s *Service) Me(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.NotFound("user not found")
	}
	return toUserResponse(user), nil
}

// --- helpers ---

func (s *Service) issueTokensFor(ctx context.Context, user *models.User) (*AuthResponse, error) {
	jti := uuid.NewString()

	var collegeID string
	if user.CollegeID != nil {
		collegeID = *user.CollegeID
	}

	in := jwt.IssueInput{
		UserID:    user.ID,
		Role:      string(user.Role),
		CollegeID: collegeID,
		TokenID:   jti,
	}

	access, err := s.jwt.IssueAccess(in)
	if err != nil {
		return nil, apperror.Internal("issue access failed", err)
	}
	refresh, err := s.jwt.IssueRefresh(in)
	if err != nil {
		return nil, apperror.Internal("issue refresh failed", err)
	}
	if err := s.sessions.Save(ctx, jti, user.ID, string(user.Role), s.jwt.RefreshTTL()); err != nil {
		return nil, apperror.Internal("save session failed", err)
	}

	return &AuthResponse{
		User:         *toUserResponse(user),
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(s.jwt.AccessTTL().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func toUserResponse(u *models.User) *UserResponse {
	r := &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FullName:      u.FullName,
		Phone:         u.Phone,
		Role:          string(u.Role),
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
	if u.CollegeID != nil {
		r.CollegeID = *u.CollegeID
	}
	return r
}

var _ = time.Second // silence import if unused
