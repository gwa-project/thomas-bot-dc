package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const PREFIX = "!T"

var (
	Token     string
	ClientID  string
	PublicKey string
	startTime time.Time
)

func main() {
	// Load .env file
	godotenv.Load()

	Token = os.Getenv("DISCORD_TOKEN")
	ClientID = os.Getenv("CLIENT_ID")
	PublicKey = os.Getenv("PUBLIC_KEY")

	if Token == "" {
		log.Fatal("DISCORD_TOKEN is required")
	}

	startTime = time.Now()

	// Create Discord session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// Register event handlers
	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)

	// Set intents
	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildVoiceStates

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer dg.Close()

	// Start HTTP server for health checks
	go startHTTPServer(dg)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("========================================")
	fmt.Printf("✅ Bot is online as %s\n", s.State.User.String())
	fmt.Printf("📝 Prefix: %s\n", PREFIX)
	fmt.Printf("🆔 Client ID: %s\n", ClientID)
	fmt.Printf("🌐 Serving %d servers\n", len(event.Guilds))
	fmt.Println("========================================")

	// Set status: Idle + Listening to new music
	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "idle",
		Activities: []*discordgo.Activity{
			{
				Name: "new music",
				Type: discordgo.ActivityTypeListening,
			},
		},
	})
	if err != nil {
		log.Println("Error setting status:", err)
	}

	fmt.Println("🎭 Status set to: Idle | Listening to new music")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check if message starts with prefix
	if !strings.HasPrefix(m.Content, PREFIX) {
		return
	}

	// Parse command
	content := strings.TrimPrefix(m.Content, PREFIX)
	args := strings.Fields(content)

	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])
	commandArgs := args[1:]

	// Execute command
	switch command {
	case "ping":
		handlePing(s, m)
	case "help":
		handleHelp(s, m)
	case "info":
		handleInfo(s, m)
	case "server":
		handleServer(s, m)
	case "status":
		handleStatus(s, m)
	case "play":
		handlePlay(s, m, commandArgs)
	case "stop":
		handleStop(s, m)
	case "skip":
		handleSkip(s, m)
	case "queue":
		handleQueue(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Unknown command. Use `%shelp` to see available commands.", PREFIX))
	}
}

func startHTTPServer(dg *discordgo.Session) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		uptime := time.Since(startTime).Seconds()
		guilds := len(dg.State.Guilds)

		fmt.Fprintf(w, `{
			"status": "online",
			"bot": "%s",
			"uptime": %.0f,
			"guilds": %d,
			"prefix": "%s"
		}`, dg.State.User.String(), uptime, guilds, PREFIX)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy"}`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 HTTP server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
