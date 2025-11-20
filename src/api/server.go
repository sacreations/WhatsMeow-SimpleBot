package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"whatsappBotGo/src/senders"

	"go.mau.fi/whatsmeow/types"
)

// Server holds references to Senders and a WebSocket Hub
type Server struct {
	senders        *senders.Senders
	Hub            *Hub
	httpServer     *http.Server
	TempDir        string
	InstanceUserID string
	startTime      time.Time
	stopChan       chan bool
}

func NewServer(s *senders.Senders, instanceUserID string) *Server {
	hub := NewHub()
	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	srv := &Server{
		senders:        s,
		Hub:            hub,
		TempDir:        os.TempDir(),
		InstanceUserID: instanceUserID,
		startTime:      time.Now(),
		stopChan:       make(chan bool, 1),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/send/text", srv.sendTextHandler)
	mux.HandleFunc("/api/send/image", srv.sendImageHandler)
	mux.HandleFunc("/api/send/video", srv.sendVideoHandler)
	mux.HandleFunc("/api/send/document", srv.sendDocumentHandler)
	mux.HandleFunc("/ws", hub.ServeWS)
	mux.HandleFunc("/api/health", srv.HandleHealth)
	mux.HandleFunc("/api/control/status", srv.HandleControlStatus)
	mux.HandleFunc("/api/control/stop", srv.HandleControlStop)

	srv.httpServer = &http.Server{Addr: addr, Handler: mux}
	return srv
}

func (s *Server) Start() {
	go s.Hub.Run()
	log.Printf("Starting API server on %s", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (s *Server) sendTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req SendTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}

	// Validate user id: if server has an instance bound, ensure the request matches
	if s.InstanceUserID != "" {
		if req.UserID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id required"})
			return
		}
		if req.UserID != s.InstanceUserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "user_id mismatch"})
			return
		}
	}

	jid, err := types.ParseJID(req.JID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid jid"})
		return
	}

	if s.senders == nil || s.senders.Text == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "text sender not configured"})
		return
	}

	if err := s.senders.Text.SendText(jid, req.Text); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("send failed: %v", err)})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// downloadFile downloads a remote file to the tempdir and returns path
func (s *Server) downloadFile(url string) (string, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp(s.TempDir, "api_media_*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func (s *Server) sendImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req SendMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}

	// Validate user id: if server has an instance bound, ensure the request matches
	if s.InstanceUserID != "" {
		if req.UserID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id required"})
			return
		}
		if req.UserID != s.InstanceUserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "user_id mismatch"})
			return
		}
	}

	jid, err := types.ParseJID(req.JID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid jid"})
		return
	}
	if s.senders == nil || s.senders.Image == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "image sender not configured"})
		return
	}

	path := req.File
	if path == "" && req.URL != "" {
		path, err = s.downloadFile(req.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("download failed: %v", err)})
			return
		}
		defer os.Remove(path)
	}

	if err := s.senders.Image.SendImage(jid, path, req.Caption); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("send failed: %v", err)})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) sendVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req SendMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}

	// Validate user id: if server has an instance bound, ensure the request matches
	if s.InstanceUserID != "" {
		if req.UserID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id required"})
			return
		}
		if req.UserID != s.InstanceUserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "user_id mismatch"})
			return
		}
	}

	jid, err := types.ParseJID(req.JID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid jid"})
		return
	}
	if s.senders == nil || s.senders.Video == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "video sender not configured"})
		return
	}

	path := req.File
	if path == "" && req.URL != "" {
		path, err = s.downloadFile(req.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("download failed: %v", err)})
			return
		}
		defer os.Remove(path)
	}

	if err := s.senders.Video.SendVideo(jid, path, req.Caption); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("send failed: %v", err)})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) sendDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req SendMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}

	// Validate user id: if server has an instance bound, ensure the request matches
	if s.InstanceUserID != "" {
		if req.UserID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id required"})
			return
		}
		if req.UserID != s.InstanceUserID {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "user_id mismatch"})
			return
		}
	}

	jid, err := types.ParseJID(req.JID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid jid"})
		return
	}
	if s.senders == nil || s.senders.Document == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "document sender not configured"})
		return
	}

	path := req.File
	if path == "" && req.URL != "" {
		path, err = s.downloadFile(req.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("download failed: %v", err)})
			return
		}
		defer os.Remove(path)
	}

	if err := s.senders.Document.SendDocument(jid, path, req.Title); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("send failed: %v", err)})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) BroadcastIncoming(msg WSMessage) {
	if s.Hub != nil {
		// annotate with instance user id if available
		if s.InstanceUserID != "" {
			msg.UserID = s.InstanceUserID
		}
		s.Hub.Broadcast(msg)
	}
}
