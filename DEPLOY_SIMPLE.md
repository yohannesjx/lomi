# 🚀 Simple Deployment Guide - Step by Step

## Your Server IP: 152.53.87.200
## Your GitHub: https://github.com/yohannesjx/lomi_mini.git

---

## STEP 1: Initial Server Setup (Run Once)

SSH into your server and run:

```bash
# Install everything needed
sudo apt update
sudo apt install -y docker.io docker-compose git curl

# Install Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install -y caddy
sudo systemctl enable caddy && sudo systemctl start caddy

# Install Node.js (for webhook)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
sudo npm install -g pm2

# Start Docker
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER
```

**Logout and login again** (or run `newgrp docker`) for Docker group to take effect.

---

## STEP 2: Clone Repository

```bash
cd /opt
sudo git clone https://github.com/yohannesjx/lomi_mini.git
sudo chown -R $USER:$USER lomi_mini
cd lomi_mini
```

---

## STEP 3: Create Environment File

```bash
nano .env.production
```

Paste this (update JWT_SECRET with a random string):

```bash
DB_USER=lomi
DB_PASSWORD=d5YhNXB5zXhT7bkbbQ7
DB_NAME=lomi_db
REDIS_PASSWORD=r5YhNXB5zXhT7bkbbQ7
JWT_SECRET=q9cN7w2Lk1xV4pF8eR0tS3zH6mJbYaUdPiGoTnW5Cs0
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
TELEGRAM_BOT_TOKEN=8453633918:AAE6UxkHrplAxyKXXBLt56bQufhZpH-rVEM
TELEGRAM_BOT_USERNAME=lomi_social_bot
S3_ENDPOINT=https://a53cdfc7c678dac2a028159bcd178da2.r2.cloudflarestorage.com
S3_ACCESS_KEY=d46ab6ad318b1127d061533769bce800
S3_SECRET_KEY=2f1a730b2b691e8fbc5a33a8595132846cb335c19c523b90a2de173705285c20
S3_USE_SSL=true
S3_REGION=auto
S3_BUCKET_PHOTOS=lomi-photos
S3_BUCKET_VIDEOS=lomi-videos
S3_BUCKET_GIFTS=lomi-gifts
S3_BUCKET_VERIFICATIONS=lomi-verifications
```

Generate JWT secret:
```bash
openssl rand -base64 32
```
Copy the output and replace `CHANGE_THIS_GENERATE_WITH_openssl_rand_base64_32` in `.env.production`

Save: `Ctrl+X`, then `Y`, then `Enter`

---

## STEP 4: Setup Caddy

```bash
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

---

## STEP 5: First Deployment

```bash
chmod +x deploy.sh
./deploy.sh
```

Wait 2-3 minutes for everything to build and start.

---

## STEP 6: Verify It Works

```bash
# Check if backend is running
curl http://localhost:8080/api/v1/health

# Should return: {"status":"ok","message":"Lomi Backend is running 🍋"}
```

---

## ✅ DONE! Your backend is live.

---

## 🔄 Future Deployments (Super Simple)

### Option A: From Your Local Computer

```bash
# Make your code changes...

# Deploy with one command
./deploy "Fixed bug"
```

That's it! It will:
1. Commit changes
2. Push to GitHub
3. Server auto-deploys (if webhook is set up)

### Option B: From Server (Manual)

```bash
ssh user@152.53.87.200
cd /opt/lomi_mini
git pull origin main
./deploy.sh
```

---

## 🎣 Optional: Setup Webhook (Auto-Deploy on Push)

### On Server:

```bash
cd /opt/lomi_mini
chmod +x setup-webhook.sh
sudo ./setup-webhook.sh
```

Copy the **webhook secret** that's displayed.

### On GitHub:

1. Go to: https://github.com/yohannesjx/lomi_mini/settings/hooks
2. Click **Add webhook**
3. Fill in:
   - **Payload URL**: `http://152.53.87.200:9000/webhook`
   - **Content type**: `application/json`
   - **Secret**: (paste from server)
   - **Events**: Just the push event
4. Click **Add webhook**

Now every `git push` will auto-deploy!

---

## 📋 Quick Commands Reference

```bash
# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Restart backend
docker-compose -f docker-compose.prod.yml restart backend

# Check health
curl http://localhost:8080/api/v1/health
```

---

## 🐛 Troubleshooting

### Backend not starting?
```bash
cd /opt/lomi_mini
docker-compose -f docker-compose.prod.yml logs backend
```

### Caddy not working?
```bash
sudo systemctl status caddy
sudo journalctl -u caddy -n 50
```

### Need to rebuild?
```bash
cd /opt/lomi_mini
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml build --no-cache backend
docker-compose -f docker-compose.prod.yml up -d
```

---

That's it! Simple and clear. 🎉

