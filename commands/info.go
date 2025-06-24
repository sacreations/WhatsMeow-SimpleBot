package commands

import (
    "fmt"
    "time"
    "go.mau.fi/whatsmeow/types"
)

type InfoCommand struct{}

func NewInfoCommand() *InfoCommand {
    return &InfoCommand{}
}

func (i *InfoCommand) Name() string {
    return "/info"
}

func (i *InfoCommand) Description() string {
    return "Get your chat info"
}

func (i *InfoCommand) Execute(args []string, sender types.JID) string {
    return fmt.Sprintf(`ℹ️ *Chat Information:*
Your JID: %s
Chat Type: %s
Timestamp: %s`, 
        sender.String(), 
        sender.Server, 
        time.Now().Format("2006-01-02 15:04:05"))
}
