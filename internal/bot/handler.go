package bot

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleReady logs when the bot is ready and registers slash commands
func (b *Bot) handleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

	// Register slash commands
	if err := b.cmdHandler.RegisterSlashCommands(s); err != nil {
		log.Printf("Error registering slash commands: %v", err)
	} else {
		log.Println("Successfully registered slash commands")
	}

	// Send startup message
	_, err := s.ChannelMessageSend(b.config.StartupChannelID, "🚀 Bot is now online and ready to annoy Tom!")
	if err != nil {
		log.Printf("Error sending startup message: %v", err)
	}

	// Check for users that are already online
	time.Sleep(2 * time.Second) // Give Discord a moment to populate presence data

	guilds := s.State.Guilds
	for _, guild := range guilds {
		// Use the guild's Presences directly
		for _, presence := range guild.Presences {
			if presence.Status == discordgo.StatusOnline || presence.Status == discordgo.StatusIdle {
				// Check if this is a user we're monitoring
				if err := b.presenceHandler.HandlePresenceUpdate(s, &discordgo.PresenceUpdate{
					Presence: *presence,
					GuildID:  guild.ID,
				}); err != nil {
					log.Printf("Error handling initial presence: %v", err)
				}
			}
		}
	}
}

// handlePresenceUpdate handles user presence updates
func (b *Bot) handlePresenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	// Log the presence update for debugging
	log.Printf("Presence update received - User: %v, Status: %v", p.User.Username, p.Status)

	if err := b.presenceHandler.HandlePresenceUpdate(s, p); err != nil {
		log.Printf("Error handling presence update: %v", err)
	}
}

// handleMessageCreate handles incoming messages
func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	log.Printf("Received message: %s", m.Content)

	// Handle commands through the command handler
	if err := b.cmdHandler.HandleMessage(s, m); err != nil {
		log.Printf("Error handling command: %v", err)
	}
}

// handleInteraction handles slash command interactions
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.cmdHandler.HandleInteraction(s, i)
}
