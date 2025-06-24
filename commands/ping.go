package commands

import (
    "go.mau.fi/whatsmeow/types"
)

type PingCommand struct{}

func NewPingCommand() *PingCommand {
    return &PingCommand{}
}

func (p *PingCommand) Name() string {
    return "/ping"
}

func (p *PingCommand) Description() string {
    return "Check if bot is alive"
}

func (p *PingCommand) Execute(args []string, sender types.JID) string {
    return "🏓 Pong! Bot is alive and running."
}
