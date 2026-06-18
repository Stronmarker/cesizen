package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cesizen/internal/domain"
	"cesizen/internal/handler"
	"cesizen/internal/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	profile *domain.PublicUser
	err     error
}

func (m *mockUserService) GetProfile(_ context.Context, _ string) (*domain.PublicUser, error) {
	return m.profile, m.err
}

func (m *mockUserService) UpdateProfile(_ context.Context, _ string, input domain.UpdateProfileInput) (*domain.PublicUser, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.PublicUser{
		ID:        m.profile.ID,
		Email:     m.profile.Email,
		FirstName: input.FirstName,
		Nickname:  input.Nickname,
		Role:      m.profile.Role,
	}, nil
}

func (m *mockUserService) DeleteAccount(_ context.Context, _ string) error {
	return m.err
}

// injecte l'userID dans le contexte comme le ferait le middleware Auth
func withUserID(r *http.Request, id string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.ContextUserID, id)
	return r.WithContext(ctx)
}

// ── GET /users/me ─────────────────────────────────────────────────────────────

func TestGetMe_200(t *testing.T) {
	svc := &mockUserService{profile: &domain.PublicUser{
		ID: "uuid-1", Email: "alice@example.com", FirstName: "Alice", Role: "user",
	}}
	h := handler.NewUserHandler(svc)

	req := withUserID(httptest.NewRequest(http.MethodGet, "/users/me", nil), "uuid-1")
	w := httptest.NewRecorder()
	h.GetMe(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.PublicUser
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "alice@example.com", resp.Email)
}

func TestGetMe_404(t *testing.T) {
	svc := &mockUserService{err: domain.ErrNotFound}
	h := handler.NewUserHandler(svc)

	req := withUserID(httptest.NewRequest(http.MethodGet, "/users/me", nil), "bad-id")
	w := httptest.NewRecorder()
	h.GetMe(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── PUT /users/me ─────────────────────────────────────────────────────────────

func TestUpdateMe_200(t *testing.T) {
	svc := &mockUserService{profile: &domain.PublicUser{
		ID: "uuid-1", Email: "alice@example.com", Role: "user",
	}}
	h := handler.NewUserHandler(svc)

	body, _ := json.Marshal(map[string]string{"first_name": "Alicia", "nickname": "ali"})
	req := withUserID(httptest.NewRequest(http.MethodPut, "/users/me", bytes.NewReader(body)), "uuid-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateMe(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.PublicUser
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Alicia", resp.FirstName)
	assert.Equal(t, "ali", resp.Nickname)
}

func TestUpdateMe_422_MissingFirstName(t *testing.T) {
	svc := &mockUserService{profile: &domain.PublicUser{ID: "uuid-1"}}
	h := handler.NewUserHandler(svc)

	body, _ := json.Marshal(map[string]string{"nickname": "ali"})
	req := withUserID(httptest.NewRequest(http.MethodPut, "/users/me", bytes.NewReader(body)), "uuid-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateMe(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUpdateMe_400_InvalidBody(t *testing.T) {
	svc := &mockUserService{profile: &domain.PublicUser{ID: "uuid-1"}}
	h := handler.NewUserHandler(svc)

	req := withUserID(httptest.NewRequest(http.MethodPut, "/users/me", bytes.NewBufferString("not-json")), "uuid-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateMe(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
