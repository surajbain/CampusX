package tests

import (
	"testing"

	"github.com/campusx/api/pkg/config"
)

func TestConfigValidation_FailsWithoutSecrets(t *testing.T) {
	t.Setenv("JWT_ACCESS_SECRET", "")
	t.Setenv("JWT_REFRESH_SECRET", "")
	t.Setenv("QR_HMAC_SECRET", "")
	t.Setenv("DB_PASSWORD", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected config validation error, got nil")
	}
}
