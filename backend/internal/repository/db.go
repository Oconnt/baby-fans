package repository

import (
	"log"

	"baby-fans/config"
	"baby-fans/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	dsn := config.Cfg.DB.GetMySQLDSN()
	log.Printf("Connecting to MySQL: %s:%d/%s", config.Cfg.DB.Host, config.Cfg.DB.Port, config.Cfg.DB.Name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// AutoMigrate (idempotent - only adds missing tables/columns)
	if err := DB.AutoMigrate(
		&model.User{},
		&model.UserBinding{},
		&model.ParentChild{},
		&model.PointsRecord{},
		&model.ShopItem{},
		&model.Redemption{},
		&model.FaceLog{},
		&model.PointsTemplate{},
		&model.TaskTemplate{},
		&model.Task{},
	); err != nil {
		log.Fatalf("Auto migrate failed: %v", err)
	}

}
