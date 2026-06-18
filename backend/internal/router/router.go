package router

import (
	"net/http"

	"cesizen/internal/handler"
	"cesizen/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminUserHandler *handler.AdminUserHandler,
	contentHandler *handler.ContentHandler,
	emotionHandler *handler.EmotionHandler,
	trackerHandler *handler.TrackerHandler,
	jwtSecret string,
	allowedOrigins string,
) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(allowedOrigins))

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
			r.Post("/logout", authHandler.Logout)
		})
		r.Get("/contents", contentHandler.List)
		r.Get("/contents/{id}", contentHandler.Get)
		r.Get("/primary-emotions", emotionHandler.ListPrimary)
		r.Get("/emotions", emotionHandler.List)

		// Authenticated
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))
			r.Get("/users/me", userHandler.GetMe)
			r.Put("/users/me", userHandler.UpdateMe)
			r.Delete("/users/me", userHandler.DeleteMe)

			// Tracker
			r.Get("/tracker/entries", trackerHandler.List)
			r.Post("/tracker/entries", trackerHandler.Create)
			r.Put("/tracker/entries/{id}", trackerHandler.Update)
			r.Delete("/tracker/entries/{id}", trackerHandler.Delete)
			r.Get("/tracker/stats", trackerHandler.Stats)

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Get("/admin/contents", contentHandler.AdminList)
				r.Post("/admin/contents", contentHandler.AdminCreate)
				r.Put("/admin/contents/{id}", contentHandler.AdminUpdate)
				r.Delete("/admin/contents/{id}", contentHandler.AdminDelete)

				r.Get("/admin/users", adminUserHandler.List)
				r.Get("/admin/users/{id}", adminUserHandler.Get)
				r.Put("/admin/users/{id}", adminUserHandler.Update)

				r.Post("/admin/primary-emotions", emotionHandler.AdminCreatePrimary)
				r.Put("/admin/primary-emotions/{id}", emotionHandler.AdminUpdatePrimary)
				r.Delete("/admin/primary-emotions/{id}", emotionHandler.AdminDeletePrimary)

				r.Post("/admin/emotions", emotionHandler.AdminCreateEmotion)
				r.Put("/admin/emotions/{id}", emotionHandler.AdminUpdateEmotion)
				r.Delete("/admin/emotions/{id}", emotionHandler.AdminDeleteEmotion)
			})
		})
	})

	return r
}
