package bot

import (
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// Command represents a bot command
type Command interface {
	Name() string
	Description() string
	Execute(args []string, sender types.JID) string
}

// CommandHandler manages all bot commands
type CommandHandler struct {
	commands         map[string]Command
	autoReplyHandler *AutoReplyHandler
}

// NewCommandHandler creates a new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands:         make(map[string]Command),
		autoReplyHandler: NewAutoReplyHandler(),
	}
}

// SetClient sets the WhatsApp client for the auto reply handler
func (ch *CommandHandler) SetClient(client *whatsmeow.Client) {
	ch.autoReplyHandler.SetClient(client)
}

// RegisterCommand registers a new command
func (ch *CommandHandler) RegisterCommand(cmd Command) {
	ch.commands[cmd.Name()] = cmd
}

// handleNonCommand handles messages that are not commands
func (ch *CommandHandler) handleNonCommand(message string, sender types.JID) string {
	return ch.autoReplyHandler.ProcessMessage(message, sender)
}

// ProcessMessage processes incoming messages and executes commands
func (ch *CommandHandler) ProcessMessage(message string, sender types.JID) string {
	originalMessage := strings.TrimSpace(message)
	message = strings.ToLower(originalMessage)

	// Handle non-command messages
	if !strings.HasPrefix(message, "/") {
		return ch.handleNonCommand(originalMessage, sender)
	}

	// Parse command and arguments
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "🤔 I didn't understand that. Type /help for available commands."
	}

	commandName := parts[0]
	args := parts[1:]

	// Find and execute command
	if cmd, exists := ch.commands[commandName]; exists {
		return cmd.Execute(args, sender)
	}
	return "" //return nothing if command not found
}

// GetAllCommands returns all registered commands
func (ch *CommandHandler) GetAllCommands() map[string]Command {
	return ch.commands
}
