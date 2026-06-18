package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"cesizen/internal/domain"
	"cesizen/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type trackerServiceIface interface {
	ListEntries(ctx context.Context, userID string, from, to time.Time) ([]domain.Entry, error)
	CreateEntry(ctx context.Context, userID string, input domain.CreateEntryInput) (*domain.Entry, error)
	UpdateEntry(ctx context.Context, userID string, id int, input domain.UpdateEntryInput) (*domain.Entry, error)
	DeleteEntry(ctx context.Context, userID string, id int) error
	Stats(ctx context.Context, userID string, period string) ([]domain.EmotionStat, error)
}

type TrackerHandler struct {
	svc trackerServiceIface
}

func NewTrackerHandler(svc trackerServiceIface) *TrackerHandler {
	return &TrackerHandler{svc: svc}
}

func (h *TrackerHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	from, to := parseRange(r)
	entries, err := h.svc.ListEntries(r.Context(), userID, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if entries == nil {
		entries = []domain.Entry{}
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *TrackerHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	var input domain.CreateEntryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.EmotionID == 0 {
		writeError(w, http.StatusUnprocessableEntity, "emotion_id is required")
		return
	}
	if input.Intensity < 1 || input.Intensity > 10 {
		writeError(w, http.StatusUnprocessableEntity, "intensity must be between 1 and 10")
		return
	}

	entry, err := h.svc.CreateEntry(r.Context(), userID, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (h *TrackerHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input domain.UpdateEntryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.EmotionID == 0 {
		writeError(w, http.StatusUnprocessableEntity, "emotion_id is required")
		return
	}
	if input.Intensity < 1 || input.Intensity > 10 {
		writeError(w, http.StatusUnprocessableEntity, "intensity must be between 1 and 10")
		return
	}

	entry, err := h.svc.UpdateEntry(r.Context(), userID, id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (h *TrackerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.DeleteEntry(r.Context(), userID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TrackerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)
	period := r.URL.Query().Get("period")

	stats, err := h.svc.Stats(r.Context(), userID, period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if stats == nil {
		stats = []domain.EmotionStat{}
	}
	writeJSON(w, http.StatusOK, stats)
}

func parseRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now()
	from := now.AddDate(0, -1, 0)
	// Borne haute = fin de journée : les entrées du jour sont horodatées à midi
	// (T12:00:00Z) côté front, donc un simple time.Now() les exclurait avant midi.
	to := endOfDay(now)

	if f := r.URL.Query().Get("from"); f != "" {
		if t, err := time.Parse("2006-01-02", f); err == nil {
			from = t
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			to = endOfDay(parsed)
		}
	}
	return from, to
}

// endOfDay renvoie le dernier instant du jour de t (23:59:59.999…), afin que la
// borne haute d'un intervalle inclue toutes les entrées datées de ce jour.
func endOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}
