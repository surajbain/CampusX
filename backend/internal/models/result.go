package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Result struct {
	BaseNoSoftDelete
	EventID     string           `gorm:"type:uuid;not null;index" json:"event_id"`
	CollegeID   string           `gorm:"type:uuid;not null" json:"college_id"`
	TeamID      *string          `gorm:"type:uuid;index" json:"team_id,omitempty"`
	UserID      *string          `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Position    int              `gorm:"not null" json:"position"`
	Title       string           `gorm:"size:200" json:"title,omitempty"`
	Score       *decimal.Decimal `gorm:"type:numeric(10,2)" json:"score,omitempty"`
	Remarks     string           `gorm:"type:text" json:"remarks,omitempty"`
	PublishedBy *string          `gorm:"type:uuid" json:"published_by,omitempty"`
	PublishedAt time.Time        `json:"published_at"`
}

func (Result) TableName() string { return "results" }

type Winner struct {
	BaseNoSoftDelete
	EventID        string    `gorm:"type:uuid;not null;index" json:"event_id"`
	ResultID       string    `gorm:"type:uuid;not null" json:"result_id"`
	UserID         string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Position       int       `gorm:"not null" json:"position"`
	CertificateURL string    `gorm:"type:text" json:"certificate_url,omitempty"`
	IssuedAt       time.Time `json:"issued_at"`
}

func (Winner) TableName() string { return "winners" }
