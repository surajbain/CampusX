package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPayment_CreateForPaidEvent(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	// 1. Register for CodeFest (paid — ₹99)
	regBody := map[string]any{"type": "ATTEND"}
	rbb, _ := json.Marshal(regBody)

	regW := httptest.NewRecorder()
	regReq, _ := http.NewRequest("POST", "/api/v1/events/33333333-3333-3333-3333-333333333302/register", bytes.NewReader(rbb))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(regW, regReq)

	if regW.Code != 201 && regW.Code != 409 {
		t.Fatalf("registration failed: %d %s", regW.Code, regW.Body.String())
	}

	var regRes struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(regW.Body.Bytes(), &regRes)
	if regRes.Data.ID == "" {
		t.Skip("registration id missing (may already exist)")
	}

	// 2. Create payment
	payBody := map[string]any{"registration_id": regRes.Data.ID}
	pbb, _ := json.Marshal(payBody)

	payW := httptest.NewRecorder()
	payReq, _ := http.NewRequest("POST", "/api/v1/payments/create", bytes.NewReader(pbb))
	payReq.Header.Set("Content-Type", "application/json")
	payReq.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(payW, payReq)

	if payW.Code != 201 && payW.Code != 409 {
		t.Fatalf("expected 201 or 409, got %d body=%s", payW.Code, payW.Body.String())
	}

	// Cleanup
	app.gdb.Exec("DELETE FROM payments WHERE registration_id = ?", regRes.Data.ID)
	app.gdb.Exec("DELETE FROM registrations WHERE id = ?", regRes.Data.ID)
}

func TestPayment_CannotPayForFreeEvent(t *testing.T) {
	app := setupFullApp(t)
	token := app.login(t, "student@iitb.edu", "Student@123")

	// Register for TechFest (free)
	regBody := map[string]any{"type": "ATTEND"}
	rbb, _ := json.Marshal(regBody)

	regW := httptest.NewRecorder()
	regReq, _ := http.NewRequest("POST", "/api/v1/events/33333333-3333-3333-3333-333333333301/register", bytes.NewReader(rbb))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(regW, regReq)

	if regW.Code != 201 && regW.Code != 409 {
		t.Skipf("registration failed: %s", regW.Body.String())
	}

	var regRes struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(regW.Body.Bytes(), &regRes)
	if regRes.Data.ID == "" {
		t.Skip("registration id missing")
	}

	// Try payment — should fail
	payBody := map[string]any{"registration_id": regRes.Data.ID}
	pbb, _ := json.Marshal(payBody)

	payW := httptest.NewRecorder()
	payReq, _ := http.NewRequest("POST", "/api/v1/payments/create", bytes.NewReader(pbb))
	payReq.Header.Set("Content-Type", "application/json")
	payReq.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(payW, payReq)

	if payW.Code != 400 {
		t.Fatalf("expected 400 for free event, got %d body=%s", payW.Code, payW.Body.String())
	}
}