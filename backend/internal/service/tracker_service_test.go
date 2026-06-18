package service_test

import (
	"context"
	"testing"
	"time"

	"cesizen/internal/domain"
	"cesizen/internal/service"
)

type mockEntryRepo struct {
	entries    []domain.Entry
	findErr    error
	createErr  error
}

func (m *mockEntryRepo) FindByUser(_ context.Context, _ string, _, _ time.Time) ([]domain.Entry, error) {
	return m.entries, nil
}

func (m *mockEntryRepo) FindByID(_ context.Context, id int) (*domain.Entry, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, e := range m.entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockEntryRepo) Create(_ context.Context, e *domain.Entry) error {
	if m.createErr != nil {
		return m.createErr
	}
	e.ID = 42
	e.CreatedAt = time.Now()
	return nil
}

func (m *mockEntryRepo) Update(_ context.Context, _ *domain.Entry) error  { return nil }
func (m *mockEntryRepo) Delete(_ context.Context, _ int) error             { return nil }
func (m *mockEntryRepo) StatsForUser(_ context.Context, _ string, _, _ time.Time) ([]domain.EmotionStat, error) {
	return []domain.EmotionStat{{EmotionID: 1, EmotionLabel: "Fierté", PrimaryLabel: "Joie", Count: 3}}, nil
}

func TestTrackerService_ListEntries(t *testing.T) {
	repo := &mockEntryRepo{
		entries: []domain.Entry{{ID: 1, UserID: "u1", EmotionID: 1, Intensity: 7}},
	}
	svc := service.NewTrackerService(repo)
	list, err := svc.ListEntries(context.Background(), "u1", time.Now().AddDate(0, -1, 0), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
}

func TestTrackerService_CreateEntry(t *testing.T) {
	svc := service.NewTrackerService(&mockEntryRepo{})
	entry, err := svc.CreateEntry(context.Background(), "u1", domain.CreateEntryInput{
		EmotionID: 1, Intensity: 8, Comment: "top", EntryDate: time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.ID != 42 || entry.EmotionID != 1 || entry.Intensity != 8 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}

func TestTrackerService_UpdateEntry_Success(t *testing.T) {
	now := time.Now()
	repo := &mockEntryRepo{
		entries: []domain.Entry{{ID: 1, UserID: "u1", EmotionID: 1, Intensity: 5, EntryDate: now}},
	}
	svc := service.NewTrackerService(repo)
	updated, err := svc.UpdateEntry(context.Background(), "u1", 1, domain.UpdateEntryInput{
		EmotionID: 2, Intensity: 9, Comment: "updated", EntryDate: now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.EmotionID != 2 || updated.Intensity != 9 {
		t.Fatalf("unexpected updated entry: %+v", updated)
	}
}

func TestTrackerService_UpdateEntry_NotOwned(t *testing.T) {
	now := time.Now()
	repo := &mockEntryRepo{
		entries: []domain.Entry{{ID: 1, UserID: "other", EmotionID: 1, Intensity: 5, EntryDate: now}},
	}
	svc := service.NewTrackerService(repo)
	_, err := svc.UpdateEntry(context.Background(), "u1", 1, domain.UpdateEntryInput{
		EmotionID: 1, Intensity: 5, EntryDate: now,
	})
	if err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound for wrong owner, got %v", err)
	}
}

func TestTrackerService_DeleteEntry_Success(t *testing.T) {
	repo := &mockEntryRepo{
		entries: []domain.Entry{{ID: 1, UserID: "u1", EmotionID: 1, Intensity: 5}},
	}
	svc := service.NewTrackerService(repo)
	if err := svc.DeleteEntry(context.Background(), "u1", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrackerService_Stats(t *testing.T) {
	svc := service.NewTrackerService(&mockEntryRepo{})
	stats, err := svc.Stats(context.Background(), "u1", "month")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 || stats[0].Count != 3 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
