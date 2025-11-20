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

// clientImageSender implements ImageSender using a whatsmeow client
type clientImageSender struct {
	client *whatsmeow.Client
}

func NewImageSender(client *whatsmeow.Client) ImageSender {
	return &clientImageSender{client: client}
}

func (s *clientImageSender) SendImage(to types.JID, imagePath, caption string) error {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}
	uploaded, err := s.client.Upload(context.Background(), data, "image")
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	imgMsg := &waE2E.ImageMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String("image/jpeg"),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uploaded.FileLength),
		Caption:       proto.String(caption),
	}

	msg := &waE2E.Message{ImageMessage: imgMsg}
	_, err = s.client.SendMessage(context.Background(), to.ToNonAD(), msg)
	if err != nil {
		return fmt.Errorf("failed to send image message: %w", err)
	}
	return nil
}
