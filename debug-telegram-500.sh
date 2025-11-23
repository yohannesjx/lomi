#!/bin/bash

# Debug script to capture real Telegram 500 errors
# Run this on the server and watch logs in real-time

echo "🔍 Debugging Telegram 500 Errors"
echo "=================================="
echo ""
echo "This script will monitor backend logs for errors when you try to login from Telegram"
echo ""
echo "Instructions:"
echo "1. Keep this script running"
echo "2. Try to login from your Telegram app"
echo "3. Watch the logs below for the exact error"
echo ""
echo "Press Ctrl+C to stop"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Follow backend logs and filter for errors
docker-compose -f docker-compose.prod.yml logs -f backend 2>&1 | grep -E "(❌|ERROR|error|500|Failed|panic|fatal)" --color=always

