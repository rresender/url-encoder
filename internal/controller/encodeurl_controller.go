package controller

import (
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
