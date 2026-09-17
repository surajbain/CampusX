package models

import "time"

type CertificateKind string

const (
	CertParticipation CertificateKind = "PARTICIPATION"
	CertWinner        CertificateKind = "WINNER"
	CertRunnerUp      CertificateKind = "RUNNER_UP"
	CertSpecial       CertificateKind = "SPECIAL"
)

type Certificate struct {
	BaseNoSoftDelete
	UserID    string          `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID   string          `gorm:"type:uuid;not null;index" json:"event_id"`
	CollegeID string          `gorm:"type:uuid;not null" json:"college_id"`
	Kind      CertificateKind `gorm:"size:30;not null" json:"kind"`
	FileURL   string          `gorm:"type:text;not null" json:"file_url"`
	Serial    string          `gorm:"size:60;not null;uniqueIndex" json:"serial"`
	IssuedAt  time.Time       `json:"issued_at"`
}

func (Certificate) TableName() string { return "certificates" }
