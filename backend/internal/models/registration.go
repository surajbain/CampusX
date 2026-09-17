package models

type RegistrationType string

const (
	RegTypeAttend      RegistrationType = "ATTEND"
	RegTypeParticipate RegistrationType = "PARTICIPATE"
)

type RegistrationStatus string

const (
	RegStatusPending    RegistrationStatus = "PENDING"
	RegStatusConfirmed  RegistrationStatus = "CONFIRMED"
	RegStatusCancelled  RegistrationStatus = "CANCELLED"
	RegStatusWaitlisted RegistrationStatus = "WAITLISTED"
)

type Registration struct {
	Base
	EventID   string             `gorm:"type:uuid;not null;index" json:"event_id"`
	UserID    string             `gorm:"type:uuid;not null;index" json:"user_id"`
	CollegeID string             `gorm:"type:uuid;not null" json:"college_id"`
	Type      RegistrationType   `gorm:"size:20;not null" json:"type"`
	Status    RegistrationStatus `gorm:"size:20;not null;default:'PENDING';index" json:"status"`
	TeamID    *string            `gorm:"type:uuid" json:"team_id,omitempty"`
	Notes     string             `gorm:"type:text" json:"notes,omitempty"`

	Event *Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Registration) TableName() string { return "registrations" }
