package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName            string
	AppEnv             string
	AppPort            string
	APIPrefix          string
	AppPublicURL       string
	FrontendURL        string
	CORSAllowedOrigins []string
	DatabaseURL        string
	RedisURL           string
	JWTAccessSecret    string
	JWTRefreshSecret   string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	PasswordHashAlgo   string
	BcryptCost         int
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	CookieSecure       bool
	CookieSameSite     string
	LogLevel           string
}

func Load() (*Config, error) {
	var missing []string
	required := map[string]*string{
		"APP_NAME":           nil,
		"APP_ENV":            nil,
		"APP_PORT":           nil,
		"API_PREFIX":         nil,
		"APP_PUBLIC_URL":     nil,
		"FRONTEND_URL":       nil,
		"DATABASE_URL":       nil,
		"REDIS_URL":          nil,
		"JWT_ACCESS_SECRET":  nil,
		"JWT_REFRESH_SECRET": nil,
	}

	values := make(map[string]string, len(required))
	for key := range required {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			missing = append(missing, key)
			continue
		}
		values[key] = value
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	accessTTL, err := parseDurationEnv("JWT_ACCESS_TTL", time.Hour)
	if err != nil {
		return nil, err
	}

	refreshTTL, err := parseDurationEnv("JWT_REFRESH_TTL", 720*time.Hour)
	if err != nil {
		return nil, err
	}

	bcryptCost, err := parseIntEnv("BCRYPT_COST", 12)
	if err != nil {
		return nil, err
	}

	cookieSecure, err := parseBoolEnv("COOKIE_SECURE", false)
	if err != nil {
		return nil, err
	}

	origins := splitCSV(getEnv("CORS_ALLOWED_ORIGINS", values["FRONTEND_URL"]))
	if len(origins) == 0 {
		return nil, errors.New("CORS_ALLOWED_ORIGINS must contain at least one origin")
	}

	return &Config{
		AppName:            values["APP_NAME"],
		AppEnv:             values["APP_ENV"],
		AppPort:            values["APP_PORT"],
		APIPrefix:          values["API_PREFIX"],
		AppPublicURL:       values["APP_PUBLIC_URL"],
		FrontendURL:        values["FRONTEND_URL"],
		CORSAllowedOrigins: origins,
		DatabaseURL:        values["DATABASE_URL"],
		RedisURL:           values["REDIS_URL"],
		JWTAccessSecret:    values["JWT_ACCESS_SECRET"],
		JWTRefreshSecret:   values["JWT_REFRESH_SECRET"],
		JWTAccessTTL:       accessTTL,
		JWTRefreshTTL:      refreshTTL,
		PasswordHashAlgo:   getEnv("PASSWORD_HASH_ALGO", "bcrypt"),
		BcryptCost:         bcryptCost,
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		CookieSecure:       cookieSecure,
		CookieSameSite:     getEnv("COOKIE_SAME_SITE", "Lax"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}

func parseDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s duration: %w", key, err)
	}
	return duration, nil
}

func parseIntEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s integer: %w", key, err)
	}
	return parsed, nil
}

func parseBoolEnv(key string, fallback bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid %s boolean: %w", key, err)
	}
	return parsed, nil
}
