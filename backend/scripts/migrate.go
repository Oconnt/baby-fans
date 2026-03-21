// Migration script - run before starting the application to initialize the database
package main

import (
	"log"

	"baby-fans/config"
	"baby-fans/internal/model"
	"baby-fans/internal/repository"
)

func main() {
	log.Println("Starting database migration...")

	// Load config to get DB path
	config.LoadConfig()

	// Initialize DB connection
	repository.InitDB()

	// AutoMigrate all models
	err := repository.DB.AutoMigrate(
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
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}
