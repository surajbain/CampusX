package models

type College struct {
	Base
	Name         string  `gorm:"size:200;not null" json:"name"`
	Slug         string  `gorm:"size:120;not null;uniqueIndex" json:"slug"`
	City         string  `gorm:"size:120;not null;index" json:"city"`
	State        string  `gorm:"size:120" json:"state,omitempty"`
	LogoURL      string  `gorm:"type:text" json:"logo_url,omitempty"`
	Website      string  `gorm:"type:text" json:"website,omitempty"`
	ContactEmail string  `gorm:"size:200" json:"contact_email,omitempty"`
	ContactPhone string  `gorm:"size:30" json:"contact_phone,omitempty"`
	IsActive     bool    `gorm:"not null;default:true" json:"is_active"`
	Users        []User  `gorm:"foreignKey:CollegeID" json:"-"`
	Events       []Event `gorm:"foreignKey:CollegeID" json:"-"`
}

func (College) TableName() string { return "colleges" }
