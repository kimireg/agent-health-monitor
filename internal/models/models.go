package models

import "time"

// Post represents a blog entry written by Jason
type Post struct {
	ID        int       `json:"id" db:"id"`
	Date      string    `json:"date" db:"date"`           // YYYY-MM-DD
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`     // Markdown supported
	Tags      string    `json:"tags" db:"tags"`           // Comma-separated
	Mood      string    `json:"mood" db:"mood"`           // happy, excited, tired, curious, etc.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PostInput for creating/updating posts
type PostInput struct {
	Date    string `json:"date" form:"date"`
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Tags    string `json:"tags" form:"tags"`
	Mood    string `json:"mood" form:"mood"`
}

// MoodOptions available for selection
var MoodOptions = []string{
	"excited", "happy", "focused", "curious", "contemplative",
	"tired", "challenged", "accomplished",
}

// AgentProfile for front page display
type AgentProfile struct {
	Name          string   `json:"name"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Mission       string   `json:"mission"`
	WorkAreas     []string `json:"work_areas"`
	Collaborators []string `json:"collaborators,omitempty"`
	LastUpdated   string   `json:"last_updated,omitempty"`
	Version       string   `json:"version,omitempty"`
}
