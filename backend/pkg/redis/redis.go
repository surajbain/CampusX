package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/logger"
)

// Connect opens a Redis client and pings it with retry.
func Connect(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	log := logger.From()

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	backoff := 500 * time.Millisecond
	maxAttempts := 10

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := client.Ping(ctx).Err(); err == nil {
			log.Info().Str("addr", cfg.Addr).Msg("redis connected")
			return client, nil
		}
		log.Warn().Int("attempt", attempt).Dur("backoff", backoff).Msg("redis not ready, retrying")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
		backoff *= 2
		if backoff > 5*time.Second {
			backoff = 5 * time.Second
		}
	}
	return nil, fmt.Errorf("redis connection failed after %d attempts", maxAttempts)
}

func Health(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}

func Close(client *redis.Client) error {
	return client.Close()
}
