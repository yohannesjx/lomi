# 🚀 Quick Deployment Guide

## One-Command Deployment Workflow

After initial setup, deploy with one command from your local machine:

```bash
./deploy "Your commit message"
```

This will:
1. ✅ Commit all changes
2. ✅ Push to GitHub  
3. ✅ Trigger server deployment via webhook

## Initial Server Setup (One Time)

### On Your Server:

```bash
# 1. Clone repository
cd /opt
git clone https://github.com/YOUR_USERNAME/lomi_mini.git
cd lomi_mini

# 2. Run initial setup (installs Docker, Caddy, Node.js, etc.)
sudo chmod +x initial-server-setup.sh
sudo ./initial-server-setup.sh

# 3. Create .env.production
# (Copy values from DEPLOYMENT_QUICK_START.md - you already have them)
nano .env.production
# Paste your environment variables

# 4. Generate secure JWT secret
openssl rand -base64 32
# Add this to JWT_SECRET in .env.production

# 5. Setup Caddy
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo nano /etc/caddy/Caddyfile  # Update domain names if needed
sudo systemctl reload caddy

# 6. Setup webhook server
sudo chmod +x setup-webhook.sh
sudo ./setup-webhook.sh
# Copy the webhook secret that's displayed

# 7. Configure GitHub webhook
# Go to: GitHub repo → Settings → Webhooks → Add webhook
# URL: http://YOUR_SERVER_IP:9000/webhook
# Secret: (paste from step 6)
# Events: Just the push event

# 8. Initial deployment
chmod +x deploy.sh
./deploy.sh
```

## Daily Usage

### From Local Machine:

```bash
# Make your code changes...

# Deploy with one command
./deploy "Fixed bug in likes feature"

# That's it! Server will auto-deploy via webhook
```

## Environment Variables

Your `.env.production` should have (you already updated these):

```bash
DB_PASSWORD=d5YhNXB5zXhT7bkbbQ7
REDIS_PASSWORD=r5YhNXB5zXhT7bkbbQ7
JWT_SECRET=<generate with: openssl rand -base64 32>
TELEGRAM_BOT_TOKEN=8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=d46ab6ad318b1127d061533769bce800
S3_SECRET_KEY=2f1a730b2b691e8fbc5a33a8595132846cb335c19c523b90a2de173705285c20
```

## Verify Deployment

```bash
# Check backend health
curl https://api.lomi.social/api/v1/health

# Check containers
docker-compose -f docker-compose.prod.yml ps

# Check webhook logs
sudo pm2 logs lomi-webhook
```

## Troubleshooting

### Webhook not working?
```bash
# Check webhook server
sudo pm2 status
sudo pm2 logs lomi-webhook

# Test webhook manually
curl -X POST http://localhost:9000/deploy
```

### Deployment fails?
```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs backend

# Manual deploy
cd /opt/lomi_mini
./deploy.sh
```

## Files Overview

- `deploy` - Local script to push and trigger deployment
- `deploy.sh` - Server script that actually deploys
- `webhook-server.js` - Webhook listener on server
- `setup-webhook.sh` - Sets up webhook server
- `initial-server-setup.sh` - One-time server setup
- `docker-compose.prod.yml` - Production Docker config
- `Caddyfile` - Reverse proxy config

