package system

import (
	"time"

	"go.mau.fi/whatsmeow/types"
)

type PingCommand struct{}

func NewPingCommand() *PingCommand         { return &PingCommand{} }
func (p *PingCommand) Name() string        { return "/ping" }
func (p *PingCommand) Description() string { return "Check if bot is alive" }

func (p *PingCommand) Execute(args []string, sender types.JID) string {
	return "PONG - " + time.Now().Format(time.RFC3339)
}
