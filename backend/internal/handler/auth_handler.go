package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cesizen/internal/domain"
)

type authServiceIface interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.AuthResponse, error)
	Login(ctx context.Context, input domain.LoginInput) (*domain.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, input domain.ForgotPasswordInput) (*domain.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error
}

type AuthHandler struct {
	auth authServiceIface
}

func NewAuthHandler(auth authServiceIface) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Email == "" || input.Password == "" || input.FirstName == "" {
		writeError(w, http.StatusUnprocessableEntity, "email, password and first_name are required")
		return
	}

	resp, err := h.auth.Register(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailTaken):
			writeError(w, http.StatusConflict, "email already taken")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, domain.ErrAccountLocked):
			writeError(w, http.StatusForbidden, "account temporarily locked, try again in 15 minutes")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input domain.RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.RefreshToken == "" {
		writeError(w, http.StatusUnprocessableEntity, "refresh_token is required")
		return
	}

	resp, err := h.auth.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			writeError(w, http.StatusUnauthorized, "refresh token expired or invalid")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input domain.RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.RefreshToken != "" {
		_ = h.auth.Logout(r.Context(), input.RefreshToken)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input domain.ForgotPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Email == "" {
		writeError(w, http.StatusUnprocessableEntity, "email is required")
		return
	}

	resp, err := h.auth.ForgotPassword(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input domain.ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Token == "" || input.NewPassword == "" {
		writeError(w, http.StatusUnprocessableEntity, "token and new_password are required")
		return
	}
	if len(input.NewPassword) < 6 {
		writeError(w, http.StatusUnprocessableEntity, "new_password must be at least 6 characters")
		return
	}

	if err := h.auth.ResetPassword(r.Context(), input); err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			writeError(w, http.StatusUnauthorized, "reset token expired or invalid")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
