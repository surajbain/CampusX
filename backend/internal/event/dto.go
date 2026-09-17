package event

import (
	"time"

	"gorm.io/datatypes"
)

type CreateEventRequest struct {
	Title              string         `json:"title"                binding:"required,min=3,max=250"`
	Slug               string         `json:"slug"                 binding:"omitempty,min=2,max=160"`
	Description        string         `json:"description"          binding:"omitempty,max=20000"`
	Category           string         `json:"category"             binding:"required,oneof=HACKATHON CULTURAL SPORTS WORKSHOP TECH_FEST OTHER"`
	PosterURL          string         `json:"poster_url"           binding:"omitempty,max=500"`
	Venue              string         `json:"venue"                binding:"omitempty,max=250"`
	City               string         `json:"city"                 binding:"omitempty,max=120"`
	StartsAt           *time.Time     `json:"starts_at"`
	EndsAt             *time.Time     `json:"ends_at"`
	RegistrationOpens  *time.Time     `json:"registration_opens"`
	RegistrationCloses *time.Time     `json:"registration_closes"`
	PricePaise         int64          `json:"price_paise"          binding:"min=0"`
	Capacity           *int           `json:"capacity"             binding:"omitempty,min=1"`
	AllowTeams         bool           `json:"allow_teams"`
	TeamSizeMin        *int           `json:"team_size_min"        binding:"omitempty,min=1,max=50"`
	TeamSizeMax        *int           `json:"team_size_max"        binding:"omitempty,min=1,max=50"`
	PrizePoolPaise     int64          `json:"prize_pool_paise"     binding:"min=0"`
	ContactEmail       string         `json:"contact_email"        binding:"omitempty,email,max=200"`
	WhatsappLink       string         `json:"whatsapp_link"        binding:"omitempty,max=500"`
	Rules              string         `json:"rules"                binding:"omitempty,max=20000"`
	ScheduleJSON       datatypes.JSON `json:"schedule"`
	FAQJSON            datatypes.JSON `json:"faq"`
	IsFeatured         bool           `json:"is_featured"`
}

type UpdateEventRequest struct {
	Title              *string    `json:"title,omitempty"                binding:"omitempty,min=3,max=250"`
	Description        *string    `json:"description,omitempty"          binding:"omitempty,max=20000"`
	Category           *string    `json:"category,omitempty"             binding:"omitempty,oneof=HACKATHON CULTURAL SPORTS WORKSHOP TECH_FEST OTHER"`
	PosterURL          *string    `json:"poster_url,omitempty"           binding:"omitempty,max=500"`
	Venue              *string    `json:"venue,omitempty"                binding:"omitempty,max=250"`
	City               *string    `json:"city,omitempty"                 binding:"omitempty,max=120"`
	StartsAt           *time.Time `json:"starts_at,omitempty"`
	EndsAt             *time.Time `json:"ends_at,omitempty"`
	RegistrationOpens  *time.Time `json:"registration_opens,omitempty"`
	RegistrationCloses *time.Time `json:"registration_closes,omitempty"`
	PricePaise         *int64     `json:"price_paise,omitempty"          binding:"omitempty,min=0"`
	Capacity           *int       `json:"capacity,omitempty"             binding:"omitempty,min=1"`
	AllowTeams         *bool      `json:"allow_teams,omitempty"`
	TeamSizeMin        *int       `json:"team_size_min,omitempty"        binding:"omitempty,min=1,max=50"`
	TeamSizeMax        *int       `json:"team_size_max,omitempty"        binding:"omitempty,min=1,max=50"`
	PrizePoolPaise     *int64     `json:"prize_pool_paise,omitempty"     binding:"omitempty,min=0"`
	ContactEmail       *string    `json:"contact_email,omitempty"        binding:"omitempty,email,max=200"`
	WhatsappLink       *string    `json:"whatsapp_link,omitempty"        binding:"omitempty,max=500"`
	Rules              *string    `json:"rules,omitempty"                binding:"omitempty,max=20000"`
	IsFeatured         *bool      `json:"is_featured,omitempty"`
}

type EventResponse struct {
	ID                 string         `json:"id"`
	CollegeID          string         `json:"college_id"`
	CreatedBy          string         `json:"created_by,omitempty"`
	Slug               string         `json:"slug"`
	Title              string         `json:"title"`
	Description        string         `json:"description,omitempty"`
	Category           string         `json:"category"`
	PosterURL          string         `json:"poster_url,omitempty"`
	Venue              string         `json:"venue,omitempty"`
	City               string         `json:"city,omitempty"`
	Status             string         `json:"status"`
	StartsAt           *time.Time     `json:"starts_at,omitempty"`
	EndsAt             *time.Time     `json:"ends_at,omitempty"`
	RegistrationOpens  *time.Time     `json:"registration_opens,omitempty"`
	RegistrationCloses *time.Time     `json:"registration_closes,omitempty"`
	PricePaise         int64          `json:"price_paise"`
	Currency           string         `json:"currency"`
	Capacity           *int           `json:"capacity,omitempty"`
	AllowTeams         bool           `json:"allow_teams"`
	TeamSizeMin        *int           `json:"team_size_min,omitempty"`
	TeamSizeMax        *int           `json:"team_size_max,omitempty"`
	PrizePoolPaise     int64          `json:"prize_pool_paise"`
	ContactEmail       string         `json:"contact_email,omitempty"`
	WhatsappLink       string         `json:"whatsapp_link,omitempty"`
	Rules              string         `json:"rules,omitempty"`
	Schedule           datatypes.JSON `json:"schedule,omitempty"`
	FAQ                datatypes.JSON `json:"faq,omitempty"`
	IsFeatured         bool           `json:"is_featured"`

	// Embedded college info (for public pages).
	CollegeName string `json:"college_name,omitempty"`
	CollegeSlug string `json:"college_slug,omitempty"`
	CollegeCity string `json:"college_city,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
