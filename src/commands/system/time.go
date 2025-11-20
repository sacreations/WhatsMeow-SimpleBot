package system

import (
	"time"

	"go.mau.fi/whatsmeow/types"
)

type TimeCommand struct{}

func NewTimeCommand() *TimeCommand         { return &TimeCommand{} }
func (t *TimeCommand) Name() string        { return "/time" }
func (t *TimeCommand) Description() string { return "Get current time" }

func (t *TimeCommand) Execute(args []string, sender types.JID) string {
	return time.Now().Format(time.RFC1123)
}
