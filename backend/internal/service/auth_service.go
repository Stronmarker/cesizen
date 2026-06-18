package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"cesizen/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const maxLoginAttempts = 3
const lockDuration = 15 * time.Minute
const accessTokenDuration = 15 * time.Minute
const refreshTokenDuration = 7 * 24 * time.Hour
const resetTokenDuration = 1 * time.Hour

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	UpdateLoginAttempts(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error
	ResetLoginAttempts(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID, hash string) error
}

type TokenRepo interface {
	CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
	CreatePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	FindPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	DeleteUserPasswordResetTokens(ctx context.Context, userID string) error
}

type AuthService struct {
	users     UserRepo
	tokens    TokenRepo
	jwtSecret []byte
}

func NewAuthService(users UserRepo, tokens TokenRepo, jwtSecret string) *AuthService {
	return &AuthService{users: users, tokens: tokens, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, input domain.RegisterInput) (*domain.AuthResponse, error) {
	_, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, domain.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		FirstName:    input.FirstName,
		Nickname:     input.Nickname,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, input.Email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, domain.ErrAccountLocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		attempts := user.LoginAttempts + 1
		var lockUntil *time.Time
		if attempts >= maxLoginAttempts {
			t := time.Now().Add(lockDuration)
			lockUntil = &t
		}
		_ = s.users.UpdateLoginAttempts(ctx, user.ID, attempts, lockUntil)
		if lockUntil != nil {
			return nil, domain.ErrAccountLocked
		}
		return nil, domain.ErrInvalidCredentials
	}

	_ = s.users.ResetLoginAttempts(ctx, user.ID)
	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	rt, err := s.tokens.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil || rt.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrTokenExpired
	}

	user, err := s.users.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Rotation: delete old token, issue new pair
	_ = s.tokens.DeleteRefreshToken(ctx, refreshToken)
	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokens.DeleteRefreshToken(ctx, refreshToken)
}

func (s *AuthService) ForgotPassword(ctx context.Context, input domain.ForgotPasswordInput) (*domain.ForgotPasswordResponse, error) {
	user, err := s.users.FindByEmail(ctx, input.Email)
	if errors.Is(err, domain.ErrNotFound) {
		// Don't reveal if email exists
		return &domain.ForgotPasswordResponse{Message: "Si cet email existe, un lien de réinitialisation a été envoyé."}, nil
	}
	if err != nil {
		return nil, err
	}

	// Invalidate previous reset tokens for this user
	_ = s.tokens.DeleteUserPasswordResetTokens(ctx, user.ID)

	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	if err := s.tokens.CreatePasswordResetToken(ctx, user.ID, token, time.Now().Add(resetTokenDuration)); err != nil {
		return nil, err
	}

	return &domain.ForgotPasswordResponse{
		ResetToken: token,
		Message:    "Token de réinitialisation généré (en production, envoyé par email).",
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	prt, err := s.tokens.FindPasswordResetToken(ctx, input.Token)
	if err != nil {
		return err
	}
	if prt == nil || prt.ExpiresAt.Before(time.Now()) {
		return domain.ErrTokenExpired
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, prt.UserID, string(hash)); err != nil {
		return err
	}

	_ = s.tokens.DeletePasswordResetToken(ctx, input.Token)
	return nil
}

func (s *AuthService) buildAuthResponse(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken()
	if err != nil {
		return nil, err
	}
	if err := s.tokens.CreateRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(refreshTokenDuration)); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         toPublicUser(user),
	}, nil
}

func (s *AuthService) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(accessTokenDuration).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func toPublicUser(u *domain.User) domain.PublicUser {
	return domain.PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		Nickname:  u.Nickname,
		Role:      u.Role,
	}
}
