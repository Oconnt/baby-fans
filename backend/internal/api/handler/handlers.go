package handler

import (
	"net/http"
	"strconv"

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
		Order("created_at desc").
		Find(&records)

	c.JSON(http.StatusOK, records)
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
	var input struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item model.ShopItem
	if input.ID > 0 {
		repository.DB.First(&item, input.ID)
	}

	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price
	item.Stock = input.Stock

	repository.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (h *ShopHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var input struct {
		Stock int `json:"stock"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item model.ShopItem
	if err := repository.DB.First(&item, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	item.Stock = input.Stock
	repository.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (h *ShopHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	repository.DB.Delete(&model.ShopItem{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ShopHandler) Exchange(c *gin.Context) {
	var input struct {
		ItemID uint `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(uint)

	err := h.Service.ExchangeItem(userID, input.ItemID)
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

func (h *ShopHandler) GetRedemptions(c *gin.Context) {
	var redemptions []model.Redemption
	repository.DB.Preload("User").Preload("Item").Order("created_at desc").Find(&redemptions)
	c.JSON(http.StatusOK, redemptions)
}
