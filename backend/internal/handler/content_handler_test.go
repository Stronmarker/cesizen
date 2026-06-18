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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockContentService struct {
	contents []domain.Content
	content  *domain.Content
	err      error
}

func (m *mockContentService) ListPublished(_ context.Context) ([]domain.Content, error) {
	return m.contents, m.err
}
func (m *mockContentService) ListAll(_ context.Context) ([]domain.Content, error) {
	return m.contents, m.err
}
func (m *mockContentService) GetByID(_ context.Context, _ int) (*domain.Content, error) {
	return m.content, m.err
}
func (m *mockContentService) Create(_ context.Context, input domain.CreateContentInput) (*domain.Content, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.Content{ID: 1, Title: input.Title, Content: input.Content, Author: input.Author, IsPublished: true}, nil
}
func (m *mockContentService) Update(_ context.Context, _ int, input domain.UpdateContentInput) (*domain.Content, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.Content{ID: 1, Title: input.Title, Content: input.Content, IsPublished: input.IsPublished}, nil
}
func (m *mockContentService) Delete(_ context.Context, _ int) error {
	return m.err
}

func withChiID(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// ── GET /contents ─────────────────────────────────────────────────────────────

func TestContentList_200(t *testing.T) {
	svc := &mockContentService{contents: []domain.Content{
		{ID: 1, Title: "Article", Content: "Corps", IsPublished: true},
	}}
	h := handler.NewContentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/contents", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.Content
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Len(t, resp, 1)
}

func TestContentList_EmptySlice(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	req := httptest.NewRequest(http.MethodGet, "/contents", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []domain.Content
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp)
}

// ── GET /contents/:id ─────────────────────────────────────────────────────────

func TestContentGet_200(t *testing.T) {
	svc := &mockContentService{content: &domain.Content{ID: 1, Title: "Article", Content: "Corps"}}
	h := handler.NewContentHandler(svc)

	req := withChiID(httptest.NewRequest(http.MethodGet, "/contents/1", nil), "1")
	w := httptest.NewRecorder()
	h.Get(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp domain.Content
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Article", resp.Title)
}

func TestContentGet_404(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{err: domain.ErrNotFound})
	req := withChiID(httptest.NewRequest(http.MethodGet, "/contents/99", nil), "99")
	w := httptest.NewRecorder()
	h.Get(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestContentGet_400_InvalidID(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	req := withChiID(httptest.NewRequest(http.MethodGet, "/contents/abc", nil), "abc")
	w := httptest.NewRecorder()
	h.Get(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ── POST /admin/contents ──────────────────────────────────────────────────────

func TestAdminCreate_201(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	body, _ := json.Marshal(domain.CreateContentInput{Title: "Titre", Content: "Corps", Author: "Admin"})
	req := httptest.NewRequest(http.MethodPost, "/admin/contents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AdminCreate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp domain.Content
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "Titre", resp.Title)
}

func TestAdminCreate_422_MissingFields(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	body, _ := json.Marshal(map[string]string{"title": "Titre"})
	req := httptest.NewRequest(http.MethodPost, "/admin/contents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AdminCreate(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// ── PUT /admin/contents/:id ───────────────────────────────────────────────────

func TestAdminUpdate_200(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	body, _ := json.Marshal(domain.UpdateContentInput{Title: "Modifié", Content: "Corps", IsPublished: true})
	req := withChiID(httptest.NewRequest(http.MethodPut, "/admin/contents/1", bytes.NewReader(body)), "1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AdminUpdate(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminUpdate_404(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{err: domain.ErrNotFound})
	body, _ := json.Marshal(domain.UpdateContentInput{Title: "X", Content: "Y"})
	req := withChiID(httptest.NewRequest(http.MethodPut, "/admin/contents/99", bytes.NewReader(body)), "99")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AdminUpdate(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── DELETE /admin/contents/:id ────────────────────────────────────────────────

func TestAdminDelete_204(t *testing.T) {
	h := handler.NewContentHandler(&mockContentService{})
	req := withChiID(httptest.NewRequest(http.MethodDelete, "/admin/contents/1", nil), "1")
	w := httptest.NewRecorder()
	h.AdminDelete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
