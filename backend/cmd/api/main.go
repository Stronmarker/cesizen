package main

import (
	"fmt"
	"log"
	"net/http"

	"cesizen/internal/config"
	"cesizen/internal/db"
	"cesizen/internal/handler"
	"cesizen/internal/repository"
	"cesizen/internal/router"
	"cesizen/internal/service"
)

func main() {
	cfg := config.Load()

	pool, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	contentRepo := repository.NewContentRepository(pool)
	emotionRepo := repository.NewEmotionRepository(pool)
	entryRepo := repository.NewEntryRepository(pool)

	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo)
	adminUserSvc := service.NewAdminUserService(userRepo)
	contentSvc := service.NewContentService(contentRepo)
	emotionSvc := service.NewEmotionService(emotionRepo)
	trackerSvc := service.NewTrackerService(entryRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	adminUserHandler := handler.NewAdminUserHandler(adminUserSvc)
	contentHandler := handler.NewContentHandler(contentSvc)
	emotionHandler := handler.NewEmotionHandler(emotionSvc)
	trackerHandler := handler.NewTrackerHandler(trackerSvc)

	r := router.New(authHandler, userHandler, adminUserHandler, contentHandler, emotionHandler, trackerHandler, cfg.JWTSecret, cfg.AllowedOrigins)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
