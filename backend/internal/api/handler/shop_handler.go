package handler

import (
	"net/http"
	"strconv"

	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"

	"github.com/gin-gonic/gin"
)

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

func (h *ShopHandler) Cancel(c *gin.Context) {
	redemptionIDStr := c.Param("id")
	redemptionID, _ := strconv.ParseUint(redemptionIDStr, 10, 32)

	err := h.Service.CancelRedemption(uint(redemptionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "redemption cancelled"})
}

func (h *ShopHandler) GetRedemptions(c *gin.Context) {
	var redemptions []model.Redemption
	repository.DB.Preload("User").Preload("Item").Order("created_at desc").Find(&redemptions)
	c.JSON(http.StatusOK, redemptions)
}
