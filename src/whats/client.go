package whats

import (
	"go.mau.fi/whatsmeow"
)

// Client wraps whatsmeow.Client so we can move client-specific helper
// functions into this package later without leaking whatsmeow types.
type Client struct {
	Client *whatsmeow.Client
}

// NewFromClient wraps an existing whatsmeow client.
func NewFromClient(c *whatsmeow.Client) *Client {
	return &Client{Client: c}
}
