# Quoted Message Implementation Guide

This document explains how to use the newly implemented quoted message feature in all senders.

## Overview

All senders now support sending messages with quoted/reply functionality:
- `TextSender` - Send text replies
- `ImageSender` - Send image replies with captions
- `VideoSender` - Send video replies with captions
- `DocumentSender` - Send document replies

## QuotedMessage Structure

```go
type QuotedMessage struct {
    MessageID  string            // The message ID of the message to reply to
    Sender     types.JID         // The JID of the sender of the original message
    Message    *waE2E.Message    // The original message proto
}
```

## Usage Examples

### Text Message with Quote

```go
// Create a quoted message reference
quotedMsg := &senders.QuotedMessage{
    MessageID:  info.MessageInfo.ID,           // From evt.Info.MessageInfo.ID
    Sender:     info.MessageInfo.Sender,       // From evt.Info.MessageInfo.Sender
    Message:    info.Message,                  // From evt.Info.Message
}

// Send text with quote
err := sender.Text.SendTextWithQuote(to, "This is a reply!", quotedMsg)
```

### Image Message with Quote

```go
quotedMsg := &senders.QuotedMessage{
    MessageID:  info.MessageInfo.ID,
    Sender:     info.MessageInfo.Sender,
    Message:    info.Message,
}

err := sender.Image.SendImageWithQuote(to, "/path/to/image.jpg", "Check this image!", quotedMsg)
```

### Video Message with Quote

```go
quotedMsg := &senders.QuotedMessage{
    MessageID:  info.MessageInfo.ID,
    Sender:     info.MessageInfo.Sender,
    Message:    info.Message,
}

err := sender.Video.SendVideoWithQuote(to, "/path/to/video.mp4", "Watch this!", quotedMsg)
```

### Document Message with Quote

```go
quotedMsg := &senders.QuotedMessage{
    MessageID:  info.MessageInfo.ID,
    Sender:     info.MessageInfo.Sender,
    Message:    info.Message,
}

err := sender.Document.SendDocumentWithQuote(to, "/path/to/document.pdf", "See attached", quotedMsg)
```

## Backward Compatibility

The original methods still work and don't require quotes:

```go
// These still work as before
sender.Text.SendText(to, "Hello!")
sender.Image.SendImage(to, "/path/to/image.jpg", "Caption")
sender.Video.SendVideo(to, "/path/to/video.mp4", "Caption")
sender.Document.SendDocument(to, "/path/to/document.pdf", "Title")
```

## Implementation Details

Each sender implementation:
1. Creates the appropriate message type (TextMessage, ImageMessage, etc.)
2. If a `QuotedMessage` is provided, adds `ContextInfo` with:
   - `StanzaID`: Message ID of the original message
   - `Participant`: Sender JID of the original message
   - `QuotedMessage`: The original message proto
3. Sends the message with the quote context

### Proto Structure Used

All senders use the `waE2E` (WhatsApp End-to-End) protocol buffers:
- Text: `ExtendedTextMessage` with `ContextInfo`
- Images: `ImageMessage` with `ContextInfo`
- Videos: `VideoMessage` with `ContextInfo`
- Documents: `DocumentMessage` with `ContextInfo`

## Integration in Bot

To integrate quoted messages in your command handlers:

```go
func (handler *MyHandler) HandleMessage(evt *events.Message) {
    // Extract quoted message info from incoming message
    quotedMsg := &senders.QuotedMessage{
        MessageID:  evt.Info.MessageInfo.ID,
        Sender:     evt.Info.MessageInfo.Sender,
        Message:    evt.Message,
    }
    
    // Send reply with quote
    handler.senders.Text.SendTextWithQuote(evt.Info.Sender, "Reply text", quotedMsg)
}
```
