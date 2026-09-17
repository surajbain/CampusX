package models

import (
	"time"

	"gorm.io/datatypes"
)

type AdSlotType string

const (
	SlotHomeBanner  AdSlotType = "HOME_BANNER"
	SlotSidebar     AdSlotType = "SIDEBAR"
	SlotEventFooter AdSlotType = "EVENT_FOOTER"
)

type Ad struct {
	Base
	AdvertiserID     *string    `gorm:"type:uuid" json:"advertiser_id,omitempty"`
	CollegeID        *string    `gorm:"type:uuid" json:"college_id,omitempty"`
	Title            string     `gorm:"size:200;not null" json:"title"`
	Body             string     `gorm:"type:text" json:"body,omitempty"`
	ImageURL         string     `gorm:"type:text" json:"image_url,omitempty"`
	TargetURL        string     `gorm:"type:text;not null" json:"target_url"`
	SlotType         AdSlotType `gorm:"size:30;not null;index" json:"slot_type"`
	StartsAt         *time.Time `json:"starts_at,omitempty"`
	EndsAt           *time.Time `json:"ends_at,omitempty"`
	IsActive         bool       `gorm:"not null;default:true;index" json:"is_active"`
	DailyBudgetPaise int64      `gorm:"not null;default:0" json:"daily_budget_paise"`
}

func (Ad) TableName() string { return "ads" }

type AdSlot struct {
	BaseNoSoftDelete
	Code        string `gorm:"size:40;not null;uniqueIndex" json:"code"`
	Description string `gorm:"size:200" json:"description,omitempty"`
	MaxActive   int    `gorm:"not null;default:1" json:"max_active"`
}

func (AdSlot) TableName() string { return "ad_slots" }

type AdImpression struct {
	BaseNoSoftDelete
	AdID     string         `gorm:"type:uuid;not null;index" json:"ad_id"`
	UserID   *string        `gorm:"type:uuid" json:"user_id,omitempty"`
	SlotCode string         `gorm:"size:40;not null" json:"slot_code"`
	Clicked  bool           `gorm:"not null;default:false" json:"clicked"`
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
}

func (AdImpression) TableName() string { return "ad_impressions" }
