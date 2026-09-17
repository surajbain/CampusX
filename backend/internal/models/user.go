package models

import "time"

type Role string

const (
	RoleStudent      Role = "STUDENT"
	RoleCollegeAdmin Role = "COLLEGE_ADMIN"
	RoleOrganizer    Role = "ORGANIZER"
	RoleVolunteer    Role = "VOLUNTEER"
	RoleSuperAdmin   Role = "SUPER_ADMIN"
	RoleAdvertiser   Role = "ADVERTISER"
)

type User struct {
	Base
	Email         string     `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash  string     `gorm:"type:text;not null" json:"-"`
	FullName      string     `gorm:"size:200;not null" json:"full_name"`
	Phone         string     `gorm:"size:30" json:"phone,omitempty"`
	Role          Role       `gorm:"size:30;not null;index" json:"role"`
	CollegeID     *string    `gorm:"type:uuid;index" json:"college_id,omitempty"`
	IsActive      bool       `gorm:"not null;default:true" json:"is_active"`
	EmailVerified bool       `gorm:"not null;default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`

	College *College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
}

func (User) TableName() string { return "users" }
