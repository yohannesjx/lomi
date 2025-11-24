#!/bin/bash

# Check CompreFace container status and connectivity

set -e

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

echo "🔍 Checking CompreFace Status"
echo "=============================="
echo ""

# Check if container exists and is running
echo "1. Checking Docker container status..."
CONTAINER_STATUS=$(docker-compose -f docker-compose.prod.yml --env-file .env.production ps compreface 2>/dev/null | grep -v "NAME" | awk '{print $7}' || echo "NOT_FOUND")

if [ "$CONTAINER_STATUS" = "Up" ] || [ "$CONTAINER_STATUS" = "running" ]; then
    echo "✅ CompreFace container is running"
else
    echo "❌ CompreFace container is NOT running (status: $CONTAINER_STATUS)"
    echo ""
    echo "💡 Try starting it:"
    echo "   docker-compose -f docker-compose.prod.yml --env-file .env.production up -d compreface"
    exit 1
fi

echo ""

# Check container logs
echo "2. Checking recent CompreFace logs..."
echo "   (Last 20 lines)"
docker-compose -f docker-compose.prod.yml --env-file .env.production logs --tail=20 compreface 2>/dev/null || echo "Could not fetch logs"

echo ""
echo ""

# Test from host (localhost)
echo "3. Testing from host (localhost:8000)..."
HOST_TEST=$(curl -s -w "\nHTTP_CODE:%{http_code}" --connect-timeout 5 "http://localhost:8000/api/v1/health" 2>&1 || echo "FAILED")
HOST_HTTP=$(echo "$HOST_TEST" | grep "HTTP_CODE:" | cut -d: -f2 || echo "000")
HOST_BODY=$(echo "$HOST_TEST" | grep -v "HTTP_CODE:" | head -1)

if [ "$HOST_HTTP" = "200" ]; then
    echo "✅ Host can reach CompreFace (localhost:8000)"
else
    echo "⚠️  Host cannot reach CompreFace (HTTP $HOST_HTTP)"
    echo "   This is OK if workers connect via Docker network"
fi

echo ""

# Test from within Docker network (via worker container)
echo "4. Testing from Docker network (compreface:8000)..."
NETWORK_TEST=$(docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T moderator-worker curl -s -w "\nHTTP_CODE:%{http_code}" --connect-timeout 5 "http://compreface:8000/api/v1/health" 2>&1 || echo "FAILED")
NETWORK_HTTP=$(echo "$NETWORK_TEST" | grep "HTTP_CODE:" | cut -d: -f2 || echo "000")
NETWORK_BODY=$(echo "$NETWORK_TEST" | grep -v "HTTP_CODE:" | head -1)

if [ "$NETWORK_HTTP" = "200" ]; then
    echo "✅ Workers can reach CompreFace (compreface:8000)"
    echo "   Response: $NETWORK_BODY"
else
    echo "❌ Workers CANNOT reach CompreFace (HTTP $NETWORK_HTTP)"
    echo "   Error: $NETWORK_BODY"
    echo ""
    echo "💡 Possible issues:"
    echo "   - CompreFace not on same Docker network"
    echo "   - CompreFace not listening on port 8000"
    echo "   - Firewall blocking connection"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Summary:"
echo "   Container: $([ "$CONTAINER_STATUS" = "Up" ] && echo "✅ Running" || echo "❌ Not running")"
echo "   Host access: $([ "$HOST_HTTP" = "200" ] && echo "✅ OK" || echo "⚠️  Limited (expected)")"
echo "   Network access: $([ "$NETWORK_HTTP" = "200" ] && echo "✅ OK" || echo "❌ FAILED")"
echo ""
echo "💡 If network access failed, check:"
echo "   1. docker-compose -f docker-compose.prod.yml --env-file .env.production ps"
echo "   2. docker network ls"
echo "   3. docker-compose -f docker-compose.prod.yml --env-file .env.production logs compreface"
echo ""

