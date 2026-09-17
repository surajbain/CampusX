package tests

import (
	"testing"

	"github.com/google/uuid"
)

func TestConstraint_UserEmailUnique(t *testing.T) {
	gdb := requireDB(t)

	email := "dup-" + uuid.NewString() + "@test.dev"
	// First insert
	err := gdb.Exec(
		`INSERT INTO users (email, password_hash, full_name, role)
		 VALUES (?, 'x', 'A', 'STUDENT')`, email,
	).Error
	if err != nil {
		t.Fatalf("first insert failed: %v", err)
	}
	// Second insert with same email should fail
	err = gdb.Exec(
		`INSERT INTO users (email, password_hash, full_name, role)
		 VALUES (?, 'x', 'B', 'STUDENT')`, email,
	).Error
	if err == nil {
		t.Fatal("expected unique violation, got nil")
	}
	// Cleanup
	gdb.Exec("DELETE FROM users WHERE email = ?", email)
}

func TestConstraint_PaymentGatewayIDUnique(t *testing.T) {
	gdb := requireDB(t)

	// Need a valid registration. Skip if seed data isn't in place.
	var regID string
	if err := gdb.Raw("SELECT id FROM registrations LIMIT 1").Scan(&regID).Error; err != nil || regID == "" {
		t.Skip("no registration row available — seed or run-time setup required")
	}

	gwID := "gw_" + uuid.NewString()
	idem1 := "idem1_" + uuid.NewString()
	idem2 := "idem2_" + uuid.NewString()

	// Insert first payment
	err := gdb.Exec(
		`INSERT INTO payments (registration_id, user_id, event_id, gateway_payment_id, idempotency_key, amount_paise)
		 SELECT ?, user_id, event_id, ?, ?, 100
		 FROM registrations WHERE id = ?`,
		regID, gwID, idem1, regID,
	).Error
	if err != nil {
		t.Fatalf("first payment insert failed: %v", err)
	}

	// Same gateway_payment_id, different idempotency_key → should fail
	err = gdb.Exec(
		`INSERT INTO payments (registration_id, user_id, event_id, gateway_payment_id, idempotency_key, amount_paise)
		 SELECT ?, user_id, event_id, ?, ?, 100
		 FROM registrations WHERE id = ?`,
		regID, gwID, idem2, regID,
	).Error
	if err == nil {
		t.Fatal("expected gateway_payment_id unique violation, got nil")
	}

	gdb.Exec("DELETE FROM payments WHERE gateway_payment_id = ?", gwID)
}

func TestConstraint_RegistrationUniquePerType(t *testing.T) {
	gdb := requireDB(t)

	var eventID, userID, collegeID string
	err := gdb.Raw(`
		SELECT e.id, u.id, e.college_id
		FROM events e, users u
		WHERE u.role = 'STUDENT'
		LIMIT 1
	`).Row().Scan(&eventID, &userID, &collegeID)
	if err != nil {
		t.Skip("seed data unavailable for unique test")
	}

	// Insert first registration
	err = gdb.Exec(
		`INSERT INTO registrations (event_id, user_id, college_id, type)
		 VALUES (?, ?, ?, 'PARTICIPATE')`,
		eventID, userID, collegeID,
	).Error
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	// Same (event_id, user_id, type) → fail
	err = gdb.Exec(
		`INSERT INTO registrations (event_id, user_id, college_id, type)
		 VALUES (?, ?, ?, 'PARTICIPATE')`,
		eventID, userID, collegeID,
	).Error
	if err == nil {
		t.Fatal("expected unique violation on (event_id, user_id, type), got nil")
	}

	// Cleanup
	gdb.Exec("DELETE FROM registrations WHERE event_id = ? AND user_id = ?", eventID, userID)
}

func TestConstraint_CheckinTicketUnique(t *testing.T) {
	gdb := requireDB(t)

	var ticketID, eventID, collegeID string
	err := gdb.Raw(`
		SELECT t.id, t.event_id, e.college_id
		FROM tickets t JOIN events e ON e.id = t.event_id
		LIMIT 1
	`).Row().Scan(&ticketID, &eventID, &collegeID)
	if err != nil {
		t.Skip("no ticket row available")
	}

	err = gdb.Exec(
		`INSERT INTO checkins (ticket_id, event_id, college_id, result)
		 VALUES (?, ?, ?, 'SUCCESS')`,
		ticketID, eventID, collegeID,
	).Error
	if err != nil {
		t.Fatalf("first checkin failed: %v", err)
	}

	err = gdb.Exec(
		`INSERT INTO checkins (ticket_id, event_id, college_id, result)
		 VALUES (?, ?, ?, 'SUCCESS')`,
		ticketID, eventID, collegeID,
	).Error
	if err == nil {
		t.Fatal("expected duplicate checkin violation, got nil")
	}

	gdb.Exec("DELETE FROM checkins WHERE ticket_id = ?", ticketID)
}
