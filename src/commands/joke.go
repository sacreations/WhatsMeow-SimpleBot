package commands

import (
	"math/rand"
	"time"

	"go.mau.fi/whatsmeow/types"
)

type JokeCommand struct {
	jokes []string
}

func NewJokeCommand() *JokeCommand {
	jokes := []string{
		"Why don't scientists trust atoms? Because they make up everything! 😄",
		"Why did the scarecrow win an award? He was outstanding in his field! 🌾",
		"Why don't eggs tell jokes? They'd crack each other up! 🥚",
		"What do you call a fake noodle? An impasta! 🍝",
		"Why did the math book look so sad? Because it had too many problems! 📚",
		"What do you call a bear with no teeth? A gummy bear! 🐻",
		"Why can't a bicycle stand up by itself? It's two tired! 🚲",
	}

	return &JokeCommand{jokes: jokes}
}

func (j *JokeCommand) Name() string {
	return "/joke"
}

func (j *JokeCommand) Description() string {
	return "Get a random joke"
}

func (j *JokeCommand) Execute(args []string, sender types.JID) string {
	rand.Seed(time.Now().UnixNano())
	randomJoke := j.jokes[rand.Intn(len(j.jokes))]
	return "😂 " + randomJoke
}
