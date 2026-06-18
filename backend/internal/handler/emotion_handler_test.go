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
)

type mockEmotionService struct {
	primaryList []domain.PrimaryEmotion
	emotionList []domain.Emotion
	findErr     error
}

func (m *mockEmotionService) ListPrimary(_ context.Context) ([]domain.PrimaryEmotion, error) {
	return m.primaryList, nil
}

func (m *mockEmotionService) ListEmotions(_ context.Context) ([]domain.Emotion, error) {
	return m.emotionList, nil
}

func (m *mockEmotionService) CreatePrimary(_ context.Context, input domain.CreatePrimaryEmotionInput) (*domain.PrimaryEmotion, error) {
	return &domain.PrimaryEmotion{ID: 1, Label: input.Label, IsActive: true}, nil
}

func (m *mockEmotionService) UpdatePrimary(_ context.Context, id int, input domain.UpdatePrimaryEmotionInput) (*domain.PrimaryEmotion, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return &domain.PrimaryEmotion{ID: id, Label: input.Label, IsActive: input.IsActive}, nil
}

func (m *mockEmotionService) DeletePrimary(_ context.Context, _ int) error { return nil }

func (m *mockEmotionService) CreateEmotion(_ context.Context, input domain.CreateEmotionInput) (*domain.Emotion, error) {
	return &domain.Emotion{ID: 1, Label: input.Label, PrimaryEmotionID: input.PrimaryEmotionID, IsActive: true}, nil
}

func (m *mockEmotionService) UpdateEmotion(_ context.Context, id int, input domain.UpdateEmotionInput) (*domain.Emotion, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return &domain.Emotion{ID: id, Label: input.Label, PrimaryEmotionID: input.PrimaryEmotionID, IsActive: input.IsActive}, nil
}

func (m *mockEmotionService) DeleteEmotion(_ context.Context, _ int) error { return nil }

func TestEmotionHandler_ListPrimary_200(t *testing.T) {
	svc := &mockEmotionService{
		primaryList: []domain.PrimaryEmotion{{ID: 1, Label: "Joie", IsActive: true}},
	}
	h := handler.NewEmotionHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/primary-emotions", nil)
	w := httptest.NewRecorder()
	h.ListPrimary(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []domain.PrimaryEmotion
	json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 1 || list[0].Label != "Joie" {
		t.Fatalf("unexpected body: %+v", list)
	}
}

func TestEmotionHandler_List_200(t *testing.T) {
	svc := &mockEmotionService{
		emotionList: []domain.Emotion{{ID: 1, Label: "Fierté", PrimaryEmotionID: 1, PrimaryLabel: "Joie", IsActive: true}},
	}
	h := handler.NewEmotionHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/emotions", nil)
	w := httptest.NewRecorder()
	h.List(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminCreatePrimary_201(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	body, _ := json.Marshal(map[string]string{"label": "Surprise"})
	req := httptest.NewRequest(http.MethodPost, "/admin/primary-emotions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.AdminCreatePrimary(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminCreatePrimary_422(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	body, _ := json.Marshal(map[string]string{"label": ""})
	req := httptest.NewRequest(http.MethodPost, "/admin/primary-emotions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.AdminCreatePrimary(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminUpdatePrimary_200(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	body, _ := json.Marshal(map[string]any{"label": "Joie modifiée", "is_active": true})
	req := withChiID(httptest.NewRequest(http.MethodPut, "/admin/primary-emotions/1", bytes.NewReader(body)), "1")
	w := httptest.NewRecorder()
	h.AdminUpdatePrimary(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminUpdatePrimary_404(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{findErr: domain.ErrNotFound})
	body, _ := json.Marshal(map[string]any{"label": "X", "is_active": true})
	req := withChiID(httptest.NewRequest(http.MethodPut, "/admin/primary-emotions/99", bytes.NewReader(body)), "99")
	w := httptest.NewRecorder()
	h.AdminUpdatePrimary(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminDeletePrimary_204(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	req := withChiID(httptest.NewRequest(http.MethodDelete, "/admin/primary-emotions/1", nil), "1")
	w := httptest.NewRecorder()
	h.AdminDeletePrimary(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminCreateEmotion_201(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	body, _ := json.Marshal(map[string]any{"label": "Fierté", "primary_emotion_id": 1})
	req := httptest.NewRequest(http.MethodPost, "/admin/emotions", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.AdminCreateEmotion(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestEmotionHandler_AdminDeleteEmotion_204(t *testing.T) {
	h := handler.NewEmotionHandler(&mockEmotionService{})
	req := withChiID(httptest.NewRequest(http.MethodDelete, "/admin/emotions/1", nil), "1")
	w := httptest.NewRecorder()
	h.AdminDeleteEmotion(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
