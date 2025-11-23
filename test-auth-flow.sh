#!/bin/bash

# Comprehensive Authentication Flow Test Script
# Run this ON THE SERVER to test all login scenarios

set -e

echo "🧪 Testing Telegram Authentication Flow"
echo "========================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to test endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local headers=$4
    local expected_status=$5
    local description=$6
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Test: $name"
    echo "Description: $description"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ -z "$headers" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" 2>&1)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" -H "$headers" 2>&1)
    fi
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n 1)
    # Extract body (all but last line)
    body=$(echo "$response" | head -n -1)
    
    echo "URL: $url"
    echo "Method: $method"
    [ -n "$headers" ] && echo "Headers: $headers"
    echo "Status Code: $status_code"
    echo "Response Body: $body"
    echo ""
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASSED${NC} - Expected status $expected_status, got $status_code"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}❌ FAILED${NC} - Expected status $expected_status, got $status_code"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        if [ "$status_code" = "500" ]; then
            echo -e "${RED}⚠️  500 Internal Server Error detected!${NC}"
        fi
    fi
    echo ""
}

# Test 1: Health Check
echo "📋 Test 1: Backend Health Check"
test_endpoint \
    "Health Check" \
    "GET" \
    "http://localhost:8080/api/v1/health" \
    "" \
    "200" \
    "Check if backend is running and healthy"

# Test 2: Health Check via Caddy
echo "📋 Test 2: Backend Health Check (via Caddy)"
test_endpoint \
    "Health Check via Caddy" \
    "GET" \
    "http://localhost/api/v1/health" \
    "" \
    "200" \
    "Check if backend is accessible through Caddy proxy"

# Test 3: Missing Authorization Header
echo "📋 Test 3: Missing Authorization Header"
test_endpoint \
    "Missing Auth Header" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "" \
    "401" \
    "Should return 401 when Authorization header is missing"

# Test 4: Invalid Authorization Format
echo "📋 Test 4: Invalid Authorization Format"
test_endpoint \
    "Invalid Auth Format" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "Authorization: invalid" \
    "401" \
    "Should return 401 for invalid Authorization format"

# Test 5: Invalid initData (just 'test')
echo "📋 Test 5: Invalid initData (test)"
test_endpoint \
    "Invalid initData" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "Authorization: tma test" \
    "401" \
    "Should return 401 for invalid initData (not 500!)"

# Test 6: Empty initData
echo "📋 Test 6: Empty initData"
test_endpoint \
    "Empty initData" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "Authorization: tma " \
    "401" \
    "Should return 401 for empty initData"

# Test 7: Malformed initData (missing parts)
echo "📋 Test 7: Malformed initData"
test_endpoint \
    "Malformed initData" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "Authorization: tma user=%7B%22id%22%3A123%7D" \
    "401" \
    "Should return 401 for malformed initData (missing hash/sign)"

# Test 8: Check backend logs for errors
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Test 8: Checking Backend Logs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Recent backend logs (last 30 lines):"
echo ""
docker-compose -f docker-compose.prod.yml logs --tail=30 backend 2>/dev/null || echo "Could not get logs"
echo ""

# Test 9: Check database connectivity
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Test 9: Database Connectivity"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Checking if database container is running..."
if docker ps | grep -q lomi_postgres; then
    echo -e "${GREEN}✅ PostgreSQL container is running${NC}"
    echo "Testing database connection..."
    docker-compose -f docker-compose.prod.yml exec -T postgres psql -U lomi -d lomi_db -c "SELECT 1;" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Database connection successful${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}❌ Database connection failed${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
else
    echo -e "${RED}❌ PostgreSQL container is not running${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
echo ""

# Test 10: Check Redis connectivity
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Test 10: Redis Connectivity"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if docker ps | grep -q lomi_redis; then
    echo -e "${GREEN}✅ Redis container is running${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}❌ Redis container is not running${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
echo ""

# Test 11: Check backend container status
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Test 11: Backend Container Status"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if docker ps | grep -q lomi_backend; then
    echo -e "${GREEN}✅ Backend container is running${NC}"
    echo "Container health status:"
    docker inspect --format='{{.State.Health.Status}}' lomi_backend 2>/dev/null || echo "No health check configured"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}❌ Backend container is not running${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
echo ""

# Test 12: Test with a more realistic (but still invalid) initData format
echo "📋 Test 12: Realistic but Invalid initData Format"
# This simulates what real initData might look like but without valid signature
INVALID_INITDATA="user=%7B%22id%22%3A123456789%2C%22first_name%22%3A%22Test%22%2C%22last_name%22%3A%22User%22%7D&auth_date=1234567890&hash=invalid_hash"
test_endpoint \
    "Realistic Invalid initData" \
    "POST" \
    "http://localhost/api/v1/auth/telegram" \
    "Authorization: tma $INVALID_INITDATA" \
    "401" \
    "Should return 401 for invalid signature (not 500!)"

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Check the output above for details.${NC}"
    echo ""
    echo "💡 Next steps:"
    echo "   1. Check backend logs: docker-compose -f docker-compose.prod.yml logs backend"
    echo "   2. Check if all containers are running: docker-compose -f docker-compose.prod.yml ps"
    echo "   3. Verify .env.production has all required variables"
    exit 1
fi

