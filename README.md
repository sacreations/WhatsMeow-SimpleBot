# WhatsMeow SimpleBot — Single WhatsApp Instance

>⚠️ Experimental — This project is in an early stage. Expect rapid changes and incomplete features. Use for testing, experimentation, or as a reference implementation

This repository provides a single-instance WhatsApp bot built with Go and the whatsmeow library. Each running instance is meant to manage a single WhatsApp session (a single WhatsApp account). If you need multi-tenant or multi-account behavior, run multiple instances — one per account.

Overview:
- REST API for sending messages (text/media)
- WebSocket server for broadcasting incoming messages & events
- Separate sender implementations for text, image, video, and document
- Configurable video downloader integration (external API)

---

## Quick start (development)

1. Clone and prepare dependencies

```bash
git clone https://github.com/sacreations/WhatsMeow-SimpleBot.git
cd WhatsMeow-SimpleBot
go mod tidy
```

2. Set required environment variables (PowerShell examples):

```pwsh
$Env:INSTANCE_USER_ID = "instance-1"         # Required to restrict API to single instance
$Env:API_ADDR = ":8080"                      # API listen address (optional)
$Env:TEMP_DIR = "./tmp"                      # Temp directory for downloads
$Env:ENABLE_VIDEO_DOWNLOAD = "true"          # Automatic video download
$Env:CLEANUP_AFTER_SEND = "true"
```

3. Run the bot (runs both bot and API):

```bash
go run ./src/main.go

# Or build a binary
go build -o whatsmeow-bot ./src
./whatsmeow-bot
```

4. Scan the QR code printed in the console with the WhatsApp app (first run only).

---

## API Reference (single-instance)

When the bot runs, it exposes a simple HTTP API. If `INSTANCE_USER_ID` is set, the API requires a `user_id` field in every request and will reject mismatched requests.

### Send text message

POST /api/send/text

Request (JSON):

```json
{
    "jid": "1234567890@s.whatsapp.net",
    "text": "Hello!",
    "user_id": "instance-1"
}
```

Success response:

```json
{ "status": "ok" }
```

Error response structure:

```json
{ "error": "<message>" }
```

### Send image / video / document

- POST `/api/send/image` — accepts `jid`, optional `url` or `file` path, and `caption`
- POST `/api/send/video` — accepts `jid`, optional `url` or `file` path, and `caption`
- POST `/api/send/document` — accepts `jid`, optional `url` or `file` path, and `title`

Request bodies follow the same pattern as the text endpoint. If `url` is provided, the API will download the file to `TEMP_DIR` and then upload it to WhatsApp.

---

## WebSocket subscription

Endpoint: `ws://<host>:<port>/ws`

The hub broadcasts incoming messages and events as JSON payloads. Example payload:

```json
{
    "from": "1234567890@s.whatsapp.net",
    "text": "Hello",
    "event": "message",
    "user_id": "instance-1"
}
```

---

## Configuration & Tuning

- `INSTANCE_USER_ID` — unique id for the instance (recommended). If set, `user_id` is required on API requests and must match.
- `API_ADDR` — API listener address, default `:8080`.
- `TEMP_DIR` — temp directory for downloads (default `./tmp`).
- `VIDEO_API_ENDPOINT` — external video downloader API endpoint (optional).
- `VIDEO_API_KEY` — API key for video downloader.
- `VIDEO_API_TIMEOUT` — timeout seconds for the video API.
- `VIDEO_QUALITY` — default `720p`.
- `VIDEO_FORMAT` — default `mp4`.
- `ENABLE_VIDEO_DOWNLOAD` — `true`/`false`.
- `CLEANUP_AFTER_SEND` — `true`/`false`.

Note: `session.db` is created in the repo root and used to persist the whatsmeow session. To run multiple independent instances, run each instance in a separate working directory (unique `session.db` per instance) or modify the source to parameterize the DB filename.

---

## Security & production

- API is not authenticated by default. Add token-based auth (or other) in the API server (`src/api`) for production.
- Use HTTPS in front of the API or a reverse proxy.
- Use firewall rules to restrict access to the API.

---

## Files & Structure (source of truth)

```
src/
├─ main.go
├─ bot/
│  ├─ bot.go
│  └─ command_handler.go
├─ handlers/
│  └─ autoreplyhandler.go
├─ senders/
│  ├─ sender.go
│  ├─ text_sender.go
│  ├─ image_sender.go
│  ├─ video_sender.go
│  ├─ document_sender.go
│  └─ factory.go
├─ commands/
│  ├─ ping.go
+│  ├─ time.go
    └─ ...
```

---

## Designed for orchestration by a main system

This repository is designed as a single-instance worker for the primary WhatsApp bot project. It is not intended to be a multi-tenant gateway — instead, a separate main system (or supervisor) runs and manages many such instances where: start, stop, restart, delete, and add new instances are orchestrated centrally.

Typical orchestration responsibilities include:

- Start a new instance: create a working directory, set `INSTANCE_USER_ID`, and start the process
- Restart an instance: stop the process and start it again (optionally with updated configuration)
- Delete an instance: stop, cleanup `session.db`, and remove the working directory
- Add new instance: allocate instance id and port/volume and start under the supervisor

This worker is designed to be used by a supervisory system that may also implement additional features such as:

- centralized logging and monitoring
- auto-restart on failures
- automated reloading of configuration
- secure secrets and API key distribution
- scaling and health checks


