package repository

import (
	"context"
	"errors"
	"time"

	"cesizen/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id::text, email, password_hash, COALESCE(first_name, ''), COALESCE(nickname, ''),
		       is_active, role, login_attempts, locked_until, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.Nickname,
		&u.IsActive, &u.Role, &u.LoginAttempts, &u.LockedUntil,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, first_name, nickname)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, role, is_active, created_at, updated_at
	`, u.Email, u.PasswordHash, u.FirstName, u.Nickname,
	).Scan(&u.ID, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) UpdateLoginAttempts(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET login_attempts = $1, locked_until = $2, updated_at = now()
		WHERE id = $3::uuid
	`, attempts, lockedUntil, userID)
	return err
}

func (r *UserRepository) ResetLoginAttempts(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET login_attempts = 0, locked_until = NULL, updated_at = now()
		WHERE id = $1::uuid
	`, userID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id::text, email, COALESCE(first_name, ''), COALESCE(nickname, ''), is_active, role, created_at, updated_at
		FROM users WHERE id = $1::uuid
	`, id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.Nickname,
		&u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return u, err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET first_name = $1, nickname = $2, updated_at = now()
		WHERE id = $3::uuid
	`, u.FirstName, u.Nickname, u.ID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2::uuid
	`, hash, userID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1::uuid`, id)
	return err
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, email, COALESCE(first_name, ''), COALESCE(nickname, ''), is_active, role, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.Nickname,
			&u.IsActive, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) AdminUpdate(ctx context.Context, id, role string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET role = $1, is_active = $2, updated_at = now()
		WHERE id = $3::uuid
	`, role, isActive, id)
	return err
}
