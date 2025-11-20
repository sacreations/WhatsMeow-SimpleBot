package senders

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

// clientTextSender implements TextSender using a whatsmeow client
type clientTextSender struct {
	client *whatsmeow.Client
}

func NewTextSender(client *whatsmeow.Client) TextSender {
	return &clientTextSender{client: client}
}

func (s *clientTextSender) SendText(to types.JID, text string) error {
	msg := &waProto.Message{
		Conversation: &text,
	}
	_, err := s.client.SendMessage(context.Background(), to.ToNonAD(), msg)
	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}
	return nil
}
