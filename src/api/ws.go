package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains active connections and broadcasts to them
type Hub struct {
	conns      map[*websocket.Conn]bool
	broadcast  chan WSMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		conns:      make(map[*websocket.Conn]bool),
		broadcast:  make(chan WSMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.conns[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				c.Close()
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			// Broadcast to all connections
			h.mu.Lock()
			for c := range h.conns {
				if err := c.WriteJSON(msg); err != nil {
					log.Printf("ws write error: %v", err)
					delConn := c
					delete(h.conns, delConn)
					delConn.Close()
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}

	h.register <- c

	// keep reading to avoid connection closure on client senders
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			h.unregister <- c
			return
		}
	}
}

func (h *Hub) Broadcast(msg WSMessage) {
	h.broadcast <- msg
}
