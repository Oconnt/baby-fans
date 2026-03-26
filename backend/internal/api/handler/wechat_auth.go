package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"baby-fans/config"
	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Role     string `json:"role"` // parent or child
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar_url"`
}

func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[WeChatLogin] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Call WeChat API
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"appid":      config.Cfg.WeChat.AppID,
			"secret":     config.Cfg.WeChat.AppSecret,
			"js_code":    req.Code,
			"grant_type": "authorization_code",
		}).
		Get("https://api.weixin.qq.com/sns/jscode2session")

	if err != nil {
		log.Printf("[WeChatLogin] call WeChat API error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call WeChat API"})
		return
	}

	var wxResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		Errcode    int    `json:"errcode"`
		Errmsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(resp.Body(), &wxResp); err != nil {
		log.Printf("[WeChatLogin] unmarshal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse WeChat response"})
		return
	}

	if wxResp.Errcode != 0 || wxResp.OpenID == "" {
		// MOCK FOR TESTING IF WECHAT CONFIG IS EMPTY (Fallback for local dev)
		if config.Cfg == nil || config.Cfg.WeChat.AppID == "" {
			wxResp.OpenID = "mock_openid_" + req.Code
		} else {
			log.Printf("[WeChatLogin] WeChat error: errcode=%d errmsg=%s", wxResp.Errcode, wxResp.Errmsg)
			c.JSON(http.StatusUnauthorized, gin.H{"error": wxResp.Errmsg})
			return
		}
	}

	// 2. Find or create user
	var user model.User
	result := repository.DB.Where("open_id = ?", wxResp.OpenID).First(&user)
	if result.Error != nil {
		// Create new user
		role := model.RoleParent
		if req.Role == string(model.RoleChild) {
			role = model.RoleChild
		}
		user = model.User{
			OpenID:    wxResp.OpenID,
			UnionID:   wxResp.UnionID,
			Nickname:  req.Nickname,
			AvatarURL: req.Avatar,
			Role:      role,
			Name:      req.Nickname,
		}
		if user.Name == "" {
			user.Name = "User_" + wxResp.OpenID[:6]
		}
		if err := repository.DB.Create(&user).Error; err != nil {
			log.Printf("[WeChatLogin] create user error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// 3. Generate JWT
	secret := "super_secret_baby_fans_key" // fallback
	expireHours := 24
	if config.Cfg != nil {
		secret = config.Cfg.JWT.Secret
		expireHours = config.Cfg.JWT.Expire
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("[WeChatLogin] sign token error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func (h *AuthHandler) GenerateBindingCode(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[GenerateBindingCode] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	parentID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[GenerateBindingCode] userID invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	// Generate random 6 character hex code
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("[GenerateBindingCode] rand.Read error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}
	code := hex.EncodeToString(bytes)

	binding := model.UserBinding{
		ParentID: parentID,
		BindCode: code,
		Status:   "pending",
	}

	if err := repository.DB.Create(&binding).Error; err != nil {
		log.Printf("[GenerateBindingCode] create binding error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save binding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bind_code": code})
}

func (h *AuthHandler) AcceptBinding(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		log.Printf("[AcceptBinding] userID not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	childID, ok := userIDVal.(uint)
	if !ok {
		log.Printf("[AcceptBinding] userID invalid type: %T", userIDVal)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID has invalid type"})
		return
	}

	var req struct {
		BindCode string `json:"bind_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AcceptBinding] bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var binding model.UserBinding
	if err := repository.DB.Where("bind_code = ? AND status = ?", req.BindCode, "pending").First(&binding).Error; err != nil {
		log.Printf("[AcceptBinding] binding not found: %s", req.BindCode)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired binding code"})
		return
	}

	binding.ChildID = childID
	binding.Status = "accepted"

	if err := repository.DB.Save(&binding).Error; err != nil {
		log.Printf("[AcceptBinding] save binding error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update binding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Binding successful"})
}
