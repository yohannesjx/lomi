# Quick Deployment Guide

## What You Need

1. **Server**: Ubuntu/Debian VPS with Docker installed
2. **Domain**: Point `api.lomi.social` to your server IP
3. **GitHub**: Repository with your code
4. **Secrets**: Telegram bot token, R2 credentials, etc.

## Quick Setup (5 minutes)

### On Your Server:

```bash
# 1. Install Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install -y caddy docker.io docker-compose git

# 2. Clone your repo
cd /opt
git clone https://github.com/YOUR_USERNAME/lomi_mini.git
cd lomi_mini

# 3. Create .env.production
nano .env.production
# (Copy from .env.production.example and fill in values)

# 4. Setup Caddy
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo nano /etc/caddy/Caddyfile  # Update domain names
sudo systemctl enable caddy && sudo systemctl start caddy

# 5. Deploy
chmod +x deploy.sh
./deploy.sh
```

## Environment Variables Needed

Create `.env.production` with:

```bash
DB_PASSWORD=d5YhNXB5zXhT7bkbbQ7
REDIS_PASSWORD=r5YhNXB5zXhT7bkbbQ7
JWT_SECRET=9f3c8b27d4a1e6f2c9b0a54e83d712af
TELEGRAM_BOT_TOKEN=8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=d46ab6ad318b1127d061533769bce800
S3_SECRET_KEY=2f1a730b2b691e8fbc5a33a8595132846cb335c19c523b90a2de173705285c20


## GitHub Deployment

### Option 1: Manual (SSH into server)
```bash
cd /opt/lomi_mini
git pull origin main
./deploy.sh
```

### Option 2: GitHub Actions (Automatic)
1. Add secrets in GitHub: Settings → Secrets
   - `SERVER_HOST`: 72.62.37.181
   - `SERVER_USER`: root
   - `SERVER_SSH_KEY`: Private SSH key
2. Push to `main` branch → Auto-deploys!

## Verify

```bash
curl https://api.lomi.social/api/v1/health
# Should return: {"status":"ok","message":"Lomi Backend is running 🍋"}
```

## Why Caddy Separate?

✅ Manages SSL automatically (Let's Encrypt)  
✅ Needs direct access to ports 80/443  
✅ Standard practice for reverse proxies  
✅ Easier to update independently  

## Troubleshooting

```bash
# Check backend logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Check Caddy
sudo journalctl -u caddy -f

# Restart services
docker-compose -f docker-compose.prod.yml restart
sudo systemctl restart caddy
```

