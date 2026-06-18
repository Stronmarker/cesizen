package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cesizen/internal/domain"
	"cesizen/internal/middleware"
)

type userServiceIface interface {
	GetProfile(ctx context.Context, userID string) (*domain.PublicUser, error)
	UpdateProfile(ctx context.Context, userID string, input domain.UpdateProfileInput) (*domain.PublicUser, error)
	DeleteAccount(ctx context.Context, userID string) error
}

type UserHandler struct {
	users userServiceIface
}

func NewUserHandler(users userServiceIface) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	profile, err := h.users.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)
	if err := h.users.DeleteAccount(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	var input domain.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.FirstName == "" {
		writeError(w, http.StatusUnprocessableEntity, "first_name is required")
		return
	}

	profile, err := h.users.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}
