#!/bin/bash

# Complete Photo Upload Test Script
# Tests the entire flow: Get URL → Upload to R2 → Create Media Record

set -e

echo "🧪 Testing Complete Photo Upload Flow"
echo "======================================"
echo ""

# Load environment variables
if [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
fi

# Step 1: Authenticate and get token
echo "Step 1: Authenticating..."
echo "⚠️  You need to provide a valid Telegram initData"
echo "   Or use an existing JWT token"
echo ""
read -p "Enter JWT token (or press Enter to skip and use initData): " TOKEN

if [ -z "$TOKEN" ]; then
    read -p "Enter Telegram initData: " INIT_DATA
    if [ -z "$INIT_DATA" ]; then
        echo "❌ Need either JWT token or initData"
        exit 1
    fi
    
    echo "Authenticating with initData..."
    AUTH_RESPONSE=$(curl -s -X POST http://localhost/api/v1/auth/telegram \
        -H "Authorization: tma $INIT_DATA")
    
    TOKEN=$(echo $AUTH_RESPONSE | jq -r '.access_token' 2>/dev/null || echo "")
    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
        echo "❌ Authentication failed"
        echo "Response: $AUTH_RESPONSE"
        exit 1
    fi
    echo "✅ Authenticated"
fi

echo "Token: ${TOKEN:0:30}..."
echo ""

# Step 2: Get upload URL
echo "Step 2: Getting presigned upload URL..."
UPLOAD_RESPONSE=$(curl -s -X GET "http://localhost/api/v1/users/media/upload-url?media_type=photo" \
    -H "Authorization: Bearer $TOKEN")

echo "Response: $UPLOAD_RESPONSE" | jq '.' 2>/dev/null || echo "$UPLOAD_RESPONSE"
echo ""

UPLOAD_URL=$(echo $UPLOAD_RESPONSE | jq -r '.upload_url' 2>/dev/null || echo "")
FILE_KEY=$(echo $UPLOAD_RESPONSE | jq -r '.file_key' 2>/dev/null || echo "")

if [ -z "$UPLOAD_URL" ] || [ "$UPLOAD_URL" = "null" ]; then
    echo "❌ Failed to get upload URL"
    exit 1
fi

echo "✅ Got upload URL"
echo "   File Key: $FILE_KEY"
echo "   URL: ${UPLOAD_URL:0:80}..."
echo ""

# Step 3: Create test image
echo "Step 3: Creating test image..."
# Create a 1x1 pixel PNG image (base64 encoded)
TEST_IMAGE_B64="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
echo "$TEST_IMAGE_B64" | base64 -d > /tmp/test_photo.jpg
echo "✅ Created test image ($(wc -c < /tmp/test_photo.jpg) bytes)"
echo ""

# Step 4: Upload to R2
echo "Step 4: Uploading to R2..."
echo "   URL: ${UPLOAD_URL:0:100}..."
UPLOAD_STATUS=$(curl -s -o /tmp/upload_response.txt -w "%{http_code}" -X PUT "$UPLOAD_URL" \
    -H "Content-Type: image/jpeg" \
    --data-binary @/tmp/test_photo.jpg)

UPLOAD_RESPONSE_BODY=$(cat /tmp/upload_response.txt 2>/dev/null || echo "")

echo "   Status Code: $UPLOAD_STATUS"
if [ -n "$UPLOAD_RESPONSE_BODY" ]; then
    echo "   Response Body: $UPLOAD_RESPONSE_BODY"
fi

if [ "$UPLOAD_STATUS" = "200" ] || [ "$UPLOAD_STATUS" = "204" ]; then
    echo "✅ Upload to R2 successful"
else
    echo "❌ Upload to R2 failed with status: $UPLOAD_STATUS"
    echo ""
    echo "Possible issues:"
    echo "  - Presigned URL expired"
    echo "  - CORS issue"
    echo "  - R2 bucket doesn't exist"
    echo "  - Wrong credentials"
    echo "  - Network issue"
    exit 1
fi
echo ""

# Step 5: Create media record
echo "Step 5: Creating media record..."
MEDIA_RESPONSE=$(curl -s -X POST "http://localhost/api/v1/users/media" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"media_type\": \"photo\",
        \"file_key\": \"$FILE_KEY\",
        \"display_order\": 0
    }")

echo "Response: $MEDIA_RESPONSE" | jq '.' 2>/dev/null || echo "$MEDIA_RESPONSE"
echo ""

MEDIA_ID=$(echo $MEDIA_RESPONSE | jq -r '.id' 2>/dev/null || echo "")
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "http://localhost/api/v1/users/media" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"media_type\": \"photo\",
        \"file_key\": \"$FILE_KEY\",
        \"display_order\": 0
    }")

if [ "$HTTP_STATUS" = "201" ] || [ -n "$MEDIA_ID" ] && [ "$MEDIA_ID" != "null" ]; then
    echo "✅ Media record created successfully"
    echo "   Media ID: $MEDIA_ID"
else
    echo "❌ Failed to create media record"
    echo "   HTTP Status: $HTTP_STATUS"
    echo "   Response: $MEDIA_RESPONSE"
    exit 1
fi

# Cleanup
rm -f /tmp/test_photo.jpg /tmp/upload_response.txt

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ All tests passed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "The photo upload flow is working correctly."

