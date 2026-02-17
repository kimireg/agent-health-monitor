package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimireg/jason-frontpage/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return database, nil
}

func (db *DB) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		tags TEXT,
		mood TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_posts_date ON posts(date DESC);
	`
	_, err := db.Exec(query)
	return err
}

func (db *DB) CreatePost(input *models.PostInput) (*models.Post, error) {
	result, err := db.Exec(
		"INSERT INTO posts (date, title, content, tags, mood) VALUES (?, ?, ?, ?, ?)",
		input.Date, input.Title, input.Content, input.Tags, input.Mood,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	id, _ := result.LastInsertId()
	return db.GetPostByID(int(id))
}

func (db *DB) GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	err := db.QueryRow(
		"SELECT id, date, title, content, tags, mood, created_at, updated_at FROM posts WHERE id = ?",
		id,
	).Scan(&post.ID, &post.Date, &post.Title, &post.Content, &post.Tags, &post.Mood, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (db *DB) GetPostByDate(date string) (*models.Post, error) {
	var post models.Post
	err := db.QueryRow(
		"SELECT id, date, title, content, tags, mood, created_at, updated_at FROM posts WHERE date = ?",
		date,
	).Scan(&post.ID, &post.Date, &post.Title, &post.Content, &post.Tags, &post.Mood, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (db *DB) GetAllPosts(limit int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := db.Query(
		"SELECT id, date, title, content, tags, mood, created_at, updated_at FROM posts ORDER BY date DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Date, &p.Title, &p.Content, &p.Tags, &p.Mood, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, rows.Err()
}

func (db *DB) GetPostsByTag(tag string) ([]*models.Post, error) {
	rows, err := db.Query(
		"SELECT id, date, title, content, tags, mood, created_at, updated_at FROM posts WHERE tags LIKE ? ORDER BY date DESC",
		"%"+tag+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Date, &p.Title, &p.Content, &p.Tags, &p.Mood, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, rows.Err()
}

func (db *DB) UpdatePost(id int, input *models.PostInput) (*models.Post, error) {
	_, err := db.Exec(
		"UPDATE posts SET date = ?, title = ?, content = ?, tags = ?, mood = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		input.Date, input.Title, input.Content, input.Tags, input.Mood, id,
	)
	if err != nil {
		return nil, err
	}
	return db.GetPostByID(id)
}

func (db *DB) DeletePost(id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

func (db *DB) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&total); err != nil {
		return nil, err
	}
	stats["total_posts"] = total

	var firstDate, latestDate string
	db.QueryRow("SELECT MIN(date) FROM posts").Scan(&firstDate)
	db.QueryRow("SELECT MAX(date) FROM posts").Scan(&latestDate)
	stats["first_post"] = firstDate
	stats["latest_post"] = latestDate

	return stats, nil
}
