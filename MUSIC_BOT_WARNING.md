# ⚠️ Music Bot - Experimental on Cloud Run

## 🚨 IMPORTANT WARNING

Music bot sudah di-implement di bot ini, TAPI dengan **MAJOR LIMITATIONS** karena berjalan di Cloud Run:

---

## ❌ Masalah yang Akan Terjadi:

### 1. **Bot Bisa Disconnect Tiba-Tiba**
- Cloud Run bisa **restart container** kapan saja
- Voice connection akan **putus** saat restart
- Music akan **berhenti** mendadak

### 2. **High Latency**
- Voice connection dari Cloud Run → Discord akan **lambat**
- Audio bisa **stutter** atau **lag**
- Quality tidak sebaik bot yang di-host di VPS

### 3. **Cost Mahal**
- Music bot butuh **CPU tinggi** untuk encode audio
- Cloud Run charge per CPU usage
- Bisa **$20-50/month** atau lebih untuk always-on music bot

### 4. **Queue Hilang Saat Restart**
- Kalau container restart, **queue akan hilang**
- Harus play ulang dari awal

### 5. **Tidak Reliable untuk Production**
- **TIDAK DISARANKAN** untuk server Discord yang ramai
- Cocok hanya untuk **testing/experimental**

---

## ✅ Cara Pakai (Experimental)

### Music Commands:

| Command | Description | Example |
|---------|-------------|---------|
| `!T play <song>` | Play music from YouTube | `!T play never gonna give you up` |
| `!T pause` | Pause current song | `!T pause` |
| `!T resume` | Resume paused song | `!T resume` |
| `!T skip` | Skip current song | `!T skip` |
| `!T stop` | Stop music & clear queue | `!T stop` |
| `!T queue` | Show music queue | `!T queue` |
| `!T nowplaying` | Show current song | `!T nowplaying` |
| `!T volume <0-100>` | Set volume | `!T volume 50` |

### Setup:

1. **Join voice channel** di Discord server
2. Run command `!T play <song name>`
3. Bot akan join voice channel dan play musik

### Contoh Usage:

```
!T play Bohemian Rhapsody
!T play https://www.youtube.com/watch?v=dQw4w9WgXcQ
!T queue
!T volume 75
!T skip
```

---

## 📝 Technical Details

### Dependencies Added:
- `@discordjs/voice` - Voice connection
- `@discordjs/opus` - Audio encoding
- `play-dl` - YouTube audio streaming
- `ffmpeg-static` - Audio processing
- `sodium-native` / `libsodium-wrappers` - Encryption

### Files Created:
- [utils/musicManager.js](utils/musicManager.js) - Music queue & player management
- [commands/play.js](commands/play.js) - Play music command
- [commands/pause.js](commands/pause.js) - Pause command
- [commands/resume.js](commands/resume.js) - Resume command
- [commands/skip.js](commands/skip.js) - Skip command
- [commands/stop.js](commands/stop.js) - Stop command
- [commands/queue.js](commands/queue.js) - Show queue
- [commands/nowplaying.js](commands/nowplaying.js) - Current song info
- [commands/volume.js](commands/volume.js) - Volume control

---

## 💡 Recommended Alternative

Untuk music bot yang **stable dan reliable**, saya **SANGAT MENYARANKAN**:

### Option 1: VPS Hosting
Deploy bot di VPS seperti:
- **Digital Ocean** ($5/month)
- **Vultr** ($3.5/month)
- **Hetzner** (€4/month)
- **Oracle Cloud** (FREE tier - 2 ARM VMs)

### Option 2: Alternative Platforms
- **Railway.app** - Free tier, support long-running apps
- **Replit** - Free tier (limited)
- **Heroku** - Paid ($7/month minimum)

### Option 3: Use Existing Music Bots
Gunakan music bot yang sudah proven:
- **Groovy** (Paid)
- **FredBoat** (Self-host gratis)
- **Lavalink** (Self-host gratis)

---

## 🔧 Jika Tetap Mau Pakai di Cloud Run

Sudah saya implement dan **bisa dicoba**, tapi:

✅ **Good for:**
- Testing music bot features
- Low-usage servers (< 10 users)
- Short sessions (< 30 menit)

❌ **NOT good for:**
- Production servers
- Long playlists
- High-quality audio requirements
- 24/7 music streaming

---

## 📊 Estimated Cloud Run Cost

Dengan music bot enabled:

| Usage | Estimated Cost/Month |
|-------|---------------------|
| 1-2 hours/day | $10-15 |
| 5-8 hours/day | $25-40 |
| Always-on streaming | $50-100+ |

Compare dengan VPS: **$3-5/month unlimited**

---

## 🐛 Known Issues

1. ⚠️ Bot might disconnect during song playback
2. ⚠️ Queue resets when container restarts
3. ⚠️ High latency possible
4. ⚠️ Some songs might fail to play
5. ⚠️ Volume control might not work properly

---

## 📚 Further Reading

- [Discord.js Voice Guide](https://discordjs.guide/voice/)
- [Cloud Run Limitations](https://cloud.google.com/run/docs/issues)
- [play-dl Documentation](https://github.com/play-dl/play-dl)

---

## ✅ Status

🟡 **Music bot is ENABLED but EXPERIMENTAL**

Kalau ada masalah atau mau pindah ke VPS hosting, let me know! 😊

---

**Tetap mau deploy ke Cloud Run? Silakan push kode ini!**
**Mau diskusi alternative hosting? Saya siap bantu setup! 🚀**
