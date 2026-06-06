package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env                string
	HTTPAddr           string
	JWTSecret          string
	DBPath             string
	CORSAllowedOrigins []string
	SMTP               SMTPConfig
	SMS                SMSConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type SMSConfig struct {
	Provider string
	APIKey   string
	From     string
}

func Load() Config {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

	return Config{
		Env:                getEnv("APP_ENV", "development"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret-change-me"),
		DBPath:             getEnv("DB_PATH", "./data/requests.db"),
		CORSAllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "")),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@requests-api.local"),
		},
		SMS: SMSConfig{
			Provider: getEnv("SMS_PROVIDER", "log"),
			APIKey:   getEnv("SMS_API_KEY", ""),
			From:     getEnv("SMS_FROM", "RequestsAPI"),
		},
	}
}

func parseOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:3000",
		}
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
