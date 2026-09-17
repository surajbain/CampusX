package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// ---- Public list ----

func TestEvent_PublicList(t *testing.T) {
	app := setupFullApp(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/events", nil)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Data    []struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	// All public events must be PUBLISHED.
	for _, e := range res.Data {
		if e.Status != "PUBLISHED" {
			t.Fatalf("public list leaked non-published event: %s", e.Status)
		}
	}
}

// ---- Create (organizer) ----

func TestEvent_OrganizerCanCreate(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "organizer@iitb.edu", "Organizer@123")

	body := map[string]any{
		"title":       "Test Hack " + uuid.NewString()[:6],
		"category":    "HACKATHON",
		"description": "Test event",
		"venue":       "Lab 1",
		"city":        "Mumbai",
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	var res struct {
		Data struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Slug   string `json:"slug"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if res.Data.Status != "DRAFT" {
		t.Fatalf("expected DRAFT, got %s", res.Data.Status)
	}
	if res.Data.Slug == "" {
		t.Fatal("expected non-empty slug")
	}

	// Cleanup.
	app.gdb.Exec("DELETE FROM events WHERE id = ?", res.Data.ID)
}

// ---- Create (student blocked) ----

func TestEvent_StudentCannotCreate(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	body := map[string]any{
		"title":    "Should Fail",
		"category": "OTHER",
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

// ---- Lifecycle ----

func TestEvent_LifecyclePublish(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "organizer@iitb.edu", "Organizer@123")

	// Create a draft with all required fields.
	startsAt := time.Now().Add(48 * time.Hour)
	body := map[string]any{
		"title":       "Lifecycle Test " + uuid.NewString()[:6],
		"category":    "WORKSHOP",
		"description": "Testing lifecycle",
		"venue":       "Room 101",
		"city":        "Mumbai",
		"starts_at":   startsAt.Format(time.RFC3339),
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}

	var created struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	// Publish
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/events/"+created.Data.ID+"/publish", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("publish failed: %d %s", w2.Code, w2.Body.String())
	}

	var published struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w2.Body.Bytes(), &published)
	if published.Data.Status != "PUBLISHED" {
		t.Fatalf("expected PUBLISHED, got %s", published.Data.Status)
	}

	// Cleanup.
	app.gdb.Exec("DELETE FROM events WHERE id = ?", created.Data.ID)
}

// ---- Filter by category ----

func TestEvent_FilterByCategory(t *testing.T) {
	app := setupFullApp(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/events?category=CULTURAL", nil)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var res struct {
		Data []struct {
			Category string `json:"category"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	for _, e := range res.Data {
		if e.Category != "CULTURAL" {
			t.Fatalf("filter leaked wrong category: %s", e.Category)
		}
	}
}

// ---- Owner can edit own draft ----

func TestEvent_OwnerCanEditDraft(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "organizer@iitb.edu", "Organizer@123")

	// Get the seeded draft event.
	var draftID string
	app.gdb.Raw(`SELECT id FROM events WHERE slug = 'workshop-ai-2026' AND college_id = '11111111-1111-1111-1111-111111111111'`).Scan(&draftID)
	if draftID == "" {
		t.Skip("no draft event seed")
	}

	newTitle := "Updated Title " + uuid.NewString()[:6]
	bb, _ := json.Marshal(map[string]any{"title": newTitle})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/events/"+draftID, bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

// ---- Cannot edit someone else's event ----

func TestEvent_CannotEditOthersEvent(t *testing.T) {
	app := setupFullApp(t)

	// IITD admin logs in and tries to edit IITB draft.
	token := app.login(t, "admin@iitd.ac.in", "Admin@123")

	var draftID string
	app.gdb.Raw(`SELECT id FROM events WHERE slug = 'workshop-ai-2026' AND college_id = '11111111-1111-1111-1111-111111111111'`).Scan(&draftID)
	if draftID == "" {
		t.Skip("no draft event seed")
	}

	bb, _ := json.Marshal(map[string]any{"title": "Hacked"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/events/"+draftID, bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Fatalf("expected 403 for non-owner edit, got %d body=%s", w.Code, w.Body.String())
	}
}

// ---- Invalid lifecycle transition ----

func TestEvent_InvalidTransition(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "organizer@iitb.edu", "Organizer@123")

	// Complete a DRAFT event (invalid: must go DRAFT → PUBLISHED first)
	var draftID string
	app.gdb.Raw(`SELECT id FROM events WHERE slug = 'workshop-ai-2026'`).Scan(&draftID)
	if draftID == "" {
		t.Skip("no draft event seed")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events/"+draftID+"/complete", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 409 {
		t.Fatalf("expected 409 for invalid transition, got %d body=%s", w.Code, w.Body.String())
	}
}

// ---- Date validation ----

func TestEvent_InvalidDatesRejected(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "organizer@iitb.edu", "Organizer@123")

	startsAt := time.Now().Add(48 * time.Hour)
	endsAt := time.Now().Add(24 * time.Hour) // ends BEFORE starts

	body := map[string]any{
		"title":     "Bad Dates " + uuid.NewString()[:6],
		"category":  "WORKSHOP",
		"starts_at": startsAt.Format(time.RFC3339),
		"ends_at":   endsAt.Format(time.RFC3339),
	}
	bb, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400 for invalid dates, got %d body=%s", w.Code, w.Body.String())
	}
}
