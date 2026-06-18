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
)

type mockTrackerService struct {
	entries []domain.Entry
	stats   []domain.EmotionStat
	findErr error
}

func (m *mockTrackerService) ListEntries(_ context.Context, _ string, _, _ time.Time) ([]domain.Entry, error) {
	return m.entries, nil
}

func (m *mockTrackerService) CreateEntry(_ context.Context, userID string, input domain.CreateEntryInput) (*domain.Entry, error) {
	return &domain.Entry{
		ID:        1,
		UserID:    userID,
		EmotionID: input.EmotionID,
		Intensity: input.Intensity,
		Comment:   input.Comment,
		EntryDate: input.EntryDate,
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockTrackerService) UpdateEntry(_ context.Context, _ string, id int, input domain.UpdateEntryInput) (*domain.Entry, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return &domain.Entry{ID: id, EmotionID: input.EmotionID, Intensity: input.Intensity}, nil
}

func (m *mockTrackerService) DeleteEntry(_ context.Context, _ string, _ int) error {
	return m.findErr
}

func (m *mockTrackerService) Stats(_ context.Context, _ string, _ string) ([]domain.EmotionStat, error) {
	return m.stats, nil
}

func TestTrackerHandler_List_200(t *testing.T) {
	svc := &mockTrackerService{
		entries: []domain.Entry{{ID: 1, UserID: "u1", EmotionID: 1, Intensity: 7}},
	}
	h := handler.NewTrackerHandler(svc)
	req := withUserID(httptest.NewRequest(http.MethodGet, "/tracker/entries", nil), "u1")
	w := httptest.NewRecorder()
	h.List(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []domain.Entry
	json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
}

func TestTrackerHandler_Create_201(t *testing.T) {
	h := handler.NewTrackerHandler(&mockTrackerService{})
	body, _ := json.Marshal(map[string]any{
		"emotion_id": 1,
		"intensity":  8,
		"comment":    "top",
		"entry_date": time.Now(),
	})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/tracker/entries", bytes.NewReader(body)), "u1")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestTrackerHandler_Create_422_MissingEmotion(t *testing.T) {
	h := handler.NewTrackerHandler(&mockTrackerService{})
	body, _ := json.Marshal(map[string]any{"intensity": 5})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/tracker/entries", bytes.NewReader(body)), "u1")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestTrackerHandler_Create_422_BadIntensity(t *testing.T) {
	h := handler.NewTrackerHandler(&mockTrackerService{})
	body, _ := json.Marshal(map[string]any{"emotion_id": 1, "intensity": 11})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/tracker/entries", bytes.NewReader(body)), "u1")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestTrackerHandler_Update_200(t *testing.T) {
	h := handler.NewTrackerHandler(&mockTrackerService{})
	body, _ := json.Marshal(map[string]any{
		"emotion_id": 2,
		"intensity":  6,
		"entry_date": time.Now(),
	})
	req := withUserID(withChiID(httptest.NewRequest(http.MethodPut, "/tracker/entries/1", bytes.NewReader(body)), "1"), "u1")
	w := httptest.NewRecorder()
	h.Update(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestTrackerHandler_Update_404(t *testing.T) {
	svc := &mockTrackerService{findErr: domain.ErrNotFound}
	h := handler.NewTrackerHandler(svc)
	body, _ := json.Marshal(map[string]any{"emotion_id": 1, "intensity": 5, "entry_date": time.Now()})
	req := withUserID(withChiID(httptest.NewRequest(http.MethodPut, "/tracker/entries/99", bytes.NewReader(body)), "99"), "u1")
	w := httptest.NewRecorder()
	h.Update(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestTrackerHandler_Delete_204(t *testing.T) {
	h := handler.NewTrackerHandler(&mockTrackerService{})
	req := withUserID(withChiID(httptest.NewRequest(http.MethodDelete, "/tracker/entries/1", nil), "1"), "u1")
	w := httptest.NewRecorder()
	h.Delete(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestTrackerHandler_Stats_200(t *testing.T) {
	svc := &mockTrackerService{
		stats: []domain.EmotionStat{{EmotionID: 1, EmotionLabel: "Fierté", PrimaryLabel: "Joie", Count: 3}},
	}
	h := handler.NewTrackerHandler(svc)
	req := withUserID(httptest.NewRequest(http.MethodGet, "/tracker/stats?period=month", nil), "u1")
	w := httptest.NewRecorder()
	h.Stats(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var stats []domain.EmotionStat
	json.NewDecoder(w.Body).Decode(&stats)
	if len(stats) != 1 || stats[0].Count != 3 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
