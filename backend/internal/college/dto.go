package college

import "time"

type CreateCollegeRequest struct {
	Name         string `json:"name"          binding:"required,min=2,max=200"`
	Slug         string `json:"slug"          binding:"required,min=2,max=120,alphanum|contains=-"`
	City         string `json:"city"          binding:"required,min=2,max=120"`
	State        string `json:"state"         binding:"omitempty,max=120"`
	LogoURL      string `json:"logo_url"      binding:"omitempty,url,max=500"`
	Website      string `json:"website"       binding:"omitempty,url,max=500"`
	ContactEmail string `json:"contact_email" binding:"omitempty,email,max=200"`
	ContactPhone string `json:"contact_phone" binding:"omitempty,max=30"`
}

type UpdateCollegeRequest struct {
	Name         *string `json:"name,omitempty"          binding:"omitempty,min=2,max=200"`
	City         *string `json:"city,omitempty"          binding:"omitempty,min=2,max=120"`
	State        *string `json:"state,omitempty"         binding:"omitempty,max=120"`
	LogoURL      *string `json:"logo_url,omitempty"      binding:"omitempty,max=500"`
	Website      *string `json:"website,omitempty"       binding:"omitempty,max=500"`
	ContactEmail *string `json:"contact_email,omitempty" binding:"omitempty,email,max=200"`
	ContactPhone *string `json:"contact_phone,omitempty" binding:"omitempty,max=30"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type CollegeResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	City         string    `json:"city"`
	State        string    `json:"state,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	Website      string    `json:"website,omitempty"`
	ContactEmail string    `json:"contact_email,omitempty"`
	ContactPhone string    `json:"contact_phone,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ---- Admin user management ----

type CreateUserRequest struct {
	Email    string `json:"email"     binding:"required,email,max=255"`
	Password string `json:"password"  binding:"required,min=8,max=128"`
	FullName string `json:"full_name" binding:"required,min=2,max=200"`
	Phone    string `json:"phone"     binding:"omitempty,max=30"`
	Role     string `json:"role"      binding:"required,oneof=ORGANIZER VOLUNTEER COLLEGE_ADMIN"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty" binding:"omitempty,min=2,max=200"`
	Phone    *string `json:"phone,omitempty"     binding:"omitempty,max=30"`
	Role     *string `json:"role,omitempty"      binding:"omitempty,oneof=ORGANIZER VOLUNTEER COLLEGE_ADMIN STUDENT"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AdminUserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name"`
	Phone         string     `json:"phone,omitempty"`
	Role          string     `json:"role"`
	CollegeID     string     `json:"college_id,omitempty"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
