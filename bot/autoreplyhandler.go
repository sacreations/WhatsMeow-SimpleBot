package bot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// Config holds configuration for the AutoReplyHandler
type Config struct {
	VideoAPIEndpoint    string
	VideoAPIKey         string
	VideoAPITimeout     int
	VideoQuality        string
	VideoFormat         string
	TempDir             string
	MaxFileSize         string
	CleanupAfterSend    bool
	EnableVideoDownload bool
}

// AutoReplyHandler handles automatic replies and special message processing
type AutoReplyHandler struct {
	client       *whatsmeow.Client
	youtubeRegex *regexp.Regexp
	tiktokRegex  *regexp.Regexp
	config       *Config
}

// NewAutoReplyHandler creates a new AutoReplyHandler instance
func NewAutoReplyHandler() *AutoReplyHandler {
	// Load configuration from environment variables
	config := &Config{
		VideoAPIEndpoint:    getEnv("VIDEO_API_ENDPOINT", "https://api.dummy-video-downloader.com/download"),
		VideoAPIKey:         getEnv("VIDEO_API_KEY", ""),
		VideoAPITimeout:     getEnvInt("VIDEO_API_TIMEOUT", 30),
		VideoQuality:        getEnv("VIDEO_QUALITY", "720p"),
		VideoFormat:         getEnv("VIDEO_FORMAT", "mp4"),
		TempDir:             getEnv("TEMP_DIR", "./tmp"),
		MaxFileSize:         getEnv("MAX_FILE_SIZE", "50MB"),
		CleanupAfterSend:    getEnvBool("CLEANUP_AFTER_SEND", true),
		EnableVideoDownload: getEnvBool("ENABLE_VIDEO_DOWNLOAD", true),
	}

	// Create temp directory if it doesn't exist
	os.MkdirAll(config.TempDir, 0755)

	// Regex patterns for YouTube and TikTok links
	youtubeRegex := regexp.MustCompile(`(?i)(youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/shorts\/)([a-zA-Z0-9_-]+)`)
	tiktokRegex := regexp.MustCompile(`(?i)(tiktok\.com\/.*\/video\/|vm\.tiktok\.com\/|tiktok\.com\/@.*\/video\/)([a-zA-Z0-9_-]+)`)

	return &AutoReplyHandler{
		youtubeRegex: youtubeRegex,
		tiktokRegex:  tiktokRegex,
		config:       config,
	}
}

// SetClient sets the WhatsApp client for sending media
func (a *AutoReplyHandler) SetClient(client *whatsmeow.Client) {
	a.client = client
}

// ProcessMessage processes incoming messages and returns appropriate responses
func (a *AutoReplyHandler) ProcessMessage(message string, sender types.JID) string {
	lowerMessage := strings.ToLower(strings.TrimSpace(message))

	// Check for video links first
	if a.containsVideoLink(message) {
		go a.handleVideoDownload(message, sender)
		return "🎥 Video link detected! I'm downloading and processing it for you. Please wait..."
	}

	// Common greetings
	greetings := []string{"hello", "hi", "hey", "good morning", "good afternoon", "good evening", "good night"}
	for _, greeting := range greetings {
		if strings.Contains(lowerMessage, greeting) {
			return "👋 Hello! How can I help you today? Type /help to see available commands."
		}
	}

	// How are you responses
	howAreYou := []string{"how are you", "how do you do", "what's up", "whats up", "wassup"}
	for _, phrase := range howAreYou {
		if strings.Contains(lowerMessage, phrase) {
			return "😊 I'm doing great, thank you for asking! I'm here and ready to help. How about you?"
		}
	}

	// Who are you responses
	whoAreYou := []string{"who are you", "what are you", "who is this", "what is this"}
	for _, phrase := range whoAreYou {
		if strings.Contains(lowerMessage, phrase) {
			return "🤖 I'm a WhatsApp bot built with Go! I can help you with various commands and tasks. Type /help to see what I can do!"
		}
	}

	// Thank you responses
	thankYou := []string{"thank you", "thanks", "thx", "thank u"}
	for _, phrase := range thankYou {
		if strings.Contains(lowerMessage, phrase) {
			return "😊 You're welcome! Happy to help. Is there anything else I can do for you?"
		}
	}

	// Goodbye responses
	goodbye := []string{"bye", "goodbye", "see you", "catch you later", "talk to you later", "ttyl"}
	for _, phrase := range goodbye {
		if strings.Contains(lowerMessage, phrase) {
			return "👋 Goodbye! Have a great day! Feel free to message me anytime you need help."
		}
	}

	// Help requests
	helpRequests := []string{"help", "what can you do", "commands", "options"}
	for _, phrase := range helpRequests {
		if strings.Contains(lowerMessage, phrase) {
			return "🆘 I can help you with many things! Type /help to see all available commands, or just chat with me!"
		}
	}

	// Default response for unrecognized messages
	defaultResponses := []string{
		"🤔 I'm not sure how to respond to that, but I'm here to help! Type /help for available commands.",
		"💭 Interesting! I'm still learning. Try typing /help to see what I can do for you.",
		"🤖 I didn't quite understand that. Type /help to see my available commands!",
	}

	// Return a random default response
	return defaultResponses[time.Now().Unix()%int64(len(defaultResponses))]
}

// containsVideoLink checks if the message contains YouTube or TikTok links
func (a *AutoReplyHandler) containsVideoLink(message string) bool {
	if !a.config.EnableVideoDownload {
		return false
	}
	return a.youtubeRegex.MatchString(message) || a.tiktokRegex.MatchString(message)
}

// handleVideoDownload downloads and sends video from YouTube or TikTok links
func (a *AutoReplyHandler) handleVideoDownload(message string, sender types.JID) {
	var videoURL string
	var platform string

	// Detect platform and extract video ID
	if a.youtubeRegex.MatchString(message) {
		platform = "YouTube"
		matches := a.youtubeRegex.FindStringSubmatch(message)
		if len(matches) > 2 {
			videoURL = fmt.Sprintf("https://www.youtube.com/watch?v=%s", matches[2])
		}
	} else if a.tiktokRegex.MatchString(message) {
		platform = "TikTok"
		// For TikTok, we'll use the original URL
		tiktokMatches := a.tiktokRegex.FindStringSubmatch(message)
		if len(tiktokMatches) > 0 {
			videoURL = tiktokMatches[0]
		}
	}

	if videoURL == "" {
		a.sendMessage(sender, "❌ Could not extract video URL from the message.")
		return
	}

	// Download video using dummy API (replace with actual implementation)
	videoPath, err := a.downloadVideoFromAPI(videoURL, platform)
	if err != nil {
		a.sendMessage(sender, fmt.Sprintf("❌ Failed to download %s video: %v", platform, err))
		return
	}

	// Send the downloaded video
	err = a.sendVideoMessage(sender, videoPath, fmt.Sprintf("🎥 Downloaded from %s", platform))
	if err != nil {
		a.sendMessage(sender, fmt.Sprintf("❌ Failed to send video: %v", err))
		return
	}

	// Clean up temporary file if configured to do so
	if a.config.CleanupAfterSend {
		os.Remove(videoPath)
	}

	a.sendMessage(sender, fmt.Sprintf("✅ Successfully downloaded and sent %s video!", platform))
}

// downloadVideoFromAPI downloads video using the configured API endpoint
func (a *AutoReplyHandler) downloadVideoFromAPI(videoURL, platform string) (string, error) {
	// Create request payload
	payload := fmt.Sprintf(`{
        "url": "%s",
        "platform": "%s",
        "quality": "%s",
        "format": "%s"
    }`, videoURL, platform, a.config.VideoQuality, a.config.VideoFormat)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(a.config.VideoAPITimeout) * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", a.config.VideoAPIEndpoint, strings.NewReader(payload))
	if err != nil {
		return a.createDummyVideo(platform)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if a.config.VideoAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+a.config.VideoAPIKey)
		// or req.Header.Set("X-API-Key", a.config.VideoAPIKey)
	}

	// Make API request
	resp, err := client.Do(req)
	if err != nil {
		// For demo purposes, create a dummy video file
		return a.createDummyVideo(platform)
	}
	defer resp.Body.Close()

	// Create temporary file in configured temp directory
	tmpFile, err := os.CreateTemp(a.config.TempDir, fmt.Sprintf("video_%s_*.%s", platform, a.config.VideoFormat))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// Copy response body to file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save video: %v", err)
	}

	return tmpFile.Name(), nil
}

// createDummyVideo creates a dummy video file for demonstration
func (a *AutoReplyHandler) createDummyVideo(platform string) (string, error) {
	// Create a dummy video file in configured temp directory
	tmpFile, err := os.CreateTemp(a.config.TempDir, fmt.Sprintf("dummy_%s_video_*.%s", platform, a.config.VideoFormat))
	if err != nil {
		return "", fmt.Errorf("failed to create dummy video file: %v", err)
	}
	defer tmpFile.Close()

	// Write some dummy content to make it a valid file
	dummyContent := fmt.Sprintf("Dummy %s video content - replace with actual video download implementation", platform)
	_, err = tmpFile.WriteString(dummyContent)
	if err != nil {
		return "", fmt.Errorf("failed to write dummy content: %v", err)
	}

	return tmpFile.Name(), nil
}

// sendVideoMessage sends a video message to the specified JID
func (a *AutoReplyHandler) sendVideoMessage(to types.JID, videoPath, caption string) error {
	if a.client == nil {
		return fmt.Errorf("WhatsApp client not set")
	}

	// Read video file
	videoData, err := os.ReadFile(videoPath)
	if err != nil {
		return fmt.Errorf("failed to read video file: %v", err)
	}

	// Upload video to WhatsApp
	uploaded, err := a.client.Upload(context.Background(), videoData, "video")
	if err != nil {
		return fmt.Errorf("failed to upload video: %v", err)
	}

	// Create video message
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

	// Send message
	msg := &waE2E.Message{
		VideoMessage: videoMsg,
	}

	_, err = a.client.SendMessage(context.Background(), to, msg)
	return err
}

// sendMessage sends a text message
func (a *AutoReplyHandler) sendMessage(to types.JID, text string) {
	if a.client == nil {
		return
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	a.client.SendMessage(context.Background(), to, msg)
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
