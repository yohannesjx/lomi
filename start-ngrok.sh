#!/bin/bash

# Start ngrok tunnel for Telegram Mini App development
# This creates an HTTPS tunnel to your local backend

echo "🚀 Starting ngrok tunnel for backend..."
echo ""
echo "Make sure your backend is running on port 8080"
echo ""

# Check if backend is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "⚠️  Warning: Backend doesn't seem to be running on port 8080"
    echo "   Start it with: docker-compose up backend"
    echo "   Or: cd backend && go run cmd/api/main.go"
    echo ""
fi

# Start ngrok
echo "Starting ngrok..."
ngrok http 8080

