package repository

import (
	"context"
	"errors"
	"time"

	"cesizen/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

// ── Refresh tokens ────────────────────────────────────────────────────────────

func (r *TokenRepository) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2::uuid, $3)
	`, token, userID, expiresAt)
	return err
}

func (r *TokenRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	rt := &domain.RefreshToken{}
	err := r.db.QueryRow(ctx, `
		SELECT token, user_id::text, expires_at FROM refresh_tokens WHERE token = $1
	`, token).Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return rt, err
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}

func (r *TokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1::uuid`, userID)
	return err
}

// ── Password reset tokens ─────────────────────────────────────────────────────

func (r *TokenRepository) CreatePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO password_reset_tokens (token, user_id, expires_at)
		VALUES ($1, $2::uuid, $3)
		ON CONFLICT (token) DO NOTHING
	`, token, userID, expiresAt)
	return err
}

func (r *TokenRepository) FindPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	prt := &domain.PasswordResetToken{}
	err := r.db.QueryRow(ctx, `
		SELECT token, user_id::text, expires_at FROM password_reset_tokens WHERE token = $1
	`, token).Scan(&prt.Token, &prt.UserID, &prt.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return prt, err
}

func (r *TokenRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM password_reset_tokens WHERE token = $1`, token)
	return err
}

func (r *TokenRepository) DeleteUserPasswordResetTokens(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM password_reset_tokens WHERE user_id = $1::uuid`, userID)
	return err
}
