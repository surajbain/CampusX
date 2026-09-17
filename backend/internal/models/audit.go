package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditSeverity string

const (
	SeverityInfo     AuditSeverity = "INFO"
	SeverityWarn     AuditSeverity = "WARN"
	SeverityCritical AuditSeverity = "CRITICAL"
)

type AuditLog struct {
	BaseNoSoftDelete
	ActorID    *string        `gorm:"type:uuid;index" json:"actor_id,omitempty"`
	CollegeID  *string        `gorm:"type:uuid" json:"college_id,omitempty"`
	Action     string         `gorm:"size:80;not null;index" json:"action"`
	EntityType string         `gorm:"size:60" json:"entity_type,omitempty"`
	EntityID   *string        `gorm:"type:uuid" json:"entity_id,omitempty"`
	Severity   AuditSeverity  `gorm:"size:20;not null;default:'INFO';index" json:"severity"`
	Metadata   datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress  string         `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
