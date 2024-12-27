package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the bot configuration
type Config struct {
	Token            string
	StartupChannelID string
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Don't return error if .env file doesn't exist
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is required")
	}

	startupChannelID := os.Getenv("DISCORD_STARTUP_CHANNEL_ID")
	if startupChannelID == "" {
		return nil, fmt.Errorf("DISCORD_STARTUP_CHANNEL_ID environment variable is required")
	}

	return &Config{
		Token:            token,
		StartupChannelID: startupChannelID,
	}, nil
}
