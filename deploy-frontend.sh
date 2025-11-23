#!/bin/bash

# Deploy Frontend to Server
# Run this LOCALLY (not on server)

set -e

echo "🚀 Deploying Lomi Social Frontend..."

# Check if we're in the right directory
if [ ! -d "frontend" ]; then
    echo "❌ Error: frontend directory not found"
    echo "Run this script from the project root"
    exit 1
fi

# Server details (update these)
SERVER_USER="${SERVER_USER:-root}"
SERVER_HOST="${SERVER_HOST:-152.53.87.200}"
FRONTEND_DIR="/var/www/lomi-frontend"

echo "📦 Building frontend..."
cd frontend

# Check if it's Expo/React Native
if [ -f "package.json" ] && grep -q "expo" package.json; then
    echo "Detected Expo project..."
    
    # Check if expo is installed
    if ! command -v expo &> /dev/null; then
        echo "Installing Expo CLI..."
        npm install -g expo-cli
    fi
    
    # Build for web
    echo "Building Expo web..."
    npx expo export:web || npm run build || npx expo export -p web
    
    BUILD_DIR="web-build"
else
    # Regular React/Next.js build
    echo "Building React app..."
    npm install
    npm run build
    
    # Find build directory
    if [ -d "build" ]; then
        BUILD_DIR="build"
    elif [ -d "dist" ]; then
        BUILD_DIR="dist"
    elif [ -d ".next" ]; then
        BUILD_DIR=".next"
    else
        echo "❌ Error: Could not find build directory"
        exit 1
    fi
fi

if [ ! -d "$BUILD_DIR" ]; then
    echo "❌ Error: Build directory '$BUILD_DIR' not found"
    exit 1
fi

echo "📤 Uploading to server..."
echo "Server: $SERVER_USER@$SERVER_HOST"
echo "Destination: $FRONTEND_DIR"

# Create directory on server
ssh $SERVER_USER@$SERVER_HOST "sudo mkdir -p $FRONTEND_DIR && sudo chown -R $SERVER_USER:$SERVER_USER $FRONTEND_DIR"

# Upload files
rsync -avz --delete "$BUILD_DIR/" $SERVER_USER@$SERVER_HOST:$FRONTEND_DIR/

echo ""
echo "✅ Frontend deployed successfully!"
echo ""
echo "🌐 Access your site:"
echo "   http://$SERVER_HOST"
echo "   http://lomi.social (after DNS is configured)"
echo ""
echo "🔄 Reload Caddy:"
echo "   ssh $SERVER_USER@$SERVER_HOST 'sudo systemctl reload caddy'"

