package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/campusx/api/internal/router"
	"github.com/campusx/api/pkg/config"
	"github.com/campusx/api/pkg/jwt"
	redisPkg "github.com/campusx/api/pkg/redis"
)

type testApp struct {
	router *gin.Engine
	gdb    *gorm.DB
	issuer *jwt.Issuer
}

func setupFullApp(t *testing.T) *testApp {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=campusx password=campusx_dev_password dbname=campusx sslmode=disable"
	}
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Skipf("db not available: %v", err)
	}
	sqlDB, _ := gdb.DB()
	if sqlDB.Ping() != nil {
		t.Skipf("db ping failed")
	}

	rdb, err := redisPkg.Connect(t.Context(), redisCfg())
	if err != nil {
		t.Skipf("redis not available: %v", err)
	}

	issuer := jwt.NewIssuer(jwt.Config{
		AccessSecret:  "test_access_secret_at_least_32_characters_long!!",
		RefreshSecret: "test_refresh_secret_at_least_32_characters_long!!",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    1 * time.Hour,
		Issuer:        "campusx-test",
	})

	// Minimal config for test
	cfg := &config.Config{
		App: config.AppConfig{
			Name: "CampusX",
			Env:  "test",
		},
		HTTP: config.HTTPConfig{
			Port:        "8080",
			CORSOrigins: []string{"http://localhost:3000"},
		},
		JWT: config.JWTConfig{
			AccessSecret:  "test_access_secret_at_least_32_characters_long!!",
			RefreshSecret: "test_refresh_secret_at_least_32_characters_long!!",
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    1 * time.Hour,
			Issuer:        "campusx-test",
		},
	}

	// Build router using shared function (same as production)
	routerEngine := router.Build(router.Options{
		Config:    cfg,
		DB:        gdb,
		Redis:     rdb,
		JWT:       issuer,
		UseTestRL: true, // no rate limits in tests
	})

	return &testApp{router: routerEngine, gdb: gdb, issuer: issuer}
}

func redisCfg() (cfg struct {
	Addr     string
	Password string
	DB       int
}) {
	cfg.Addr = "localhost:6379"
	return
}

func (a *testApp) login(t *testing.T, email, password string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	a.router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("login failed: %d %s", w.Code, w.Body.String())
	}
	var res struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	return res.Data.AccessToken
}

func TestCollege_PublicList(t *testing.T) {
	app := setupFullApp(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/colleges", nil)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCollege_SuperAdminCanCreate(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "super@campusx.dev", "SuperAdmin@123")

	body := map[string]any{
		"name":          "IIT Madras",
		"slug":          "iitm",
		"city":          "Chennai",
		"state":         "Tamil Nadu",
		"contact_email": "events@iitm.ac.in",
	}
	bb, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/colleges", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	// Cleanup.
	app.gdb.Exec("DELETE FROM colleges WHERE slug = 'iitm'")
}

func TestCollege_StudentCannotCreate(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	body := map[string]any{"name": "X", "slug": "x-c", "city": "Y"}
	bb, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/colleges", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}
