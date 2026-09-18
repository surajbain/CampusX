package registration

import "time"

type RegisterEventRequest struct {
	Type string `json:"type" binding:"required,oneof=ATTEND PARTICIPATE"`
	Notes string `json:"notes" binding:"omitempty,max=1000"`
}

type UpdateRegistrationRequest struct {
	Status *string `json:"status,omitempty" binding:"omitempty,oneof=PENDING CONFIRMED CANCELLED WAITLISTED"`
	Notes  *string `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

type RegistrationResponse struct {
	ID        string     `json:"id"`
	EventID   string     `json:"event_id"`
	UserID    string     `json:"user_id"`
	CollegeID string     `json:"college_id"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	Notes     string     `json:"notes,omitempty"`

	// Event details (denormalized for convenience)
	EventTitle       string     `json:"event_title,omitempty"`
	EventSlug        string     `json:"event_slug,omitempty"`
	EventPosterURL   string     `json:"event_poster_url,omitempty"`
	EventCategory    string     `json:"event_category,omitempty"`
	EventStartsAt    *time.Time `json:"event_starts_at,omitempty"`
	EventEndsAt      *time.Time `json:"event_ends_at,omitempty"`
	EventVenue       string     `json:"event_venue,omitempty"`
	EventCity        string     `json:"event_city,omitempty"`
	EventPricePaise  int64      `json:"event_price_paise"`
	EventCurrency    string     `json:"event_currency,omitempty"`
	EventCollegeName string     `json:"event_college_name,omitempty"`

	// User details (for admin views)
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListRegistrationsResponse struct {
	Data []RegistrationResponse `json:"data"`
	Meta interface{}            `json:"meta"`
}