package models

import (
	"encoding/json"
	"os"
)

// AgentProfile represents the public profile of the Jason agent
type AgentProfile struct {
	// Identity information (public)
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Mission     string   `json:"mission"`
	WorkAreas   []string `json:"work_areas"`

	// Collaboration information (public)
	Collaborators []string `json:"collaborators,omitempty"`

	// Metadata
	LastUpdated string `json:"last_updated,omitempty"`
	Version     string `json:"version,omitempty"`
}

// LoadProfile loads a profile from JSON file
func LoadProfile(path string) (*AgentProfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var profile AgentProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// SaveProfile saves a profile to JSON file
func (p *AgentProfile) SaveProfile(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}