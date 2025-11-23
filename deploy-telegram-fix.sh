#!/bin/bash

# Quick deployment script for Telegram login fix
# Run this on your server

set -e

echo "🚀 Deploying Telegram Login Fix..."

# 1. Stop services
echo "📦 Stopping services..."
docker-compose -f docker-compose.prod.yml down

# 2. Rebuild backend
echo "🔨 Building backend..."
docker-compose -f docker-compose.prod.yml build backend

# 3. Start services
echo "▶️  Starting services..."
docker-compose -f docker-compose.prod.yml up -d

# 4. Wait for backend to be ready
echo "⏳ Waiting for backend to start..."
sleep 5

# 5. Test backend health
echo "🏥 Testing backend health..."
curl -f http://localhost:8080/api/v1/health || {
    echo "❌ Backend health check failed!"
    docker-compose -f docker-compose.prod.yml logs backend
    exit 1
}

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📋 Next steps:"
echo "1. Rebuild frontend: cd frontend && npm run build"
echo "2. Deploy frontend: scp -r dist/* user@server:/var/www/lomi-frontend/"
echo "3. Test in Telegram app"
echo ""
echo "📊 Check logs:"
echo "   docker-compose -f docker-compose.prod.yml logs backend -f"
