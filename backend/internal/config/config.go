package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	AllowedOrigins string
}

func Load() Config {
	return Config{
		// Aucun secret en dur : ces variables sont fournies par l'environnement
		// (voir .env / docker-compose). L'app refuse de démarrer si elles manquent.
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// mustEnv retourne la variable d'environnement ou arrête l'app si elle est absente.
// Évite tout identifiant/secret codé en dur (bonne pratique OWASP A05/A07).
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("variable d'environnement requise manquante : %s", key)
	}
	return v
}
