#!/bin/bash

# Diagnostic script to check Telegram login setup
# Run this ON THE SERVER

echo "🔍 Telegram Login Diagnostic"
echo "=============================="
echo ""

# 1. Check backend is running
echo "1️⃣ Checking backend..."
if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ Backend is running"
    curl -s http://localhost:8080/api/v1/health | jq . || curl -s http://localhost:8080/api/v1/health
else
    echo "❌ Backend is NOT running!"
fi
echo ""

# 2. Test auth endpoint directly (should get 401, not 405)
echo "2️⃣ Testing auth endpoint (expecting 401, NOT 405)..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -H "Authorization: tma test_data")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "HTTP Status: $http_code"
echo "Response: $body"

if [ "$http_code" = "405" ]; then
    echo "❌ ERROR: Getting 405 (Method Not Allowed)"
    echo "   This means the route is not configured correctly"
elif [ "$http_code" = "401" ]; then
    echo "✅ CORRECT: Getting 401 (Unauthorized)"
    echo "   Route is working, just needs valid initData"
else
    echo "⚠️  Unexpected status: $http_code"
fi
echo ""

# 3. Check Caddy is proxying correctly
echo "3️⃣ Testing through Caddy..."
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost/api/v1/auth/telegram \
  -H "Content-Type: application/json" \
  -H "Authorization: tma test_data")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo "HTTP Status: $http_code"
echo "Response: $body"

if [ "$http_code" = "405" ]; then
    echo "❌ ERROR: Caddy returning 405"
    echo "   Caddy might be serving static files instead of proxying"
elif [ "$http_code" = "401" ]; then
    echo "✅ CORRECT: Caddy is proxying correctly"
else
    echo "⚠️  Unexpected status: $http_code"
fi
echo ""

# 4. Check Caddy configuration
echo "4️⃣ Checking Caddy configuration..."
if [ -f "/etc/caddy/Caddyfile" ]; then
    echo "Caddyfile exists"
    if sudo caddy validate --config /etc/caddy/Caddyfile 2>/dev/null; then
        echo "✅ Caddyfile is valid"
    else
        echo "❌ Caddyfile has errors:"
        sudo caddy validate --config /etc/caddy/Caddyfile
    fi
    
    echo ""
    echo "Checking for handle_path directive..."
    if grep -q "handle_path /api/\*" /etc/caddy/Caddyfile; then
        echo "✅ Found handle_path /api/* (CORRECT)"
    elif grep -q "handle /api/\*" /etc/caddy/Caddyfile; then
        echo "⚠️  Found handle /api/* (should be handle_path)"
    else
        echo "❌ No API handler found!"
    fi
else
    echo "❌ Caddyfile not found!"
fi
echo ""

# 5. Check frontend build
echo "5️⃣ Checking frontend..."
if [ -d "/var/www/lomi-frontend" ]; then
    echo "Frontend directory exists"
    file_count=$(find /var/www/lomi-frontend -type f | wc -l)
    echo "Files in frontend: $file_count"
    
    if [ -f "/var/www/lomi-frontend/index.html" ]; then
        echo "✅ index.html exists"
        
        # Check if it's a recent build
        mod_time=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" /var/www/lomi-frontend/index.html 2>/dev/null || stat -c "%y" /var/www/lomi-frontend/index.html 2>/dev/null)
        echo "Last modified: $mod_time"
    else
        echo "❌ index.html NOT found!"
    fi
else
    echo "❌ Frontend directory NOT found!"
fi
echo ""

# 6. Check Docker containers
echo "6️⃣ Checking Docker containers..."
docker-compose -f docker-compose.prod.yml ps
echo ""

# 7. Check backend logs for errors
echo "7️⃣ Recent backend logs (last 20 lines)..."
docker-compose -f docker-compose.prod.yml logs backend --tail=20
echo ""

# 8. Check Caddy logs
echo "8️⃣ Recent Caddy logs..."
sudo journalctl -u caddy -n 20 --no-pager
echo ""

# Summary
echo "=============================="
echo "📊 SUMMARY"
echo "=============================="
echo ""
echo "Next steps:"
echo "1. If backend shows 401 but Caddy shows 405:"
echo "   → Caddy is not proxying correctly"
echo "   → Update Caddyfile and reload: sudo systemctl reload caddy"
echo ""
echo "2. If both show 405:"
echo "   → Backend route not configured"
echo "   → Rebuild backend: docker-compose -f docker-compose.prod.yml build backend --no-cache"
echo ""
echo "3. If frontend is old:"
echo "   → Rebuild frontend: cd frontend && npx expo export -p web"
echo "   → Deploy: sudo cp -r dist/* /var/www/lomi-frontend/"
