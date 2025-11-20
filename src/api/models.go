package api

// JSON models for API requests

type SendTextRequest struct {
	JID    string `json:"jid"`
	Text   string `json:"text"`
	UserID string `json:"user_id,omitempty"`
}

type SendMediaRequest struct {
	JID     string `json:"jid"`
	URL     string `json:"url,omitempty"`  // if provided, server will download
	File    string `json:"file,omitempty"` // local file path
	Caption string `json:"caption,omitempty"`
	Title   string `json:"title,omitempty"` // for documents
	UserID  string `json:"user_id,omitempty"`
}

// Message struct broadcasted over WebSocket
type WSMessage struct {
	From    string `json:"from"`
	Text    string `json:"text,omitempty"`
	Event   string `json:"event"`
	RawType string `json:"raw_type,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

// Helper to convert string JID to types.JID is done in handlers
