package handler

import (
	"baby-fans/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct{}

func (h *VersionHandler) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":      config.Cfg.App.Version,
		"build":        config.Cfg.App.BuildNumber,
		"update_url":   config.Cfg.App.UpdateURL,
		"force_update": config.Cfg.App.ForceUpdate,
	})
}
