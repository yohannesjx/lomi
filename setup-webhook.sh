#!/bin/bash

# Setup Webhook Server for GitHub Deployment
# Run this on your server after initial deployment

set -e

echo "🎣 Setting up GitHub webhook server..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "📦 Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

# Install PM2 for process management
if ! command -v pm2 &> /dev/null; then
    echo "📦 Installing PM2..."
    sudo npm install -g pm2
fi

# Create webhook directory
WEBHOOK_DIR="/opt/lomi-webhook"
sudo mkdir -p $WEBHOOK_DIR
sudo cp webhook-server.js $WEBHOOK_DIR/
sudo chmod +x $WEBHOOK_DIR/webhook-server.js

# Create environment file for webhook
WEBHOOK_SECRET=$(openssl rand -hex 32)
sudo tee $WEBHOOK_DIR/.env > /dev/null <<EOF
WEBHOOK_PORT=9000
WEBHOOK_SECRET=$WEBHOOK_SECRET
DEPLOY_PATH=/opt/lomi_mini
EOF

# Create log directory
sudo mkdir -p /var/log
sudo touch /var/log/lomi-webhook.log
sudo chmod 666 /var/log/lomi-webhook.log

# Start webhook server with PM2
cd $WEBHOOK_DIR
sudo pm2 start webhook-server.js --name lomi-webhook --env .env
sudo pm2 save
sudo pm2 startup

echo "✅ Webhook server setup complete!"
echo ""
echo "📋 Webhook Details:"
echo "   URL: http://$(hostname -I | awk '{print $1}'):9000/webhook"
echo "   Secret: $WEBHOOK_SECRET"
echo ""
echo "🔧 Configure in GitHub:"
echo "   1. Go to your repo → Settings → Webhooks"
echo "   2. Add webhook"
echo "   3. Payload URL: http://YOUR_SERVER_IP:9000/webhook"
echo "   4. Content type: application/json"
echo "   5. Secret: $WEBHOOK_SECRET"
echo "   6. Events: Just the push event"
echo ""
echo "📊 View logs: sudo pm2 logs lomi-webhook"
echo "🔄 Restart: sudo pm2 restart lomi-webhook"

