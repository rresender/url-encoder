package db

import (
	"fmt"
	"strings"

	"github.com/rresender/url-enconder/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.DBDriver))

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch driver {
	case "", "sqlite":
		return gorm.Open(sqlite.Open(cfg.DatabaseURL), gormCfg)
	case "postgres", "postgresql":
		return gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q", cfg.DBDriver)
	}
}

