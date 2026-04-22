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
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRefreshToken  string
	GoogleDriveFolderID string
	GoogleRedirectURI   string
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
		GoogleDriveFolderID: os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		GoogleRedirectURI:   os.Getenv("GOOGLE_REDIRECT_URI"),
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
