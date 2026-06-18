package service

import (
	"context"

	"cesizen/internal/domain"
)

type AdminUserRepo interface {
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	AdminUpdate(ctx context.Context, id, role string, isActive bool) error
}

type AdminUserService struct {
	users AdminUserRepo
}

func NewAdminUserService(users AdminUserRepo) *AdminUserService {
	return &AdminUserService{users: users}
}

func (s *AdminUserService) ListUsers(ctx context.Context) ([]domain.AdminUserView, error) {
	users, err := s.users.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	views := make([]domain.AdminUserView, len(users))
	for i, u := range users {
		views[i] = toAdminView(u)
	}
	return views, nil
}

func (s *AdminUserService) GetUser(ctx context.Context, id string) (*domain.AdminUserView, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	v := toAdminView(u)
	return &v, nil
}

func (s *AdminUserService) UpdateUser(ctx context.Context, id string, input domain.AdminUpdateUserInput) (*domain.AdminUserView, error) {
	if input.Role != "user" && input.Role != "admin" {
		return nil, domain.ErrInvalidInput
	}
	if err := s.users.AdminUpdate(ctx, id, input.Role, input.IsActive); err != nil {
		return nil, err
	}
	return s.GetUser(ctx, id)
}

func toAdminView(u *domain.User) domain.AdminUserView {
	return domain.AdminUserView{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		Nickname:  u.Nickname,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
