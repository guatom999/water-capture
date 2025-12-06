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
)

func LoadConfig(path string) *Config {

	if err := godotenv.Load(path); err != nil {
		log.Fatalf("Error loading env file %s", err.Error())
		return nil
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
	}
}
