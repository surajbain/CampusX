package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"uid"`
	Role      string    `json:"role"`
	CollegeID string    `json:"cid,omitempty"`
	TokenType TokenType `json:"typ"`
	TokenID   string    `json:"jti"`
	jwt.RegisteredClaims
}

type Issuer struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
}

type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
}

var (
	ErrInvalidToken   = errors.New("jwt: invalid token")
	ErrExpiredToken   = errors.New("jwt: token expired")
	ErrWrongTokenType = errors.New("jwt: wrong token type")
)

func NewIssuer(cfg Config) *Issuer {
	return &Issuer{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
		issuer:        cfg.Issuer,
	}
}

type IssueInput struct {
	UserID    string
	Role      string
	CollegeID string
	TokenID   string // jti — pass a UUID; caller manages rotation
}

// IssueAccess generates an access token.
func (i *Issuer) IssueAccess(in IssueInput) (string, error) {
	return i.sign(in, AccessToken, i.accessTTL, i.accessSecret)
}

// IssueRefresh generates a refresh token.
func (i *Issuer) IssueRefresh(in IssueInput) (string, error) {
	return i.sign(in, RefreshToken, i.refreshTTL, i.refreshSecret)
}

func (i *Issuer) sign(in IssueInput, typ TokenType, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    in.UserID,
		Role:      in.Role,
		CollegeID: in.CollegeID,
		TokenType: typ,
		TokenID:   in.TokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   in.UserID,
			ID:        in.TokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("jwt: sign: %w", err)
	}
	return signed, nil
}

// ParseAccess validates an access token.
func (i *Issuer) ParseAccess(token string) (*Claims, error) {
	return i.parse(token, AccessToken, i.accessSecret)
}

// ParseRefresh validates a refresh token.
func (i *Issuer) ParseRefresh(token string) (*Claims, error) {
	return i.parse(token, RefreshToken, i.refreshSecret)
}

func (i *Issuer) parse(tokenStr string, want TokenType, secret []byte) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	}, jwt.WithIssuer(i.issuer), jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims.TokenType != want {
		return nil, ErrWrongTokenType
	}
	return claims, nil
}

// TTLs for informational purposes.
func (i *Issuer) AccessTTL() time.Duration  { return i.accessTTL }
func (i *Issuer) RefreshTTL() time.Duration { return i.refreshTTL }
