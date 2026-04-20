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
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}
	if cfg.OwnerEmail == "" {
		cfg.OwnerEmail = "vuongnguyenbinh@gmail.com"
	}

	return cfg
}
