package models

import "time"

type CheckinResult string

const (
	CheckinSuccess   CheckinResult = "SUCCESS"
	CheckinDuplicate CheckinResult = "DUPLICATE"
	CheckinInvalid   CheckinResult = "INVALID"
	CheckinExpired   CheckinResult = "EXPIRED"
)

type Checkin struct {
	BaseNoSoftDelete
	TicketID  string        `gorm:"type:uuid;not null;uniqueIndex" json:"ticket_id"`
	ScannerID *string       `gorm:"type:uuid;index" json:"scanner_id,omitempty"`
	EventID   string        `gorm:"type:uuid;not null;index" json:"event_id"`
	CollegeID string        `gorm:"type:uuid;not null" json:"college_id"`
	Result    CheckinResult `gorm:"size:20;not null" json:"result"`
	Reason    string        `gorm:"type:text" json:"reason,omitempty"`
	ScannedAt time.Time     `gorm:"index" json:"scanned_at"`
}

func (Checkin) TableName() string { return "checkins" }
