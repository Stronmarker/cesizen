package domain

import (
	"errors"
	"time"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
)

type User struct {
	ID            string
	Email         string
	PasswordHash  string
	FirstName     string
	Nickname      string
	IsActive      bool
	Role          string
	LoginAttempts int
	LockedUntil   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	Nickname  string `json:"nickname,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	User         PublicUser `json:"user"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	ResetToken string `json:"reset_token"`
	Message    string `json:"message"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

var ErrTokenExpired = errors.New("token expired or invalid")

type PublicUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Nickname  string `json:"nickname"`
	Role      string `json:"role"`
}

type UpdateProfileInput struct {
	FirstName string `json:"first_name"`
	Nickname  string `json:"nickname"`
}

type AdminUserView struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	Nickname  string    `json:"nickname"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUpdateUserInput struct {
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}
