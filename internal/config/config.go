package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port        string
	DBDriver    string
	DatabaseURL string
	CacheTTL    string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DBDriver:    getEnv("DB_DRIVER", "sqlite"),
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
	d, err := time.ParseDuration(c.CacheTTL)
	if err != nil || d <= 0 {
		log.Printf("invalid CACHE_TTL %q, using default 30m: %v", c.CacheTTL, err)
		return 30 * time.Minute
	}
	return d
}
