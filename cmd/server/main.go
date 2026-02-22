package main

import (
	"time"

	"baby-fans/internal/api"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"
)

func main() {
	repository.InitDB()

	// Seed some initial data if needed (e.g., a parent user)

	// Start background cleanup task
	shopService := &service.ShopService{}
	go func() {
		for {
			shopService.CleanupEmptyStockItems()
			time.Sleep(1 * time.Hour)
		}
	}()

	r := api.SetupRouter()
	r.Run(":18081")
}
