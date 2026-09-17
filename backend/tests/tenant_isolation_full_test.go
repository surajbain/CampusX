package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// This is the critical test: IITB admin cannot access IITD users.
func TestTenantIsolation_CrossCollegeBlocked(t *testing.T) {
	app := setupFullApp(t)

	iitbAdmin := app.login(t, "admin@iitb.edu", "Admin@123")

	// IITD college ID.
	iitdCollegeID := "44444444-4444-4444-4444-444444444401"

	// IITB admin tries to list IITD users → should be 403.
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/colleges/"+iitdCollegeID+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+iitbAdmin)
	app.router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403 for cross-tenant, got %d body=%s", w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Error   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if res.Error.Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %s", res.Error.Code)
	}
}

func TestTenantIsolation_OwnCollegeAllowed(t *testing.T) {
	app := setupFullApp(t)

	iitbAdmin := app.login(t, "admin@iitb.edu", "Admin@123")
	iitbCollegeID := "11111111-1111-1111-1111-111111111111"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/colleges/"+iitbCollegeID+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+iitbAdmin)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 for own college, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestTenantIsolation_SuperAdminBypasses(t *testing.T) {
	app := setupFullApp(t)

	superToken := app.login(t, "super@campusx.dev", "SuperAdmin@123")
	iitdCollegeID := "44444444-4444-4444-4444-444444444401"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/colleges/"+iitdCollegeID+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+superToken)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("SUPER_ADMIN should bypass tenant, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestTenantIsolation_CreateUserInOwnCollege(t *testing.T) {
	app := setupFullApp(t)

	iitbAdmin := app.login(t, "admin@iitb.edu", "Admin@123")
	iitbCollegeID := "11111111-1111-1111-1111-111111111111"

	email := "test-vol-" + uuid.NewString() + "@iitb.edu"
	body := map[string]any{
		"email":     email,
		"password":  "Volunteer@123",
		"full_name": "Test Volunteer",
		"role":      "VOLUNTEER",
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/colleges/"+iitbCollegeID+"/users", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+iitbAdmin)
	app.router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	// Cleanup.
	app.gdb.Exec("DELETE FROM users WHERE email = ?", email)
}
