package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cesizen/internal/domain"
	"cesizen/internal/handler"
	"cesizen/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo local au package handler_test.
type mockUserRepo struct {
	findByEmailFn         func(ctx context.Context, email string) (*domain.User, error)
	createFn              func(ctx context.Context, u *domain.User) error
	updateLoginAttemptsFn func(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.findByEmailFn(ctx, email)
}
func (m *mockUserRepo) FindByID(_ context.Context, _ string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	return m.createFn(ctx, u)
}
func (m *mockUserRepo) UpdateLoginAttempts(ctx context.Context, id string, a int, lock *time.Time) error {
	if m.updateLoginAttemptsFn != nil {
		return m.updateLoginAttemptsFn(ctx, id, a, lock)
	}
	return nil
}
func (m *mockUserRepo) ResetLoginAttempts(_ context.Context, _ string) error  { return nil }
func (m *mockUserRepo) UpdatePassword(_ context.Context, _, _ string) error   { return nil }

// noopTokenRepo satisfait service.TokenRepo dans les tests handler.
type noopTokenRepo struct{}

func (n *noopTokenRepo) CreateRefreshToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (n *noopTokenRepo) FindRefreshToken(_ context.Context, _ string) (*domain.RefreshToken, error) {
	return nil, nil
}
func (n *noopTokenRepo) DeleteRefreshToken(_ context.Context, _ string) error          { return nil }
func (n *noopTokenRepo) DeleteUserRefreshTokens(_ context.Context, _ string) error     { return nil }
func (n *noopTokenRepo) CreatePasswordResetToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (n *noopTokenRepo) FindPasswordResetToken(_ context.Context, _ string) (*domain.PasswordResetToken, error) {
	return nil, nil
}
func (n *noopTokenRepo) DeletePasswordResetToken(_ context.Context, _ string) error        { return nil }
func (n *noopTokenRepo) DeleteUserPasswordResetTokens(_ context.Context, _ string) error { return nil }

func newHandler(repo service.UserRepo) *handler.AuthHandler {
	svc := service.NewAuthService(repo, &noopTokenRepo{}, "test-secret")
	return handler.NewAuthHandler(svc)
}

func postJSON(t *testing.T, h http.HandlerFunc, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)
	return w
}

// ── POST /auth/register ───────────────────────────────────────────────────────

func TestRegisterHandler_201(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
		createFn: func(_ context.Context, u *domain.User) error {
			u.ID = "new-uuid"
			u.Role = "user"
			return nil
		},
	}

	w := postJSON(t, newHandler(repo).Register, "/auth/register", map[string]string{
		"email": "alice@example.com", "password": "password123", "first_name": "Alice",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.AuthResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "alice@example.com", resp.User.Email)
}

func TestRegisterHandler_409_EmailTaken(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{Email: "alice@example.com"}, nil
		},
	}

	w := postJSON(t, newHandler(repo).Register, "/auth/register", map[string]string{
		"email": "alice@example.com", "password": "password123", "first_name": "Alice",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterHandler_422_MissingFields(t *testing.T) {
	cases := []map[string]string{
		{"email": "alice@example.com", "password": "password123"}, // pas de first_name
		{"email": "alice@example.com", "first_name": "Alice"},     // pas de password
		{"password": "password123", "first_name": "Alice"},        // pas d'email
	}

	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	h := newHandler(repo)

	for _, body := range cases {
		w := postJSON(t, h.Register, "/auth/register", body)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "body: %v", body)
	}
}

// ── POST /auth/login ──────────────────────────────────────────────────────────

func TestLoginHandler_200(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: "user-uuid", Role: "user", PasswordHash: string(hash)}, nil
		},
	}

	w := postJSON(t, newHandler(repo).Login, "/auth/login", map[string]string{
		"email": "alice@example.com", "password": "password123",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.AuthResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Token)
}

func TestLoginHandler_401_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: "user-uuid", PasswordHash: string(hash)}, nil
		},
	}

	w := postJSON(t, newHandler(repo).Login, "/auth/login", map[string]string{
		"email": "alice@example.com", "password": "wrong",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandler_403_AccountLocked(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: "user-uuid", PasswordHash: string(hash), LockedUntil: &future}, nil
		},
	}

	w := postJSON(t, newHandler(repo).Login, "/auth/login", map[string]string{
		"email": "alice@example.com", "password": "password123",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLoginHandler_400_InvalidBody(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrNotFound
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newHandler(repo).Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
