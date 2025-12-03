#!/bin/bash

# TikTok API Endpoint Test Script
# Tests all 6 critical endpoints

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="$BASE_URL/api"

echo "🧪 Testing TikTok API Endpoints"
echo "================================"
echo "Base URL: $API_URL"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0

# Helper function to test endpoint
test_endpoint() {
    local name=$1
    local endpoint=$2
    local data=$3
    
    echo -e "${YELLOW}Testing: $name${NC}"
    echo "Endpoint: POST $endpoint"
    
    response=$(curl -s -X POST "$API_URL$endpoint" \
        -H "Content-Type: application/json" \
        -d "$data")
    
    code=$(echo "$response" | jq -r '.code // empty')
    
    if [ "$code" = "200" ]; then
        echo -e "${GREEN}✅ PASSED${NC}"
        echo "Response: $(echo "$response" | jq -c '.msg' | head -c 100)..."
        ((PASSED++))
    else
        echo -e "${RED}❌ FAILED${NC}"
        echo "Response: $response"
        ((FAILED++))
    fi
    echo ""
}

# 1. Test registerUser
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Register User (Social Login)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Register User" "/registerUser" '{
  "username": "testuser",
  "first_name": "Test",
  "last_name": "User",
  "email": "test@example.com",
  "phone": "+251912345678",
  "social_id": "google_test123",
  "social": "google",
  "device_token": "fcm_test_token"
}'

# 2. Test showUserDetail
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  Show User Detail"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Show User Detail" "/showUserDetail" '{
  "user_id": "1"
}'

# 3. Test showRelatedVideos
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  Show Related Videos (Feed)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Show Related Videos" "/showRelatedVideos" '{
  "user_id": "1",
  "device_id": "device_test123",
  "starting_point": 0,
  "lat": 9.0320,
  "long": 38.7469
}'

# 4. Test liveStream
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  Start Live Stream"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Start Live Stream" "/liveStream" '{
  "user_id": "1",
  "started_at": "2024-12-03 16:00:00"
}'

# 5. Test showCoinWorth
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5️⃣  Show Coin Packages"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Show Coin Worth" "/showCoinWorth" '{}'

# 6. Test showGifts
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6️⃣  Show Gifts Catalog"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_endpoint "Show Gifts" "/showGifts" '{}'

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
