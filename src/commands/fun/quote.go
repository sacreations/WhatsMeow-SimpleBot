package fun

import (
	"math/rand"

	"go.mau.fi/whatsmeow/types"
)

type QuoteCommand struct{}

func NewQuoteCommand() *QuoteCommand        { return &QuoteCommand{} }
func (q *QuoteCommand) Name() string        { return "/quote" }
func (q *QuoteCommand) Description() string { return "Get an inspirational quote" }

func (q *QuoteCommand) Execute(args []string, sender types.JID) string {
	quotes := []string{
		"✨ \"The best time to plant a tree was 20 years ago. The second best time is now.\" - Chinese Proverb",
		"🌟 \"Your time is limited, don't waste it living someone else's life.\" - Steve Jobs",
		"💪 \"The only way to do great work is to love what you do.\" - Steve Jobs",
		"🎯 \"Success is not final, failure is not fatal.\" - Winston Churchill",
		"🚀 \"Believe you can and you're halfway there.\" - Theodore Roosevelt",
	}
	return quotes[rand.Intn(len(quotes))]
}
