package repository

import (
	"log"

	"baby-fans/internal/model"
	"github.com/glebarez/sqlite" // 暂时保留作为备选
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 如果您有密码，请将下面的 DSN 修改为: root:您的密码@tcp(127.0.0.1:3306)/baby_fans...
	dsn := "root:@tcp(127.0.0.1:3306)/baby_fans?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 如果 MySQL 连接失败，自动降级回 SQLite 确保系统可用
	if err != nil {
		log.Println("MySQL 连接失败，正在使用 SQLite 作为备选存储...")
		DB, err = gorm.Open(sqlite.Open("baby-fans.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 自动迁移
	DB.AutoMigrate(
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
	)

	// 强制初始化持久账号
	initPersistentAccounts()
}

func initPersistentAccounts() {
	// 1. 强制初始化家长 (使用 Save 确保即使姓名冲突也会更新登录码)
	var parent model.User
	DB.Where("name = ?", "超级家长").First(&parent)
	parent.Name = "超级家长"
	parent.Role = model.RoleParent
	parent.LoginCode = "888888"
	DB.Save(&parent)

	// 2. 强制初始化孩子
	var child model.User
	DB.Where("name = ?", "小明同学").First(&child)
	child.Name = "小明同学"
	child.Role = model.RoleChild
	child.LoginCode = "666666"
	if child.ID == 0 {
		child.Points = 200
	}
	DB.Save(&child)
}
