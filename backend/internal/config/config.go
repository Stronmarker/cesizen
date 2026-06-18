package config

import "os"

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	AllowedOrigins string
}

func Load() Config {
	return Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://cesizen:cesizen@localhost:5432/cesizen?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "changeme"),
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
