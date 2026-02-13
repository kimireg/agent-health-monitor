package main

import (
	"log"
	"net/http"

	"github.com/kimireg/jason-frontpage/internal/config"
	"github.com/kimireg/jason-frontpage/internal/renderer"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize renderer
	r := renderer.New(cfg)

	// Setup HTTP routes
	http.HandleFunc("/", r.HomeHandler)
	http.HandleFunc("/health", r.HealthHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Start server
	log.Printf("Starting Jason Front Page server on %s", cfg.ServerAddr)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, nil))
}