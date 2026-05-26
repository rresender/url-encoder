package main

import (
	"github.com/rresender/url-enconder/internal/cache"
	"github.com/rresender/url-enconder/internal/config"
	"github.com/rresender/url-enconder/internal/controller"
	"github.com/rresender/url-enconder/internal/db"
	"github.com/rresender/url-enconder/internal/model"
	"github.com/rresender/url-enconder/internal/repository"
	"github.com/rresender/url-enconder/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	db, err := db.Open(cfg)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.EncodeURL{}, &model.SequenceCounter{})

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
		api.GET("/resolve/:short_url", controller.ResolveEncodeURL)
	}

	router.Run(":" + cfg.Port)
}
