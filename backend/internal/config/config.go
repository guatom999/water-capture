package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   Server
		Database Database
		App      App
		JWT      JWT
	}

	Server struct {
		Port int
	}

	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
		SSLMode  string
	}

	App struct {
		BaseURL            string
		UploadDir          string
		ImageProcessingDir string
	}

	JWT struct {
		Secret             string
		AccessTokenExpiry  int // in minutes
		RefreshTokenExpiry int // in days
	}
)

func LoadConfig(path string) *Config {
	// Try to load .env file, but don't fail if not found
	// (Docker Compose provides environment variables directly)
	if err := godotenv.Load(path); err != nil {
		log.Printf("Note: .env file not found at %s, using environment variables", path)
	}

	return &Config{
		Server: Server{
			Port: func() int {
				port, err := strconv.Atoi(os.Getenv("APP_PORT"))
				if err != nil {
					return 8080
				}
				return port
			}(),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		App: App{
			BaseURL: func() string {
				url := os.Getenv("BASE_URL")
				if url == "" {
					return "http://localhost:8080"
				}
				return url
			}(),
			UploadDir:          os.Getenv("UPLOAD_DIR"),
			ImageProcessingDir: os.Getenv("IMAGE_PROCESSING_DIR"),
		},
		JWT: JWT{
			Secret: func() string {
				secret := os.Getenv("JWT_SECRET")
				if secret == "" {
					return "your-super-secret-key-change-in-production"
				}
				return secret
			}(),
			AccessTokenExpiry: func() int {
				expiry, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY"))
				if err != nil {
					return 15 // default 15 minutes
				}
				return expiry
			}(),
			RefreshTokenExpiry: func() int {
				expiry, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY"))
				if err != nil {
					return 1 // default 1 day
				}
				return expiry
			}(),
		},
	}
}
