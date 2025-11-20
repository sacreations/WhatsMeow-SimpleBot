package system

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

type InfoCommand struct{}

func NewInfoCommand() *InfoCommand         { return &InfoCommand{} }
func (i *InfoCommand) Name() string        { return "/info" }
func (i *InfoCommand) Description() string { return "Get your chat info" }

func (i *InfoCommand) Execute(args []string, sender types.JID) string {
	return fmt.Sprintf("👤 *Your Chat Info:*\n\nJID: %s\nUser: %s", sender, strings.Split(sender.String(), "@")[0])
}
