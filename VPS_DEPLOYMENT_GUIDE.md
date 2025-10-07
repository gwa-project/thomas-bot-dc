# 🚀 Deploy Discord Music Bot ke Alibaba VPS

**VPS Specs:**
- IP: 147.139.241.73
- RAM: 1GB
- CPU: 1 Core
- OS: (akan kita cek)

---

## 📋 Step-by-Step Deployment

### Step 1: Connect ke VPS

**Windows (PowerShell/CMD):**
```powershell
ssh root@147.139.241.73
# Password: Jamban99@
```

**Atau pakai PuTTY:**
- Host: 147.139.241.73
- Port: 22
- Username: root
- Password: Jamban99@

---

### Step 2: Update System & Install Dependencies

```bash
# Update system
apt update && apt upgrade -y

# Install Node.js 18.x (LTS)
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# Install Git
apt install -y git

# Install FFmpeg (untuk music bot)
apt install -y ffmpeg

# Install PM2 (process manager)
npm install -g pm2

# Verify installations
node -v        # Should show v18.x.x
npm -v         # Should show npm version
ffmpeg -version # Should show ffmpeg version
pm2 -v         # Should show pm2 version
```

---

### Step 3: Clone Bot dari GitHub

```bash
# Pergi ke home directory
cd ~

# Clone repository (ganti dengan URL repo Anda)
git clone https://github.com/YOUR_USERNAME/thomas-bot-dc.git

# Masuk ke folder bot
cd thomas-bot-dc

# Install dependencies
npm install
```

---

### Step 4: Setup Environment Variables

```bash
# Buat file .env
nano .env
```

**Paste ini ke dalam file .env:**
```env
DISCORD_TOKEN=NzQ3NDM4NDkxNzYzMjEyMzc4.GJm7Vs.r7ViLEi7iQIR27xqrio3PTFzQN3LlCK4kt-NdY
CLIENT_ID=747438491763212378
PUBLIC_KEY=4a594400b4d623f6465ecd74efa9b7259152320fe2604d7019d58a5592686dda
NODE_ENV=production
PORT=8080
```

**Save & Exit:**
- Press `CTRL + X`
- Press `Y`
- Press `Enter`

---

### Step 5: Test Bot Lokal

```bash
# Test run bot
node index.js
```

**Kalau berhasil, akan muncul:**
```
🌐 HTTP server listening on port 8080
========================================
✅ Bot is online as THOMAS#5268
📝 Prefix: !T
🆔 Client ID: 747438491763212378
🌐 Serving X servers
🔧 Loaded X commands
========================================
🎭 Status set to: Idle | Listening to new music
```

**Test di Discord:**
- Join voice channel
- Ketik: `!T play never gonna give you up`

**Kalau sudah OK, stop bot:**
- Press `CTRL + C`

---

### Step 6: Setup PM2 untuk Auto-Start

```bash
# Start bot dengan PM2
pm2 start index.js --name thomas-bot

# Save PM2 process list
pm2 save

# Setup PM2 auto-start on reboot
pm2 startup

# Copy paste command yang muncul, contoh:
# sudo env PATH=$PATH:/usr/bin pm2 startup systemd -u root --hp /root

# Enable PM2 startup
systemctl enable pm2-root
```

---

### Step 7: PM2 Commands

```bash
# Lihat status bot
pm2 status

# Lihat logs real-time
pm2 logs thomas-bot

# Restart bot
pm2 restart thomas-bot

# Stop bot
pm2 stop thomas-bot

# Delete bot dari PM2
pm2 delete thomas-bot

# Monitor (dashboard)
pm2 monit
```

---

### Step 8: Update Bot di Future

Kalau ada update code:

```bash
# SSH ke VPS
ssh root@147.139.241.73

# Masuk ke folder bot
cd ~/thomas-bot-dc

# Pull latest changes
git pull

# Install new dependencies (kalau ada)
npm install

# Restart bot
pm2 restart thomas-bot
```

---

## 🔥 Firewall Setup (Optional)

Kalau port 8080 perlu dibuka (untuk health check):

```bash
# Install ufw
apt install -y ufw

# Allow SSH (PENTING! Jangan lupa ini!)
ufw allow 22/tcp

# Allow HTTP server (optional)
ufw allow 8080/tcp

# Enable firewall
ufw enable

# Check status
ufw status
```

---

## 🗑️ Hapus Cloud Run Service di GCP

### Method 1: Via gcloud CLI

```bash
# Install gcloud CLI dulu (kalau belum)
# https://cloud.google.com/sdk/docs/install

# Login
gcloud auth login

# Set project
gcloud config set project gwa-project-472118

# Delete Cloud Run service
gcloud run services delete thomas-bot-dc --region=asia-southeast1

# Konfirmasi dengan ketik: Y
```

### Method 2: Via Cloud Console (Web)

1. Buka: https://console.cloud.google.com/run?project=gwa-project-472118
2. Cari service: **thomas-bot-dc**
3. Klik checkbox di sebelah service name
4. Klik **DELETE** di atas
5. Konfirmasi delete

### Method 3: Via GitHub Actions (Disable Workflow)

Agar tidak auto-deploy lagi:

1. Buka: `.github/workflows/deploy.yml`
2. Tambahkan di line pertama:
```yaml
# DISABLED - Bot moved to VPS
name: Deploy Discord Bot to Google Cloud Run (DISABLED)
```
3. Atau delete file `deploy.yml` sepenuhnya

---

## 📊 Monitoring di VPS

### Check Resource Usage

```bash
# Check CPU & RAM usage
htop

# Check disk usage
df -h

# Check bot process
ps aux | grep node

# Check PM2 status
pm2 status

# Check logs
pm2 logs thomas-bot --lines 100
```

---

## 🐛 Troubleshooting

### Bot Tidak Connect ke Discord

```bash
# Check logs
pm2 logs thomas-bot

# Check environment variables
cat .env

# Restart bot
pm2 restart thomas-bot
```

### Music Tidak Play

```bash
# Check ffmpeg installed
ffmpeg -version

# Check bot logs
pm2 logs thomas-bot

# Restart bot
pm2 restart thomas-bot
```

### VPS Kehabisan Memory

```bash
# Check memory usage
free -h

# Restart bot
pm2 restart thomas-bot

# Kalau perlu, tambah swap
dd if=/dev/zero of=/swapfile bs=1M count=1024
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

### Bot Crash Terus

```bash
# Check error logs
pm2 logs thomas-bot --err

# Try run manual untuk debug
cd ~/thomas-bot-dc
node index.js
```

---

## ✅ Keuntungan Deploy di VPS

| Feature | Cloud Run | VPS Alibaba |
|---------|-----------|-------------|
| **Music Bot Stability** | ❌ Sering disconnect | ✅ Stable 24/7 |
| **Latency** | ❌ High | ✅ Low |
| **Cost** | ❌ $20-50/month | ✅ Sudah bayar VPS |
| **Control** | ❌ Limited | ✅ Full control |
| **Reliability** | ❌ Container restart | ✅ Always on |

---

## 🎯 Next Steps After Deployment

1. ✅ Test semua commands (`!T ping`, `!T help`, etc)
2. ✅ Test music commands (`!T play`, `!T queue`, etc)
3. ✅ Monitor logs untuk errors
4. ✅ Setup backup/update workflow
5. ✅ Delete Cloud Run service

---

## 📞 Support

Kalau ada masalah:
1. Check PM2 logs: `pm2 logs thomas-bot`
2. Check bot status: `pm2 status`
3. Restart bot: `pm2 restart thomas-bot`

---

🎉 **Bot sekarang running di VPS yang dedicated dan stable!**
