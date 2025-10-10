package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var (
	lavalinkClient disgolink.Client
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

	userID, err := snowflake.Parse(session.State.User.ID)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %v", err)
	}

	lavalinkClient = disgolink.New(userID,
		disgolink.WithListenerFunc(onPlayerUpdate),
		disgolink.WithListenerFunc(onTrackStart),
		disgolink.WithListenerFunc(onTrackEnd),
		disgolink.WithListenerFunc(onTrackException),
	)

	// Add node
	node, err := lavalinkClient.AddNode(context.Background(), disgolink.NodeConfig{
		Name:     "local",
		Address:  "127.0.0.1:2333",
		Password: "youshallnotpass",
		Secure:   false,
	})

	if err != nil {
		return fmt.Errorf("failed to add node: %v", err)
	}

	log.Printf("✅ Lavalink node added: %s", node.Config().Name)
	return nil
}

func onPlayerUpdate(player disgolink.Player, event lavalink.PlayerUpdateMessage) {
	// Player position update
}

func onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	log.Printf("🎵 Track started: %s", event.Track.Info.Title)
}

func onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	log.Printf("⏹ Track ended: %s (reason: %s)", event.Track.Info.Title, event.Reason)

	if event.Reason == lavalink.TrackEndReasonFinished || event.Reason == lavalink.TrackEndReasonLoadFailed {
		go handleNextTrack(player.GuildID().String())
	}
}

func onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	log.Printf("❌ Track exception: %s - %s", event.Track.Info.Title, event.Exception.Message)
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
	searchPrefix := lavalink.SearchTypeYouTube
	if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
		searchPrefix = ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var track lavalink.Track

	// Build query string
	queryStr := query
	if searchPrefix != "" {
		queryStr = searchPrefix.Apply(query)
	}

	// Load tracks with handler
	var foundTrack *lavalink.Track
	var loadErr error

	lavalinkClient.BestNode().LoadTracksHandler(ctx, queryStr, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			foundTrack = &track
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) > 0 {
				foundTrack = &playlist.Tracks[0]
			}
		},
		func(tracks []lavalink.Track) {
			if len(tracks) > 0 {
				foundTrack = &tracks[0]
			}
		},
		func() {
			// No matches
		},
		func(err error) {
			log.Printf("ERROR: Load failed: %v", err)
			loadErr = err
		},
	))

	if loadErr != nil {
		errMsg := fmt.Sprintf("❌ Failed to search!\n```\nError: %v```", loadErr)
		log.Printf("ERROR: Failed to search: %v", loadErr)
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, errMsg)
		return
	}

	if foundTrack == nil {
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, "❌ No results found!")
		return
	}

	track = *foundTrack

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
		duration := formatDuration(time.Duration(track.Info.Length) * time.Millisecond)
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
	if _, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, false); err != nil {
		log.Printf("ERROR: Failed to join voice channel: %v", err)
		s.ChannelMessageSend(textChannelID, "❌ Failed to join voice channel!")
		player.mu.Lock()
		player.IsPlaying = false
		player.mu.Unlock()
		return
	}

	// Voice state is updated via Discord gateway events automatically (voiceStateUpdate handler)
	guildSnowflake, _ := snowflake.Parse(guildID)

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
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Update(ctx, lavalink.WithTrack(track)); err != nil {
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
	guildSnowflake, _ := snowflake.Parse(guildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Update(ctx, lavalink.WithTrack(nextTrack)); err != nil {
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
	guildSnowflake, _ := snowflake.Parse(m.GuildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Destroy(ctx); err != nil {
		log.Printf("ERROR: Failed to stop player: %v", err)
	}

	// Clear queue
	player.Queue = []lavalink.Track{}
	player.IsPlaying = false

	// Leave voice channel
	s.ChannelVoiceJoin(m.GuildID, "", false, false)

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
	guildSnowflake, _ := snowflake.Parse(m.GuildID)
	ctx := context.Background()

	if err := lavalinkClient.Player(guildSnowflake).Update(ctx, lavalink.WithNullTrack()); err != nil {
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
		duration := formatDuration(time.Duration(track.Info.Length) * time.Millisecond)
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
