package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	AI     AIConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type AIConfig struct {
	GeminiAPIKey string
	Model        string
	Temperature  float32
	MaxTokens    int32
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "3000"),
			Env:  getEnv("APP_ENV", "development"),
		},
		AI: AIConfig{
			GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
			Model:        getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
			Temperature:  0.8,
			MaxTokens:    1024,
		},
	}

	if cfg.AI.GeminiAPIKey == "" {
		return nil, errors.New("GEMINI_API_KEY is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

// Optional helper if you later want configurable numeric settings.
func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
