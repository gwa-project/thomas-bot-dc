package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
)

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	msg, _ := s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")

	latency := time.Since(start).Milliseconds()

	s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("🏓 Pong! Latency: `%dms`", latency))
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📚 Bot Commands",
		Description: fmt.Sprintf("Prefix: `%s`", PREFIX),
		Color:       0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "📌 General Commands",
				Value:  "`ping` - Check bot latency\n`help` - Show this message\n`info` - Bot information\n`server` - Server information\n`status` - Bot status",
				Inline: false,
			},
			{
				Name:   "🎵 Music Commands",
				Value:  "`play <song>` - Play a song from YouTube\n`stop` - Stop music and clear queue\n`skip` - Skip current song\n`queue` - Show music queue",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made with Golang 🚀",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	uptime := time.Since(startTime)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	embed := &discordgo.MessageEmbed{
		Title: "ℹ️ Bot Information",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Bot Name",
				Value:  s.State.User.Username,
				Inline: true,
			},
			{
				Name:   "Bot ID",
				Value:  s.State.User.ID,
				Inline: true,
			},
			{
				Name:   "Prefix",
				Value:  fmt.Sprintf("`%s`", PREFIX),
				Inline: true,
			},
			{
				Name:   "Servers",
				Value:  fmt.Sprintf("%d", len(s.State.Guilds)),
				Inline: true,
			},
			{
				Name:   "Uptime",
				Value:  formatDuration(uptime),
				Inline: true,
			},
			{
				Name:   "Memory Usage",
				Value:  fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024),
				Inline: true,
			},
			{
				Name:   "Language",
				Value:  "Golang " + runtime.Version(),
				Inline: true,
			},
			{
				Name:   "Library",
				Value:  "DiscordGo",
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("256"),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Error getting server information")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "🏠 Server Information",
		Color: 0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server Name",
				Value:  guild.Name,
				Inline: true,
			},
			{
				Name:   "Server ID",
				Value:  guild.ID,
				Inline: true,
			},
			{
				Name:   "Owner",
				Value:  fmt.Sprintf("<@%s>", guild.OwnerID),
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprintf("%d", guild.MemberCount),
				Inline: true,
			},
			{
				Name:   "Created At",
				Value:  guild.ID, // Simplified - can parse snowflake for real date
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL("256"),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	uptime := time.Since(startTime)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	embed := &discordgo.MessageEmbed{
		Title: "📊 Bot Status",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  "🟢 Online",
				Inline: true,
			},
			{
				Name:   "Uptime",
				Value:  formatDuration(uptime),
				Inline: true,
			},
			{
				Name:   "Memory",
				Value:  fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024),
				Inline: true,
			},
			{
				Name:   "Goroutines",
				Value:  fmt.Sprintf("%d", runtime.NumGoroutine()),
				Inline: true,
			},
			{
				Name:   "Servers",
				Value:  fmt.Sprintf("%d", len(s.State.Guilds)),
				Inline: true,
			},
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
