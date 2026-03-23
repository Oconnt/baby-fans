package main

import (
	"log"
	"time"

	"baby-fans/config"
	"baby-fans/internal/api"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"
)

func main() {
	config.LoadConfig()
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
	port := config.Cfg.Server.Port
	if port == "" {
		port = "18081"
	}
	certFile := "certs/server.crt"
	keyFile := "certs/server.key"
	log.Printf("Starting HTTPS server on port %s", port)
	log.Printf("Certificate: %s, Key: %s", certFile, keyFile)
	if err := r.RunTLS(":"+port, certFile, keyFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

