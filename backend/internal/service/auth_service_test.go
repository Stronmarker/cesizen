package service_test

import (
	"context"
	"testing"
	"time"

	"cesizen/internal/domain"
	"cesizen/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implémente service.UserRepo avec des champs fonction pour contrôler le comportement par test.
type mockUserRepo struct {
	findByEmailFn         func(ctx context.Context, email string) (*domain.User, error)
	findByIDFn            func(ctx context.Context, id string) (*domain.User, error)
	createFn              func(ctx context.Context, u *domain.User) error
	updateLoginAttemptsFn func(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error
	resetLoginAttemptsFn  func(ctx context.Context, userID string) error
	updatePasswordFn      func(ctx context.Context, userID, hash string) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.findByEmailFn(ctx, email)
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	return m.createFn(ctx, u)
}
func (m *mockUserRepo) UpdateLoginAttempts(ctx context.Context, id string, attempts int, lock *time.Time) error {
	if m.updateLoginAttemptsFn != nil {
		return m.updateLoginAttemptsFn(ctx, id, attempts, lock)
	}
	return nil
}
func (m *mockUserRepo) ResetLoginAttempts(ctx context.Context, id string) error {
	if m.resetLoginAttemptsFn != nil {
		return m.resetLoginAttemptsFn(ctx, id)
	}
	return nil
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id, hash string) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, id, hash)
	}
	return nil
}

// noopTokenRepo satisfait service.TokenRepo sans rien stocker.
type noopTokenRepo struct{}

func (n *noopTokenRepo) CreateRefreshToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (n *noopTokenRepo) FindRefreshToken(_ context.Context, _ string) (*domain.RefreshToken, error) {
	return nil, nil
}
func (n *noopTokenRepo) DeleteRefreshToken(_ context.Context, _ string) error  { return nil }
func (n *noopTokenRepo) DeleteUserRefreshTokens(_ context.Context, _ string) error { return nil }
func (n *noopTokenRepo) CreatePasswordResetToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (n *noopTokenRepo) FindPasswordResetToken(_ context.Context, _ string) (*domain.PasswordResetToken, error) {
	return nil, nil
}
func (n *noopTokenRepo) DeletePasswordResetToken(_ context.Context, _ string) error        { return nil }
func (n *noopTokenRepo) DeleteUserPasswordResetTokens(_ context.Context, _ string) error { return nil }

func newService(repo service.UserRepo) *service.AuthService {
	return service.NewAuthService(repo, &noopTokenRepo{}, "test-secret")
}

func hashedPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	require.NoError(t, err)
	return string(h)
}

// ── Register ─────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(_ context.Context, u *domain.User) error {
			u.ID = "new-uuid"
			u.Role = "user"
			u.IsActive = true
			return nil
		},
	}

	resp, err := newService(repo).Register(context.Background(), domain.RegisterInput{
		Email:     "alice@example.com",
		Password:  "password123",
		FirstName: "Alice",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "alice@example.com", resp.User.Email)
	assert.Equal(t, "user", resp.User.Role)
}

func TestRegister_EmailTaken(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{Email: "alice@example.com"}, nil
		},
	}

	_, err := newService(repo).Register(context.Background(), domain.RegisterInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, domain.ErrEmailTaken)
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	resetCalled := false
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{
				ID:           "user-uuid",
				Email:        "alice@example.com",
				Role:         "user",
				PasswordHash: hashedPassword(t, "password123"),
			}, nil
		},
		resetLoginAttemptsFn: func(_ context.Context, _ string) error {
			resetCalled = true
			return nil
		},
	}

	resp, err := newService(repo).Login(context.Background(), domain.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.True(t, resetCalled, "login_attempts doit être remis à zéro après succès")
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{
				ID:           "user-uuid",
				PasswordHash: hashedPassword(t, "password123"),
			}, nil
		},
	}

	_, err := newService(repo).Login(context.Background(), domain.LoginInput{
		Email:    "alice@example.com",
		Password: "wrong",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}

	_, err := newService(repo).Login(context.Background(), domain.LoginInput{
		Email:    "nobody@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_LockedAccount(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{
				ID:           "user-uuid",
				PasswordHash: hashedPassword(t, "password123"),
				LockedUntil:  &future,
			}, nil
		},
	}

	_, err := newService(repo).Login(context.Background(), domain.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	assert.ErrorIs(t, err, domain.ErrAccountLocked)
}

func TestLogin_LockAfterThreeFailures(t *testing.T) {
	attempts := 0
	var capturedLock *time.Time

	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{
				ID:            "user-uuid",
				PasswordHash:  hashedPassword(t, "password123"),
				LoginAttempts: attempts,
			}, nil
		},
		updateLoginAttemptsFn: func(_ context.Context, _ string, a int, lock *time.Time) error {
			attempts = a
			capturedLock = lock
			return nil
		},
	}

	svc := newService(repo)
	ctx := context.Background()
	badCreds := domain.LoginInput{Email: "alice@example.com", Password: "wrong"}

	_, err := svc.Login(ctx, badCreds)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	assert.Nil(t, capturedLock)

	_, err = svc.Login(ctx, badCreds)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	assert.Nil(t, capturedLock)

	_, err = svc.Login(ctx, badCreds)
	assert.ErrorIs(t, err, domain.ErrAccountLocked, "3ème échec doit verrouiller le compte")
	assert.NotNil(t, capturedLock, "locked_until doit être défini après 3 échecs")
}
