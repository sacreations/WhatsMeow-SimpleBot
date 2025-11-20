package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents health check data
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
}

// HandleHealth returns a health check response
func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(s.startTime).String(),
		Version:   "1.0.0",
	})
}

// ControlResponse represents control command response
type ControlResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// HandleControlStop stops the bot
func (s *Server) HandleControlStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ControlResponse{
		Message: "Bot stopping...",
		Success: true,
	})

	// Signal stop (if needed)
	if s.stopChan != nil {
		go func() { s.stopChan <- true }()
	}
}

// HandleControlStatus returns bot status
func (s *Server) HandleControlStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ControlResponse{
		Message: "Bot is running",
		Success: true,
	})
}
