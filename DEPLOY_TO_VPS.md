# 🚀 Complete Guide: Deploy Discord Bot ke Alibaba VPS

---

## 📋 GitHub Secrets yang Dibutuhkan

Tambahkan secrets berikut di GitHub repository:

**Path:** `Settings` → `Secrets and variables` → `Actions` → `New repository secret`

| Secret Name | Value | Cara Dapat |
|-------------|-------|------------|
| `VPS_SSH_KEY` | Private SSH key | Dari VPS: `cat /root/.ssh/github_actions` |
| `DISCORD_TOKEN` | Discord bot token | Regenerate di Discord Developer Portal |
| `CLIENT_ID` | `747438491763212378` | Discord Developer Portal |
| `PUBLIC_KEY` | `4a594400b4d623f6465ecd74efa9b7259152320fe2604d7019d58a5592686dda` | Discord Developer Portal |

---

## 🔧 PART 1: Setup VPS (One-Time Setup)

### 1. Connect ke VPS

```bash
ssh root@147.139.241.73
# Password: Jamban99@
```

### 2. Copy Paste Script Ini (All-in-One Setup)

```bash
# Update system
apt update && apt upgrade -y

# Install Node.js 18.x
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# Install dependencies
apt install -y git ffmpeg

# Install PM2
npm install -g pm2

# Create deploy directory
mkdir -p /var/www/discord-bot
cd /var/www/discord-bot

# Generate SSH key for GitHub Actions
ssh-keygen -t ed25519 -C "github-actions" -f /root/.ssh/github_actions -N ""

# Add to authorized_keys
cat /root/.ssh/github_actions.pub >> /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys
chmod 700 /root/.ssh

# Setup PM2 auto-start
pm2 startup systemd
# (Copy paste command yang muncul jika ada)

echo "=========================================="
echo "✅ VPS Setup Complete!"
echo "=========================================="
echo ""
echo "📋 COPY THESE KEYS TO GITHUB SECRETS:"
echo ""
echo "1️⃣ VPS_SSH_KEY (PRIVATE KEY):"
echo "---"
cat /root/.ssh/github_actions
echo "---"
echo ""
echo "2️⃣ PUBLIC KEY (for reference):"
cat /root/.ssh/github_actions.pub
echo ""
echo "=========================================="
```

### 3. Copy SSH Private Key

Setelah script selesai, akan muncul **PRIVATE KEY**. Copy seluruh isi key (dari `-----BEGIN` sampai `-----END`).

**Contoh output:**
```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtz
...
...
-----END OPENSSH PRIVATE KEY-----
```

**SIMPAN INI!** Akan dipakai untuk GitHub Secret `VPS_SSH_KEY`.

---

## 🔐 PART 2: Setup GitHub Secrets

### 1. Buka GitHub Repository

```
https://github.com/YOUR_USERNAME/thomas-bot-dc/settings/secrets/actions
```

### 2. Add Secrets Satu per Satu

#### Secret 1: VPS_SSH_KEY

- Click **"New repository secret"**
- Name: `VPS_SSH_KEY`
- Secret: **Paste ENTIRE private key** (dari `-----BEGIN` sampai `-----END`)
- Click **"Add secret"**

#### Secret 2: DISCORD_TOKEN

**⚠️ IMPORTANT:** Regenerate token dulu!

1. Buka https://discord.com/developers/applications
2. Pilih aplikasi bot Anda (ID: 747438491763212378)
3. Pergi ke **Bot** section
4. Klik **"Reset Token"**
5. Konfirmasi & copy token baru
6. Add ke GitHub:
   - Name: `DISCORD_TOKEN`
   - Secret: Token baru (contoh: `NzQ3NDM4NDkxNzYzMjEyMzc4.GJm7Vs.xxxxxxxxxxxxx`)
   - Click **"Add secret"**

#### Secret 3: CLIENT_ID

- Name: `CLIENT_ID`
- Secret: `747438491763212378`
- Click **"Add secret"**

#### Secret 4: PUBLIC_KEY

- Name: `PUBLIC_KEY`
- Secret: `4a594400b4d623f6465ecd74efa9b7259152320fe2604d7019d58a5592686dda`
- Click **"Add secret"**

### 3. Verify Secrets

Setelah semua ditambahkan, harusnya ada 4 secrets:
- ✅ VPS_SSH_KEY
- ✅ DISCORD_TOKEN
- ✅ CLIENT_ID
- ✅ PUBLIC_KEY

---

## 🚀 PART 3: Deploy Bot

### Method 1: Push ke GitHub (Auto-Deploy)

```bash
cd thomas-bot-dc

git add .
git commit -m "Setup VPS deployment workflow"
git push origin main
```

GitHub Actions akan otomatis:
1. Connect ke VPS via SSH
2. Clone/pull repository
3. Install dependencies
4. Update .env file
5. Restart bot dengan PM2

**Monitor deployment:**
```
https://github.com/YOUR_USERNAME/thomas-bot-dc/actions
```

### Method 2: Manual Trigger

1. Buka: `https://github.com/YOUR_USERNAME/thomas-bot-dc/actions`
2. Pilih workflow: **"Deploy Discord Bot to VPS"**
3. Click **"Run workflow"**
4. Select branch: `main`
5. Click **"Run workflow"**

---

## ✅ PART 4: Verify Bot Running

### 1. Check via SSH

```bash
ssh root@147.139.241.73

# Check PM2 status
pm2 status

# View logs
pm2 logs thomas-bot --lines 50

# Monitor real-time
pm2 monit
```

### 2. Test di Discord

```
!T ping
!T help
!T play never gonna give you up
!T queue
!T nowplaying
```

---

## 🔧 PM2 Management Commands

```bash
# SSH ke VPS
ssh root@147.139.241.73

# View status
pm2 status

# View logs
pm2 logs thomas-bot

# Restart bot
pm2 restart thomas-bot

# Stop bot
pm2 stop thomas-bot

# Start bot
pm2 start thomas-bot

# Delete from PM2
pm2 delete thomas-bot

# Save configuration
pm2 save

# Monitor dashboard
pm2 monit
```

---

## 🐛 Troubleshooting

### Bot Tidak Start Setelah Deploy

```bash
ssh root@147.139.241.73
cd /var/www/discord-bot

# Check logs
pm2 logs thomas-bot --err

# Check .env file
cat .env

# Try manual start
node index.js
```

### GitHub Actions Deploy Gagal

**Error: "Permission denied (publickey)"**

Solution: Check `VPS_SSH_KEY` secret, pastikan full key ter-copy (termasuk BEGIN dan END).

**Error: "npm: command not found"**

Solution: Node.js belum terinstall di VPS
```bash
ssh root@147.139.241.73
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs
```

**Error: "pm2: command not found"**

Solution: PM2 belum terinstall
```bash
ssh root@147.139.241.73
npm install -g pm2
```

### Music Bot Tidak Play

```bash
ssh root@147.139.241.73

# Check ffmpeg
ffmpeg -version

# If not installed
apt install -y ffmpeg

# Restart bot
pm2 restart thomas-bot
```

---

## 📊 Monitoring

### Resource Usage

```bash
ssh root@147.139.241.73

# CPU & Memory
htop
# (Press q to quit)

# Disk usage
df -h

# Bot memory usage
pm2 monit
```

### Bot Logs

```bash
# Real-time logs
pm2 logs thomas-bot

# Last 100 lines
pm2 logs thomas-bot --lines 100

# Error logs only
pm2 logs thomas-bot --err

# Save logs to file
pm2 logs thomas-bot --lines 1000 > bot-logs.txt
```

---

## 🔄 Update Bot (Future Updates)

Setiap kali ada perubahan code:

```bash
# Local
cd thomas-bot-dc
git add .
git commit -m "Update bot features"
git push origin main
```

GitHub Actions akan otomatis deploy ke VPS! ✅

---

## 💾 Backup & Restore

### Backup Bot Data

```bash
ssh root@147.139.241.73

# Backup entire bot directory
tar -czf bot-backup-$(date +%Y%m%d).tar.gz /var/www/discord-bot

# Download to local
scp root@147.139.241.73:~/bot-backup-*.tar.gz ./
```

### Restore Bot

```bash
ssh root@147.139.241.73

# Extract backup
tar -xzf bot-backup-YYYYMMDD.tar.gz -C /

# Restart bot
cd /var/www/discord-bot
pm2 restart thomas-bot
```

---

## 🎯 Summary

| Step | Command | Status |
|------|---------|--------|
| 1. Setup VPS | Run setup script | ✅ Done once |
| 2. Get SSH Keys | `cat /root/.ssh/github_actions` | ✅ |
| 3. Add GitHub Secrets | 4 secrets | ✅ |
| 4. Push to GitHub | `git push origin main` | ✅ Auto-deploy |
| 5. Verify | `pm2 status` | ✅ Bot running |

---

## 🎉 Done!

Bot sekarang:
- ✅ Running 24/7 di VPS
- ✅ Auto-restart on crash (PM2)
- ✅ Auto-deploy on push (GitHub Actions)
- ✅ Stable music playback
- ✅ Low latency
- ✅ Full control

**Cost: $0 (already have VPS)** 🎊

---

Need help? Check logs: `pm2 logs thomas-bot`
