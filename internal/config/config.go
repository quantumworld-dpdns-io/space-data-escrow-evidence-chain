package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	APIKey         string
	ShutdownTimout time.Duration
	Version        string
	Commit         string
	BuildDate      string
}

func Load() Config {
	return Config{
		Port:           getEnv("APP_PORT", "8080"),
		ReadTimeout:    getDuration("APP_READ_TIMEOUT_SEC", 10),
		WriteTimeout:   getDuration("APP_WRITE_TIMEOUT_SEC", 10),
		APIKey:         getEnv("APP_API_KEY", "dev-api-key"),
		ShutdownTimout: getDuration("APP_SHUTDOWN_TIMEOUT_SEC", 10),
		Version:        getEnv("APP_VERSION", "0.1.0"),
		Commit:         getEnv("APP_COMMIT", "dev"),
		BuildDate:      getEnv("APP_BUILD_DATE", "unknown"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDuration(key string, fallback int) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return time.Duration(fallback) * time.Second
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(parsed) * time.Second
}
