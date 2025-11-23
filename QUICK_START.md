# ⚡ Quick Start - 5 Minutes

## Server IP: 152.53.87.200
## GitHub: https://github.com/yohannesjx/lomi_mini.git

---

## 🖥️ On Your Server (SSH)

### 1. Install Everything
```bash
sudo apt update && sudo apt install -y docker.io docker-compose git curl
sudo systemctl enable docker && sudo systemctl start docker
sudo usermod -aG docker $USER
# Logout and login again
```

### 2. Install Caddy
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install -y caddy
sudo systemctl enable caddy && sudo systemctl start caddy
```

### 3. Clone Repo
```bash
cd /opt
sudo git clone https://github.com/yohannesjx/lomi_mini.git
sudo chown -R $USER:$USER lomi_mini
cd lomi_mini
```

### 4. Create .env.production
```bash
nano .env.production
```

Paste this (generate JWT_SECRET with `openssl rand -base64 32`):
```bash
DB_USER=lomi
DB_PASSWORD=d5YhNXB5zXhT7bkbbQ7
DB_NAME=lomi_db
REDIS_PASSWORD=r5YhNXB5zXhT7bkbbQ7
JWT_SECRET=PASTE_GENERATED_SECRET_HERE
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

Save: `Ctrl+X`, `Y`, `Enter`

### 5. Setup Caddy
```bash
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

### 6. Deploy
```bash
chmod +x deploy.sh
./deploy.sh
```

Wait 2-3 minutes, then test:
```bash
curl http://localhost:8080/api/v1/health
```

---

## 💻 On Your Local Computer

### Deploy Changes
```bash
# Make your code changes...

# Deploy with one command
./deploy "Your changes"
```

That's it! 🎉

---

## 📝 Need Help?

- View logs: `docker-compose -f docker-compose.prod.yml logs -f backend`
- Restart: `docker-compose -f docker-compose.prod.yml restart backend`
- Check status: `docker-compose -f docker-compose.prod.yml ps`

