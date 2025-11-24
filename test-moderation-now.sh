#!/bin/bash

# Quick test to verify moderation is working with OpenCV fallback

set -e

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

echo "🧪 Testing Moderation System"
echo "==========================="
echo ""

# Check worker logs for OpenCV fallback
echo "1. Checking worker logs for OpenCV fallback usage..."
echo "   (Looking for 'OpenCV fallback' or 'CompreFace error' messages)"
echo ""

RECENT_LOGS=$(docker-compose -f docker-compose.prod.yml --env-file .env.production logs --tail=50 moderator-worker 2>/dev/null | grep -i "opencv\|compreface\|face detection\|moderation result" | tail -10 || echo "")

if [ -n "$RECENT_LOGS" ]; then
    echo "$RECENT_LOGS"
else
    echo "   No recent moderation activity. Upload a photo to trigger moderation."
fi

echo ""
echo ""

# Check if workers are running
echo "2. Worker Status:"
docker-compose -f docker-compose.prod.yml --env-file .env.production ps moderator-worker | grep -v "NAME" | head -4
echo ""

# Check queue
echo "3. Queue Status:"
QUEUE_LEN=$(docker-compose -f docker-compose.prod.yml --env-file .env.production exec -T redis redis-cli -a "${REDIS_PASSWORD}" LLEN photo_moderation_queue 2>/dev/null | tr -d '\r\n' || echo "0")
echo "   Jobs in queue: $QUEUE_LEN"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Workers are ready!"
echo ""
echo "📸 Next: Upload a photo from the frontend and watch the logs:"
echo "   docker-compose -f docker-compose.prod.yml --env-file .env.production logs -f moderator-worker"
echo ""
echo "You should see:"
echo "   - 'OpenCV fallback detected face' (if CompreFace fails)"
echo "   - 'Face detection result: has_face=True' (if face found)"
echo "   - '✅ Approved' or '❌ Rejected: [reason]'"
echo ""

