package commands

import (
	"math/rand"
	"time"

	"go.mau.fi/whatsmeow/types"
)

type QuoteCommand struct {
	quotes []string
}

func NewQuoteCommand() *QuoteCommand {
	quotes := []string{
		"\"The only way to do great work is to love what you do.\" - Steve Jobs ✨",
		"\"Innovation distinguishes between a leader and a follower.\" - Steve Jobs 🚀",
		"\"Life is what happens to you while you're busy making other plans.\" - John Lennon 🌟",
		"\"The future belongs to those who believe in the beauty of their dreams.\" - Eleanor Roosevelt 💭",
		"\"It is during our darkest moments that we must focus to see the light.\" - Aristotle 💡",
		"\"Success is not final, failure is not fatal: it is the courage to continue that counts.\" - Winston Churchill 💪",
		"\"The only impossible journey is the one you never begin.\" - Tony Robbins 🛤️",
	}

	return &QuoteCommand{quotes: quotes}
}

func (q *QuoteCommand) Name() string {
	return "/quote"
}

func (q *QuoteCommand) Description() string {
	return "Get an inspirational quote"
}

func (q *QuoteCommand) Execute(args []string, sender types.JID) string {
	rand.Seed(time.Now().UnixNano())
	randomQuote := q.quotes[rand.Intn(len(q.quotes))]
	return "💫 " + randomQuote
}
