#!/bin/bash

# Test Photo Moderation System on VPS
# Tests the complete flow via API endpoints only (no Docker access needed)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing Photo Moderation on VPS${NC}"
echo "=============================================="
echo ""

# Configuration
API_BASE="${API_BASE:-https://lomi.social}"
NUM_PHOTOS="${1:-5}"  # Default to 5 photos, can override: ./test-moderation-vps.sh 9
NUM_PHOTOS=$((NUM_PHOTOS > 9 ? 9 : NUM_PHOTOS))  # Max 9 photos per batch

echo -e "${YELLOW}📸 Testing with $NUM_PHOTOS photos${NC}"
echo -e "${YELLOW}🌐 API Base: $API_BASE${NC}"
echo ""

# Step 1: Authenticate
echo -e "${BLUE}Step 1: Authentication${NC}"
echo "─────────────────────────────────"
read -p "Enter JWT token (or press Enter to use initData): " TOKEN

if [ -z "$TOKEN" ]; then
    read -p "Enter Telegram initData: " INIT_DATA
    if [ -z "$INIT_DATA" ]; then
        echo -e "${RED}❌ Need either JWT token or initData${NC}"
        exit 1
    fi
    
    echo "Authenticating with initData..."
    AUTH_RESPONSE=$(curl -s -X POST "$API_BASE/api/v1/auth/telegram" \
        -H "Authorization: tma $INIT_DATA")
    
    # Extract access_token (works without jq)
    TOKEN=$(echo "$AUTH_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
        echo -e "${RED}❌ Authentication failed${NC}"
        echo "Response: $AUTH_RESPONSE"
        exit 1
    fi
    echo -e "${GREEN}✅ Authenticated${NC}"
fi

echo "Token: ${TOKEN:0:30}..."
echo ""

# Step 2: Upload photos to R2 and collect file keys
echo -e "${BLUE}Step 2: Uploading $NUM_PHOTOS photos to R2${NC}"
echo "─────────────────────────────────"

PHOTOS_JSON="["
TEMP_DIR="/tmp/vps_test_photos"
mkdir -p "$TEMP_DIR"

# Create a simple test image (1x1 pixel PNG)
TEST_IMAGE_B64="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
echo "$TEST_IMAGE_B64" | base64 -d > "$TEMP_DIR/test_base.jpg"

for i in $(seq 1 $NUM_PHOTOS); do
    echo "  📤 Uploading photo $i/$NUM_PHOTOS..."
    
    # Get presigned upload URL
    UPLOAD_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/users/media/upload-url?media_type=photo" \
        -H "Authorization: Bearer $TOKEN")
    
    # Check for errors in response
    if echo "$UPLOAD_RESPONSE" | grep -q '"error"'; then
        echo -e "    ${RED}❌ Error getting upload URL: $UPLOAD_RESPONSE${NC}"
        continue
    fi
    
    # Extract upload_url and file_key (works without jq)
    # Handle escaped characters in JSON (like \u0026 for &)
    UPLOAD_URL=$(echo "$UPLOAD_RESPONSE" | grep -o '"upload_url":"[^"]*"' | cut -d'"' -f4 | sed 's/\\u0026/\&/g' | sed 's/\\\//\//g')
    FILE_KEY=$(echo "$UPLOAD_RESPONSE" | grep -o '"file_key":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$UPLOAD_URL" ] || [ "$UPLOAD_URL" = "null" ]; then
        echo -e "    ${RED}❌ Failed to get upload URL for photo $i${NC}"
        echo "    Response: $UPLOAD_RESPONSE"
        continue
    fi
    
    # Upload to R2
    UPLOAD_RESPONSE_FILE="/tmp/upload_response_$i.txt"
    UPLOAD_STATUS=$(curl -s -o "$UPLOAD_RESPONSE_FILE" -w "%{http_code}" -X PUT "$UPLOAD_URL" \
        -H "Content-Type: image/jpeg" \
        --data-binary @"$TEMP_DIR/test_base.jpg")
    
    if [ "$UPLOAD_STATUS" = "200" ] || [ "$UPLOAD_STATUS" = "204" ]; then
        echo -e "    ${GREEN}✅ Photo $i uploaded (key: ${FILE_KEY:0:40}...)${NC}"
        
        # Add to photos array
        if [ "$i" -gt 1 ]; then
            PHOTOS_JSON+=","
        fi
        PHOTOS_JSON+="{\"file_key\":\"$FILE_KEY\",\"media_type\":\"photo\"}"
    else
        UPLOAD_ERROR=$(cat "$UPLOAD_RESPONSE_FILE" 2>/dev/null || echo "")
        echo -e "    ${RED}❌ Upload failed with status: $UPLOAD_STATUS${NC}"
        if [ -n "$UPLOAD_ERROR" ]; then
            echo "    Error: $UPLOAD_ERROR"
        fi
    fi
    rm -f "$UPLOAD_RESPONSE_FILE"
    
    # Small delay to avoid rate limiting
    sleep 0.5
done

PHOTOS_JSON+="]"

# Wrap in proper request format: {"photos": [...]}
if [ "$PHOTOS_JSON" = "[]" ]; then
    echo -e "${RED}❌ No photos were uploaded successfully${NC}"
    exit 1
fi

UPLOAD_COMPLETE_BODY="{\"photos\":$PHOTOS_JSON}"
echo ""

# Step 3: Call upload-complete endpoint
echo -e "${BLUE}Step 3: Calling upload-complete endpoint${NC}"
echo "─────────────────────────────────"

UPLOAD_COMPLETE_RESPONSE=$(curl -s -X POST "$API_BASE/api/v1/users/media/upload-complete" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$UPLOAD_COMPLETE_BODY")

echo "Response:"
echo "$UPLOAD_COMPLETE_RESPONSE"
echo ""

# Extract batch_id (works without jq)
BATCH_ID=$(echo "$UPLOAD_COMPLETE_RESPONSE" | grep -o '"batch_id":"[^"]*"' | cut -d'"' -f4)
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_BASE/api/v1/users/media/upload-complete" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$UPLOAD_COMPLETE_BODY")

if [ "$HTTP_STATUS" = "200" ] && [ -n "$BATCH_ID" ] && [ "$BATCH_ID" != "null" ]; then
    echo -e "${GREEN}✅ Upload-complete successful!${NC}"
    echo "   Batch ID: $BATCH_ID"
else
    echo -e "${RED}❌ Upload-complete failed${NC}"
    echo "   HTTP Status: $HTTP_STATUS"
    exit 1
fi
echo ""

# Step 4: Monitor the queue using admin endpoints
echo -e "${BLUE}Step 4: Monitoring Queue & Moderation Status${NC}"
echo "─────────────────────────────────"
echo ""
echo -e "${YELLOW}📊 Initial Queue Status:${NC}"

# Check queue stats via admin endpoint
QUEUE_STATS=$(curl -s -X GET "$API_BASE/api/v1/admin/queue-stats" \
    -H "Authorization: Bearer $TOKEN")

if [ -n "$QUEUE_STATS" ] && ! echo "$QUEUE_STATS" | grep -q '"error"'; then
    QUEUE_LEN=$(echo "$QUEUE_STATS" | grep -o '"queue_length":[0-9]*' | cut -d':' -f2 || echo "0")
    PENDING_COUNT=$(echo "$QUEUE_STATS" | grep -o '"pending_media":[0-9]*' | cut -d':' -f2 || echo "0")
    echo "   Queue length: $QUEUE_LEN"
    echo "   Pending media: $PENDING_COUNT"
else
    echo "   ⚠️  Could not fetch queue stats (admin endpoint may require admin role)"
    PENDING_COUNT="unknown"
fi
echo ""

# Wait and monitor using dashboard endpoint
echo -e "${YELLOW}⏳ Waiting for moderation to complete (max 60 seconds)...${NC}"
MAX_WAIT=60
ELAPSED=0
CHECK_INTERVAL=3

while [ $ELAPSED -lt $MAX_WAIT ]; do
    sleep $CHECK_INTERVAL
    ELAPSED=$((ELAPSED + CHECK_INTERVAL))
    
    # Check moderation dashboard for this batch
    DASHBOARD_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/admin/moderation/dashboard?status=all&limit=100" \
        -H "Authorization: Bearer $TOKEN")
    
    # Count pending photos in this batch
    if [ -n "$DASHBOARD_RESPONSE" ] && ! echo "$DASHBOARD_RESPONSE" | grep -q '"error"'; then
        # Try to extract pending count for this batch
        PENDING_IN_BATCH=$(echo "$DASHBOARD_RESPONSE" | grep -o "\"batch_id\":\"$BATCH_ID\"" | wc -l | tr -d ' ' || echo "0")
        # This is a rough check - we'll use a different approach
    fi
    
    # Alternative: Check queue stats again
    QUEUE_STATS=$(curl -s -X GET "$API_BASE/api/v1/admin/queue-stats" \
        -H "Authorization: Bearer $TOKEN" 2>/dev/null || echo "")
    
    if [ -n "$QUEUE_STATS" ] && ! echo "$QUEUE_STATS" | grep -q '"error"'; then
        PENDING_COUNT=$(echo "$QUEUE_STATS" | grep -o '"pending_media":[0-9]*' | cut -d':' -f2 || echo "?")
        if [ "$PENDING_COUNT" = "0" ] || [ "$PENDING_COUNT" = "?" ]; then
            echo -e "${GREEN}✅ Queue appears empty - moderation likely complete!${NC}"
            break
        fi
    fi
    
    # Show progress
    echo "   Progress: Checking... (${ELAPSED}s elapsed)"
done

if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo -e "${YELLOW}⚠️  Timeout reached. Checking final status...${NC}"
fi
echo ""

# Step 5: Check results via dashboard
echo -e "${BLUE}Step 5: Checking Moderation Results${NC}"
echo "─────────────────────────────────"
echo ""

# Get moderation dashboard
echo "📊 Moderation Dashboard Results:"
echo ""

DASHBOARD_ALL=$(curl -s -X GET "$API_BASE/api/v1/admin/moderation/dashboard?status=all&limit=100" \
    -H "Authorization: Bearer $TOKEN")

if [ -n "$DASHBOARD_ALL" ] && ! echo "$DASHBOARD_ALL" | grep -q '"error"'; then
    echo "$DASHBOARD_ALL" | python3 -m json.tool 2>/dev/null || echo "$DASHBOARD_ALL"
else
    echo -e "${YELLOW}⚠️  Could not fetch dashboard (may require admin role)${NC}"
    echo "Response: $DASHBOARD_ALL"
fi

echo ""

# Get queue stats summary
echo "📈 Queue Statistics:"
QUEUE_STATS=$(curl -s -X GET "$API_BASE/api/v1/admin/queue-stats" \
    -H "Authorization: Bearer $TOKEN")

if [ -n "$QUEUE_STATS" ] && ! echo "$QUEUE_STATS" | grep -q '"error"'; then
    echo "$QUEUE_STATS" | python3 -m json.tool 2>/dev/null || echo "$QUEUE_STATS"
else
    echo -e "${YELLOW}⚠️  Could not fetch queue stats${NC}"
fi

echo ""

# Cleanup
rm -rf "$TEMP_DIR"

# Final summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ VPS Moderation Test Complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Batch ID: $BATCH_ID"
echo ""
echo "Next steps:"
echo "  1. Check moderation dashboard: curl -X GET \"$API_BASE/api/v1/admin/moderation/dashboard?status=all\" -H \"Authorization: Bearer \$TOKEN\""
echo "  2. Check queue stats: curl -X GET \"$API_BASE/api/v1/admin/queue-stats\" -H \"Authorization: Bearer \$TOKEN\""
echo "  3. SSH to VPS and check worker logs: docker-compose -f docker-compose.prod.yml logs moderator-worker"
echo ""

