package service_test

import (
	"context"
	"testing"

	"cesizen/internal/domain"
	"cesizen/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProfileRepo struct {
	user          *domain.User
	err           error
	updateCalled  bool
	updatedUser   *domain.User
}

func (m *mockProfileRepo) FindByID(_ context.Context, _ string) (*domain.User, error) {
	return m.user, m.err
}

func (m *mockProfileRepo) UpdateProfile(_ context.Context, u *domain.User) error {
	m.updateCalled = true
	m.updatedUser = u
	return m.err
}

func (m *mockProfileRepo) Delete(_ context.Context, _ string) error {
	return m.err
}

func newUserService(repo service.UserProfileRepo) *service.UserService {
	return service.NewUserService(repo)
}

// ── GetProfile ────────────────────────────────────────────────────────────────

func TestGetProfile_Success(t *testing.T) {
	repo := &mockProfileRepo{user: &domain.User{
		ID: "uuid-1", Email: "alice@example.com", FirstName: "Alice", Role: "user",
	}}

	profile, err := newUserService(repo).GetProfile(context.Background(), "uuid-1")

	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", profile.Email)
	assert.Equal(t, "Alice", profile.FirstName)
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := &mockProfileRepo{err: domain.ErrNotFound}

	_, err := newUserService(repo).GetProfile(context.Background(), "unknown")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ── UpdateProfile ─────────────────────────────────────────────────────────────

func TestUpdateProfile_Success(t *testing.T) {
	repo := &mockProfileRepo{user: &domain.User{
		ID: "uuid-1", Email: "alice@example.com", FirstName: "Alice", Nickname: "",
	}}

	profile, err := newUserService(repo).UpdateProfile(context.Background(), "uuid-1", domain.UpdateProfileInput{
		FirstName: "Alicia",
		Nickname:  "ali",
	})

	require.NoError(t, err)
	assert.Equal(t, "Alicia", profile.FirstName)
	assert.Equal(t, "ali", profile.Nickname)
	assert.True(t, repo.updateCalled)
	assert.Equal(t, "Alicia", repo.updatedUser.FirstName)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	repo := &mockProfileRepo{err: domain.ErrNotFound}

	_, err := newUserService(repo).UpdateProfile(context.Background(), "unknown", domain.UpdateProfileInput{
		FirstName: "Alice",
	})

	assert.ErrorIs(t, err, domain.ErrNotFound)
}
