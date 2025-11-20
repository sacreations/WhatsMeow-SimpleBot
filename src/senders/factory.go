package senders

import "go.mau.fi/whatsmeow"

// NewSendersFromClient creates Senders using a whatsmeow client
func NewSendersFromClient(client *whatsmeow.Client) *Senders {
	return &Senders{
		Text:     NewTextSender(client),
		Image:    NewImageSender(client),
		Video:    NewVideoSender(client),
		Document: NewDocumentSender(client),
	}
}
