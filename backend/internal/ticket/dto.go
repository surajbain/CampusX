package ticket

import "time"

type TicketResponse struct {
	ID             string     `json:"id"`
	RegistrationID string     `json:"registration_id"`
	UserID         string     `json:"user_id"`
	EventID        string     `json:"event_id"`
	TicketCode     string     `json:"ticket_code"`
	QRPayload      string     `json:"qr_payload"`
	QRSignature    string     `json:"qr_signature"`
	QRImageDataURL string     `json:"qr_image_data_url,omitempty"` // optional base64 PNG
	ValidFrom      *time.Time `json:"valid_from,omitempty"`
	ValidUntil     *time.Time `json:"valid_until,omitempty"`
	Status         string     `json:"status"`
	IssuedAt       time.Time  `json:"issued_at"`
	UsedAt         *time.Time `json:"used_at,omitempty"`

	// Event details (denormalized)
	EventTitle       string     `json:"event_title,omitempty"`
	EventSlug        string     `json:"event_slug,omitempty"`
	EventPosterURL   string     `json:"event_poster_url,omitempty"`
	EventCategory    string     `json:"event_category,omitempty"`
	EventStartsAt    *time.Time `json:"event_starts_at,omitempty"`
	EventEndsAt      *time.Time `json:"event_ends_at,omitempty"`
	EventVenue       string     `json:"event_venue,omitempty"`
	EventCity        string     `json:"event_city,omitempty"`
	EventCollegeName string     `json:"event_college_name,omitempty"`
	EventWhatsappLink string `json:"event_whatsapp_link,omitempty"`  

	// User details
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

// QRPayload is the data encoded inside the QR code (before signing).
type QRPayload struct {
	TicketID string `json:"tid"`
	UserID   string `json:"uid"`
	EventID  string `json:"eid"`
	IssuedAt int64  `json:"iat"`
	Nonce    string `json:"nonce"` // random per-ticket, prevents duplication
}