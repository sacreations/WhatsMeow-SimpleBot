package senders

import (
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// QuotedMessage contains information about a message being replied to
type QuotedMessage struct {
	MessageID string         // The message ID of the message to reply to
	Sender    types.JID      // The JID of the sender of the original message
	Message   *waE2E.Message // The original message proto
}

// TextSender sends text messages
type TextSender interface {
	SendText(to types.JID, text string) error
	SendTextWithQuote(to types.JID, text string, quotedMsg *QuotedMessage) error
}

// ImageSender sends image messages
type ImageSender interface {
	SendImage(to types.JID, imagePath, caption string) error
	SendImageWithQuote(to types.JID, imagePath, caption string, quotedMsg *QuotedMessage) error
}

// VideoSender sends video messages
type VideoSender interface {
	SendVideo(to types.JID, videoPath, caption string) error
	SendVideoWithQuote(to types.JID, videoPath, caption string, quotedMsg *QuotedMessage) error
}

// DocumentSender sends document messages
type DocumentSender interface {
	SendDocument(to types.JID, docPath, title string) error
	SendDocumentWithQuote(to types.JID, docPath, title string, quotedMsg *QuotedMessage) error
}

// Senders aggregates all sender interfaces
type Senders struct {
	Text     TextSender
	Image    ImageSender
	Video    VideoSender
	Document DocumentSender
}
