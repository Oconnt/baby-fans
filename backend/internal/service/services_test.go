package service

import (
	"fmt"
	"os"
	"testing"
	"time"

	"baby-fans/internal/model"
	"baby-fans/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) string {
	dbName := fmt.Sprintf("test_%d.db", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	repository.DB = db
	err = repository.DB.AutoMigrate(
		&model.User{},
		&model.PointsRecord{},
		&model.FaceLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return dbName
}

func TestPhotoRotation(t *testing.T) {
	dbName := setupTestDB(t)
	defer os.Remove(dbName)

	// Ensure uploads directory is clean for test
	uploadDir := "uploads"
	os.RemoveAll(uploadDir)
	os.Mkdir(uploadDir, 0755)
	defer os.RemoveAll(uploadDir)

	authService := &AuthService{}
	name := "RotationChild"

	// Simulate 8 logins
	for i := 1; i <= 8; i++ {
		_, _, err := authService.LoginByFace(name, []byte(fmt.Sprintf("photo-%d", i)))
		if err != nil {
			t.Fatalf("Login %d failed: %v", i, err)
		}
		// Ensure different timestamps for stable sorting in SQLite
		time.Sleep(10 * time.Millisecond)
	}

	// Verify database record count
	var user model.User
	repository.DB.Where("name = ?", name).First(&user)

	var count int64
	repository.DB.Model(&model.FaceLog{}).Where("user_id = ?", user.ID).Count(&count)

	if count != 7 {
		t.Errorf("Expected 7 face logs in DB, got %d", count)
	}

	// Verify physical file count
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		t.Fatalf("Failed to read uploads dir: %v", err)
	}

	if len(files) != 7 {
		t.Errorf("Expected 7 files in uploads, got %d", len(files))
	}
}

func TestPointsTransaction(t *testing.T) {
	dbName := setupTestDB(t)
	defer os.Remove(dbName)

	pointsService := &PointsService{}

	// Create a user
	user := model.User{Name: "PointsChild", Role: model.RoleChild, Points: 100}
	repository.DB.Create(&user)

	// Add points
	err := pointsService.UpdatePoints(user.ID, 50, "Good behavior", 1)
	if err != nil {
		t.Fatalf("UpdatePoints failed: %v", err)
	}

	// Verify points
	var updatedUser model.User
	repository.DB.First(&updatedUser, user.ID)
	if updatedUser.Points != 150 {
		t.Errorf("Expected 150 points, got %d", updatedUser.Points)
	}

	// Verify record
	var record model.PointsRecord
	repository.DB.Where("user_id = ?", user.ID).First(&record)
	if record.Amount != 50 || record.Reason != "Good behavior" {
		t.Errorf("Record mismatch: %+v", record)
	}

	// Test insufficient points
	err = pointsService.UpdatePoints(user.ID, -200, "Spending too much", 1)
	if err == nil {
		t.Error("Expected error for insufficient points, got nil")
	}
}
