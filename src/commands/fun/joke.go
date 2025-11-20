package fun

import (
	"math/rand"

	"go.mau.fi/whatsmeow/types"
)

type JokeCommand struct{}

func NewJokeCommand() *JokeCommand         { return &JokeCommand{} }
func (j *JokeCommand) Name() string        { return "/joke" }
func (j *JokeCommand) Description() string { return "Get a random joke" }

func (j *JokeCommand) Execute(args []string, sender types.JID) string {
	jokes := []string{
		"😄 Why did the scarecrow win an award? Because he was outstanding in his field!",
		"😂 What do you call a fake noodle? An impasta!",
		"😆 Why don't scientists trust atoms? Because they make up everything!",
		"🤣 What did the ocean say to the beach? Nothing, it just waved!",
		"😅 Why don't eggs tell jokes? They'd crack each other up!",
	}
	return jokes[rand.Intn(len(jokes))]
}
