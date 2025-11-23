#!/bin/bash

# Restore Lomi Mini from GitHub to /opt/lomi_mini
# Run this on your server if files were deleted

set -e

echo "🔄 Restoring Lomi Mini from GitHub..."

# Navigate to /opt
cd /opt

# Remove old directory if it exists (empty or corrupted)
if [ -d "lomi_mini" ]; then
    echo "⚠️  Removing old /opt/lomi_mini directory..."
    rm -rf lomi_mini
fi

# Clone fresh from GitHub
echo "📥 Cloning repository from GitHub..."
git clone https://github.com/yohannesjx/lomi_mini.git

cd lomi_mini

echo "✅ Repository cloned successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Copy your .env.production file back (if you have a backup)"
echo "2. Run: cd /opt/lomi_mini && ./deploy-all.sh"
echo ""
echo "📁 Files restored to: /opt/lomi_mini"
ls -la | head -20

