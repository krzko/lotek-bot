package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// UserConfig represents the configuration for a user
type UserConfig struct {
	Aliases []string
	Message string
}

// PresenceHandler manages user presence monitoring
type PresenceHandler struct {
	users     map[string]UserConfig // Map of Discord User IDs to their configs
	channelID string                // Channel ID to send messages to
}

// NewPresenceHandler creates a new presence handler
func NewPresenceHandler(channelID string) *PresenceHandler {
	handler := &PresenceHandler{
		users:     make(map[string]UserConfig),
		channelID: channelID,
	}

	// Add default users and their configs
	handler.AddUser("Kristof", []string{
		"Kristof",
		"krzko",
		"Tommy",
		"TommyBoy",
	}, "Shuddup %s!")

	log.Printf("Initialized presence handler with monitored users: %v", handler.GetUserAliases())
	log.Printf("Will send notifications to channel: %s", channelID)

	return handler
}

// AddUser adds a new user to monitor
func (h *PresenceHandler) AddUser(name string, aliases []string, message string) {
	// Convert all aliases to lowercase for case-insensitive matching
	lowerAliases := make([]string, len(aliases))
	for i, alias := range aliases {
		lowerAliases[i] = strings.ToLower(alias)
	}

	h.users[name] = UserConfig{
		Aliases: lowerAliases,
		Message: message,
	}
}

// HandlePresenceUpdate processes a presence update event
func (h *PresenceHandler) HandlePresenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) error {
	if p.User == nil {
		return fmt.Errorf("presence update has no user information")
	}

	// Log the presence update for debugging
	log.Printf("Processing presence update for user: %s (ID: %s), Status: %s",
		p.User.Username, p.User.ID, p.Status)

	// Check if the user matches any of our monitored users
	username := strings.ToLower(p.User.Username)
	for name, config := range h.users {
		log.Printf("Checking against monitored user %s with aliases: %v", name, config.Aliases)
		if containsUsername(config.Aliases, username) {
			// Check if the user is online or their status turned green
			if isUserActive(p) {
				message := fmt.Sprintf(config.Message, p.User.Username)
				log.Printf("User %s is active, sending message to channel %s: %s",
					p.User.Username, h.channelID, message)
				_, err := s.ChannelMessageSend(h.channelID, message)
				if err != nil {
					return fmt.Errorf("failed to send message: %w", err)
				}
			}
			break
		}
	}
	return nil
}

// containsUsername checks if a username matches any of the aliases
func containsUsername(aliases []string, username string) bool {
	for _, alias := range aliases {
		if strings.ToLower(alias) == username {
			return true
		}
	}
	return false
}

// isUserActive checks if a user is active (online or green status)
func isUserActive(p *discordgo.PresenceUpdate) bool {
	return p.Status == discordgo.StatusOnline || p.Status == discordgo.StatusIdle
}

// GetUserAliases returns all configured aliases for monitoring
func (h *PresenceHandler) GetUserAliases() map[string][]string {
	result := make(map[string][]string)
	for name, config := range h.users {
		result[name] = config.Aliases
	}
	return result
}
