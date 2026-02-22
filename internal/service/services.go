package service

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID uint       `json:"user_id"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct{}

func (s *AuthService) LoginByFace(name string, photoContent []byte) (string, string, error) {
	var user model.User
	err := repository.DB.Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto-register as child for mock face login simplicity
			user = model.User{
				Name: name,
				Role: model.RoleChild,
			}
			if err := repository.DB.Create(&user).Error; err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	// Mock face recognition: Save photo
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%d.jpg", user.ID, timestamp)
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	photoPath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(photoPath, photoContent, 0644); err != nil {
		return "", "", err
	}

	// Create FaceLog
	faceLog := model.FaceLog{
		UserID:    user.ID,
		PhotoPath: photoPath,
	}
	if err := repository.DB.Create(&faceLog).Error; err != nil {
		return "", "", err
	}

	// Photo rotation: Keep only 7 most recent
	var logs []model.FaceLog
	repository.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&logs)
	if len(logs) > 7 {
		for i := 7; i < len(logs); i++ {
			os.Remove(logs[i].PhotoPath)
			repository.DB.Delete(&logs[i])
		}
	}

	// Generate login code
	loginCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	user.LoginCode = loginCode
	repository.DB.Save(&user)

	token, err := s.generateToken(user)
	return token, loginCode, err
}

func (s *AuthService) LoginByCode(code string) (string, error) {
	var user model.User
	if err := repository.DB.Where("login_code = ?", code).First(&user).Error; err != nil {
		return "", errors.New("invalid login code")
	}

	// 演示版本：不再清空登录码，方便重复测试
	// user.LoginCode = ""
	// repository.DB.Save(&user)

	return s.generateToken(user)
}

func (s *AuthService) generateToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

type PointsService struct{}

func (s *PointsService) UpdatePoints(userID uint, amount int, reason string, operatorID uint) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		user.Points += amount
		if user.Points < 0 {
			return errors.New("insufficient points")
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		record := model.PointsRecord{
			UserID:     userID,
			Amount:     amount,
			Reason:     reason,
			OperatorID: operatorID,
		}
		return tx.Create(&record).Error
	})
}

type ShopService struct{}

func (s *ShopService) ExchangeItem(userID uint, itemID uint) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		var item model.ShopItem
		if err := tx.First(&item, itemID).Error; err != nil {
			return err
		}

		if item.Stock <= 0 {
			return errors.New("item out of stock")
		}

		if user.Points < item.Price {
			return errors.New("insufficient points")
		}

		// Deduct points
		user.Points -= item.Price
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Update stock
		item.Stock -= 1
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		// Create redemption
		redemption := model.Redemption{
			UserID: userID,
			ItemID: itemID,
			Status: model.RedemptionPending,
		}
		if err := tx.Create(&redemption).Error; err != nil {
			return err
		}

		// Create points record
		record := model.PointsRecord{
			UserID: userID,
			Amount: -item.Price,
			Reason: "Exchange: " + item.Name,
		}
		return tx.Create(&record).Error
	})
}

func (s *ShopService) ConfirmRedemption(redemptionID uint) error {
	return repository.DB.Model(&model.Redemption{}).Where("id = ?", redemptionID).Update("status", model.RedemptionCompleted).Error
}

func (s *ShopService) CleanupEmptyStockItems() {
	// Items with 0 stock for more than 24 hours should be deleted
	threshold := time.Now().Add(-24 * time.Hour)
	repository.DB.Where("stock = 0 AND updated_at < ?", threshold).Delete(&model.ShopItem{})
}
