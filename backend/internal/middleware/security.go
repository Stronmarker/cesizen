package middleware

import "net/http"

// SecurityHeaders ajoute les en-têtes de sécurité HTTP sur toutes les réponses
// de l'API (durcissement OWASP A05 — Security Misconfiguration).
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		// Empêche le MIME-sniffing (ZAP 10021).
		h.Set("X-Content-Type-Options", "nosniff")
		// L'API n'a jamais vocation à être affichée dans une frame.
		h.Set("X-Frame-Options", "DENY")
		// Ne fuite pas l'URL référente vers des tiers.
		h.Set("Referrer-Policy", "no-referrer")
		// Réponses d'API : jamais mises en cache (données potentiellement sensibles, ZAP 10049).
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
