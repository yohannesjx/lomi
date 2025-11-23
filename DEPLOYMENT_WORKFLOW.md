# One-Command Deployment Workflow

## Setup (One Time)

### 1. On Your Server

```bash
# Navigate to project
cd /opt/lomi_mini

# Create .env.production (already done, but verify)
cat .env.production

# Set up webhook server
chmod +x setup-webhook.sh
sudo ./setup-webhook.sh
```

This will:
- Install Node.js and PM2
- Set up webhook server on port 9000
- Generate a webhook secret
- Start the webhook service

### 2. Configure GitHub Webhook

1. Go to your GitHub repo → **Settings** → **Webhooks**
2. Click **Add webhook**
3. Fill in:
   - **Payload URL**: `http://YOUR_SERVER_IP:9000/webhook`
   - **Content type**: `application/json`
   - **Secret**: (Copy from server output after running setup-webhook.sh)
   - **Events**: Select "Just the push event"
4. Click **Add webhook**

### 3. Update Caddyfile (if needed)

Make sure Caddyfile allows webhook port:
```bash
sudo nano /etc/caddy/Caddyfile
```

Add this if webhook needs HTTPS:
```
webhook.lomi.social {
    reverse_proxy localhost:9000
}
```

## Daily Usage

### From Your Local Machine:

```bash
# Make changes to code...

# Deploy with one command
./deploy "Your commit message here"

# Or just
./deploy
```

That's it! The script will:
1. ✅ Stage all changes
2. ✅ Commit with your message
3. ✅ Push to GitHub
4. ✅ GitHub webhook triggers server deployment automatically

## Manual Deployment (If Needed)

### On Server:
```bash
cd /opt/lomi_mini
./deploy.sh
```

### Trigger Webhook Manually:
```bash
curl -X POST http://localhost:9000/deploy
```

## Check Deployment Status

### On Server:
```bash
# Check webhook logs
sudo pm2 logs lomi-webhook

# Check deployment logs
tail -f /var/log/lomi-webhook.log

# Check container status
docker-compose -f docker-compose.prod.yml ps

# Check backend health
curl http://localhost:8080/api/v1/health
```

## Troubleshooting

### Webhook not triggering?
```bash
# Check if webhook server is running
sudo pm2 status

# Restart webhook
sudo pm2 restart lomi-webhook

# Check webhook logs
sudo pm2 logs lomi-webhook
```

### GitHub webhook not receiving?
- Check GitHub webhook delivery logs in repo settings
- Verify server IP is accessible
- Check firewall allows port 9000

### Deployment fails?
```bash
# Check deployment script logs
cd /opt/lomi_mini
./deploy.sh

# Check container logs
docker-compose -f docker-compose.prod.yml logs backend
```

