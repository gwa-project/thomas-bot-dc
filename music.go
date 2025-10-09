package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type MusicQueue struct {
	GuildID    string
	Songs      []Song
	IsPlaying  bool
	VoiceConn  *discordgo.VoiceConnection
	mu         sync.Mutex
}

type Song struct {
	Title     string
	URL       string
	Requester string
}

var queues = make(map[string]*MusicQueue)
var queuesMu sync.Mutex

func getQueue(guildID string) *MusicQueue {
	queuesMu.Lock()
	defer queuesMu.Unlock()

	if q, exists := queues[guildID]; exists {
		return q
	}

	q := &MusicQueue{
		GuildID: guildID,
		Songs:   []Song{},
	}
	queues[guildID] = q
	return q
}

func handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Check if user is in voice channel
	guild, _ := s.State.Guild(m.GuildID)
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

	query := joinArgs(args)

	log.Printf("Play command from %s in guild %s: %s", m.Author.Username, m.GuildID, query)

	// Send searching message
	searchMsg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🔍 Searching for **%s**...", query))

	// Search with yt-dlp
	title, url, err := searchYouTube(query)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to find song!\n```\nError: %v```", err)
		log.Printf("ERROR: Failed to search YouTube for '%s': %v", query, err)
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, errMsg)
		return
	}

	log.Printf("Found song: %s (%s)", title, url)

	// Add to queue
	queue := getQueue(m.GuildID)
	song := Song{
		Title:     title,
		URL:       url,
		Requester: m.Author.Username,
	}

	queue.mu.Lock()
	queue.Songs = append(queue.Songs, song)
	queueLength := len(queue.Songs)
	isPlaying := queue.IsPlaying
	queue.mu.Unlock()

	if queueLength == 1 && !isPlaying {
		// Start playing
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, fmt.Sprintf("🎵 Now playing: **%s**", title))
		go playMusic(s, m.GuildID, voiceChannelID, m.ChannelID, queue)
	} else {
		// Added to queue
		s.ChannelMessageEdit(m.ChannelID, searchMsg.ID, fmt.Sprintf("➕ Added to queue: **%s** (Position: %d)", title, queueLength))
	}
}

func handleStop(s *discordgo.Session, m *discordgo.MessageCreate) {
	queue := getQueue(m.GuildID)

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if !queue.IsPlaying {
		s.ChannelMessageSend(m.ChannelID, "❌ Nothing is playing!")
		return
	}

	queue.Songs = []Song{}
	queue.IsPlaying = false

	if queue.VoiceConn != nil {
		queue.VoiceConn.Disconnect()
		queue.VoiceConn = nil
	}

	s.ChannelMessageSend(m.ChannelID, "⏹ Stopped the music!")
}

func handleSkip(s *discordgo.Session, m *discordgo.MessageCreate) {
	queue := getQueue(m.GuildID)

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if !queue.IsPlaying || len(queue.Songs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ Nothing is playing!")
		return
	}

	// Skip current song
	if len(queue.Songs) > 0 {
		queue.Songs = queue.Songs[1:]
	}

	s.ChannelMessageSend(m.ChannelID, "⏭ Skipped the song!")

	// Will be handled by playMusic loop
}

func handleQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	queue := getQueue(m.GuildID)

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if len(queue.Songs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "❌ Queue is empty!")
		return
	}

	queueText := ""
	for i, song := range queue.Songs {
		if i == 0 {
			queueText += fmt.Sprintf("**Now Playing:**\n🎵 %s\n\n", song.Title)
		} else if i < 10 {
			queueText += fmt.Sprintf("**%d.** %s\n", i, song.Title)
		}
	}

	if len(queue.Songs) > 10 {
		queueText += fmt.Sprintf("\n...and %d more songs", len(queue.Songs)-10)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📜 Music Queue",
		Description: queueText,
		Color:       0x0099ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Total: %d song(s)", len(queue.Songs)),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func searchYouTube(query string) (string, string, error) {
	log.Printf("Searching YouTube for: %s", query)

	// Use yt-dlp to search with cookies
	cmd := exec.Command("yt-dlp",
		"--cookies", "/var/www/thomas-bot/youtube_cookies.txt",
		"--default-search", "ytsearch",
		"--get-title",
		"--get-id",
		"--no-playlist",
		query,
	)

	log.Printf("Running command: yt-dlp --default-search ytsearch --get-title --get-id --no-playlist %s", query)

	output, err := cmd.CombinedOutput() // Get both stdout and stderr
	if err != nil {
		log.Printf("ERROR: yt-dlp command failed: %v", err)
		log.Printf("ERROR: yt-dlp output: %s", string(output))
		return "", "", fmt.Errorf("yt-dlp failed: %v (output: %s)", err, string(output))
	}

	log.Printf("yt-dlp raw output: %s", string(output))

	// Parse output (first line = title, second line = ID)
	lines := splitLines(string(output))
	if len(lines) < 2 {
		log.Printf("ERROR: Invalid yt-dlp output, expected 2+ lines, got %d lines", len(lines))
		log.Printf("ERROR: Lines received: %v", lines)
		return "", "", fmt.Errorf("invalid output: expected 2+ lines, got %d (output: %s)", len(lines), string(output))
	}

	title := lines[0]
	videoID := lines[1]
	url := "https://youtube.com/watch?v=" + videoID

	log.Printf("Parsed - Title: %s, Video ID: %s", title, videoID)

	return title, url, nil
}

func playMusic(s *discordgo.Session, guildID, voiceChannelID, textChannelID string, queue *MusicQueue) {
	log.Printf("Starting playback for guild %s in voice channel %s", guildID, voiceChannelID)

	// Join voice channel
	vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, true)
	if err != nil {
		log.Printf("ERROR: Failed to join voice channel %s: %v", voiceChannelID, err)
		s.ChannelMessageSend(textChannelID, "❌ Failed to join voice channel!")
		return
	}

	log.Printf("Successfully joined voice channel %s", voiceChannelID)

	queue.mu.Lock()
	queue.VoiceConn = vc
	queue.IsPlaying = true
	queue.mu.Unlock()

	defer func() {
		queue.mu.Lock()
		queue.IsPlaying = false
		queue.VoiceConn = nil
		queue.mu.Unlock()
		vc.Disconnect()
	}()

	for {
		queue.mu.Lock()
		if len(queue.Songs) == 0 {
			queue.mu.Unlock()
			s.ChannelMessageSend(textChannelID, "✅ Queue finished!")
			break
		}
		currentSong := queue.Songs[0]
		queue.mu.Unlock()

		// Download and play with yt-dlp + ffmpeg
		// This is simplified - in production you'd want better audio handling
		_ = exec.Command("yt-dlp",
			"-f", "bestaudio",
			"-o", "-",
			currentSong.URL,
		)

		// Note: Actual audio streaming would require DCA encoding
		// This is a placeholder for the concept
		s.ChannelMessageSend(textChannelID, fmt.Sprintf("⚠️ Playing: **%s** (Music streaming coming soon - basic implementation)", currentSong.Title))

		// Simulate playback (in production, stream audio to voice)
		// time.Sleep(3 * time.Minute)

		// Remove played song
		queue.mu.Lock()
		if len(queue.Songs) > 0 {
			queue.Songs = queue.Songs[1:]
		}
		queue.mu.Unlock()
	}
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, c := range s {
		if c == '\n' {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
		} else if c != '\r' {
			current += string(c)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
