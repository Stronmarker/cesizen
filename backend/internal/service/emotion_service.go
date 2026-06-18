package service

import (
	"context"

	"cesizen/internal/domain"
)

type EmotionRepo interface {
	FindAllPrimary(ctx context.Context) ([]domain.PrimaryEmotion, error)
	FindPrimaryByID(ctx context.Context, id int) (*domain.PrimaryEmotion, error)
	CreatePrimary(ctx context.Context, p *domain.PrimaryEmotion) error
	UpdatePrimary(ctx context.Context, p *domain.PrimaryEmotion) error
	DeletePrimary(ctx context.Context, id int) error

	FindAll(ctx context.Context) ([]domain.Emotion, error)
	FindByID(ctx context.Context, id int) (*domain.Emotion, error)
	CreateEmotion(ctx context.Context, e *domain.Emotion) error
	UpdateEmotion(ctx context.Context, e *domain.Emotion) error
	DeleteEmotion(ctx context.Context, id int) error
}

type EmotionService struct {
	repo EmotionRepo
}

func NewEmotionService(repo EmotionRepo) *EmotionService {
	return &EmotionService{repo: repo}
}

func (s *EmotionService) ListPrimary(ctx context.Context) ([]domain.PrimaryEmotion, error) {
	return s.repo.FindAllPrimary(ctx)
}

func (s *EmotionService) ListEmotions(ctx context.Context) ([]domain.Emotion, error) {
	return s.repo.FindAll(ctx)
}

func (s *EmotionService) CreatePrimary(ctx context.Context, input domain.CreatePrimaryEmotionInput) (*domain.PrimaryEmotion, error) {
	p := &domain.PrimaryEmotion{Label: input.Label}
	if err := s.repo.CreatePrimary(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *EmotionService) UpdatePrimary(ctx context.Context, id int, input domain.UpdatePrimaryEmotionInput) (*domain.PrimaryEmotion, error) {
	p, err := s.repo.FindPrimaryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Label = input.Label
	p.IsActive = input.IsActive
	if err := s.repo.UpdatePrimary(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *EmotionService) DeletePrimary(ctx context.Context, id int) error {
	return s.repo.DeletePrimary(ctx, id)
}

func (s *EmotionService) CreateEmotion(ctx context.Context, input domain.CreateEmotionInput) (*domain.Emotion, error) {
	e := &domain.Emotion{Label: input.Label, PrimaryEmotionID: input.PrimaryEmotionID}
	if err := s.repo.CreateEmotion(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EmotionService) UpdateEmotion(ctx context.Context, id int, input domain.UpdateEmotionInput) (*domain.Emotion, error) {
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	e.Label = input.Label
	e.PrimaryEmotionID = input.PrimaryEmotionID
	e.IsActive = input.IsActive
	if err := s.repo.UpdateEmotion(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EmotionService) DeleteEmotion(ctx context.Context, id int) error {
	return s.repo.DeleteEmotion(ctx, id)
}
