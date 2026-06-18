package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cesizen/internal/domain"

	"github.com/go-chi/chi/v5"
)

type adminUserServiceIface interface {
	ListUsers(ctx context.Context) ([]domain.AdminUserView, error)
	GetUser(ctx context.Context, id string) (*domain.AdminUserView, error)
	UpdateUser(ctx context.Context, id string, input domain.AdminUpdateUserInput) (*domain.AdminUserView, error)
}

type AdminUserHandler struct {
	users adminUserServiceIface
}

func NewAdminUserHandler(users adminUserServiceIface) *AdminUserHandler {
	return &AdminUserHandler{users: users}
}

func (h *AdminUserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *AdminUserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u, err := h.users.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *AdminUserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input domain.AdminUpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Role == "" {
		writeError(w, http.StatusUnprocessableEntity, "role is required")
		return
	}
	u, err := h.users.UpdateUser(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusUnprocessableEntity, "role must be 'user' or 'admin'")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, u)
}
