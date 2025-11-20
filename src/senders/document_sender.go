package senders

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// clientDocumentSender implements DocumentSender using a whatsmeow client
type clientDocumentSender struct {
	client *whatsmeow.Client
}

func NewDocumentSender(client *whatsmeow.Client) DocumentSender {
	return &clientDocumentSender{client: client}
}

func (s *clientDocumentSender) SendDocument(to types.JID, docPath, title string) error {
	return s.SendDocumentWithQuote(to, docPath, title, nil)
}

func (s *clientDocumentSender) SendDocumentWithQuote(to types.JID, docPath, title string, quotedMsg *QuotedMessage) error {
	data, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("failed to read document: %w", err)
	}
	uploaded, err := s.client.Upload(context.Background(), data, "document")
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}

	docMsg := &waE2E.DocumentMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String("application/octet-stream"),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uploaded.FileLength),
		Title:         proto.String(title),
	}

	if quotedMsg != nil {
		docMsg.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(quotedMsg.MessageID),
			Participant:   proto.String(quotedMsg.Sender.String()),
			QuotedMessage: quotedMsg.Message,
		}
	}

	msg := &waE2E.Message{DocumentMessage: docMsg}
	_, err = s.client.SendMessage(context.Background(), to.ToNonAD(), msg)
	if err != nil {
		return fmt.Errorf("failed to send document message: %w", err)
	}
	return nil
}
