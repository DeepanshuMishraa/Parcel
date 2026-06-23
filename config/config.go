package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT           string
	DATABASE_URL   string
	REDIS_URL      string
	JWT_SECRET     string
	RESEND_API_KEY string
	FROM_EMAIL     string
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return &Config{}, errors.New("Env Not Recognized")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	PORT := os.Getenv("PORT")
	REDIS_URL := os.Getenv("REDIS_URL")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	RESEND_API_KEY := os.Getenv("RESEND_API_KEY")
	FROM_EMAIL := os.Getenv("FROM_EMAIL")

	if DATABASE_URL == "" || PORT == "" || REDIS_URL == "" {
		return &Config{}, errors.New("env vars cannot be empty")
	}

	log.Println("All Env Vars Loaded")

	return &Config{
		DATABASE_URL:   DATABASE_URL,
		PORT:           PORT,
		REDIS_URL:      REDIS_URL,
		JWT_SECRET:     JWT_SECRET,
		RESEND_API_KEY: RESEND_API_KEY,
		FROM_EMAIL:     FROM_EMAIL,
	}, nil
}
