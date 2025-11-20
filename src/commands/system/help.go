package system

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type HelpCommand struct{}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (h *HelpCommand) Name() string        { return "/help" }
func (h *HelpCommand) Description() string { return "Show this help message" }

func (h *HelpCommand) Execute(args []string, sender types.JID) string {
	var response strings.Builder
	response.WriteString("🤖 *WhatsApp Bot Commands:*\n\n")

	// Predefined commands help text
	commands := map[string]string{
		"/help":  "Show this help message",
		"/ping":  "Check if bot is alive",
		"/time":  "Get current time",
		"/echo":  "Echo your message",
		"/info":  "Get your chat info",
		"/joke":  "Get a random joke",
		"/quote": "Get an inspirational quote",
	}

	for cmd, desc := range commands {
		response.WriteString(fmt.Sprintf("%s - %s\n", cmd, desc))
	}

	response.WriteString("\nSend any message to interact with the bot!")
	return response.String()
}

func (h *HelpCommand) ExecuteWithContext(args []string, evt *events.Message, sender types.JID) string {
	return h.Execute(args, sender)
}
