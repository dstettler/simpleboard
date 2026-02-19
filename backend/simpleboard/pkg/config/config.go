package config

import (
	"os"
	"strconv"
	"strings"
)

// Config type
type Config struct {
	Port          int
	ServerAddress string
	DBPath        string
	CORSOrigins   []string
}

// Loads the config by getting env defined variables
func Load() *Config {
	port := getEnvInt("PORT", 8080)

	return &Config{
		Port:          port,
		ServerAddress: ":" + strconv.Itoa(port),
		DBPath:        getEnv("DB_PATH", "./simpleboard.db"),
		CORSOrigins:   getEnvList("CORS_ORIGINS", []string{"http://localhost:4200"}),
	}
}

// String env get
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Int env get
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// String list parser from string env get
func getEnvList(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return fallback
}
