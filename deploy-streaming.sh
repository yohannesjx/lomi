#!/bin/bash

# Deploy TikTok Streaming Backend to Production Server
# This script deploys the new streaming endpoints to your VPS

set -e

echo "🚀 Deploying TikTok Streaming Backend"
echo "======================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
SERVER_USER="${SERVER_USER:-root}"
SERVER_HOST="${SERVER_HOST:-your-server-ip}"
SERVER_PATH="${SERVER_PATH:-/opt/lomi-backend}"

echo -e "${YELLOW}📋 Deployment Configuration:${NC}"
echo "Server: $SERVER_USER@$SERVER_HOST"
echo "Path: $SERVER_PATH"
echo ""

# Step 1: Git add and commit
echo -e "${YELLOW}1️⃣  Committing changes to git...${NC}"
git add backend/internal/handlers/streaming.go
git add backend/internal/routes/streaming_routes.go
git add backend/cmd/api/main.go
git add TIKTOK_*.md
git add DEPLOYMENT_CHECKLIST.md
git add test-tiktok-api.sh

git commit -m "feat: Add TikTok-style streaming endpoints

- Implement 6 critical endpoints (registerUser, showUserDetail, showRelatedVideos, liveStream, sendGift, purchaseCoin)
- Add streaming handler with exact TikTok API response format
- Integrate with existing User/Gift/Wallet models
- Add comprehensive documentation and test scripts
- Ready for production deployment"

echo -e "${GREEN}✅ Changes committed${NC}"
echo ""

# Step 2: Push to origin
echo -e "${YELLOW}2️⃣  Pushing to origin...${NC}"
git push origin main
echo -e "${GREEN}✅ Pushed to origin${NC}"
echo ""

# Step 3: Deploy to server
echo -e "${YELLOW}3️⃣  Deploying to server...${NC}"
echo "Connecting to $SERVER_HOST..."

ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
set -e

echo "📥 Pulling latest code..."
cd /opt/lomi-backend
git pull origin main

echo "🔨 Building backend..."
cd backend
go build -o lomi-backend cmd/api/main.go

echo "🔄 Restarting backend service..."
sudo systemctl restart lomi-backend

echo "✅ Deployment complete!"
echo ""
echo "📊 Service status:"
sudo systemctl status lomi-backend --no-pager -l

echo ""
echo "📝 Recent logs:"
sudo journalctl -u lomi-backend -n 20 --no-pager
ENDSSH

echo ""
echo -e "${GREEN}🎉 Deployment successful!${NC}"
echo ""
echo "Next steps:"
echo "1. Test endpoints: curl https://api.lomilive.com/api/registerUser"
echo "2. Update Android/iOS apps with production URL"
echo "3. Build and test apps"
echo ""
