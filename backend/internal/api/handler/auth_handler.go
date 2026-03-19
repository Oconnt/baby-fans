package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		"name":       user.Name,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Role     string `json:"role"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	role := model.RoleParent
	if input.Role == "child" {
		role = model.RoleChild
	}

	// Generate login code
	loginCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Generate a random unique name
	randomName := fmt.Sprintf("User%d", rand.Intn(999999))

	user := model.User{
		Name:      randomName,
		Role:      role,
		Nickname:  input.Nickname,
		LoginCode: loginCode,
		Points:    0,
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "registered successfully",
		"login_code": loginCode,
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
		"token":      token,
		"role":       user.Role,
		"user_id":    user.ID,
		"name":       user.Name,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
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
	repository.DB.Preload("Operator").Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&records)

	// Get all bound parent names
	var parentNames []string
	var bindings []model.ParentChild
	repository.DB.Where("child_id = ?", userID).Find(&bindings)
	for _, binding := range bindings {
		var parent model.User
		if err := repository.DB.First(&parent, binding.ParentID).Error; err == nil {
			name := parent.Nickname
			if name == "" {
				name = parent.Name
			}
			parentNames = append(parentNames, name)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"points":       user.Points,
		"records":      records,
		"parent_names": parentNames,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	repository.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Generate filename: avatars/{user_id}.{ext}
	ext := ".jpg"
	name := file.Filename
	if len(name) > 4 {
		ext = "." + strings.ToLower(name[len(name)-4:])
		if ext == ".jpeg" {
			ext = ".jpg"
		}
	}
	filename := fmt.Sprintf("avatars/%d%s", userID, ext)
	dst := fmt.Sprintf("./storage/uploads/%s", filename)

	// Ensure directory exists
	os.MkdirAll("./storage/uploads/avatars", 0755)

	// Remove existing file if exists
	os.Remove(dst)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user avatar_url
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	avatarURL := "http://localhost:18081/storage/uploads/" + filename
	user.AvatarURL = avatarURL
	repository.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"url": avatarURL})
}

func (h *AuthHandler) GetChildren(c *gin.Context) {
	parentID := c.MustGet("userID").(uint)
	var bindings []model.ParentChild
	repository.DB.Where("parent_id = ?", parentID).Find(&bindings)

	var children []model.User
	if len(bindings) > 0 {
		childIDs := make([]uint, len(bindings))
		for i, b := range bindings {
			childIDs[i] = b.ChildID
		}
		repository.DB.Where("id IN ?", childIDs).Find(&children)
	}
	c.JSON(http.StatusOK, children)
}

func (h *AuthHandler) BindChildByCode(c *gin.Context) {
	parentID := c.MustGet("userID").(uint)
	var input struct {
		LoginCode string `json:"login_code"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var child model.User
	if err := repository.DB.Where("login_code = ? AND role = ?", input.LoginCode, model.RoleChild).First(&child).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该登录码对应的孩子"})
		return
	}

	var existing model.ParentChild
	if err := repository.DB.Where("parent_id = ? AND child_id = ?", parentID, child.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "已绑定该孩子"})
		return
	}

	binding := model.ParentChild{ParentID: parentID, ChildID: child.ID}
	repository.DB.Create(&binding)
	c.JSON(http.StatusOK, gin.H{"message": "绑定成功", "child": child})
}

func (h *AuthHandler) UnbindChild(c *gin.Context) {
	parentID := c.MustGet("userID").(uint)
	childIDStr := c.Param("id")
	childID, _ := strconv.ParseUint(childIDStr, 10, 32)

	result := repository.DB.Where("parent_id = ? AND child_id = ?", parentID, uint(childID)).Delete(&model.ParentChild{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到绑定关系"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "解绑成功"})
}

func (h *AuthHandler) GetFaceLogs(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	var logs []model.FaceLog
	repository.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&logs)
	c.JSON(http.StatusOK, logs)
}
