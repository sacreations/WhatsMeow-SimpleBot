package system

import (
	"strings"

	"go.mau.fi/whatsmeow/types"
)

type EchoCommand struct{}

func NewEchoCommand() *EchoCommand         { return &EchoCommand{} }
func (e *EchoCommand) Name() string        { return "/echo" }
func (e *EchoCommand) Description() string { return "Echo your message" }

func (e *EchoCommand) Execute(args []string, sender types.JID) string {
	return "📢 " + strings.Join(args, " ")
}
