package models

import "time"

type TicketStatus string

const (
	TicketActive  TicketStatus = "ACTIVE"
	TicketUsed    TicketStatus = "USED"
	TicketRevoked TicketStatus = "REVOKED"
	TicketExpired TicketStatus = "EXPIRED"
)

type Ticket struct {
	Base
	RegistrationID string `gorm:"type:uuid;not null;uniqueIndex" json:"registration_id"`
	UserID         string `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID        string `gorm:"type:uuid;not null;index" json:"event_id"`

	TicketCode  string       `gorm:"size:40;not null;uniqueIndex" json:"ticket_code"`
	QRPayload   string       `gorm:"type:text;not null" json:"-"`
	QRSignature string       `gorm:"type:text;not null" json:"-"`
	ValidFrom   *time.Time   `json:"valid_from,omitempty"`
	ValidUntil  *time.Time   `json:"valid_until,omitempty"`
	Status      TicketStatus `gorm:"size:20;not null;default:'ACTIVE';index" json:"status"`
	IssuedAt    time.Time    `json:"issued_at"`
	UsedAt      *time.Time   `json:"used_at,omitempty"`

	// ---- Relations ----
	Registration *Registration `gorm:"foreignKey:RegistrationID" json:"registration,omitempty"`
}

func (Ticket) TableName() string { return "tickets" }