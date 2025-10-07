#!/bin/bash

# Discord Bot VPS Setup Script
# Alibaba Cloud VPS - 1GB RAM / 1 Core
# Run this on VPS: bash setup-vps.sh

set -e

echo "=========================================="
echo "🚀 Discord Music Bot - VPS Setup"
echo "=========================================="
echo ""

# Update system
echo "📦 Updating system..."
apt update && apt upgrade -y

# Install Node.js 18.x
echo "📦 Installing Node.js 18.x..."
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# Install dependencies
echo "📦 Installing dependencies..."
apt install -y git ffmpeg

# Install PM2
echo "📦 Installing PM2..."
npm install -g pm2

# Verify installations
echo ""
echo "✅ Verifying installations..."
node -v
npm -v
ffmpeg -version | head -n 1
pm2 -v

echo ""
echo "=========================================="
echo "✅ System setup complete!"
echo "=========================================="
echo ""
echo "📝 Next steps:"
echo "1. Clone your bot repository:"
echo "   git clone https://github.com/YOUR_USERNAME/thomas-bot-dc.git"
echo ""
echo "2. Go to bot directory:"
echo "   cd thomas-bot-dc"
echo ""
echo "3. Install dependencies:"
echo "   npm install"
echo ""
echo "4. Create .env file:"
echo "   nano .env"
echo "   (Copy content from .env.example)"
echo ""
echo "5. Start bot with PM2:"
echo "   pm2 start index.js --name thomas-bot"
echo "   pm2 save"
echo "   pm2 startup"
echo ""
echo "=========================================="
