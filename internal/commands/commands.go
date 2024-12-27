package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/krzko/lotek-bot/internal/interfaces"
	"github.com/krzko/lotek-bot/pkg/jokes"
)

// Command represents a bot command
type Command interface {
	Name() string
	Aliases() []string
	Description() string
	SlashCommand() *discordgo.ApplicationCommand
	HandleText(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
	HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

// Handler manages bot commands
type Handler struct {
	prefix   string
	jokes    *jokes.Service
	commands map[string]Command
	bot      interfaces.Bot
}

// NewHandler creates a new command handler
func NewHandler(prefix string, jokes *jokes.Service, bot interfaces.Bot) *Handler {
	h := &Handler{
		prefix:   prefix,
		jokes:    jokes,
		commands: make(map[string]Command),
		bot:      bot,
	}

	// Register commands
	h.registerCommands()
	return h
}

// RegisterSlashCommands registers all slash commands with Discord
func (h *Handler) RegisterSlashCommands(s *discordgo.Session) error {
	for _, cmd := range h.commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd.SlashCommand())
		if err != nil {
			return fmt.Errorf("failed to create slash command %s: %w", cmd.Name(), err)
		}
	}
	return nil
}

// HandleInteraction processes slash command interactions
func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmdName := i.ApplicationCommandData().Name
	if cmd, exists := h.commands[cmdName]; exists {
		if err := cmd.HandleSlash(s, i); err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error executing command: %v", err),
				},
			})
		}
	}
}

// HandleMessage processes text commands
func (h *Handler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !strings.HasPrefix(m.Content, h.prefix) {
		return nil
	}

	parts := strings.Fields(m.Content)
	if len(parts) == 0 {
		return nil
	}

	cmdName := strings.TrimPrefix(parts[0], h.prefix)
	// Check for both the command name and its aliases
	for name, cmd := range h.commands {
		if cmdName == name || contains(cmd.Aliases(), cmdName) {
			return cmd.HandleText(s, m, parts[1:])
		}
	}

	return fmt.Errorf("unknown command: %s", cmdName)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// registerCommands registers all available commands
func (h *Handler) registerCommands() {
	h.commands["joke"] = &JokeCommand{jokes: h.jokes}
	h.commands["help"] = &HelpCommand{commands: h.commands}
	h.commands["aliases"] = &AliasesCommand{bot: h.bot}
}

// JokeCommand handles the joke command
type JokeCommand struct {
	jokes *jokes.Service
}

func (c *JokeCommand) Name() string {
	return "joke"
}

func (c *JokeCommand) Aliases() []string {
	return []string{"jokes"}
}

func (c *JokeCommand) Description() string {
	return "Tells a random nerdy joke"
}

func (c *JokeCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
	}
}

func (c *JokeCommand) HandleText(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	joke := c.jokes.GetRandomJoke()
	_, err := s.ChannelMessageSend(m.ChannelID, joke)
	return err
}

func (c *JokeCommand) HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	joke := c.jokes.GetRandomJoke()
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: joke,
		},
	})
}

// HelpCommand handles the help command
type HelpCommand struct {
	commands map[string]Command
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Aliases() []string {
	return []string{}
}

func (c *HelpCommand) Description() string {
	return "Shows available commands and their descriptions"
}

func (c *HelpCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
	}
}

func (c *HelpCommand) HandleText(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	helpMsg := c.buildHelpMessage()
	_, err := s.ChannelMessageSend(m.ChannelID, helpMsg)
	return err
}

func (c *HelpCommand) HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	helpMsg := c.buildHelpMessage()
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMsg,
		},
	})
}

func (c *HelpCommand) buildHelpMessage() string {
	var helpMsg strings.Builder
	helpMsg.WriteString("Available commands:\n")

	for _, cmd := range c.commands {
		aliases := ""
		if len(cmd.Aliases()) > 0 {
			aliases = fmt.Sprintf(" (aliases: %s)", strings.Join(cmd.Aliases(), ", "))
		}
		helpMsg.WriteString(fmt.Sprintf("- /%s%s: %s\n", cmd.Name(), aliases, cmd.Description()))
	}

	return helpMsg.String()
}

// AliasesCommand shows all monitored users and their aliases
type AliasesCommand struct {
	bot interfaces.Bot
}

func (c *AliasesCommand) Name() string {
	return "aliases"
}

func (c *AliasesCommand) Aliases() []string {
	return []string{"users", "monitored"}
}

func (c *AliasesCommand) Description() string {
	return "Shows all monitored users and their aliases"
}

func (c *AliasesCommand) SlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Description: c.Description(),
	}
}

func (c *AliasesCommand) HandleText(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	aliases := c.bot.GetUserAliases()
	var msg strings.Builder
	msg.WriteString("Monitored users and their aliases:\n")

	for name, userAliases := range aliases {
		msg.WriteString(fmt.Sprintf("- %s: %s\n", name, strings.Join(userAliases, ", ")))
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg.String())
	return err
}

func (c *AliasesCommand) HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	aliases := c.bot.GetUserAliases()
	var msg strings.Builder
	msg.WriteString("Monitored users and their aliases:\n")

	for name, userAliases := range aliases {
		msg.WriteString(fmt.Sprintf("- %s: %s\n", name, strings.Join(userAliases, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg.String(),
		},
	})
}
