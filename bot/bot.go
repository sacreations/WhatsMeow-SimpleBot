package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"whatsappBotGo/commands"
)

type WhatsAppBot struct {
	client         *whatsmeow.Client
	store          *sqlstore.Container
	commandHandler *CommandHandler
}

// NewBot creates a new WhatsApp bot instance
func NewBot() (*WhatsAppBot, error) {
	ctx := context.Background()
	dbLog := waLog.Stdout("DB", "INFO", true)

	// Create SQLite store
	container, err := sqlstore.New(ctx, "sqlite3", "file:session.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %v", err)
	}

	// Get first device store
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// Create command handler and register commands
	commandHandler := NewCommandHandler()
	registerCommands(commandHandler)

	bot := &WhatsAppBot{
		client:         client,
		store:          container,
		commandHandler: commandHandler,
	}

	// Set client for auto reply handler
	commandHandler.SetClient(client)

	// Add event handler
	client.AddEventHandler(bot.eventHandler)

	return bot, nil
}

// registerCommands registers all available commands
func registerCommands(handler *CommandHandler) {
	handler.RegisterCommand(commands.NewHelpCommand())
	handler.RegisterCommand(commands.NewPingCommand())
	handler.RegisterCommand(commands.NewTimeCommand())
	handler.RegisterCommand(commands.NewEchoCommand())
	handler.RegisterCommand(commands.NewInfoCommand())
	handler.RegisterCommand(commands.NewJokeCommand())
	handler.RegisterCommand(commands.NewQuoteCommand())
}

// Start starts the WhatsApp bot
func (bot *WhatsAppBot) Start() error {
	if bot.client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := bot.client.GetQRChannel(context.Background())
		err := bot.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				// Print QR code visually in terminal
				qr := qrcodeTerminal.New()
				qr.Get(evt.Code).Print()
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := bot.client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
	}

	fmt.Println("Bot is running...")

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	bot.client.Disconnect()
	return nil
}

// eventHandler handles incoming WhatsApp events
func (bot *WhatsAppBot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		bot.handleMessage(v)
	case *events.Receipt:
		if v.Type == types.ReceiptTypeRead || v.Type == types.ReceiptTypeReadSelf {
			fmt.Printf("Message %s was read\n", v.MessageIDs[0])
		}
	case *events.Presence:
		if v.Unavailable {
			fmt.Printf("%s went offline\n", v.From)
		} else {
			fmt.Printf("%s is online\n", v.From)
		}
	}
}

// handleMessage processes incoming messages and responds accordingly
func (bot *WhatsAppBot) handleMessage(evt *events.Message) {
	if evt.Info.IsFromMe {
		return // Ignore own messages
	}

	msg := evt.Message
	if msg == nil {
		return
	}

	var messageText string
	if msg.GetConversation() != "" {
		messageText = msg.GetConversation()
	} else if msg.GetExtendedTextMessage() != nil {
		messageText = msg.GetExtendedTextMessage().GetText()
	} else {
		return // Not a text message
	}

	messageText = strings.TrimSpace(messageText)

	fmt.Printf("Received message from %s: %s\n", evt.Info.Sender, messageText)

	// Process message using command handler with sender information
	response := bot.commandHandler.ProcessMessage(messageText, evt.Info.Sender)
	if response != "" {
		bot.sendMessage(evt.Info.Sender, response)
	}
}

// sendMessage sends a text message to specified JID
func (bot *WhatsAppBot) sendMessage(to types.JID, text string) {
	msg := &waProto.Message{
		Conversation: &text,
	}

	_, err := bot.client.SendMessage(context.Background(), to, msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		fmt.Printf("Sent message to %s: %s\n", to, text)
	}
}

// Disconnect gracefully disconnects the bot
func (bot *WhatsAppBot) Disconnect() {
	bot.client.Disconnect()
}
