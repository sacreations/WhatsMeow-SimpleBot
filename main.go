package main

import (
	"log"
	"whatsappBotGo/src/bot"
)

func main() {
	// Create new bot instance
	whatsappBot, err := bot.NewBot()
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start the bot
	if err := whatsappBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
