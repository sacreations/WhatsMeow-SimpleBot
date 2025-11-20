package senders

import "go.mau.fi/whatsmeow/types"

// TextSender sends text messages
type TextSender interface {
	SendText(to types.JID, text string) error
}

// ImageSender sends image messages
type ImageSender interface {
	SendImage(to types.JID, imagePath, caption string) error
}

// VideoSender sends video messages
type VideoSender interface {
	SendVideo(to types.JID, videoPath, caption string) error
}

// DocumentSender sends document messages
type DocumentSender interface {
	SendDocument(to types.JID, docPath, title string) error
}

// Senders aggregates all sender interfaces
type Senders struct {
	Text     TextSender
	Image    ImageSender
	Video    VideoSender
	Document DocumentSender
}
