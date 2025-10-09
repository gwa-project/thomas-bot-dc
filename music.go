package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lavalink-devs/lavalink-go/lavalink"
)

var (
	lavalinkClient *lavalink.Client
	players        = make(map[string]*MusicPlayer)
	playersMu      sync.Mutex
)

type MusicPlayer struct {
	GuildID   string
	Queue     []lavalink.Track
	IsPlaying bool
	mu        sync.Mutex
}

func initLavalink(session *discordgo.Session) error {
	log.Println("🎵 Initializing Lavalink client...")

	lavalinkClient = lavalink.NewClient(
		snowflake.MustParse(session.State.User.ID),
	)

	node := lavalink.NodeConfig{
		Name:      "local",
		Address:   "127.0.0.1:2333",
		Password:  "youshallnotpass",
		Secure:    false,
		SessionID: "",
	}

	if err := lavalinkClient.AddNode(context.Background(), node); err != nil {
		return fmt.Errorf("failed to add Lavalink node: %v", err)
	}

	// Event handlers
	lavalinkClient.AddEventHandler(func(p lavalink.Player, track lavalink.Track) {
		log.Printf("🎵 Track started: %s", track.Info.Title)
	})

	lavalinkClient.AddEventHandler(func(p lavalink.Player, track lavalink.Track, err lavalink.Exception) {
		log.Printf("❌ Track exception: %s - %s", track.Info.Title, err.Message)
	})

	lavalinkClient.AddEventHandler(func(p lavalink.Player, track lavalink.Track, reason lavalink.TrackEndReason) {
		log.Printf("Track ended: %s (reason: %s)", track.Info.Title, reason)

		if reason == lavalink.TrackEndReasonFinished || reason == lavalink.TrackEndReasonLoadFailed {
			// Play next in queue
			go handleNextTrack(p.GuildID().String())
		}
	})

	log.Println("✅ Lavalink client initialized")
	return nil
}

func getPlayer(guildID string) *MusicPlayer {
	playersMu.Lock()
	defer playersMu.Unlock()

	if p, exists := players[guildID]; exists {
		return p
	}

	p := &MusicPlayer{
		GuildID: guildID,
		Queue:   []lavalink.Track{},
	}
	players[guildID] = p
	return p
}

func handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Check if user is in voice channel
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to get guild info!")
		return
	}

	var voiceChannelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	if voiceChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, "❌ You need to be in a voice channel!")
		return
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Please provide a song name or URL!\nUsage: `%splay <song name or URL>`", PREFIX))
		return
	}

	query := strings.Join(args, " ")
	log.Printf("Play command from %s in guild %s: %s", m.Author.Username, m.GuildID, query)

	// Send searching message
	searchMsg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🔍 Searching for **%s**...", query))

	// Search with Lavalink
	searchPrefix := "ytsearch:"
	if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
		searchPrefix = ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := lavalinkClient.LoadTracks(ctx, searchPrefix+query)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to search!\n```\nError: %v```", err)
		log.Printf("ERROR: Failed to search: %v", err)
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, errMsg)
		return
	}

	var track lavalink.Track

	switch result.LoadType {
	case lavalink.LoadTypeTrack:
		track = result.Data.(lavalink.Track)
	case lavalink.LoadTypeSearch:
		tracks := result.Data.([]lavalink.Track)
		if len(tracks) == 0 {
			s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, "❌ No results found!")
			return
		}
		track = tracks[0]
	case lavalink.LoadTypePlaylist:
		playlist := result.Data.(lavalink.Playlist)
		if len(playlist.Tracks) == 0 {
			s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, "❌ Playlist is empty!")
			return
		}
		track = playlist.Tracks[0]
	default:
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, "❌ Failed to load track!")
		return
	}

	log.Printf("Found track: %s - %s", track.Info.Title, track.Info.Author)

	// Get player
	player := getPlayer(m.GuildID)
	player.mu.Lock()

	// Add to queue
	player.Queue = append(player.Queue, track)
	queueLength := len(player.Queue)
	isPlaying := player.IsPlaying

	player.mu.Unlock()

	if queueLength == 1 && !isPlaying {
		// Start playing
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, fmt.Sprintf("🎵 Now playing: **%s** by %s", track.Info.Title, track.Info.Author))
		go playTrack(s, m.GuildID, voiceChannelID, m.ChannelID)
	} else {
		// Added to queue
		duration := formatDuration(track.Info.Length)
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, fmt.Sprintf("➕ Added to queue: **%s** by %s `[%s]` (Position: %d)", track.Info.Title, track.Info.Author, duration, queueLength))
	}
}

func playTrack(s *discordgo.Session, guildID, voiceChannelID, textChannelID string) {
	player := getPlayer(guildID)

	player.mu.Lock()
	if len(player.Queue) == 0 {
		player.IsPlaying = false
		player.mu.Unlock()
		return
	}
	player.IsPlaying = true
	player.mu.Unlock()

	// Join voice channel
	if err := s.ChannelVoiceJoinManual(guildID, voiceChannelID, false, false); err != nil {
		log.Printf("ERROR: Failed to join voice channel: %v", err)
		s.ChannelMessageSend(textChannelID, "❌ Failed to join voice channel!")
		player.mu.Lock()
		player.IsPlaying = false
		player.mu.Unlock()
		return
	}

	// Update voice server
	lavalinkClient.UpdateVoiceState(context.Background(), snowflake.MustParse(guildID), voiceChannelID, false, false)

	// Get current track
	player.mu.Lock()
	if len(player.Queue) == 0 {
		player.IsPlaying = false
		player.mu.Unlock()
		return
	}
	track := player.Queue[0]
	player.mu.Unlock()

	// Play track
	guildSnowflake := snowflake.MustParse(guildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Play(ctx, track); err != nil {
		log.Printf("ERROR: Failed to play track: %v", err)
		s.ChannelMessageSend(textChannelID, fmt.Sprintf("❌ Failed to play: %s", track.Info.Title))

		// Remove failed track and try next
		player.mu.Lock()
		if len(player.Queue) > 0 {
			player.Queue = player.Queue[1:]
		}
		player.mu.Unlock()

		go handleNextTrack(guildID)
		return
	}

	log.Printf("▶️ Playing: %s - %s", track.Info.Title, track.Info.Author)
}

func handleNextTrack(guildID string) {
	player := getPlayer(guildID)

	player.mu.Lock()
	if len(player.Queue) > 0 {
		player.Queue = player.Queue[1:]
	}

	if len(player.Queue) == 0 {
		player.IsPlaying = false
		player.mu.Unlock()
		log.Printf("Queue finished for guild %s", guildID)
		return
	}

	nextTrack := player.Queue[0]
	player.mu.Unlock()

	// Play next track
	guildSnowflake := snowflake.MustParse(guildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Play(ctx, nextTrack); err != nil {
		log.Printf("ERROR: Failed to play next track: %v", err)

		// Remove failed track and try next
		player.mu.Lock()
		if len(player.Queue) > 0 {
			player.Queue = player.Queue[1:]
		}
		player.mu.Unlock()

		go handleNextTrack(guildID)
	}
}

func handleStop(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := getPlayer(m.GuildID)

	player.mu.Lock()
	defer player.mu.Unlock()

	if !player.IsPlaying {
		s.ChannelMessageSend(m.ChannelID, "❌ Nothing is playing!")
		return
	}

	// Stop playback
	guildSnowflake := snowflake.MustParse(m.GuildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Destroy(ctx); err != nil {
		log.Printf("ERROR: Failed to stop player: %v", err)
	}

	// Clear queue
	player.Queue = []lavalink.Track{}
	player.IsPlaying = false

	// Leave voice channel
	if err := s.ChannelVoiceJoinManual(m.GuildID, "", false, false); err != nil {
		log.Printf("ERROR: Failed to leave voice channel: %v", err)
	}

	s.ChannelMessageSend(m.ChannelID, "⏹ Stopped the music!")
}

func handleSkip(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := getPlayer(m.GuildID)

	player.mu.Lock()
	defer player.mu.Unlock()

	if !player.IsPlaying || len(player.Queue) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ Nothing is playing!")
		return
	}

	// Skip current song
	guildSnowflake := snowflake.MustParse(m.GuildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Stop(ctx); err != nil {
		log.Printf("ERROR: Failed to skip: %v", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to skip!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "⏭ Skipped the song!")
}

func handleQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := getPlayer(m.GuildID)

	player.mu.Lock()
	defer player.mu.Unlock()

	if len(player.Queue) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ Queue is empty!")
		return
	}

	queueText := ""
	for i, track := range player.Queue {
		duration := formatDuration(track.Info.Length)
		if i == 0 {
			queueText += fmt.Sprintf("**Now Playing:**\n🎵 %s - %s `[%s]`\n\n", track.Info.Title, track.Info.Author, duration)
		} else if i < 10 {
			queueText += fmt.Sprintf("**%d.** %s - %s `[%s]`\n", i, track.Info.Title, track.Info.Author, duration)
		}
	}

	if len(player.Queue) > 10 {
		queueText += fmt.Sprintf("\n...and %d more songs", len(player.Queue)-10)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📜 Music Queue",
		Description: queueText,
		Color:       0x0099ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Total: %d song(s)", len(player.Queue)),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func formatDuration(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
