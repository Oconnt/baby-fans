package handler

import (
	"log"
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
		log.Printf("[SaveTemplate] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.DB.Save(&template).Error; err != nil {
		log.Printf("[SaveTemplate] save error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *PointsHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DB.Delete(&model.PointsTemplate{}, id).Error; err != nil {
		log.Printf("[DeleteTemplate] delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *PointsHandler) ManagePoints(c *gin.Context) {
	var input struct {
		UserID uint   `json:"user_id"`
		Amount int    `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[ManagePoints] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operatorIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[ManagePoints] userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}
	operatorID, ok := operatorIDVal.(uint)
	if !ok {
		log.Printf("[ManagePoints] userID has invalid type: %T", operatorIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	err := h.Service.UpdatePoints(input.UserID, input.Amount, input.Reason, operatorID)
	if err != nil {
		log.Printf("[ManagePoints] UpdatePoints error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "points updated"})
}

func (h *PointsHandler) GetPointsRecords(c *gin.Context) {
	parentIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[GetPointsRecords] userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	parentID, ok := parentIDVal.(uint)
	if !ok {
		log.Printf("[GetPointsRecords] userID has invalid type: %T", parentIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	// Get children IDs
	var bindings []model.ParentChild
	if err := repository.DB.Where("parent_id = ?", parentID).Find(&bindings).Error; err != nil {
		log.Printf("[GetPointsRecords] query bindings error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	childIDs := make([]uint, 0)
	for _, b := range bindings {
		childIDs = append(childIDs, b.ChildID)
	}

	if len(childIDs) == 0 {
		c.JSON(http.StatusOK, []model.PointsRecord{})
		return
	}

	var records []model.PointsRecord
	if err := repository.DB.Where("user_id IN ? AND operator_id = ?", childIDs, parentID).
		Preload("User").
		Preload("Operator").
		Order("created_at desc").
		Find(&records).Error; err != nil {
		log.Printf("[GetPointsRecords] query records error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *PointsHandler) GetPointsHistory(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[GetPointsHistory] userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[GetPointsHistory] userID has invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var records []model.PointsRecord
	if err := repository.DB.Preload("Operator").Where("user_id = ?", userID).Order("created_at desc").Find(&records).Error; err != nil {
		log.Printf("[GetPointsHistory] query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}
