package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/krzko/lotek-bot/internal/commands"
	"github.com/krzko/lotek-bot/internal/config"
	"github.com/krzko/lotek-bot/internal/interfaces"
	"github.com/krzko/lotek-bot/pkg/jokes"
)

// Bot represents the Discord bot
type Bot struct {
	session         *discordgo.Session
	config          *config.Config
	jokes           *jokes.Service
	cmdHandler      *commands.Handler
	presenceHandler *PresenceHandler
}

// GetUserAliases returns all configured user aliases
func (b *Bot) GetUserAliases() map[string][]string {
	return b.presenceHandler.GetUserAliases()
}

// Bot should implement the interfaces.Bot interface
var _ interfaces.Bot = (*Bot)(nil)

// New creates a new instance of the bot
func New(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	// Add required intents
	session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildPresences

	bot := &Bot{
		session: session,
		config:  cfg,
	}

	// Initialize handlers with the startup channel ID
	bot.presenceHandler = NewPresenceHandler(cfg.StartupChannelID)
	bot.jokes = jokes.NewService()
	bot.cmdHandler = commands.NewHandler("!", bot.jokes, bot)

	// Register handlers
	session.AddHandler(bot.handleReady)
	session.AddHandler(bot.handlePresenceUpdate)
	session.AddHandler(bot.handleMessageCreate)
	session.AddHandler(bot.handleInteraction)

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start(ctx context.Context) error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	log.Println("Bot is now running")
	return nil
}

// Stop gracefully shuts down the bot
func (b *Bot) Stop(ctx context.Context) error {
	return b.session.Close()
}
