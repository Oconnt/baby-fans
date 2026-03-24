package main

import (
	"context"
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

	// Setup autocert for Let's Encrypt HTTP-01 challenge
	m := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		Email:      email,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			return nil // Allow all hosts for certificate generation
		},
	}

	// Start HTTP server on port 80 for ACME HTTP-01 challenge and health check
	go func() {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		httpMux.HandleFunc("/.well-known/acme-challenge/", m.HTTPHandler(nil).ServeHTTP)
		log.Printf("HTTP server for ACME HTTP-01 challenge on port 80")
		if err := http.ListenAndServe(":80", httpMux); err != nil {
			log.Printf("ACME HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server with Let's Encrypt certificates
	tlsConfig := m.TLSConfig()
	srv := &http.Server{
		Addr:      ":" + port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
