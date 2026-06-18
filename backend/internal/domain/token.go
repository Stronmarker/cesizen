package domain

import "time"

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type PasswordResetToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}
