package auth

import "time"

type RegisterRequest struct {
	Email     string `json:"email"     binding:"required,email,max=255"`
	Password  string `json:"password"  binding:"required,min=8,max=128"`
	FullName  string `json:"full_name" binding:"required,min=2,max=200"`
	Phone     string `json:"phone"     binding:"omitempty,max=30"`
	CollegeID string `json:"college_id" binding:"required,uuid"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=8,max=128"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone,omitempty"`
	Role          string    `json:"role"`
	CollegeID     string    `json:"college_id,omitempty"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"` // seconds
	TokenType    string       `json:"token_type"`
}
