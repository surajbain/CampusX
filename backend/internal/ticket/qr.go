package ticket

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidQR      = errors.New("ticket: invalid QR")
	ErrQRSignatureBad = errors.New("ticket: QR signature mismatch")
	ErrQRExpired      = errors.New("ticket: QR expired")
	ErrQRNotYetValid  = errors.New("ticket: QR not yet valid")
)

// QRManager handles signing and verifying QR payloads.
type QRManager struct {
	secret []byte
}

func NewQRManager(secret string) *QRManager {
	return &QRManager{secret: []byte(secret)}
}

// Sign generates a signed QR payload for a ticket.
// Returns (payload JSON string, signature hex).
func (m *QRManager) Sign(p QRPayload) (string, string, error) {
	// Add nonce if missing
	if p.Nonce == "" {
		var b [16]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", "", fmt.Errorf("nonce: %w", err)
		}
		p.Nonce = hex.EncodeToString(b[:])
	}

	// Marshal payload
	body, err := json.Marshal(p)
	if err != nil {
		return "", "", fmt.Errorf("marshal: %w", err)
	}

	// Compute HMAC-SHA256 over payload
	sig := m.computeHMAC(body)

	// Encode payload as base64 (URL-safe, no padding)
	payloadB64 := base64.RawURLEncoding.EncodeToString(body)

	return payloadB64, sig, nil
}

// Verify checks a signed QR and returns the decoded payload.
// `validFrom` and `validUntil` are optional time windows.
func (m *QRManager) Verify(payloadB64, signature string, validFrom, validUntil *time.Time) (*QRPayload, error) {
	// Decode
	body, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, ErrInvalidQR
	}

	// Verify HMAC
	expectedSig := m.computeHMAC(body)
	if !hmac.Equal([]byte(expectedSig), []byte(signature)) {
		return nil, ErrQRSignatureBad
	}

	// Parse payload
	var p QRPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, ErrInvalidQR
	}

	// Time window check
	now := time.Now().UTC()
	if validFrom != nil && now.Before(*validFrom) {
		return nil, ErrQRNotYetValid
	}
	if validUntil != nil && now.After(*validUntil) {
		return nil, ErrQRExpired
	}

	return &p, nil
}

func (m *QRManager) computeHMAC(body []byte) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// GenerateTicketCode produces a human-readable code like "CX-A1B2-C3D4".
func GenerateTicketCode() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "CX-" + hex.EncodeToString(b[:2]) + "-" + hex.EncodeToString(b[2:])
}