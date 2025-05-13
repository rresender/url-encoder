package main

import (
	"github.com/rresender/url-enconder/internal/cache"
	"github.com/rresender/url-enconder/internal/config"
	"github.com/rresender/url-enconder/internal/controller"
	"github.com/rresender/url-enconder/internal/model"
	"github.com/rresender/url-enconder/internal/repository"
	"github.com/rresender/url-enconder/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	cfg := config.LoadConfig()

	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.EncodeURL{})

	repo := repository.NewEncodeURLRepository(db)
	cache := cache.NewInMemoryTTLCache(cfg.GetCacheTTL())
	service := service.NewEncodeURLService(repo, cache)
	controller := controller.NewEncodeURLController(service)

	router := gin.Default()

	// Middleware for tenant logging
	router.Use(func(ctx *gin.Context) {
		tenantID := ctx.GetHeader("X-Tenant-ID")
		if tenantID != "" {
			ctx.Set("tenant_id", tenantID)
		}
		ctx.Next()
	})

	api := router.Group("encoder/api/v1")
	{
		api.POST("/encode", controller.CreateEncodeURL)
	}

	router.Run(":" + cfg.Port)
}
