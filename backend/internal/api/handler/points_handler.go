package handler

import (
	"net/http"

	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	Service *service.PointsService
}

func (h *PointsHandler) GetTemplates(c *gin.Context) {
	var templates []model.PointsTemplate
	repository.DB.Find(&templates)
	c.JSON(http.StatusOK, templates)
}

func (h *PointsHandler) SaveTemplate(c *gin.Context) {
	var template model.PointsTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repository.DB.Save(&template)
	c.JSON(http.StatusOK, template)
}

func (h *PointsHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.PointsTemplate{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *PointsHandler) ManagePoints(c *gin.Context) {
	var input struct {
		UserID uint   `json:"user_id"`
		Amount int    `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operatorID := c.MustGet("userID").(uint)
	err := h.Service.UpdatePoints(input.UserID, input.Amount, input.Reason, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "points updated"})
}

func (h *PointsHandler) GetPointsRecords(c *gin.Context) {
	parentID := c.MustGet("userID").(uint)

	// Get children IDs
	var bindings []model.ParentChild
	repository.DB.Where("parent_id = ?", parentID).Find(&bindings)
	childIDs := make([]uint, 0)
	for _, b := range bindings {
		childIDs = append(childIDs, b.ChildID)
	}

	if len(childIDs) == 0 {
		c.JSON(http.StatusOK, []model.PointsRecord{})
		return
	}

	var records []model.PointsRecord
	// Find records where UserID is in children AND OperatorID is the parent
	repository.DB.Where("user_id IN ? AND operator_id = ?", childIDs, parentID).
		Preload("User"). // Load the child user details
		Preload("Operator"). // Load the parent (operator) details
		Order("created_at desc").
		Find(&records)

	c.JSON(http.StatusOK, records)
}

func (h *PointsHandler) GetPointsHistory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var records []model.PointsRecord
	repository.DB.Preload("Operator").Where("user_id = ?", userID).Order("created_at desc").Find(&records)

	c.JSON(http.StatusOK, records)
}
