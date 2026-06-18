package service

import (
	"context"

	"cesizen/internal/domain"
)

type UserProfileRepo interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	UpdateProfile(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	users UserProfileRepo
}

func NewUserService(users UserProfileRepo) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.PublicUser, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	p := userToPublic(user)
	return &p, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, input domain.UpdateProfileInput) (*domain.PublicUser, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.FirstName = input.FirstName
	user.Nickname = input.Nickname
	if err := s.users.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	p := userToPublic(user)
	return &p, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	return s.users.Delete(ctx, userID)
}

func userToPublic(u *domain.User) domain.PublicUser {
	return domain.PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		Nickname:  u.Nickname,
		Role:      u.Role,
	}
}
