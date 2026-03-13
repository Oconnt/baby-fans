package main

import (
	"fmt"
	"log"

	"baby-fans/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("baby-fans.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var users []model.User
	db.Find(&users)

	fmt.Println("=== 当前数据库中的用户与登录码 ===")
	for _, u := range users {
		fmt.Printf("姓名: %-10s | 角色: %-8s | 登录码: %s\n", u.Name, u.Role, u.LoginCode)
	}
	fmt.Println("====================================")
}
