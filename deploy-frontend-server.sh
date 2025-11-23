#!/bin/bash

# Deploy Frontend from Server
# Run this ON THE SERVER after building

set -e

echo "🚀 Deploying Frontend from Server..."

FRONTEND_DIR="/var/www/lomi-frontend"
PROJECT_DIR="${PROJECT_DIR:-~/lomi_mini}"

cd "$PROJECT_DIR/frontend"

# Check if dist directory exists
if [ ! -d "dist" ]; then
    echo "❌ Error: dist directory not found"
    echo "Run 'npx expo export -p web' first"
    exit 1
fi

echo "📤 Copying files to $FRONTEND_DIR..."

# Create directory
sudo mkdir -p $FRONTEND_DIR

# Copy files
sudo cp -r dist/* $FRONTEND_DIR/

# Set permissions
sudo chown -R www-data:www-data $FRONTEND_DIR
sudo chmod -R 755 $FRONTEND_DIR

echo "✅ Frontend deployed!"
echo ""
echo "Files in $FRONTEND_DIR:"
ls -la $FRONTEND_DIR | head -10

echo ""
echo "🔄 Reloading Caddy..."
sudo systemctl reload caddy

echo ""
echo "✅ Done! Test at: http://152.53.87.200"

