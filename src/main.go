package main

import (
	"context"
	"log"
	"time"

	"whatsappBotGo/src/api"
	"whatsappBotGo/src/bot"
)

func main() {
	// Create new bot instance
	whatsappBot, err := bot.NewBot()
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start API server (if senders provided) and wire WS hub
	if whatsappBot != nil && whatsappBot.Senders() != nil {
		srv := api.NewServer(whatsappBot.Senders(), whatsappBot.InstanceUserID())
		// Attach to bot so it can broadcast incoming messages
		whatsappBot.SetAPIServer(srv)
		srv.Start()
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}()
	}

	// Start the bot (blocking)
	if err := whatsappBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
