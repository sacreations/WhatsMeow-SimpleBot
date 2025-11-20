package handlers

import (
	"fmt"
	"whatsappBotGo/src/senders"

	"go.mau.fi/whatsmeow/types"
)

// MediaHandler handles media-related operations (document uploads, image processing, etc.)
type MediaHandler struct {
	senders *senders.Senders
}

// NewMediaHandler creates a new media handler
func NewMediaHandler(s *senders.Senders) *MediaHandler {
	return &MediaHandler{senders: s}
}

// HandleDocumentUpload processes document uploads
func (m *MediaHandler) HandleDocumentUpload(to types.JID, filePath, fileName string) error {
	if m.senders == nil || m.senders.Document == nil {
		return fmt.Errorf("document sender not configured")
	}
	return m.senders.Document.SendDocument(to, filePath, fileName)
}

// HandleImageUpload processes image uploads
func (m *MediaHandler) HandleImageUpload(to types.JID, filePath, caption string) error {
	if m.senders == nil || m.senders.Image == nil {
		return fmt.Errorf("image sender not configured")
	}
	return m.senders.Image.SendImage(to, filePath, caption)
}

// HandleVideoUpload processes video uploads
func (m *MediaHandler) HandleVideoUpload(to types.JID, filePath, caption string) error {
	if m.senders == nil || m.senders.Video == nil {
		return fmt.Errorf("video sender not configured")
	}
	return m.senders.Video.SendVideo(to, filePath, caption)
}
