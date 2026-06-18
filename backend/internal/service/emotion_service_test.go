package service_test

import (
	"context"
	"testing"

	"cesizen/internal/domain"
	"cesizen/internal/service"
)

type mockEmotionRepo struct {
	primaryEmotions []domain.PrimaryEmotion
	emotions        []domain.Emotion
	findPrimaryErr  error
	findEmotionErr  error
}

func (m *mockEmotionRepo) FindAllPrimary(_ context.Context) ([]domain.PrimaryEmotion, error) {
	return m.primaryEmotions, nil
}

func (m *mockEmotionRepo) FindPrimaryByID(_ context.Context, id int) (*domain.PrimaryEmotion, error) {
	if m.findPrimaryErr != nil {
		return nil, m.findPrimaryErr
	}
	for _, p := range m.primaryEmotions {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockEmotionRepo) CreatePrimary(_ context.Context, p *domain.PrimaryEmotion) error {
	p.ID = 99
	p.IsActive = true
	return nil
}

func (m *mockEmotionRepo) UpdatePrimary(_ context.Context, _ *domain.PrimaryEmotion) error { return nil }
func (m *mockEmotionRepo) DeletePrimary(_ context.Context, _ int) error                     { return nil }

func (m *mockEmotionRepo) FindAll(_ context.Context) ([]domain.Emotion, error) {
	return m.emotions, nil
}

func (m *mockEmotionRepo) FindByID(_ context.Context, id int) (*domain.Emotion, error) {
	if m.findEmotionErr != nil {
		return nil, m.findEmotionErr
	}
	for _, e := range m.emotions {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockEmotionRepo) CreateEmotion(_ context.Context, e *domain.Emotion) error {
	e.ID = 99
	e.IsActive = true
	return nil
}

func (m *mockEmotionRepo) UpdateEmotion(_ context.Context, _ *domain.Emotion) error { return nil }
func (m *mockEmotionRepo) DeleteEmotion(_ context.Context, _ int) error              { return nil }

func TestEmotionService_ListPrimary(t *testing.T) {
	repo := &mockEmotionRepo{
		primaryEmotions: []domain.PrimaryEmotion{{ID: 1, Label: "Joie", IsActive: true}},
	}
	svc := service.NewEmotionService(repo)
	list, err := svc.ListPrimary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].Label != "Joie" {
		t.Fatalf("expected 1 primary emotion Joie, got %+v", list)
	}
}

func TestEmotionService_ListEmotions(t *testing.T) {
	repo := &mockEmotionRepo{
		emotions: []domain.Emotion{{ID: 1, Label: "Fierté", PrimaryEmotionID: 1, PrimaryLabel: "Joie", IsActive: true}},
	}
	svc := service.NewEmotionService(repo)
	list, err := svc.ListEmotions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].Label != "Fierté" {
		t.Fatalf("expected 1 emotion Fierté, got %+v", list)
	}
}

func TestEmotionService_CreatePrimary(t *testing.T) {
	svc := service.NewEmotionService(&mockEmotionRepo{})
	p, err := svc.CreatePrimary(context.Background(), domain.CreatePrimaryEmotionInput{Label: "Joie"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 99 || p.Label != "Joie" {
		t.Fatalf("unexpected result: %+v", p)
	}
}

func TestEmotionService_UpdatePrimary_Success(t *testing.T) {
	repo := &mockEmotionRepo{
		primaryEmotions: []domain.PrimaryEmotion{{ID: 1, Label: "Joie", IsActive: true}},
	}
	svc := service.NewEmotionService(repo)
	p, err := svc.UpdatePrimary(context.Background(), 1, domain.UpdatePrimaryEmotionInput{Label: "Joie modifiée", IsActive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Label != "Joie modifiée" {
		t.Fatalf("expected label Joie modifiée, got %s", p.Label)
	}
}

func TestEmotionService_UpdatePrimary_NotFound(t *testing.T) {
	svc := service.NewEmotionService(&mockEmotionRepo{})
	_, err := svc.UpdatePrimary(context.Background(), 99, domain.UpdatePrimaryEmotionInput{Label: "X"})
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEmotionService_CreateEmotion(t *testing.T) {
	svc := service.NewEmotionService(&mockEmotionRepo{})
	e, err := svc.CreateEmotion(context.Background(), domain.CreateEmotionInput{Label: "Fierté", PrimaryEmotionID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.ID != 99 || e.Label != "Fierté" {
		t.Fatalf("unexpected result: %+v", e)
	}
}

func TestEmotionService_UpdateEmotion_NotFound(t *testing.T) {
	svc := service.NewEmotionService(&mockEmotionRepo{})
	_, err := svc.UpdateEmotion(context.Background(), 99, domain.UpdateEmotionInput{Label: "X", PrimaryEmotionID: 1})
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestEmotionService_DeleteEmotion(t *testing.T) {
	svc := service.NewEmotionService(&mockEmotionRepo{})
	if err := svc.DeleteEmotion(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
