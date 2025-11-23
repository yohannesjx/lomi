# Server Commands Cheat Sheet

## Initial Setup (Run Once on Server)

```bash
# 1. Clone repo
cd /opt
git clone https://github.com/YOUR_USERNAME/lomi_mini.git
cd lomi_mini

# 2. Run initial setup
sudo chmod +x initial-server-setup.sh
sudo ./initial-server-setup.sh

# 3. Create .env.production with your values
nano .env.production
# (Paste your environment variables)

# 4. Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET"
# Add this to .env.production

# 5. Setup Caddy
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo nano /etc/caddy/Caddyfile  # Update domains
sudo systemctl reload caddy

# 6. Setup webhook
sudo chmod +x setup-webhook.sh
sudo ./setup-webhook.sh
# Copy the webhook secret shown

# 7. Configure GitHub webhook
# Go to: GitHub → Settings → Webhooks → Add webhook
# URL: http://YOUR_SERVER_IP:9000/webhook
# Secret: (from step 6)

# 8. Deploy
chmod +x deploy.sh
./deploy.sh
```

## Daily Operations

### Check Status
```bash
# Containers
docker-compose -f docker-compose.prod.yml ps

# Backend health
curl http://localhost:8080/api/v1/health

# Webhook logs
sudo pm2 logs lomi-webhook

# All logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Restart Services
```bash
# Restart backend
docker-compose -f docker-compose.prod.yml restart backend

# Restart webhook
sudo pm2 restart lomi-webhook

# Restart Caddy
sudo systemctl restart caddy
```

### Manual Deployment
```bash
cd /opt/lomi_mini
git pull origin main
./deploy.sh
```

### Trigger Webhook Manually
```bash
curl -X POST http://localhost:9000/deploy
```

## Troubleshooting

### View Logs
```bash
# Backend logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Webhook logs
sudo pm2 logs lomi-webhook
tail -f /var/log/lomi-webhook.log

# Caddy logs
sudo journalctl -u caddy -f
```

### Check Services
```bash
# Docker containers
docker ps
docker-compose -f docker-compose.prod.yml ps

# PM2 processes
sudo pm2 status
sudo pm2 list

# System services
sudo systemctl status caddy
sudo systemctl status docker
```

### Fix Issues
```bash
# Rebuild containers
docker-compose -f docker-compose.prod.yml build --no-cache backend
docker-compose -f docker-compose.prod.yml up -d

# Restart everything
docker-compose -f docker-compose.prod.yml restart
sudo systemctl restart caddy
sudo pm2 restart lomi-webhook
```

