#!/bin/bash

# Quick Caddy fix script
# Run this ON THE SERVER if Caddy is stuck

echo "🔧 Caddy Emergency Fix"
echo "======================"
echo ""

# 1. Stop the stuck reload
echo "1️⃣ Stopping Caddy..."
sudo systemctl stop caddy
sleep 2

# 2. Check Caddyfile syntax
echo "2️⃣ Validating Caddyfile..."
if sudo caddy validate --config /etc/caddy/Caddyfile; then
    echo "✅ Caddyfile is valid"
else
    echo "❌ Caddyfile has errors!"
    echo ""
    echo "Restoring backup..."
    if [ -f "/etc/caddy/Caddyfile.backup" ]; then
        sudo cp /etc/caddy/Caddyfile.backup /etc/caddy/Caddyfile
        echo "✅ Restored backup"
    else
        echo "⚠️  No backup found, using minimal config..."
        sudo tee /etc/caddy/Caddyfile > /dev/null << 'EOF'
:80 {
    handle /api/* {
        reverse_proxy localhost:8080
    }
    
    handle {
        root * /var/www/lomi-frontend
        try_files {path} /index.html
        file_server
    }
}
EOF
        echo "✅ Created minimal working config"
    fi
fi

# 3. Start Caddy
echo ""
echo "3️⃣ Starting Caddy..."
sudo systemctl start caddy
sleep 2

# 4. Check status
echo ""
echo "4️⃣ Checking Caddy status..."
if sudo systemctl is-active --quiet caddy; then
    echo "✅ Caddy is running"
else
    echo "❌ Caddy failed to start!"
    echo ""
    echo "Checking logs..."
    sudo journalctl -u caddy -n 50 --no-pager
fi

echo ""
echo "5️⃣ Testing API..."
if curl -s http://localhost/api/v1/health > /dev/null; then
    echo "✅ API is accessible"
else
    echo "❌ API not accessible"
fi

echo ""
echo "Done! Caddy should be running now."
