package senders

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// clientTextSender implements TextSender using a whatsmeow client
type clientTextSender struct {
	client *whatsmeow.Client
}

func NewTextSender(client *whatsmeow.Client) TextSender {
	return &clientTextSender{client: client}
}

func (s *clientTextSender) SendText(to types.JID, text string) error {
	return s.SendTextWithQuote(to, text, nil)
}

func (s *clientTextSender) SendTextWithQuote(to types.JID, text string, quotedMsg *QuotedMessage) error {
	textMsg := &waE2E.ExtendedTextMessage{
		Text: proto.String(text),
	}

	if quotedMsg != nil {
		textMsg.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(quotedMsg.MessageID),
			Participant:   proto.String(quotedMsg.Sender.String()),
			QuotedMessage: quotedMsg.Message,
		}
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: textMsg,
	}

	_, err := s.client.SendMessage(context.Background(), to.ToNonAD(), msg)
	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}
	return nil
}
