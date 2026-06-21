package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rresender/url-enconder/internal/model"
	"github.com/rresender/url-enconder/internal/service"
)

type EncodeURLController struct {
	service service.EncodeURLService
}

func NewEncodeURLController(service service.EncodeURLService) *EncodeURLController {
	return &EncodeURLController{service: service}
}

func (c *EncodeURLController) CreateEncodeURL(ctx *gin.Context) {
	var request model.CreateEncodeURLRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.TenantID == "" {
		request.TenantID = ctx.GetHeader("X-Tenant-ID")
	}

	if request.TenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tenant ID is required"})
		return
	}

	response, err := c.service.CreateEncodeURL(&request)
	if err != nil {
		if errors.Is(err, service.ErrInvalidEncodingStrategy) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *EncodeURLController) ResolveEncodeURL(ctx *gin.Context) {
	shortURL := ctx.Param("short_url")
	if shortURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "short_url is required"})
		return
	}

	originalURL, err := c.service.GetOriginalURL(shortURL)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"short_url":    shortURL,
		"original_url": originalURL,
	})
}
