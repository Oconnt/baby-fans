package handler

import (
	"net/http"
	"time"

	"baby-fans/internal/model"
	"baby-fans/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskTemplateService *service.TaskTemplateService
	TaskService         *service.TaskService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		TaskTemplateService: &service.TaskTemplateService{},
		TaskService:         &service.TaskService{},
	}
}

// Parent APIs

// GetTaskTemplates 获取任务模版列表
func (h *TaskHandler) GetTaskTemplates(c *gin.Context) {
	templates, err := h.TaskTemplateService.GetTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// CreateTaskTemplate 创建任务模版
func (h *TaskHandler) CreateTaskTemplate(c *gin.Context) {
	var template model.TaskTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.TaskTemplateService.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, template)
}

// DeleteTaskTemplate 删除任务模版
func (h *TaskHandler) DeleteTaskTemplate(c *gin.Context) {
	var id struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.TaskTemplateService.DeleteTemplate(id.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetParentTasks 获取家长发布的所有任务
func (h *TaskHandler) GetParentTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	tasks, err := h.TaskService.GetTasksByPublisher(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// CreateTask 创建/派发任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Points      int    `json:"points" binding:"required"`
		HandlerID   uint   `json:"handler_id" binding:"required"`
		ExpireTime  string `json:"expire_time" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publisherID := c.GetUint("user_id")
	expireTime, err := time.Parse(time.RFC3339, req.ExpireTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expire_time format"})
		return
	}

	task := &model.Task{
		Name:        req.Name,
		Description: req.Description,
		Points:      req.Points,
		Status:      model.TaskPending,
		PublisherID: publisherID,
		HandlerID:   req.HandlerID,
		PublishTime: time.Now(),
		ExpireTime:  expireTime,
	}

	if err := h.TaskService.PublishTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UpdateTaskStatus 更新任务状态
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	var id struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TaskService.UpdateTaskStatus(id.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Child APIs

// GetTodayTasks 获取今日任务
func (h *TaskHandler) GetTodayTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	tasks, err := h.TaskService.GetTodayTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetChildTasks 获取孩子的所有任务
func (h *TaskHandler) GetChildTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	tasks, err := h.TaskService.GetTasksByChild(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTaskDetail 获取任务详情
func (h *TaskHandler) GetTaskDetail(c *gin.Context) {
	var id struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.TaskService.GetTaskByID(id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CompleteTask 完成任务
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	var id struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get task first to check ownership
	task, err := h.TaskService.GetTaskByID(id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify the task belongs to the current child
	userID := c.GetUint("user_id")
	if task.HandlerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	// Update task status to completed
	if err := h.TaskService.UpdateTaskStatus(id.ID, model.TaskCompleted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Award points to child
	pointsService := &service.PointsService{}
	if err := pointsService.UpdatePoints(task.HandlerID, task.Points, "完成任务: "+task.Name, task.PublisherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task completed", "points_awarded": task.Points})
}
