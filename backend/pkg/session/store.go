package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
}

var (
	ErrSessionNotFound = errors.New("session: not found")
	ErrSessionRevoked  = errors.New("session: revoked")
)

func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

// refreshKey namespaces refresh tokens in Redis.
func refreshKey(jti string) string { return "refresh:" + jti }

// userSessionsKey tracks all active jtis per user (for mass revoke).
func userSessionsKey(userID string) string { return "user:sessions:" + userID }

// Save persists a refresh token id (jti) for a user until ttl.
// Stores minimal metadata: userID + role.
func (s *Store) Save(ctx context.Context, jti, userID, role string, ttl time.Duration) error {
	pipe := s.rdb.TxPipeline()
	pipe.Set(ctx, refreshKey(jti), userID+"|"+role, ttl)
	pipe.SAdd(ctx, userSessionsKey(userID), jti)
	pipe.Expire(ctx, userSessionsKey(userID), ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("session: save: %w", err)
	}
	return nil
}

// Verify checks jti is still valid for userID. Returns stored metadata.
func (s *Store) Verify(ctx context.Context, jti, userID string) (string, error) {
	v, err := s.rdb.Get(ctx, refreshKey(jti)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrSessionNotFound
		}
		return "", fmt.Errorf("session: verify: %w", err)
	}
	stored := userID + "|"
	if len(v) < len(stored) || v[:len(stored)] != stored {
		return "", ErrSessionNotFound
	}
	return v, nil
}

// Revoke removes a single refresh token.
func (s *Store) Revoke(ctx context.Context, jti, userID string) error {
	pipe := s.rdb.TxPipeline()
	pipe.Del(ctx, refreshKey(jti))
	pipe.SRem(ctx, userSessionsKey(userID), jti)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("session: revoke: %w", err)
	}
	return nil
}

// RevokeAll removes every refresh token for a user (logout everywhere).
func (s *Store) RevokeAll(ctx context.Context, userID string) error {
	jtis, err := s.rdb.SMembers(ctx, userSessionsKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("session: members: %w", err)
	}
	pipe := s.rdb.TxPipeline()
	for _, jti := range jtis {
		pipe.Del(ctx, refreshKey(jti))
	}
	pipe.Del(ctx, userSessionsKey(userID))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("session: revoke all: %w", err)
	}
	return nil
}

// CountActive returns how many active refresh tokens a user has.
func (s *Store) CountActive(ctx context.Context, userID string) (int64, error) {
	return s.rdb.SCard(ctx, userSessionsKey(userID)).Result()
}
