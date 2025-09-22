// Package config carga variables de entorno (.env o sistema) y las expone tipadas.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	SQLiteDSN       string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AESKey          string
}

func Load() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	sqlite := getEnv("SQLITE_DSN", "data/app.db")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-me")
	accessTTL := getEnvDuration("ACCESS_TOKEN_TTL", "15m")
	refreshTTL := getEnvDuration("REFRESH_TOKEN_TTL", "168h") // 7d
	aesKey := getEnv("AES_KEY", "0123456789abcdef0123456789abcdef")

	return &Config{
		Port:            port,
		SQLiteDSN:       sqlite,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		AESKey:          aesKey,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvDuration(key, def string) time.Duration {
	s := getEnv(key, def)
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("invalid duration for %s, using default %s", key, def)
		d, _ = time.ParseDuration(def)
	}
	return d
}
