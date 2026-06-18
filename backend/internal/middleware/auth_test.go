package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cesizen/internal/middleware"
)

// Vérifie que la barrière de rôle protégeant les routes /admin/* :
//   - laisse passer un administrateur (200) ;
//   - refuse un utilisateur authentifié non-admin avec 403 Forbidden ;
//   - refuse une requête sans rôle avec 403.
func TestRequireRoleAdmin(t *testing.T) {
	tests := []struct {
		name       string
		role       any // nil = aucun rôle dans le contexte
		wantStatus int
		wantNext   bool
	}{
		{name: "admin autorisé", role: "admin", wantStatus: http.StatusOK, wantNext: true},
		{name: "utilisateur non-admin refusé", role: "user", wantStatus: http.StatusForbidden, wantNext: false},
		{name: "rôle absent refusé", role: nil, wantStatus: http.StatusForbidden, wantNext: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.RequireRole("admin")(next)

			// Simule un POST /admin/contents (la route est protégée par RequireRole).
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/contents", nil)
			if tc.role != nil {
				ctx := context.WithValue(req.Context(), middleware.ContextUserRole, tc.role)
				req = req.WithContext(ctx)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("code HTTP = %d, attendu %d", rec.Code, tc.wantStatus)
			}
			if nextCalled != tc.wantNext {
				t.Errorf("handler protégé appelé = %v, attendu %v", nextCalled, tc.wantNext)
			}
		})
	}
}
