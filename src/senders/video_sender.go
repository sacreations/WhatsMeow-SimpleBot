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

// clientVideoSender implements VideoSender using a whatsmeow client
type clientVideoSender struct {
	client *whatsmeow.Client
}

func NewVideoSender(client *whatsmeow.Client) VideoSender {
	return &clientVideoSender{client: client}
}

func (s *clientVideoSender) SendVideo(to types.JID, videoPath, caption string) error {
	data, err := os.ReadFile(videoPath)
	if err != nil {
		return fmt.Errorf("failed to read video: %w", err)
	}
	uploaded, err := s.client.Upload(context.Background(), data, "video")
	if err != nil {
		return fmt.Errorf("failed to upload video: %w", err)
	}

	videoMsg := &waE2E.VideoMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String("video/mp4"),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uploaded.FileLength),
		Caption:       proto.String(caption),
	}

	msg := &waE2E.Message{VideoMessage: videoMsg}
	_, err = s.client.SendMessage(context.Background(), to.ToNonAD(), msg)
	if err != nil {
		return fmt.Errorf("failed to send video message: %w", err)
	}
	return nil
}
