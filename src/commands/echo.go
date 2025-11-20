package commands

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

type EchoCommand struct{}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (e *EchoCommand) Name() string {
	return "/echo"
}

func (e *EchoCommand) Description() string {
	return "Echo your message"
}

func (e *EchoCommand) Execute(args []string, sender types.JID) string {
	if len(args) == 0 {
		return "📢 Please provide text to echo. Usage: /echo <your message>"
	}

	text := strings.Join(args, " ")
	return fmt.Sprintf("📢 Echo: %s", text)
}
