package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"baby-fans/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("baby-fans.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	// 1. 创建家长账号
	parentCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	parent := model.User{
		Name:      "超级家长",
		Role:      model.RoleParent,
		LoginCode: parentCode,
	}
	db.Where(model.User{Name: "超级家长"}).FirstOrCreate(&parent)
	// 强制更新 Code 以便提供给用户
	db.Model(&parent).Update("LoginCode", parentCode)

	// 2. 创建孩子账号
	childCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	child := model.User{
		Name:      "小明同学",
		Role:      model.RoleChild,
		LoginCode: childCode,
		Points:    200,
	}
	db.Where(model.User{Name: "小明同学"}).FirstOrCreate(&child)
	db.Model(&child).Update("LoginCode", childCode)

	// 3. 关联关系 (可选)
	// 在本演示版本中，家长端默认可以管理所有孩子角色

	fmt.Println("======================================")
	fmt.Printf("👨‍👩‍👧 家长账号已就绪！\n姓名: %s\n登录码: %s\n", parent.Name, parentCode)
	fmt.Println("--------------------------------------")
	fmt.Printf("🧒 孩子账号已就绪！\n姓名: %s\n登录码: %s\n", child.Name, childCode)
	fmt.Println("======================================")
}
