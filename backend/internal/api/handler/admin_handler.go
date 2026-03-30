package handler

import (
	"baby-fans/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func (h *AdminHandler) ReloadConfig(c *gin.Context) {
	if err := config.ReloadConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "config reloaded",
		"app":     config.Cfg.App,
	})
}
