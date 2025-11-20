package admin

import (
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// StatusCommand returns bot status for admins
type StatusCommand struct{}

func NewStatusCommand() *StatusCommand       { return &StatusCommand{} }
func (s *StatusCommand) Name() string        { return "/status" }
func (s *StatusCommand) Description() string { return "Get bot status (admin)" }

func (s *StatusCommand) Execute(args []string, sender types.JID) string {
	return "🟢 *Bot Status*\nUptime: Running\nConnected: Yes\nVersion: 1.0.0"
}

func (s *StatusCommand) ExecuteWithContext(args []string, evt *events.Message, sender types.JID) string {
	return s.Execute(args, sender)
}
