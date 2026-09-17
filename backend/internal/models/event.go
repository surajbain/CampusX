package models

import (
	"time"

	"gorm.io/datatypes"
)

type EventCategory string

const (
	CatHackathon EventCategory = "HACKATHON"
	CatCultural  EventCategory = "CULTURAL"
	CatSports    EventCategory = "SPORTS"
	CatWorkshop  EventCategory = "WORKSHOP"
	CatTechFest  EventCategory = "TECH_FEST"
	CatOther     EventCategory = "OTHER"
)

type EventStatus string

const (
	StatusDraft     EventStatus = "DRAFT"
	StatusPublished EventStatus = "PUBLISHED"
	StatusOngoing   EventStatus = "ONGOING"
	StatusCompleted EventStatus = "COMPLETED"
	StatusCancelled EventStatus = "CANCELLED"
)

type Event struct {
	Base
	CollegeID         string         `gorm:"type:uuid;not null;index" json:"college_id"`
	CreatedBy         *string        `gorm:"type:uuid" json:"created_by,omitempty"`
	Slug              string         `gorm:"size:160;not null" json:"slug"`
	Title             string         `gorm:"size:250;not null" json:"title"`
	Description       string         `gorm:"type:text" json:"description,omitempty"`
	Category          EventCategory  `gorm:"size:40;not null;index" json:"category"`
	PosterURL         string         `gorm:"type:text" json:"poster_url,omitempty"`
	Venue             string         `gorm:"size:250" json:"venue,omitempty"`
	City              string         `gorm:"size:120;index" json:"city,omitempty"`
	Status            EventStatus    `gorm:"size:30;not null;default:'DRAFT';index" json:"status"`
	StartsAt          *time.Time     `gorm:"index" json:"starts_at,omitempty"`
	EndsAt            *time.Time     `json:"ends_at,omitempty"`
	RegistrationOpens *time.Time     `gorm:"column:registration_opens" json:"registration_opens,omitempty"`
	RegistrationClose *time.Time     `gorm:"column:registration_closes" json:"registration_closes,omitempty"`
	PricePaise        int64          `gorm:"not null;default:0" json:"price_paise"`
	Currency          string         `gorm:"size:8;not null;default:'INR'" json:"currency"`
	Capacity          *int           `json:"capacity,omitempty"`
	AllowTeams        bool           `gorm:"not null;default:false" json:"allow_teams"`
	TeamSizeMin       *int           `json:"team_size_min,omitempty"`
	TeamSizeMax       *int           `json:"team_size_max,omitempty"`
	PrizePoolPaise    int64          `gorm:"not null;default:0" json:"prize_pool_paise"`
	ContactEmail      string         `gorm:"size:200" json:"contact_email,omitempty"`
	WhatsappLink      string         `gorm:"type:text" json:"whatsapp_link,omitempty"`
	Rules             string         `gorm:"type:text" json:"rules,omitempty"`
	ScheduleJSON      datatypes.JSON `gorm:"type:jsonb" json:"schedule,omitempty"`
	FAQJSON           datatypes.JSON `gorm:"type:jsonb" json:"faq,omitempty"`
	IsFeatured        bool           `gorm:"not null;default:false" json:"is_featured"`

	College *College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
}

func (Event) TableName() string { return "events" }
