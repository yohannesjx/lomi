#!/bin/bash

# Initial Server Setup Script for Lomi Social
# Run this ONCE on your server to set everything up
# Usage: sudo ./initial-server-setup.sh

set -e

echo "🚀 Lomi Social - Initial Server Setup"
echo "======================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root or with sudo"
    exit 1
fi

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. Update system
echo "📦 Updating system packages..."
apt update && apt upgrade -y

# 2. Install Docker
if ! command -v docker &> /dev/null; then
    echo "🐳 Installing Docker..."
    apt install -y docker.io docker-compose
    systemctl enable docker
    systemctl start docker
else
    echo -e "${GREEN}✓ Docker already installed${NC}"
fi

# 3. Install Caddy
if ! command -v caddy &> /dev/null; then
    echo "🌐 Installing Caddy..."
    apt install -y debian-keyring debian-archive-keyring apt-transport-https
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
    apt update
    apt install -y caddy
    systemctl enable caddy
    systemctl start caddy
else
    echo -e "${GREEN}✓ Caddy already installed${NC}"
fi

# 4. Install Node.js (for webhook server)
if ! command -v node &> /dev/null; then
    echo "📦 Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt install -y nodejs
else
    echo -e "${GREEN}✓ Node.js already installed${NC}"
fi

# 5. Install PM2
if ! command -v pm2 &> /dev/null; then
    echo "📦 Installing PM2..."
    npm install -g pm2
else
    echo -e "${GREEN}✓ PM2 already installed${NC}"
fi

# 6. Install Git
if ! command -v git &> /dev/null; then
    echo "📦 Installing Git..."
    apt install -y git
else
    echo -e "${GREEN}✓ Git already installed${NC}"
fi

# 7. Create project directory
PROJECT_DIR="/opt/lomi_mini"
if [ ! -d "$PROJECT_DIR" ]; then
    echo "📁 Creating project directory..."
    mkdir -p $PROJECT_DIR
    echo -e "${YELLOW}⚠️  Project directory created at $PROJECT_DIR"
    echo "   Please clone your repository there:"
    echo "   cd /opt && git clone https://github.com/YOUR_USERNAME/lomi_mini.git${NC}"
else
    echo -e "${GREEN}✓ Project directory exists${NC}"
fi

# 8. Setup firewall
echo "🔥 Configuring firewall..."
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw allow 9000/tcp # Webhook (optional, can be behind Caddy)
ufw --force enable

# 9. Create log directory
mkdir -p /var/log
touch /var/log/lomi-webhook.log
chmod 666 /var/log/lomi-webhook.log

echo ""
echo -e "${GREEN}✅ Initial server setup complete!${NC}"
echo ""
echo "📋 Next steps:"
echo "   1. Clone your repository:"
echo "      cd /opt && git clone https://github.com/YOUR_USERNAME/lomi_mini.git"
echo ""
echo "   2. Create .env.production:"
echo "      cd /opt/lomi_mini"
echo "      cp .env.production.example .env.production"
echo "      nano .env.production  # Fill in your values"
echo ""
echo "   3. Setup Caddy:"
echo "      sudo cp Caddyfile /etc/caddy/Caddyfile"
echo "      sudo nano /etc/caddy/Caddyfile  # Update domain names"
echo "      sudo systemctl reload caddy"
echo ""
echo "   4. Setup webhook:"
echo "      cd /opt/lomi_mini"
echo "      sudo ./setup-webhook.sh"
echo ""
echo "   5. Initial deployment:"
echo "      ./deploy.sh"

