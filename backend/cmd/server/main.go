package main

import (
	"log"
	"net/http"
	"time"

	"baby-fans/config"
	"baby-fans/internal/api"
	"baby-fans/internal/repository"
	"baby-fans/internal/service"

	"golang.org/x/crypto/acme/autocert"
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

	domain := config.Cfg.Server.Domain
	certDir := config.Cfg.Server.CertDir
	email := config.Cfg.Server.Email

	log.Printf("Starting HTTPS server on port %s with Let's Encrypt", port)
	log.Printf("Domain: %s, CertDir: %s, Email: %s", domain, certDir, email)

	// Setup autocert for Let's Encrypt
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Email:      email,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
	}

	// Start HTTP server on port 80 for ACME challenge
	go func() {
		log.Printf("HTTP server for ACME challenge on port 80")
		if err := http.ListenAndServe(":80", m.HTTPHandler(nil)); err != nil {
			log.Printf("ACME HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server with Let's Encrypt certificates
	tlsConfig := m.TLSConfig()
	tlsConfig.MinVersion = 0 // Allow default
	srv := &http.Server{
		Addr:      ":" + port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

