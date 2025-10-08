# Thomas Discord Bot (Golang)

🚀 **Lightweight Discord bot built with Golang** - Perfect for 1GB RAM VPS!

## ✨ Features

- **General Commands**: ping, help, info, server, status
- **Music Commands**: play, stop, skip, queue (YouTube support)
- **Super Lightweight**: ~30-40MB RAM usage (vs Node.js ~200MB)
- **Fast Startup**: <1 second
- **Single Binary**: No dependencies to install

## 📊 Performance

| Metric | Node.js | Golang |
|--------|---------|--------|
| RAM Usage | ~200MB | ~40MB |
| Startup Time | ~3s | <1s |
| Binary Size | N/A | ~12MB |
| CPU Usage | Medium | Low |

## 🎵 Music System

Uses:
- `yt-dlp` for YouTube download
- `ffmpeg` for audio processing
- Native Golang audio streaming

**Note**: Music playback is simplified in current version. Full DCA (Discord Audio) implementation coming soon for better quality.

## 🚀 Deployment

### Automatic (GitHub Actions)

1. Push to `main` branch
2. Workflow will:
   - Build Golang binary
   - Stop old Node.js bot (PM2)
   - Deploy to `/var/www/thomas-bot`
   - Install yt-dlp + ffmpeg
   - Create systemd service
   - Start bot

### Manual

```bash
# Build
go build -o thomas-bot .

# Run
./thomas-bot
```

## 🔧 Configuration

Create `.env` file:

```env
DISCORD_TOKEN=your_token_here
CLIENT_ID=your_client_id
PUBLIC_KEY=your_public_key
PORT=8080
```

## 📝 Commands

**Prefix:** `!T`

### General
- `!T ping` - Check bot latency
- `!T help` - Show commands
- `!T info` - Bot information
- `!T server` - Server info
- `!T status` - Bot status

### Music
- `!T play <song>` - Play from YouTube
- `!T stop` - Stop and clear queue
- `!T skip` - Skip current song
- `!T queue` - Show queue

## 🛠️ Development

```bash
# Install dependencies
go mod download

# Run locally
go run .

# Build for production
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o thomas-bot .
```

## 📦 Dependencies

- `github.com/bwmarrin/discordgo` - Discord API
- `github.com/joho/godotenv` - Environment variables
- `yt-dlp` (binary) - YouTube download
- `ffmpeg` (binary) - Audio processing

## 🔄 Migration from Node.js

The workflow automatically:
1. ✅ Stops PM2 process `thomas-bot`
2. ✅ Removes `/var/www/discord-bot` (Node.js)
3. ✅ Creates `/var/www/thomas-bot` (Golang)
4. ✅ Installs systemd service

**No manual intervention needed!**

## 📋 System Requirements

- **OS**: Linux (Ubuntu/Debian)
- **RAM**: 512MB minimum (1GB recommended)
- **Go**: 1.21+ (for building)
- **Runtime**: Only binary + yt-dlp + ffmpeg

## 🎯 Why Golang?

- **5x lighter** than Node.js music bot
- **10x faster** startup
- **Zero npm install** nightmares
- **Single binary** deployment
- **Stable** for 1GB RAM VPS

## 🆘 Troubleshooting

### Check bot status
```bash
sudo systemctl status thomas-bot
```

### View logs
```bash
sudo journalctl -u thomas-bot -f
```

### Restart bot
```bash
sudo systemctl restart thomas-bot
```

## 📜 License

ISC

---

🎉 **Made with Golang & ❤️**
