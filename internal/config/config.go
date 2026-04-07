package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DBDSN           string
	RateLimitRPS    float64
	RateLimitBurst  int
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
	AppEnv          string
	JWTSecret       string
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPass        string
	MailFrom        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	rps, err := strconv.ParseFloat(getEnv("RATE_LIMIT_RPS", "5"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}

	burst, err := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_BURST: %w", err)
	}

	shutdownTimeout, err := time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	allowedOrigins := strings.Split(
		getEnv("ALLOWED_ORIGINS", "http://127.0.0.1:5500,http://localhost:5500"),
		",",
	)

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		DBDSN:           getEnv("DB_DSN", ""),
		RateLimitRPS:    rps,
		RateLimitBurst:  burst,
		ShutdownTimeout: shutdownTimeout,
		AllowedOrigins:  allowedOrigins,
		AppEnv:          getEnv("APP_ENV", "development"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		SMTPHost:        getEnv("SMTP_HOST", ""),
		SMTPPort:        getEnv("SMTP_PORT", ""),
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPass:        getEnv("SMTP_PASS", ""),
		MailFrom:        getEnv("MAIL_FROM", ""),
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
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
