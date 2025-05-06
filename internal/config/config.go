package config

import (
	"os"
	"time"
)

type Config struct {
	Port        string
	DatabaseURL string
	CacheTTL    string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "file:encodeurl.db?cache=shared&mode=rwc"),
		CacheTTL:    getEnv("CACHE_TTL", "30m"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) GetCacheTTL() time.Duration {
	timeVal, _ := time.ParseDuration(c.CacheTTL)
	return timeVal
}
