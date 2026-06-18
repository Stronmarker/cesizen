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

type emotionServiceIface interface {
	ListPrimary(ctx context.Context) ([]domain.PrimaryEmotion, error)
	ListEmotions(ctx context.Context) ([]domain.Emotion, error)
	CreatePrimary(ctx context.Context, input domain.CreatePrimaryEmotionInput) (*domain.PrimaryEmotion, error)
	UpdatePrimary(ctx context.Context, id int, input domain.UpdatePrimaryEmotionInput) (*domain.PrimaryEmotion, error)
	DeletePrimary(ctx context.Context, id int) error
	CreateEmotion(ctx context.Context, input domain.CreateEmotionInput) (*domain.Emotion, error)
	UpdateEmotion(ctx context.Context, id int, input domain.UpdateEmotionInput) (*domain.Emotion, error)
	DeleteEmotion(ctx context.Context, id int) error
}

type EmotionHandler struct {
	svc emotionServiceIface
}

func NewEmotionHandler(svc emotionServiceIface) *EmotionHandler {
	return &EmotionHandler{svc: svc}
}

func (h *EmotionHandler) ListPrimary(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListPrimary(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if list == nil {
		list = []domain.PrimaryEmotion{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *EmotionHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListEmotions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if list == nil {
		list = []domain.Emotion{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *EmotionHandler) AdminCreatePrimary(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePrimaryEmotionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Label == "" {
		writeError(w, http.StatusUnprocessableEntity, "label is required")
		return
	}
	p, err := h.svc.CreatePrimary(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *EmotionHandler) AdminUpdatePrimary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input domain.UpdatePrimaryEmotionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Label == "" {
		writeError(w, http.StatusUnprocessableEntity, "label is required")
		return
	}
	p, err := h.svc.UpdatePrimary(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *EmotionHandler) AdminDeletePrimary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeletePrimary(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EmotionHandler) AdminCreateEmotion(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateEmotionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Label == "" || input.PrimaryEmotionID == 0 {
		writeError(w, http.StatusUnprocessableEntity, "label and primary_emotion_id are required")
		return
	}
	e, err := h.svc.CreateEmotion(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *EmotionHandler) AdminUpdateEmotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input domain.UpdateEmotionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Label == "" || input.PrimaryEmotionID == 0 {
		writeError(w, http.StatusUnprocessableEntity, "label and primary_emotion_id are required")
		return
	}
	e, err := h.svc.UpdateEmotion(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *EmotionHandler) AdminDeleteEmotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteEmotion(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
