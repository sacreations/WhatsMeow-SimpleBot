package handlers

import (
	"fmt"

	"whatsappBotGo/src/senders"

	"go.mau.fi/whatsmeow/types"
)

// AdminHandler handles admin-level operations
type AdminHandler struct {
	senders *senders.Senders
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(s *senders.Senders) *AdminHandler {
	return &AdminHandler{senders: s}
}

// SendSystemMessage sends a system-level message
func (a *AdminHandler) SendSystemMessage(to types.JID, msg string) error {
	if a.senders == nil || a.senders.Text == nil {
		return fmt.Errorf("text sender not configured")
	}
	return a.senders.Text.SendText(to, fmt.Sprintf("🔧 *SYSTEM*: %s", msg))
}

// SendErrorAlert sends an error alert message
func (a *AdminHandler) SendErrorAlert(to types.JID, err string) error {
	if a.senders == nil || a.senders.Text == nil {
		return fmt.Errorf("text sender not configured")
	}
	return a.senders.Text.SendText(to, fmt.Sprintf("⚠️ *ERROR*: %s", err))
}

// BroadcastMessage broadcasts a message (placeholder for multi-recipient support)
func (a *AdminHandler) BroadcastMessage(recipients []types.JID, msg string) error {
	if a.senders == nil || a.senders.Text == nil {
		return fmt.Errorf("text sender not configured")
	}
	for _, to := range recipients {
		if err := a.senders.Text.SendText(to, msg); err != nil {
			return fmt.Errorf("failed to broadcast to %s: %v", to, err)
		}
	}
	return nil
}
