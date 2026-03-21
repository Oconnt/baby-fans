package repository

import (
	"log"

	"baby-fans/config"
	"baby-fans/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Determine database type from config
	dbType := "sqlite"
	if config.Cfg != nil && config.Cfg.DB.Type != "" {
		dbType = config.Cfg.DB.Type
	}

	if dbType == "mysql" && config.Cfg != nil {
		// Connect to MySQL
		dsn := config.Cfg.DB.GetMySQLDSN()
		log.Printf("Connecting to MySQL: %s:%d/%s", config.Cfg.DB.Host, config.Cfg.DB.Port, config.Cfg.DB.Name)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		// Fallback to SQLite
		sqliteDsn := "baby-fans.db"
		if config.Cfg != nil {
			sqliteDsn = config.Cfg.DB.GetSQLiteDSN()
		}
		log.Printf("Using SQLite: %s", sqliteDsn)
		DB, err = gorm.Open(sqlite.Open(sqliteDsn), &gorm.Config{})
	}

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

	// Initialize persistent accounts (idempotent)
	initPersistentAccounts()
}

func initPersistentAccounts() {
	// 1. Upsert parent account
	parent := model.User{
		Name:      "超级家长",
		Role:      model.RoleParent,
		LoginCode: "888888",
	}
	// Use ON DUPLICATE KEY UPDATE for MySQL or INSERT OR IGNORE for SQLite
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "login_code", "updated_at"}),
	}).Create(&parent)

	// 2. Upsert child account
	child := model.User{
		Name:      "小明同学",
		Role:      model.RoleChild,
		LoginCode: "666666",
		Points:    200,
	}
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "login_code", "points", "updated_at"}),
	}).Create(&child)

	log.Println("Persistent accounts initialized")
}
