#!/bin/bash

# Diagnose CompreFace issues

set -e

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

echo "🔍 CompreFace Diagnosis"
echo "======================"
echo ""

# Check container status
echo "1. Container Status:"
docker-compose -f docker-compose.prod.yml --env-file .env.production ps compreface
echo ""

# Check if it's on the network
echo "2. Network Check:"
NETWORK_NAME=$(docker-compose -f docker-compose.prod.yml --env-file .env.production config | grep -A 5 "networks:" | grep -v "networks:" | head -1 | awk '{print $1}' | tr -d ':')
if [ -z "$NETWORK_NAME" ]; then
    NETWORK_NAME="lomi_mini_lomi_network"
fi
echo "   Network: $NETWORK_NAME"
docker network inspect "$NETWORK_NAME" 2>/dev/null | grep -A 3 "compreface" || echo "   ⚠️  CompreFace not found in network"
echo ""

# Check recent logs
echo "3. Recent Logs (last 30 lines):"
docker-compose -f docker-compose.prod.yml --env-file .env.production logs --tail=30 compreface
echo ""

# Check if port is listening inside container
echo "4. Checking if CompreFace is listening inside container:"
LISTENING=$(docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T compreface netstat -tlnp 2>/dev/null | grep ":8000" || docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T compreface ss -tlnp 2>/dev/null | grep ":8000" || echo "Could not check")
if echo "$LISTENING" | grep -q ":8000"; then
    echo "   ✅ CompreFace is listening on port 8000"
    echo "   $LISTENING"
else
    echo "   ❌ CompreFace is NOT listening on port 8000"
    echo "   This means CompreFace hasn't started its web server"
fi
echo ""

# Check CompreFace process
echo "5. CompreFace processes:"
docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T compreface ps aux 2>/dev/null | head -10 || echo "   Could not check processes"
echo ""

# Try to access from inside CompreFace container itself
echo "6. Testing from inside CompreFace container:"
INTERNAL_TEST=$(docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T compreface curl -s http://localhost:8000/api/v1/health 2>&1 || echo "FAILED")
if echo "$INTERNAL_TEST" | grep -q "FAILED\|Connection refused"; then
    echo "   ❌ CompreFace cannot reach itself on localhost:8000"
    echo "   This means CompreFace web server is not running"
else
    echo "   ✅ CompreFace can reach itself"
    echo "   Response: $INTERNAL_TEST"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 Based on the diagnosis above:"
echo ""
echo "If CompreFace is not listening:"
echo "  1. CompreFace might still be starting (can take 2-3 minutes)"
echo "  2. CompreFace might need additional environment variables"
echo "  3. CompreFace image might need a different configuration"
echo ""
echo "Try:"
echo "  docker-compose -f docker-compose.prod.yml --env-file .env.production logs -f compreface"
echo "  (Watch for 'started' or 'listening' messages)"
echo ""

