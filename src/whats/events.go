package whats

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

// EventHandler is called when WhatsApp events arrive
type EventHandler interface {
	HandleMessage(*events.Message)
	HandleReceipt(*events.Receipt)
	HandlePresence(*events.Presence)
}

// Events package skeleton for future event routing and processing
type Events struct {
	handlers []EventHandler
}

// NewEvents creates a new Events dispatcher
func NewEvents() *Events {
	return &Events{handlers: make([]EventHandler, 0)}
}

// RegisterHandler registers an event handler
func (e *Events) RegisterHandler(h EventHandler) {
	e.handlers = append(e.handlers, h)
}

// Dispatch sends an event to all registered handlers
func (e *Events) Dispatch(ctx context.Context, evt interface{}) error {
	switch v := evt.(type) {
	case *events.Message:
		for _, h := range e.handlers {
			h.HandleMessage(v)
		}
	case *events.Receipt:
		for _, h := range e.handlers {
			h.HandleReceipt(v)
		}
	case *events.Presence:
		for _, h := range e.handlers {
			h.HandlePresence(v)
		}
	default:
		return fmt.Errorf("unknown event type: %T", evt)
	}
	return nil
}
