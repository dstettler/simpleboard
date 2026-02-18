package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       int
	DBPath     string
	CORSOrigins []string
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	dbPath := "./simpleboard.db"
	if d := os.Getenv("DB_PATH"); d != "" {
		dbPath = d
	}

	origins := []string{"http://localhost:4200"}
	if o := os.Getenv("CORS_ORIGINS"); o != "" {
		origins = splitAndTrim(o)
	}

	return &Config{
		Port:        port,
		DBPath:      dbPath,
		CORSOrigins: origins,
	}
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range splitOn(s, ',') {
		t := trim(part)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func splitOn(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func trim(s string) string {
	i, j := 0, len(s)
	for i < j && s[i] == ' ' {
		i++
	}
	for j > i && s[j-1] == ' ' {
		j--
	}
	return s[i:j]
}
