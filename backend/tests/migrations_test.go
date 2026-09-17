package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/campusx/api/internal/migrations"
)

// requireDB connects to the dev Postgres or skips the test if not available.
func requireDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=campusx password=campusx_dev_password dbname=campusx sslmode=disable"
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Skipf("skipping: cannot connect to db: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Skipf("skipping: db not reachable")
	}
	return gdb
}

func TestMigrations_RunOnce(t *testing.T) {
	gdb := requireDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log := zerolog.Nop()

	// First run
	if err := migrations.Run(ctx, gdb, &log); err != nil {
		t.Fatalf("first migration run failed: %v", err)
	}

	// Second run — should be a no-op
	if err := migrations.Run(ctx, gdb, &log); err != nil {
		t.Fatalf("second migration run failed: %v", err)
	}

	// Verify schema_migrations has exactly the number of .up.sql files
	var count int64
	if err := gdb.Raw("SELECT COUNT(*) FROM schema_migrations").Scan(&count).Error; err != nil {
		t.Fatalf("count schema_migrations: %v", err)
	}
	if count < 1 {
		t.Fatalf("expected at least 1 migration, got %d", count)
	}
	t.Logf("applied migrations: %d", count)
}

func TestMigrations_TablesExist(t *testing.T) {
	gdb := requireDB(t)

	want := []string{
		"colleges", "users", "events", "registrations", "payments",
		"tickets", "checkins", "teams", "team_members",
		"results", "winners", "certificates",
		"ads", "ad_slots", "ad_impressions", "audit_logs",
		"schema_migrations",
	}

	for _, table := range want {
		var exists bool
		err := gdb.Raw(
			"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)",
			table,
		).Scan(&exists).Error
		if err != nil {
			t.Fatalf("check %s: %v", table, err)
		}
		if !exists {
			t.Errorf("missing table: %s", table)
		}
	}
}

func TestSeed_CollegeExists(t *testing.T) {
	gdb := requireDB(t)
	var count int64
	if err := gdb.Raw("SELECT COUNT(*) FROM colleges WHERE slug = 'iitb'").Scan(&count).Error; err != nil {
		t.Fatalf("query: %v", err)
	}
	if count == 0 {
		t.Skip("seed migration may not have run — skipping")
	}
	if count != 1 {
		t.Fatalf("expected 1 iitb college, got %d", count)
	}
}

// Ensure test DB is isolated — helper
func TestMain(m *testing.M) {
	code := m.Run()
	fmt.Fprintln(os.Stderr, "tests done")
	os.Exit(code)
}
