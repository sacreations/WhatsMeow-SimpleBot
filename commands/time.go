package commands

import (
    "fmt"
    "time"
    "go.mau.fi/whatsmeow/types"
)

type TimeCommand struct{}

func NewTimeCommand() *TimeCommand {
    return &TimeCommand{}
}

func (t *TimeCommand) Name() string {
    return "/time"
}

func (t *TimeCommand) Description() string {
    return "Get current time"
}

func (t *TimeCommand) Execute(args []string, sender types.JID) string {
    now := time.Now()
    return fmt.Sprintf("🕐 Current time: %s", now.Format("2006-01-02 15:04:05 MST"))
}
