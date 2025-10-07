# 🤖 Thomas Discord Bot

Discord bot dengan prefix `!T` yang di-deploy ke Google Cloud Platform menggunakan Cloud Run.

## ✨ Features

- ✅ **Command Handler** - Menangani berbagai command dengan prefix !T
- ✅ **Ping Command** - Cek latency bot
- ✅ **Help Command** - Menampilkan daftar command
- ✅ **Info Command** - Informasi tentang bot
- ✅ **Server Command** - Informasi server Discord
- ✅ **Auto Deploy** - GitHub Actions workflow untuk deploy ke GCP

## 📋 Prerequisites

- Node.js 18 atau lebih tinggi
- Discord Bot Token dari [Discord Developer Portal](https://discord.com/developers/applications)
- Google Cloud Platform account untuk deployment

## 🚀 Quick Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd thomas-bot-dc
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Setup Environment Variables

Buat file `.env` dengan menyalin dari `.env.example`:

```bash
cp .env.example .env
```

Edit file `.env` dan isi dengan credentials Anda:

```env
DISCORD_TOKEN=your_discord_bot_token_here
CLIENT_ID=your_application_id_here
PUBLIC_KEY=your_public_key_here
NODE_ENV=development
```

⚠️ **PENTING:** Jangan pernah commit file `.env` ke repository!

### 4. Run Bot Locally

```bash
# Development mode dengan auto-reload
npm run dev

# Production mode
npm start
```

## 📁 Project Structure

```
thomas-bot-dc/
├── index.js                    # Main bot file
├── package.json                # Node.js dependencies
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
├── README.md                   # This file
└── .github/
    └── workflows/
        └── deploy.yml          # GitHub Actions deployment workflow
```

## 🎮 Available Commands

| Command | Description |
|---------|-------------|
| `!T ping` | Check bot latency |
| `!T help` | Show command list |
| `!T info` | Show bot information |
| `!T server` | Show server information |

## 🌐 Deploy to Google Cloud Platform

### Prerequisites

1. **Google Cloud Project** dengan billing enabled
2. **Service Account** dengan permissions:
   - Cloud Run Admin
   - Service Account User
   - Cloud Build Editor
   - Artifact Registry Administrator

### Setup GitHub Secrets

Tambahkan secrets berikut di GitHub repository (`Settings` → `Secrets and variables` → `Actions`):

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `GOOGLE_CREDENTIALS` | JSON key dari GCP Service Account | `{"type": "service_account"...}` |
| `DISCORD_TOKEN` | Token Discord bot | `NzQ3NDM4NDkxNzYzMjEyMzc4...` |
| `CLIENT_ID` | Application ID Discord | `747438491763212378` |
| `PUBLIC_KEY` | Public Key Discord | `4a594400b4d623f6465...` |

### Deploy Configuration

Workflow ini akan deploy bot ke:
- **Service Name**: `thomas-bot-dc`
- **Region**: `asia-northeast1` (Tokyo)
- **Memory**: 512Mi
- **CPU**: 1 core
- **Min Instances**: 1 (bot selalu running)
- **Max Instances**: 3

### Deploy

Push ke branch `main` atau trigger workflow secara manual:

```bash
git add .
git commit -m "Deploy Thomas Discord Bot"
git push origin main
```

Atau trigger manual di GitHub:
1. Buka repository → `Actions`
2. Pilih workflow `Deploy Discord Bot to Google Cloud Run`
3. Klik `Run workflow`

GitHub Actions akan otomatis:
1. Build Docker image dari Dockerfile
2. Push ke Google Container Registry
3. Deploy ke Cloud Run
4. Bot langsung aktif dan running 24/7

## 🔧 Configuration

### Bot Intents

Bot ini menggunakan intents berikut:
- `Guilds` - Akses informasi server
- `GuildMessages` - Membaca messages di server
- `MessageContent` - Membaca isi message (privileged intent)
- `GuildMembers` - Akses informasi member

⚠️ **Penting:** Enable `Message Content Intent` di Discord Developer Portal untuk bot berfungsi dengan baik.

### Discord Developer Portal Setup

1. Buka [Discord Developer Portal](https://discord.com/developers/applications)
2. Pilih aplikasi bot Anda
3. Pergi ke **Bot** section
4. Enable **Message Content Intent** di Privileged Gateway Intents
5. Copy **Token** untuk environment variable `DISCORD_TOKEN`
6. Pergi ke **General Information**
7. Copy **Application ID** untuk environment variable `CLIENT_ID`
8. Copy **Public Key** untuk environment variable `PUBLIC_KEY`

### Invite Bot ke Server

Generate invite link dengan permissions yang diperlukan:

```
https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=2147519488&scope=bot
```

Ganti `YOUR_CLIENT_ID` dengan Application ID bot Anda.

## 📝 Development

### Menambahkan Command Baru

Edit file `index.js` dan tambahkan command handler di dalam event `messageCreate`:

```javascript
// Command: your-command
else if (command === 'your-command') {
  message.reply('Your response here!');
}
```

### Testing

```bash
# Run in development mode
npm run dev
```

Bot akan auto-reload saat ada perubahan code.

## 🐛 Troubleshooting

### Bot tidak merespon command

1. Pastikan **Message Content Intent** sudah enabled
2. Cek apakah bot memiliki permission `Read Messages` dan `Send Messages`
3. Pastikan prefix benar (`!T`)

### Deploy gagal di GCP

1. Cek apakah semua GitHub Secrets sudah diset dengan benar
2. Pastikan GCP Project sudah enable billing
3. Cek logs di GitHub Actions untuk error detail
4. Pastikan Service Account memiliki permissions yang tepat

## 📄 License

Lihat file [LICENSE](LICENSE) untuk detail.

## 🤝 Contributing

Pull requests welcome! Untuk perubahan major, mohon buka issue terlebih dahulu.

---

**Created with ❤️ using Discord.js and Google Cloud Platform**
