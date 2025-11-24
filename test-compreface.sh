#!/bin/bash

# Test CompreFace API connectivity and face detection

set -e

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

COMPREFACE_URL="${COMPREFACE_URL:-http://localhost:8000}"

echo "🔍 Testing CompreFace API"
echo "========================="
echo "CompreFace URL: $COMPREFACE_URL"
echo ""

# Test 1: Health check
echo "1. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$COMPREFACE_URL/api/v1/health" || echo "FAILED")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
HEALTH_BODY=$(echo "$HEALTH_RESPONSE" | grep -v "HTTP_CODE:")

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Health check passed"
    echo "Response: $HEALTH_BODY"
else
    echo "❌ Health check failed (HTTP $HTTP_CODE)"
    echo "Response: $HEALTH_BODY"
    echo ""
    echo "⚠️  CompreFace might not be running or accessible"
    exit 1
fi

echo ""

# Test 2: List services (if admin API available)
echo "2. Testing services endpoint..."
SERVICES_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$COMPREFACE_URL/api/v1/services" || echo "FAILED")
SERVICES_HTTP=$(echo "$SERVICES_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
SERVICES_BODY=$(echo "$SERVICES_RESPONSE" | grep -v "HTTP_CODE:")

if [ "$SERVICES_HTTP" = "200" ]; then
    echo "✅ Services endpoint accessible"
    echo "Response: $SERVICES_BODY"
else
    echo "⚠️  Services endpoint returned HTTP $SERVICES_HTTP (might need API key)"
fi

echo ""

# Test 3: Test detection with a sample image
echo "3. Testing face detection endpoint..."
echo "   (This requires a test image - creating a simple test...)"

# Create a simple test image (1x1 pixel - won't have a face, but tests the endpoint)
TEST_IMAGE_B64="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
echo "$TEST_IMAGE_B64" | base64 -d > /tmp/test_face.jpg

DETECT_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST \
    -F "file=@/tmp/test_face.jpg" \
    "$COMPREFACE_URL/api/v1/detection/detect" || echo "FAILED")

DETECT_HTTP=$(echo "$DETECT_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
DETECT_BODY=$(echo "$DETECT_RESPONSE" | grep -v "HTTP_CODE:")

echo "Detection endpoint HTTP code: $DETECT_HTTP"
echo "Response: $DETECT_BODY"

if [ "$DETECT_HTTP" = "200" ]; then
    echo "✅ Detection endpoint is working!"
elif [ "$DETECT_HTTP" = "401" ] || [ "$DETECT_HTTP" = "403" ]; then
    echo "⚠️  Detection endpoint requires authentication (API key)"
    echo "   You may need to create a detection service and get an API key"
elif [ "$DETECT_HTTP" = "404" ]; then
    echo "⚠️  Detection endpoint not found - might need different endpoint"
    echo "   Try: /api/v1/recognition/recognize or check CompreFace docs"
else
    echo "❌ Detection endpoint failed (HTTP $DETECT_HTTP)"
fi

rm -f /tmp/test_face.jpg

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Summary:"
echo "   - Health: $([ "$HTTP_CODE" = "200" ] && echo "✅ OK" || echo "❌ FAILED")"
echo "   - Detection: $([ "$DETECT_HTTP" = "200" ] && echo "✅ OK" || echo "⚠️  Check logs above")"
echo ""
echo "💡 Next steps:"
echo "   1. Check worker logs: docker-compose logs moderator-worker"
echo "   2. Check CompreFace logs: docker-compose logs compreface"
echo "   3. If detection needs API key, create service via CompreFace admin UI"
echo ""

