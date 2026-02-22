package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

func (h *AuthHandler) LoginFace(c *gin.Context) {
	name := c.PostForm("name")
	file, _ := c.FormFile("photo")
	if file == nil || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and photo are required"})
		return
	}

	f, _ := file.Open()
	defer f.Close()
	content := make([]byte, file.Size)
	f.Read(content)

	token, code, err := h.Service.LoginByFace(name, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询用户信息以返回角色
	var user model.User
	repository.DB.Where("name = ?", name).First(&user)

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"login_code": code,
		"role":       user.Role,
		"user_id":    user.ID,
	})
}

func (h *AuthHandler) LoginCode(c *gin.Context) {
	code := c.Query("code")
	var user model.User
	if err := repository.DB.Where("login_code = ?", code).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login code"})
		return
	}

	token, err := h.Service.LoginByCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"role":    user.Role,
		"user_id": user.ID,
	})
}

func (h *AuthHandler) GetOverview(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var records []model.PointsRecord
	repository.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"points":  user.Points,
		"records": records,
	})
}

func (h *AuthHandler) GetChildren(c *gin.Context) {
	var children []model.User
	// 简单演示：查询所有 child 角色
	repository.DB.Where("role = ?", model.RoleChild).Find(&children)
	c.JSON(http.StatusOK, children)
}

func (h *AuthHandler) GetFaceLogs(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	var logs []model.FaceLog
	repository.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&logs)
	c.JSON(http.StatusOK, logs)
}

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

type ShopHandler struct {
	Service *service.ShopService
}

func (h *ShopHandler) GetItems(c *gin.Context) {
	var items []model.ShopItem
	repository.DB.Find(&items)
	c.JSON(http.StatusOK, items)
}

func (h *ShopHandler) SaveItem(c *gin.Context) {
	name := c.PostForm("name")
	price, _ := strconv.Atoi(c.PostForm("price"))
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	description := c.PostForm("description")
	idStr := c.PostForm("id")

	var item model.ShopItem
	if idStr != "" {
		id, _ := strconv.ParseUint(idStr, 10, 32)
		repository.DB.First(&item, uint(id))
	}

	item.Name = name
	item.Price = price
	item.Stock = stock
	item.Description = description

	file, err := c.FormFile("image")
	if err == nil {
		filename := fmt.Sprintf("item_%d_%s", time.Now().Unix(), file.Filename)
		dst := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, dst); err == nil {
			item.ImagePath = "/uploads/" + filename
		}
	}

	repository.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (h *ShopHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.ShopItem{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ShopHandler) Exchange(c *gin.Context) {
	itemIDStr := c.Param("id")
	itemID, _ := strconv.ParseUint(itemIDStr, 10, 32)
	userID := c.MustGet("userID").(uint)

	err := h.Service.ExchangeItem(userID, uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exchange successful"})
}

func (h *ShopHandler) Confirm(c *gin.Context) {
	redemptionIDStr := c.Param("id")
	redemptionID, _ := strconv.ParseUint(redemptionIDStr, 10, 32)

	err := h.Service.ConfirmRedemption(uint(redemptionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "redemption confirmed"})
}
