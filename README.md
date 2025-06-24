# WhatsApp Bot using whatsmeow

A simple and well-structured WhatsApp bot built with Go and whatsmeow library.

## Features

- 🤖 Automated message responses with intelligent pattern matching
- 📱 QR code authentication
- 💾 Session persistence with SQLite
- 🔧 Command-based interactions
- 📝 Comprehensive logging
- 💬 Smart autoreply for natural conversations
- 🎥 **NEW:** Video download from YouTube and TikTok links

## Available Commands

- `/help` - Show available commands
- `/ping` - Check if bot is alive
- `/time` - Get current time
- `/echo <text>` - Echo your message
- `/info` - Get chat information
- `/joke` - Get a random joke
- `/quote` - Get an inspirational quote

## Autoreply Features

The bot automatically responds to common greetings and phrases:

- Hello, Hi, Good morning/evening/night
- How are you, What's up, Who are you
- Thank you, Thanks, Bye, Goodbye
- And many more natural conversation patterns!

## Video Download Feature

🎥 **Automatic Video Downloads**: Simply send a YouTube or TikTok link and the bot will:

- Detect the video link automatically
- Download the video using a configurable API
- Send the video back to you
- Support for both YouTube and TikTok platforms

Supported link formats:

- YouTube: `youtube.com/watch?v=`, `youtu.be/`, `youtube.com/shorts/`
- TikTok: `tiktok.com/@user/video/`, `vm.tiktok.com/`, etc.

## Setup

1. Initialize Go module (if not already done):

    ```bash
    go mod init whatsappBotGo
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    go get go.mau.fi/whatsmeow@latest
    go get github.com/mattn/go-sqlite3
    ```

3. Run the bot:

    ```bash
    go run main.go
    ```

4. Scan the QR code with your WhatsApp to authenticate

## Project Structure

```text
whatsappBotGo/
├── main.go                    # Entry point
├── bot/
│   ├── bot.go                # Main bot logic
│   ├── command_handler.go    # Command processing
│   └── autoreplyhandler.go   # Automatic reply and video download handling
├── commands/                 # Individual command implementations
│   ├── help.go
│   ├── ping.go
│   ├── time.go
│   ├── echo.go
│   ├── info.go
│   ├── joke.go
│   └── quote.go
├── go.mod                    # Go module file
├── session.db               # SQLite session storage (auto-created)
└── README.md                # This file
```

## Dependencies

- `go.mau.fi/whatsmeow@latest` - WhatsApp Web API library (Version: v0.0.0-...-947866b)
- `go.mau.fi/libsignal` - Signal protocol implementation
- `github.com/mattn/go-sqlite3` - SQLite driver for session storage
- `google.golang.org/protobuf` - Protocol Buffers for message handling

## Usage

1. Start the bot
2. Scan QR code with WhatsApp
3. Send messages to the bot number
4. Use commands like `/help` to interact
5. **NEW:** Send YouTube or TikTok links to download videos automatically

The bot will respond to various commands and messages automatically.

## Video Download Configuration

To configure the video download API:

1. Edit `bot/autoreplyhandler.go`
2. Replace the dummy API endpoint in `downloadVideoFromAPI()` function
3. Update the API request format according to your chosen video download service

## Notes

- The bot uses the latest version of whatsmeow package
- Session data is automatically saved to `session.db`
- QR code authentication is required only on first run
- Bot supports both individual and group chats
- Video download feature uses a dummy API by default - replace with actual implementation
- Downloaded videos are automatically cleaned up after sending
