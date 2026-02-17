package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kimireg/jason-frontpage/internal/config"
	"github.com/kimireg/jason-frontpage/internal/db"
	"github.com/kimireg/jason-frontpage/internal/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/workspace/data/blog.db"
	}
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize handlers
	h, err := handlers.New(cfg, database)
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Setup HTTP routes
	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/health", h.HealthHandler)

	// Blog web routes
	http.HandleFunc("/blog", h.BlogListHandler)
	http.HandleFunc("/blog/new", h.BlogNewHandler)
	http.HandleFunc("/blog/edit/", h.BlogEditHandler)
	http.HandleFunc("/blog/delete", h.BlogDeleteHandler)
	http.HandleFunc("/blog/post/", h.BlogPostHandler)
	http.HandleFunc("/blog/tag/", h.BlogTagHandler)

	// Blog API routes
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h.APICreatePost(w, r)
		} else {
			h.APIGetPosts(w, r)
		}
	})
	http.HandleFunc("/api/post", h.APIGetPost)
	http.HandleFunc("/api/stats", h.APIGetStats)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Start server
	log.Printf("Starting Jason Front Page server on %s", cfg.ServerAddr)
	log.Printf("Database: %s", dbPath)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, nil))
}
