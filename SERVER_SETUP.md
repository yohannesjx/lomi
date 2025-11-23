# Server Setup Guide for Lomi Social

This guide will help you set up your production server with Caddy, Docker, and GitHub deployment.

## Prerequisites

- Ubuntu 20.04+ or Debian 11+ server
- Domain names pointing to your server:
  - `api.lomi.social` → Backend API
  - `lomi.social` → Frontend (optional, if hosting static files)
- SSH access to server
- Root or sudo access

## Step 1: Initial Server Setup

### Update system
```bash
sudo apt update && sudo apt upgrade -y
```

### Install required packages
```bash
sudo apt install -y \
    git \
    curl \
    wget \
    docker.io \
    docker-compose \
    ufw \
    fail2ban
```

### Start Docker
```bash
sudo systemctl enable docker
sudo systemctl start docker
```

### Add user to docker group (if not root)
```bash
sudo usermod -aG docker $USER
newgrp docker  # Or logout and login again
```

## Step 2: Install Caddy

### Add Caddy repository
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
```

### Install Caddy
```bash
sudo apt update
sudo apt install -y caddy
```

### Enable and start Caddy
```bash
sudo systemctl enable caddy
sudo systemctl start caddy
```

## Step 3: Clone Repository

```bash
cd /opt  # Or your preferred directory
git clone https://github.com/YOUR_USERNAME/lomi_mini.git
cd lomi_mini
```

## Step 4: Configure Environment Variables

Create `.env.production` file:

```bash
cp .env.example .env.production  # If you have an example file
nano .env.production
```

Add all required variables (see `.env.production.example` below):

```bash
# Database
DB_USER=lomi
DB_PASSWORD=your_secure_password_here
DB_NAME=lomi_db

# Redis
REDIS_PASSWORD=your_secure_redis_password_here

# JWT
JWT_SECRET=your_super_secret_jwt_key_min_32_chars_long
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Telegram
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_BOT_USERNAME=lomi_social_bot

# Cloudflare R2
S3_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
S3_ACCESS_KEY=your_r2_access_key
S3_SECRET_KEY=your_r2_secret_key
S3_USE_SSL=true
S3_REGION=auto
S3_BUCKET_PHOTOS=lomi-photos
S3_BUCKET_VIDEOS=lomi-videos
S3_BUCKET_GIFTS=lomi-gifts
S3_BUCKET_VERIFICATIONS=lomi-verifications

# Payment Gateways (optional)
TELEBIRR_API_KEY=
CBE_BIRR_API_KEY=
HELLOCASH_API_KEY=
AMOLE_API_KEY=

# Platform Settings
PLATFORM_FEE_PERCENTAGE=25
MIN_PAYOUT_AMOUNT=1000
COIN_TO_BIRR_RATE=0.10

# Notifications (optional)
ONESIGNAL_APP_ID=
ONESIGNAL_API_KEY=
FIREBASE_SERVER_KEY=
```

**Important:** Never commit `.env.production` to GitHub!

## Step 5: Configure Caddy

### Copy Caddyfile
```bash
sudo cp Caddyfile /etc/caddy/Caddyfile
```

### Edit Caddyfile (if needed)
```bash
sudo nano /etc/caddy/Caddyfile
```

Update domain names to match your setup.

### Test Caddy configuration
```bash
sudo caddy validate --config /etc/caddy/Caddyfile
```

### Reload Caddy
```bash
sudo systemctl reload caddy
```

## Step 6: Make Deploy Script Executable

```bash
chmod +x deploy.sh
```

## Step 7: Initial Deployment

```bash
./deploy.sh
```

This will:
- Pull latest code
- Build Docker images
- Start all services
- Check health
- Reload Caddy

## Step 8: Configure Firewall

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS (Caddy handles this)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

## Step 9: Set Up GitHub Actions (Optional)

### On GitHub:
1. Go to your repository → Settings → Secrets and variables → Actions
2. Add these secrets:
   - `SERVER_HOST`: Your server IP or domain
   - `SERVER_USER`: SSH username (usually `root` or `ubuntu`)
   - `SERVER_SSH_KEY`: Your private SSH key

### Generate SSH key (if needed)
```bash
ssh-keygen -t ed25519 -C "github-actions"
# Copy public key to server
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@your-server
```

### Update workflow file
Edit `.github/workflows/deploy.yml` and update the path:
```yaml
script: |
  cd /opt/lomi_mini  # Update to your actual path
  ./deploy.sh
```

## Step 10: Manual Deployment (Alternative)

If not using GitHub Actions, deploy manually:

```bash
# SSH into server
ssh user@your-server

# Navigate to project
cd /opt/lomi_mini

# Pull latest code
git pull origin main

# Run deploy script
./deploy.sh
```

## Step 11: Verify Deployment

### Check backend
```bash
curl https://api.lomi.social/api/v1/health
```

Should return:
```json
{"status":"ok","message":"Lomi Backend is running 🍋"}
```

### Check containers
```bash
docker-compose -f docker-compose.prod.yml ps
```

### Check logs
```bash
docker-compose -f docker-compose.prod.yml logs -f backend
```

## Step 12: Configure Telegram Mini App

1. Go to [@BotFather](https://t.me/BotFather)
2. Send `/newapp`
3. Select your bot: `lomi_social_bot`
4. Provide:
   - **Title**: Lomi Social
   - **Short name**: lomi
   - **Description**: Find your Lomi in Ethiopia
   - **Photo**: Upload your app icon
   - **Web App URL**: `https://lomi.social` (or your frontend URL)
   - **Short name**: lomi

## Maintenance

### View logs
```bash
# Backend logs
docker-compose -f docker-compose.prod.yml logs -f backend

# All logs
docker-compose -f docker-compose.prod.yml logs -f

# Caddy logs
sudo journalctl -u caddy -f
```

### Restart services
```bash
docker-compose -f docker-compose.prod.yml restart backend
sudo systemctl restart caddy
```

### Update code
```bash
git pull origin main
./deploy.sh
```

### Backup database
```bash
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U lomi lomi_db > backup_$(date +%Y%m%d).sql
```

## Troubleshooting

### Caddy not starting
```bash
sudo caddy validate --config /etc/caddy/Caddyfile
sudo journalctl -u caddy -n 50
```

### Backend not responding
```bash
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml ps
```

### SSL certificate issues
Caddy automatically manages SSL certificates. If issues occur:
```bash
sudo caddy reload --config /etc/caddy/Caddyfile
```

## Security Checklist

- [ ] Firewall configured (UFW)
- [ ] Fail2ban installed and configured
- [ ] Strong passwords in `.env.production`
- [ ] `.env.production` not in Git
- [ ] SSH key authentication (disable password)
- [ ] Regular system updates
- [ ] Database backups scheduled
- [ ] Caddy automatically handles HTTPS

## Why Caddy Separate (Not in Docker)?

✅ **Better SSL Management**: Caddy needs direct access to ports 80/443 for Let's Encrypt  
✅ **Easier Updates**: Update Caddy independently without rebuilding containers  
✅ **Standard Practice**: Reverse proxies are typically on the host  
✅ **Performance**: No Docker network overhead for SSL termination  
✅ **Logging**: Easier access to Caddy logs on the host  

## Next Steps

1. Set up monitoring (optional): Prometheus, Grafana, or simple health checks
2. Set up backups: Automated database backups
3. Set up monitoring alerts: Get notified of downtime
4. Configure CDN: For frontend static assets (Cloudflare, etc.)

