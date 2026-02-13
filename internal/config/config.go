package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	ServerAddr string `mapstructure:"server_addr"`

	// Workspace paths
	WorkspacePath string `mapstructure:"workspace_path"`

	// Update schedule
	UpdateInterval time.Duration `mapstructure:"update_interval"`

	// Privacy settings
	PrivacyEnabled bool `mapstructure:"privacy_enabled"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/jason-frontpage/")
	viper.AddConfigPath("$HOME/.jason-frontpage/")

	// Set defaults
	viper.SetDefault("server_addr", ":8080")
	viper.SetDefault("workspace_path", "/workspace")
	viper.SetDefault("update_interval", "168h") // 1 week
	viper.SetDefault("privacy_enabled", true)

	// Read config file (if exists)
	_ = viper.ReadInConfig()

	// Override with environment variables
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Validate required fields
	if cfg.WorkspacePath == "" {
		cfg.WorkspacePath = os.Getenv("WORKSPACE_PATH")
		if cfg.WorkspacePath == "" {
			cfg.WorkspacePath = "/workspace"
		}
	}

	return &cfg, nil
}