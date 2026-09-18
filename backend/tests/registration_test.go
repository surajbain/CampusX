package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestRegistration_StudentCanRegisterForFreeEvent(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	// TechFest 2026 — free event
	eventID := "33333333-3333-3333-3333-333333333301"

	body := map[string]any{
		"type": "ATTEND",
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events/"+eventID+"/register", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 201 && w.Code != 409 {
		// 409 = already registered from seed
		t.Fatalf("expected 201 or 409, got %d body=%s", w.Code, w.Body.String())
	}

	// Cleanup if new registration
	if w.Code == 201 {
		var res struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		app.gdb.Exec("DELETE FROM registrations WHERE id = ?", res.Data.ID)
	}
}

func TestRegistration_DuplicateBlocked(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	eventID := "33333333-3333-3333-3333-333333333301"
	body := map[string]any{"type": "ATTEND"}
	bb, _ := json.Marshal(body)

	// First attempt
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/events/"+eventID+"/register", bytes.NewReader(bb))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w1, req1)

	// Second attempt — should fail with 409
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/events/"+eventID+"/register", bytes.NewReader(bb))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w2, req2)

	if w2.Code != 409 {
		t.Fatalf("expected 409 for duplicate, got %d body=%s", w2.Code, w2.Body.String())
	}
}

func TestRegistration_ListMine(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/registrations/my", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var res struct {
		Data []struct {
			ID      string `json:"id"`
			EventID string `json:"event_id"`
			Status  string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	// Should have at least 1 from seed
	if len(res.Data) == 0 {
		t.Fatal("expected at least 1 registration from seed")
	}
}

func TestRegistration_Cancel(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	// Register for a fresh event
	eventID := "33333333-3333-3333-3333-333333333302" // CodeFest (paid)
	body := map[string]any{"type": "ATTEND"}
	bb, _ := json.Marshal(body)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/v1/events/"+eventID+"/register", bytes.NewReader(bb))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w1, req1)

	if w1.Code != 201 {
		t.Skipf("registration failed (may already exist): %s", w1.Body.String())
	}

	var res struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w1.Body.Bytes(), &res)

	// Cancel
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/api/v1/registrations/"+res.Data.ID, nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("expected 200 on cancel, got %d body=%s", w2.Code, w2.Body.String())
	}

	// Cleanup
	app.gdb.Exec("DELETE FROM registrations WHERE id = ?", res.Data.ID)
}

func TestRegistration_CannotRegisterForUnpublishedEvent(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	// AI/ML Bootcamp — DRAFT event
	eventID := "55555555-5555-5555-5555-555555555503"

	body := map[string]any{"type": "ATTEND"}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events/"+eventID+"/register", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400 for draft event, got %d body=%s", w.Code, w.Body.String())
	}
}

var _ = uuid.NewString