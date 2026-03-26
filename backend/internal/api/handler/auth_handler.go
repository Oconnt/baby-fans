package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"baby-fans/config"
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
		log.Printf("[LoginFace] missing name or photo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and photo are required"})
		return
	}

	f, _ := file.Open()
	defer f.Close()
	content := make([]byte, file.Size)
	f.Read(content)

	token, code, err := h.Service.LoginByFace(name, content)
	if err != nil {
		log.Printf("[LoginFace] LoginByFace error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询用户信息以返回角色
	var user model.User
	if err := repository.DB.Where("name = ?", name).First(&user).Error; err != nil {
		log.Printf("[LoginFace] query user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
		log.Printf("[Register] bind error: %v", err)
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
		log.Printf("[Register] create user error: %v", err)
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
		log.Printf("[LoginCode] invalid login code: %s", code)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login code"})
		return
	}

	token, err := h.Service.LoginByCode(code)
	if err != nil {
		log.Printf("[LoginCode] LoginByCode error: %v", err)
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
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[GetOverview] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[GetOverview] userID invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		log.Printf("[GetOverview] First user error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var records []model.PointsRecord
	if err := repository.DB.Preload("Operator").Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&records).Error; err != nil {
		log.Printf("[GetOverview] query records error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all bound parent names
	var parentNames []string
	var bindings []model.ParentChild
	if err := repository.DB.Where("child_id = ?", userID).Find(&bindings).Error; err != nil {
		log.Printf("[GetOverview] query bindings error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[UpdateProfile] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[UpdateProfile] userID invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var input struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[UpdateProfile] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		log.Printf("[UpdateProfile] First user error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	if err := repository.DB.Save(&user).Error; err != nil {
		log.Printf("[UpdateProfile] Save user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[UploadAvatar] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[UploadAvatar] userID invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("[UploadAvatar] FormFile error: %v", err)
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
		log.Printf("[UploadAvatar] SaveUploadedFile error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user avatar_url
	var user model.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		log.Printf("[UploadAvatar] First user error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	avatarURL := "https://" + config.Cfg.Server.Domain + "/storage/uploads/" + filename
	user.AvatarURL = avatarURL
	if err := repository.DB.Save(&user).Error; err != nil {
		log.Printf("[UploadAvatar] Save user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": avatarURL})
}

func (h *AuthHandler) GetChildren(c *gin.Context) {
	parentIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[GetChildren] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	parentID, ok := parentIDVal.(uint)
	if !ok {
		log.Printf("[GetChildren] userID invalid type: %T", parentIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var bindings []model.ParentChild
	if err := repository.DB.Where("parent_id = ?", parentID).Find(&bindings).Error; err != nil {
		log.Printf("[GetChildren] query bindings error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var children []model.User
	if len(bindings) > 0 {
		childIDs := make([]uint, len(bindings))
		for i, b := range bindings {
			childIDs[i] = b.ChildID
		}
		if err := repository.DB.Where("id IN ?", childIDs).Find(&children).Error; err != nil {
			log.Printf("[GetChildren] query children error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, children)
}

func (h *AuthHandler) BindChildByCode(c *gin.Context) {
	parentIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[BindChildByCode] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	parentID, ok := parentIDVal.(uint)
	if !ok {
		log.Printf("[BindChildByCode] userID invalid type: %T", parentIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var input struct {
		LoginCode string `json:"login_code"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[BindChildByCode] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var child model.User
	if err := repository.DB.Where("login_code = ? AND role = ?", input.LoginCode, model.RoleChild).First(&child).Error; err != nil {
		log.Printf("[BindChildByCode] child not found by code: %s", input.LoginCode)
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该登录码对应的孩子"})
		return
	}

	var existing model.ParentChild
	if err := repository.DB.Where("parent_id = ? AND child_id = ?", parentID, child.ID).First(&existing).Error; err == nil {
		log.Printf("[BindChildByCode] already bound: parent=%d child=%d", parentID, child.ID)
		c.JSON(http.StatusConflict, gin.H{"error": "已绑定该孩子"})
		return
	}

	binding := model.ParentChild{ParentID: parentID, ChildID: child.ID}
	if err := repository.DB.Create(&binding).Error; err != nil {
		log.Printf("[BindChildByCode] create binding error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "绑定成功", "child": child})
}

func (h *AuthHandler) UnbindChild(c *gin.Context) {
	parentIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[UnbindChild] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	parentID, ok := parentIDVal.(uint)
	if !ok {
		log.Printf("[UnbindChild] userID invalid type: %T", parentIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	childIDStr := c.Param("id")
	childID, _ := strconv.ParseUint(childIDStr, 10, 32)

	result := repository.DB.Where("parent_id = ? AND child_id = ?", parentID, uint(childID)).Delete(&model.ParentChild{})
	if result.RowsAffected == 0 {
		log.Printf("[UnbindChild] binding not found: parent=%d child=%d", parentID, childID)
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到绑定关系"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "解绑成功"})
}

func (h *AuthHandler) GetFaceLogs(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	var logs []model.FaceLog
	if err := repository.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&logs).Error; err != nil {
		log.Printf("[GetFaceLogs] query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
