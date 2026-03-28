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
			UserID:     userID,
			Amount:     -item.Price,
			Reason:     "Exchange: " + item.Name,
			OperatorID: userID,
		}
		return tx.Create(&record).Error
	})
}

func (s *ShopService) ConfirmRedemption(redemptionID uint) error {
	return repository.DB.Model(&model.Redemption{}).Where("id = ?", redemptionID).Update("status", model.RedemptionCompleted).Error
}

func (s *ShopService) CancelRedemption(redemptionID uint) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var redemption model.Redemption
		if err := tx.Preload("Item").First(&redemption, redemptionID).Error; err != nil {
			return err
		}

		// 1. 检查状态，防止重复取消
		if redemption.Status == model.RedemptionCancelled || redemption.Status == model.RedemptionCompleted {
			return errors.New("redemption already processed")
		}

		// 2. 更新状态为已取消
		if err := tx.Model(&redemption).Update("status", model.RedemptionCancelled).Error; err != nil {
			return err
		}

		// 3. 恢复库存 (因为之前兑换扣除了库存)
		if redemption.ItemID > 0 {
			if err := tx.Model(&redemption.Item).Update("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
				return err
			}
		}

		// 4. 返还积分给孩子
		if redemption.UserID > 0 && redemption.ItemID > 0 {
			var user model.User
			if err := tx.First(&user, redemption.UserID).Error; err != nil {
				return err
			}
			user.Points += redemption.Item.Price
			if err := tx.Save(&user).Error; err != nil {
				return err
			}

			// 5. 记录积分变动
			record := model.PointsRecord{
				UserID:     redemption.UserID,
				Amount:     redemption.Item.Price,
				Reason:     "取消兑换: " + redemption.Item.Name,
				OperatorID: redemption.UserID,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ShopService) CleanupEmptyStockItems() {
	// Items with 0 stock for more than 24 hours should be deleted
	threshold := time.Now().Add(-24 * time.Hour)
	repository.DB.Where("stock = 0 AND updated_at < ?", threshold).Delete(&model.ShopItem{})
}

// TaskTemplateService
type TaskTemplateService struct{}

func (s *TaskTemplateService) GetTemplates() ([]model.TaskTemplate, error) {
	var templates []model.TaskTemplate
	err := repository.DB.Order("created_at desc").Find(&templates).Error
	return templates, err
}

func (s *TaskTemplateService) CreateTemplate(template *model.TaskTemplate) error {
	return repository.DB.Create(template).Error
}

func (s *TaskTemplateService) DeleteTemplate(id uint) error {
	return repository.DB.Delete(&model.TaskTemplate{}, id).Error
}

// TaskService
type TaskService struct{}

func (s *TaskService) CreateTask(task *model.Task) error {
	return repository.DB.Create(task).Error
}

func (s *TaskService) GetTasksByChild(childID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := repository.DB.Where("handler_id = ?", childID).
		Preload("Publisher").
		Order("publish_time desc").
		Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) GetTodayTasks(childID uint) ([]model.Task, error) {
	var tasks []model.Task
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := repository.DB.Where("handler_id = ? AND publish_time >= ? AND publish_time < ?", childID, startOfDay, endOfDay).
		Preload("Publisher").
		Order("publish_time desc").
		Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) GetTaskByID(id uint) (*model.Task, error) {
	var task model.Task
	err := repository.DB.Preload("Publisher").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) GetTasksByPublisher(publisherID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := repository.DB.Where("publisher_id = ?", publisherID).
		Preload("Handler").
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func (s *TaskService) UpdateTaskStatus(id uint, status int) error {
	updates := map[string]interface{}{"status": status}
	if status == model.TaskCompleted {
		now := time.Now()
		updates["complete_time"] = &now
	}
	return repository.DB.Model(&model.Task{}).Where("id = ?", id).Updates(updates).Error
}

func (s *TaskService) ExpireOldTasks() error {
	now := time.Now()
	return repository.DB.Model(&model.Task{}).
		Where("status = ? AND expire_time < ?", model.TaskPending, now).
		Updates(map[string]interface{}{"status": model.TaskExpired}).Error
}

func (s *TaskService) PublishTask(task *model.Task) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		// Create the task
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		return nil
	})
}
