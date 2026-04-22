package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	JWTSecret           string
	APIKey              string
	Port                string
	OwnerEmail          string
	// Drive OAuth app (offline refresh_token flow) — used by drive.UploadFile.
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRefreshToken  string
	GoogleDriveFolderID string
	// Login OAuth app (interactive web flow) — separate app to avoid scope/consent collision.
	// Falls back to the Drive app credentials when unset, so projects with a single OAuth app still work.
	GoogleLoginClientID     string
	GoogleLoginClientSecret string
	GoogleRedirectURI       string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		APIKey:              os.Getenv("API_KEY"),
		Port:                os.Getenv("PORT"),
		OwnerEmail:          os.Getenv("OWNER_EMAIL"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRefreshToken:  os.Getenv("GOOGLE_REFRESH_TOKEN"),
		GoogleDriveFolderID:     os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		GoogleLoginClientID:     os.Getenv("GOOGLE_LOGIN_CLIENT_ID"),
		GoogleLoginClientSecret: os.Getenv("GOOGLE_LOGIN_CLIENT_SECRET"),
		GoogleRedirectURI:       os.Getenv("GOOGLE_REDIRECT_URI"),
	}

	// Default login credentials to Drive credentials when login-specific ones are not set.
	if cfg.GoogleLoginClientID == "" {
		cfg.GoogleLoginClientID = cfg.GoogleClientID
	}
	if cfg.GoogleLoginClientSecret == "" {
		cfg.GoogleLoginClientSecret = cfg.GoogleClientSecret
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}
	if cfg.OwnerEmail == "" {
		cfg.OwnerEmail = "vuongnguyenbinh@gmail.com"
	}
	if cfg.GoogleRedirectURI == "" {
		cfg.GoogleRedirectURI = "https://os.binhvuong.vn/auth/google/callback"
	}

	return cfg
}
