#!/bin/bash

# Test Backend and Caddy
# Run this on your server

echo "🔍 Testing Backend and Caddy..."

echo ""
echo "1️⃣ Checking Docker containers..."
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "2️⃣ Testing backend directly (port 8080)..."
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health | grep -q "200\|404"; then
    echo "✅ Backend is responding on port 8080"
    curl http://localhost:8080/api/v1/health
else
    echo "❌ Backend is NOT responding on port 8080"
    echo "Checking backend logs..."
    docker-compose -f docker-compose.prod.yml logs backend --tail 20
fi

echo ""
echo "3️⃣ Testing via Caddy (IP address)..."
if curl -s -o /dev/null -w "%{http_code}" http://152.53.87.200/api/v1/health | grep -q "200\|404"; then
    echo "✅ Caddy is proxying correctly"
    curl http://152.53.87.200/api/v1/health
else
    echo "❌ Caddy cannot reach backend"
    echo "Checking Caddy status..."
    sudo systemctl status caddy --no-pager -l | head -10
fi

echo ""
echo "4️⃣ Checking what's listening on port 8080..."
sudo netstat -tulpn | grep 8080 || sudo ss -tulpn | grep 8080

echo ""
echo "5️⃣ Backend container logs (last 10 lines)..."
docker-compose -f docker-compose.prod.yml logs backend --tail 10

echo ""
echo "✅ Testing complete!"

