package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cesizen/internal/domain"

	"github.com/go-chi/chi/v5"
)

type contentServiceIface interface {
	ListPublished(ctx context.Context) ([]domain.Content, error)
	ListAll(ctx context.Context) ([]domain.Content, error)
	GetByID(ctx context.Context, id int) (*domain.Content, error)
	Create(ctx context.Context, input domain.CreateContentInput) (*domain.Content, error)
	Update(ctx context.Context, id int, input domain.UpdateContentInput) (*domain.Content, error)
	Delete(ctx context.Context, id int) error
}

type ContentHandler struct {
	svc contentServiceIface
}

func NewContentHandler(svc contentServiceIface) *ContentHandler {
	return &ContentHandler{svc: svc}
}

func (h *ContentHandler) List(w http.ResponseWriter, r *http.Request) {
	contents, err := h.svc.ListPublished(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if contents == nil {
		contents = []domain.Content{}
	}
	writeJSON(w, http.StatusOK, contents)
}

func (h *ContentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ContentHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	contents, err := h.svc.ListAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if contents == nil {
		contents = []domain.Content{}
	}
	writeJSON(w, http.StatusOK, contents)
}

func (h *ContentHandler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateContentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Title == "" || input.Content == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and content are required")
		return
	}
	c, err := h.svc.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *ContentHandler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input domain.UpdateContentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Title == "" || input.Content == "" {
		writeError(w, http.StatusUnprocessableEntity, "title and content are required")
		return
	}
	c, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *ContentHandler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
