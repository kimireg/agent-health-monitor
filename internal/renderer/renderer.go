package renderer

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kimireg/jason-frontpage/internal/config"
	"github.com/kimireg/jason-frontpage/pkg/models"
)

// Renderer handles HTML rendering and data extraction
type Renderer struct {
	cfg     *config.Config
	profile *models.AgentProfile
}

// New creates a new renderer instance
func New(cfg *config.Config) *Renderer {
	return &Renderer{
		cfg: cfg,
	}
}

// HomeHandler serves the main Jason front page
func (r *Renderer) HomeHandler(w http.ResponseWriter, req *http.Request) {
	// Load or extract agent profile data
	if err := r.loadProfile(); err != nil {
		log.Printf("Failed to load profile: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	tmplPath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, r.profile); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Render Error", http.StatusInternalServerError)
	}
}

// HealthHandler provides health check endpoint
func (r *Renderer) HealthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "jason-frontpage"}`))
}

// loadProfile loads or extracts the agent profile with privacy filtering
func (r *Renderer) loadProfile() error {
	// Check if we have cached profile
	cachePath := filepath.Join(r.cfg.WorkspacePath, "cache", "agent-profile.json")
	if _, err := os.Stat(cachePath); err == nil {
		// Load from cache
		profile, err := models.LoadProfile(cachePath)
		if err != nil {
			return err
		}
		r.profile = profile
		return nil
	}

	// Extract fresh profile (with privacy filtering)
	profile, err := r.extractProfile()
	if err != nil {
		return err
	}
	r.profile = profile
	return nil
}

// extractProfile extracts agent profile data from workspace with privacy filtering
func (r *Renderer) extractProfile() (*models.AgentProfile, error) {
	// This will be implemented in Phase 2
	// For now, return a basic profile structure
	return &models.AgentProfile{
		Name:        "Jason",
		Title:       "Digital Agent",
		Description: "Your AI companion and digital presence",
		Mission:     "To assist, collaborate, and showcase responsible AI interaction",
		WorkAreas:   []string{"AI Assistance", "Digital Presence", "Collaboration"},
	}, nil
}