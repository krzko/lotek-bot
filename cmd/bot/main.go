package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/krzko/lotek-bot/internal/bot"
	"github.com/krzko/lotek-bot/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	discordBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := discordBot.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	if err := discordBot.Stop(ctx); err != nil {
		log.Printf("Error stopping bot: %v", err)
	}
}
