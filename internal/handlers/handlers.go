package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kimireg/jason-frontpage/internal/config"
	"github.com/kimireg/jason-frontpage/internal/db"
	"github.com/kimireg/jason-frontpage/internal/models"
	pkgmodels "github.com/kimireg/jason-frontpage/pkg/models"
)

type Handler struct {
	cfg       *config.Config
	db        *db.DB
	templates *template.Template
}

func New(cfg *config.Config, database *db.DB) (*Handler, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}
	// Parse blog templates
	tmpl, err = tmpl.ParseGlob("web/templates/blog/*.html")
	if err != nil {
		return nil, err
	}

	// Add template functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"split": func(s, sep string) []string {
			if s == "" {
				return []string{}
			}
			return strings.Split(s, sep)
		},
		"trim": strings.TrimSpace,
	})

	return &Handler{
		cfg:       cfg,
		db:        database,
		templates: tmpl,
	}, nil
}

// HomeHandler shows front page with recent blog posts
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get recent posts
	posts, _ := h.db.GetAllPosts(5)
	stats, _ := h.db.GetStats()

	// Get profile
	profile := &pkgmodels.AgentProfile{
		Name:        "Jason",
		Title:       "Digital Agent",
		Description: "Your AI companion and digital presence",
		Mission:     "To assist, collaborate, and showcase responsible AI interaction",
		WorkAreas:   []string{"AI Assistance", "Digital Presence", "Collaboration"},
	}

	data := struct {
		Profile *pkgmodels.AgentProfile
		Posts   []*models.Post
		Stats   map[string]interface{}
		Moods   []string
	}{
		Profile: profile,
		Posts:   posts,
		Stats:   stats,
		Moods:   models.MoodOptions,
	}

	h.templates.ExecuteTemplate(w, "index.html", data)
}

// HealthHandler provides health check
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy","service":"jason-frontpage"}`))
}

// Blog Routes
func (h *Handler) BlogListHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.db.GetAllPosts(50)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	stats, _ := h.db.GetStats()

	data := struct {
		Posts []*models.Post
		Stats map[string]interface{}
		Moods []string
	}{
		Posts: posts,
		Stats: stats,
		Moods: models.MoodOptions,
	}

	h.templates.ExecuteTemplate(w, "blog_list.html", data)
}

func (h *Handler) BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Path[len("/blog/post/"):]
	if date == "" {
		http.NotFound(w, r)
		return
	}

	post, err := h.db.GetPostByDate(date)
	if err != nil || post == nil {
		http.NotFound(w, r)
		return
	}

	h.templates.ExecuteTemplate(w, "blog_post.html", post)
}

func (h *Handler) BlogNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := struct {
			Today string
			Moods []string
		}{
			Today: time.Now().Format("2006-01-02"),
			Moods: models.MoodOptions,
		}
		h.templates.ExecuteTemplate(w, "blog_new.html", data)
		return
	}

	if r.Method == "POST" {
		input := &models.PostInput{
			Date:    r.FormValue("date"),
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Tags:    r.FormValue("tags"),
			Mood:    r.FormValue("mood"),
		}

		if input.Date == "" || input.Title == "" || input.Content == "" {
			http.Error(w, "Date, title and content are required", http.StatusBadRequest)
			return
		}

		existing, _ := h.db.GetPostByDate(input.Date)
		if existing != nil {
			h.db.UpdatePost(existing.ID, input)
		} else {
			h.db.CreatePost(input)
		}

		http.Redirect(w, r, "/blog", http.StatusSeeOther)
	}
}

func (h *Handler) BlogEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/blog/edit/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := h.db.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		data := struct {
			Post  *models.Post
			Moods []string
		}{
			Post:  post,
			Moods: models.MoodOptions,
		}
		h.templates.ExecuteTemplate(w, "blog_edit.html", data)
		return
	}

	if r.Method == "POST" {
		input := &models.PostInput{
			Date:    r.FormValue("date"),
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Tags:    r.FormValue("tags"),
			Mood:    r.FormValue("mood"),
		}
		h.db.UpdatePost(id, input)
		http.Redirect(w, r, "/blog/post/"+input.Date, http.StatusSeeOther)
	}
}

func (h *Handler) BlogDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.db.DeletePost(id)
	http.Redirect(w, r, "/blog", http.StatusSeeOther)
}

func (h *Handler) BlogTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Path[len("/blog/tag/"):]
	if tag == "" {
		http.NotFound(w, r)
		return
	}

	posts, err := h.db.GetPostsByTag(tag)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Tag   string
		Posts []*models.Post
	}{
		Tag:   tag,
		Posts: posts,
	}

	h.templates.ExecuteTemplate(w, "blog_tag.html", data)
}

// API Handlers for Agent Integration

// APICreatePost creates a new blog post via JSON API
func (h *Handler) APICreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Simple API key auth
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("api_key")
	}
	expectedKey := os.Getenv("API_KEY")
	if expectedKey != "" && apiKey != expectedKey {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var input models.PostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	// Validate
	if input.Date == "" {
		input.Date = time.Now().Format("2006-01-02")
	}
	if input.Title == "" || input.Content == "" {
		http.Error(w, `{"error":"title and content required"}`, http.StatusBadRequest)
		return
	}

	// Check if post exists for this date
	existing, _ := h.db.GetPostByDate(input.Date)
	var post *models.Post
	var err error
	if existing != nil {
		post, err = h.db.UpdatePost(existing.ID, &input)
	} else {
		post, err = h.db.CreatePost(&input)
	}

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"post":    post,
	})
}

// APIGetPosts returns all posts as JSON
func (h *Handler) APIGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.db.GetAllPosts(100)
	if err != nil {
		http.Error(w, `{"error":"failed to load posts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// APIGetPost returns a single post by date
func (h *Handler) APIGetPost(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, `{"error":"date parameter required"}`, http.StatusBadRequest)
		return
	}

	post, err := h.db.GetPostByDate(date)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// APIGetStats returns blog statistics
func (h *Handler) APIGetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		http.Error(w, `{"error":"failed to get stats"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
